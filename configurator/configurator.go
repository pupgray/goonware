package configurator

import (
	g "github.com/AllenDang/giu"
)

func ConfiguratorUI() {
	wnd := g.NewMasterWindow("Goonware", 700, 700, g.MasterWindowFlagsNotResizable)
	wnd.Run(func() {
		g.SingleWindow().Layout(
			g.TabBar().TabItems(
				g.TabItem("General").Layout(GeneralTab()...),
				g.TabItem("About").Layout(AboutTab()...),
			),
		)
	})
}