package complex

import (
	"image"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/lucasb-eyer/go-colorful"

	"currents/internal/log"
	"currents/internal/xgio"
	"currents/internal/xmaterial"
	"currents/pkg/audio"
	"currents/pkg/gui/simple"
	"currents/pkg/session"
)

type Visualisation struct {
	// Audio
	audio           *audio.Audio
	audioConfig     *audio.Config
	currentColour   colorful.Color
	devices         []audio.Device
	fft             *audio.FFT
	gradients       *audio.Gradients
	currentDevice   string
	currentGradient string
	started         bool
	drawMode        *audio.InterpolateMode
	defaultDamp     float32

	// Session
	session *session.Server

	// Widgets
	startBtn          widget.Clickable
	stopBtn           widget.Clickable
	gradientsCombobox xgio.Combo
	devicesCombobox   xgio.Combo
	dampCheckbox      widget.Bool
	drawModes         widget.Enum
	dampSlider        widget.Float
	dampReset         widget.Clickable
}

func NewVisualisation(gradients *audio.Gradients, redraw func(), drawMode *audio.InterpolateMode, server *session.Server) *Visualisation {
	v := &Visualisation{
		audio:             audio.MustCreateNewAudio(),
		audioConfig:       audio.DefaultConfig(),
		gradientsCombobox: xgio.Combo{},
		devicesCombobox:   xgio.Combo{},
		gradients:         gradients,
		session:           server,
		drawMode:          drawMode,
	}

	// Load the possible gradients
	v.gradientsCombobox = xgio.MakeCombo(v.gradients.List(), "Select a gradient")
	if v.gradientsCombobox.Len() > 0 {
		v.gradientsCombobox.SelectIndex(0)
		v.currentGradient = v.gradientsCombobox.SelectedText()
	}

	// Load the possible devices
	v.devices = v.audio.MustParseDevices()
	deviceList := make([]string, 0, len(v.devices))
	for _, d := range v.devices {
		deviceList = append(deviceList, d.Name)
	}
	v.devicesCombobox = xgio.MakeCombo(deviceList, "Select an audio device")
	if v.devicesCombobox.Len() > 0 {
		v.devicesCombobox.SelectIndex(0)
		v.currentDevice = v.devicesCombobox.SelectedText()
	}

	// Create the fft context and start waiting for colours to input
	v.fft = audio.MustCreateNewFFT(v.audioConfig)
	v.dampCheckbox.Value = v.fft.Damp
	g := v.gradients.Get(v.gradientsCombobox.SelectedText())
	v.fft.Gradient = &g
	v.fft.DrawMode = *drawMode
	v.drawModes.Value = drawMode.String()
	v.dampSlider.Value = float32(v.fft.SampleRate.Milliseconds())
	v.defaultDamp = v.dampSlider.Value

	go func() {
	loop:
		for {
			select {
			case err := <-v.audio.Done:
				log.Debug().Err(err).Msg("audio capture done")
			case err := <-v.fft.Done:
				log.Debug().Err(err).Msg("fft done")
				if err != nil {
					log.Fatal().Err(err).Msg("fft buffer failed")
				}
				break loop
			case hue := <-v.fft.Hues:
				v.currentColour = hue
				v.session.SendColour(hue)
				redraw()
				log.Trace().Str("hue", hue.Hex()).Msg("fft hue received")
			}
		}
	}()

	return v
}

func (v *Visualisation) Layout(th *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Handles logic i.e. controlling the audio context
		v.handleLogic()

		// Draw the layout
		controls := simple.Inset(unit.Dp(5),
			func(gtx layout.Context) layout.Dimensions {
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
							layout.Flexed(3, material.Button(th, &v.startBtn, "Start").Layout),
							layout.Flexed(1, layout.Spacer{}.Layout),
							layout.Flexed(3, material.Button(th, &v.stopBtn, "Stop").Layout),
						)
						return dim
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
					layout.Rigid(material.H6(th, "Gradient:").Layout),
					layout.Rigid(xmaterial.Combo(th, &v.gradientsCombobox).Layout),
					layout.Rigid(material.RadioButton(th, &v.drawModes, audio.Blended.String(), audio.Blended.String()).Layout),
					layout.Rigid(material.RadioButton(th, &v.drawModes, audio.Blocky.String(), audio.Blocky.String()).Layout),
					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
					layout.Rigid(material.H6(th, "Device:").Layout),
					layout.Rigid(xmaterial.Combo(th, &v.devicesCombobox).Layout),
					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
					layout.Rigid(material.H6(th, "Damping:").Layout),
					layout.Rigid(material.CheckBox(th, &v.dampCheckbox, "On/Off").Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if !(v.dampCheckbox.Value) {
							gtx = gtx.Disabled()
						}
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Flexed(1, material.Slider(th, &v.dampSlider, 1, 600).Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.UniformInset(unit.Dp(8)).Layout(gtx,
									material.Body2(th, "Strength").Layout,
								)
							}),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if !(v.dampCheckbox.Value) {
							gtx = gtx.Disabled()
						}
						return material.Button(th, &v.dampReset, "Reset Damping").Layout(gtx)
					}),
				)
			},
		)

		colourBox := simple.Inset(unit.Dp(5),
			func(gtx layout.Context) layout.Dimensions {
				dr := image.Rectangle{Max: gtx.Constraints.Min}
				defer op.Save(gtx.Ops).Load()
				paint.ColorOp{Color: convertNRGBA(v.currentColour)}.Add(gtx.Ops)
				clip.Rect(dr).Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)

				return layout.Dimensions{
					Size: gtx.Constraints.Max,
				}
			},
		)

		page := layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(gtx,
			layout.Flexed(4, controls),
			layout.Flexed(7, colourBox),
		)

		return page
	}
}

func (v *Visualisation) GradientsCombobox() *xgio.Combo {
	return &v.gradientsCombobox
}

// Audio Control

func (v *Visualisation) handleLogic() {
	// Controls
	if v.startBtn.Clicked() {
		v.startCapture()
	}
	if v.stopBtn.Clicked() {
		v.stopCapture()
	}

	// Gradient
	if v.gradientsCombobox.SelectedText() != v.currentGradient {
		v.currentGradient = v.gradientsCombobox.SelectedText()
		gradient := v.gradients.Get(v.currentGradient)
		v.fft.Gradient = &gradient
		log.Debug().Str("name", v.currentGradient).Msg("fft gradient changed")
	}
	if v.drawModes.Changed() {
		switch v.drawModes.Value {
		case audio.Blended.String():
			*v.drawMode = audio.Blended
			v.fft.DrawMode = audio.Blended
		case audio.Blocky.String():
			*v.drawMode = audio.Blocky
			v.fft.DrawMode = audio.Blocky
		}

	}

	// Device
	if v.started && v.devicesCombobox.SelectedText() != v.currentDevice {
		v.currentDevice = v.devicesCombobox.SelectedText()
		v.startCapture()
	}

	// Options
	if v.dampCheckbox.Changed() {
		v.fft.Damp = v.dampCheckbox.Value
		log.Debug().Bool("value", v.fft.Damp).Msg("fft damp toggled changed")
	}
	if v.dampSlider.Changed() {
		v.fft.ChangeSampleRate(time.Duration(v.dampSlider.Value) * time.Millisecond)
		log.Debug().Float32("value", v.dampSlider.Value).Msg("fft damp value changed")
	}
	if v.dampReset.Clicked() {
		v.dampSlider.Value = v.defaultDamp
		v.fft.ChangeSampleRate(time.Duration(v.dampSlider.Value) * time.Millisecond)
		log.Debug().Float32("value", v.dampSlider.Value).Msg("fft damp reset")
	}
}

func (v *Visualisation) startCapture() {
	name := v.devicesCombobox.SelectedText()
	var device audio.Device
	for _, d := range v.devices {
		if d.Name == name {
			device = d
		}
	}

	err := v.audio.StartCapture(device, v.fft, v.audioConfig)
	if err != nil {
		log.Error().Err(err).Str("device", name).Msg("failed to start capture")
	}
	v.started = true
	log.Debug().Err(err).Str("device", name).Msg("started capture")
}

func (v *Visualisation) stopCapture() {
	go v.audio.StopCapture()
	// We get Frame events because we are clicking the button so no need to call invalidate,
	// we only have to change the current colour
	v.started = false
	v.currentColour = colorful.Color{R: 0, G: 0, B: 0}
	log.Debug().Msg("stopped capture")
}
