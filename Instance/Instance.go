package Goonware

import (
	"os"
)

type Instance struct {
	WorkingDirectoryPath        string
	PackageExtractDirectoryPath string
	ConfigFilePath              string
	configurator                Configurator
}

func New() Instance {
	userConfigDirectory, errUserConfigDirectory := os.UserConfigDir()

	if errUserConfigDirectory != nil {
		panic(errUserConfigDirectory)
	}

	workingDirectoryPath := userConfigDirectory + "/goonware"

	instance := Instance{
		WorkingDirectoryPath:        workingDirectoryPath,
		PackageExtractDirectoryPath: workingDirectoryPath + "/package",
		ConfigFilePath:              workingDirectoryPath + "/goonware.json",
	}

	instance.configurator = NewConfigurator(&instance)

	return instance
}

/*func (i Instance) LoadConfiguration() types.Config {
	config, err := i.configurator.NewOrLoadConfig()
	if err != nil {
		panic(err)
	}
	return config
}*/

func (i Instance) LaunchConfigurator() {
	i.configurator.Launch()
}
