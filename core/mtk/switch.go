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
	//"fmt"
	"image/color"
	
	"github.com/faiface/pixel"
	//"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/faiface/pixel/imdraw"
)

// Swtch struct represents value switch control.
type Switch struct {
	bgDraw        *imdraw.IMDraw
	bgSpr         *pixel.Sprite
	prevB, nextB  *Button
	label         *text.Text
	drawArea      pixel.Rect // updated on each draw
	size          Size
	color         color.Color
	value         string
	values        []string
}

// NewSwitch return new instance of switch with IMDraw
// background with specified values to switch.
func NewSwitch(size Size, color color.Color, values []string) *Switch {
	s := new(Switch)
	// Background.
	s.bgDraw = imdraw.New(nil)
	s.size = size
	s.color = color
	// TODO: buttons, values, label.
	return s 
}

// Draw draws switch.
func (s *Switch) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Calculating draw area.
	s.drawArea = MatrixToDrawArea(matrix, s.Frame())
	// Background.
	if s.bgSpr == nil {
		s.bgSpr.Draw(t, matrix)
	} else {
		s.drawIMBackground(t)
	}
}

// Update updates switch and all elements.
func (s *Switch) Update(win *pixel.Window) {
	// TODO: update buttons.
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
	if s.bgSpr == nil {
		return s.bgSpr.Frame()
	} else {
		return s.size.SwitchSize()
	}
}
