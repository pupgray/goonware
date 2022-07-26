package configurator

import (
	types "goonware/types"

	"os"
	"fmt"

	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
)

func AnnoyancesTab(c *types.Config) []g.Widget {
	return []g.Widget{
		g.Checkbox("On", &c.Annoyances),
		ConditionOrNothing(c.Annoyances, g.Layout{
			g.Row(LabelSliderTooltip("Seconds per tick", &c.TimerDelay, 1, 60, 250,
				"Number of seconds between attempts to spawn annoyance", FormatSecondSlider)),
			StandardSeparation(),
				
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
									Format(FormatSecondSlider(c.PopupDenialChance)).Size(150),
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
									Format(FormatSecondSlider(c.PopupTimeoutDelay)).Size(150),
								g.Tooltip("How long popups will remain on the screen" +
									" until disppearing"),
							}),
						),
				).Size(g.Auto, 120),
			}),
			StandardSeparation(),

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
			StandardSeparation(),
			
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
			StandardSeparation(),

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
			StandardSeparation(),
		}),
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