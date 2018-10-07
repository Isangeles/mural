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
	values                  []string
	onChange                func(s *Switch)
}

// NewSwitch return new instance of switch with IMDraw
// background with specified values to switch.
func NewSwitch(size Size, color color.Color, label string, values []string) *Switch {
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
	atlas := text.NewAtlas(font, text.ASCII)
	s.label = text.New(pixel.V(0, 0), atlas)
	fmt.Fprint(s.label, label)
	// Values.
	s.values = values
	s.index = 0
	s.valueText = text.New(pixel.V(0,0), atlas)
	s.updateValueText()
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
func (s *Switch) Value() string {
	return s.values[s.index]
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
	valueMariginX := (-s.valueText.BoundsOf(s.Value()).Max.X) / 2
	s.valueText.Orig.X = valueMariginX
	s.valueText.Clear()
	fmt.Fprintf(s.valueText, s.Value())
}

// Triggered after next button clicked.
func (s *Switch) onNextButtonClicked(b *Button) {
	if s.index < len(s.values)-1 {
		s.index++
		s.updateValueText()
	} else {
		s.index = 0
		s.updateValueText()
	}
	if s.onChange != nil {
		s.onChange(s)
	}
}

// Triggered after prev button clicked.
func (s *Switch) onPrevButtonClicked(b *Button) {
	if s.index > 0 {
		s.index--
		s.updateValueText()
	} else {
		s.index = len(s.values)-1
		s.updateValueText()
	}
	if s.onChange != nil {
		s.onChange(s)
	}
}
