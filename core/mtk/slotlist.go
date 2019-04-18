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

	"golang.org/x/image/colornames"
                                       
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
	sl.upButton = NewButton(SIZE_MINI, SHAPE_SQUARE, colornames.Red)
	sl.upButton.SetLabel("+")
	sl.upButton.SetOnClickFunc(sl.onUpButtonClicked)
	sl.downButton = NewButton(SIZE_MINI, SHAPE_SQUARE, colornames.Red)
	sl.downButton.SetLabel("-")
	sl.downButton.SetOnClickFunc(sl.onDownButtonClicked)
	// Calculating amount of slots based on background and
	// buttons sizes.
	slotBounds := slotSize.SlotSize()
	bgWidth := sl.Size().X - sl.upButton.Size().X
	sl.spl = int(bgWidth/(slotBounds.W() + ConvSize(2)))
	sl.lines = int(sl.Size().Y / (slotBounds.H() + ConvSize(2))) - 1
	return sl
}

// Draw draws list.
func (sl *SlotList) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Draw area.
	sl.drawArea = MatrixToDrawArea(matrix, sl.bgSize)
	// Background.
	if sl.bgSpr != nil {
		sl.bgSpr.Draw(t, matrix)
	} else {
		DrawRectangle(t, sl.drawArea, sl.bgColor)
	}
	// Buttons.
	upButtonPos := MoveTR(sl.Size(), sl.upButton.Size())
	downButtonPos := MoveBR(sl.Size(), sl.downButton.Size())
	sl.upButton.Draw(t, matrix.Moved(upButtonPos))
	sl.downButton.Draw(t, matrix.Moved(downButtonPos))
	// Slots.
	// TODO: too slow.
	if len(sl.slots) < 1 {
		return
	}
	slotStart := pixel.V(-sl.Size().X/2+sl.slots[0].Size().X/2,
		sl.Size().Y/2-sl.slots[0].Size().Y/2)
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
		slotMove.X += s.Size().X + ConvSize(2)
		splCount += 1
		if splCount >= sl.spl {
			slotMove.X = slotStart.X
			slotMove.Y -= s.Size().Y + ConvSize(2)
			splCount = 0
			lineCount += 1
		}
		if lineCount > sl.lines {
			break
		}
	}
}

// SetUpButtonBackground sets specified sprite as scroll up button
// background.
func (sl *SlotList) SetUpButtonBackground(s *pixel.Sprite) {
	sl.upButton.SetBackground(s)
	sl.upButton.SetColor(nil)
}

// SetDownButtonBackground sets specified sprite as scroll down button
// background.
func (sl *SlotList) SetDownButtonBackground(s *pixel.Sprite) {
	sl.downButton.SetBackground(s)
	sl.downButton.SetColor(nil)
}

// Update updates list.
func (sl *SlotList) Update(win *Window) {
	// Buttons.
	sl.upButton.Update(win)
	sl.downButton.Update(win)
	// Slots.
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

// EmptySlot returns first empty slot.
func (sl *SlotList) EmptySlot() *Slot {
	for _, s := range sl.slots {
		if len(s.Values()) < 1 {
			return s
		}
	}
	return nil
}

// Clear clears all slots on the list.
func (sl *SlotList) Clear() {
	for _, s := range sl.slots {
		s.Clear()
	}
}

// Bounds retruns background size.
func (sl *SlotList) Size() pixel.Vec {
	if sl.bgSpr == nil {
		return sl.bgSize
	}
	return sl.bgSpr.Frame().Size()
}

// setStartLine sets specified line ID as current
// line ID.
func (sl *SlotList) setStartLine(line int) {
	if line*sl.spl > len(sl.slots) || line*sl.spl < 0 {
		return
	}
	sl.lineID = line
}

// Triggered after up button clicked.
func (sl *SlotList) onUpButtonClicked(b *Button) {
	sl.setStartLine(sl.lineID - 1)
}

// Triggered after down button clicked.
func (sl *SlotList) onDownButtonClicked(b *Button) {
	sl.setStartLine(sl.lineID + 1)
}
