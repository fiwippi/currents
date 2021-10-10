package simple

import (
	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var AddIcon *widget.Icon
var DeleteIcon *widget.Icon

func init() {
	var err error

	AddIcon, err = widget.NewIcon(icons.ContentAdd)
	if err != nil {
		panic(err)
	}
	DeleteIcon, err = widget.NewIcon(icons.NavigationClose)
	if err != nil {
		panic(err)
	}
}
