package gui

import (
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/widget/material"
	"github.com/rs/zerolog/log"
)

func Run() {
	// Create the window
	w := defaultWindow()

	// Create the theme and the gradients
	th := material.NewTheme(gofont.Collection())
	gradients := loadGradients()

	// Create the tabs
	tabs := createTabs(th, w, gradients)
	drawFunc := tabs.Layout(th)

	go func() {
		if err := loop(w, drawFunc, gradients); err != nil {
			log.Fatal().Err(err).Msg("gui error")
		}
		os.Exit(0)
	}()

	app.Main()
}
