/*
 * progressbar.go
 *
 * Copyright 2019 Dariusz Sikora <dev@isangeles.pl>
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
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Struct for progress bars.
type ProgressBar struct {
	value, max int
	bgDraw     *imdraw.IMDraw
	bgSpr      *pixel.Sprite
	size       Size
	color      color.Color
	label      *Text
	drawArea   pixel.Rect
}

// NewProgressBar creates new progress bar with IMDraw
// background bar with specified size, color and label text.
func NewProgressBar(size Size, color color.Color,
	labelText string) *ProgressBar {
	pb := new(ProgressBar)
	pb.size = size
	pb.color = color
	pb.bgDraw = imdraw.New(nil)
	pb.label = NewText(labelText, pb.size-1, 0)
	return pb
}

// NewProgressBarSprite creates new progress bar with
// specified background texture, and label.
func NewProgressBarSprite(bgPic pixel.Picture,
	labelSize Size, labelText string) *ProgressBar {
	pb := new(ProgressBar)
	pb.bgSpr = pixel.NewSprite(bgPic, bgPic.Bounds())
	pb.label = NewText(labelText, labelSize, 0)
	return pb
}

// Draw draws progress bar.
func (pb *ProgressBar) Draw(t pixel.Target, matrix pixel.Matrix) {
	pb.drawArea = MatrixToDrawArea(matrix, pb.Frame())
	// Background.
	if pb.bgSpr != nil {
		pb.bgSpr.Draw(t, matrix)
	} else {
		pb.drawIMBackground(t)
	}
}

// Update updates progress bar.
func (pb *ProgressBar) Update(win *Window) {

}

// Frame returns bar size bounds.
func (pb *ProgressBar) Frame() pixel.Rect {
	if pb.bgSpr != nil {
		return pb.bgSpr.Frame()
	}
	return pb.size.BarSize()
}

// drawIMBackground draws bar background with pixel IMDraw.
func (pb *ProgressBar) drawIMBackground(t pixel.Target) {
	pb.bgDraw.Color = pixel.ToRGBA(pb.color)
	pb.bgDraw.Push(pb.drawArea.Min)
	pb.bgDraw.Push(pb.drawArea.Max)
	pb.bgDraw.Rectangle(0)
	pb.bgDraw.Draw(t)
}
