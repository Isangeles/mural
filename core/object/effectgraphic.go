/*
 * effectgraphic.go
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

	"github.com/isangeles/flame/core/module/object/effect"

	"github.com/isangeles/mural/core/data/res"
)

// Struct for graphical wrapper
// for effects.
type EffectGraphic struct {
	*effect.Effect
	icon *pixel.Sprite
}

// NewEffectGraphic creates new graphical wrapper for specified effect.
func NewEffectGraphic(effect *effect.Effect, data *res.EffectGraphicData) *EffectGraphic {
	eg := new(EffectGraphic)
	eg.Effect = effect
	eg.icon = pixel.NewSprite(data.IconPic, data.IconPic.Bounds())
	return eg
}

// Icon returns effect icon.
func (eg *EffectGraphic) Icon() *pixel.Sprite {
	return eg.icon
}
