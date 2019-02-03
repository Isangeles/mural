/*
 * hud.go
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

// Package with 'head-up display' elements.
package hud

import (
	"fmt"
	"image/color"
	"path/filepath"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame"
	flamecore "github.com/isangeles/flame/core"
	flamedata "github.com/isangeles/flame/core/data"
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/scenario"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/areamap"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/exp"
	"github.com/isangeles/mural/core/data/imp"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

var (
	// HUD colors.
	main_color   color.Color = colornames.Grey
	sec_color    color.Color = colornames.Blue
	accent_color color.Color = colornames.Red
)

// Struct for 'head-up display'.
type HUD struct {
	loadScreen *LoadingScreen
	camera     *Camera
	bar        *MenuBar
	menu       *Menu
	pcFrame    *CharFrame
	chat       *Chat
	game       *flamecore.Game
	pcs        []*object.Avatar
	activePC   *object.Avatar
	destPos    pixel.Vec
	loading    bool
	exiting    bool
	loaderr    error
}

// NewHUD creates new HUD instance.
// HUD loads all game data and setup specified
// PC area as current HUD area.
func NewHUD(g *flamecore.Game, pcs []*object.Avatar) (*HUD, error) {
	hud := new(HUD)
	hud.game = g
	if len(pcs) < 1 {
		return nil, fmt.Errorf("no player characters")
	}
	// Players.
	hud.pcs = pcs
	hud.activePC = hud.pcs[0]
	// Loading screen.
	hud.loadScreen = newLoadingScreen(hud)
	// Camera.
	hud.camera = newCamera(hud, config.Resolution())
	hud.camera.CenterAt(hud.ActivePlayer().Position())
	// Bar.
	hud.bar = newMenuBar(hud)
	// Menu.
	hud.menu = newMenu(hud)
	// Active player character frame.
	pcFrame, err := newCharFrame(hud, hud.ActivePlayer())
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_active_pc_frame:%v",
			err)
	}
	hud.pcFrame = pcFrame
	// Chat window.
	hud.chat = newChat(hud)
	// Start game loading.
	go hud.LoadGame(g)
	return hud, nil
}

// Draw draws HUD elements.
func (hud *HUD) Draw(win *mtk.Window) {
	if hud.loading {
		hud.loadScreen.Draw(win)
		return
	}
	// Elements positions.
	pcFramePos := mtk.DrawPosTL(win.Bounds(), hud.pcFrame.Bounds())
	barPos := mtk.DrawPosBC(win.Bounds(), hud.bar.Bounds())
	// Draw elements.
	hud.camera.Draw(win)
	hud.bar.Draw(win, mtk.Matrix().Moved(barPos))
	hud.chat.Draw(win)
	hud.pcFrame.Draw(win, mtk.Matrix().Moved(pcFramePos))
	if hud.menu.Opened() {
		menuPos := win.Bounds().Center()
		hud.menu.Draw(win, mtk.Matrix().Moved(menuPos))
	}
}

// Update updated HUD elements.
func (hud *HUD) Update(win *mtk.Window) {
	if hud.exiting {
		// TODO: exit back to menu.
		win.SetClosed(true)
	}
	if hud.loading {
		if hud.loaderr != nil { // on loading error
			log.Err.Printf("hud_loading_fail:%v", hud.loaderr)
			hud.Exit()
		}
	}
	// Key events.
	if win.JustPressed(pixelgl.KeyGraveAccent) {
		// Toggle chat activity.
		if !hud.chat.Activated() {
			hud.chat.Active(true)
			hud.camera.Lock(true)

		} else {
			hud.chat.Active(false)
			hud.camera.Lock(false)
		}
	}
	if !hud.chat.Activated() { // block rest of key events if chat is active
		if win.JustPressed(pixelgl.KeySpace) {
			// Pause game.
			if !hud.game.Paused() {
				hud.game.Pause(true)
				hud.camera.Lock(true)
			} else {
				hud.game.Pause(false)
				hud.camera.Lock(false)
			}
		}
	}
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		destPos := hud.camera.ConvCameraPos(win.MousePosition())
		if !hud.game.Paused() && hud.camera.Map().Passable(destPos) {
			// Move active PC.
			hud.destPos = destPos
			hud.ActivePlayer().SetDestPoint(hud.destPos.X, hud.destPos.Y)
		}
	}
	hud.loadScreen.Update(win)
	hud.camera.Update(win)
	hud.bar.Update(win)
	hud.chat.Update(win)
	hud.pcFrame.Update(win)
	hud.menu.Update(win)
}

// Camera position returns current position of
// HUD camera.
func (hud *HUD) CameraPosition() pixel.Vec {
	return hud.camera.Position()
}

// Players returns player characters
// (game characters of HUD user).
func (hud *HUD) Players() []*object.Avatar {
	return hud.pcs
}

// ActivePlayers retruns currently active PC.
func (hud *HUD) ActivePlayer() *object.Avatar {
	return hud.activePC
}

// AreaAvatars returns all avatars from current
// HUD area.
func (hud *HUD) AreaAvatars() []*object.Avatar {
	return hud.camera.Avatars()
}

// Exit sends exit request to HUD.
func (hud *HUD) Exit() {
	hud.exiting = true
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
	err := imp.LoadResources()
	if err != nil {
		hud.loaderr = fmt.Errorf("fail_to_load_resources:%v", err)
		return
	}
	pcArea, err := hud.game.PlayerArea(hud.ActivePlayer().SerialID())
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
	// TODO: sometimes SetLoadInfo causes panic on load text draw.
	//hud.loadScreen.SetLoadInfo(lang.Text("gui", "load_map_info"))
	mapsPath := filepath.FromSlash(hud.game.Module().Chapter().AreasPath() +
		"/maps")
	tmxMap, err := data.Map(mapsPath, area.ID())
	if err != nil {
		hud.loaderr = fmt.Errorf("fail_to_retrieve_tmx_map:%v", err)
		return
	}
	areaMap, err := areamap.NewMap(tmxMap, mapsPath)
	if err != nil {
		hud.loaderr = fmt.Errorf("fail_to_create_pc_area_map:%v", err)
		return
	}
	hud.camera.SetMap(areaMap)
	// Objects.
	hud.loadScreen.SetLoadInfo(lang.Text("gui", "load_avatars_info"))
	avatars := make([]*object.Avatar, 0)
	npcPath := hud.Game().Module().Chapter().NPCPath()
	for _, c := range area.Characters() {
		var pcAvatar *object.Avatar
		for _, pc := range hud.Players() {
			if c == pc.Character {
				pcAvatar = pc
				break
			}
		}
		if pcAvatar != nil { // skip players, PCs already have avatars
			avatars = append(avatars, pcAvatar)
			continue
		}
		avData, err := imp.CharacterAvatarData(c, npcPath)
		if err != nil {
			log.Err.Printf("hud_area_change:char:%s:fail_to_retrieve_avatar:%v",
				c.ID(), err)
			continue
		}
		av := object.NewAvatar(avData)
		avatars = append(avatars, av)
	}
	hud.camera.SetAvatars(avatars)
	hud.loading = false
}

// Save saves GUI and game state to
// savegames directory.
func (hud *HUD) Save(saveName string) error {
	// Retrieve saves path.
	savesPath := flame.SavegamesPath()
	// Save current game.
	err := flamedata.SaveGame(hud.Game(), savesPath, saveName)
	if err != nil {
		return fmt.Errorf("fail_to_save_game:%v",
			err)
	}
	// Save GUI state.
	guisav := hud.NewGUISave()
	err = exp.ExportGUISave(guisav, savesPath, saveName)
	if err != nil {
		return fmt.Errorf("fail_to_save_gui:%v",
			err)
	}
	return nil
}

// Saves GUI to save struct.
func (hud *HUD) NewGUISave() *res.GUISave {
	sav := new(res.GUISave)
	// Save players avatars.
	for _, pc := range hud.Players() {
		avData := res.AvatarData{
			Character:      pc.Character,
			PortraitName:   pc.PortraitName(),
			SSHeadName:     pc.HeadSpritesheetName(),
			SSTorsoName:    pc.TorsoSpritesheetName(),
			SSFullBodyName: pc.FullBodySpritesheetName(),
		}
		sav.PlayersData = append(sav.PlayersData, &avData)
	}
	// Save camera XY position.
	sav.CameraPosX = hud.CameraPosition().X
	sav.CameraPosY = hud.CameraPosition().Y
	return sav
}

// LoadGUISave load specified saved GUI state.
func (hud *HUD) LoadGUISave(save *res.GUISave) error {
	// Players.
	for _, pcData := range save.PlayersData {
		pc := object.NewAvatar(pcData)
		hud.pcs = append(hud.pcs, pc)
	}
	// Camera position.
	hud.camera.SetPosition(pixel.V(save.CameraPosX, save.CameraPosY))
	return nil
}
