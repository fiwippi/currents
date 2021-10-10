package xmaterial

import (
	l "gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	giomat "gioui.org/widget/material"
)

// RigidButton returns layout function for a button with inset
func RigidButton(th *giomat.Theme, caption string, button *widget.Clickable) l.FlexChild {
	inset := l.UniformInset(unit.Dp(3))
	return l.Rigid(func(gtx l.Context) l.Dimensions {
		return inset.Layout(gtx, func(gtx l.Context) l.Dimensions {
			return giomat.Button(th, button, caption).Layout(gtx)
		})
	})
}
