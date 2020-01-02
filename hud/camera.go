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
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	flameconf "github.com/isangeles/flame/config"
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/area"
	"github.com/isangeles/flame/core/module/character"

	"github.com/isangeles/stone"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

var (
	FOWColor pixel.RGBA = pixel.RGBA{0.1, 0.1, 0.1, 0.7}
)

const (
	LootRange   = 50
	DialogRange = 50
	ActionRange = 50
)

// Struct for HUD camera.
type Camera struct {
	hud      *HUD
	position pixel.Vec
	size     pixel.Vec
	locked   bool
	area     *area.Area
	// Map & objects.
	areaMap *stone.Map
	fow     *imdraw.IMDraw
	avatars map[string]*object.Avatar
	objects map[string]*object.ObjectGraphic
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
	c.avatars = make(map[string]*object.Avatar)
	c.objects = make(map[string]*object.ObjectGraphic)
	// Debug info.
	textParams := mtk.Params{
		FontSize: mtk.SizeMedium,
	}
	c.cameraInfo = mtk.NewText(textParams)
	c.cursorInfo = mtk.NewText(textParams)
	return c
}

// Draw draws camera on specified map.
func (c *Camera) Draw(win *mtk.Window) {
	// Map.
	if c.areaMap != nil {
		c.areaMap.DrawPart(win.Window, mtk.Matrix().Moved(c.Position()), c.Size())
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
	if c.areaMap != nil && config.MapFOW {
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
	c.updateAreaObjects()
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
		if c.position.X < mapSizePlus.X && (win.JustPressed(pixelgl.KeyD) ||
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
			if mtk.Range(pc.Position(), av.Position()) > pc.SightRange() {
				av.Silence(true)
				continue
			}
			av.Silence(false)
			av.Update(win)
		}
	}
	// Objects.
	for _, ob := range c.objects {
		for _, pc := range c.hud.Players() {
			if mtk.Range(pc.Position(), ob.Position()) > pc.SightRange() {
				ob.Silence(true)
				continue
			}
			ob.Silence(false)
			ob.Update(win)
		}
	}
	c.cameraInfo.SetText(fmt.Sprintf("camera_pos:%v",
		c.Position()))
	c.cursorInfo.SetText(fmt.Sprintf("cursor_pos:%v",
		win.MousePosition()))
}

// SetPosition sets camera position.
func (c *Camera) SetPosition(pos pixel.Vec) {
	c.position = pos
}

// SetArea sets area for camera to display.
func (c *Camera) SetArea(a *area.Area) error {
	c.area = a
	if c.area == nil {
		c.clear()
		return nil
	}
	// Set map.
	chapter := c.hud.game.Module().Chapter()
	mapPath := fmt.Sprintf("%s/gui/chapters/%s/areas/%s/map.tmx",
		chapter.Conf().ModulePath, chapter.ID(), a.ID())
	areaMap, err := stone.NewMap(mapPath)
	if err != nil {
		return fmt.Errorf("fail to create pc area map: %v", err)
	}
	c.areaMap = areaMap
	// PC avatars.
	c.avatars = make(map[string]*object.Avatar)
	for _, char := range a.Characters() {
		var pcAvatar *object.Avatar
		for _, pc := range c.hud.Players() {
			if char == pc.Character {
				pcAvatar = pc
				break
			}
		}
		if pcAvatar == nil {
			continue
		}
		c.avatars[char.ID()+char.Serial()] = pcAvatar
	}
	// Update objects graphics.
	c.updateAreaObjects()
	// Center camera at player.
	pc := c.hud.ActivePlayer()
	c.CenterAt(pc.Position())
	return nil
}

// CenterAt centers camera at specified position.
func (c *Camera) CenterAt(pos pixel.Vec) {
	center := pixel.V(pos.X-c.Size().X/2, pos.Y-c.Size().Y/2)
	c.SetPosition(center)
}

// Map returns current map.
func (c *Camera) Map() *stone.Map {
	return c.areaMap
}

// Avatars returns all avatars from current area.
func (c *Camera) Avatars() (avatars []*object.Avatar) {
	for _, av := range c.avatars {
		avatars = append(avatars, av)
	}
	return
}

// AreaObjects returns all objects in current area.
func (c *Camera) AreaObjects() (objects []*object.ObjectGraphic) {
	for _, ob := range c.objects {
		objects = append(objects, ob)
	}
	return
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

// Area retuns current area.
func (c *Camera) Area() *area.Area {
	return c.area
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

// PassablePosition checks if specified position is 'passable',
// i.e. map there is visible layer on this position where player
// is allowed to move(like 'ground' layer').
func (c *Camera) PassablePosition(pos pixel.Vec) bool {
	layer := c.Map().PositionLayer(pos)
	if layer == nil {
		return false
	}
	return layer.Name() == "ground"
}

// ConvAreaPos translates specified area
// position to camera position.
func (c *Camera) ConvAreaPos(pos pixel.Vec) pixel.Vec {
	drawMatrix := mtk.Matrix().Moved(c.Position())
	drawPos := pixel.V(drawMatrix[4], drawMatrix[5]) 
	drawScale := drawMatrix[0]
	posX := pos.X * drawScale
	posY := pos.Y * drawScale
	drawX := drawPos.X //* drawScale
	drawY := drawPos.Y //* drawScale
	return pixel.V(posX - drawX, posY - drawY)
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
			c.fow.Color = FOWColor
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
	//c.fow.Color = FOWColor
	//c.fow.Push(c.ConvAreaPos(c.hud.ActivePlayer().Position()))
	//c.fow.Circle(c.hud.ActivePlayer().SightRange(), 10)
	c.fow.Draw(t)
}

// updateAreaObjects updates lists with avatars
// and graphical objects for current area.
func (c *Camera) updateAreaObjects() {
	if c.area == nil {
		return
	}
	c.clearAreaObjects()
	// Add new objects & characters.
	for _, char := range c.area.Characters() {
		if c.avatars[char.ID()+char.Serial()] != nil {
			continue
		}
		avData := res.Avatar(char.ID())
		if avData == nil {
			log.Err.Printf("hud camera: update area objects: avatar data not found: %s",
				char.ID())
			continue
		}
		av := object.NewAvatar(char, avData)
		c.avatars[char.ID()+char.Serial()] = av
	}
	for _, ob := range c.area.Objects() {
		if c.objects[ob.ID()+ob.Serial()] != nil {
			continue
		}
		ogData := res.Object(ob.ID())
		if ogData == nil {
			log.Err.Printf("hud camera: update area objects: object data not found: %s",
				ob.ID())
			continue
		}
		og := object.NewObjectGraphic(ob, ogData)
		c.objects[ob.ID()+ob.Serial()] = og
	}
}

// clearAreaObjecs removes avatars & graphics without
// corresponding objects in current area.
func (c *Camera) clearAreaObjects() {
	for _, av := range c.Avatars() {
		found := false
		for _, char := range c.area.Characters() {
			if char.ID() == av.ID() && char.Serial() == av.Serial() {
				found = true
				break
			}
		}
		if found {
			continue
		}
		delete(c.avatars, av.ID()+av.Serial())
	}
	for _, gob := range c.AreaObjects() {
		found := false
		for _, ob := range c.area.Objects() {
			if ob.ID() == gob.ID() && ob.Serial() == gob.Serial() {
				found = true
				break
			}
		}
		if found {
			continue
		}
		delete(c.objects, gob.ID()+gob.Serial())
	}
}

// clear removes removes current area and all objects.
func (c *Camera) clear() {
	c.area = nil
	c.areaMap = nil
	c.avatars = make(map[string]*object.Avatar)
	c.objects = make(map[string]*object.ObjectGraphic)
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
		log.Dbg.Printf("hud: set target: %s", av.ID()+"_"+av.Serial())
		c.hud.ActivePlayer().SetTarget(av.Character)
		return
	}
	for _, ob := range c.AreaObjects() {
		if !ob.DrawArea().Contains(pos) {
			continue
		}
		log.Dbg.Printf("hud: set target: %s", ob.ID()+"_"+ob.Serial())
		c.hud.ActivePlayer().SetTarget(ob.Object)
		return
	}
	c.hud.ActivePlayer().SetTarget(nil)
}

// Triggered after left mouse button was pressed.
func (c *Camera) onMouseLeftPressed(pos pixel.Vec) {
	pc := c.hud.ActivePlayer()
	langPath := flameconf.LangPath()
	// Action.
	for _, ob := range c.AreaObjects() {
		if !ob.DrawArea().Contains(pos) || !ob.Live() || ob.Action() == nil {
			continue
		}
		// Range check.
		r := math.Hypot(ob.Position().X-pc.Position().X, ob.Position().Y-pc.Position().Y)
		if r > ActionRange {
			pc.SendPrivate(lang.TextDir(langPath, "tar_too_far"))
			continue
		}
		log.Dbg.Printf("hud: action: %s#%s", ob.ID(), ob.Serial())
		pc.TakeModifiers(ob.Object, ob.Action().UserMods()...)
		ob.TakeModifiers(ob.Object, ob.Action().SelfMods()...)
		return
	}
	// Loot.
	for _, av := range c.Avatars() {
		if !av.DrawArea().Contains(pos) || av.Live() || av == pc {
			continue
		}
		// Range check.
		r := math.Hypot(av.Position().X-pc.Position().X, av.Position().Y-pc.Position().Y)
		if r > LootRange {
			pc.SendPrivate(lang.TextDir(langPath, "tar_too_far"))
			continue
		}
		// Show loot window.
		log.Dbg.Printf("hud: loot: %s#%s", av.ID(), av.Serial())
		c.hud.loot.SetTarget(av)
		c.hud.loot.Show(true)
		return
	}
	for _, ob := range c.AreaObjects() {
		if !ob.DrawArea().Contains(pos) || ob.Live() {
			continue
		}
		// Range check.
		r := math.Hypot(ob.Position().X-pc.Position().X, ob.Position().Y-pc.Position().Y)
		if r > LootRange {
			pc.SendPrivate(lang.TextDir(langPath, "tar_too_far"))
			continue
		}
		// Show loot window.
		log.Dbg.Printf("hud: loot: %s#%s", ob.ID(), ob.Serial())
		c.hud.loot.SetTarget(ob)
		c.hud.loot.Show(true)
		return
	}
	// Dialog.
	for _, av := range c.Avatars() {
		if !av.DrawArea().Contains(pos) || !av.Live() || av == pc ||
			av.AttitudeFor(pc) == character.Hostile || len(av.Dialogs()) < 1 {
			continue
		}
		// Range check.
		r := math.Hypot(av.Position().X-pc.Position().X, av.Position().Y-pc.Position().Y)
		if r > DialogRange {
			pc.SendPrivate(lang.TextDir(langPath, "tar_too_far"))
			continue
		}
		// Show dialog window.
		log.Dbg.Printf("hud: dialog: %s#%s", av.ID(), av.Serial())
		dialog := av.Dialogs()[0]
		c.hud.dialog.SetDialog(dialog)
		c.hud.dialog.Show(true)
	}
	// Move active PC.
	destPos := c.ConvCameraPos(pos)
	if !c.hud.game.Paused() && c.PassablePosition(destPos) && !c.hud.containsPos(pos) {
		c.hud.ActivePlayer().SetDestPoint(destPos.X, destPos.Y)
	}
}


