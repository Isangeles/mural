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
	"fmt"
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"
)

var (
	def_slot_color = pixel.RGBA{0.1, 0.1, 0.1, 0.5} // default color
)

// Struct for slot.
type Slot struct {
	bgSpr              *pixel.Sprite
	bgDraw             *imdraw.IMDraw
	drawArea           pixel.Rect
	size               Size
	color              color.Color
	fontSize           Size
	label              *Text
	countLabel         *Text
	info               *InfoWindow
	icon               *pixel.Sprite
	values             []interface{}
	mousePos           pixel.Vec
	specialKey         pixelgl.Button
	onRightClick       func(s *Slot)
	onLeftClick        func(s *Slot)
	onSpecialLeftClick func(s *Slot)
	hovered            bool
	dragged            bool
}

// NewSlot creates new slot without background.
func NewSlot(size, fontSize Size) *Slot {
	s := new(Slot)
	// Background.
	s.bgDraw = imdraw.New(nil)
	s.size = size
	s.color = def_slot_color
	// Labels & info.
	s.fontSize = fontSize
	s.label = NewText(fontSize, 0)
	s.countLabel = NewText(fontSize, 0)
	s.countLabel.JustCenter()
	s.info = NewInfoWindow(SIZE_SMALL, colornames.Grey)
	return s
}

// SlotSwitch transfers all contant of slot A
// (value, icon, label, info) to slot B and
// vice versa.
func SlotSwitch(slotA, slotB *Slot) {
	slotC := *slotA
	slotA.SetValues(slotB.Values())
	slotA.SetIcon(slotB.Icon())
	slotA.SetInfo(slotB.info.String())
	slotA.SetLabel(slotB.label.String())
	slotB.SetValues(slotC.Values())
	slotB.SetIcon(slotC.Icon())
	slotB.SetInfo(slotC.info.String())
	slotB.SetLabel(slotC.label.String())	
}

// SlotCopy copies content from A to
// slot B(overwrites current content).
func SlotCopy(slotA, slotB *Slot) {
	slotB.SetValues(slotA.Values())
	slotB.SetIcon(slotA.Icon())
	slotB.SetInfo(slotA.info.String())
	slotB.SetLabel(slotA.info.String())
}

// Draw draws slot.
func (s *Slot) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Draw area.
	s.drawArea = MatrixToDrawArea(matrix, s.Bounds())
	// Icon.
	if s.icon != nil && s.icon.Picture() != nil {
		if s.dragged {
			s.icon.Draw(t, Matrix().Moved(s.mousePos))
		} else {
			s.icon.Draw(t, matrix)
		}
	}
	// Slot.
	if s.bgSpr != nil {
		s.bgSpr.Draw(t, matrix)
	} else {
		s.drawBG(t)
	}
	// Labels.
	if len(s.values) > 0 {
		countPos := MoveTL(s.DrawArea(), s.countLabel.Bounds().Size())
		s.countLabel.Draw(t, matrix.Moved(countPos))
	}
	// Info window.
	if s.hovered {
		s.info.Draw(t)
	}
}

// Update updates slot.
func (s *Slot) Update(win *Window) {
	// Mouse position.
	s.mousePos = win.MousePosition()
	// Mouse events.
	if s.DrawArea().Contains(s.mousePos) {
		switch {
		case win.JustPressed(pixelgl.MouseButtonRight):
			if s.onRightClick != nil {
				s.onRightClick(s)
			}
		case s.specialKey != 0 && win.Pressed(s.specialKey) &&
			win.JustPressed(pixelgl.MouseButtonLeft):
			if s.onSpecialLeftClick != nil {
				s.onSpecialLeftClick(s)
			}
		case win.JustPressed(pixelgl.MouseButtonLeft):
			if s.onLeftClick != nil {
				s.onLeftClick(s)
			}
		}
	}
	// On-hover.
	s.hovered = s.DrawArea().Contains(s.mousePos)
	// Count label.
	s.countLabel.SetText(fmt.Sprintf("%d", len(s.values)))
	// Elements update.
	s.info.Update(win)
}

// Values returns all slot values.
func (s *Slot) Values() []interface{} {
	return s.values
}

// Pop removes and returns first value
// from slot. Clears slot if removed value
// was last value in slot.
func (s *Slot) Pop() interface{} {
	if s.values == nil {
		return nil
	}
	lastID := len(s.values)-1
	v := s.values[lastID]
	s.values = s.values[:lastID]
	if len(s.values) < 1 {
		s.Clear()
	}
	return v
}

// Icon returns current slot icon
// picture.
func (s *Slot) Icon() pixel.Picture {
	if s.icon == nil {
		return nil
	}
	return s.icon.Picture()
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
func (s *Slot) SetIcon(pic pixel.Picture) {
	s.icon = pixel.NewSprite(pic, s.Bounds())
}

// AddValue adds specified interface to slot
// values list.
func (s *Slot) AddValues(vls... interface{}) {
	s.values = append(s.values, vls...)
}

// SetValues replaces current values with
// specified ones.
func (s *Slot) SetValues(vls []interface{}) {
	s.values = vls
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

// Clear removes slot value, icon,
// label and text.
func (s *Slot) Clear() {
	s.SetValues(nil)
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

// SetSpecialKey sets special key for slot click
// events.
func (s *Slot) SetSpecialKey(k pixelgl.Button) {
	s.specialKey = k
}

// SetOnClickFunc set speicfied function as function
// triggered after right mouse click event.
func (s *Slot) SetOnRightClickFunc(f func(s *Slot)) {
	s.onRightClick = f
}

// SetOnLeftClickFunc set speicfied function as function
// triggered after left mouse click event.
func (s *Slot) SetOnLeftClickFunc(f func(s *Slot)) {
	s.onLeftClick = f
}

// SetOnSpecialLeftClickFunc set speicfied function as function
// triggered after special key pressed + left mouse click event.
func (s *Slot) SetOnSpecialLeftClickFunc(f func(s *Slot)) {
	s.onSpecialLeftClick = f
}

// drawBG draw background within current draw area
// with IMDraw.
func (s *Slot) drawBG(t pixel.Target) {
	s.bgDraw.Clear()
	s.bgDraw.Color = s.color
	s.bgDraw.Push(s.DrawArea().Min)
	s.bgDraw.Push(s.DrawArea().Max)
	s.bgDraw.Rectangle(0)
	s.bgDraw.Draw(t)
}
