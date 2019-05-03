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
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/areamap"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

var (
	FOW_color pixel.RGBA = pixel.RGBA{0.1, 0.1, 0.1, 0.7}
)

// Struct for HUD camera.
type Camera struct {
	hud      *HUD
	position pixel.Vec
	size     pixel.Vec
	locked   bool
	// Map & objects.
	areaMap *areamap.Map
	fow     *imdraw.IMDraw
	avatars []*object.Avatar
	objects []*object.ObjectGraphic
	// Debug mode.
	cameraInfo *mtk.Text
	cursorInfo *mtk.Text
}

// newCamera creates new instance of camera.
func newCamera(hud *HUD, size pixel.Vec) *Camera {
	c := new(Camera)
	c.hud = hud
	c.size = size
	c.position = pixel.V(0, 0)
	c.fow = imdraw.New(nil)
	c.cameraInfo = mtk.NewText(mtk.SIZE_MEDIUM, 0)
	c.cursorInfo = mtk.NewText(mtk.SIZE_MEDIUM, 0)
	return c
}

// Draw draws camera on specified map.
func (c *Camera) Draw(win *mtk.Window) {
	// Map.
	if c.areaMap != nil {
		//c.areaMap.Draw(win.Window, mtk.Matrix().Moved(c.Position()), c.Size())
		c.areaMap.DrawFull(win.Window, mtk.Matrix().Moved(c.Position()))
	}
	// Avatars.
	for _, av := range c.avatars {
		for _, pc := range c.hud.Players() {
			if mtk.Range(pc.Position(),
				av.Position()) > pc.SightRange() {
				continue
			}
			avPos := c.ConvAreaPos(av.Position())
			av.Draw(win, mtk.Matrix().Moved(avPos))
		}
	}
	// Objects.
	for _, ob := range c.objects {
		for _, pc := range c.hud.Players() {
			if mtk.Range(pc.Position(),
				ob.Position()) > pc.SightRange() {
				continue
			}
			obPos := c.ConvAreaPos(ob.Position())
			ob.Draw(win, mtk.Matrix().Moved(obPos))
		}
	}
	// FOW effect.
	if c.areaMap != nil && config.MapFOW() {
		c.drawMapFOW(win.Window)
	}
	// Debug mode.
	if config.Debug() {
		camInfoPos := mtk.DrawPosBR(win.Bounds(), c.cameraInfo.Size())
		c.cameraInfo.Draw(win, mtk.Matrix().Moved(camInfoPos))
		curInfoPos := mtk.TopOf(c.cameraInfo.DrawArea(), c.cursorInfo.Size(), 10)
		c.cursorInfo.Draw(win, mtk.Matrix().Moved(curInfoPos))
	}
}

// Update updates camera.
func (c *Camera) Update(win *mtk.Window) {
	if c.areaMap == nil {
		return
	}
	if !c.locked {
		mSize := c.areaMap.Size()
		mTileSize := c.areaMap.TileSize()
		offset := pixel.V(mTileSize.X*16, mTileSize.Y*16)
		mapSizePlus := pixel.V(mSize.X+offset.X, mSize.Y+offset.Y)
		// Key events.
		if c.position.Y < mapSizePlus.Y &&
			(win.JustPressed(pixelgl.KeyW) ||
				win.JustPressed(pixelgl.KeyUp)) {
			c.position.Y += mTileSize.Y
		}
		if c.position.X < mapSizePlus.X &&
			(win.JustPressed(pixelgl.KeyD) ||
				win.JustPressed(pixelgl.KeyRight)) {
			c.position.X += mTileSize.X
		}
		if c.position.Y > 0-offset.Y &&
			(win.JustPressed(pixelgl.KeyS) ||
				win.JustPressed(pixelgl.KeyDown)) {
			c.position.Y -= mTileSize.Y
		}
		if c.position.X > 0-offset.X &&
			win.JustPressed(pixelgl.KeyA) ||
			win.JustPressed(pixelgl.KeyLeft) {
			c.position.X -= mTileSize.X
		}
	}
	// Mouse events.	
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		c.onMouseLeftPressed(win.MousePosition())
	}
	if win.JustPressed(pixelgl.MouseButtonRight) {
		c.onMouseRightPressed(win.MousePosition())
	}
	// Avatars.
	for _, av := range c.avatars {
		for _, pc := range c.hud.Players() {
			if mtk.Range(pc.Position(),
				av.Position()) > pc.SightRange() {
				continue
			}
			av.Update(win)
		}
	}
	// Objects.
	for _, ob := range c.objects {
		for _, pc := range c.hud.Players() {
			if mtk.Range(pc.Position(),
				ob.Position()) > pc.SightRange() {
				continue
			}
			ob.Update(win)
		}
	}
	c.cameraInfo.SetText(fmt.Sprintf("camera_pos:%v",
		c.Position()))
	c.cursorInfo.SetText(fmt.Sprintf("cursor_pos:%v",
		win.MousePosition()))
}

// SetMap sets maps for camera.
func (c *Camera) SetMap(m *areamap.Map) {
	c.areaMap = m
}

// SetAvatars sets avatars to draw.
func (c *Camera) SetAvatars(avs []*object.Avatar) {
	c.avatars = avs
}

// SetObjects sets area objects to draw.
func (c *Camera) SetObjects(obs []*object.ObjectGraphic) {
	c.objects = obs
}

// SetPosition sets camera position.
func (c *Camera) SetPosition(pos pixel.Vec) {
	c.position = pos
}

// CenterAt centers camera at specified position.
func (c *Camera) CenterAt(pos pixel.Vec) {
	center := pixel.V(pos.X-c.Size().X/2, pos.Y-c.Size().Y/2)
	c.SetPosition(center)
}

// Map returns current map.
func (c *Camera) Map() *areamap.Map {
	return c.areaMap
}

// Avatars returns all avatars from current area.
func (c *Camera) Avatars() []*object.Avatar {
	return c.avatars
}

// DrawObjects returns all objects with 'drawable'
// objects from current area.
func (c *Camera) DrawObjects() []object.Drawer {
	objects := make([]object.Drawer, 0)
	for _, av := range c.avatars {
		objects = append(objects, av)
	}
	for _, ob := range c.objects {
		objects = append(objects, ob)
	}
	return objects
}

// AreaObjects returns all objects in current area.
func (c *Camera) AreaObjects() []*object.ObjectGraphic {
	return c.objects
}

// Position return camera position.
func (c *Camera) Position() pixel.Vec {
	return c.position
}

// Size returns camera size.
func (c *Camera) Size() pixel.Vec {
	return c.size
}

// Frame returns current camera frame bounds.
func (c *Camera) Frame() pixel.Rect {
	cpos := c.Position()
	csize := c.Size()
	return pixel.R(cpos.X, cpos.Y, cpos.X+csize.X,
		cpos.Y+csize.Y)
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
	return pixel.V(posX+camX, posY+camY)
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
			tileDrawMax := pixel.V(tileDrawMin.X+tileSizeX,
				tileDrawMin.Y+tileSizeY)
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

// Triggered after right mouse button was pressed.
func (c *Camera) onMouseRightPressed(pos pixel.Vec) {
	// Set target.
	if c.hud.containsPos(pos) {
		return
	}
	for _, av := range c.Avatars() {
		if !av.DrawArea().Contains(pos) {
			continue
		}
		log.Dbg.Printf("hud:set_target:%s", av.ID()+"_"+av.Serial())
		c.hud.ActivePlayer().SetTarget(av.Character)
		return
	}
	for _, ob := range c.AreaObjects() {
		if !ob.DrawArea().Contains(pos) {
			continue
		}
		log.Dbg.Printf("hud:set_target:%s", ob.ID()+"_"+ob.Serial())
		c.hud.ActivePlayer().SetTarget(ob.Object)
		return
	}
	c.hud.ActivePlayer().SetTarget(nil)
}

// Triggered after left mouse button was pressed.
func (c *Camera) onMouseLeftPressed(pos pixel.Vec) {
	// Loot.
	for _, av := range c.Avatars() {
		if !av.DrawArea().Contains(pos) || av.Live() || av == c.hud.ActivePlayer() {
			continue
		}
		log.Dbg.Printf("hud:loot:%s_%s", av.ID(), av.Serial())
		c.hud.loot.SetTarget(av)
		c.hud.loot.Show(true)
		return
	}
	for _, ob := range c.AreaObjects() {
		if !ob.DrawArea().Contains(pos) || ob.Live() {
			continue
		}
		log.Dbg.Printf("hud:loot:%s_%s", ob.ID(), ob.Serial())
		c.hud.loot.SetTarget(ob)
		c.hud.loot.Show(true)
		return
	}
	// Dialog.
	for _, av := range c.Avatars() {
		if !av.DrawArea().Contains(pos) || !av.Live() || av == c.hud.ActivePlayer() ||
			len(av.Dialogs()) < 1 {
			continue
		}
		log.Dbg.Printf("hud:dialog:%s_%s", av.ID(), av.Serial())
		dialog := av.Dialogs()[0]
		c.hud.dialog.SetDialog(dialog)
		c.hud.dialog.Show(true)
	}
	// Move active PC.
	destPos := c.ConvCameraPos(pos)
	if !c.hud.game.Paused() && c.Map().Passable(destPos) && !c.hud.containsPos(pos) {
		c.hud.ActivePlayer().SetDestPoint(destPos.X, destPos.Y)
	}
}
