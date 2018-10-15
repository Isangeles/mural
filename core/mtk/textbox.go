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

package mtk

import (
	"fmt"
	"image/color"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	//"golang.org/x/image/colornames"
)

// Struct for textboxes.
type Textbox struct {
	bg          *imdraw.IMDraw
	color       color.Color
	textarea    *text.Text
	drawArea    pixel.Rect // updated at every draw
	textContent []string
	visibleText []string
	startID     int
}

// NewTextbox creates new textbox.
func NewTextbox(fontSize Size, color color.Color) (*Textbox) {
	t := new(Textbox)
	// Background.
	t.bg = imdraw.New(nil)
	t.color = color
	// Text.
	font := MainFont(fontSize)
	atlas := Atlas(&font)
	t.textarea = text.New(pixel.V(0, 0), atlas)
	
	return t
}

// Draw draws textbox.
func (tb *Textbox) Draw(drawArea pixel.Rect, t pixel.Target) {
	// Background.
	tb.drawArea = drawArea
	tb.drawIMBackground(t)
	// Text content.
	tb.textarea.Draw(t, pixel.IM.Moved(pixel.V(drawArea.Min.X,
		drawArea.Max.Y - tb.textarea.BoundsOf("AA").H()))) 
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

// drawIMBackground draws IMDraw background in size of draw area.
func (tb *Textbox) drawIMBackground(t pixel.Target) {
	// TODO: use color from constructor.
	tb.bg.Clear()
	tb.bg.Color = pixel.RGBA{0.1, 0.1, 0.1, 0.5}
	tb.bg.Push(tb.drawArea.Min)
	tb.bg.Color = pixel.RGBA{0.1, 0.1, 0.1, 0.5}
	tb.bg.Push(tb.drawArea.Max)
	tb.bg.Rectangle(0)
	tb.bg.Draw(t)
}

// Bounds returns size parameters of textbox textarea.
func (t *Textbox) Bounds() pixel.Rect {
	return t.textarea.Bounds()
}

// Insert clears textbox and inserts specified text.
func (t *Textbox) Insert(text []fmt.Stringer) {
	t.Clear()
	for _, txt := range text {
		t.Add(txt.String())
	}
	t.updateTextVisibility()
}

// InsertText clears textbox and inserts specified text.
func (t *Textbox) InsertText(text []string) {
	t.textContent = text
	t.updateTextVisibility()
}

// Add adds specified text to box.
func (t *Textbox) Add(line string) {
	t.textContent = append(t.textContent, line)
	t.updateTextVisibility()
}

// Clear clears textbox.
func (t *Textbox) Clear() {
	t.textContent = []string{}
	t.textarea.Clear()
}

// updateTextVisibility updates content of visible
// text area.
func (t *Textbox) updateTextVisibility() {
	t.textarea.Clear()
	var (
		visibleText       []string
		visibleTextHeight float64 
	)
	
	for i := 0; i < len(t.textContent); i++ {
		line := t.textContent[i]
		if i < t.startID {
			continue
		}
		if visibleTextHeight > t.drawArea.Max.Y {
			break;
		}
		
		visibleText = append(visibleText, line)
		visibleTextHeight += t.textarea.BoundsOf(line).W()/2
	}
	for _, txt := range visibleText {
		fmt.Fprintln(t.textarea, ">" + txt)
	}
}
