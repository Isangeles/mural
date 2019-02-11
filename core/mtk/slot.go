/*
 * slot.go
 *
 * Copyright 2019 Dariusz Sikora <dev@isangeles.pl>
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

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	def_slot_color = pixel.RGBA{0.1, 0.1, 0.1, 0.5} // default color
)

// Struct for slot.
type Slot struct {
	bgSpr        *pixel.Sprite
	drawArea     pixel.Rect
	size         Size
	color        color.Color
	fontSize     Size
	label        *Text
	info         *InfoWindow
	icon         *pixel.Sprite
	value        interface{}
	mousePos     pixel.Vec
	onRightClick func(s *Slot)
	onLeftClick  func(s *Slot)
	hovered      bool
	dragged      bool
}

// NewSlot creates new slot without background.
func NewSlot(size, fontSize Size) *Slot {
	s := new(Slot)
	s.size = size
	s.fontSize = fontSize
	s.color = def_slot_color
	// Label & info.
	s.label = NewText(fontSize, 0)
	s.info = NewInfoWindow(SIZE_SMALL, colornames.Grey)
	return s
}

// Draw draws slot.
func (s *Slot) Draw(t pixel.Target, matrix pixel.Matrix) {
	s.drawArea = MatrixToDrawArea(matrix, s.Bounds())
	if s.bgSpr != nil {
		s.bgSpr.Draw(t, matrix)
	} else {
		DrawRectangle(t, s.DrawArea(), s.color)
	}
	if s.icon != nil {
		if s.dragged {
			s.icon.Draw(t, Matrix().Moved(s.mousePos))
		} else {
			s.icon.Draw(t, matrix)
		}
	}
	if s.hovered {
		s.info.Draw(t)
	}
}

// Update updates slot.
func (s *Slot) Update(win *Window) {
	// Mouse events.
	if s.DrawArea().Contains(win.MousePosition()) {
		if win.JustPressed(pixelgl.MouseButtonRight) {
			if s.onRightClick != nil {
				s.onRightClick(s)
			}
		}
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if s.onLeftClick != nil {
				s.onLeftClick(s)
			}
		}
	}
	s.mousePos = win.MousePosition()
	// On-hover.
	s.hovered = s.DrawArea().Contains(s.mousePos)
	// Elements update.
	s.info.Update(win)
}

// Value returns current slot value.
func (s *Slot) Value() interface{} {
	return s.value
}

// Icon returns current slot icon.
func (s *Slot) Icon() *pixel.Sprite {
	return s.icon
}

// Drag toggles slot drag mode(icon
// follows mouse cursor).
func (s *Slot) Drag(drag bool) {
	s.dragged = drag
}

// Dragged checks whether slot is in
// drag mode(icon follows mouse cursor).
func (s *Slot) Dragged() bool {
	return s.dragged
}

// SetColor sets specified color as slot
// color.
func (s *Slot) SetColor(c color.Color) {
	s.color = c
}

// SetIcon sets specified sprite as current
// slot icon.
func (s *Slot) SetIcon(spr *pixel.Sprite) {
	s.icon = spr
}

// SetValue sets specified interface as current
// slot value.
func (s *Slot) SetValue(val interface{}) {
	s.value = val
}

// SetLabel sets specified text as slot label.
func (s *Slot) SetLabel(text string) {
	s.label.SetText(text)
}

// SetInfo sets specified text as content of
// slot info window.
func (s *Slot) SetInfo(text string) {
	s.info.InsertText(text)
}

// Transfer transfers all content(value, icon,
// labels, info) from specified slot.
func (s *Slot) Transfer(oldSlot *Slot) {
	s.SetValue(oldSlot.Value())
	s.SetIcon(oldSlot.Icon())
	s.SetInfo(oldSlot.info.String())
	s.SetLabel(oldSlot.label.String())
	oldSlot.Clear()
}

// Clear removes slot value, icon,
// label and text.
func (s *Slot) Clear() {
	s.SetValue(nil)
	s.SetIcon(nil)
	s.SetLabel("")
	s.SetInfo("")
	s.Drag(false)
}

// DrawArea returns current slot background
// draw area.
func (s *Slot) DrawArea() pixel.Rect {
	return s.drawArea
}

// Bounds returns slot size bounds.
func (s *Slot) Bounds() pixel.Rect {
	if s.bgSpr == nil {
		return s.size.SlotSize()
	}
	return s.bgSpr.Frame()
}

// SetOnLeftClickFunc set speicfied function as function
// triggered after on-click event.
func (s *Slot) SetOnLeftClickFunc(f func(s *Slot)) {
	s.onLeftClick = f
}

// SetOnClickFunc set speicfied function as function
// triggered after on-click event.
func (s *Slot) SetOnRightClickFunc(f func(s *Slot)) {
	s.onRightClick = f
}
