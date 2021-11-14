/*
 * avatarsprite.go
 *
 * Copyright 2018-2021 Dariusz Sikora <dev@isangeles.pl>
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

package internal

import (
	"github.com/faiface/pixel"

	"github.com/isangeles/mtk"
)

// Struct for avatar sprite animations.
type AvatarSprite struct {
	head      *AvatarBodyPart
	torso     *AvatarBodyPart
	weapon    *AvatarBodyPart
	baseHead  *AvatarBodyPart
	baseTorso *AvatarBodyPart
	fullBody  *AvatarBodyPart
}

// NewAvatarSprite creates new sprite for specified
// body and head spritesheets.
func NewAvatarSprite(bodySpritesheet, headSpritesheet pixel.Picture) *AvatarSprite {
	as := new(AvatarSprite)
	as.baseHead = newAvatarBodyPart(headSpritesheet)
	as.baseTorso = newAvatarBodyPart(bodySpritesheet)
	as.head = as.baseHead
	as.torso = as.baseTorso
	return as
}

// NewFullBodyAvatarSprite creates new sprite with only one
// full body animated part.
func NewFullBodyAvatarSprite(spritesheet pixel.Picture) *AvatarSprite {
	as := new(AvatarSprite)
	as.fullBody = newAvatarBodyPart(spritesheet)
	return as
}

// Draw draws current sprite elements.
func (as *AvatarSprite) Draw(t pixel.Target, matrix pixel.Matrix) {
	if as.weapon != nil {
		as.weapon.Draw(t, matrix)
	}
	if as.fullBody != nil {
		as.fullBody.Draw(t, matrix)
		return
	}
	as.head.Draw(t, matrix)
	as.torso.Draw(t, matrix)
}

// Update updates current sprite elements.
func (as *AvatarSprite) Update(win *mtk.Window) {
	if as.weapon != nil {
		as.weapon.Update(win)
	}
	if as.fullBody != nil {
		as.fullBody.Update(win)
		return
	}
	as.head.Update(win)
	as.torso.Update(win)
}

// Up turns all current animataions up.
func (as *AvatarSprite) Up() {
	if as.weapon != nil {
		as.weapon.Up()
	}
	if as.fullBody != nil {
		as.fullBody.Up()
		return
	}
	as.head.Up()
	as.torso.Up()
}

// Right turns all current animataions right.
func (as *AvatarSprite) Right() {
	if as.weapon != nil {
		as.weapon.Right()
	}
	if as.fullBody != nil {
		as.fullBody.Right()
		return
	}
	as.head.Right()
	as.torso.Right()
}

// Down turns all current animataions down.
func (as *AvatarSprite) Down() {
	if as.weapon != nil {
		as.weapon.Down()
	}
	if as.fullBody != nil {
		as.fullBody.Down()
		return
	}
	as.head.Down()
	as.torso.Down()
}

// Left turns all current animataions left.
func (as *AvatarSprite) Left() {
	if as.weapon != nil {
		as.weapon.Left()
	}
	if as.fullBody != nil {
		as.fullBody.Left()
		return
	}
	as.head.Left()
	as.torso.Left()
}

// Idle sets idle animations as current
// draw animations.
func (as *AvatarSprite) Idle() {
	if as.weapon != nil {
		as.weapon.Idle()
	}
	if as.fullBody != nil {
		as.fullBody.Idle()
		return
	}
	as.head.Idle()
	as.torso.Idle()
}

// Move sets move animations as current
// draw animations.
func (as *AvatarSprite) Move() {
	if as.weapon != nil {
		as.weapon.Move()
	}
	if as.fullBody != nil {
		as.fullBody.Move()
		return
	}
	as.head.Move()
	as.torso.Move()
}

// Melee sets melee animations as current
// draw animations.
func (as *AvatarSprite) Melee() {
	if as.weapon != nil {
		as.weapon.Melee()
	}
	if as.fullBody != nil {
		as.fullBody.Melee()
		return
	}
	as.head.Melee()
	as.torso.Melee()
}

// Shoot sets shoot animations as current
// draw animations.
func (as *AvatarSprite) Shoot() {
	if as.weapon != nil {
		as.weapon.Shoot()
	}
	if as.fullBody != nil {
		as.fullBody.Shoot()
		return
	}
	as.head.Shoot()
	as.torso.Shoot()
}

// SpellCast sets spell cast animations as
// current draw animations.
func (as *AvatarSprite) SpellCast() {
	if as.weapon != nil {
		as.weapon.SpellCast()
	}
	if as.fullBody != nil {
		as.fullBody.SpellCast()
		return
	}
	as.head.SpellCast()
	as.torso.SpellCast()
}

// Kneel sets kneel animations as current draw
// animtaions.
func (as *AvatarSprite) Kneel() {
	if as.weapon != nil {
		as.weapon.Kneel()
	}
	if as.fullBody != nil {
		as.fullBody.Kneel()
		return
	}
	as.head.Kneel()
	as.torso.Kneel()
}

// Lie sets lie animations as current draw
// animtaions.
func (as *AvatarSprite) Lie() {
	if as.weapon != nil {
		as.weapon.Lie()
	}
	if as.fullBody != nil {
		as.fullBody.Lie()
		return
	}
	as.head.Lie()
	as.torso.Lie()
}

// CraftCast sets craft cast animations as
// current draw animations.
func (as *AvatarSprite) CraftCast() {
	if as.weapon != nil {
		as.weapon.SpellCast()
	}
	if as.fullBody != nil {
		as.fullBody.SpellCast()
		return
	}
	as.head.SpellCast()
	as.torso.SpellCast()
}

// MeleeOnce starts one melee animation
// for all sprite parts.
func (as *AvatarSprite) MeleeOnce() {
	if as.weapon != nil {
		as.weapon.MeleeOnce()
	}
	if as.fullBody != nil {
		as.fullBody.MeleeOnce()
		return
	}
	as.head.MeleeOnce()
	as.torso.MeleeOnce()
}

// ShootOnce starts one shoot animation
// for all sprite parts.
func (as *AvatarSprite) ShootOnce() {
	if as.weapon != nil {
		as.weapon.ShootOnce()
	}
	if as.fullBody != nil {
		as.fullBody.ShootOnce()
		return
	}
	as.head.ShootOnce()
	as.torso.ShootOnce()
}

// SetHead creates head animations from specified
// avatar spritesheet.
func (as *AvatarSprite) SetHead(spritesheet pixel.Picture) {
	as.head = newAvatarBodyPart(spritesheet)
}

// SetTorso creates body animations from specified
// avatar spritesheet.
func (as *AvatarSprite) SetTorso(spritesheet pixel.Picture) {
	if spritesheet == nil {
		as.torso = as.baseTorso
		return
	}
	as.torso = newAvatarBodyPart(spritesheet)
}

// SetFullBody creates new full body animations from specified
// avatar spritesheet.
func (as *AvatarSprite) SetFullBody(spritesheet pixel.Picture) {
	as.fullBody = newAvatarBodyPart(spritesheet)
}

// SetWeapon creates new weapon animations from specified
// avatar spritesheet.
func (as *AvatarSprite) SetWeapon(spritesheet pixel.Picture) {
	if spritesheet == nil {
		as.weapon = nil
		return
	}
	as.weapon = newAvatarBodyPart(spritesheet)
}

// Clear sets base body parts as current body parts.
func (as *AvatarSprite) Clear() {
	as.head = as.baseHead
	as.torso = as.baseTorso
}

// DrawArea returns current draw area.
func (as *AvatarSprite) DrawArea() pixel.Rect {
	if as.fullBody != nil {
		return as.fullBody.DrawArea()
	}
	return as.torso.DrawArea()
}
