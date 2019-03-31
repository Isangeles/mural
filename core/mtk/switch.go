/*
 * switch.go
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
	"fmt"
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
)

// Tuple for switch values, contains value to
// display(view) and real value.
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
func (s SwitchValue) Picture() (pixel.Picture, error) {
	pic, ok := s.View.(pixel.Picture)
	if !ok {
		return nil, fmt.Errorf("fail_to_retrieve_view_picture")
	}
	return pic, nil
}

// IntValue returns value as integer, or error if not
// possible.
func (s SwitchValue) IntValue() (int, error) {
	num, ok := s.Value.(int)
	if !ok {
		return 0, fmt.Errorf("fail_to_retrieve_switch_integer_value")		
	}
	return num, nil
}

// TextValue returns value as text, or error if not
// possible.
func (s SwitchValue) TextValue() (string, error) {
	txt, ok := s.Value.(string)
	if !ok {
		return "", fmt.Errorf("fail_to_retrieve_switch_text_value")
	}
	return txt, nil
}

// Switch struct represents graphical switch for values.
type Switch struct {
	bgSpr                  *pixel.Sprite
	prevButton, nextButton *Button
	valueText              *Text
	valueSprite            *pixel.Sprite
	label                  *Text
	info                   *InfoWindow
	drawArea               pixel.Rect // updated on each draw
	size                   Size
	color                  color.Color
	index                  int
	focused                bool
	disabled               bool
	hovered                bool
	values                 []SwitchValue
	onChange               func(s *Switch, old, new *SwitchValue)
}

// NewSwitch creates new instance of switch with IMDraw
// background with specified values to switch.
func NewSwitch(size Size, color color.Color, label, info string,
	values []SwitchValue) *Switch {
	s := new(Switch)
	// Background.
	s.size = size
	s.color = color
	// Buttons.
	s.prevButton = NewButton(s.size-2, SHAPE_SQUARE, colornames.Red, "-", "")
	s.prevButton.SetOnClickFunc(s.onPrevButtonClicked)
	s.nextButton = NewButton(s.size-2, SHAPE_SQUARE, colornames.Red, "+", "")
	s.nextButton.SetOnClickFunc(s.onNextButtonClicked)
	// Label & info.
	s.label = NewText(s.size-1, s.Bounds().W())
	s.label.JustCenter()
	s.label.SetText(label)
	if len(info) > 0 {
		s.info = NewInfoWindow(SIZE_SMALL, colornames.Grey)
		s.info.Add(info)
	}
	// Values.
	s.values = values
	s.index = 0
	s.valueText = NewText(s.size, 100)
	s.updateValueView()

	return s
}

// Draw draws switch.
func (s *Switch) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Calculating draw area.
	s.drawArea = MatrixToDrawArea(matrix, s.Bounds())
	// Background.
	if s.bgSpr != nil {
		s.bgSpr.Draw(t, matrix)
	} else {
		DrawRectangle(t, s.DrawArea(), s.color)
	}
	// Value view.
	valueDA := s.valueText.DrawArea()
	if s.valueSprite == nil {
		s.valueText.Draw(t, matrix)
	} else {
		s.valueSprite.Draw(t, matrix)
		valueDA = MatrixToDrawArea(matrix, s.valueSprite.Frame())
	}
	// Label & info window.
	labelPos := MoveBC(s.Bounds(), s.label.Bounds().Max)
	s.label.Draw(t, matrix.Moved(labelPos))
	if s.info != nil && s.hovered {
		s.info.Draw(t)
	}
	// Buttons.
	prevButtonPos := LeftOf(valueDA, s.prevButton.Frame(), 10)
	nextButtonPos := RightOf(valueDA, s.nextButton.Frame(), 10)
	s.prevButton.Draw(t, Matrix().Moved(prevButtonPos))
	s.nextButton.Draw(t, Matrix().Moved(nextButtonPos))
}

// Update updates switch and all elements.
func (s *Switch) Update(win *Window) {
	if s.Disabled() {
		return
	}
	// Mouse events.
	if s.DrawArea().Contains(win.MousePosition()) {
		s.hovered = true
		if s.info != nil {	
			s.info.Update(win)
		}
	} else {
		s.hovered = false
	}
	// Elements update.
	s.prevButton.Update(win)
	s.nextButton.Update(win)
}

// SetBackground sets specified sprite as switch
// background, also removes background color.
func (s *Switch) SetBackground(spr *pixel.Sprite) {
	s.bgSpr = spr
	s.color = nil
}

// SetColor sets specified color as switch background
// color.
func (s *Switch) SetColor(c color.Color) {
	s.color = c
}

// SetNextButtonBackground sets specified sprite as next
// button background.
func (s *Switch) SetNextButtonBackground(spr *pixel.Sprite) {
	s.nextButton.SetBackground(spr)
}

// SetPrevButtonBackground sets specified sprite as previous
// button background.
func (s *Switch) SetPrevButtonBackground(spr *pixel.Sprite) {
	s.prevButton.SetBackground(spr)
}

// SetValues sets specified list with values as switch values.
func (s *Switch) SetValues(values []SwitchValue) {
	s.values = values
	s.updateValueView()
}

// SetTextValues sets specified textual values as switch
// values
func (s *Switch) SetTextValues(values []string) {
	// All string values to switchString helper struct.
	strValues := make([]SwitchValue, len(values))
	for i, v := range values {
		ss := SwitchValue{v, v}
		strValues[i] = ss
	}
	s.SetValues(strValues)
}

// SetIntValue sets all integer values from specified range as
// switch values.
func (s *Switch) SetIntValues(min, max int) {
	intValues := make([]SwitchValue, max+1)
	for i := min; i < max+1; i++ {
		value := i
		intVal := SwitchValue{fmt.Sprint(value), value}
		intValues[i] = intVal
	}
	s.SetValues(intValues)
}

// SetPictureValues sets specified pictures as switch values.
func (s *Switch) SetPictureValues(pics map[string]pixel.Picture) {
	var picValues []SwitchValue
	for name, pic := range pics {
		val := SwitchValue{pic, name}
		picValues = append(picValues, val)
	}
	s.SetValues(picValues)
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

// Bounds returns switch background size, in form
// of rectangle.
func (s *Switch) Bounds() pixel.Rect {
	if s.bgSpr == nil {
		return s.size.SwitchSize()
	}
	return s.bgSpr.Frame()
}

// DrawArea returns current switch background position and size.
func (s *Switch) DrawArea() pixel.Rect {
	return s.drawArea
}

// Value returns current switch value.
func (s *Switch) Value() *SwitchValue {
	if s.index >= len(s.values) || s.index < 0 {
		return nil
	}
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
	if index >= len(s.values) || index < 0 {
		return nil
	}
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
	if s.values == nil || len(s.values) < 1 {
		return
	}
	if pic, err := s.Value().Picture(); err != nil {
		s.valueText.SetText(s.Value().Label())
	} else {
		s.valueSprite = pixel.NewSprite(pic, pic.Bounds())
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
