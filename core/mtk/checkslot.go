/*
 * checkslot.go
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
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"
)

// Struct for 'chackable' slots.
type CheckSlot struct {
	checked    bool
	bgColor    color.Color
	checkColor color.Color
	drawArea   pixel.Rect
	label      *Text
	value      interface{}
	onCheck    func(s *CheckSlot)
}

// NewCheckSlot creates new item for list.
func NewCheckSlot(label string, value interface{},
	color, checkColor color.Color) *CheckSlot {
	cs := new(CheckSlot)
	cs.bgColor = color
	cs.checkColor = checkColor
	cs.label = NewText(label, SIZE_MEDIUM, 0);
	cs.label.JustLeft()
	cs.value = value
	return cs
}

// Draw draws slot.
func (cs *CheckSlot) Draw(t pixel.Target, drawArea pixel.Rect) {
	cs.drawArea = drawArea
	cs.drawIMBackground(t)
	cs.label.Draw(t, Matrix().Moved(drawArea.Min))
}

// Update updates slot.
func (cs *CheckSlot) Update(win *Window) {
	// Mouse events.
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		if cs.DrawArea().Contains(win.MousePosition()) {
			if cs.Checked() {
				cs.Check(false)
			} else {
				cs.Check(true)
				if cs.onCheck != nil {
					cs.onCheck(cs)
				}
			}
		}
	}
}

// Label returns slot label.
func (cs *CheckSlot) Label() *Text {
	return cs.label
}

// Value returns slot value.
func (cs *CheckSlot) Value() interface{} {
	return cs.value
}

// Bounds returns slot size bounds.
func (cs *CheckSlot) Bounds() pixel.Rect {
	//return SIZE_SMALL.ButtonSize(SHAPE_RECTANGLE)
	return cs.label.Bounds()
}

// DrawArea returns current slot draw area.
func (cs *CheckSlot) DrawArea() pixel.Rect {
	return cs.drawArea
}

// Check toggles slot selection.
func (cs *CheckSlot) Check(check bool) {
	cs.checked = check
}

// Checked checks whether slot is checked.
func (cs *CheckSlot) Checked() bool {
	return cs.checked
}

// SetOnCheckFunc sets specified function as function triggered
// after slot was selected.
func (cs *CheckSlot) SetOnCheckFunc(f func(s *CheckSlot)) {
	cs.onCheck = f
}

// drawIMBackground draws IMDraw background in siaze of
// current draw area.
func (cs *CheckSlot) drawIMBackground(t pixel.Target) {
	color := cs.bgColor
	if cs.Checked() {
		color = cs.checkColor
	}
	draw := imdraw.New(nil)
	draw.Color = pixel.ToRGBA(color)
	draw.Push(cs.drawArea.Min)
	draw.Color = pixel.ToRGBA(color)
	draw.Push(cs.drawArea.Max)
	draw.Rectangle(0)
	draw.Draw(t)
}

