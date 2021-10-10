package simple

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Tab struct {
	btn     widget.Clickable
	Title   string
	Content layout.Widget
}

func (t *Tab) Layout(th *material.Theme) layout.Widget {
	return func(gtx C) D {
		c := t.Content
		if c == nil {
			return layout.Center.Layout(gtx,
				material.H1(th, fmt.Sprintf("Page type: %s", t.Title)).Layout,
			)
		}
		return c(gtx)
	}
}
