/*
 * list.go
 *
 * Copyright 2018-2019 Dariusz Sikora <dev@isangeles.pl>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston,
 * MA 02110-1301, USA.
 *
 *
 */

package mtk

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Struct for list with 'selectable' items.
type List struct {
	bgSpr       *pixel.Sprite
	size        Size
	bgColor     color.Color
	secColor    color.Color
	accentColor color.Color
	drawArea    pixel.Rect
	upButton    *Button
	downButton  *Button
	items       []*CheckSlot
	startIndex  int
	selectedVal interface{}
	focused     bool
	disabled    bool
}

// NewList creates new list with specified size
// and colors.
func NewList(size Size, bgColor, secColor,
	accentColor color.Color) *List {
	l := new(List)
	// Background.
	l.size = size
	l.bgColor = bgColor
	l.secColor = secColor
	l.accentColor = accentColor
	// Buttons.
	l.upButton = NewButton(l.size  , SHAPE_SQUARE, accentColor,
		"^", "")
	l.upButton.SetOnClickFunc(l.onButtonUpClicked)
	l.downButton = NewButton(l.size  , SHAPE_SQUARE, accentColor,
		".", "")
	l.downButton.SetOnClickFunc(l.onButtonDownClicked)
	return l
}

// Draw draws list.
func (l *List) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Calculating draw area.
	l.drawArea = MatrixToDrawArea(matrix, l.Bounds())
	// Background.
	if l.bgSpr != nil {
		l.bgSpr.Draw(t, matrix)
	} else {
		DrawRectangle(t, l.DrawArea(), l.bgColor)
	}
	// List.
	l.drawListItems(t)
	// Buttons.
	upButtonPos := MoveTR(l.Bounds(), l.upButton.Frame().Max)
	downButtonPos := MoveBR(l.Bounds(), l.downButton.Frame().Max)
	l.upButton.Draw(t, matrix.Moved(upButtonPos))
	l.downButton.Draw(t, matrix.Moved(downButtonPos))
}

// Update updates list.
func (l *List) Update(win *Window) {
	if l.Disabled() {
		return
	}
	if l.Focused() {
		if win.JustPressed(pixelgl.KeyUp) {
			l.SetStartIndex(l.startIndex - 1)
		}
		if win.JustPressed(pixelgl.KeyDown) {
			l.SetStartIndex(l.startIndex + 1)
		}
	}
	// Buttons.
	l.upButton.Update(win)
	l.downButton.Update(win)
	// List items.
	for _, i := range l.items {
		i.Update(win)
	}
}

// Focus toggles focus on element.
func (l *List) Focus(focus bool) {
	l.focused = focus
}

// Focused checks whether list is focused.
func (l *List) Focused() bool {
	return l.focused
}

// Active toggles list activity.
func (l *List) Active(active bool) {
	l.upButton.Active(active)
	l.downButton.Active(active)
	l.disabled = !active
}

// Disabled checks whether list is disabled.
func (l *List) Disabled() bool {
	return l.disabled
}

// Bounds returns list background size, in from
// of rectangle.
func (l *List) Bounds() pixel.Rect {
	if l.bgSpr == nil {
		return l.size.ListSize()
	}
	return l.bgSpr.Frame()
}

// SetStartIndex sets specified integer as index
// of first item to display. If specified value is
// bigger than last item index then first index(0)
// is set, if is smaller than 0 then last index is
// set.
func (l *List) SetStartIndex(index int) {
	if index > len(l.items)-1 {
		l.startIndex = 0
	} else if index < 0 {
		l.startIndex = len(l.items)-1
	} else {
		l.startIndex = index
	}
}

// InsertItems sets specified values with labels as
// current list content.
func (l *List) InsertItems(content map[string]interface{}) {
	for label, val := range content {
		l.AddItem(label, val)
	}
}

// AddItem adds specified value with label to current
// list content.
func (l *List) AddItem(label string, value interface{}) {
	itemSlot := NewCheckSlot(label, value, l.secColor, l.accentColor)
	itemSlot.SetOnCheckFunc(l.onItemSelected)
	l.items = append(l.items, itemSlot)
}

// SelectedValue returns value of currently selected
// list item.
func (l *List) SelectedValue() interface{} {
	return l.selectedVal
}

// DrawArea returns current list background position
// and size.
func (l *List) DrawArea() pixel.Rect {
	return l.drawArea
}

// drawListItems draws visible list content.
func (l *List) drawListItems(t pixel.Target) {
	if len(l.items) < 1 { // if list empty
		return
	}
	bgTLPos := pixel.V(l.DrawArea().Min.X, l.DrawArea().Max.Y)
	listH := l.DrawArea().H()
	var contentH float64
	// Draw first visible item.
	item := l.items[l.startIndex]
	drawMin := bgTLPos//PosTL(item.Bounds(), bgTLPos)
	drawMax := pixel.V(l.DrawArea().Max.X, drawMin.Y +
		item.Bounds().H())
	drawArea := pixel.R(drawMin.X, drawMin.Y,
		drawMax.X, drawMax.Y)
	item.Draw(t, drawArea)
	contentH += item.DrawArea().H() + ConvSize(15)
	lastItemDA := item.Label().DrawArea()
	// Draw rest of visible items.
	for i := l.startIndex+1; i < len(l.items) && contentH + lastItemDA.H() < listH; i ++ {
		item := l.items[i]
		drawMin := BottomOf(lastItemDA, lastItemDA, 15)
		drawMax := pixel.V(l.DrawArea().Max.X, drawMin.Y +
			item.Bounds().H())
		drawArea := pixel.R(drawMin.X, drawMin.Y,
			drawMax.X, drawMax.Y)
		item.Draw(t, drawArea)
		contentH += item.DrawArea().H() + ConvSize(15)
		lastItemDA = item.Label().DrawArea()
	}
}

// unselectAll unselects all list items.
func (l *List) unselectAll() {
	for _, i := range l.items {
		i.Check(false)
	}
}

// Triggered after button up clicked.
func (l *List) onButtonUpClicked(b *Button) {
	l.SetStartIndex(l.startIndex - 1)
}

// Triggered after button down clicked.
func (l *List) onButtonDownClicked(b *Button) {
	l.SetStartIndex(l.startIndex + 1)
}

// Triggered after list item selected.
func (l *List) onItemSelected(s *CheckSlot) {
	l.unselectAll()
	s.Check(true)
	l.selectedVal = s.Value()
}
