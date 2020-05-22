/*
 * itemgraphic.go
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

package res

import (
	"encoding/xml"
)

// Struct for item graphics data.
type ItemGraphicsData struct {
	XMLName xml.Name          `xml:"item-graphics" json:"-"`
	Items   []ItemGraphicData `xml:"item-graphic" json:"item-graphics"`
}

// Struct for item graphic data.
type ItemGraphicData struct {
	ItemID       string             `xml:"id,attr" json:"id"`
	Icon         string             `xml:"icon,attr" json:"icon"`
	MaxStack     int                `xml:"stack,attr" json:"stack"`
	Spritesheets []*SpritesheetData `xml:"spritesheets>spritesheet" json:"spritesheets"`
}

// Struct for avatar spritesheet data.
type SpritesheetData struct {
	Texture string `xml:"texture,attr" json:"texture"`
	Race    string `xml:"race,attr" json:"race"`
	Gender  string `xml:"gender,attr" json:"gender"`
}
