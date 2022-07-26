package configurator

import (
	g "github.com/AllenDang/giu"
)

func ConfiguratorUI() error {
	c, err := NewOrLoadConfig()
	if err != nil {
		return err
	}

	wnd := g.NewMasterWindow("Goonware", 700, 700, g.MasterWindowFlagsNotResizable)
	wnd.Run(func() {
		g.SingleWindow().Layout(
			g.TabBar().TabItems(
				g.TabItem("General").Layout(GeneralTab(&c)...),
				g.TabItem("About").Layout(AboutTab()...),
			),
		)
	})

	return nil
}