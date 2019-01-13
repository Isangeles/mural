/*
 * itemgraphic.go
 *
 * Copyright 2019 Dariusz Sikora <dev@isangeles.pl>
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

package object

import (
	"github.com/faiface/pixel"
	
	flameitem "github.com/isangeles/flame/core/module/object/item"
	
	"github.com/isangeles/mural/core/data/res"
)

// Struct for items with graphical
// representation.
type ItemGraphic struct {
	flameitem.Item
	spritesheet pixel.Picture
	icon        *pixel.Sprite
}

// NewItemGraphic creates new graphical wrapper from
// specified data.
func NewItemGraphic(data *res.ItemGraphicData) *ItemGraphic {
	item := new(ItemGraphic)
	item.Item = data.Item
	item.spritesheet = data.SpritesheetPic
	// TODO: icon sprite.
	return item
}

// Spritesheet returns item spritesheet.
func (it *ItemGraphic) Spritesheet() pixel.Picture {
	return it.spritesheet
}
