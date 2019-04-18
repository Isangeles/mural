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
	"fmt"
	"image/color"
	
	"github.com/faiface/pixel"
)

// Struct for progress bars.
type ProgressBar struct {
	value, max int
	labelText  string
	hovered    bool
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
func NewProgressBar(size Size, color color.Color) *ProgressBar {
	pb := new(ProgressBar)
	pb.size = size
	pb.color = color
	pb.maxBounds = pb.size.BarSize()
	pb.label = NewText(pb.size-1, 0)
	return pb
}

// Draw draws progress bar.
func (pb *ProgressBar) Draw(t pixel.Target, matrix pixel.Matrix) {
	widthDiff := ConvSize(pb.maxBounds.W() - pb.Size().X)
	barPos := pixel.V(-widthDiff/2, 0)
	mx := matrix.Moved(barPos)
	pb.drawArea = MatrixToDrawArea(mx, pb.Size())
	// Background.
	if pb.bgSpr != nil {
		pb.bgSpr.Draw(t, mx)
	} else {
		DrawRectangle(t, pb.DrawArea(), pb.color)
	}
	// Label.
	if pb.hovered {
		pb.label.Draw(t, matrix)
	} 
}

// Update updates progress bar.
func (pb *ProgressBar) Update(win *Window) {
	// On-hover.
	if pb.DrawArea().Contains(win.MousePosition()) {
		pb.hovered = true
	} else {
		pb.hovered = false
	}
}

// SetBackground sets specified sprite as bar
// background, also removes current background color.
func (pb *ProgressBar) SetBackground(p pixel.Picture) {
	bounds := pixel.R(0, p.Bounds().Min.Y, 0, p.Bounds().Max.Y)
	pb.bgSpr = pixel.NewSprite(p, bounds)
	pb.maxBounds = pb.bgSpr.Picture().Bounds()
	pb.SetColor(nil)
}

// SetColor sets specified color as background color.
func (pb *ProgressBar) SetColor(c color.Color) {
	pb.color = c
}

// SetLabel sets specified text as progress label.
func (pb *ProgressBar) SetLabel(t string) {
	pb.labelText = t
}

// Size returns bar size.
func (pb *ProgressBar) Size() pixel.Vec {
	return pb.bounds.Size()
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

// DrawArea returns current last draw area of
// this element.
func (pb *ProgressBar) DrawArea() pixel.Rect {
	return pb.drawArea
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
	pb.label.SetText(fmt.Sprintf("%s:%d/%d", pb.labelText, pb.Value(),
		pb.Max()))
}
