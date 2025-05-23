/*
 * game.go
 *
 * Copyright 2020-2025 Dariusz Sikora <ds@isangeles.dev>
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
	"time"

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/dialog"
	"github.com/isangeles/flame/flag"
	"github.com/isangeles/flame/item"

	"github.com/isangeles/fire/request"

	"github.com/isangeles/ignite/ai"

	"github.com/isangeles/mural/log"
)

const (
	aiCharFlag = flag.Flag("igniteNpc")
)

// Wrapper struct for game.
type Game struct {
	*flame.Module
	chars              map[string]*Character
	playerChars        []*Character
	activePC           *Character
	server             *Server
	localAI            *ai.AI
	closing            bool
	pause              bool
	onPlayerCharChange func(c *Character)
}

// New creates new wrapper for specified module.
func New(module *flame.Module) *Game {
	g := Game{
		Module: module,
		chars:  make(map[string]*Character),
	}
	g.localAI = ai.New(ai.NewGame(module))
	return &g
}

// Update updates game.
// It executes game update loop until game is closed.
func (g *Game) Update() {
	update := time.Now()
	for !g.closing {
		if g.pause {
			continue
		}
		delta := time.Since(update).Milliseconds()
		g.Module.Update(delta)
		if g.Server() == nil {
			g.updateAIChars()
			g.localAI.Update(delta)
		}
		g.updateChars()
		update = time.Now()
		time.Sleep(time.Duration(16) * time.Millisecond)
	}
}

// Char returns game character with specified ID and serial.
func (g *Game) Char(id, serial string) *Character {
	return g.chars[id+serial]
}

// AddPlayerChar adds specified character to player characters list.
func (g *Game) AddPlayerChar(char *Character) error {
	g.playerChars = append(g.playerChars, char)
	g.SetActivePlayerChar(char)
	return nil
}

// PlayerChars returns all player characters.
func (g *Game) PlayerChars() []*Character {
	return g.playerChars
}

// SetActivePlayer sets specified avatar as active player avatar.
func (g *Game) SetActivePlayerChar(char *Character) {
	g.activePC = char
	if g.onPlayerCharChange != nil {
		g.onPlayerCharChange(char)
	}
}

// ActivePlayerChar returns active player character.
func (g *Game) ActivePlayerChar() *Character {
	return g.activePC
}

// Pause checks if the game pause is active.
func (g *Game) Pause() bool {
	return g.pause
}

// SetPause pauses/unpauses the game.
// If game uses the remote server this function will only send
// the pause request to the server and return, changing the game
// pause variable value should happend while handling the server
// response in such a case.
func (g *Game) SetPause(pause bool) {
	if g.server != nil {
		req := request.Request{Pause: pause}
		err := g.server.Send(req)
		if err != nil {
			log.Err.Printf("Game: unable to send pause request: %v", err)
		}
		return
	}
	g.pause = pause
	log.Inf.Print(pauseMessage(g.pause))
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

// SetOnPlayerChangeFunc sets function triggered on active player change.
func (g *Game) SetOnPlayerCharChangeFunc(f func(c *Character)) {
	g.onPlayerCharChange = f
}

// SpawnChar sets start area and position of current chapter for specified
// character.
func (g *Game) SpawnChar(char *Character) error {
	// Set start position.
	char.SetPosition(g.Chapter().Conf().StartPosX,
		g.Chapter().Conf().StartPosY)
	// Set start area.
	startArea := g.Chapter().Area(g.Chapter().Conf().StartArea)
	if startArea == nil {
		return fmt.Errorf("chapter start area not found: %s",
			g.Chapter().Conf().StartArea)
	}
	startArea.AddObject(char.Character)
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
		to.Inventory().AddItem(i)
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
	if !fairTrade(sellItems, buyItems) {
		return
	}
	for _, it := range sellItems {
		buyer.Inventory().RemoveItem(it)
		seller.Inventory().AddItem(it)
	}
	for _, it := range buyItems {
		seller.Inventory().RemoveItem(it)
		buyer.Inventory().AddItem(it)
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

// EndDialog ends specified dialog.
func (g *Game) EndDialog(dialog *dialog.Dialog) {
	defer dialog.SetTarget(nil)
	if g.Server() == nil || dialog.Owner() == nil {
		return
	}
	dialogReq := request.DialogEnd{
		TargetID:     dialog.Target().ID(),
		TargetSerial: dialog.Target().Serial(),
		OwnerID:      dialog.Owner().ID(),
		OwnerSerial:  dialog.Owner().Serial(),
		DialogID:     dialog.ID(),
	}
	req := request.Request{DialogEnd: []request.DialogEnd{dialogReq}}
	err := g.Server().Send(req)
	if err != nil {
		log.Err.Printf("Game: end dialog: unable to send end dialog request: %v",
			err)
	}
}

// AnswerDialog answers dialog with specified answer.
func (g *Game) AnswerDialog(dialog *dialog.Dialog, answer *dialog.Answer) {
	if g.Server() != nil && dialog.Owner() != nil && dialog.Target() != nil {
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
	dialog.Next(answer)
}

// VisibleForPlayer checks whether specified position is
// in visibility range of any PC.
func (g *Game) VisibleForPlayer(x, y float64) bool {
	for _, pc := range g.PlayerChars() {
		if pc.InSight(x, y) {
			return true
		}
	}
	return false
}

// updateAIChars updates list of characters controlled by the AI.
func (g *Game) updateAIChars() {
outer:
	for _, c := range g.chars {
		for _, aic := range g.localAI.Game().Characters() {
			if aic.ID() == c.ID() && aic.Serial() == c.Serial() {
				continue outer
			}
		}
		if !c.HasFlag(aiCharFlag) {
			continue
		}
		aiChar := ai.NewCharacter(c.Character, g.localAI.Game())
		g.localAI.Game().AddCharacter(aiChar)
	}
}

// updateChars updates list of game characters.
func (g *Game) updateChars() {
	for _, p := range g.playerChars {
		gameChar := g.chars[p.ID()+p.Serial()]
		if gameChar != nil {
			continue
		}
		g.chars[p.ID()+p.Serial()] = p
	}
	for _, c := range g.Chapter().Characters() {
		gameChar := g.chars[c.ID()+c.Serial()]
		if gameChar != nil {
			continue
		}
		gameChar = NewCharacter(c, g)
		g.chars[gameChar.ID()+gameChar.Serial()] = gameChar
	}
}

// fairTrade checks if value of all items to sell is greater or
// equal to the value of items to buy.
func fairTrade(sell, buy []item.Item) bool {
	buyValue := 0
	for _, it := range buy {
		buyValue += it.Value()
	}
	sellValue := 0
	for _, it := range sell {
		sellValue += it.Value()
	}
	return sellValue >= buyValue
}

// pauseMessage returns the pause message depending on the pause
// value(true = game pause message, false = game unpaused message).
func pauseMessage(pause bool) string {
	id := "game_paused"
	if !pause {
		id = "game_unpaused"
	}
	return lang.Text(id)
}
