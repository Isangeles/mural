/*
 * guisave.go
 *
 * Copyright 2018-2020 Dariusz Sikora <dev@isangeles.pl>
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

package res

import (
	"encoding/xml"
)

// Struct for GUI state save.
type GUISave struct {
	XMLName xml.Name     `xml:"save" json:"-"`
	Name    string       `xml:"name,attr" json:"name,attr"`
	Players []PlayerSave `xml:"players>player" json:"players"`
	Camera  CameraSave   `xml:"camera" json:"camera"`
}

// Struct for saved camera data.
type CameraSave struct {
	X float64 `xml:"x,attr" json:"x"`
	Y float64 `xml:"y,attr" json:"y"`
}

// Struct for saved GUI user data
// (avatar, inventory layout, etc.).
type PlayerSave struct {
	Avatar   AvatarData `xml:"avatar" json:"avatar"`
	InvSlots []SlotSave `xml:"inventory>slot" json:"inv-slots"`
	BarSlots []SlotSave `xml:"bar>slot" json:"bar-slots"`
}

// Struct for saved slot.
type SlotSave struct {
	ID      int    `xml:"id,attr"`
	Content string `xml:"content,attr"`
}
