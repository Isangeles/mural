/*
 * textbox.go
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
	"bytes"
	"fmt"
	"image/color"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	//"golang.org/x/image/colornames"
)

// Struct for textboxes.
type Textbox struct {
	bg          *imdraw.IMDraw
	color       color.Color
	textarea    *Text
	textSize    pixel.Vec
	drawArea    pixel.Rect // updated at every draw
	textContent []string
	visibleText []string
	startID     int
}

// NewTextbox creates new textbox with specified font size,
// background color and maximal size of text content (0 for
// no maximal values).
func NewTextbox(textSize pixel.Vec, fontSize Size,
	color color.Color) (*Textbox) {
	t := new(Textbox)
	t.textSize = textSize
	// Background.
	t.bg = imdraw.New(nil)
	t.color = color
	// Text.
	t.textarea = NewText("", fontSize, t.textSize.X)
	t.textarea.JustLeft()
	return t
}

// Draw draws textbox.
func (tb *Textbox) Draw(drawArea pixel.Rect, t pixel.Target) {
	// Background.
	tb.drawArea = drawArea
	tb.drawIMBackground(t)
	// Text content.
	tb.textarea.Draw(t, Matrix().Moved(pixel.V(tb.DrawArea().Min.X,
		tb.DrawArea().Max.Y - tb.textarea.BoundsOf("AA").H()))) 
}

// Update handles key events.
func (t *Textbox) Update(win *Window) {
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

// drawIMBackground draws IMDraw background in size of
// current draw area.
func (tb *Textbox) drawIMBackground(t pixel.Target) {
	// TODO: use color from constructor.
	tb.bg.Clear()
	tb.bg.Color = pixel.RGBA{0.1, 0.1, 0.1, 0.5}
	tb.bg.Push(tb.drawArea.Min)
	tb.bg.Push(tb.drawArea.Max)
	tb.bg.Rectangle(0)
	tb.bg.Draw(t)
}

// Bounds returns size parameters of textbox textarea.
func (t *Textbox) Bounds() pixel.Rect {
	return t.textarea.Bounds()
}

// DrawArea returns current draw area of text box
// background.
func (t *Textbox) DrawArea() pixel.Rect {
	return t.drawArea
}

// SetMaxTextWidth sets maximal width of single
// line in text area.
func (tb *Textbox) SetMaxTextWidth(width float64) {
	tb.textarea.SetMaxWidth(width)
}

// Insert clears textbox and inserts text from
// String() function of specified stringers.
func (t *Textbox) Insert(text []fmt.Stringer) {
	t.Clear()
	for _, txt := range text {
		t.Add(txt.String())
	}
	t.updateTextVisibility()
}

// InsertText clears textbox and inserts specified text.
func (t *Textbox) InsertText(text ...string) {
	t.Clear()
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
	t.updateTextVisibility()
}

// updateTextVisibility updates content of visible
// text area.
// TODO: height of visible text is calculated wrong, value is too big.
// TODO: break too wide text into more lines.
func (t *Textbox) updateTextVisibility() {
	var (
		visibleText       []string
		visibleTextHeight float64 
	)
	
	for i := 0; i < len(t.textContent); i++ {
		if i < t.startID {
			continue
		}
		if visibleTextHeight > t.drawArea.H() {
			break
		}
		
		line := t.textContent[i]
		visibleText = append(visibleText, line)
		visibleTextHeight += t.textarea.BoundsOf(line).H()
	}
	t.textarea.Clear()
	for _, txt := range visibleText {
		fmt.Fprintf(t.textarea, txt)
	}
	//t.textarea.JustLeft()
}

// Splits string at specified index.
// Author: mozey(@stackoverflow).
func SplitSubN(s string, n int) []string {
    sub := ""
    subs := []string{}

    runes := bytes.Runes([]byte(s))
    l := len(runes)
    for i, r := range runes {
        sub = sub + string(r)
        if (i + 1) % n == 0 {
            subs = append(subs, sub)
            sub = ""
        } else if (i + 1) == l {
            subs = append(subs, sub)
        }
    }

    return subs
}
