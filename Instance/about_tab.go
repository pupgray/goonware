package Goonware

import (
	g "github.com/AllenDang/giu"
)

func AboutTab() []g.Widget {
	return []g.Widget{
		g.Row(
			g.Label("Goonware by zoomasochist#2530"),
		),
	}
}
