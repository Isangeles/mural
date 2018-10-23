/*
 * textedit.go
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

	"golang.org/x/image/colornames"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
)

// Struct for text edit fields.
type Textedit struct {
	bg         *imdraw.IMDraw
	drawArea   pixel.Rect
	color      color.Color
	colorFocus color.Color
	label      *text.Text
	input      *text.Text
	text       string
	focused    bool
	disabled   bool
	onInput    func(t *Textedit)
}

// NewTextecit creates new instance of textedit with specified
// font size, background color and optional label(empty string == no label).
func NewTextedit(fontSize Size, color color.Color, label string) *Textedit {
	t := new(Textedit)
	// Background.
	t.bg = imdraw.New(nil)
	t.color = color
	t.colorFocus = colornames.Crimson
	// Label & Text input.
	font := MainFont(fontSize)
	atlas := Atlas(&font)
	if len(label) > 0 {	
		t.label = text.New(pixel.V(0, 0), atlas)
		t.label.WriteString(label)
	}
	t.input = text.New(pixel.V(0, 0), atlas)

	return t
}

// Draw draws text edit.
func (te *Textedit) Draw(drawArea pixel.Rect, t pixel.Target) {
	// Background & label.
	if te.label != nil {
		te.label.Draw(t, pixel.IM.Moved(drawArea.Min))
		te.drawArea = pixel.R(drawArea.Min.X,
			drawArea.Min.Y - te.label.Bounds().H(),
			drawArea.Max.X, drawArea.Min.Y)
	} else {	
		te.drawArea = drawArea
	}
	color := te.color
	if te.Focused() {
		color = te.colorFocus
	}
	te.drawIMBackground(t, color)
	// Text input.
	te.input.Draw(t, pixel.IM.Moved(te.drawArea.Min))
}

// Update updates text edit.
func (te *Textedit) Update(win *pixelgl.Window) {
	if te.Disabled() {
		return
	}
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		if te.DrawArea().Contains(win.MousePosition()) {
			te.Focus(true)
		} else {
			te.Focus(false)
		}
	}
	if te.focused {
		if win.JustPressed(pixelgl.KeyEnter) {
			if te.onInput != nil {	
				te.onInput(te)
			}
		}
		if win.JustPressed(pixelgl.KeyBackspace) {
			if len(te.text) > 0 {
		 		te.text = te.text[:len(te.text)-1]
				te.input.Clear()
				te.input.WriteString(te.text)
			}
		}
		te.text += win.Typed()
		te.input.WriteString(win.Typed())
	}
}

// Focus sets or removes focus from text edit.
func (te *Textedit) Focus(focus bool) {
	te.focused = focus
}

// Focused checks whether text edit is focused.
func (te *Textedit) Focused() bool {
	return te.focused
}

// Active toggles field activity.
func (te *Textedit) Active(active bool) {
	te.disabled = !active
}

// Disabled checks whether field is disabled.
func (te *Textedit) Disabled() bool {
	return te.disabled
}

// Clear clears text edit input.
func (te *Textedit) Clear() {
	te.input.Clear()
	te.text = ""
}

// Text return current value of text edit.
func (te *Textedit) Text() string {
	return te.text
}

// SetText sets specified text as current value of
// text edit field.
func (te *Textedit) SetText(text string) {
	te.text = text
}

// DrawArea returns current draw area rectangle.
func (te *Textedit) DrawArea() pixel.Rect {
	return te.drawArea
}

// SetOnInputFunc sets callback function triggered after
// input in text edit was accepted(i.e. enter key was pressed).
func (te *Textedit) SetOnInputFunc(f func(t *Textedit)) {
	te.onInput = f
}

// drawIMBackground draws IMDraw background in size of current draw area
// on specified target and with specified color.
func (te *Textedit) drawIMBackground(t pixel.Target, color color.Color) {
	te.bg.Clear()
	te.bg.Color = pixel.ToRGBA(color)
	te.bg.Push(te.drawArea.Min)
	te.bg.Color = pixel.ToRGBA(color)
	te.bg.Push(te.drawArea.Max)
	te.bg.Rectangle(0)
	te.bg.Draw(t)
}
