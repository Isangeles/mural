/*
 * multianimation.go
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
	"github.com/faiface/pixel"
)

// Struct with sparate animation for each
// direction(up, right, down and left).
type MultiAnimation struct {
	drawAnim  *Animation
	upAnim    *Animation
	rightAnim *Animation
	downAnim  *Animation
	leftAnim  *Animation
}

// NewMultiAnimation creates new multi direction animation
// from specified animations(up, right, down, left).
func NewMultiAnimation(up, right, down, left *Animation) *MultiAnimation {
	ma := new(MultiAnimation)
	ma.upAnim = up
	ma.rightAnim = right
	ma.downAnim = down
	ma.leftAnim = left
	ma.drawAnim = ma.upAnim
	return ma
}

// Draw draws animation for current dirtection.
func (ma *MultiAnimation) Draw(t pixel.Target, matrix pixel.Matrix) {
	ma.drawAnim.Draw(t, matrix)
}

// Update updates animation for current direction.
func (ma *MultiAnimation) Update(win *Window) {
	ma.drawAnim.Update(win)
}

// Up sets animation direction to up.
func (ma *MultiAnimation) Up() {
	ma.drawAnim = ma.upAnim
}

// Right sets animation direction to right.
func (ma *MultiAnimation) Right() {
	ma.drawAnim = ma.rightAnim
}

// Down sets animation direction to down.
func (ma *MultiAnimation) Down() {
	ma.drawAnim = ma.downAnim
}

// Left sets animation direction to left.
func (ma *MultiAnimation) Left() {
	ma.drawAnim = ma.leftAnim
}

// DrawArea returns current draw area.
func (ma *MultiAnimation) DrawArea() pixel.Rect {
	return ma.drawAnim.DrawArea()
}
