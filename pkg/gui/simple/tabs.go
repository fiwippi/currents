package simple

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Tabs struct {
	Tabs     []Tab
	list     layout.List
	selected int
}

func (t *Tabs) Layout(th *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(t.tabsBarLayout(th)),
			layout.Flexed(1, t.currentTabLayout(th)),
		)
	}
}

func (t *Tabs) currentTabLayout(th *material.Theme) layout.Widget {
	return t.Tabs[t.selected].Layout(th)
}

func (t *Tabs) tabsBarLayout(th *material.Theme) layout.Widget {
	return func(gtx C) D {
		return t.list.Layout(gtx, len(t.Tabs), func(gtx C, tabIdx int) D {
			tab := &t.Tabs[tabIdx]
			if tab.btn.Clicked() {
				t.selected = tabIdx
			}
			var tabWidth int
			return layout.Stack{Alignment: layout.S}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					dims := material.Clickable(gtx, &tab.btn, func(gtx C) D {
						return layout.UniformInset(unit.Sp(12)).Layout(gtx,
							material.Body1(th, tab.Title).Layout,
						)
					})
					tabWidth = dims.Size.X
					return dims
				}),
				layout.Stacked(func(gtx C) D {
					if t.selected != tabIdx {
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
		})
	}
}
