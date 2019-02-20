/*
 * text.go
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
	"strings"
	"image/color"

	"golang.org/x/image/colornames"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

// Text struct for short text like labels, names, etc.
type Text struct {
	text     *text.Text
	content  string
	drawArea pixel.Rect // updated on each draw
	color    color.Color
	fontSize Size
	width    float64
}

// NewText creates new text with specified text content,
// font size and with max width of text line(0 for no max width).
func NewText(fontSize Size, width float64) *Text {
	t := new(Text)
	// Parameters.
	t.fontSize = fontSize
	t.width = width
 	// Text.
	font := MainFont(t.fontSize)
	atlas := Atlas(&font)
	t.text = text.New(pixel.V(0, 0), atlas)
	t.color = colornames.White // default color white
	
	return t
}

// SetText sets specified text as text to display.
func (tx *Text) SetText(text string) {
	tx.content = text
	// If text too wide, then split to more lines.
	if tx.width != 0 && tx.text.BoundsOf(tx.content).W() > tx.width {
		// TODO: figure out something better. Now its only
		// replaces first blank line with '\n'.
		tx.content = strings.Replace(tx.content, " ", "\n", 1)
	}
	tx.JustCenter()
}

// SetColor sets specified color as
// current text color.
func (tx *Text) SetColor(c color.Color) {
	tx.color = c
}

// SetMaxWidth sets maximal width of single text line.
func (tx *Text) SetMaxWidth(width float64) {
	tx.width = width
}

// AddText adds specified text to current text
// content.
func (tx *Text) AddText(text string) {
	tx.content += "\n" + text
	tx.SetText(tx.content)
}

// Write writes specified data as text to text
// area.
func (tx *Text) Write(p []byte) (n int, err error) {
	return tx.text.Write(p)
}

// JustRight adjusts text origin position to right.
func (tx *Text) JustRight() {
	mariginX := (-tx.text.BoundsOf(tx.content).Max.X)
	tx.text.Orig = pixel.V(mariginX, 0)
	tx.text.Clear()
	tx.text.WriteString(tx.content)
}

// JustCenter adjusts text origin position to center.
func (tx *Text) JustCenter() {
	mariginX := (-tx.text.BoundsOf(tx.content).Max.X) / 2
	tx.text.Orig = pixel.V(mariginX, 0)
	tx.text.Clear()
	tx.text.WriteString(tx.content)
}

// JustLeft adjusts text origin to left.
func (tx *Text) JustLeft() {
	tx.text.Orig = pixel.V(0, 0)
	tx.text.Clear()
	tx.text.WriteString(tx.content)
}

// Draw draws text.
func (tx *Text) Draw(t pixel.Target, matrix pixel.Matrix) {
	tx.drawArea = MatrixToDrawArea(matrix, tx.Bounds())
	tx.text.DrawColorMask(t, matrix, tx.color)
}

// Bounds return size of text.
func (tx *Text) Bounds() pixel.Rect {
	return tx.text.Bounds()
}

// BoundsOf returns bounds of specified text
// while displayed.
func (tx *Text) BoundsOf(text string) pixel.Rect {
	return tx.text.BoundsOf(text)
}

// Clear clears texts,
func (tx *Text) Clear() {
	tx.text.Clear()
}

// DrawArea returns current draw area of text.
func (tx *Text) DrawArea() pixel.Rect {
	return tx.drawArea
}

// String returns text content.
func (tx *Text) String() string {
	return tx.content
}
