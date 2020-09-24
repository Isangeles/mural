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
	server       *Server
	closing      bool
}

// New creates new wrapper for specified game.
func New(game *flame.Game, server *Server) (*Game, error) {
	g := Game{
		Game:   game,
		server: server,
	}
	if g.Server() != nil {
		g.Server().SetOnResponseFunc(g.handleResponse)
	}
	return &g, nil
}

// AddPlayer adds specified avatar to player avatars list.
func (g *Game) AddPlayer(avatar *object.Avatar) error {
	if g.Server() != nil {
		err := g.Server().NewCharacter(avatar.Character.Data())
		if err != nil {
			return fmt.Errorf("Unable to send new character request: %v",
				err)
		}
		return nil
	}
	player, err := g.newPlayer(avatar)
	if err != nil {
		return fmt.Errorf("Unable to create player: %v", err)
	}
	g.players = append(g.players, player)
	return nil
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

// Server returns game server connection.
func (g *Game) Server() *Server {
	return g.server
}

// Closing checks if game should be closed.
func (g *Game) Closing() bool {
	return g.closing
}

// newPlayer creates new character avatar for the player from specified data and
// places this character in the start area of game module.
func (g *Game) newPlayer(avatar *object.Avatar) (*Player, error) {
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
