/*
 * textbox.go
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
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	//"golang.org/x/image/colornames"
	
	"github.com/isangeles/mural/core/data"
)

// Struct for textboxes
type Textbox struct {
	bg          *imdraw.IMDraw
	textarea    *text.Text
	bgHeight    float64
	textContent []fmt.Stringer
	visibleText []string
	startID     int
}

// NewTextbox creates new textbox.
func NewTextbox() (*Textbox, error) {
	t := new(Textbox)
	
	// Background.
	t.bg = imdraw.New(nil)

	// Text.
	font := data.MainFontNormal()
	atlas := text.NewAtlas(font, text.ASCII)
	t.textarea = text.New(pixel.V(0, 0), atlas)
	
	return t, nil
}

// Draw draws textbox.
func (t *Textbox) Draw(drawMin, drawMax pixel.Vec, win *pixelgl.Window) {
	// Background.
	t.bg.Color = pixel.RGB(0.4, 0.4, 0.4)
	t.bg.Push(drawMin)
	t.bg.Color = pixel.RGB(0.4, 0.4, 0.4)
	t.bg.Push(drawMax)
	t.bg.Rectangle(0)
	t.bg.Draw(win)
	t.bgHeight = drawMax.Y

	// Text content.
	t.textarea.Draw(win, pixel.IM.Moved(pixel.V(win.Bounds().Min.X, win.Bounds().Max.Y - 50)))
}

// Update handles key events.
func (t *Textbox) Update(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyDown) {
		if t.startID < len(t.textContent) - 1 {
			t.startID ++
			t.updateTextVisibility()
		}
	}
	if win.JustPressed(pixelgl.KeyUp) {
		if t.startID > 0 {
			t.startID --
			t.updateTextVisibility()
		}
	}
}

// Clears textbox and inserts specified text.
func (t *Textbox) Insert(text []fmt.Stringer) {
	t.textContent = text
	t.updateTextVisibility()
}

// updateTextVisibility updates content of visible
// text area.
func (t *Textbox) updateTextVisibility() {
	t.textarea.Clear()
	var (
		visibleText       []string
		visibleTextHeight float64 
	)
	
	for i, line := range t.textContent {
		if i < t.startID {
			continue
		}
		if visibleTextHeight > t.bgHeight {
			break;
		}
		
		visibleText = append(visibleText, line.String())
		visibleTextHeight += t.textarea.BoundsOf(line.String()).W()
	}
	for _, txt := range visibleText {
		fmt.Fprintln(t.textarea, txt)
	}
}
