/*
 * avatarsprite.go
 *
 * Copyright 2018 Dariusz Sikora <dev@isangeles.pl>
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
}

// NewAvatarSprite creates new sprite for specified
// body and head spritesheets.
func NewAvatarSprite(bodySpritesheet,
	headSpritesheet pixel.Picture) *AvatarSprite {
	spr := new(AvatarSprite)
	spr.baseHead = newAvatarBodyPart(bodySpritesheet)
	spr.baseTorso = newAvatarBodyPart(bodySpritesheet)
	spr.head = spr.baseHead
	spr.torso = spr.baseTorso
	return spr
}

// Draw draws current sprite elements.
func (spr *AvatarSprite) Draw(t pixel.Target, matrix pixel.Matrix) {
	spr.torso.Draw(t, matrix)
}

// Update updates current sprite elements.
func (spr *AvatarSprite) Update(win *mtk.Window) {
	spr.torso.Update(win)
}

// Up turns all current animataions up.
func (spr *AvatarSprite) Up() {
	spr.torso.Up()
}

// Right turns all current animataions right.
func (spr *AvatarSprite) Right() {
	spr.torso.Right()
}

// Down turns all current animataions down.
func (spr *AvatarSprite) Down() {
	spr.torso.Down()
}

// Left turns all current animataions left.
func (spr *AvatarSprite) Left() {
	spr.torso.Left()
}

// Idle sets idle animations as current
// draw animations.
func (spr *AvatarSprite) Idle() {
	spr.torso.Idle()
}

// Move sets move animations as current
// draw animations.
func (spr *AvatarSprite) Move() {
	spr.torso.Move()
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

// Clear sets base body parts as current body parts.
func (spr *AvatarSprite) Clear() {
	spr.head = spr.baseHead
	spr.torso = spr.baseTorso
}

