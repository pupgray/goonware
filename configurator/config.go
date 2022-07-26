package configurator

import (
	types "goonware/types"

	"os"
	"errors"
	"encoding/json"
)

var workingDirectory   = Expect(os.UserConfigDir()) + "/goonware"
var packageExtractDirectory = workingDirectory + "/package"
var configFileLocation = workingDirectory + "/goonware.json"

func Expect(s string, err error) string {
	if err != nil {
		panic(err)
	}

	return s
}

func NewOrLoadConfig() types.Config {
	if _, err := os.Stat(configFileLocation); errors.Is(err, os.ErrNotExist) {
		_ = os.MkdirAll(packageExtractDirectory, os.ModePerm)

		return types.Config{
			WorkingDirectory: workingDirectory,
			// Mode
			Mode: 0,
			HibernateMinWait: 120,
			HibernateMaxWait: 3600,
			HibernateActivityLength: 20,
			// Annoyances
			Annoyances: true,
			TimerDelay: 10,
			AnnoyancePopups: true,

			PopupChance: 50,
			PopupOpacity: 85,
			PopupDenialMode: false,
			PopupDenialChance: 50,
			PopupMitosis: true,
			PopupMitosisStrength: 4,
			PopupTimeout: false,
			PopupTimeoutDelay: 30,

			AnnoyanceVideos: true,
			VideoChance: 25,
			VideoVolume: 25,

			AnnoyancePrompts: true,
			PromptChance: 25,
			MaxMistakesToggle: true,
			MaxMistakes: 1,

			AnnoyanceAudio: true,
			AudioChance: 25,
			AudioVolume: 25,
			// Package
			LoadedPackage: "",
			LoadedPackageData: nil,
			// Other
			StartOnBoot: false,
			RunOnExit: false,
		}
	}

	return LoadConfig()
}

func SaveConfig(c types.Config) {
	structBytes, err := json.Marshal(c)
	if err != nil {
		// Todo: Handle better
		panic(err)
	}

	err = os.WriteFile(configFileLocation, structBytes, 0644)
	if err != nil {
		// Todo: Handle better
		panic(err)
	}
}

func LoadConfig() types.Config {
	structBytes, err := os.ReadFile(configFileLocation)
	if err != nil {
		// Todo: Handle better
		panic(err)
	}

	var c types.Config
	err = json.Unmarshal(structBytes, &c)
	if err != nil {
		// Todo: Handle better
		panic(err)
	}

	return c
}