/*
 * switch.go
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
	"fmt"
	"image/color"

	"golang.org/x/image/colornames"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/faiface/pixel/imdraw"
)

// Plain string for regular string to satisfy stringer
// interface.
type SwitchValue struct {
	Label string
	Value interface{}
}

// String returns label of switch value.
func (s SwitchValue) String() string {
	return s.Label
}

// Switch struct represents graphical switch for values.
type Switch struct {
	bgDraw                  *imdraw.IMDraw
	bgSpr                   *pixel.Sprite
	prevButton, nextButton  *Button
	valueText               *text.Text
	label                   *text.Text
	drawArea                pixel.Rect // updated on each draw
	size                    Size
	color                   color.Color
	index                   int
	values                  []SwitchValue
	onChange                func(s *Switch)
}

// NewSwitch creates new instance of switch with IMDraw
// background with specified values to switch.
func NewSwitch(size Size, color color.Color, label string, values []SwitchValue) *Switch {
	s := new(Switch)
	// Background.
	s.bgDraw = imdraw.New(nil)
	s.size = size
	s.color = color
	// Buttons.
	s.prevButton = NewButton(size-2, SHAPE_SQUARE, colornames.Red, "-")
	s.nextButton = NewButton(size-2, SHAPE_SQUARE, colornames.Red, "+")
	s.prevButton.SetOnClickFunc(s.onPrevButtonClicked)
	s.nextButton.SetOnClickFunc(s.onNextButtonClicked)
	// Label.
	font := MainFont(s.size)
	atlas := Atlas(&font)
	s.label = text.New(pixel.V(0, 0), atlas)
	fmt.Fprint(s.label, label)
	// Values.
	s.values = values
	s.index = 0
	s.valueText = text.New(pixel.V(0,0), atlas)
	s.updateValueText()
	return s 
}

// NewStringSwitch creates new instance of switch with IMDraw
// background with specified string values to switch.
func NewStringSwitch(size Size, color color.Color, label string, values []string) *Switch {
	// All string values to switchString helper struct.
	strValues := make([]SwitchValue, len(values))
	for i, v := range values {
		ss := SwitchValue{Label: v, Value: v}
		strValues[i] = ss
	}
	
	s := NewSwitch(size, color, label, strValues)
	return s
}

// Draw draws switch.
func (s *Switch) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Calculating draw area.
	s.drawArea = MatrixToDrawArea(matrix, s.Frame())
	// Background.
	if s.bgSpr != nil {
		s.bgSpr.Draw(t, matrix)
	} else {
		s.drawIMBackground(t)
	}
	// Value & label.
	s.valueText.Draw(t, matrix)
	s.label.Draw(t, pixel.IM.Moved(PosBL(s.label.Bounds(), s.drawArea.Min)))
	// Buttons.
	s.prevButton.Draw(t, pixel.IM.Moved(DisTL(s.drawArea, 0.03)))
	s.nextButton.Draw(t, pixel.IM.Moved(DisTR(s.drawArea, 0.03)))
}

// Update updates switch and all elements.
func (s *Switch) Update(win *pixelgl.Window) {
	s.prevButton.Update(win)
	s.nextButton.Update(win)
}

// drawIMBackground Draws IMDraw background.
func (s *Switch) drawIMBackground(t pixel.Target) {
	s.bgDraw.Color = pixel.ToRGBA(s.color)
	s.bgDraw.Push(s.drawArea.Min)
	s.bgDraw.Color = pixel.ToRGBA(s.color)
	s.bgDraw.Push(s.drawArea.Max)
	s.bgDraw.Rectangle(0)
	s.bgDraw.Draw(t)
}

// Frame returns switch background size, in form
// of rectangle.
func (s *Switch) Frame() pixel.Rect {
	if s.bgSpr != nil {
		return s.bgSpr.Frame()
	} else {
		return s.size.SwitchSize()
	}
}

// DrawArea returns current switch background position and size.
func (s *Switch) DrawArea() pixel.Rect {
	return s.drawArea
}

// Value returns current switch value.
func (s *Switch) Value() SwitchValue {
	return s.values[s.index]
}

// Find if switch constains specified value and returns
// index of this value or -1 if switch does not contains
// such value.
func (s *Switch) Find(value interface{}) int {
	for i, v := range s.values {
		if value == v.Value {
			return i
		}
	}
	return -1
}

// SetIndex sets value with specified index as current value
// of this switch. If specified value is bigger than maximal
// possible index then index of last value is set, if specified
// index is smaller than minimal then index of first value is set. 
func (s *Switch) SetIndex(index int) {
	if index > len(s.values)-1 {
		s.index = len(s.values)-1
	} else if index < 0 {
		s.index = 0
	} else {
		s.index = index
	}
	s.updateValueText()
}

// Sets specified function as function triggered on on switch value change.
func (s *Switch) SetOnChangeFunc(f func(s *Switch)) {
	s.onChange = f
}

// updateValueText updates text with current switch value.
func (s *Switch) updateValueText() {
	valueMariginX := (-s.valueText.BoundsOf(s.Value().String()).Max.X) / 2
	s.valueText.Orig.X = valueMariginX
	s.valueText.Clear()
	fmt.Fprintf(s.valueText, s.Value().String())
}

// Triggered after next button clicked.
func (s *Switch) onNextButtonClicked(b *Button) {
	s.SetIndex(s.index+1)
	if s.onChange != nil {
		s.onChange(s)
	}
}

// Triggered after prev button clicked.
func (s *Switch) onPrevButtonClicked(b *Button) {
	s.SetIndex(s.index-1)
	if s.onChange != nil {
		s.onChange(s)
	}
}
