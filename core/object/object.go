/*
 * object.go
 *
 * Copyright 2019-2021 Dariusz Sikora <dev@isangeles.pl>
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
	"github.com/isangeles/flame/effect"
	"github.com/isangeles/flame/item"

	"github.com/faiface/pixel"

	"github.com/isangeles/mural/core/data/res"
)

// Interface for all objects with
// drawable objects.
type Drawer interface {
	Position() pixel.Vec
	DrawArea() pixel.Rect
}

// DefaultEffectGraphic returns default graphic for specified effect.
func DefaultEffectGraphic(eff *effect.Effect) *res.EffectGraphicData {
	return &res.EffectGraphicData{
		EffectID: eff.ID(),
		Icon:     defaultEffectIcon,
	}
}

// DefaultItemGraphic returns default graphic for specified item.
func DefaultItemGraphic(it item.Item) *res.ItemGraphicData {
	return &res.ItemGraphicData{
		ItemID:   it.ID(),
		Icon:     defaultItemIcon,
		MaxStack: 100,
	}
}
