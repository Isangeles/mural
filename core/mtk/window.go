/*
 * window.go
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
	"time"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Wrapper struct for pixel window, to provide scalability.
type Window struct {
	*pixelgl.Window

	lastUpdate time.Time
	delta      int64 // time from last update in millis
	frameCount int
	fps        int
}

// NewWindow creates new MTK window.
func NewWindow(conf pixelgl.WindowConfig) (*Window, error) {
	initScale(conf.Bounds.Max)
	w := new(Window)
	win, err := pixelgl.NewWindow(conf)
	if err != nil {
		return nil, err
	}
	win.SetSmooth(true)
	w.Window = win
	return w, nil
}

// Update updates window.
func (w *Window) Update() {
	w.Window.Update()
	w.frameCount++
	select {
	case <-sec_timer:
		w.fps = w.frameCount
		w.frameCount = 0
	default:
	}
	dtNano := time.Since(w.lastUpdate).Nanoseconds()
	w.delta = dtNano / int64(time.Millisecond) // delta to milliseconds
	w.lastUpdate = time.Now()
}

// FPS returns current frame per second value.
func (w *Window) FPS() int {
	return w.fps
}

// Delta returns time from last window update
// in milliseconds.
func (w *Window) Delta() int64 {
	return w.delta
}

// PointTL returns position of top left corner of
// window.
func (w *Window) PointTL() pixel.Vec {
	return pixel.V(w.Bounds().Min.X,
		w.Bounds().Max.Y)
}

// PointBR returns position of bottom right corner
// of window.
func (w *Window) PointBR() pixel.Vec {
	return pixel.V(w.Bounds().Max.X,
		w.Bounds().Min.Y)
}
