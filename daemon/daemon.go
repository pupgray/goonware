package daemon

import (
	types "goonware/types"

	"time"
	"math/rand"
	g "github.com/AllenDang/giu"
)

var showingPrompt bool
var showingPopup bool

func Tick(c types.Config, pkg types.EdgewarePackage) {
	// We take gooning seriously
	rand.Seed(time.Now().UnixNano())
	inHibernation := (c.Mode == 1)
	tickCounter := 0

	HibernationLoop:
	for {
		hibernationCountdown :=
			rand.Intn(int(c.HibernateMaxWaitMinutes) - int(c.HibernateMinWaitMinutes) + 1) +
				int(c.HibernateMinWaitMinutes)

		for !inHibernation {
			tickCounter++

			if c.Mode == 1 && (tickCounter * int(c.TimerDelay)) > int(c.HibernateActivityLength) {
				inHibernation = true
				continue HibernationLoop
			}

			if c.Annoyances {
				DoAnnoyances(c, pkg)
			}

			time.Sleep(time.Duration(c.TimerDelay) * time.Second)
		}

		time.Sleep(time.Duration(hibernationCountdown) * time.Second)
		inHibernation = false
	}
}

func DoAnnoyances(c types.Config, pkg types.EdgewarePackage) {
	if c.AnnoyancePopups && (rand.Intn(101) + 1) < int(c.PopupChance) {
		image := pkg.ImageFiles[rand.Intn(len(pkg.ImageFiles))]

		go MakePopup(image)
	}

	if c.AnnoyanceVideos && (rand.Intn(101) + 1) < int(c.VideoChance) {
		// Todo
	}

	if c.AnnoyancePrompts && !showingPrompt && (rand.Intn(101) + 1) < int(c.PromptChance) {
		prompts := pkg.Prompts[rand.Intn(len(pkg.Prompts))].Prompts
		prompt := prompts[rand.Intn(len(prompts))]

		go MakePrompt(prompt)
	}

	if c.AnnoyanceAudio && (rand.Intn(101) + 1) < int(c.AudioChance) {
		// Todo
	}
}

func MakePopup(imagePath string) {
	showingPopup = true
	wnd := g.NewMasterWindow("Goonware" + string(rand.Int()), 500, 600, 
		g.MasterWindowFlagsNotResizable | g.MasterWindowFlagsFrameless |
		g.MasterWindowFlagsFloating)
	wnd.Run(func () {
		g.SingleWindow().Layout(
			g.ImageWithFile(imagePath).Size(g.Auto, g.Auto),
			g.Button("Submit <3").OnClick(func() { wnd.Close(); showingPopup = false }),
		)
	})
}

func MakePrompt(text string) {
	showingPrompt = true
	wnd := g.NewMasterWindow("Goonware" + string(rand.Int()), 500, 300, 
		g.MasterWindowFlagsNotResizable | g.MasterWindowFlagsFrameless |
		g.MasterWindowFlagsFloating)
	largerFont := g.GetDefaultFonts()[0].SetSize(20)
	var input string
	wnd.Run(func() {
		g.SingleWindow().Layout(
			g.Label("Repeat.").Font(largerFont),
			g.Dummy(g.Auto/3, 1),
			g.Label(text).Wrapped(true),
			g.Dummy(1, 20),
			g.InputTextMultiline(&input).Size(g.Auto, g.Auto).OnChange(func() {
				if text == input { wnd.Close(); showingPrompt = false}
			}),
		)
	})
}