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
	effects  map[string]*EffectGraphic
}

// NewObject creates new graphical wrapper for specified object.
func NewObjectGraphic(ob *area.Object, data *res.ObjectGraphicData) *ObjectGraphic {
	og := new(ObjectGraphic)
	og.Object = ob
	og.data = data
	og.sprite = mtk.NewAnimation(buildSpriteFrames(data.SpritePic), 2)
	return og
}

// Draw draws object sprite.
func (og *ObjectGraphic) Draw(t pixel.Target, matrix pixel.Matrix) {
	og.sprite.Draw(t, matrix)
}

// Update updates object.
func (og *ObjectGraphic) Update(win *mtk.Window) {
	og.sprite.Update(win)
}

// DrawArea returns current draw area of
// object sprite.
func (og *ObjectGraphic) DrawArea() pixel.Rect {
	return og.sprite.DrawArea()
}

// Portrait returns portrait picture.
func (og *ObjectGraphic) Portrait() pixel.Picture {
	return og.data.PortraitPic
}

// Position return object position in form of
// pixel XY vector.
func (og *ObjectGraphic) Position() pixel.Vec {
	x, y := og.Object.Position()
	return pixel.V(x, y)
}

// Effects returns all object effects in form of
// graphical wrappers.
func (og *ObjectGraphic) Effects() []*EffectGraphic {
	effs := make([]*EffectGraphic, 0)
	for _, e := range og.effects {
		effs = append(effs, e)
	}
	return effs
}

// Mana returns 0, objects does not have mana.
// Function to sadisfy HUD frame object interface.
func (og *ObjectGraphic) Mana() int {
	return 0
}

// MaxMana returns 0, objects does not have mana.
// Function to sadisfy HUD frame object interface.
func (og *ObjectGraphic) MaxMana() int {
	return 0
}

// buildSpriteFrames creates animation frames from specified
// spritesheet.
func buildSpriteFrames(ss pixel.Picture) []*pixel.Sprite {
	frames := []*pixel.Sprite{
		pixel.NewSprite(ss, ss.Bounds()),
	}
	return frames
}
