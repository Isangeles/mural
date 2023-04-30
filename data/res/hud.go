/*
 * hud.go
 *
 * Copyright 2018-2023 Dariusz Sikora <ds@isangeles.dev>
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

// Struct for HUD data.
type HUDData struct {
	XMLName xml.Name `xml:"hud" json:"-"`
	Name    string   `xml:"name,attr" json:"name,attr"`
	Players []Player `xml:"players>player" json:"players"`
	Camera  Camera   `xml:"camera" json:"camera"`
}

// Struct for HUD camera data.
type Camera struct {
	X float64 `xml:"x,attr" json:"x"`
	Y float64 `xml:"y,attr" json:"y"`
}

// Struct for HUD player data (avatar, inventory layout, etc.).
type Player struct {
	ID       string `xml:"id" json:"id"`
	Serial   string `xml:"serial" json:"serial"`
	InvSlots []Slot `xml:"inventory>slot" json:"inv-slots"`
	BarSlots []Slot `xml:"bar>slot" json:"bar-slots"`
}

// Struct for HUD slot data.
type Slot struct {
	ID      int    `xml:"id,attr"`
	Content string `xml:"content,attr"`
}
