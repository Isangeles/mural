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
)

// Struct for HUD layout.
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
func (l *Layout) SaveInvSlot(ob flameobject.Object, slotID int) {
	l.invSlots[ob.ID()+ob.Serial()] = slotID
}

// SaveBarSlot saves position(slot ID) of specified object
// at menu bar.
func (l *Layout) SaveBarSlot(ob flameobject.Object, slotID int) {
	l.barSlots[ob.ID()+ob.Serial()] = slotID
}

// InvSlotID returns saved inventory slot ID for specified
// object.
func (l *Layout) InvSlotID(ob flameobject.Object) int {
	id, prs := l.invSlots[ob.ID()+ob.Serial()]
	if !prs {
		return -1
	}
	return id
}

// BarSlotID returns saved menu bar slot ID for specified
// object.
func (l *Layout) BarSlotID(ob flameobject.Object) int {
	id, prs := l.barSlots[ob.ID()+ob.Serial()]
	if !prs {
		return -1
	}
	return id
}
