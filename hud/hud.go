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
	
	flamecore "github.com/isangeles/flame/core"
	"github.com/isangeles/flame/core/module/object/character"

	"github.com/isangeles/mural/core/areamap"
	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/data"
)

// Struct for 'head-up display'.
type HUD struct {
	camera *Camera

	game *flamecore.Game
	pc   *character.Character
}

// NewHUD creates new HUD instance.
func NewHUD(g *flamecore.Game, pc *character.Character) (*HUD, error) {
	hud := new(HUD)
	hud.camera = newCamera(config.Resolution())
	pcArea, err := g.PlayerArea(pc.Id())
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_player_area:%v", err)
	}
	m, err := data.Map(pcArea.Id(), g.Module().Chapter().AreasPath())
	if err != nil {
		return nil, fmt.Errorf("fail_to_load_pc_area_map:%v", err)
	}
	areaMap := areamap.NewMap(pcArea, m)
	hud.camera.SetMap(areaMap)
	hud.game = g
	hud.pc = pc
	return hud, nil
}

// Draw draws HUD elements.
func (hud *HUD) Draw(win *mtk.Window) {
	hud.camera.Draw(win)
}
