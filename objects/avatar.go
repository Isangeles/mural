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
)

// Avatar struct for graphical representation of
// game character.
type Avatar struct {
	*character.Character

	sprite          *mtk.Animation
	portrait        *pixel.Sprite
	portraitName    string
	spritesheetName string
}

// NewAvatar creates new avatar for specified game character.
// Portrait and spritesheet names are required for saving and
// loading avatar file.
func NewAvatar(char *character.Character, portraitPic,
	spritesheetPic pixel.Picture, portraitName,
	spritesheetName string) (*Avatar, error) {
	av := new(Avatar)
	av.Character = char
	av.portraitName = portraitName
	av.spritesheetName = spritesheetName
	// Sprite animation.
	spriteFrames := []*pixel.Sprite{
		pixel.NewSprite(spritesheetPic, pixel.R(0, 0, 32, 32)),
		pixel.NewSprite(spritesheetPic, pixel.R(32, 0, 64, 32)),
	}
	av.sprite = mtk.NewAnimation(spriteFrames, 1)
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

// SpritesheetName returns name of spritesheet picture
// file.
func (av *Avatar) SpritesheetName() string {
	return av.spritesheetName
}

// Position return current position of avatar.
func (av *Avatar) Position() pixel.Vec {
	x, y := av.Character.Position()
	return pixel.V(x, y)
}

