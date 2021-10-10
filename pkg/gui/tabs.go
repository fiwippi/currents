package gui

import (
	"gioui.org/app"
	"gioui.org/widget/material"

	"currents/pkg/audio"
	"currents/pkg/gui/complex"
	"currents/pkg/gui/simple"
	"currents/pkg/session"
)

func createTabs(th *material.Theme, w *app.Window, gradients *audio.Gradients) simple.Tabs {
	drawMode := audio.Blended
	server := session.NewServer()
	// Redrawing happens outside a frame event so we need to call
	// window.Invalidate instead of using op.InvalidateOp
	v := complex.NewVisualisation(gradients, func() { w.Invalidate() }, &drawMode, server)
	ge := complex.NewGradientEditor(gradients, v.GradientsCombobox(), &drawMode)
	ac := complex.NewArduinoController(server)

	var tabs simple.Tabs
	tabs.Tabs = append(tabs.Tabs,
		simple.Tab{Title: "Visualisation", Content: v.Layout(th)},
		simple.Tab{Title: "Gradients", Content: ge.Layout(th)},
		simple.Tab{Title: "Arduino", Content: ac.Layout(th)},
	)

	return tabs
}
