/*
 * object.go
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
	
	"github.com/isangeles/flame/core/module/object/area"

	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/mtk"
)

// Struct for graphical representation
// of area object.
type ObjectGraphic struct {
	*area.Object
	data     *res.ObjectGraphicData
	sprite   *mtk.Animation
	portrait *pixel.Sprite
	effects  map[string]*EffectGraphic
}

// NewObject creates new graphical wrapper for specified object.
func NewObject(ob *area.Object, data *res.ObjectGraphicData) *ObjectGraphic {
	og := new(ObjectGraphic)
	og.Object = ob
	og.data = data
	og.portrait = pixel.NewSprite(data.PortraitPic, data.PortraitPic.Bounds())
	og.sprite = mtk.NewAnimation(buildSpriteFrames(data.SpritePic), 2)
	return og
}

// buildSpriteFrames creates animation frames from specified
// spritesheet.
func buildSpriteFrames(ss pixel.Picture) []*pixel.Sprite {
	frames := []*pixel.Sprite{
		pixel.NewSprite(ss, ss.Bounds()),
	}
	return frames
}
