/*
 * checkslot.go
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

// Struct for 'chackable' slots.
type CheckSlot struct {
	checked    bool
	bgSize     pixel.Vec
	bgColor    color.Color
	checkColor color.Color
	drawArea   pixel.Rect
	label      *Text
	value      interface{}
	onCheck    func(s *CheckSlot)
}

// NewCheckSlot creates new item for list.
func NewCheckSlot(label string, value interface{}, bgSize pixel.Vec,
	color, checkColor color.Color) *CheckSlot {
	cs := new(CheckSlot)
	cs.bgSize = bgSize
	cs.bgColor = color
	cs.checkColor = checkColor
	cs.label = NewText(SIZE_MEDIUM, 0);
	cs.label.SetText(label)
	cs.value = value
	return cs
}

// Draw draws slot.
func (cs *CheckSlot) Draw(t pixel.Target, matrix pixel.Matrix) {
	cs.drawArea = MatrixToDrawArea(matrix, cs.Bounds())
	color := cs.bgColor
	if cs.Checked() {
		color = cs.checkColor
	}
	DrawRectangle(t, cs.DrawArea(), color)
	labelMove := MoveBL(cs.Bounds(), cs.label.Bounds().Size())
	cs.label.Draw(t, matrix.Moved(labelMove))
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

// SetSize sets specified vector as
// background size.
func (cs *CheckSlot) SetSize(s pixel.Vec) {
	cs.bgSize = s
}

// Label returns slot label.
func (cs *CheckSlot) Label() string {
	return cs.label.String()
}

// Value returns slot value.
func (cs *CheckSlot) Value() interface{} {
	return cs.value
}

// Bounds returns slot size bounds.
func (cs *CheckSlot) Bounds() pixel.Rect {
	return pixel.R(0, 0, cs.bgSize.X, cs.bgSize.Y)
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

