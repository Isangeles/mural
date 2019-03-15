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
// TODO: animations for kneel and lie.
type AvatarBodyPart struct {
	drawAnim  *mtk.MultiAnimation
	secAnim   *mtk.MultiAnimation
	idleAnim  *mtk.MultiAnimation
	moveAnim  *mtk.MultiAnimation
	meleeAnim *mtk.MultiAnimation
	shootAnim *mtk.MultiAnimation
	castAnim  *mtk.MultiAnimation
}

const (
	anim_fps = 2
)

// newAvatarBodyPart creates new body part from
// specified avatar sdpritesheet.
func newAvatarBodyPart(spritesheet pixel.Picture) *AvatarBodyPart {
	part := new(AvatarBodyPart)
	part.idleAnim = buildIdleAnim(spritesheet)
	part.moveAnim = buildMoveAnim(spritesheet)
	part.meleeAnim = buildMeleeAnim(spritesheet)
	part.shootAnim = buildShootAnim(spritesheet)
	part.castAnim = buildCastAnim(spritesheet)
	part.drawAnim = part.idleAnim
	return part
}

// Draw draws current body animation.
func (body *AvatarBodyPart) Draw(t pixel.Target, matrix pixel.Matrix) {
	if body.secAnim != nil {
		body.secAnim.Draw(t, matrix)
		return
	}
	body.drawAnim.Draw(t, matrix)
}

// Update updates current body animation.
func (body *AvatarBodyPart) Update(win *mtk.Window) {
	if body.secAnim != nil {
		body.secAnim.Update(win)
		if body.secAnim.Finished() {
			body.secAnim.Loop(true)
			body.secAnim = nil
		}
		return
	}
	body.drawAnim.Update(win)
}

// Up turns current animataion up.
func (part *AvatarBodyPart) Up() {
	part.idleAnim.Up()
	part.moveAnim.Up()
	part.meleeAnim.Up()
	part.shootAnim.Up()
	part.castAnim.Up()
}

// Right turns current animataion right.
func (part *AvatarBodyPart) Right() {
	part.idleAnim.Right()
	part.moveAnim.Right()
	part.meleeAnim.Right()
	part.shootAnim.Right()
	part.castAnim.Right()
}

// Down turns current animataion down.
func (part *AvatarBodyPart) Down() {
	part.idleAnim.Down()
	part.moveAnim.Down()
	part.meleeAnim.Down()
	part.shootAnim.Down()
	part.castAnim.Down()
}

// Left turns current animataion left.
func (part *AvatarBodyPart) Left() {
	part.idleAnim.Left()
	part.moveAnim.Left()
	part.meleeAnim.Left()
	part.shootAnim.Left()
	part.castAnim.Left()
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

// Melee sets melee animation as current
// draw animation.
func (part *AvatarBodyPart) Melee() {
	part.drawAnim = part.meleeAnim
}

// Shoot sets shoot animation as current
// draw animation.
func (part *AvatarBodyPart) Shoot() {
	part.drawAnim = part.shootAnim
}

// Cast sets cast animation as current
// draw animation.
func (part *AvatarBodyPart) Cast() {
	part.drawAnim = part.castAnim
}

// MeleeOnce restarts and sets melee animation
// as secondary animation to show only once.
func (part *AvatarBodyPart) MeleeOnce() {
	part.meleeAnim.Loop(false)
	part.meleeAnim.Restart()
	part.secAnim = part.meleeAnim
}

// ShootOnce restarts and sets shoot animation
// as secondary animation to show only once.
func (part *AvatarBodyPart) ShootOnce() {
	part.shootAnim.Loop(false)
	part.shootAnim.Restart()
	part.secAnim = part.shootAnim
}

// DrawArea returns current draw area.
func (part *AvatarBodyPart) DrawArea() pixel.Rect {
	return part.drawAnim.DrawArea()
}

// buildIdleAnim creates idle animations for each direction(up, right, down,
// left) with frames from specified spritesheet.
func buildIdleAnim(spritesheet pixel.Picture) *mtk.MultiAnimation {
	framesUp := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(0, 270, 80, 360)),
	}
	animUp := mtk.NewAnimation(framesUp, anim_fps)
	framesRight := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(0, 180, 80, 270)),
	}
	animRight := mtk.NewAnimation(framesRight, anim_fps)
	framesDown := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(0, 90, 80, 180)),
	}
	animDown := mtk.NewAnimation(framesDown, anim_fps)
	framesLeft := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(0, 0, 80, 90)),
	}
	animLeft := mtk.NewAnimation(framesLeft, anim_fps)
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
	animUp := mtk.NewAnimation(framesUp, anim_fps)
	framesRight := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(80, 180, 160, 270)),
		pixel.NewSprite(spritesheet, pixel.R(160, 180, 240, 270)),
	}
	animRight := mtk.NewAnimation(framesRight, anim_fps)
	framesDown := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(80, 90, 160, 180)),
		pixel.NewSprite(spritesheet, pixel.R(160, 90, 240, 180)),
	}
	animDown := mtk.NewAnimation(framesDown, anim_fps)
	framesLeft := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(80, 0, 160, 90)),
		pixel.NewSprite(spritesheet, pixel.R(160, 0, 240, 90)),
	}
	animLeft := mtk.NewAnimation(framesLeft, anim_fps)
	anim := mtk.NewMultiAnimation(animUp, animRight, animDown, animLeft)
	return anim
}

// buildMeleeAnim creates melee animations for each direction(up, right, down,
// left) with frames from specified spritesheet.
func buildMeleeAnim(spritesheet pixel.Picture) *mtk.MultiAnimation {
	framesUp := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(240, 270, 320, 360)),
		pixel.NewSprite(spritesheet, pixel.R(320, 270, 400, 360)),
	}
	animUp := mtk.NewAnimation(framesUp, anim_fps)
	framesRight := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(240, 180, 320, 270)),
		pixel.NewSprite(spritesheet, pixel.R(320, 180, 400, 270)),
	}
	animRight := mtk.NewAnimation(framesRight, anim_fps)
	framesDown := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(240, 90, 320, 180)),
		pixel.NewSprite(spritesheet, pixel.R(320, 90, 400, 180)),
	}
	animDown := mtk.NewAnimation(framesDown, anim_fps)
	framesLeft := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(240, 0, 320, 90)),
		pixel.NewSprite(spritesheet, pixel.R(320, 0, 400, 90)),
	}
	animLeft := mtk.NewAnimation(framesLeft, anim_fps)
	anim := mtk.NewMultiAnimation(animUp, animRight, animDown, animLeft)
	return anim
}

// buildShootAnim creates shoot animations for each direction(up, right, down,
// left) with frames from specified spritesheet.
func buildShootAnim(spritesheet pixel.Picture) *mtk.MultiAnimation {
	framesUp := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(400, 270, 480, 360)),
		pixel.NewSprite(spritesheet, pixel.R(480, 270, 560, 360)),
	}
	animUp := mtk.NewAnimation(framesUp, anim_fps)
	framesRight := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(400, 180, 480, 270)),
		pixel.NewSprite(spritesheet, pixel.R(480, 180, 560, 270)),
	}
	animRight := mtk.NewAnimation(framesRight, anim_fps)
	framesDown := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(400, 90, 480, 180)),
		pixel.NewSprite(spritesheet, pixel.R(480, 90, 560, 180)),
	}
	animDown := mtk.NewAnimation(framesDown, anim_fps)
	framesLeft := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(400, 0, 480, 90)),
		pixel.NewSprite(spritesheet, pixel.R(480, 0, 560, 90)),
	}
	animLeft := mtk.NewAnimation(framesLeft, anim_fps)
	anim := mtk.NewMultiAnimation(animUp, animRight, animDown, animLeft)
	return anim
}

// builCastAnim creates cast animations for each direction(up, right, down,
// left) with frames from specified spritesheet.
func buildCastAnim(spritesheet pixel.Picture) *mtk.MultiAnimation {
	framesUp := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(560, 270, 640, 360)),
		pixel.NewSprite(spritesheet, pixel.R(640, 270, 720, 360)),
	}
	animUp := mtk.NewAnimation(framesUp, anim_fps)
	framesRight := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(560, 180, 640, 270)),
		pixel.NewSprite(spritesheet, pixel.R(640, 180, 720, 270)),
	}
	animRight := mtk.NewAnimation(framesRight, anim_fps)
	framesDown := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(560, 90, 640, 180)),
		pixel.NewSprite(spritesheet, pixel.R(640, 90, 720, 180)),
	}
	animDown := mtk.NewAnimation(framesDown, anim_fps)
	framesLeft := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(560, 0, 640, 90)),
		pixel.NewSprite(spritesheet, pixel.R(640, 0, 720, 90)),
	}
	animLeft := mtk.NewAnimation(framesLeft, anim_fps)
	anim := mtk.NewMultiAnimation(animUp, animRight, animDown, animLeft)
	return anim
}

