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
	"github.com/faiface/pixel/imdraw"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/areamap"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/object"
)

var (
	FOW_color pixel.RGBA = pixel.RGBA{0.1, 0.1, 0.1, 0.7}
)

// Struct for HUD camera.
type Camera struct {
	hud        *HUD
	position   pixel.Vec
	size       pixel.Vec
	locked     bool
	// Map & objects.
	areaMap    *areamap.Map
	fow        *imdraw.IMDraw
	avatars    []*object.Avatar
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
	c.fow = imdraw.New(nil)
	return c
}

// Draw draws camera on specified map.
func (c *Camera) Draw(win *mtk.Window) {
	// Map.
	if c.areaMap != nil {
		//c.areaMap.Draw(win.Window, mtk.Matrix().Moved(c.Position()), c.Size())
		c.areaMap.DrawFull(win.Window, mtk.Matrix().Moved(c.Position()))
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
	// FOW effect.
	if config.MapFOW() {
		c.drawMapFOW(win.Window)
	}
	// Debug mode.
	if config.Debug() {
		c.cameraInfo.Draw(win, mtk.Matrix().Moved(mtk.PosBR(
			c.cameraInfo.Bounds(), win.PointBR())))
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
	return areamap.MapDrawPos(pos, mtk.Matrix().Moved(c.Position()))
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

// VisibleForPlayers checks whether specified position is
// in visibility range of any HUD PCs.
func (c *Camera) VisibleForPlayers(pos pixel.Vec) bool {
	for _, pc := range c.hud.Players() {
		if mtk.Range(pc.Position(), pos) <= pc.SightRange() {
			return true
		}
	}
	return false
}

// drawMapFOW draws 'Fog Of War' effect on current area map.
func (c *Camera) drawMapFOW(t pixel.Target) {
	c.fow.Clear()
	w, h := 0.0, 0.0
	for h < c.areaMap.Size().Y {
		pos := pixel.V(w, h)
		if !c.VisibleForPlayers(pos) {
			// Draw FOW tile.
			tileSizeX := mtk.ConvSize(c.areaMap.TileSize().X)
			tileSizeY := mtk.ConvSize(c.areaMap.TileSize().Y)
			tileDrawMin := c.ConvAreaPos(pos)
			tileDrawMax := pixel.V(tileDrawMin.X + tileSizeX,
				tileDrawMin.Y + tileSizeY)
			c.fow.Color = FOW_color
			c.fow.Push(tileDrawMin)
			c.fow.Push(tileDrawMax)
			c.fow.Rectangle(0)
		}
		// Next tile.
		w += c.areaMap.TileSize().X
		if w > c.areaMap.Size().X {
			w = 0.0
			h += c.areaMap.TileSize().Y
		}
	}
	// 'FOW ring' effect.
	//c.fow.Color = FOW_color
	//c.fow.Push(c.ConvAreaPos(c.hud.ActivePlayer().Position()))
	//c.fow.Circle(c.hud.ActivePlayer().SightRange(), 10)
	c.fow.Draw(t)
}
