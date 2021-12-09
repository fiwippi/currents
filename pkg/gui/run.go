package gui

import (
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/widget/material"
	"github.com/rs/zerolog/log"

	"currents/pkg/session"
)

func Run() {
	// Create the window
	w := defaultWindow()

	// Create the theme and the gradients and the server
	server := session.NewServer()
	th := material.NewTheme(gofont.Collection())
	gradients := loadGradients()

	// Create the tabs
	tabs := createTabs(th, w, gradients, server)
	drawFunc := tabs.Layout(th)

	go func() {
		// Run the event loop until finish/error
		err := loop(w, drawFunc, gradients)

		// Always try to close the connection to the arduino
		arduinoErr := server.Disconnect()
		if arduinoErr != nil {
			log.Error().Err(arduinoErr).Msg("failed to close arduino connection on exit")
		}

		// Close the program based on the gui's exit status
		if err != nil {
			log.Fatal().Err(err).Msg("gui error")
		}
		os.Exit(0)
	}()

	app.Main()
}
