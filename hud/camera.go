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
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/areamap"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/objects"
)

// Struct for HUD camera.
type Camera struct {
	hud      *HUD
	position pixel.Vec
	size     pixel.Vec
	locked   bool
	// Map.
	areaMap  *areamap.Map
	avatars  []*objects.Avatar
	// Debug mode.
	cameraInfo *mtk.Text
}

// newCamera creates new instance of camera.
func newCamera(hud *HUD, size pixel.Vec) *Camera {
	c := new(Camera)
	c.hud = hud
	c.size = size
	c.position = pixel.V(0, 0)
	c.cameraInfo = mtk.NewText("", mtk.SIZE_MEDIUM, 0)
	return c
}

// Draw draws camera on specified map.
func (c *Camera) Draw(win *mtk.Window) {
	// Map.
	if c.areaMap != nil {
		c.areaMap.DrawWithFOW(win, c.position, c.size,
			c.hud.Player().Position(), c.hud.Player().SightRange())
	}
	// Objects.
	for _, a := range c.avatars {
		if mtk.Range(c.hud.Player().Position(),
			a.Position()) <= c.hud.Player().SightRange() {
				avPos := c.ConvAreaPos(a.Position())
				a.Draw(win, mtk.Matrix().Moved(avPos))
		}
	}
	// Debug mode.
	if config.Debug() {
		c.cameraInfo.Draw(win, mtk.Matrix().Moved(mtk.PosBL(
			c.cameraInfo.Bounds(), win.Bounds().Center())))
	}
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
	c.cameraInfo.SetText(fmt.Sprintf("camera_pos:%v",
		c.Position()))
}

// SetMap sets maps for camera.
func (c *Camera) SetMap(m *areamap.Map) {
	c.areaMap = m
}

// SetAvatars sets avatars to draw.
func (c *Camera) SetAvatars(avs []*objects.Avatar) {
	c.avatars = avs
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

// ConvAreaPos translates specified area
// position to camera position.
func (c *Camera) ConvAreaPos(pos pixel.Vec) pixel.Vec {
	return pixel.V(mtk.ConvSize(pos.X)-c.Position().X,
		mtk.ConvSize(pos.Y)-c.Position().Y)
}

// ConvCameraPos translates specified camera
// position to area position.
func (c *Camera) ConvCameraPos(pos pixel.Vec) pixel.Vec {
	return pixel.V(pos.X+c.Position().X, pos.Y+c.Position().Y)
}
