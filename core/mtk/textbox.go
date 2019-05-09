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
	bgSize      pixel.Vec
	color       color.Color
	textarea    *Text
	drawArea    pixel.Rect // updated at every draw
	upButton    *Button
	downButton  *Button
	textContent []string   // every line of text content
	visibleText []string
	startID     int
	buttons     bool
}

// NewTextbox creates new textbox with specified font size,
// and background color.
func NewTextbox(size pixel.Vec, buttonSize, fontSize Size,
	color, accentColor color.Color) *Textbox {
	t := new(Textbox)
	// Background.
	t.bgSize = size
	t.color = color
	// Text.
	t.textarea = NewText(fontSize, t.bgSize.X)
	t.textarea.JustLeft()
	// TODO: Buttons.
	return t
}

// Draw draws textbox.
func (tb *Textbox) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Background.
	tb.drawArea = MatrixToDrawArea(matrix, tb.Size())
	DrawRectangle(t, tb.DrawArea(), pixel.RGBA{0.1, 0.1, 0.1, 0.5})
	// Text content.
	tb.textarea.Draw(t, Matrix().Moved(pixel.V(tb.DrawArea().Min.X,
		tb.DrawArea().Max.Y-tb.textarea.BoundsOf("AA").H())))
}

// Update handles key events.
func (tb *Textbox) Update(win *Window) {
	// Key events.
	if win.JustPressed(pixelgl.KeyDown) {
		if tb.startID < len(tb.textContent)-1 {
			tb.startID++
			tb.updateTextVisibility()
		}
	}
	if win.JustPressed(pixelgl.KeyUp) {
		if tb.startID > 0 {
			tb.startID--
			tb.updateTextVisibility()
		}
	}
}

// SetSize sets background size.
func (tb *Textbox) SetSize(s pixel.Vec) {
	tb.bgSize = s
}

// Size returns size of textbox background.
func (tb *Textbox) Size() pixel.Vec {
	return tb.bgSize
}

// TextSize returns size of text content.
func (tb *Textbox) TextSize() pixel.Vec {
	return ConvVec(tb.textarea.Size())
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

// SetText clears textbox and inserts specified
// lines of text.
func (tb *Textbox) SetText(text ...string) {
	tb.Clear()
	tb.textContent = text
	tb.updateTextVisibility()
}

// AddText adds specified text to box.
func (tb *Textbox) AddText(line string) {
	tb.textContent = append(tb.textContent, line)
	tb.updateTextVisibility()
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

// updateTextVisibility updates conte nt of visible
// text area.
func (t *Textbox) updateTextVisibility() {
	var (
		visibleText       []string
		visibleTextHeight float64
	)
	boxWidth := t.Size().X
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
		breakLines := t.breakLine(line, boxWidth)
		visibleText = append(visibleText, breakLines...)
		visibleTextHeight += t.textarea.BoundsOf(line).H() * float64(len(breakLines))
	}
	t.textarea.Clear()
	for _, txt := range visibleText {
		fmt.Fprintf(t.textarea, txt)
	}
}

// breakLine breaks specified line into few lines with specified
// maximal width.
func (t *Textbox) breakLine(line string, width float64) []string {
	lines := make([]string, 0)
	lineWidth := t.textarea.BoundsOf(line).W()
	if width > 0 && lineWidth > width {
		breakPoint := t.breakPoint(line, width)
		breakLines := SplitSubN(line, breakPoint)
		for i, l := range breakLines {
			if !strings.HasSuffix(l, "\n") {
				breakLines[i] += "\n"
			}	
		}
		lines = append(lines, breakLines...)
	} else {
		lines = append(lines, line)
	}
	return lines
}

// breakPoint return break position for specified line and width.
func (t *Textbox) breakPoint(line string, width float64) int {
	checkLine := ""
	for i, c := range line {
		checkLine += string(c)
		if t.textarea.BoundsOf(checkLine).W() >= width {
			return i
		}
	}
	return len(line)-1
}

// ScrollBottom scrolls textbox to last lines
// of text content.
func (t *Textbox) ScrollBottom() {
	t.startID = len(t.textContent)-1
}

// Triggered after button up clicked.
func (t *Textbox) onButtonUpClicked(b *Button) {	
}

// Triggered after button down clicked.
func (t *Textbox) onButtonDownClicked(b *Button) {
}

// Splits string to chunks with n as max chunk width.
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
