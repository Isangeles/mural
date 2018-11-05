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
	"fmt"
	"github.com/faiface/pixel"

	"github.com/isangeles/flame/core/module/object/character"

	"github.com/isangeles/mural/core/data"
)

// Avatar struct for graphical representation of
// game character.
type Avatar struct {
	*character.Character
	
	portrait    *pixel.Sprite
	portraitName string
}

// NewAvatar creates new avatar for specified game character.
func NewAvatar(char *character.Character, portraitName string) (*Avatar, error) {
	av := new(Avatar)
	av.Character = char
	av.portraitName = portraitName
	portraitPic, err := data.Portrait(av.portraitName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_portrait")
	}
	av.portrait = pixel.NewSprite(portraitPic, portraitPic.Bounds())
	return av, nil
}

// Portrait returns avatar portrait.
func (av *Avatar) Portrait() *pixel.Sprite {
	return av.portrait
}

