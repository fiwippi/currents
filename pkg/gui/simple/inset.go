package simple

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

func Inset(inset unit.Value, w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(inset).Layout(gtx, w)
	}
}
