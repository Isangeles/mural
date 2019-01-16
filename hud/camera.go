/*
 * camera.go
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

package hud

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/areamap"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/object"
)

// Struct for HUD camera.
type Camera struct {
	hud      *HUD
	position pixel.Vec
	size     pixel.Vec
	locked   bool
	// Map.
	areaMap  *areamap.Map
	avatars  []*object.Avatar
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
		for _, pc := range c.hud.Players() {
			c.areaMap.DrawWithFOW(win, c.position, c.size,
				pc.Position(), pc.SightRange())
		}
	}
	// Objects.
	for _, av := range c.avatars {
		for _, pc := range c.hud.Players() {
			if mtk.Range(pc.Position(),
				av.Position()) <= pc.SightRange() {
					avPos := c.ConvAreaPos(av.Position())
					av.Draw(win, mtk.Matrix().Moved(avPos))
				}
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
		offset := pixel.V(c.areaMap.TileSize().X*16,
			c.areaMap.TileSize().Y*16)
		mapSizePlus := pixel.V(c.areaMap.Size().X + offset.X,
			c.areaMap.Size().Y + offset.Y)
		// Key events.
		if c.position.Y < mapSizePlus.Y &&
			(win.JustPressed(pixelgl.KeyW) ||
				win.JustPressed(pixelgl.KeyUp)) {
			c.position.Y += c.areaMap.TileSize().Y
		}
		if c.position.X < mapSizePlus.X &&
			(win.JustPressed(pixelgl.KeyD) ||
				win.JustPressed(pixelgl.KeyRight)) {
			c.position.X += c.areaMap.TileSize().X
		}
		if c.position.Y > 0 - offset.Y &&
			(win.JustPressed(pixelgl.KeyS) ||
				win.JustPressed(pixelgl.KeyDown)) {
			c.position.Y -= c.areaMap.TileSize().Y
		}
		if c.position.X > 0 - offset.X &&
			win.JustPressed(pixelgl.KeyA) ||
			win.JustPressed(pixelgl.KeyLeft) {
			c.position.X -= c.areaMap.TileSize().X
		}
	}
	// Objects.
	for _, av := range c.avatars {
		for _, pc := range c.hud.Players() {
			if mtk.Range(pc.Position(),
				av.Position()) <= pc.SightRange() {
					av.Update(win)
				}
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
func (c *Camera) SetAvatars(avs []*object.Avatar) {
	c.avatars = avs
}

// SetPosition sets camera position.
func (c *Camera) SetPosition(pos pixel.Vec) {
	c.position = pos
}

// CenterAt centers camera at specified position.
func (c *Camera) CenterAt(pos pixel.Vec) {
	c.SetPosition(pixel.V(pos.X - c.Size().X/2,
		pos.Y - c.Size().Y/2))
}

// Avatars returns all avatars from current
// area.
func (c *Camera) Avatars() []*object.Avatar {
	return c.avatars
}

// Position return camera position.
func (c *Camera) Position() pixel.Vec {
	return c.position
}

// Size returns camera size.
func (c *Camera) Size() pixel.Vec {
	return c.size
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
	posX := mtk.ConvSize(pos.X)
	posY := mtk.ConvSize(pos.Y)
	camX := mtk.ConvSize(c.Position().X)
	camY := mtk.ConvSize(c.Position().Y)
	return pixel.V(posX - camX, posY - camY)
}

// ConvCameraPos translates specified camera
// position to area position.
func (c *Camera) ConvCameraPos(pos pixel.Vec) pixel.Vec {
	posX := pos.X
	posY := pos.Y
	camX := c.Position().X
	camY := c.Position().Y
	return pixel.V(posX + camX, posY + camY)
}
