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
	maxBounds  pixel.Rect
	bounds     pixel.Rect
}

// NewProgressBar creates new progress bar with IMDraw
// background bar with specified size, color and label text.
func NewProgressBar(size Size, color color.Color,
	labelText string, max int) *ProgressBar {
	pb := new(ProgressBar)
	pb.size = size
	pb.color = color
	pb.bgDraw = imdraw.New(nil)
	pb.label = NewText(labelText, pb.size-1, 0)
	pb.maxBounds = pb.size.BarSize()
	pb.SetMax(max)
	return pb
}

// NewProgressBarSprite creates new progress bar with
// specified background texture, and label.
func NewProgressBarSprite(bgPic pixel.Picture, labelSize Size,
	labelText string, max int) *ProgressBar {
	pb := new(ProgressBar)
	pb.bgSpr = pixel.NewSprite(bgPic, bgPic.Bounds())
	pb.maxBounds = pb.bgSpr.Picture().Bounds()
	pb.label = NewText(labelText, labelSize, 0)
	pb.SetMax(max)
	return pb
}

// Draw draws progress bar.
func (pb *ProgressBar) Draw(t pixel.Target, matrix pixel.Matrix) {
	widthDiff := ConvSize(pb.maxBounds.W() - pb.Bounds().W())
	barPos := pixel.V(-widthDiff/2, 0)
	mx := matrix.Moved(barPos)
	pb.drawArea = MatrixToDrawArea(mx, pb.Bounds())
	// Background.
	if pb.bgSpr != nil {
		pb.bgSpr.Draw(t, mx)
	} else {
		pb.drawIMBackground(t)
	}
}

// Update updates progress bar.
func (pb *ProgressBar) Update(win *Window) {
}

// Bounds returns bar size bounds.
func (pb *ProgressBar) Bounds() pixel.Rect {
	return pb.bounds
}

// Value retruns current progress value.
func (pb *ProgressBar) Value() int {
	return pb.value
}

// SetValue sets specified value as
// current progress value.
func (pb *ProgressBar) SetValue(val int) {
	pb.value = val
	pb.updateProgress()
}

// Max retruns maximal progress value.
func (pb *ProgressBar) Max() int {
	return pb.max
}

// SetMax sets specified value as progress
// maximal value.
func (pb *ProgressBar) SetMax(max int) {
	pb.max = max
	pb.updateProgress()
}

// drawIMBackground draws bar background with pixel IMDraw.
func (pb *ProgressBar) drawIMBackground(t pixel.Target) {
	pb.bgDraw.Color = pixel.ToRGBA(pb.color)
	pb.bgDraw.Push(pb.drawArea.Min)
	pb.bgDraw.Push(pb.drawArea.Max)
	pb.bgDraw.Rectangle(0)
	pb.bgDraw.Draw(t)
}

// updateProgress updates bar progress to current
// progress value.
func (pb *ProgressBar) updateProgress() {
	valPercent := float64(pb.Value()) * 1.0 / float64(pb.Max())
	bgWidth := pb.maxBounds.W() * valPercent
	pb.bounds = pixel.R(pb.maxBounds.Min.X, pb.maxBounds.Min.Y,
		bgWidth, pb.maxBounds.Max.Y)
	if pb.bgSpr != nil {
		pb.bgSpr = pixel.NewSprite(pb.bgSpr.Picture(), pb.bounds)
	}
}
