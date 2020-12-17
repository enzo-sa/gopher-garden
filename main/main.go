package main

import (
	"gioui.org/app"
	"gioui.org/unit"
	"github.com/enzo-sa/gopher-garden/ui"
	"log"
	"os"
)

// driver for game
// inits window and calls ui loop
// loop runs until window exits or player exits through menu
func main() {
	go func() {
		w := app.NewWindow(
			app.Title("Gopher-Garden"),
			app.Size(unit.Dp(ui.WidthPx+500), unit.Dp(ui.HeightPx)))
		u := ui.NewUi(w)
		if err := u.Loop(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

