/*
 * character.go
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
	"math"

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

// Struct for game character.
type Character struct {
	*object.Avatar
	game       *Game
	privateLog *objects.Log
}

// NewCharacter creates game wrapper for avatar.
func NewCharacter(avatar *object.Avatar, game *Game) *Character {
	c := Character{
		Avatar:     avatar,
		game:       game,
		privateLog: objects.NewLog(),
	}
	return &c
}

// PrivateLog returns avatar private log.
func (c *Character) PrivateLog() *objects.Log {
	return c.privateLog
}

// InSight checks if specified XY position is in sight range
// of the character.
func (c *Character) InSight(x, y float64) bool {
	return math.Hypot(c.Position().X-x, c.Position().Y-y) <= c.SightRange()
}

// SetDestPoint sets destination point for player character.
func (c *Character) SetDestPoint(x, y float64) {
	c.Character.SetDestPoint(x, y)
	if c.game.Server() == nil {
		return
	}
	moveReq := request.Move{c.ID(), c.Serial(), x, y}
	req := request.Request{Move: []request.Move{moveReq}}
	err := c.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Character: %s %s: unable to send move request to the server: %v",
			c.ID(), c.Serial(), err)
	}
}

// AddChatMessage adds new message to the character chat log.
func (c *Character) AddChatMessage(message objects.Message) {
	c.ChatLog().Add(message)
	if c.game.Server() == nil {
		return
	}
	chatReq := request.Chat{c.ID(), c.Serial(), message.String(), true}
	req := request.Request{Chat: []request.Chat{chatReq}}
	err := c.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Character: %s %s: unable to send chat request to the server: %v",
			c.ID(), c.Serial(), err)
	}
}

// SetTarget sets specified targetable object as the current target.
func (c *Character) SetTarget(tar effect.Target) {
	c.Character.SetTarget(tar)
	if c.game.Server() == nil {
		return
	}
	targetReq := request.Target{
		ObjectID:     c.ID(),
		ObjectSerial: c.Serial(),
	}
	if tar != nil {
		targetReq.TargetID, targetReq.TargetSerial = tar.ID(), tar.Serial()
	}
	req := request.Request{Target: []request.Target{targetReq}}
	err := c.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Character: %s %s: unable to send target request to the server: %v",
			c.ID(), c.Serial(), err)
	}
}

// Use uses specified usable object.
func (c *Character) Use(ob useaction.Usable) {
	if c.Casted() != nil {
		return
	}
	err := c.Character.Use(ob)
	if err != nil {
		c.PrivateLog().Add(objects.Message{Text: "cant_do_right_now"})
		if !c.meetTargetRangeReqs(ob.UseAction().Requirements()...) {
			tar := c.Targets()[0]
			tarPosX, tarPosY := tar.Position()
			c.moveCloseTo(tarPosX, tarPosY, ob.UseAction().MinRange())
		}
		return
	}
	if c.game.Server() == nil {
		return
	}
	useReq := request.Use{
		UserID:     c.ID(),
		UserSerial: c.Serial(),
		ObjectID:   ob.ID(),
	}
	if ob, ok := ob.(serial.Serialer); ok {
		useReq.ObjectSerial = ob.Serial()
	}
	req := request.Request{Use: []request.Use{useReq}}
	err = c.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Character: %s %s: unable to send use request: %v",
			c.ID(), c.Serial(), err)
	}
}

// Equip inserts specified equipable item to all
// compatible slots in active PC equipment.
func (c *Character) Equip(it item.Equiper) error {
	if !c.MeetReqs(it.EquipReqs()...) {
		return fmt.Errorf(lang.Text("reqs_not_meet"))
	}
	slots := make([]*character.EquipmentSlot, 0)
	for _, itSlot := range it.Slots() {
		equiped := false
		for _, eqSlot := range c.Equipment().Slots() {
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
			c.Equipment().Unequip(it)
			return fmt.Errorf(lang.Text("equip_no_free_slot_error"))
		}
	}
	if !c.Equipment().Equiped(it) {
		return fmt.Errorf(lang.Text("equip_no_valid_slot_error"))
	}
	if c.game.Server() == nil {
		return nil
	}
	eqReq := request.Equip{
		CharID:     c.ID(),
		CharSerial: c.Serial(),
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
	err := c.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Character: %s %s: unable to send equip request: %v",
			c.ID(), c.Serial(), err)
	}
	return nil
}

// Unequip removes specified item from player equipment.
func (c *Character) Unequip(it item.Equiper) {
	c.Equipment().Unequip(it)
	if c.game.Server() == nil {
		return
	}
	uneqReq := request.Unequip{
		CharID:     c.ID(),
		CharSerial: c.Serial(),
		ItemID:     it.ID(),
		ItemSerial: it.Serial(),
	}
	req := request.Request{Unequip: []request.Unequip{uneqReq}}
	err := c.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Character: %s %s: unable to send unequip request: %v",
			c.ID(), c.Serial(), err)
	}
}

// meetTargetRangeReqs check if all target range requirements are meet.
// Returns true, if none of specified requirements is a target range
// requirement.
func (c *Character) meetTargetRangeReqs(reqs ...req.Requirement) bool {
	tarRangeReqs := make([]req.Requirement, 0)
	for _, r := range reqs {
		if r, ok := r.(*req.TargetRange); ok {
			tarRangeReqs = append(tarRangeReqs, r)
		}
	}
	if !c.MeetReqs(tarRangeReqs...) {
		return false
	}
	return true
}

// moveCloseTo moves player to the position at minimal range
// to the specified position.
func (c *Character) moveCloseTo(x, y, minRange float64) {
	switch {
	case x > c.Position().X:
		x -= minRange
	case x < c.Position().X:
		x += minRange
	case y > c.Position().Y:
		y -= minRange
	case y < c.Position().Y:
		y += minRange
	}
	c.SetDestPoint(x, y)
}
