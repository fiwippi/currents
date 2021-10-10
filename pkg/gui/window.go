package gui

import (
	"gioui.org/app"
	"gioui.org/unit"
)

const windowWidth = 800
const windowHeight = 600

func defaultWindow() *app.Window {
	return app.NewWindow(
		app.Title("Currents"),
		app.MaxSize(unit.Dp(windowWidth), unit.Dp(windowHeight)),
		app.MinSize(unit.Dp(windowWidth), unit.Dp(windowHeight)),
	)
}
