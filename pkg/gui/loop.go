package gui

import (
	"encoding/json"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"

	"currents/internal/log"
	"currents/pkg/audio"
)

func loop(w *app.Window, drawLayout layout.Widget, gradients *audio.Gradients) error {
	var ops op.Ops

	for {
		e := <-w.Events()

		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			drawLayout(gtx)
			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			// Save gradients on exit
			data, err := json.MarshalIndent(gradients, "", "    ")
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal gradients")
			}

			err = os.WriteFile("gradients.json", data, 0644)
			if err != nil {
				log.Error().Err(err).Msg("failed to save to gradients.json")
			}

			return e.Err
		}
	}
}
