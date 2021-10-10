package xgio

// TODO license for giox

import (
	"errors"

	"gioui.org/widget"
)

// Combo holds combobox state
type Combo struct {
	items        []string
	hint         string
	selected     int
	expanded     bool
	selectButton widget.Clickable
	buttons      []widget.Clickable
}

// MakeCombo Creates new combobox widget
func MakeCombo(items []string, hint string) Combo {
	c := Combo{
		items:    items,
		hint:     hint,
		selected: -1,
		expanded: false,
		buttons:  make([]widget.Clickable, len(items)),
	}

	return c
}

// HasSelected returns true if an item is selected
func (c *Combo) HasSelected() bool {
	return c.selected != -1
}

// IsExpanded checks wheather box is expanded
func (c *Combo) IsExpanded() bool {
	return c.expanded
}

// Toggle expands and collapses a combobox
func (c *Combo) Toggle() bool {
	c.expanded = !c.expanded
	return c.expanded
}

// Len returns number of items
func (c *Combo) Len() int {
	return len(c.items)
}

// Items returns current list of items
func (c *Combo) Items() []string {
	return c.items
}

// Hint returns control's hint test
func (c *Combo) Hint() string {
	return c.hint
}

// Item returns a text for corresponding item
func (c *Combo) Item(index int) string {
	return c.items[index]
}

// SelectButton returns a points to main (open) combobox button
func (c *Combo) SelectButton() *widget.Clickable {
	return &c.selectButton
}

// Button returns a pointer to correspoding button widget
func (c *Combo) Button(index int) *widget.Clickable {
	return &(c.buttons[index])
}

// SelectedText returns currently selected item
func (c *Combo) SelectedText() string {
	if c.selected == -1 {
		return c.hint
	}

	return c.items[c.selected]
}

// SelectIndex sets currently selected item by index
func (c *Combo) SelectIndex(index int) error {
	N := len(c.items)
	if index != -1 && (index < 0 || index >= N) {
		return errors.New("Combobox: bad index")
	}

	c.selected = index
	return nil
}

// SelectItem sets currently selected item by value
func (c *Combo) SelectItem(item string) error {
	for i, val := range c.items {
		if val == item {
			c.selected = i
			return nil
		}
	}

	return errors.New("Combobox: bad index")
}

// Unselect removes current selection
func (c *Combo) Unselect() {
	c.selected = -1
}

// Add adds an additional item to the combobox
func (c *Combo) Add(item string) {
	c.items = append(c.items, item)
	c.buttons = append(c.buttons, widget.Clickable{})
}

// Remove deletes an element from the combobox
func (c *Combo) Remove(item string) {
	if len(c.items) == 0 {
		return
	}

	var index int
	for i := range c.items {
		if c.items[i] == item {
			index = i
		}
	}

	newItems := make([]string, 0)
	newItems = append(newItems, c.items[:index]...)
	c.items = append(newItems, c.items[index+1:]...)

	newBtns := make([]widget.Clickable, 0)
	newBtns = append(newBtns, c.buttons[:index]...)
	c.buttons = append(newBtns, c.buttons[index+1:]...)

	if len(c.items) == 0 {
		c.selected = -1
	} else {
		c.selected = 0
	}
}

// ChangeName changes the name of an item if it exists
func (c *Combo) ChangeName(old, new string) {
	var index int
	for i := range c.items {
		if c.items[i] == old {
			index = i
			break
		}
	}

	c.items[index] = new
}
