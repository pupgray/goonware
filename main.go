package main

import (
	configurator "goonware/configurator"
	daemon "goonware/daemon"
	types "goonware/types"

	"flag"
)

var enterDaemon = flag.Bool("daemon", false, "Spawn Goonware, not the configurator")

func main() {
	flag.Parse()

	if *enterDaemon {
		c, err := configurator.NewOrLoadConfig()
		if err != nil {
			panic(err)
		}
		pkg, err := types.LoadEdgewarePackage(c.LoadedPackage, c.WorkingDirectory)
		if err != nil {
			panic(err)
		}

		daemon.Tick(c, pkg)
	} else {
		panic(configurator.ConfiguratorUI())
	}
}
