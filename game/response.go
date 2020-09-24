/*
 * response.go
 *
 * Copyright 2020 Dariusz Sikora <dev@isangeles.pl>
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

package game

import (
	"github.com/isangeles/flame"
	flameres "github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/module"
	"github.com/isangeles/flame/module/character"

	"github.com/isangeles/fire/response"

	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// handleResponse handles specified response from Fire server.
func (g *Game) handleResponse(resp response.Response) {
	g.handleUpdateResponse(resp.Update)
	g.handleNewCharResponse(resp.NewChar)
	for _, r := range resp.Error {
		log.Err.Printf("Game: server error response: %s", r)
	}
}

// handleUpdateResponse handles update response.
func (g *Game) handleUpdateResponse(resp response.Update) {
	if g.Game == nil {
		mod := module.New(resp.Module)
		g.Game = flame.NewGame(mod)
		return
	}
	g.Module().Apply(resp.Module)
}

// handleNewCharResponse handles new characters from server response.
func (g *Game) handleNewCharResponse(resp []flameres.CharacterData) {
	for _, cd := range resp {
		char := character.New(cd)
		avData := res.Avatar(char.ID())
		if avData == nil {
			log.Err.Printf("Game: no avatar data for new character: %s",
				char.ID())
			continue
		}
		avatar := object.NewAvatar(char, avData)
		player := Player{avatar, g}
		g.players = append(g.players, &player)
	}
}
