/*
 * layout.go
 *
 * Copyright 2019-2020 Dariusz Sikora <dev@isangeles.pl>
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

package hud

import (
	flameobject "github.com/isangeles/flame/module/objects"

	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/core/object"
)

// Struct for HUD layout.
// Stores layout of items and skills
// in menu bar and inventory menu.
type Layout struct {
	invSlots map[string]int
	barSlots map[string]int
}

// NewLayout creates new HUD layout.
func NewLayout() *Layout {
	l := new(Layout)
	l.invSlots = make(map[string]int)
	l.barSlots = make(map[string]int)
	return l
}

// SetInvSlots sets specified layout map as
// current inventory slots content layout.
func (l *Layout) SetInvSlots(m map[string]int) {
	l.invSlots = m
}

// SetBarSlot sets specified layout map as
// current menu bar slots content layout.
func (l *Layout) SetBarSlots(m map[string]int) {
	l.barSlots = m
}

// InvSlots returns map with saved inventory
// slots.
func (l *Layout) InvSlots() map[string]int {
	return l.invSlots
}

// BarSlots returns map witg saved bar slots.
func (l *Layout) BarSlots() map[string]int {
	return l.barSlots
}

// SaveInvSlot saves position(slot ID) of specified item at
// inventory slot list.
func (l *Layout) SaveInvSlot(ob *object.ItemGraphic, slotID int) {
	l.invSlots[ob.ID()+ob.Serial()] = slotID
}

// SaveBarSlot saves position(slot ID) of specified object
// at menu bar.
func (l *Layout) SaveBarSlot(ob interface{}, id int) {
	switch ob := ob.(type) {
	case *object.ItemGraphic:
		l.barSlots[ob.ID()+ob.Serial()] = id
	case *object.SkillGraphic:
		l.barSlots[ob.ID()] = id
	default:
		log.Err.Printf("hud: layout: save bar slot: unsupported object type: %v",
			ob)
	}
}

// InvSlotID returns saved inventory slot ID for specified
// object.
func (l *Layout) InvSlotID(ob flameobject.Object) int {
	id, ok := l.invSlots[ob.ID()+ob.Serial()]
	if !ok {
		return -1
	}
	return id
}

// BarSlotID returns saved menu bar slot ID for specified
// object.
func (l *Layout) BarSlotID(ob interface{}) int {
	id := -1
	switch ob := ob.(type) {
	case *object.ItemGraphic:
		i, ok := l.barSlots[ob.ID()+ob.Serial()]
		if ok {
			id = i
		}
	case *object.SkillGraphic:
		i, ok := l.barSlots[ob.ID()]
		if ok {
			id = i
		}
	}
	return id
}
