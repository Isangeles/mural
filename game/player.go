/*
 * player.go
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
	"fmt"

	"github.com/isangeles/flame/character"
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/effect"
	"github.com/isangeles/flame/item"
	"github.com/isangeles/flame/objects"
	"github.com/isangeles/flame/req"
	"github.com/isangeles/flame/serial"
	"github.com/isangeles/flame/useaction"

	"github.com/isangeles/fire/request"

	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/object"
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
func (p *Player) AddChatMessage(message objects.Message) {
	p.ChatLog().Add(message)
	if p.game.Server() == nil {
		return
	}
	chatReq := request.Chat{p.ID(), p.Serial(), message.String(), true}
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
	if p.Casted() != nil {
		return
	}
	err := p.Character.Use(ob)
	if err != nil {
		p.PrivateLog().Add(objects.Message{Text: "cant_do_right_now"})
		if !p.meetTargetRangeReqs(ob.UseAction().Requirements()...) {
			tar := p.Targets()[0]
			tarPosX, tarPosY := tar.Position()
			p.moveCloseTo(tarPosX, tarPosY, ob.UseAction().MinRange())
		}
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

// Equip inserts specified equipable item to all
// compatible slots in active PC equipment.
func (p *Player) Equip(it item.Equiper) error {
	if !p.MeetReqs(it.EquipReqs()...) {
		return fmt.Errorf(lang.Text("reqs_not_meet"))
	}
	slots := make([]*character.EquipmentSlot, 0)
	for _, itSlot := range it.Slots() {
		equiped := false
		for _, eqSlot := range p.Equipment().Slots() {
			if eqSlot.Item() != nil {
				continue
			}
			if eqSlot.Type() == itSlot {
				eqSlot.SetItem(it)
				equiped = true
				slots = append(slots, eqSlot)
				break
			}
		}
		if !equiped {
			p.Equipment().Unequip(it)
			return fmt.Errorf(lang.Text("equip_no_free_slot_error"))
		}
	}
	if !p.Equipment().Equiped(it) {
		return fmt.Errorf(lang.Text("equip_no_valid_slot_error"))
	}
	if p.game.Server() == nil {
		return nil
	}
	eqReq := request.Equip{
		CharID:     p.ID(),
		CharSerial: p.Serial(),
		ItemID:     it.ID(),
		ItemSerial: it.Serial(),
	}
	for _, s := range slots {
		slotReq := request.EquipmentSlot{
			Type: string(s.Type()),
			ID:   s.ID(),
		}
		eqReq.Slots = append(eqReq.Slots, slotReq)
	}
	req := request.Request{Equip: []request.Equip{eqReq}}
	err := p.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Player: %s %s: unable to send equip request: %v",
			p.ID(), p.Serial(), err)
	}
	return nil
}

// Unequip removes specified item from player equipment.
func (p *Player) Unequip(it item.Equiper) {
	p.Equipment().Unequip(it)
	if p.game.Server() == nil {
		return
	}
	uneqReq := request.Unequip{
		CharID:     p.ID(),
		CharSerial: p.Serial(),
		ItemID:     it.ID(),
		ItemSerial: it.Serial(),
	}
	req := request.Request{Unequip: []request.Unequip{uneqReq}}
	err := p.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Player: %s %s: unable to send unequip request: %v",
			p.ID(), p.Serial(), err)
	}
}

// meetTargetRangeReqs check if all target range requirements are meet.
// Returns true, if none of specified requirements is a target range
// requirement.
func (p *Player) meetTargetRangeReqs(reqs ...req.Requirement) bool {
	tarRangeReqs := make([]req.Requirement, 0)
	for _, r := range reqs {
		if r, ok := r.(*req.TargetRange); ok {
			tarRangeReqs = append(tarRangeReqs, r)
		}
	}
	if !p.MeetReqs(tarRangeReqs...) {
		return false
	}
	return true
}

// moveCloseTo moves player to the position at minimal range
// to the specified position.
func (p *Player) moveCloseTo(x, y, minRange float64) {
	switch {
	case x > p.Position().X:
		x -= minRange
	case x < p.Position().X:
		x += minRange
	case y > p.Position().Y:
		y -= minRange
	case y < p.Position().Y:
		y += minRange
	}
	p.SetDestPoint(x, y)
}

