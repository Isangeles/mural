/*
 * button.go
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
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/faiface/beep"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"
)

// Button struct for UI button.
type Button struct {
	bgSpr      *pixel.Sprite
	bgDraw     *imdraw.IMDraw
	label      *Text
	info       *InfoWindow
	size       Size
	shape      Shape
	color      color.Color
	colorPush  color.Color
	colorHover color.Color
	pressed    bool
	focused    bool
	hovered    bool
	disabled   bool
	drawArea   pixel.Rect // updated on each draw
	onClick    func(b *Button)
	clickSound *beep.Buffer
}

// NewButton creates new instance of button with specified size, color and
// label text.
func NewButton(size Size, shape Shape, color color.Color,
	labelText, infoText string) *Button {
	button := new(Button)
	// Background.
	button.bgDraw = imdraw.New(nil)
	button.size = size
	button.shape = shape
	button.color = color
	button.colorPush = colornames.Grey
	button.colorHover = colornames.Crimson
	// Label.
	button.label = NewText(labelText, size, button.Frame().W()) 
	// Info window.
	if len(infoText) > 0 {	
		button.info = NewInfoWindow(infoText)
	}

	return button
}

// NewButtonSprite creates new instance of button with specified
// background image and label text.
func NewButtonSprite(bgPic pixel.Picture, fontSize Size, labelText,
	infoText string) *Button {
	button := new(Button)
	// Backround.
	button.bgSpr = pixel.NewSprite(bgPic, bgPic.Bounds())
	button.colorPush = colornames.Grey
	button.colorHover = colornames.Crimson
	// Label.
	button.label = NewText(labelText, fontSize, button.Frame().W()) 
	// Info window.
	if len(infoText) > 0 {	
		button.info = NewInfoWindow(infoText)
	}

	return button
}

// Draw draws button.
func (b *Button) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Calculating draw area.
	b.drawArea = MatrixToDrawArea(matrix, b.Frame())
	// Drawing background.
	bgColor := b.color
	if b.pressed || b.Disabled() {
		bgColor = b.colorPush
	} else if b.hovered {
		bgColor = b.colorHover
	}
	if b.bgSpr != nil {
		if bgColor == nil {
			b.bgSpr.Draw(t, matrix)
		} else {
			b.bgSpr.DrawColorMask(t, matrix, bgColor)
		}
	} else {
		b.drawIMBackground(t, bgColor)
	}
	// Drawing label.
	if b.label != nil {
		b.label.Draw(t, matrix)
	}
	// Info window.
	if b.info != nil && b.hovered {
		b.info.Draw(t)
	}
}

// Update updates button.
func (b *Button) Update(win *Window) {
	if b.Disabled() {
		return
	}
	// Mouse events.
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		if b.DrawArea().Contains(win.MousePosition()) {
			b.pressed = true
		}
	}
	if win.JustReleased(pixelgl.MouseButtonLeft) {
		if b.pressed && b.DrawArea().Contains(win.MousePosition()) {
			if b.onClick != nil {
				b.onClick(b)
			}
			if Audio() != nil && b.clickSound != nil {
				Audio().Play(b.clickSound)
			}
		}
		b.pressed = false
	}
	// On-hover.
	if b.DrawArea().Contains(win.MousePosition()) || b.Focused() {
		b.hovered = true
		if b.info != nil {	
			b.info.Update(win)
		}
	} else {
		b.hovered = false
	}
	// On-focus events.
	if b.Focused() {
		if win.JustPressed(pixelgl.KeyEnter) {
			b.onClick(b)
		}
	}
}

// Focus sets/removes focus from button
func (b *Button) Focus(focus bool) {
	b.focused = focus
}

// Focused checks whether buttons is focused.
func (b *Button) Focused() bool {
	return b.focused
}

// Active toggles button active state.
func (b *Button) Active(active  bool) {
	b.disabled = !active
}

// Disabled checks whether button is disabled.
func (b *Button) Disabled() bool {
	return b.disabled
}

// SetOnClickFunc sets specified function as on-click
// callback function.
func (b *Button) SetOnClickFunc(callback func(b *Button)) {
	b.onClick = callback
}

// SetClickSound sets specified audio buffer as
// on-click audio effect.
func (b *Button) SetClickSound(s *beep.Buffer) {
	b.clickSound = s
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

// Draws button background with IMDraw.
func (b *Button) drawIMBackground(t pixel.Target, color color.Color) {
	b.bgDraw.Clear()
	b.bgDraw.Color = pixel.ToRGBA(color)
	b.bgDraw.Push(b.DrawArea().Min)
	b.bgDraw.Push(b.DrawArea().Max)
	b.bgDraw.Rectangle(0)
	b.bgDraw.Draw(t)
}
