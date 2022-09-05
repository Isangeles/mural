/*
 * response.go
 *
 * Copyright 2020-2022 Dariusz Sikora <ds@isangeles.dev>
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
	"sync"

	flameres "github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/serial"
	"github.com/isangeles/flame/useaction"

	"github.com/isangeles/fire/response"

	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/log"
)

var addPlayerMutex sync.Mutex

// handleResponse handles specified response from Fire server.
func (g *Game) handleResponse(resp response.Response) {
	g.handleUpdateResponse(resp.Update)
	for _, r := range resp.Character {
		g.handleCharacterResponse(r)
	}
	for _, r := range resp.Use {
		g.handleUseResponse(r)
	}
	for _, r := range resp.Command {
		log.Inf.Printf("[%d]: %s", r.Result, r.Out)
	}
	for _, r := range resp.Error {
		log.Err.Printf("Game server: error response: %s", r)
	}
}

// handleUpdateResponse handles update response.
func (g *Game) handleUpdateResponse(resp response.Update) {
	flameres.Clear()
	flameres.TranslationBases = res.TranslationBases()
	g.Apply(resp.Module)
}

// handleCharacterResponse handles new characters from server response.
func (g *Game) handleCharacterResponse(resp response.Character) {
	addPlayerMutex.Lock()
	defer addPlayerMutex.Unlock()
	for _, p := range g.PlayerChars() {
		if p.ID() == resp.ID && p.Serial() == resp.Serial {
			return
		}
	}
	gameChar := g.Char(resp.ID, resp.Serial)
	if gameChar != nil {
		g.AddPlayerChar(gameChar)
		return
	}
	char := g.Chapter().Character(resp.ID, resp.Serial)
	if char == nil {
		log.Err.Printf("Game: character from new char response not found in current module: %s %s",
			resp.ID, resp.Serial)
		return
	}
	gameChar = NewCharacter(char, g)
	g.AddPlayerChar(gameChar)
}

// handleUseResponse handles use response.
func (g *Game) handleUseResponse(resp response.Use) {
	char := g.Char(resp.UserID, resp.UserSerial)
	if char == nil {
		return
	}
	if char.onUse == nil {
		return
	}
	usable := char.Usable(resp.ObjectID, resp.ObjectSerial)
	if usable == nil {
		// Search for item or area object.
		ob := serial.Object(resp.ObjectID, resp.ObjectSerial)
		if ob == nil {
			log.Err.Printf("Object not found: %s %s", resp.ObjectID,
				resp.ObjectSerial)
			return
		}
		u, ok := ob.(useaction.Usable)
		if !ok {
			log.Err.Printf("Object is not usable: %s %s", resp.ObjectID,
				resp.ObjectSerial)
			return
		}
		usable = u
	}
	char.onUse(usable)
}
