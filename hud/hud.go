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
	
	flamecore "github.com/isangeles/flame/core"
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/scenario"

	"github.com/isangeles/mural/core/areamap"
	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/mtk"
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

	game    *flamecore.Game
	pc      *objects.Avatar
	loading bool
}

// NewHUD creates new HUD instance.
func NewHUD(g *flamecore.Game, pc *objects.Avatar) (*HUD, error) {
	hud := new(HUD)
	hud.game = g
	hud.pc = pc
	hud.loadScreen = newLoadingScreen(hud)
	hud.camera = newCamera(hud, config.Resolution())
	pcArea, err := hud.game.PlayerArea(pc.Id())
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_pc_area:%v", err)
	}
	go hud.ChangeArea(pcArea, hud.game)
	return hud, nil
}

// Draw draws HUD elements.
func (hud *HUD) Draw(win *mtk.Window) {
	if hud.loading {
		hud.loadScreen.Draw(win)
	} else {
		hud.camera.Draw(win)
	}
}

// Update updated HUD elements.
func (hud *HUD) Update(win *mtk.Window) {
	hud.loadScreen.Update(win)
	hud.camera.Update(win)
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

// ChangeArea changes current HUD area.
func (hud *HUD) ChangeArea(area *scenario.Area, game *flamecore.Game) {
	hud.loading = true
	hud.loadScreen.SetLoadInfo(lang.Text("gui", "load_area_info"))
	areaMap, err := areamap.NewMap(area, game.Module().Chapter().AreasPath())
	if err != nil {
		log.Err.Printf("fail_to_create_pc_area_map:%v", err)
		hud.loading = false
		return
	}
	hud.camera.SetMap(areaMap)
	hud.loading = false
}
