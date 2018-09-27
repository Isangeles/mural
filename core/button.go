/*
 * button.go
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

package core

import (
	"fmt"
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/faiface/pixel/imdraw"

	"github.com/isangeles/mural/core/data"
)

// Button struct for UI button.
type Button struct {
	bg       *pixel.Sprite
	bgDraw   *imdraw.IMDraw
	label    *text.Text
	size     Size
	pressed  bool
	drawArea pixel.Rect // updated on each draw
}

// NewButton returns new instance of button with specified
// background image and label text.
func NewButton(bgPic pixel.Picture, labelText string) *Button {
	button := new(Button)
	// Backround.
	bg := pixel.NewSprite(bgPic, bgPic.Bounds())
	button.bg = bg
	// Label.
	font := data.MainFontSmall()
	atlas := text.NewAtlas(font, text.ASCII)
	button.label = text.New(pixel.V(0, 0), atlas)
	fmt.Fprint(button.label, labelText)

	return button
}

func NewButtonDraw(size Size, labelText string) *Button {
	button := new(Button)
	// Background.
	button.bgDraw = imdraw.New(nil)
	button.size = size
	// Label.
	font := data.MainFontSmall()
	atlas := text.NewAtlas(font, text.ASCII)
	button.label = text.New(pixel.V(0, 0), atlas)
	fmt.Fprint(button.label, labelText)

	return button
}

// Draw draws button.
func (b *Button) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Calculating draw area.
	// (there should be some more elegant way)
	bgBottomX := matrix[4] - (b.Frame().Size().X / 2)
	bgBottomY := matrix[5] - (b.Frame().Size().Y / 2)
	b.drawArea.Min = pixel.V(bgBottomX, bgBottomY)
	b.drawArea.Max = b.drawArea.Min.Add(b.Frame().Size())
	fmt.Printf("\nb_draw_area:%s\n", b.drawArea) // test
	// Drawing background.
	if b.pressed {
		if b.bg != nil {
			b.bg.DrawColorMask(t, matrix, colornames.Gray)
		} else {
			b.drawIMBackground(t, colornames.Grey, matrix)
		}
	} else {
		if b.bg != nil {
			b.bg.Draw(t, matrix)
		} else {
			b.drawIMBackground(t, colornames.Red, matrix)
		}
	}
	// Drawing label.
	if b.label != nil {
		b.label.Draw(t, matrix)
	}
}

// Draws button background with IMDraw.
func (b *Button) drawIMBackground(t pixel.Target, color color.Color, matrix pixel.Matrix) {
	b.bgDraw.Color = pixel.ToRGBA(color)
	b.bgDraw.Push(b.drawArea.Min)
	b.bgDraw.Color = pixel.ToRGBA(color)
	b.bgDraw.Push(b.drawArea.Max)
	b.bgDraw.Rectangle(0)
	b.bgDraw.Draw(t)
}

// Update updates button.
func (b *Button) Update(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		if b.drawArea.Contains(win.MousePosition()) {
			b.pressed = true
		}
	}
	if win.JustReleased(pixelgl.MouseButtonLeft) {
		b.pressed = false
	}
}

// DrawArea returns button background position and size.
func (b *Button) DrawArea() pixel.Rect {
	return b.drawArea
}

// Frame returns button background size.
func (b *Button) Frame() pixel.Rect {
	if b.bg != nil {
		return b.bg.Frame()
	} else {
		return b.size.ButtonSize()
	}
}

// ContainsPosition checks if specified position is
// within current draw area.
func (b *Button) ContainsPosition(pos pixel.Vec) bool {
	return b.drawArea.Contains(pos)
}
