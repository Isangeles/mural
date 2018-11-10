/*
 * camera.go
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

package hud

import (
	"github.com/faiface/pixel"

	"github.com/isangeles/mural/core/areamap"
	"github.com/isangeles/mural/core/mtk"
)

// Struct for HUD camera.
type Camera struct {
	size pixel.Vec
	areaMap *areamap.Map
}

// newCamera creates new instance of camera.
func newCamera(size pixel.Vec) (*Camera) {
	c := new(Camera)
	c.size = size
	return c
}

// Draw draws camera on specified map.
func (c *Camera) Draw(win *mtk.Window) {
	// TODO: draw visible map part.
}

// SetMap sets maps for camera.
func (c *Camera) SetMap(m *areamap.Map) {
	c.areaMap = m
}
