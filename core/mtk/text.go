/*
 * text.go
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
	"strings"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

// Text struct for short text like labels, names, etc.
type Text struct {
	text     *text.Text
	content  string
	drawArea pixel.Rect // updated on each draw
	fontSize Size
	width    float64
}

// NewText creates new text with specified text content,
// font size and with max width of text line(0 for no max width).
// Note that text is adjusted to center by default.
// Adjsust to center must be done here, beacuse doing this
// by function don't work.
func NewText(content string, fontSize Size, width float64) *Text {
	t := new(Text)
	// Parameters.
	t.fontSize = fontSize
	t.width = width
 	// Text.
	t.content = content
	font := MainFont(t.fontSize)
	atlas := Atlas(&font)
	t.text = text.New(pixel.V(0, 0), atlas)
	// If text too wide, then split to more lines.
	if t.width != 0 && t.text.BoundsOf(t.content).W() > t.width {
		t.content = strings.Replace(t.content, " ", "\n", 1)
	}
	mariginX := (-t.text.BoundsOf(t.content).Max.X) / 2
	t.text.Orig = pixel.V(mariginX, 0)
	t.text.Clear()
	t.text.WriteString(t.content)
	
	return t
}

// SetText sets specified text as text to display.
func (tx *Text) SetText(text string) {
	tx.content = text
	// If text too wide, then split to more lines.
	if tx.width != 0 && tx.text.BoundsOf(tx.content).W() > tx.width {
		tx.content = strings.Replace(tx.content, " ", "\n", 1)
	}
	mariginX := (-tx.text.BoundsOf(tx.content).Max.X) / 2
	tx.text.Orig = pixel.V(mariginX, 0)
	tx.text.Clear()
	tx.text.WriteString(tx.content)
}

// Adjust text origin position to center.
// TODO: don't work well.
func (tx *Text) JustCenter() {
	mariginX := (-tx.text.BoundsOf(tx.content).Max.X) / 2
	tx.text.Orig = pixel.V(mariginX, 0)
	tx.text.Clear()
	tx.text.WriteString(tx.content)
}

// Draw draws text.
func (tx *Text) Draw(t pixel.Target, matrix pixel.Matrix) {
	tx.drawArea = MatrixToDrawArea(matrix, tx.Bounds())
	tx.text.Draw(t, matrix)
}

// Bounds return size of text.
func (tx *Text) Bounds() pixel.Rect {
	return tx.text.Bounds()
}

// DrawArea returns current draw area of text.
func (tx *Text) DrawArea() pixel.Rect {
	return tx.drawArea
}
