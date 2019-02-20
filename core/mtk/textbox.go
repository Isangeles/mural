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
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Struct for textboxes.
type Textbox struct {
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
func NewTextbox(textSize pixel.Vec, fontSize Size, color color.Color) *Textbox {
	t := new(Textbox)
	t.textSize = textSize
	// Background.
	t.color = color
	// Text.
	t.textarea = NewText(fontSize, t.textSize.X)
	t.textarea.JustLeft()
	return t
}

// Draw draws textbox.
func (tb *Textbox) Draw(drawArea pixel.Rect, t pixel.Target) {
	// Background.
	tb.drawArea = drawArea
	DrawRectangle(t, tb.DrawArea(), pixel.RGBA{0.1, 0.1, 0.1, 0.5})
	// Text content.
	tb.textarea.Draw(t, Matrix().Moved(pixel.V(tb.DrawArea().Min.X,
		tb.DrawArea().Max.Y-tb.textarea.BoundsOf("AA").H())))
}

// Update handles key events.
func (t *Textbox) Update(win *Window) {
	if win.JustPressed(pixelgl.KeyDown) {
		if t.startID < len(t.textContent)-1 {
			t.startID++
			t.updateTextVisibility()
		}
	}
	if win.JustPressed(pixelgl.KeyUp) {
		if t.startID > 0 {
			t.startID--
			t.updateTextVisibility()
		}
	}
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

// String returns textbox content.
func (t *Textbox) String() string {
	content := ""
	for _, line := range t.textContent {
		content = fmt.Sprintf("%s\n%s", content, line)
	}
	return strings.TrimSpace(content)
}

// updateTextVisibility updates content of visible
// text area.
func (t *Textbox) updateTextVisibility() {
	var (
		visibleText       []string
		visibleTextHeight float64
	)
	boxWidth := t.DrawArea().W()
	for i := 0; i < len(t.textContent); i++ {
		if i < t.startID {
			continue
		}
		if visibleTextHeight > t.drawArea.H() {
			break
		}
		line := t.textContent[i]
		if len(line) < 1 {
			continue
		}
		breakLines := t.Break(line, boxWidth)
		visibleText = append(visibleText, breakLines...)
		visibleTextHeight += t.textarea.BoundsOf(line).H()*float64(len(breakLines))
	}
	t.textarea.Clear()
	for _, txt := range visibleText {
		fmt.Fprintf(t.textarea, txt)
	}
}

// Break breaks specified line into few lines with specified
// maximal width.
func (t *Textbox) Break(line string, width float64) []string {
	lines := make([]string, 0)
	lineWidth := t.textarea.BoundsOf(line).W()
	if width > 0 && lineWidth > width {
		breakLines := SplitSubN(line, len(line)/2)
		if len(breakLines) < 2 {
			return breakLines
		}
		lines = append(lines, breakLines[0] + "\n")
		breakLineWidth := t.textarea.BoundsOf(breakLines[1]).W()
		if breakLineWidth > width {
			for _, l := range t.Break(breakLines[1], width) {
				lines = append(lines, l + "\n")
			}
		} else {
			lines = append(lines, breakLines[1])
		}
	} else {
		lines = append(lines, line)
	}
	return lines
}

// Splits string at specified index.
// Author: mozey(@stackoverflow).
func SplitSubN(s string, n int) []string {
	if n == 0 {
		return []string{s}
	}
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}
