package configurator

import (
	types "goonware/types"

	"os"
	"fmt"

	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
)

func FormatTimeSlider(value int32) string {
	if value < 60 {
		if value == 1 {
			return fmt.Sprintf("1 second")
		}

		return fmt.Sprintf("%d seconds", value)
	}

	if value % 60 == 0 {
		if value == 60 {
			return fmt.Sprintf("1 minute")
		}

		return fmt.Sprintf("%d minutes", value / 60)
	}

	if value < 120 {
		return fmt.Sprintf("1 minute %d seconds", value % 60)
	}

	return fmt.Sprintf("%d minutes %d seconds", value / 60, value % 60)
}

func FormatPercentSlider(value int32) string {
	// I think having to do this is a consequence of me not exactly using the g.Slider* .Format
	// method correctly. It expects a string that it can pass to C's sprintf along with the value of
	// the slider. With FormatTimeSlider it so happens that the format it gets is along the lines of
	// sprintf("5 seconds");, so all is well. But here, %% would usually be used to escape the first
	// % into a literal %, but because giu passes it into sprintf, it ends up with sprintf("5%");,
	// which is invalid. so %%%% results in sprintf("5%%");, which is okay again.
	return fmt.Sprintf("%d%%%%", value)
}

func ConditionOrNothing(condition bool, layout g.Layout) g.Layout {
	if condition { 
		return layout
	}

	return g.Layout{}
}

func LabelSliderTooltip(label string, value *int32, min, max int32, size float32, tooltip string,
		format func(int32) string) g.Layout {
	return g.Layout{
		g.Label(label),
		g.SliderInt(value, min, max).Format(format(*value)).Size(size),
		g.Tooltip(tooltip),
	}
}

func FrequencySlider(value *int32, tooltip string) g.Layout {
	return g.Layout{
		g.Label("| Frequency"),
		g.SliderInt(value, 1, 100).Format(FormatPercentSlider(*value)).Size(150),
		g.Tooltip(tooltip),
	}
}

func GeneralTab(c *types.Config) []g.Widget {
	largerFont := g.GetDefaultFonts()[0].SetSize(20)

	return []g.Widget{
		g.Layout{
			g.Button("Save").OnClick(func() { SaveConfig(c) }),
			g.Button("Save and exit").OnClick(func() { SaveAndExit(c) }),
		},			
		g.Separator(),

		g.Label("Mode").Font(largerFont),
		g.Row(
			g.RadioButton("Normal", c.Mode == 0).OnChange(func() { c.Mode = 0 }),
			g.Tooltip("As soon as Goonware starts, it will start running payloads."),

			g.RadioButton("Hibernate", c.Mode == 1).OnChange(func() { c.Mode = 1 }),
			g.Tooltip("Goonware will wait a random amount of time (within given limits) before" +
				" spamming payloads, then stop and start waiting again."),
		),
		ConditionOrNothing(c.Mode == 1, g.Layout{
			g.Child().Layout(
				g.Row(
					LabelSliderTooltip("Min. hibernation wait", &c.HibernateMinWait, 0, 3600, 150,
						"The minimum amount of time Goonware will hibernate", FormatTimeSlider),
					LabelSliderTooltip("Max. hibernation wait", &c.HibernateMaxWait, 0, 3600, 150,
						"The maximum amount of time Goonware will hibernate", FormatTimeSlider),
				),
				
				g.Row(
					LabelSliderTooltip("Hibernation activity time", &c.HibernateActivityLength, 1,
						3600, 150, "How long Goonware will spam payloads", FormatTimeSlider),
				),
			).Size(g.Auto, 65),
		}),
		g.Separator(),

		g.Label("Annoyances").Font(largerFont),
		g.Checkbox("On", &c.Annoyances),
		ConditionOrNothing(c.Annoyances, g.Layout{
			g.Row(LabelSliderTooltip("Seconds per tick", &c.TimerDelay, 1, 60, 250,
				"Number of seconds between attempts to spawn annoyance", FormatTimeSlider)),

			g.Row(
				g.Checkbox("Popups", &c.AnnoyancePopups),
				ConditionOrNothing(c.AnnoyancePopups, FrequencySlider(&c.PopupChance,
					"The percent chance a popup will be displayed every tick")),
			),
			ConditionOrNothing(c.AnnoyancePopups, g.Layout{
				g.Child().Layout(
						g.Row(
							LabelSliderTooltip("Popup opacity", &c.PopupOpacity, 1, 100, 150,
								"The opacity of the popup. 100 is opaque, 1 is almost invisible.",
								FormatPercentSlider),
						),

						g.Row(
							g.Checkbox("Denial mode", &c.PopupDenialMode),
							g.Tooltip("Popups may show up blurred, and with captions"),
							ConditionOrNothing(c.PopupDenialMode, g.Layout{
								g.SliderInt(&c.PopupDenialChance, 1, 100).
									Format(FormatTimeSlider(c.PopupDenialChance)).Size(150),
								g.Tooltip("Percent chance of triggering denial mode"),
							}),
						),

						g.Row(
							g.Checkbox("Mitosis", &c.PopupMitosis),
							g.Tooltip("Whenever a popup is closed, more popups will appear" +
								" in its place."),
							ConditionOrNothing(c.PopupMitosis, g.Layout{
								g.SliderInt(&c.PopupMitosisStrength, 2, 10).Size(75),
								g.Tooltip("Number of popups to spawn in the place of a" +
									" closed one."),
							}),
						),

						g.Row(
							g.Checkbox("Popup timeout", &c.PopupTimeout),
							g.Tooltip("Whether popups will timeout (disappear) by themselves." +
								" When disabled, they must be closed manually."),
							ConditionOrNothing(c.PopupTimeout, g.Layout{
								g.SliderInt(&c.PopupTimeoutDelay, 1, 360).
									Format(FormatTimeSlider(c.PopupTimeoutDelay)).Size(150),
								g.Tooltip("How long popups will remain on the screen" +
									" until disppearing"),
							}),
						),
				).Size(g.Auto, 120),
			}),
			g.Dummy(g.Auto, 10),

			g.Row(
				g.Checkbox("Videos",  &c.AnnoyanceVideos),
				ConditionOrNothing(c.AnnoyanceVideos, FrequencySlider(&c.VideoChance,
					"The percent chance a video will be displayed every tick")),
			),
			ConditionOrNothing(c.AnnoyanceVideos, g.Layout{
				g.Child().Layout(g.Row(
					LabelSliderTooltip("Video volume", &c.VideoVolume, 0, 100, 150,
						"Volume of the played video", FormatPercentSlider))).
				Size(g.Auto, 40),
			}),
			g.Dummy(g.Auto, 10),
			
			g.Row(
				g.Checkbox("Prompts", &c.AnnoyancePrompts),
				ConditionOrNothing(c.AnnoyancePrompts, g.Layout{ FrequencySlider(&c.PromptChance,
						"The percent chance a prompt will be displayed every tick") }),
			),
			ConditionOrNothing(c.AnnoyancePrompts, g.Layout{
				g.Child().Layout(
					g.Row(
						g.Checkbox("Max mistakes", &c.MaxMistakesToggle),
						g.Tooltip("With this enabled, there will be punishment for making a" + 
							" specified number of mistakes in the prompt"),
						g.SliderInt(&c.MaxMistakes, 0, 150).Size(150),
					),
				).Size(g.Auto, 40),
			}),
			g.Dummy(g.Auto, 10),

			g.Row(
				g.Checkbox("Audio", &c.AnnoyanceAudio),
				ConditionOrNothing(c.AnnoyanceAudio, FrequencySlider(&c.AudioChance,
					"The percent chance audio will play every tick")),
			),
			ConditionOrNothing(c.AnnoyanceAudio, g.Layout{
				g.Child().Layout(
					g.Row(
						LabelSliderTooltip("Audio volume", &c.AudioVolume, 0, 100, 150,
							"Audio volume", FormatPercentSlider),
					),
				).Size(g.Auto, 40),
			}),
			g.Separator(),
		}),
		g.Separator(),

		g.Label("Package").Font(largerFont),
		g.Condition(len(c.LoadedPackage) != 0,
			g.Layout{ g.Label("Currently loaded package is " + c.LoadedPackage) },
			g.Layout{ g.Label("No package loaded") },
		),
		g.Button("Load package").OnClick(func() { LoadPackage(c) }),
		g.Separator(),

		g.Label("Other").Font(largerFont),
		g.Row(
			g.Checkbox("Launch on startup", &c.StartOnBoot),
			g.Tooltip("Silently start Goonware every time your computer starts, executing" +
				" whatever package was running last time."),
			g.Checkbox("Run Goonware on save and exit", &c.RunOnExit),
		),
	}
}

func LoadPackage(c *types.Config) {
	filename, err := dialog.File().Filter("Edgeware package (.zip)", "zip").Load()

	if err != nil && err != dialog.Cancelled {
		dialog.Message("%s", fmt.Sprintf("Failed to load package; %s", err.Error())).Error()
	} else if err == nil {
		c.LoadedPackage = filename
		//pkg := LoadEdgewarePackage(filename)
	}
}

func SaveAndExit(c *types.Config) {
	fmt.Println("TODO: Shell out")
	SaveConfig(c)
	os.Exit(0)
}