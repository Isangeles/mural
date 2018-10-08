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

// Button struct for UI button.
type Button struct {
	bgSpr     *pixel.Sprite
	bgDraw    *imdraw.IMDraw
	label     *text.Text
	size      Size
	shape     Shape
	color     color.Color
	colorPush color.Color
	pressed   bool
	drawArea  pixel.Rect // updated on each draw
	onClick   func(b *Button)
}

// NewButton creates new instance of button with specified size, color and
// label text.
func NewButton(size Size, shape Shape, color color.Color, labelText string) *Button {
	button := new(Button)
	// Background.
	button.bgDraw = imdraw.New(nil)
	button.size = size
	button.shape = shape
	button.color = color
	button.colorPush = colornames.Grey
	// Label.
	font := MainFont(button.size)
	atlas := text.NewAtlas(font, text.ASCII)
	button.label = text.New(pixel.V(0, 0), atlas)
	labelMariginX := (-button.label.BoundsOf(labelText).Max.X) / 2
	button.label.Orig = pixel.V(labelMariginX, 0)
	button.label.Clear()
	fmt.Fprint(button.label, labelText)

	return button
}

// NewButtonSprite creates new instance of button with specified
// background image and label text.
func NewButtonSprite(bgPic pixel.Picture, labelText string) *Button {
	button := new(Button)
	// Backround.
	bg := pixel.NewSprite(bgPic, bgPic.Bounds())
	button.bgSpr = bg
	// Label.
	font := MainFont(SIZE_SMALL)
	atlas := text.NewAtlas(font, text.ASCII)
	button.label = text.New(pixel.V(0, 0), atlas)
	labelMarigin := (-button.label.BoundsOf(labelText).Max.X) / 2
	button.label.Orig = pixel.V(labelMarigin, 0)
	button.label.Clear()
	fmt.Fprint(button.label, labelText)

	return button
}

// Draw draws button.
func (b *Button) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Calculating draw area.
	b.drawArea = MatrixToDrawArea(matrix, b.Frame())
	// Drawing background.
	if b.pressed {
		if b.bgSpr != nil {
			b.bgSpr.DrawColorMask(t, matrix, colornames.Gray)
		} else {
			b.drawIMBackground(t, b.colorPush)
		}
	} else {
		if b.bgSpr != nil {
			b.bgSpr.Draw(t, matrix)
		} else {
			b.drawIMBackground(t, b.color)
		}
	}
	// Drawing label.
	if b.label != nil {
		b.label.Draw(t, matrix)
	}
}

// Update updates button.
func (b *Button) Update(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		if b.ContainsPosition(win.MousePosition()) {
			b.pressed = true
		}
	}
	if win.JustReleased(pixelgl.MouseButtonLeft) {
		if b.pressed && b.ContainsPosition(win.MousePosition()) {
			if b.onClick != nil {
				b.onClick(b)
			}
		}
		b.pressed = false
	}
}

// Draws button background with IMDraw.
func (b *Button) drawIMBackground(t pixel.Target, color color.Color) {
	b.bgDraw.Color = pixel.ToRGBA(color)
	b.bgDraw.Push(b.drawArea.Min)
	b.bgDraw.Color = pixel.ToRGBA(color)
	b.bgDraw.Push(b.drawArea.Max)
	b.bgDraw.Rectangle(0)
	b.bgDraw.Draw(t)
}

// OnClick sets specified function as on-click
// callback function.
func (b *Button) SetOnClickFunc(callback func(b *Button)) {
	b.onClick = callback
}

// DrawArea returns current button background position and size.
func (b *Button) DrawArea() pixel.Rect {
	return b.drawArea
}

// Frame returns button background size, in form
// of rectangle.
func (b *Button) Frame() pixel.Rect {
	if b.bgSpr != nil {
		return b.bgSpr.Frame()
	} else {
		return b.size.ButtonSize(b.shape)
	}
}

// ContainsPosition checks if specified position is
// within current draw area.
func (b *Button) ContainsPosition(pos pixel.Vec) bool {
	return b.drawArea.Contains(pos)
}
