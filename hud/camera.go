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
	locked   bool
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
	if c.areaMap != nil {
		//c.areaMap.Draw(win, c.position, c.size)
		// Drawing map in circular form is faster and simulates FOW.
		// TODO: use player position and sight range.
		//playerPos := c.translatePos(c.hud.Player().Position())
		//c.areaMap.DrawCircle(win, playerPos, c.hud.Player().SightRange())
		c.areaMap.DrawForChar(win, c.position, c.size, c.hud.Player())
	}
	playerPos := c.translatePos(c.hud.Player().Position())
	c.hud.Player().Draw(win, mtk.Matrix().Moved(playerPos))
}

// Update updates camera.
func (c *Camera) Update(win *mtk.Window) {
	if c.areaMap == nil {
		return
	}
	if !c.locked {
		// Key events.
		if c.position.Y < c.areaMap.Size().Y &&
			(win.JustPressed(pixelgl.KeyW) ||
			win.JustPressed(pixelgl.KeyUp)) {
			c.position.Y += c.areaMap.TileSize().Y
		}
		if c.position.X < c.areaMap.Size().X &&
			(win.JustPressed(pixelgl.KeyD) ||
			win.JustPressed(pixelgl.KeyRight)) {
			c.position.X += c.areaMap.TileSize().X
		}
		if c.position.Y > 0 &&
			(win.JustPressed(pixelgl.KeyS) ||
			win.JustPressed(pixelgl.KeyDown)) {
			c.position.Y -= c.areaMap.TileSize().Y
		}
		if c.position.X > 0 &&
			win.JustPressed(pixelgl.KeyA) ||
			win.JustPressed(pixelgl.KeyLeft) {
			c.position.X -= c.areaMap.TileSize().X
		}
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

// Lock toggles camera lock.
func (c *Camera) Lock(lock bool) {
	c.locked = lock
}

// Locked checks whether camera is locked.
func (c *Camera) Locked() bool {
	return c.locked
}

// TranslatePos translates specified position to
// position on camera.
func (c *Camera) translatePos(pos pixel.Vec) pixel.Vec {
	return pixel.V(pos.X - c.Position().X, pos.Y - c.Position().Y)
}
