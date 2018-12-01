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

	sprite      *pixel.Sprite
	portrait    *pixel.Sprite
	portraitName string
}

// NewAvatar creates new avatar for specified game character.
func NewAvatar(char *character.Character, spritePic,
	portraitPic pixel.Picture) (*Avatar) {
	av := new(Avatar)
	av.Character = char
	// Sprite.
	// TODO: handling spritesheets(frames, animations, etc.).
	av.sprite = pixel.NewSprite(spritePic, spritePic.Bounds())
	// Portrait.
	av.portrait = pixel.NewSprite(portraitPic, portraitPic.Bounds())
	return av
}

// Draw draws avatar.
func (av *Avatar) Draw(win *mtk.Window, matrix pixel.Matrix) {
	av.sprite.Draw(win, matrix)
}

// Portrait returns avatar portrait.
func (av *Avatar) Portrait() *pixel.Sprite {
	return av.portrait
}

// Position return current position of avatar.
func (av *Avatar) Position() pixel.Vec {
	x, y := av.Character.Position()
	return pixel.V(x, y)
}

