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

package object

import (
	"github.com/faiface/pixel"
	
	"github.com/isangeles/flame/module/item"
	
	"github.com/isangeles/mural/core/data/res"
)

// Struct for graphical wrapper
// for items.
type ItemGraphic struct {
	item.Item
	data *res.ItemGraphicData
}

// NewItemGraphic creates new graphical wrapper for specified item.
func NewItemGraphic(item item.Item, data *res.ItemGraphicData) *ItemGraphic {
	itg := new(ItemGraphic)
	itg.Item = item
	itg.data = data
	return itg
}

// Spritesheets returns all spritesheets data.
func (itg *ItemGraphic) Spritesheets() []*res.SpritesheetData {
	return itg.data.Spritesheets
}

// Icon returns item icon.
func (itg *ItemGraphic) Icon() pixel.Picture {
	return itg.data.IconPic
}

// MaxStack returns maximal number of stacked items
// with same ID.
func (itg *ItemGraphic) MaxStack() int {
	return itg.data.MaxStack
}
