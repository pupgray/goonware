package main

import (
	"flag"
	Goonware "goonware/Instance"
)

var enterDaemon = flag.Bool("daemon", false, "Spawn Goonware, not the configurator")

var instance = Goonware.New()

func main() {
	flag.Parse()

	/*if *enterDaemon {
		c, err := configurator.NewOrLoadConfig()
		if err != nil {
			panic(err)
		}
		pkg, err := types.LoadEdgewarePackage(c.LoadedPackage, packageExtractDirectory)
		if err != nil {
			panic(err)
		}

		daemon.Tick(c, pkg)
	} else {*/
	instance.LaunchConfigurator()
	/*}*/
}
