/*
 * camera.go
 *
 * Copyright 2018-2025 Dariusz Sikora <ds@isangeles.dev>
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

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"

	"github.com/isangeles/flame/area"
	"github.com/isangeles/flame/character"
	"github.com/isangeles/flame/objects"

	"github.com/isangeles/stone"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/object"
)

var (
	FOWColor     = pixel.RGBA{0.1, 0.1, 0.1, 0.7}
	debugMoveKey = pixelgl.KeyLeftShift
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
	area     *object.Area
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
	c.area.Draw(win, mtk.Matrix().Moved(c.Position()), c.Size())
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
	c.size = win.Bounds().Size()
	// Key events.
	if !c.locked && c.area != nil {
		tileSize := c.area.Map().TileSize()
		offset := pixel.V(tileSize.X*16, tileSize.Y*16)
		mapSize := c.area.Map().Size()
		mapSize = pixel.V(mapSize.X+offset.X, mapSize.Y+offset.Y)
		// Key events.
		if c.position.Y < mapSize.Y && win.Pressed(pixelgl.KeyW) ||
			win.Pressed(pixelgl.KeyUp) {
			c.position.Y += tileSize.Y
		}
		if c.position.X < mapSize.X && win.Pressed(pixelgl.KeyD) ||
			win.Pressed(pixelgl.KeyRight) {
			c.position.X += tileSize.X
		}
		if c.position.Y > 0-offset.Y && win.Pressed(pixelgl.KeyS) ||
			win.Pressed(pixelgl.KeyDown) {
			c.position.Y -= tileSize.Y
		}
		if c.position.X > 0-offset.X && win.Pressed(pixelgl.KeyA) ||
			win.Pressed(pixelgl.KeyLeft) {
			c.position.X -= tileSize.X
		}
	}
	//Area.
	if c.area != nil {
		c.area.Update(win)
	}
	// Mouse events.
	if config.Debug && win.JustPressed(pixelgl.MouseButtonLeft) && win.Pressed(debugMoveKey) {
		c.onDebugMouseLeftPressed(win.MousePosition())
	} else if win.JustPressed(pixelgl.MouseButtonLeft) {
		c.onMouseLeftPressed(win.MousePosition())
	}
	if win.JustPressed(pixelgl.MouseButtonRight) {
		c.onMouseRightPressed(win.MousePosition())
	}
	// Debug.
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
	if a == nil {
		a = nil
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
	c.area = object.NewArea(c.hud.game, a, areaMap)
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

// Area retuns current area.
func (c *Camera) Area() *object.Area {
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

// Triggered after right mouse button was pressed.
func (c *Camera) onMouseRightPressed(pos pixel.Vec) {
	// Set target.
	if c.hud.containsPos(pos) {
		return
	}
	for _, av := range c.area.Avatars() {
		if !av.DrawArea().Contains(pos) {
			continue
		}
		log.Dbg.Printf("hud: set target: %s", av.ID()+"_"+av.Serial())
		c.hud.Game().ActivePlayerChar().SetTarget(av.Character)
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
	for _, ob := range c.area.Avatars() {
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
		pc.Use(ob.Character)
		return
	}
	// Loot.
	for _, av := range c.area.Avatars() {
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
	// Dialog.
	for _, av := range c.area.Avatars() {
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
		dialog := av.Dialog(pc)
		c.hud.dialog.SetDialog(dialog)
		c.hud.dialog.Show()
	}
	// Move active PC.
	destPos := c.ConvCameraPos(pos)
	if !c.hud.game.Pause() && c.area.PassablePosition(destPos) {
		c.hud.Game().ActivePlayerChar().SetDestPoint(destPos.X, destPos.Y)
	}
}

// Triggered after pressing left mouse button with move debug key.
func (c *Camera) onDebugMouseLeftPressed(pos pixel.Vec) {
	movePos := c.ConvCameraPos(pos)
	c.hud.Game().ActivePlayerChar().SetPosition(movePos.X, movePos.Y)
}
