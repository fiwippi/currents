package complex

import (
	"fmt"
	"image"

	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/colorpicker"
	"gioui.org/x/component"
	"gioui.org/x/eventx"
	"github.com/lucasb-eyer/go-colorful"

	"currents/internal/xgio"
	"currents/pkg/audio"
	"currents/pkg/gui/simple"
)

type GradientEditor struct {
	// Gradients that the visualisation widget also has access to, this can b
	// e edited so that changes are also displayed in the visualisation widget
	gradients *audio.Gradients
	// Data for each gradient and its widgets
	data []*gradientData
	// Combobox the visualiser widget is holding
	combobox *xgio.Combo

	// Widgets for the list
	list        layout.List        // List holds buttons to select which gradient to edit
	selected    int                // Index of the selected gradient
	addGradient widget.Clickable   // Button to add a new gradient
	removeBtns  []widget.Clickable // Buttons to remove a gradient from the set

	// Widgets for the editor
	drawMode    *audio.InterpolateMode
	nameField   component.TextField // Editor to change the gradient name
	showPicker  bool                // Whether to show the picker
	closePicker widget.Clickable    // Button to close the colour picker
	pickerIndex int                 // Which colour in the gradient is the picker editing
	pickerState colorpicker.State   // State holds the colour currently in the picker
	addColour   widget.Clickable    // Button to add a new colour to the gradient
	pageScroll  *widget.List        // Embedding all child widgets in this enables a scrollable page
}

func NewGradientEditor(gradients *audio.Gradients, combobox *xgio.Combo, drawMode *audio.InterpolateMode) *GradientEditor {
	creator := &GradientEditor{
		gradients: gradients,
		drawMode:  drawMode,
		list: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.End,
		},
		pageScroll: &widget.List{List: layout.List{Axis: layout.Vertical}},
		nameField: component.TextField{
			Editor: widget.Editor{
				SingleLine: true,
				Submit:     true,
			},
		},
		removeBtns: make([]widget.Clickable, gradients.Size()),
		combobox:   combobox,
	}

	data := make([]*gradientData, 0, gradients.Size())
	for _, name := range gradients.List() {
		d := newGradientData(name, gradients.Get(name))
		data = append(data, d)
	}
	creator.data = data

	return creator
}

func (ge *GradientEditor) Layout(th *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		spy, spyGtx := eventx.Enspy(gtx)

		dims := layout.Flex{
			Axis: layout.Horizontal,
		}.Layout(spyGtx,
			layout.Flexed(3, ge.listWidget(th)),
			layout.Flexed(9, ge.editorWidget(th)),
		)

		// If we are focused on the name editor but the click was
		// not on it then we want to unfocus it
		if ge.disableFocus(spy.AllEvents()) {
			key.FocusOp{}.Add(spyGtx.Ops)
		}

		return dims
	}
}

func (ge *GradientEditor) selectedEntry() *gradientData {
	if len(ge.data) == 0 {
		return nil
	}
	return ge.data[ge.selected]
}

func (ge *GradientEditor) listWidget(th *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Check if we need to add a gradient
		if ge.addGradient.Clicked() {
			var count = 1
			var newName string
		loop:
			for {
				newName = fmt.Sprintf("Untitled #%d", count)
				if !ge.gradients.Has(newName) {
					g := audio.DefaultGradient()
					ge.gradients.Add(newName, g)
					d := newGradientData(newName, g)
					ge.combobox.Add(newName)
					ge.data = append(ge.data, d)
					ge.removeBtns = append(ge.removeBtns, widget.Clickable{})
					break loop
				}
				count++
			}

		}
		// Check if we need to remove a gradient
		for i := range ge.removeBtns {
			if ge.removeBtns[i].Clicked() {
				name := ge.data[i].name
				ge.gradients.Delete(name)
				ge.combobox.Remove(name)

				newData := make([]*gradientData, 0)
				newData = append(newData, ge.data[:i]...)
				ge.data = append(newData, ge.data[i+1:]...)

				if ge.selected == i {
					ge.selected = 0
				}

				break
			}
		}

		// Layout
		return ge.list.Layout(gtx, len(ge.data)+1,
			func(gtx layout.Context, i int) layout.Dimensions {
				var w layout.Widget

				if i < len(ge.data) {
					entry := ge.data[i]

					btn := &entry.selectBtn
					if btn.Clicked() {
						ge.showPicker = false
						ge.selected = i
					}

					var tabWidth int
					w = func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis:      layout.Horizontal,
							Alignment: layout.Middle,
							Spacing:   layout.SpaceSides,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								ico := material.IconButton(th, &ge.removeBtns[i], simple.DeleteIcon)
								ico.Size = unit.Dp(10)
								return ico.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Stack{Alignment: layout.S}.Layout(gtx,
									layout.Stacked(func(gtx layout.Context) layout.Dimensions {
										dims := material.Clickable(gtx, btn, func(gtx layout.Context) layout.Dimensions {
											return layout.UniformInset(unit.Sp(12)).Layout(gtx,
												material.Body1(th, entry.name).Layout,
											)
										})
										tabWidth = dims.Size.X
										return dims
									}),
									layout.Stacked(func(gtx layout.Context) layout.Dimensions {
										if ge.selected != i {
											return layout.Dimensions{}
										}
										tabHeight := gtx.Px(unit.Dp(4))
										tabRect := image.Rect(0, 0, tabWidth, tabHeight)
										paint.FillShape(gtx.Ops, th.Palette.ContrastBg, clip.Rect(tabRect).Op())
										return layout.Dimensions{
											Size: image.Point{X: tabWidth, Y: tabHeight},
										}
									}),
								)
							}),
						)
					}

				} else {
					w = simple.Inset(unit.Dp(9), material.IconButton(th, &ge.addGradient, simple.AddIcon).Layout)
				}

				return w(gtx)
			})
	}
}

func (ge *GradientEditor) editorWidget(th *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		e := ge.selectedEntry()
		if e == nil {
			return layout.Dimensions{}
		}

		// First check if a colour needs to be removed or added
	remove:
		for i := range e.gradient {
			if e.removeBtns[i].Clicked() {
				e.RemoveColour(i)
				if len(e.gradient) == 0 {
					ge.showPicker = false
				}
				// Leave the loop to avoid index error since length of
				// e.gradients has been reduced by 1 but the loop continues
				// for the original length
				break remove
			}
		}
		if ge.addColour.Clicked() {
			e.AddColour()
		}

		// Update the colour the picker is editing
		if ge.pickerState.Changed() {
			clr, _ := colorful.MakeColor(ge.pickerState.Color())
			e.gradient[ge.pickerIndex].Col = clr
		}

		// Set the colour positions to the ones configured right now
		for i, p := range e.positions {
			e.gradient[i].Pos = float64(p.Value)
		}

		// Stops sorting colours and changing their position
		// at the same which causes buggy behaviour
		var dragging bool
	loop:
		for _, s := range e.positions {
			if s.Dragging() {
				dragging = true
				break loop
			}
		}
		if !dragging {
			// Sort the colour so they're in order based on position,
			// this is needed for correct blending
			e.Sort()

			// Ensure the current position sliders reflect the new position values
			for i := range e.positions {
				e.positions[i].Value = float32(e.gradient[i].Pos)
			}
		}

		// Save gradient back to the audio.Gradients the visualiser uses
		ge.gradients.Add(e.name, e.gradient)

		// Edit the gradient name if needed
		if !ge.nameField.Focused() {
			ge.nameField.SetText(e.name)
		} else {
			oldName := e.name
			newName := ge.nameField.Text()

			if !ge.gradients.Has(newName) {
				ge.gradients.Delete(oldName)
				ge.gradients.Add(newName, e.gradient)
				ge.combobox.ChangeName(oldName, newName)
				e.name = newName
			} else if newName != e.name {
				ge.nameField.SetText(e.name)
			}
		}

		// Handle layout
		children := []layout.FlexChild{
			// Edit the gradient name
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ge.nameField.Layout(gtx, th, "Name")
			}),
			// Visualise the gradient
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				w := func(gtx layout.Context) layout.Dimensions {
					defer op.Save(gtx.Ops).Load()
					dr := image.Rectangle{Max: image.Point{X: gtx.Constraints.Max.X - 10, Y: 80}}
					total := dr.Max.X - dr.Min.X
					for x := dr.Max.X - 1; x >= dr.Min.X; x-- {
						c := ge.drawMode.Interpolate(float64(x)/float64(total), e.gradient)
						paint.ColorOp{Color: convertNRGBA(c)}.Add(gtx.Ops)
						clip.Rect(image.Rectangle{Max: image.Point{X: x, Y: dr.Max.Y}}).Add(gtx.Ops)
						paint.PaintOp{}.Add(gtx.Ops)
					}

					return layout.Dimensions{
						Size: dr.Max,
					}
				}
				return simple.Inset(unit.Sp(10), w)(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		}

		if ge.closePicker.Clicked() {
			ge.showPicker = false
		}
		if ge.showPicker {
			picker := func(gtx layout.Context) layout.Dimensions {
				ico := material.IconButton(th, &ge.closePicker, simple.DeleteIcon)
				ico.Size = unit.Dp(15)

				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(ico.Layout),
					layout.Rigid(colorpicker.Picker(th, &ge.pickerState, "").Layout),
				)
			}
			children = append(children, layout.Rigid(picker))
			children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout))
		}

		for idx := range e.gradient {
			i := idx
			col := convertNRGBA((e.gradient)[i].Col)

			if e.pickerBtns[i].Clicked() {
				ge.pickerIndex = i
				ge.pickerState.SetColor(convertNRGBA(e.gradient[i].Col))
				ge.showPicker = true
			}

			// Draw a small box to show the gradient colour at the point
			colourBox := func(gtx layout.Context) layout.Dimensions {
				var tabSize image.Point
				return layout.Stack{Alignment: layout.S}.Layout(gtx,
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						dims := material.Button(th, &e.pickerBtns[i], "######").Layout(gtx)
						tabSize = dims.Size
						return dims
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						margin := 4
						fullRect := image.Rect(0, 0, tabSize.X, tabSize.Y)
						colourRect := image.Rect(margin, margin, tabSize.X-margin, tabSize.Y-margin)
						paint.FillShape(gtx.Ops, th.Palette.Bg, clip.Rect(fullRect).Op())
						paint.FillShape(gtx.Ops, col, clip.Rect(colourRect).Op())
						return layout.Dimensions{
							Size: tabSize,
						}
					}),
				)
			}

			// Slider to edit the position of the gradient colour
			colourSlider := func(gtx layout.Context) layout.Dimensions {
				slider := material.Slider(th, &e.positions[i], 0, 1)
				return slider.Layout(gtx)
			}

			// Button to remove the colour from the gradient
			removeBtn := func(gtx layout.Context) layout.Dimensions {
				ico := material.IconButton(th, &e.removeBtns[i], simple.DeleteIcon)
				ico.Size = unit.Dp(15)
				return ico.Layout(gtx)
			}

			// Ties all the separate widgets into a single flex line
			line := func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(removeBtn),
					layout.Flexed(2, colourBox),
					layout.Flexed(8, colourSlider),
				)
			}

			children = append(children, layout.Rigid(line))
			if len(e.gradient) != i-1 {
				children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout))
			}
		}
		children = append(children, layout.Rigid(simple.Inset(unit.Dp(10), material.IconButton(th, &ge.addColour, simple.AddIcon).Layout)))

		return material.List(th, ge.pageScroll).Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Baseline}.Layout(gtx, children...)
		})
	}
}

func (ge *GradientEditor) disableFocus(groups []eventx.EventGroup) bool {
	action := func() bool {
		for _, g := range groups {
			for _, ev := range g.Items {
				switch e := ev.(type) {
				case pointer.Event:
					if e.Type == pointer.Press {
						return true
					}
				case key.Event:
					if e.Name == key.NameEscape || e.Name == key.NameReturn {
						return true
					}
				}
			}
		}
		return false
	}()

	var selected bool
	for _, ev := range ge.nameField.Events() {
		switch ev.(type) {
		case widget.SelectEvent:
			selected = true
		}
	}

	focused := ge.nameField.Focused()

	return focused && action && !selected
}
