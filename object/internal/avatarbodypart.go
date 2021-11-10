/*
 * avatarbodypart.go
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

// Struct for avatar sprite body part
// (e.g. body, head, weapon).
// TODO: animations for kneel and lie.
type AvatarBodyPart struct {
	drawAnim      *mtk.MultiAnimation
	secAnim       *mtk.MultiAnimation
	idleAnim      *mtk.MultiAnimation
	moveAnim      *mtk.MultiAnimation
	meleeAnim     *mtk.MultiAnimation
	shootAnim     *mtk.MultiAnimation
	spellCastAnim *mtk.MultiAnimation
	craftCastAnim *mtk.MultiAnimation
	lieAnim       *mtk.MultiAnimation
}

const (
	animFPS = 2
)

// newAvatarBodyPart creates new body part from
// specified avatar spritesheet.
func newAvatarBodyPart(spritesheet pixel.Picture) *AvatarBodyPart {
	abp := new(AvatarBodyPart)
	abp.idleAnim = buildIdleAnim(spritesheet)
	abp.moveAnim = buildMoveAnim(spritesheet)
	abp.meleeAnim = buildMeleeAnim(spritesheet)
	abp.shootAnim = buildShootAnim(spritesheet)
	abp.spellCastAnim = buildCastAnim(spritesheet)
	abp.craftCastAnim = buildCastAnim(spritesheet)
	abp.lieAnim = buildLieAnim(spritesheet)
	abp.drawAnim = abp.idleAnim
	return abp
}

// Draw draws current body animation.
func (abp *AvatarBodyPart) Draw(t pixel.Target, matrix pixel.Matrix) {
	if abp.secAnim != nil {
		abp.secAnim.Draw(t, matrix)
		return
	}
	abp.drawAnim.Draw(t, matrix)
}

// Update updates current body animation.
func (abp *AvatarBodyPart) Update(win *mtk.Window) {
	if abp.secAnim != nil {
		abp.secAnim.Update(win)
		if abp.secAnim.Finished() {
			abp.secAnim.Loop(true)
			abp.secAnim = nil
		}
		return
	}
	abp.drawAnim.Update(win)
}

// Up turns current animataion up.
func (abp *AvatarBodyPart) Up() {
	abp.idleAnim.Up()
	abp.moveAnim.Up()
	abp.meleeAnim.Up()
	abp.shootAnim.Up()
	abp.spellCastAnim.Up()
	abp.craftCastAnim.Up()
}

// Right turns current animataion right.
func (abp *AvatarBodyPart) Right() {
	abp.idleAnim.Right()
	abp.moveAnim.Right()
	abp.meleeAnim.Right()
	abp.shootAnim.Right()
	abp.spellCastAnim.Right()
	abp.craftCastAnim.Right()
}

// Down turns current animataion down.
func (abp *AvatarBodyPart) Down() {
	abp.idleAnim.Down()
	abp.moveAnim.Down()
	abp.meleeAnim.Down()
	abp.shootAnim.Down()
	abp.spellCastAnim.Down()
	abp.craftCastAnim.Down()
}

// Left turns current animataion left.
func (abp *AvatarBodyPart) Left() {
	abp.idleAnim.Left()
	abp.moveAnim.Left()
	abp.meleeAnim.Left()
	abp.shootAnim.Left()
	abp.spellCastAnim.Left()
	abp.craftCastAnim.Left()
}

// Idle sets idle animation as current
// draw animation.
func (abp *AvatarBodyPart) Idle() {
	abp.drawAnim = abp.idleAnim
}

// Move sets move animation as current
// draw animation.
func (abp *AvatarBodyPart) Move() {
	abp.drawAnim = abp.moveAnim
}

// Melee sets melee animation as current
// draw animation.
func (abp *AvatarBodyPart) Melee() {
	abp.drawAnim = abp.meleeAnim
}

// Shoot sets shoot animation as current
// draw animation.
func (abp *AvatarBodyPart) Shoot() {
	abp.drawAnim = abp.shootAnim
}

// SpellCast sets spell cast animation as
// current draw animation.
func (abp *AvatarBodyPart) SpellCast() {
	abp.drawAnim = abp.spellCastAnim
}

// Lie sets lie animation as current draw
// animation.
func (abp *AvatarBodyPart) Lie() {
	abp.drawAnim = abp.lieAnim
}

// MeleeOnce restarts and sets melee animation
// as secondary animation to show only once.
func (abp *AvatarBodyPart) MeleeOnce() {
	abp.meleeAnim.Loop(false)
	abp.meleeAnim.Restart()
	abp.secAnim = abp.meleeAnim
}

// ShootOnce restarts and sets shoot animation
// as secondary animation to show only once.
func (abp *AvatarBodyPart) ShootOnce() {
	abp.shootAnim.Loop(false)
	abp.shootAnim.Restart()
	abp.secAnim = abp.shootAnim
}

// DrawArea returns current draw area.
func (abp *AvatarBodyPart) DrawArea() pixel.Rect {
	return abp.drawAnim.DrawArea()
}

// buildIdleAnim creates idle animations for each direction(up, right, down,
// left) with frames from specified spritesheet.
func buildIdleAnim(spritesheet pixel.Picture) *mtk.MultiAnimation {
	framesUp := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(0, 270, 80, 360)),
	}
	animUp := mtk.NewAnimation(framesUp, animFPS)
	framesRight := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(0, 180, 80, 270)),
	}
	animRight := mtk.NewAnimation(framesRight, animFPS)
	framesDown := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(0, 90, 80, 180)),
	}
	animDown := mtk.NewAnimation(framesDown, animFPS)
	framesLeft := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(0, 0, 80, 90)),
	}
	animLeft := mtk.NewAnimation(framesLeft, animFPS)
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
	animUp := mtk.NewAnimation(framesUp, animFPS)
	framesRight := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(80, 180, 160, 270)),
		pixel.NewSprite(spritesheet, pixel.R(160, 180, 240, 270)),
	}
	animRight := mtk.NewAnimation(framesRight, animFPS)
	framesDown := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(80, 90, 160, 180)),
		pixel.NewSprite(spritesheet, pixel.R(160, 90, 240, 180)),
	}
	animDown := mtk.NewAnimation(framesDown, animFPS)
	framesLeft := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(80, 0, 160, 90)),
		pixel.NewSprite(spritesheet, pixel.R(160, 0, 240, 90)),
	}
	animLeft := mtk.NewAnimation(framesLeft, animFPS)
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
	animUp := mtk.NewAnimation(framesUp, animFPS)
	framesRight := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(240, 180, 320, 270)),
		pixel.NewSprite(spritesheet, pixel.R(320, 180, 400, 270)),
	}
	animRight := mtk.NewAnimation(framesRight, animFPS)
	framesDown := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(240, 90, 320, 180)),
		pixel.NewSprite(spritesheet, pixel.R(320, 90, 400, 180)),
	}
	animDown := mtk.NewAnimation(framesDown, animFPS)
	framesLeft := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(240, 0, 320, 90)),
		pixel.NewSprite(spritesheet, pixel.R(320, 0, 400, 90)),
	}
	animLeft := mtk.NewAnimation(framesLeft, animFPS)
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
	animUp := mtk.NewAnimation(framesUp, animFPS)
	framesRight := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(400, 180, 480, 270)),
		pixel.NewSprite(spritesheet, pixel.R(480, 180, 560, 270)),
	}
	animRight := mtk.NewAnimation(framesRight, animFPS)
	framesDown := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(400, 90, 480, 180)),
		pixel.NewSprite(spritesheet, pixel.R(480, 90, 560, 180)),
	}
	animDown := mtk.NewAnimation(framesDown, animFPS)
	framesLeft := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(400, 0, 480, 90)),
		pixel.NewSprite(spritesheet, pixel.R(480, 0, 560, 90)),
	}
	animLeft := mtk.NewAnimation(framesLeft, animFPS)
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
	animUp := mtk.NewAnimation(framesUp, animFPS)
	framesRight := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(560, 180, 640, 270)),
		pixel.NewSprite(spritesheet, pixel.R(640, 180, 720, 270)),
	}
	animRight := mtk.NewAnimation(framesRight, animFPS)
	framesDown := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(560, 90, 640, 180)),
		pixel.NewSprite(spritesheet, pixel.R(640, 90, 720, 180)),
	}
	animDown := mtk.NewAnimation(framesDown, animFPS)
	framesLeft := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(560, 0, 640, 90)),
		pixel.NewSprite(spritesheet, pixel.R(640, 0, 720, 90)),
	}
	animLeft := mtk.NewAnimation(framesLeft, animFPS)
	anim := mtk.NewMultiAnimation(animUp, animRight, animDown, animLeft)
	return anim
}

// buildLieAnim creates lie animations for each directions with frames
// form specified spritesheet.
func buildLieAnim(spritesheet pixel.Picture) *mtk.MultiAnimation {
	framesUp := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(800, 270, 880, 360)),
	}
	animUp := mtk.NewAnimation(framesUp, animFPS)
	framesRight := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(800, 180, 880, 270)),
	}
	animRight := mtk.NewAnimation(framesRight, animFPS)
	framesDown := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(800, 90, 880, 180)),
	}
	animDown := mtk.NewAnimation(framesDown, animFPS)
	framesLeft := []*pixel.Sprite{
		pixel.NewSprite(spritesheet, pixel.R(800, 0, 880, 90)),
	}
	animLeft := mtk.NewAnimation(framesLeft, animFPS)
	anim := mtk.NewMultiAnimation(animUp, animRight, animDown, animLeft)
	return anim
}
