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
	"golang.org/x/image/colornames"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	slot_color = pixel.RGBA{0.1, 0.1, 0.1, 0.5}
)

// Struct for slot.
type Slot struct {
	bgSpr    *pixel.Sprite
	drawArea pixel.Rect
	size     Size
	fontSize Size
	shape    Shape
	label    *Text
	info     *InfoWindow
	icon     *pixel.Sprite
	value    interface{}
	onClick  func(s *Slot)
	hovered  bool
}

// NewSlot creates new slot without background.
func NewSlot(size, fontSize Size, shape Shape) *Slot {
	s := new(Slot)
	s.size = size
	s.fontSize = fontSize
	s.shape = shape
	// Label & info.
	s.label = NewText("", fontSize, 0)
	s.info = NewInfoWindow(SIZE_SMALL, colornames.Grey)
	return s
}

// NewSlotSprite creates new slot with specified picture
// as background.
func NewSlotSprite(bg pixel.Picture, fontSize Size) *Slot {
	s := new(Slot)
	s.fontSize = fontSize
	s.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	// Label & info.
	s.label = NewText("", fontSize, 0)
	s.info = NewInfoWindow(SIZE_SMALL, colornames.Grey)
	return s
}

// Draw draws slot.
func (s *Slot) Draw(t pixel.Target, matrix pixel.Matrix) {
	s.drawArea = MatrixToDrawArea(matrix, s.Bounds())
	if s.bgSpr != nil {
		s.bgSpr.Draw(t, matrix)
	} else {
		DrawRectangle(t, s.DrawArea(), slot_color)
	}
	if s.icon != nil {
		s.icon.Draw(t, matrix)
	}
	if s.hovered {
		s.info.Draw(t)
	}
}

// Update updates slot.
func (s *Slot) Update(win *Window) {
	// Mouse events.
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		if s.onClick != nil {
			s.onClick(s)
		}
	}
	// On-hover.
	s.hovered = s.DrawArea().Contains(win.MousePosition())
	// Elements update.
	s.info.Update(win)
}

// Value returns current slot value.
func (s *Slot) Value() interface{} {
	return s.value
}

// SetIcon sets specified sprite as current
// slot icon.
func (s *Slot) SetIcon(icon pixel.Picture) {
	s.icon = pixel.NewSprite(icon, s.Bounds())
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

// DrawArea returns current slot background
// draw area.
func (s *Slot) DrawArea() pixel.Rect {
	return s.drawArea
}

// Bounds returns slot size bounds.
func (s *Slot) Bounds() pixel.Rect {
	if s.bgSpr == nil {
		return s.size.SlotSize(s.shape)
	}
	return s.bgSpr.Frame()
}

// SetOnClickFunc set speicfied function as function
// triggered after on-click event.
func (s *Slot) SetOnClickFunc(f func(s *Slot)) {
	s.onClick = f
}
