package Goonware

import (
	types "goonware/types"

	"encoding/json"
	"errors"
	"os"

	g "github.com/AllenDang/giu"
)

type Configurator struct {
	parentInstance *Instance
}

func NewConfigurator(parentInstance *Instance) Configurator {
	return Configurator{
		parentInstance: parentInstance,
	}
}

func (config Configurator) Launch() error {
	c, err := config.NewOrLoadConfig()
	if err != nil {
		return err
	}

	wnd := g.NewMasterWindow("Goonware", 700, 700, 0)
	wnd.Run(func() {
		g.SingleWindow().Layout(
			g.TabBar().TabItems(
				g.TabItem("Annoyances").Layout(AnnoyancesTab(&c)...),
				g.TabItem("Drive Filler").Layout(DriveFillterTab(&c)...),
				// Comes last
				g.TabItem("About").Layout(AboutTab()...),
			),

			g.Row(
				g.Button("Load package").OnClick(func() { LoadPackage(&c) }),
				g.Condition(len(c.LoadedPackage) != 0,
					g.Layout{g.Label("(Using package " + c.LoadedPackage + ")")},
					g.Layout{g.Label("(No package loaded)")},
				),
			),

			g.Row(
				g.Button("Save").OnClick(func() { config.SaveConfig(&c) }),
				g.Button("Save and Exit").OnClick(func() { SaveAndExit(config, &c) }),
				g.Row(
					g.Checkbox("Launch on startup", &c.StartOnBoot),
					g.Tooltip("Silently start Goonware every time your computer starts, executing"+
						" whatever package was running last time."),
					g.Checkbox("Run Goonware on save and exit", &c.RunOnExit),
				),
			),

			g.Row(
				g.Label("Mode"),
				g.RadioButton("Normal", c.Mode == 0).OnChange(func() { c.Mode = 0 }),
				g.Tooltip("As soon as Goonware starts, it will start running payloads."),

				g.RadioButton("Hibernate", c.Mode == 1).OnChange(func() { c.Mode = 1 }),
				g.Tooltip("Goonware will wait a random amount of time (within given limits) before"+
					" spamming payloads, then stop and start waiting again."),

				ConditionOrNothing(c.Mode == 1, g.Layout{
					LabelSliderTooltip("Min. wait", &c.HibernateMinWaitMinutes, 0, 120, 50,
						"The minimum amount of time Goonware will hibernate", FormatMinuteSlider),
					LabelSliderTooltip("Max. wait", &c.HibernateMaxWaitMinutes, 0, 120, 50,
						"The maximum amount of time Goonware will hibernate", FormatMinuteSlider),
					LabelSliderTooltip("Wake for", &c.HibernateActivityLength, 1,
						600, 50, "How long Goonware will spam payloads", FormatSecondSlider),
				}),
			),
		)
	})

	return nil
}

func (config Configurator) NewOrLoadConfig() (types.Config, error) {
	if _, err := os.Stat(config.parentInstance.ConfigFilePath); errors.Is(err, os.ErrNotExist) {
		_ = os.MkdirAll(config.parentInstance.PackageExtractDirectoryPath, os.ModePerm)

		return types.Config{
			WorkingDirectory: config.parentInstance.WorkingDirectoryPath,
			// General
			Mode:                    0,
			HibernateMinWaitMinutes: 120,
			HibernateMaxWaitMinutes: 3600,
			HibernateActivityLength: 20,
			StartOnBoot:             false,
			RunOnExit:               false,
			LoadedPackage:           "",
			LoadedPackageData:       nil,

			// Annoyances
			Annoyances:      true,
			TimerDelay:      10,
			AnnoyancePopups: true,

			PopupChance:          50,
			PopupOpacity:         85,
			PopupDenialMode:      false,
			PopupDenialChance:    50,
			PopupMitosis:         true,
			PopupMitosisStrength: 4,
			PopupTimeout:         false,
			PopupTimeoutDelay:    30,

			AnnoyanceVideos: true,
			VideoChance:     25,
			VideoVolume:     25,

			AnnoyancePrompts:  true,
			PromptChance:      25,
			MaxMistakesToggle: true,
			MaxMistakes:       1,

			AnnoyanceAudio: true,
			AudioChance:    25,
			AudioVolume:    25,

			// Drive Filler
			DriveFiller:                              false,
			DriveFillerDelay:                         10,
			DriveFillerBase:                          Expect(os.UserHomeDir()),
			DriveFillerTags:                          []string{"feral+paws", "feral+rimming"},
			DriveFillerImageSource:                   1,
			DriveFillerImageUseTags:                  true,
			DriveFillerDownloadMinimumScoreToggle:    false,
			DriveFillerDownloadMinimumScoreThreshold: 0,
		}, nil
	}

	return config.LoadConfig()
}

func (config Configurator) SaveConfig(c *types.Config) error {
	structBytes, err := json.Marshal(*c)
	if err != nil {
		return err
	}

	err = os.WriteFile(config.parentInstance.ConfigFilePath, structBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (config Configurator) LoadConfig() (types.Config, error) {
	structBytes, err := os.ReadFile(config.parentInstance.ConfigFilePath)
	if err != nil {
		return types.Config{}, err
	}

	var c types.Config
	err = json.Unmarshal(structBytes, &c)
	if err != nil {
		return types.Config{}, err
	}

	return c, nil
}

func Expect(s string, err error) string {
	if err != nil {
		panic(err)
	}

	return s
}
