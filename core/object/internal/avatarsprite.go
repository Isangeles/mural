/*
 * avatarsprite.go
 *
 * Copyright 2018-2019 Dariusz Sikora <dev@isangeles.pl>
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
	
	"github.com/isangeles/mural/core/mtk"
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
func NewAvatarSprite(bodySpritesheet,
	headSpritesheet pixel.Picture) *AvatarSprite {
	spr := new(AvatarSprite)
	spr.baseHead = newAvatarBodyPart(headSpritesheet)
	spr.baseTorso = newAvatarBodyPart(bodySpritesheet)
	spr.head = spr.baseHead
	spr.torso = spr.baseTorso
	return spr
}

// NewFullBodyAvatarSprite creates new sprite with only one
// full body animated part.
func NewFullBodyAvatarSprite(spritesheet pixel.Picture) *AvatarSprite {
	spr := new(AvatarSprite)
	spr.fullBody = newAvatarBodyPart(spritesheet)
	return spr
}

// Draw draws current sprite elements.
func (spr *AvatarSprite) Draw(t pixel.Target, matrix pixel.Matrix) {
	if spr.weapon != nil {
		spr.weapon.Draw(t, matrix)
	}
	if spr.fullBody != nil {
		spr.fullBody.Draw(t, matrix)
		return
	}
	spr.head.Draw(t, matrix)
	spr.torso.Draw(t, matrix)
}

// Update updates current sprite elements.
func (spr *AvatarSprite) Update(win *mtk.Window) {
	if spr.weapon != nil {
		spr.weapon.Update(win)
	}
	if spr.fullBody != nil {
		spr.fullBody.Update(win)
		return
	}
	spr.head.Update(win)
	spr.torso.Update(win)
}

// Up turns all current animataions up.
func (spr *AvatarSprite) Up() {
	if spr.weapon != nil {
		spr.weapon.Up()
	}
	if spr.fullBody != nil {
		spr.fullBody.Up()
		return
	}
	spr.head.Up()
	spr.torso.Up()
}

// Right turns all current animataions right.
func (spr *AvatarSprite) Right() {
	if spr.weapon != nil {
		spr.weapon.Right()
	}
	if spr.fullBody != nil {
		spr.fullBody.Right()
		return
	}
	spr.head.Right()
	spr.torso.Right()
}

// Down turns all current animataions down.
func (spr *AvatarSprite) Down() {
	if spr.weapon != nil {
		spr.weapon.Down()
	}
	if spr.fullBody != nil {
		spr.fullBody.Down()
		return
	}
	spr.head.Down()
	spr.torso.Down()
}

// Left turns all current animataions left.
func (spr *AvatarSprite) Left() {
	if spr.weapon != nil {
		spr.weapon.Left()
	}
	if spr.fullBody != nil {
		spr.fullBody.Left()
		return
	}
	spr.head.Left()
	spr.torso.Left()
}

// Idle sets idle animations as current
// draw animations.
func (spr *AvatarSprite) Idle() {
	if spr.weapon != nil {
		spr.weapon.Idle()
	}
	if spr.fullBody != nil {
		spr.fullBody.Idle()
		return
	}
	spr.head.Idle()
	spr.torso.Idle()
}

// Move sets move animations as current
// draw animations.
func (spr *AvatarSprite) Move() {
	if spr.weapon != nil {
		spr.weapon.Move()
	}
	if spr.fullBody != nil {
		spr.fullBody.Move()
		return
	}
	spr.head.Move()
	spr.torso.Move()
}

// Melee sets melee animations as current
// draw animations.
func (spr *AvatarSprite) Melee() {
	if spr.weapon != nil {
		spr.weapon.Melee()
	}
	if spr.fullBody != nil {
		spr.fullBody.Melee()
		return
	}
	spr.head.Melee()
	spr.torso.Melee()
}
// Shoot sets shoot animations as current
// draw animations.
func (spr *AvatarSprite) Shoot() {
	if spr.weapon != nil {
		spr.weapon.Shoot()
	}
	if spr.fullBody != nil {
		spr.fullBody.Shoot()
		return
	}
	spr.head.Shoot()
	spr.torso.Shoot()
}

// Cast sets cast animations as current
// draw animations.
func (spr *AvatarSprite) Cast() {
	if spr.weapon != nil {
		spr.weapon.Cast()
	}
	if spr.fullBody != nil {
		spr.fullBody.Cast()
		return
	}
	spr.head.Cast()
	spr.torso.Cast()
}

// MeleeOnce starts one melee animation
// for all sprite parts.
func (spr *AvatarSprite) MeleeOnce() {
	if spr.weapon != nil {
		spr.weapon.MeleeOnce()
	}
	if spr.fullBody != nil {
		spr.fullBody.MeleeOnce()
		return
	}
	spr.head.MeleeOnce()
	spr.torso.MeleeOnce()
}

// ShootOnce starts one shoot animation
// for all sprite parts.
func (spr *AvatarSprite) ShootOnce() {
	if spr.weapon != nil {
		spr.weapon.ShootOnce()
	}
	if spr.fullBody != nil {
		spr.fullBody.ShootOnce()
		return
	}
	spr.head.ShootOnce()
	spr.torso.ShootOnce()
}

// SetHead creates head animations from specified
// avatar spritesheet.
func (spr *AvatarSprite) SetHead(spritesheet pixel.Picture) {
	spr.head = newAvatarBodyPart(spritesheet)
}

// SetTorso creates body animations from specified
// avatar spritesheet.
func (spr *AvatarSprite) SetTorso(spritesheet pixel.Picture) {
	spr.torso = newAvatarBodyPart(spritesheet)
}

// SetFullBody creates new full body animations from specified
// avatar spritesheet.
func (spr *AvatarSprite) SetFullBody(spritesheet pixel.Picture) {
	spr.fullBody = newAvatarBodyPart(spritesheet)
}

// SetWeapon creates new weapon animations from specified
// avatar spritesheet.
func (spr *AvatarSprite) SetWeapon(spritesheet pixel.Picture) {
	if spritesheet == nil {
		spr.weapon = nil
		return
	}
	spr.weapon = newAvatarBodyPart(spritesheet)
}

// Clear sets base body parts as current body parts.
func (spr *AvatarSprite) Clear() {
	spr.head = spr.baseHead
	spr.torso = spr.baseTorso
}

// DrawArea returns current draw area.
func (spr *AvatarSprite) DrawArea() pixel.Rect {
	if spr.fullBody != nil {
		return spr.fullBody.DrawArea()
	}
	return spr.torso.DrawArea()
}

