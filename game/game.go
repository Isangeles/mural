/*
 * game.go
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

// Package with game wrapper struct.
package game

import (
	"fmt"

	"github.com/faiface/pixel"

	"github.com/isangeles/flame"

	"github.com/isangeles/mural/core/object"
)

// Wrapper struct for game.
type Game struct {
	*flame.Game
	players      []*Player
	activePlayer *Player
}

// New creates new wrapper for specified game.
func New(game *flame.Game) *Game {
	g := Game{Game: game}
	return &g
}

// AddPlayer adds specified avatar to player avatars list.
func (g *Game) AddPlayer(avatar *Player) {
	g.players = append(g.players, avatar)
}

// Players returns all player avatars.
func (g *Game) Players() []*Player {
	return g.players
}

// SetActivePlayer sets specified avatar as active player avatar.
func (g *Game) SetActivePlayer(player *Player) {
	g.activePlayer = player
}

// ActivePlayer returns active player avatar.
func (g *Game) ActivePlayer() *Player {
	return g.activePlayer
}

// NewPlayer creates new character avatar for the player from specified data and
// places this character in the start area of game module.
func (g *Game) NewPlayer(avatar *object.Avatar) (*Player, error) {
	// Set start position.
	startPos := pixel.V(g.Module().Chapter().Conf().StartPosX,
		g.Module().Chapter().Conf().StartPosY)
	avatar.SetPosition(startPos)
	// Set start area.
	startArea := g.Module().Chapter().Area(g.Module().Chapter().Conf().StartArea)
	if startArea == nil {
		return nil, fmt.Errorf("game: start area not found: %s",
			g.Module().Chapter().Conf().StartArea)
	}
	startArea.AddCharacter(avatar.Character)
	player := Player{avatar, g}
	return &player, nil
}