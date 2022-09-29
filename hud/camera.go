/*
 * camera.go
 *
 * Copyright 2018-2022 Dariusz Sikora <ds@isangeles.dev>
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
	"path/filepath"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame/area"
	"github.com/isangeles/flame/character"
	flameob "github.com/isangeles/flame/object"
	"github.com/isangeles/flame/objects"

	"github.com/isangeles/stone"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/data"
	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/object"
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
	avatars *sync.Map
	objects *sync.Map
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
	c.avatars = new(sync.Map)
	c.objects = new(sync.Map)
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
		if config.MapFull {
			c.areaMap.Draw(win.Window, mtk.Matrix().Moved(c.Position()))
		} else {
			c.areaMap.DrawPart(win.Window, mtk.Matrix().Moved(c.Position()), c.Size())
		}
	}
	// Avatars.
	for _, av := range c.Avatars() {
		if !c.VisibleForPlayer(av.Position()) {
			continue
		}
		avPos := c.ConvAreaPos(av.Position())
		av.Draw(win, mtk.Matrix().Moved(avPos))
	}
	// Objects.
	for _, ob := range c.AreaObjects() {
		if !c.VisibleForPlayer(ob.Position()) {
			continue
		}
		obPos := c.ConvAreaPos(ob.Position())
		ob.Draw(win, mtk.Matrix().Moved(obPos))
	}
	// FOW effect.
	if c.areaMap != nil && config.MapFOW {
		c.drawMapFOW(win.Window)
	}
	// Debug mode.
	if config.Debug {
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
	c.size = win.Bounds().Size()
	c.updateAreaObjects()
	if !c.locked {
		mSize := c.areaMap.Size()
		mTileSize := c.areaMap.TileSize()
		offset := pixel.V(mTileSize.X*16, mTileSize.Y*16)
		mapSizePlus := pixel.V(mSize.X+offset.X, mSize.Y+offset.Y)
		// Key events.
		if c.position.Y < mapSizePlus.Y && win.Pressed(pixelgl.KeyW) ||
			win.Pressed(pixelgl.KeyUp) {
			c.position.Y += mTileSize.Y
		}
		if c.position.X < mapSizePlus.X && win.Pressed(pixelgl.KeyD) ||
			win.Pressed(pixelgl.KeyRight) {
			c.position.X += mTileSize.X
		}
		if c.position.Y > 0-offset.Y && win.Pressed(pixelgl.KeyS) ||
			win.Pressed(pixelgl.KeyDown) {
			c.position.Y -= mTileSize.Y
		}
		if c.position.X > 0-offset.X && win.Pressed(pixelgl.KeyA) ||
			win.Pressed(pixelgl.KeyLeft) {
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
	for _, av := range c.Avatars() {
		if !c.VisibleForPlayer(av.Position()) {
			continue
		}
		av.Silence(false)
		av.Update(win)
	}
	// Objects.
	for _, ob := range c.AreaObjects() {
		if !c.VisibleForPlayer(ob.Position()) {
			continue
		}
		ob.Silence(false)
		ob.Update(win)
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
	chapter := c.hud.game.Chapter()
	mapPath := filepath.Join(config.GUIPath, "chapters", chapter.ID(),
		"areas", a.ID(), "map.tmx")
	areaMap, err := stone.NewMap(mapPath)
	if err != nil {
		return fmt.Errorf("unable to create pc area map: %v", err)
	}
	c.areaMap = areaMap
	// PC avatars.
	c.avatars = new(sync.Map)
	for _, ob := range a.Objects() {
		char, ok := ob.(*character.Character)
		if !ok {
			continue
		}
		var pcAvatar *object.Avatar
		for _, av := range c.hud.playerAvatars {
			if char.ID() == av.ID() && char.Serial() == av.Serial() {
				pcAvatar = av
				break
			}
		}
		if pcAvatar == nil {
			continue
		}
		c.avatars.Store(char.ID()+char.Serial(), pcAvatar)
	}
	// Update objects graphics.
	c.updateAreaObjects()
	// Center camera at player
	pcAvatar := c.hud.PCAvatar()
	if pcAvatar != nil {
		c.CenterAt(pcAvatar.Position())
	}
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
	addAvatar := func(k, v interface{}) bool {
		av, ok := v.(*object.Avatar)
		if !ok {
			return true
		}
		avatars = append(avatars, av)
		return true
	}
	c.avatars.Range(addAvatar)
	return
}

// AreaObjects returns all objects in current area.
func (c *Camera) AreaObjects() (objects []*object.ObjectGraphic) {
	addObject := func(k, v interface{}) bool {
		ob, ok := v.(*object.ObjectGraphic)
		if !ok {
			return true
		}
		objects = append(objects, ob)
		return true
	}
	c.objects.Range(addObject)
	return
}

// DrawObjects returns all objects with 'drawable'
// objects from current area.
func (c *Camera) DrawObjects() []object.Drawer {
	objects := make([]object.Drawer, 0)
	for _, av := range c.Avatars() {
		objects = append(objects, av)
	}
	for _, ob := range c.AreaObjects() {
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
	return pixel.V(posX-drawX, posY-drawY)
}

// ConvCameraPos translates specified camera
// position to area position.
func (c *Camera) ConvCameraPos(pos pixel.Vec) pixel.Vec {
	areaPos := pixel.V(pos.X+c.Position().X, pos.Y+c.Position().Y)
	// Unscale position.
	drawScale := mtk.Matrix()[0]
	areaPos.X = math.Round(areaPos.X / drawScale)
	areaPos.Y = math.Round(areaPos.Y / drawScale)
	return areaPos
}

// VisibleForPlayer checks whether specified position is
// in visibility range of any PC.
func (c *Camera) VisibleForPlayer(pos pixel.Vec) bool {
	for _, pc := range c.hud.Game().PlayerChars() {
		if pc.InSight(pos.X, pos.Y) {
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
		if !c.VisibleForPlayer(pos) {
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
	for _, ob := range c.area.Objects() {
		char, isChar := ob.(*character.Character)
		if !isChar {
			continue
		}
		_, avExists := c.avatars.Load(char.ID() + char.Serial())
		if avExists {
			continue
		}
		var av *object.Avatar
		// Search players first.
		for _, p := range c.hud.playerAvatars {
			if p.ID() == char.ID() && p.Serial() == char.Serial() {
				av = p
			}
		}
		if av == nil {
			avData := res.Avatar(char.ID())
			if avData == nil {
				defData := data.DefaultAvatarData(char)
				res.SetAvatars(append(res.Avatars(), defData))
				avData = &defData
			}
			gameChar := c.hud.game.Char(char.ID(), char.Serial())
			if gameChar == nil {
				return
			}
			av = object.NewAvatar(gameChar, avData)
		}
		c.avatars.Store(char.ID()+char.Serial(), av)
	}
	for _, ob := range c.area.Objects() {
		ob, isOb := ob.(*flameob.Object)
		if !isOb {
			continue
		}
		_, obExists := c.objects.Load(ob.ID() + ob.Serial())
		if !obExists {
			continue
		}
		ogData := res.Object(ob.ID())
		if ogData == nil {
			log.Err.Printf("hud camera: update area objects: object data not found: %s",
				ob.ID())
			continue
		}
		og := object.NewObjectGraphic(ob, ogData)
		c.objects.Store(ob.ID()+ob.Serial(), og)
	}
}

// clearAreaObjecs removes avatars & graphics without
// corresponding objects in current area.
func (c *Camera) clearAreaObjects() {
	for _, av := range c.Avatars() {
		found := false
		for _, ob := range c.area.Objects() {
			_, isChar := ob.(*character.Character)
			if isChar && ob.ID() == av.ID() && ob.Serial() == av.Serial() {
				found = true
				break
			}
		}
		if found {
			continue
		}
		c.avatars.Delete(av.ID() + av.Serial())
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
		c.objects.Delete(gob.ID() + gob.Serial())
	}
}

// clear removes removes current area and all objects.
func (c *Camera) clear() {
	c.area = nil
	c.areaMap = nil
	c.avatars = new(sync.Map)
	c.objects = new(sync.Map)
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
		c.hud.Game().ActivePlayerChar().SetTarget(av.Character)
		return
	}
	for _, ob := range c.AreaObjects() {
		if !ob.DrawArea().Contains(pos) {
			continue
		}
		log.Dbg.Printf("hud: set target: %s", ob.ID()+"_"+ob.Serial())
		c.hud.Game().ActivePlayerChar().SetTarget(ob.Object)
		return
	}
	c.hud.Game().ActivePlayerChar().SetTarget(nil)
}

// Triggered after left mouse button was pressed.
func (c *Camera) onMouseLeftPressed(pos pixel.Vec) {
	if c.hud.containsPos(pos) {
		return
	}
	pc := c.hud.PCAvatar()
	// Action.
	for _, ob := range c.AreaObjects() {
		if !ob.DrawArea().Contains(pos) || !ob.Live() || ob.UseAction() == nil {
			continue
		}
		// Range check.
		r := math.Hypot(ob.Position().X-pc.Position().X, ob.Position().Y-pc.Position().Y)
		if r > ActionRange {
			pc.PrivateLog().Add(objects.Message{Text: "tar_too_far"})
			continue
		}
		log.Dbg.Printf("hud: action: %s#%s", ob.ID(), ob.Serial())
		pc.Use(ob.Object)
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
			pc.PrivateLog().Add(objects.Message{Text: "tar_too_far"})
			continue
		}
		// Show loot window.
		log.Dbg.Printf("hud: loot: %s#%s", av.ID(), av.Serial())
		c.hud.loot.SetTarget(av)
		c.hud.loot.Show()
		return
	}
	for _, ob := range c.AreaObjects() {
		if !ob.DrawArea().Contains(pos) || ob.Live() {
			continue
		}
		// Range check.
		r := math.Hypot(ob.Position().X-pc.Position().X, ob.Position().Y-pc.Position().Y)
		if r > LootRange {
			pc.PrivateLog().Add(objects.Message{Text: "tar_too_far"})
			continue
		}
		// Show loot window.
		log.Dbg.Printf("hud: loot: %s#%s", ob.ID(), ob.Serial())
		c.hud.loot.SetTarget(ob)
		c.hud.loot.Show()
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
			pc.PrivateLog().Add(objects.Message{Text: "tar_too_far"})
			continue
		}
		// Show dialog window.
		log.Dbg.Printf("hud: dialog: %s#%s", av.ID(), av.Serial())
		dialog := av.Dialogs()[0]
		c.hud.dialog.SetDialog(dialog)
		c.hud.dialog.Show()
	}
	// Move active PC.
	destPos := c.ConvCameraPos(pos)
	if !c.hud.game.Pause && c.PassablePosition(destPos) {
		c.hud.Game().ActivePlayerChar().SetDestPoint(destPos.X, destPos.Y)
	}
}
