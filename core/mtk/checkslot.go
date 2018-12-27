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
	"github.com/faiface/pixel/imdraw"
)

// Struct for 'chackable' slots.
type CheckSlot struct {
	selected bool
	bgColor  color.Color
	drawArea pixel.Rect
	label    *Text
	value    interface{}
}

// NewCheckSlot creates new item for list.
func NewCheckSlot(label string, value interface{},
	color color.Color) *CheckSlot {
	cs := new(CheckSlot)
	cs.bgColor = color
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

// drawIMBackground draws IMDraw background in siaze of
// current draw area.
func (cs *CheckSlot) drawIMBackground(t pixel.Target) {
	draw := imdraw.New(nil)
	draw.Color = pixel.ToRGBA(cs.bgColor)
	draw.Push(cs.drawArea.Min)
	draw.Color = pixel.ToRGBA(cs.bgColor)
	draw.Push(cs.drawArea.Max)
	draw.Rectangle(0)
	draw.Draw(t)
}

