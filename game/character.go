/*
 * character.go
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
	"github.com/isangeles/flame/training"
	"github.com/isangeles/flame/useaction"

	"github.com/isangeles/fire/request"

	"github.com/isangeles/ignite/ai"

	"github.com/isangeles/mural/log"
)

// Struct for game character.
type Character struct {
	*character.Character
	game       *Game
	name       string
	combatLog  *objects.Log
	privateLog *objects.Log
	onUse      func(o useaction.Usable)
}

// NewCharacter creates game wrapper for module character.
func NewCharacter(char *character.Character, game *Game) *Character {
	c := Character{
		Character:  char,
		game:       game,
		name:       lang.Text(char.ID()),
		combatLog:  objects.NewLog(),
		privateLog: objects.NewLog(),
	}
	return &c
}

// Name returns character name.
func (c *Character) Name() string {
	return c.name
}

// CombatLog returns character combat log.
func (c *Character) CombatLog() *objects.Log {
	return c.combatLog
}

// PrivateLog returns character private log.
func (c *Character) PrivateLog() *objects.Log {
	return c.privateLog
}

// InSight checks if specified XY position is in sight range
// of the character.
func (c *Character) InSight(x, y float64) bool {
	charX, charY := c.Position()
	return math.Hypot(charX-x, charY-y) <= c.SightRange()
}

// SetOnUseFunc sets function to trigger after using an object.
func (c *Character) SetOnUseFunc(f func(o useaction.Usable)) {
	c.onUse = f
	aiChar := c.aiChar()
	if aiChar != nil {
		aiChar.AddOnUseEvent(f)
	}
}

// SetPosition sets character position and destination point.
func (c *Character) SetPosition(x, y float64) {
	c.Character.SetPosition(x, y)
	c.SetDestPoint(x, y)
	if c.game.Server() == nil {
		return
	}
	setPosReq := request.SetPos{c.ID(), c.Serial(), x, y}
	req := request.Request{SetPos: []request.SetPos{setPosReq}}
	err := c.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Character: %s %s: unable to send set pos request to the server: %v",
			c.ID(), c.Serial(), err)
	}
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
	if c.game.Server() == nil {
		c.ChatLog().Add(message)
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

// Train uses specified training from specified trainer.
func (c *Character) Train(training *training.TrainerTraining, trainer training.Trainer) {
	err := c.Character.Use(training)
	if err != nil {
		c.PrivateLog().Add(objects.NewMessage("cant_do_right_now", false))
	}
	if c.game.Server() == nil {
		// If no server then trigger onUse event and return.
		// With server this event will be triggered after
		// user response from the server.
		if c.onUse != nil {
			c.onUse(training)
		}
		return
	}
	trainReq := request.Training{
		TrainingID:    training.ID(),
		TrainerID:     trainer.ID(),
		TrainerSerial: trainer.Serial(),
		UserID:        c.ID(),
		UserSerial:    c.Serial(),
	}
	req := request.Request{Training: []request.Training{trainReq}}
	err = c.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Character: %s %s: unable to send training request: %v",
			c.ID(), c.Serial(), err)
	}
}

// Use uses specified usable object.
func (c *Character) Use(ob useaction.Usable) {
	if c.Casted() != nil || ob.UseAction() == nil {
		return
	}
	err := c.Character.Use(ob)
	if err != nil {
		c.PrivateLog().Add(objects.NewMessage("cant_do_right_now", false))
		if !c.meetTargetRangeReqs(ob.UseAction().Requirements()...) {
			tar := c.Targets()[0]
			tarPosX, tarPosY := tar.Position()
			c.moveCloseTo(tarPosX, tarPosY, ob.UseAction().MinRange())
		}
		return
	}
	if c.game.Server() == nil {
		// If no server then trigger onUse event and return.
		// With server this event will be triggered after
		// user response from the server.
		if c.onUse != nil {
			c.onUse(ob)
		}
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

// RemoveItem removes specified item from player inventory
// and creates loot area object.
func (c *Character) RemoveItems(items ...item.Item) {
	reqItems := make(map[string][]string, 0)
	for _, it := range items {
		c.Inventory().RemoveItem(it)
		reqItems[it.ID()] = append(reqItems[it.ID()], it.Serial())
	}
	// Server request to remove items.
	if c.game.Server() == nil {
		return
	}
	throwItemsReq := request.ThrowItems{
		ObjectID:     c.ID(),
		ObjectSerial: c.Serial(),
		Items:        reqItems,
	}
	req := request.Request{ThrowItems: []request.ThrowItems{throwItemsReq}}
	err := c.game.Server().Send(req)
	if err != nil {
		log.Err.Printf("Character: %s %s: unable to send throw items request: %v",
			c.ID(), c.Serial(), err)
	}
}

// Usable returns usable object with specified ID and serial.
func (c *Character) Usable(id, serial string) useaction.Usable {
	for _, s := range c.Skills() {
		if s.ID() == id {
			return s
		}
	}
	for _, r := range c.Crafting().Recipes() {
		if r.ID() == id {
			return r
		}
	}
	return c.Inventory().Item(id, serial)
}

// Targeted checks if character is targeted by the active
// player character.
func (c *Character) Targeted() bool {
	if c.game.ActivePlayerChar() == nil {
		return false
	}
	for _, tar := range c.game.ActivePlayerChar().Targets() {
		if tar.ID() == c.ID() && tar.Serial() == c.Serial() {
			return true
		}
	}
	return false
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
	charX, charY := c.Position()
	switch {
	case x > charX:
		x -= minRange
	case x < charX:
		x += minRange
	case y > charY:
		y -= minRange
	case y < charY:
		y += minRange
	}
	c.SetDestPoint(x, y)
}

// aiChar returns AI wrapper for the character, or nil
// if AI doesn't control the character.
func (c *Character) aiChar() *ai.Character {
	for _, aiChar := range c.game.localAI.Game().Characters() {
		if c.ID() == aiChar.ID() && c.Serial() == aiChar.Serial() {
			return aiChar
		}
	}
	return nil
}
