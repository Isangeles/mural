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
	"github.com/isangeles/flame/module/dialog"
	"github.com/isangeles/flame/module/item"

	"github.com/isangeles/fire/request"

	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// Wrapper struct for game.
type Game struct {
	*flame.Game
	players              []*Player
	activePlayer         *Player
	server               *Server
	closing              bool
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

// TransferItems transfer items between specified objects.
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
	if g.Server() == nil {
		return nil
	}
	transferReq := request.TransferItems{
		ObjectFromID:     from.ID(),
		ObjectFromSerial: from.Serial(),
		ObjectToID:       to.ID(),
		ObjectToSerial:   to.Serial(),
		Items:            make(map[string][]string),
	}
	for _, i := range items {
		transferReq.Items[i.ID()] = append(transferReq.Items[i.ID()], i.Serial())
	}
	req := request.Request{TransferItems: []request.TransferItems{transferReq}}
	err := g.Server().Send(req)
	if err != nil {
		log.Err.Printf("Game: transfer items: unable to send transfer items request: %v",
			err)
	}
	return nil
}

// Trade exchanges items between specified containers.
func (g *Game) Trade(seller, buyer item.Container, sellItems, buyItems []item.Item) {
	for _, it := range sellItems {
		buyer.Inventory().RemoveItem(it)
		err := seller.Inventory().AddItem(it)
		if err  != nil {
			log.Err.Printf("Game: trade items: unable to add sell item: %s %s: %v",
				it.ID(), it.Serial(), err)
		}
	}
	for _, it := range buyItems {
		seller.Inventory().RemoveItem(it)
		err := buyer.Inventory().AddItem(it)
		if err != nil {
			log.Err.Printf("Game: trade items: unable to add buy item: %s %s: %v",
				it.ID(), it.Serial(), err)
		}
	}
	if g.Server() == nil {
		return
	}
	transferReqSell := request.TransferItems{
		ObjectFromID:     buyer.ID(),
		ObjectFromSerial: buyer.Serial(),
		ObjectToID:       seller.ID(),
		ObjectToSerial:   seller.Serial(),
		Items:            make(map[string][]string),
	}
	for _, i := range sellItems {
		transferReqSell.Items[i.ID()] = append(transferReqSell.Items[i.ID()], i.Serial())
	}
	transferReqBuy := request.TransferItems{
		ObjectFromID:     seller.ID(),
		ObjectFromSerial: seller.Serial(),
		ObjectToID:       buyer.ID(),
		ObjectToSerial:   buyer.Serial(),
		Items:            make(map[string][]string),
	}
	for _, i := range buyItems {
		transferReqBuy.Items[i.ID()] = append(transferReqBuy.Items[i.ID()], i.Serial())
	}
	tradeReq := request.Trade{Sell: transferReqSell, Buy: transferReqBuy}
	req := request.Request{Trade: []request.Trade{tradeReq}}
	err := g.Server().Send(req)
	if err != nil {
		log.Err.Printf("Game: trade items: unable to send trade request: %v",
			err)
	}
}

// StartDialog starts dialog with specified object as dialog target.
func (g *Game) StartDialog(dialog *dialog.Dialog, target dialog.Talker) {
	dialog.SetTarget(target)
	if g.Server() == nil || dialog.Owner() == nil {
		return
	}
	dialogReq := request.Dialog{
		TargetID:     target.ID(),
		TargetSerial: target.Serial(),
		OwnerID:      dialog.Owner().ID(),
		OwnerSerial:  dialog.Owner().Serial(),
		DialogID:     dialog.ID(),
	}
	req := request.Request{Dialog: []request.Dialog{dialogReq}}
	err := g.Server().Send(req)
	if err != nil {
		log.Err.Printf("Game: start dialog: unable to send dialog request: %v",
			err)
	}
}

// AnswerDialog answers dialog with specified answer.
func (g *Game) AnswerDialog(dialog *dialog.Dialog, answer *dialog.Answer) {
	dialog.Next(answer)
	if g.Server() == nil || dialog.Owner() == nil || dialog.Target() == nil {
		return
	}
	dialogReq := request.Dialog{
		TargetID:     dialog.Target().ID(),
		TargetSerial: dialog.Target().Serial(),
		OwnerID:      dialog.Owner().ID(),
		OwnerSerial:  dialog.Owner().Serial(),
		DialogID:     dialog.ID(),
	}
	dialogAnswerReq := request.DialogAnswer{
		Dialog:   dialogReq,
		AnswerID: answer.ID(),
	}
	req := request.Request{DialogAnswer: []request.DialogAnswer{dialogAnswerReq}}
	err := g.Server().Send(req)
	if err != nil {
		log.Err.Printf("Game: answer dialog: unable to send dialog answer: %v",
			err)
	}
}
