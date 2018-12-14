/*
 * hud.go
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

// Package with 'head-up display' elements.
package hud

import (
	"fmt"
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	
	flamecore "github.com/isangeles/flame/core"
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/scenario"

	"github.com/isangeles/mural/core/areamap"
	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/objects"
)

var (
	main_color   color.Color = colornames.Grey
	sec_color    color.Color = colornames.Blue
	accent_color color.Color = colornames.Red
)

// Struct for 'head-up display'.
type HUD struct {
	loadScreen *LoadingScreen
	camera     *Camera
	chat       *Chat

	game    *flamecore.Game
	pc      *objects.Avatar
	destPos pixel.Vec
	loading bool
	exitReq bool
	loaderr error
}

// NewHUD creates new HUD instance.
// HUD loads all game data and setup specified
// PC area as current HUD area.
func NewHUD(g *flamecore.Game, pc *objects.Avatar) (*HUD, error) {
	hud := new(HUD)
	hud.game = g
	hud.pc = pc
	hud.loadScreen = newLoadingScreen(hud)
	hud.camera = newCamera(hud, config.Resolution())
	hud.chat = newChat(hud)
	go hud.LoadGame(g)
	return hud, nil
}

// Draw draws HUD elements.
func (hud *HUD) Draw(win *mtk.Window) {
	if hud.loading {
		hud.loadScreen.Draw(win)
	} else {
		hud.camera.Draw(win)
		hud.chat.Draw(win)
	}
}

// Update updated HUD elements.
func (hud *HUD) Update(win *mtk.Window) {
	if hud.exitReq {
		// TODO: exit back to menu.
		win.SetClosed(true)
	}
	if hud.loading {
		if hud.loaderr != nil {
			log.Err.Printf("loading_fail:%v", hud.loaderr)
			hud.Exit()
		}
	}
	// Key events.
	if win.JustPressed(pixelgl.KeyGraveAccent) {
		if !hud.chat.Active() {
			hud.chat.SetActive(true)
			hud.camera.Lock(true)
			
		} else {
			hud.chat.SetActive(false)
			hud.camera.Lock(false)
		}
	}
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		hud.destPos = hud.camera.ConvCameraPos(win.MousePosition())
		hud.Player().SetPosition(hud.destPos.X, hud.destPos.Y)
	}
	hud.loadScreen.Update(win)
	hud.camera.Update(win)
	hud.chat.Update(win)
}

// Camera position returns current position of
// HUD camera.
func (hud *HUD) CameraPosition() pixel.Vec {
	return hud.camera.Position()
}

// Player returns player character
// (game character of HUD user).
func (hud *HUD) Player() *objects.Avatar {
	return hud.pc
}

// AreaAvatars returns all avatars from current
// HUD area.
func (hud *HUD) AreaAvatars() []*objects.Avatar {
	return hud.camera.Avatars()
}

// Exit sends exit request to HUD.
func (hud *HUD) Exit() {
	hud.exitReq = true
}

// Chat returns HUD chat.
func (hud *HUD) Chat() *Chat {
	return hud.chat
}

// Game returns HUD current game.
func (hud *HUD) Game() *flamecore.Game {
	return hud.game
}

// LoadNewGame load all game data.
func (hud *HUD) LoadGame(game *flamecore.Game) {
	hud.loading = true
	hud.loadScreen.SetLoadInfo(lang.Text("gui", "load_game_data_info"))
	data.LoadGameData()
	pcArea, err := hud.game.PlayerArea(hud.Player().SerialID())
	if err != nil {
		hud.loaderr = fmt.Errorf("fail_to_retrieve_pc_area:%v", err)
		return
	}
	hud.ChangeArea(pcArea)
}

// ChangeArea changes current HUD area.
func (hud *HUD) ChangeArea(area *scenario.Area) {
	hud.loading = true
	// Map.
	hud.loadScreen.SetLoadInfo(lang.Text("gui", "load_map_info"))
	areaMap, err := areamap.NewMap(area, hud.game.Module().Chapter().AreasPath())
	if err != nil {
		hud.loaderr = fmt.Errorf("fail_to_create_pc_area_map:%v", err)
		return
	}
	hud.camera.SetMap(areaMap)
	// Objects.
	hud.loadScreen.SetLoadInfo(lang.Text("gui", "load_avatars_info"))
	avatars := make([]*objects.Avatar, 0)
	for _, c := range area.Characters() {
		if c == hud.pc.Character { // skip player, PC already has avatar
			avatars = append(avatars, hud.pc)
			continue
		}
		av, err := data.CharacterAvatar(
			hud.game.Module().Chapter().NPCPath(), c)
		if err != nil {
			log.Err.Printf("hud_area_change:char:%s:fail_to_retrieve_avatar:%v",
				c.ID(), err)
			continue
		}
		avatars = append(avatars, av)
	}
	hud.camera.SetAvatars(avatars)
	hud.loading = false
}
