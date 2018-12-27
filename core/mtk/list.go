/*
 * list.go
 *
 * Copyright 2018 Dariusz Sikora <dev@isangeles.pl>
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
	"github.com/faiface/pixel/imdraw"
)

// Struct for list with 'selectable'
// items.
type List struct {
	bgDraw      *imdraw.IMDraw
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
	focused     bool
	disabled    bool
}

// NewList creates new list with specified size
// and colors.
func NewList(size Size, bgColor, secColor,
	accentColor color.Color) *List {
	l := new(List)
	// Background.
	l.bgDraw = imdraw.New(nil)
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
	l.drawArea = MatrixToDrawArea(matrix, l.Frame())
	// Background.
	if l.bgSpr != nil {
		l.bgSpr.Draw(t, matrix)
	} else {
		l.drawIMBackground(t)
	}
	// List.
	l.drawListItems(t)
	// Buttons.
	bgBRPos := pixel.V(l.DrawArea().Max.X, l.DrawArea().Min.Y)
	l.upButton.Draw(t, Matrix().Moved(PosTR(l.upButton.Frame(),
		l.DrawArea().Max)))
	l.downButton.Draw(t, Matrix().Moved(PosBR(l.upButton.Frame(),
		bgBRPos)))
}

// Update updates list.
func (l *List) Update(win *Window) {
	if l.Disabled() {
		return
	}
	l.upButton.Update(win)
	l.downButton.Update(win)
}

// Focus toggles focus on element.
func (l *List) Focus(focus bool) {
	l.focused = focus
}

// Focused checks whether list is focused.
func (l *List) Focused() bool {
	return l.focused
}

func (l *List) Active(active bool) {
	l.upButton.Active(active)
	l.downButton.Active(active)
	l.disabled = !active
}

// Disabled checks whether list is disabled.
func (l *List) Disabled() bool {
	return l.disabled
}

// Frame returns list background size, in from
// of rectangle.
func (l *List) Frame() pixel.Rect {
	if l.bgSpr != nil {
		return l.bgSpr.Frame()
	} else {
		return l.size.ListSize()
	}
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
	itemSlot := NewCheckSlot(label, value, l.secColor)
	l.items = append(l.items, itemSlot)
}

// DrawArea returns current list background position
// and size.
func (l *List) DrawArea() pixel.Rect {
	return l.drawArea
}

// drawIMBackground draws IMDraw background.
func (l *List) drawIMBackground(t pixel.Target) {
	l.bgDraw.Color = pixel.ToRGBA(l.bgColor)
	l.bgDraw.Push(l.DrawArea().Min)
	l.bgDraw.Color = pixel.ToRGBA(l.bgColor)
	l.bgDraw.Push(l.DrawArea().Max)
	l.bgDraw.Rectangle(0)
	l.bgDraw.Draw(t)
}

// drawListItems draws visible list content.
func (l *List) drawListItems(t pixel.Target) {
	bgTLPos := pixel.V(l.DrawArea().Min.X, l.DrawArea().Max.Y)
	listH := l.DrawArea().H()
	var contentH float64
	var lastItemDA pixel.Rect
	for i := l.startIndex; i < len(l.items) && contentH + lastItemDA.H() < listH; i ++ {
		item := l.items[i]
		if i == l.startIndex {
			drawMin := PosTL(item.Bounds(), bgTLPos)
			//itemPos.Y -= item.Label().Bounds().H()
			l.drawItem(t, item, drawMin)
			contentH += item.DrawArea().H() + ConvSize(15)
			lastItemDA = item.Label().DrawArea()
			continue
		}
		drawMin := BottomOf(lastItemDA, item.Bounds(), 15)
		//itemPos.Y -= item.Label().Bounds().H()
		l.drawItem(t, item, drawMin)
		contentH += item.DrawArea().H() + ConvSize(15)
		lastItemDA = item.Label().DrawArea()
	}
}

// drawItemBackground draws list item background.
func (l *List) drawItem(t pixel.Target, item *CheckSlot,
	drawAreaMin pixel.Vec) {
	drawAreaMax := pixel.V(l.DrawArea().Max.X, drawAreaMin.Y +
		item.Bounds().H())
	drawArea := pixel.R(drawAreaMin.X, drawAreaMin.Y,
		drawAreaMax.X, drawAreaMax.Y)
	item.Draw(t, drawArea)
}

// Triggered after button up clicked.
func (l *List) onButtonUpClicked(b *Button) {
	l.SetStartIndex(l.startIndex - 1)
}

// Triggered after button down clicked.
func (l *List) onButtonDownClicked(b *Button) {
	l.SetStartIndex(l.startIndex + 1)
}
