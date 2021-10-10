package complex

import (
	"sort"

	"gioui.org/widget"

	"currents/pkg/audio"
)

type gradientData struct {
	name     string
	gradient audio.Gradient

	// Button to select the gradient in the editor list
	selectBtn widget.Clickable
	// Picker buttons to display each colour at each position
	// and shows a color picker to change the colour specified
	pickerBtns []widget.Clickable
	// Sliders to change the position of the given colour
	positions []widget.Float
	// Buttons to remove the specified colour from the gradient
	removeBtns []widget.Clickable
}

func newGradientData(name string, gradient audio.Gradient) *gradientData {
	d := &gradientData{
		name:       name,
		gradient:   gradient,
		pickerBtns: make([]widget.Clickable, len(gradient)),
		positions:  make([]widget.Float, len(gradient)),
		removeBtns: make([]widget.Clickable, len(gradient)),
	}

	for i := range d.positions {
		d.positions[i].Value = float32(d.gradient[i].Pos)
	}

	return d
}

func (gd *gradientData) Sort() {
	posSort := func(i, j int) bool {
		return gd.gradient[i].Pos < gd.gradient[j].Pos
	}
	sort.Slice(gd.pickerBtns, posSort)
	sort.Slice(gd.positions, posSort)
	sort.Slice(gd.removeBtns, posSort)
	sort.Slice(gd.gradient, posSort)
}

func (gd *gradientData) RemoveColour(index int) {
	a := make(audio.Gradient, 0)
	a = append(a, gd.gradient[:index]...)
	gd.gradient = append(a, gd.gradient[index+1:]...)

	b := make([]widget.Clickable, 0)
	b = append(b, gd.pickerBtns[:index]...)
	gd.pickerBtns = append(b, gd.pickerBtns[index+1:]...)

	c := make([]widget.Float, 0)
	c = append(c, gd.positions[:index]...)
	gd.positions = append(c, gd.positions[index+1:]...)

	d := make([]widget.Clickable, 0)
	d = append(d, gd.removeBtns[:index]...)
	gd.removeBtns = append(d, gd.removeBtns[index+1:]...)
}

func (gd *gradientData) AddColour() {
	gd.gradient = append(gd.gradient, audio.DefaultGradient()...)
	gd.pickerBtns = append(gd.pickerBtns, widget.Clickable{})
	gd.positions = append(gd.positions, widget.Float{Value: 1.0})
	gd.removeBtns = append(gd.removeBtns, widget.Clickable{})
}
