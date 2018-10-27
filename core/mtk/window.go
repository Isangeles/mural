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
	//"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Wrapper struct for pixel window, to provide scalability.
type Window struct {
	*pixelgl.Window
}

// NewWindow creates new MTK window.
func NewWindow(conf pixelgl.WindowConfig) (*Window, error) {
	initScale(conf.Bounds.Max)
	w := new(Window)
	win, err := pixelgl.NewWindow(conf)
	if err != nil {
		return nil, err
	}
	w.Window = win
	return w, nil
}
