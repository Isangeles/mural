/*
 * core.go
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

package core

import (
	//"fmt"

	"github.com/faiface/pixel"
)

// Size type.
// Sizes: small(0), normal(1), big(2).
type Size int

const (
	SMALL Size = iota
	NORMAL
	BIG
)

// ButtonSize returns button rectange parameters
// for this size.
func (s Size) ButtonSize() pixel.Rect {
	switch(s) {
	case SMALL:
		return pixel.R(0, 0, ConvSize(70), ConvSize(35))
	default:
		return pixel.R(0, 0, ConvSize(70), ConvSize(35))
	}
}
