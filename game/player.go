/*
 * player.go
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

package game

import (
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/module/effect"
	"github.com/isangeles/flame/module/serial"
	"github.com/isangeles/flame/module/useaction"

	"github.com/isangeles/fire/request"

	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// Struct for game player.
type Player struct {
	*object.Avatar
	game *Game
}

// NewPlayer creates game wrapper for player avatar.
func NewPlayer(avatar *object.Avatar, game *Game) *Player {
	player := Player{avatar, game}
	return &player
}

// SetDestPoint sets destination point for player character.
func (p *Player) SetDestPoint(x, y float64) {
	p.Character.SetDestPoint(x, y)
	if p.game.Server() == nil {
		return
	}
	moveReq := request.Move{p.ID(), p.Serial(), x, y}
	req := request.Request{Move: []request.Move{moveReq}}
	err := p.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Player: %s %s: unable to send move request to the server: %v",
			p.ID(), p.Serial(), err)
	}
}

// AddChatMessage adds new message to the character chat log.
func (p *Player) AddChatMessage(message string) {
	p.ChatLog().Add(message)
	if p.game.Server() == nil {
		return
	}
	chatReq := request.Chat{p.ID(), p.Serial(), message}
	req := request.Request{Chat: []request.Chat{chatReq}}
	err := p.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Player: %s %s: unable to send chat request to the server: %v",
			p.ID(), p.Serial(), err)
	}
}

// SetTarget sets specified targetable object as the current target.
func (p *Player) SetTarget(tar effect.Target) {
	p.Character.SetTarget(tar)
	if p.game.Server() == nil {
		return
	}
	targetReq := request.Target{
		ObjectID:     p.ID(),
		ObjectSerial: p.Serial(),
	}
	if tar != nil {
		targetReq.TargetID, targetReq.TargetSerial = tar.ID(), tar.Serial()
	}
	req := request.Request{Target: []request.Target{targetReq}}
	err := p.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Player: %s %s: unable to send target request to the server: %v",
			p.ID(), p.Serial(), err)
	}
}

// Use uses specified usable object.
func (p *Player) Use(ob useaction.Usable) {
	err := p.Character.Use(ob)
	if err != nil {
		p.PrivateLog().Add(lang.Text("cant_do_right_now"))
		return
	}
	if p.game.Server() == nil {
		return
	}
	useReq := request.Use{
		UserID:     p.ID(),
		UserSerial: p.Serial(),
		ObjectID:   ob.ID(),
	}
	if ob, ok := ob.(serial.Serialer); ok {
		useReq.ObjectSerial = ob.Serial()
	}
	req := request.Request{Use: []request.Use{useReq}}
	err = p.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Player: %s %s: unable to send use request: %v",
			p.ID(), p.Serial(), err)
	}
}
