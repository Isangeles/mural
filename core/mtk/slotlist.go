/*
 * slotlist.go
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
	
	"github.com/faiface/pixel"
)

// Struct for list with slots.
type SlotList struct {
	bgSpr      *pixel.Sprite
	bgSize     pixel.Vec
	bgColor    color.Color
	drawArea   pixel.Rect
	upButton   *Button
	downButton *Button
	slots      []*Slot
	spl        int // slots per line
	lines      int
	lineID     int // current line
}

// NewSlotList creates new list with slots.
func NewSlotList(bgSize pixel.Vec, bgColor color.Color, slotSize Size) *SlotList {
	sl := new(SlotList)
	sl.bgSize = bgSize
	sl.bgColor = bgColor
	sl.spl = int(bgSize.X / slotSize.SlotSize(SHAPE_SQUARE).W()) - 1
	sl.lines = int(bgSize.Y / slotSize.SlotSize(SHAPE_SQUARE).H())
	return sl
}

// Draw draws list.
func (sl *SlotList) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Draw area.
	sl.drawArea = MatrixToDrawArea(matrix,
		pixel.R(0, 0, sl.bgSize.X, sl.bgSize.Y))
	// Background.
	if sl.bgSpr != nil {
		sl.bgSpr.Draw(t, matrix)
	} else {
		DrawRectangle(t, sl.drawArea, sl.bgColor)
	}
	// Slots.
	if len(sl.slots) < 1 {
		return
	}
	slotStart := pixel.V(-sl.Bounds().W()/2 + sl.slots[0].Bounds().W()/2,
		sl.Bounds().H()/2 - sl.slots[0].Bounds().H()/2)
	slotMove := slotStart
	splCount := 0
	lineCount := 0
	startSlot := 0
	startSlot = sl.lineID * sl.spl
	for i, s := range sl.slots {
		if i < startSlot {
			continue
		}
		s.Draw(t, matrix.Moved(slotMove))
		slotMove.X += s.Bounds().W() + ConvSize(2)
		splCount += 1
		if splCount >= sl.spl {
			slotMove.X = slotStart.X
			slotMove.Y -= s.Bounds().H() + ConvSize(2)
			splCount = 0
			lineCount += 1
		}
		if lineCount > sl.lines {
			break
		}
	}
}

// Update updates list.
func (sl *SlotList) Update(win *Window) {
	for _, s := range sl.slots {
		s.Update(win)
	}
}

// Add adds specified slot to list.
func (sl *SlotList) Add(s *Slot) {
	sl.slots = append(sl.slots, s)
}

// Slots returns all slots from list.
func (sl *SlotList) Slots() []*Slot {
	return sl.slots
}

// Bounds retruns background size bounds.
func (sl *SlotList) Bounds() pixel.Rect {
	if sl.bgSpr == nil {
		return pixel.R(0, 0, sl.bgSize.X, sl.bgSize.Y)
	}
	return sl.bgSpr.Frame()
}
