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
	"github.com/isangeles/flame/module/item"

	"github.com/isangeles/mural/core/object"
)

// Wrapper struct for game.
type Game struct {
	*flame.Game
	players      []*Player
	activePlayer *Player
	server       *Server
	closing      bool
	onActivePlayerChange func(p *Player)
}

// New creates new wrapper for specified game.
func New(game *flame.Game) *Game {
	g := Game{Game: game}
	return &g
}

// AddPlayer adds specified avatar to player avatars list.
func (g *Game) AddPlayer(avatar *Player) error {
	g.players = append(g.players, avatar)
	g.SetActivePlayer(avatar)
	return nil
}

// Players returns all player avatars.
func (g *Game) Players() []*Player {
	return g.players
}

// SetActivePlayer sets specified avatar as active player avatar.
func (g *Game) SetActivePlayer(player *Player) {
	g.activePlayer = player
	if g.onActivePlayerChange != nil {
		g.onActivePlayerChange(player)
	}
}

// ActivePlayer returns active player avatar.
func (g *Game) ActivePlayer() *Player {
	return g.activePlayer
}

// SetServer sets game server.
func (g *Game) SetServer(server *Server) {
	g.server = server
	if g.Server() != nil {
		g.Server().SetOnResponseFunc(g.handleResponse)
	}
}

// Server returns game server connection.
func (g *Game) Server() *Server {
	return g.server
}

// Closing checks if game should be closed.
func (g *Game) Closing() bool {
	return g.closing
}

// SetOnActivePlayerChangeFunc sets function triggered on active player change.
func (g *Game) SetOnActivePlayerChangeFunc(f func(p *Player)) {
	g.onActivePlayerChange = f
}

// SpawnChar sets start area and position of current chapter for specified avatar.
func (g *Game) SpawnChar(avatar *object.Avatar) error {
	// Set start position.
	startPos := pixel.V(g.Module().Chapter().Conf().StartPosX,
		g.Module().Chapter().Conf().StartPosY)
	avatar.SetPosition(startPos)
	// Set start area.
	startArea := g.Module().Chapter().Area(g.Module().Chapter().Conf().StartArea)
	if startArea == nil {
		return fmt.Errorf("chapter start area not found: %s",
			g.Module().Chapter().Conf().StartArea)
	}
	startArea.AddCharacter(avatar.Character)
	return nil
}

// transferItems transfer items between specified objects.
// Items are in the form of a map with IDs as keys and serial values as values.
func (g *Game) TransferItems(from, to item.Container, items ...item.Item) error {
	for _, i := range items {
		if from.Inventory().Item(i.ID(), i.Serial()) == nil {
			return fmt.Errorf("Item not found: %s %s",
				i.ID(), i.Serial())
		}
		from.Inventory().RemoveItem(i)
		err := to.Inventory().AddItem(i)
		if err != nil {
			return fmt.Errorf("Unable to add item inventory: %v",
				err)
		}
	}
	return nil
}
