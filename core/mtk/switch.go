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
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

// Tuple for switch values, contains value to display and
// real value.
type SwitchValue struct {
	View  interface{}
	Value interface{}
}

// Label returns string representation of switch value.
func (s SwitchValue) Label() string {
	switch v := s.View.(type) {
	case string:
		return v
	default:
		return "none"
	}
}

// Sprite returns graphical representation of switch value.
func (s SwitchValue) Sprite() (*pixel.Sprite, error) {
	spr, ok := s.View.(*pixel.Sprite)
	if !ok {
		return nil, fmt.Errorf("fail_to_retrieve_view_sprite")
	}
	return spr, nil
}

// Switch struct represents graphical switch for values.
type Switch struct {
	bgDraw                 *imdraw.IMDraw
	bgSpr                  *pixel.Sprite
	prevButton, nextButton *Button
	valueText              *Text
	valueSprite            *pixel.Sprite
	label                  *Text
	drawArea               pixel.Rect // updated on each draw
	size                   Size
	color                  color.Color
	index                  int
	focused                bool
	disabled               bool
	values                 []SwitchValue
	onChange               func(s *Switch, old, new *SwitchValue)
}

// NewSwitch creates new instance of switch with IMDraw
// background with specified values to switch.
func NewSwitch(size Size, color color.Color, label string,
	values []SwitchValue) *Switch {
	s := new(Switch)
	// Background.
	s.bgDraw = imdraw.New(nil)
	s.size = size
	s.color = color
	// Buttons.
	s.prevButton = NewButton(size-2, SHAPE_SQUARE, colornames.Red, "-", "")
	s.nextButton = NewButton(size-2, SHAPE_SQUARE, colornames.Red, "+", "")
	s.prevButton.SetOnClickFunc(s.onPrevButtonClicked)
	s.nextButton.SetOnClickFunc(s.onNextButtonClicked)
	// Label.
	s.label = NewText(label, size-1, s.Frame().W())
	// Values.
	s.values = values
	s.index = 0
	s.valueText = NewText("", size, 100)
	s.updateValueView()
	return s
}

// NewStringSwitch creates new instance of switch with IMDraw
// background with specified string values to switch.
func NewStringSwitch(size Size, color color.Color, label string,
	values []string) *Switch {
	// All string values to switchString helper struct.
	strValues := make([]SwitchValue, len(values))
	for i, v := range values {
		ss := SwitchValue{v, v}
		strValues[i] = ss
	}

	s := NewSwitch(size, color, label, strValues)
	return s
}

// NewIntSwitch creates new instance of switch with IMDraw
// background and with specified int values to switch.
func NewIntSwitch(size Size, color color.Color, label string,
	min, max int) *Switch {
	// All int values from specified min max range.
	length := max - min + 1
	intValues := make([]SwitchValue, length)
	for i := min; i <= max; i++ {
		value := i //+ 1
		intValues[i] = SwitchValue{fmt.Sprint(value), value}
	}

	s := NewSwitch(size, color, label, intValues)
	return s
}

// NewPictureSwitch creates new switch with IMDraw background
// and with specified pictures as switch values.
func NewPictureSwitch(size Size, color color.Color, label string,
	pics map[string]pixel.Picture) *Switch {
	var picValues []SwitchValue
	for name, pic := range pics {
		spr := pixel.NewSprite(pic, pic.Bounds())
		val := SwitchValue{spr, name}
		picValues = append(picValues, val)
	}

	s := NewSwitch(size, color, label, picValues)
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
	// Value view.
	valueDA := s.valueText.DrawArea()
	if s.valueSprite == nil {
		s.valueText.Draw(t, matrix)
	} else {
		s.valueSprite.Draw(t, matrix.Scaled(s.DrawArea().Center(), Scale()))
		valueDA = MatrixToDrawArea(matrix.Scaled(s.DrawArea().Center(),
			Scale()), s.valueSprite.Frame())
	}
	s.label.Draw(t, pixel.IM.Moved(PosBL(s.label.Bounds(), s.drawArea.Min)))
	// Buttons.
	s.prevButton.Draw(t, pixel.IM.Moved(LeftOf(valueDA, s.prevButton.Frame(),
		10)))
	s.nextButton.Draw(t, pixel.IM.Moved(RightOf(valueDA, s.nextButton.Frame(),
		10)))
}

// Update updates switch and all elements.
func (s *Switch) Update(win *pixelgl.Window) {
	if s.Disabled() {
		return
	}
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

// SetValues sets specified list with values as switch values.
func (s *Switch) SetValues(values []SwitchValue) {
	s.values = values
}

// SetIntValue sets all integer values from specified range as
// switch values.
func (s *Switch) SetIntValues(min, max int) {
	intValues := make([]SwitchValue, max)
	for i := min; i < max; i++ {
		value := i + 1
		intVal := SwitchValue{fmt.Sprint(value), value}
		intValues[i] = intVal
	}
	s.SetValues(intValues)
}

// Focus toggles focus on element.
func (s *Switch) Focus(focus bool) {
	s.focused = focus
}

// Focused checks whether switch is focused.
func (s *Switch) Focused() bool {
	return s.focused
}

// Active toggles switch activity.
func (s *Switch) Active(active bool) {
	s.prevButton.Active(active)
	s.nextButton.Active(active)
	s.disabled = !active
}

// Disabled checks whether switch is active.
func (s *Switch) Disabled() bool {
	return s.disabled
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
func (s *Switch) Value() *SwitchValue {
	return &s.values[s.index]
}

// Reset sets value with first index as current
// value.
func (s *Switch) Reset() {
	s.SetIndex(0)
}

// Find checks if switch constains specified value and returns
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

// Find searches switch values for value with specified index
// and returns this value or nil if switch does not contains
// value with such index.
func (s *Switch) FindValue(index int) *SwitchValue {
	return &s.values[index]
}

// SetIndex sets value with specified index as current value
// of this switch. If specified value is bigger than maximal
// possible index, then index of first value is set, if specified
// index is smaller than minimal, then index of last value is set.
func (s *Switch) SetIndex(index int) {
	if index > len(s.values)-1 {
		s.index = 0
	} else if index < 0 {
		s.index = len(s.values) - 1
	} else {
		s.index = index
	}
	s.updateValueView()
}

// Sets specified function as function triggered on on switch value change.
func (s *Switch) SetOnChangeFunc(f func(s *Switch, old, new *SwitchValue)) {
	s.onChange = f
}

// updateValueView updates value view with current switch value.
func (s *Switch) updateValueView() {
	if spr, err := s.Value().Sprite(); err != nil {
		s.valueText.SetText(s.Value().Label())
	} else {
		s.valueSprite = spr
	}
}

// Triggered after next button clicked.
func (s *Switch) onNextButtonClicked(b *Button) {
	oldIndex := s.index
	s.SetIndex(s.index + 1)
	if s.onChange != nil {
		oldValue := s.FindValue(oldIndex)
		s.onChange(s, oldValue, s.Value())
	}
}

// Triggered after prev button clicked.
func (s *Switch) onPrevButtonClicked(b *Button) {
	oldIndex := s.index
	s.SetIndex(s.index - 1)
	if s.onChange != nil {
		s.onChange(s, s.FindValue(oldIndex), s.Value())
	}
}
