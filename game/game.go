/*
 * game.go
 *
 * Copyright 2020-2021 Dariusz Sikora <dev@isangeles.pl>
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
	"github.com/isangeles/flame/module/flag"

	"github.com/isangeles/ignite/ai"

	"github.com/isangeles/mural/core/object"
)

const (
	aiCharFlag = flag.Flag("igniteNpc")
)

// Wrapper struct for game.
type Game struct {
	*flame.Game
	players      []*Player
	localAI      *ai.AI
	activePlayer *Player
	Pause        bool
}

// New creates new wrapper for specified game.
func New(game *flame.Game) *Game {
	g := Game{Game: game}
	g.localAI = ai.New(ai.NewGame(game))
	return &g
}

// Update updates game.
func (g *Game) Update(delta int64) {
	if g.Pause {
		return
	}
	g.Game.Update(delta)
	g.updateAIChars()
	g.localAI.Update(delta)
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

// updateAIChars updates list of characters controlled by the AI.
func (g *Game) updateAIChars() {
outer:
	for _, c := range g.Module().Chapter().Characters() {
		for _, aic := range g.localAI.Game().Characters() {
			if aic.ID() == c.ID() && aic.Serial() == c.Serial() {
				continue outer
			}
		}
		if !c.HasFlag(aiCharFlag) {
			continue
		}
		aiChar := ai.NewCharacter(c, g.localAI.Game())
		g.localAI.Game().AddCharacter(aiChar)
	}
}
