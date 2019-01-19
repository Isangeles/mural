/*
 * avatarbodypart.go
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

// Struct for avatar sprite body part
// (e.g. body, head, weapon).
// TODO: animations for kneel, melee,
// range shoot and spell cast.
type AvatarBodyPart struct {
	drawAnim *mtk.MultiAnimation
	idleAnim *mtk.MultiAnimation
	moveAnim *mtk.MultiAnimation
}

// newAvatarBodyPart creates new body part from
// specified avatar sdpritesheet.
func newAvatarBodyPart(spritesheet pixel.Picture) *AvatarBodyPart {
	part := new(AvatarBodyPart)
	part.idleAnim = buildIdleAnim(spritesheet)
	part.moveAnim = buildMoveAnim(spritesheet)
	part.drawAnim = part.idleAnim
	return part
}

// Draw draws current body animation.
func (body *AvatarBodyPart) Draw(t pixel.Target, matrix pixel.Matrix) {
	body.drawAnim.Draw(t, matrix)
}

// Update updates current body animation.
func (body *AvatarBodyPart) Update(win *mtk.Window) {
	body.drawAnim.Update(win)
}

// Up turns current animataion up.
func (body *AvatarBodyPart) Up() {
	body.idleAnim.Up()
	body.moveAnim.Up()	
}

// Right turns current animataion right.
func (part *AvatarBodyPart) Right() {
	part.idleAnim.Right()
	part.moveAnim.Right()
}

// Down turns current animataion down.
func (part *AvatarBodyPart) Down() {
	part.idleAnim.Down()
	part.moveAnim.Down()
}

// Left turns current animataion left.
func (part *AvatarBodyPart) Left() {
	part.idleAnim.Left()
	part.moveAnim.Left()
}

// Idle sets idle animation as current
// draw animation.
func (part *AvatarBodyPart) Idle() {
	part.drawAnim = part.idleAnim
}

// Move sets move animation as current
// draw animation.
func (part *AvatarBodyPart) Move() {
	part.drawAnim = part.moveAnim
}

// buildIdleAnim creates idle animations for each direction(up, right, down,
// left) with frames from specified spritesheet.
func buildIdleAnim(spritesheet pixel.Picture) *mtk.MultiAnimation {
	framesUp := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(0, 270, 80, 360)),
	}
	animUp := mtk.NewAnimation(framesUp, 2)
	framesRight := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(0, 180, 80, 270)),
	}
	animRight := mtk.NewAnimation(framesRight, 2)
	framesDown := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(0, 90, 80, 180)),
	}
	animDown := mtk.NewAnimation(framesDown, 2)
	framesLeft := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(0, 0, 80, 90)),
	}
	animLeft := mtk.NewAnimation(framesLeft, 2)
	anim := mtk.NewMultiAnimation(animUp, animRight, animDown, animLeft)
	return anim
}

// buildMoveAnim creates move animations for each direction(up, right, down,
// left) with frames from specified spritesheet.
func buildMoveAnim(spritesheet pixel.Picture) *mtk.MultiAnimation {
	framesUp := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(80, 270, 160, 360)),
		pixel.NewSprite(spritesheet, pixel.R(160, 270, 240, 360)),
	}
	animUp := mtk.NewAnimation(framesUp, 2)
	framesRight := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(80, 180, 160, 270)),
		pixel.NewSprite(spritesheet, pixel.R(160, 180, 240, 270)),
	}
	animRight := mtk.NewAnimation(framesRight, 2)
	framesDown := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(80, 90, 160, 180)),
		pixel.NewSprite(spritesheet, pixel.R(160, 90, 240, 180)),
	}
	animDown := mtk.NewAnimation(framesDown, 2)
	framesLeft := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(80, 0, 160, 90)),
		pixel.NewSprite(spritesheet, pixel.R(160, 0, 240, 90)),
	}
	animLeft := mtk.NewAnimation(framesLeft, 2)
	anim := mtk.NewMultiAnimation(animUp, animRight, animDown, animLeft)
	return anim
}

