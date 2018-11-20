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
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/mural/core/areamap"
	"github.com/isangeles/mural/core/mtk"
)

// Struct for HUD camera.
type Camera struct {
	hud      *HUD
	position pixel.Vec
	size     pixel.Vec
	areaMap  *areamap.Map
}

// newCamera creates new instance of camera.
func newCamera(hud *HUD, size pixel.Vec) (*Camera) {
	c := new(Camera)
	c.hud = hud
	c.size = size
	c.position = pixel.V(0, 0)
	return c
}

// Draw draws camera on specified map.
func (c *Camera) Draw(win *mtk.Window) {
	//c.areaMap.Draw(win, c.position, c.size)
}

// Update updates camera.
func (c *Camera) Update(win *mtk.Window) {
	if win.JustPressed(pixelgl.KeyW) ||
		win.JustPressed(pixelgl.KeyUp) {
		c.position.Y -= c.areaMap.TileSize().Y
	}
	if win.JustPressed(pixelgl.KeyD) ||
		win.JustPressed(pixelgl.KeyRight) {
		c.position.X += c.areaMap.TileSize().X
	}
	if win.JustPressed(pixelgl.KeyS) ||
		win.JustPressed(pixelgl.KeyDown) {
		c.position.Y += c.areaMap.TileSize().Y
	}
	if win.JustPressed(pixelgl.KeyA) ||
		win.JustPressed(pixelgl.KeyLeft) {
		c.position.X -= c.areaMap.TileSize().X
	}
}

// SetMap sets maps for camera.
func (c *Camera) SetMap(m *areamap.Map) {
	c.areaMap = m
}

// Position return camera position.
func (c *Camera) Position() pixel.Vec {
	return c.position
}
