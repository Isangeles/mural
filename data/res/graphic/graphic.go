/*
 * graphic.go
 *
 * Copyright 2020-2024 Dariusz Sikora <ds@isangeles.dev>
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

// Package with pre-loaded game graphics.
package graphic

import (
	"github.com/golang/freetype/truetype"

	"github.com/gopxl/pixel"
)

var (
	Textures           map[string]pixel.Picture
	AvatarSpritesheets map[string]pixel.Picture
	Icons              map[string]pixel.Picture
	Portraits          map[string]pixel.Picture
	Fonts              map[string]*truetype.Font
)

// On init.
func init() {
	Textures = make(map[string]pixel.Picture)
	AvatarSpritesheets = make(map[string]pixel.Picture)
	Icons = make(map[string]pixel.Picture)
	Portraits = make(map[string]pixel.Picture)
	Fonts = make(map[string]*truetype.Font)
}
