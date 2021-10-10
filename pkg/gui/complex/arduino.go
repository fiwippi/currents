package complex

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"currents/internal/log"
	"currents/internal/xgio"
	"currents/internal/xmaterial"
	"currents/pkg/gui/simple"
	"currents/pkg/session"
)

type ArduinoController struct {
	// Port data
	ports    []string
	portName string
	server   *session.Server

	// Widgets
	startBtn      widget.Clickable
	stopBtn       widget.Clickable
	portsCombobox xgio.Combo
}

func NewArduinoController(server *session.Server) *ArduinoController {
	ac := &ArduinoController{
		server: server,
	}

	// Load the possible ports
	ports, err := session.GetAvailablePorts()
	if err != nil {
		log.Fatal().Err(err).Msg("could not get port data")
	}
	ac.ports = ports

	// Create the ports combobox
	ac.portsCombobox = xgio.MakeCombo(ports, "Select a port")
	if ac.portsCombobox.Len() > 0 {
		ac.portsCombobox.SelectIndex(0)
		ac.portName = ac.portsCombobox.SelectedText()
	}

	return ac
}

func (ac *ArduinoController) Layout(th *material.Theme) layout.Widget {
	return simple.Inset(unit.Dp(5),
		func(gtx layout.Context) layout.Dimensions {
			// Handle logic
			if ac.startBtn.Clicked() {
				err := ac.server.Connect(ac.portName)
				if err != nil {
					log.Error().Err(err).Msg("failed to connect to arduino")
				}
			}

			if ac.stopBtn.Clicked() {
				err := ac.server.Disconnect()
				if err != nil {
					log.Error().Err(err).Msg("failed to disconnect from arduino")
				}
			}

			// Layout
			return layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Start,
			}.Layout(gtx,
				layout.Rigid(material.H6(th, "Controls:").Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					dim := layout.Flex{
						Axis:    layout.Horizontal,
						Spacing: layout.SpaceEvenly,
					}.Layout(gtx,
						layout.Flexed(3, material.Button(th, &ac.startBtn, "Start").Layout),
						layout.Flexed(1, layout.Spacer{}.Layout),
						layout.Flexed(3, material.Button(th, &ac.stopBtn, "Stop").Layout),
					)
					return dim
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(material.H6(th, "Port:").Layout),
				layout.Rigid(xmaterial.Combo(th, &ac.portsCombobox).Layout),
			)
		},
	)
}
