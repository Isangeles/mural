/*
 * avatar.go
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

package objects

import (
	"github.com/faiface/pixel"

	"github.com/isangeles/flame/core/module/object/character"

	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/objects/internal"
)

// Avatar struct for graphical representation of
// game character.
type Avatar struct {
	*character.Character

	sprite         *internal.AvatarSprite
	portrait       *pixel.Sprite
	portraitName   string
	ssHeadName     string
	ssTorsoName    string
	ssFullBodyName string
}

// NewAvatar creates new avatar for specified game character.
// Portrait and spritesheet names are required for saving and
// loading avatar file.
func NewAvatar(char *character.Character, portraitPic,
	ssHeadPic, ssTorsoPic pixel.Picture, portraitName,
	ssHeadName, ssTorsoName string) (*Avatar, error) {
	av := new(Avatar)
	av.Character = char
	// Portrait & spritesheets names.
	av.portraitName = portraitName
	av.ssHeadName = ssHeadName
	av.ssTorsoName = ssTorsoName
	// Sprite.
	av.sprite = internal.NewAvatarSprite(ssTorsoPic, ssHeadPic)
	// Portrait.
	av.portrait = pixel.NewSprite(portraitPic, portraitPic.Bounds())
	return av, nil
}

// NewStaticAvatar creates new avatar with static(not affected by
// equipped items) body sprite.
// Portrait and spritesheet names are required for saving and
// loading avatar file.
func NewStaticAvatar(char *character.Character, portraitPic,
	ssBodyPic pixel.Picture, portraitName,
	ssBodyName string) (*Avatar, error) {
	av := new(Avatar)
	// Portrait & spritesheet names.
	av.portraitName = portraitName
	av.ssFullBodyName = ssBodyName
	// Sprite.
	av.sprite = internal.NewFullBodyAvatarSprite(ssBodyPic)
	// Portrait.
	av.portrait = pixel.NewSprite(portraitPic, portraitPic.Bounds())
	return av, nil
}

// Draw draws avatar.
func (av *Avatar) Draw(win *mtk.Window, matrix pixel.Matrix) {
	av.sprite.Draw(win, matrix)
}

// Update updates avatar.
func (av *Avatar) Update(win *mtk.Window) {
	if av.InMove() {
		av.sprite.Move()
		pos := av.Position()
		dest := av.DestPoint()
		switch {
		case pos.X < dest.X:
			av.sprite.Right()
		case pos.Y < dest.Y:
			av.sprite.Up()
		case pos.X > dest.X:
			av.sprite.Left()
		case pos.Y > dest.Y:
			av.sprite.Down()
		}
	} else {
		av.sprite.Idle()
	}
	av.sprite.Update(win)
}

// Portrait returns avatar portrait.
func (av *Avatar) Portrait() *pixel.Sprite {
	return av.portrait
}

// PortraitName returns name of portrait picture
// file.
func (av *Avatar) PortraitName() string {
	return av.portraitName
}

// TorsoSpritesheetName returns name of base torso
// spritesheet picture file.
func (av *Avatar) TorsoSpritesheetName() string {
	return av.ssTorsoName
}

// HeadSpritesheetName returns name of base head
// spritesheet picture file.
func (av *Avatar) HeadSpritesheetName() string {
	return av.ssHeadName
}

// FullBodySpritesheetName returns name of base
// full body spritesheet picture file.
func (av *Avatar) FullBodySpritesheetName() string {
	return av.ssFullBodyName
}

// Position return current position of avatar.
func (av *Avatar) Position() pixel.Vec {
	x, y := av.Character.Position()
	return pixel.V(x, y)
}

// DestPoint returns current destination point of
// avatar.
func (av *Avatar) DestPoint() pixel.Vec {
	x, y := av.Character.DestPoint()
	return pixel.V(x, y)
}

