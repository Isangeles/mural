/*
 * animation.go
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

// Struct for animations.
type Animation struct {
	drawFrameID int
	frames      []*pixel.Sprite
	drawArea    pixel.Rect
	fps         int
	lastChange  int64
	looping     bool
	finished    bool
}

// NewAnimation creates new animation with specified
// frames and FPS value.
// Animation is looping by default.
func NewAnimation(frames []*pixel.Sprite, fps int) *Animation {
	anim := new(Animation)
	anim.frames = frames
	anim.drawFrameID = 0
	anim.fps = fps
	anim.looping = true
	return anim
}

// Draw draws current animation frame.
func (anim *Animation) Draw(t pixel.Target, matrix pixel.Matrix) {
	frame := anim.frames[anim.drawFrameID]
	anim.drawArea = MatrixToDrawArea(matrix, frame.Frame().Size())
	frame.Draw(t, matrix)
}

// Update updates animation.
func (anim *Animation) Update(win *Window) {
	if anim.Finished() {
		return
	}
	anim.lastChange += win.Delta()
	if anim.lastChange >= int64(1000 / anim.fps) {
		if !anim.looping && anim.drawFrameID == len(anim.frames)-1 {
			anim.finished = true
			return
		}
		anim.SetCurrentFrameID(anim.drawFrameID + 1)
		anim.lastChange = 0
	}
}

// DrawArea returns current frame draw area.
func (anim *Animation) DrawArea() pixel.Rect {
	return anim.drawArea
}

// SetCurrentFrameID sets frame with specified ID as
// current draw frame of animation. If specified index is
// bigger than maximal frame index then first index is set,
// if smaller than minimal then last index is set.
func (anim *Animation) SetCurrentFrameID(id int) {
	switch {
	case id < 0:
		anim.drawFrameID = len(anim.frames)-1
	case id > len(anim.frames)-1:
		anim.drawFrameID = 0
	default:
		anim.drawFrameID = id
	}
}

// Restarts restarts animation.
func (anim *Animation) Restart() {
	anim.finished = false
	anim.SetCurrentFrameID(0)
}

// Loop toggles animation looping.
func (anim *Animation) Loop(loop bool) {
	anim.looping = loop
}

// Finished checks whether animation is finished, i.e.
// current frame is last frame of animation.
func (anim *Animation) Finished() bool {
	return anim.finished
}

// SetFPS sets number of frames per second.
func (anim *Animation) SetFPS(fps int) {
	anim.fps = fps
}



