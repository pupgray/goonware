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
		c := configurator.NewOrLoadConfig()
		pkg := types.LoadEdgewarePackage(c.LoadedPackage, c.WorkingDirectory)
		daemon.Tick(c, pkg)
	} else {
		configurator.ConfiguratorUI()
	}
}