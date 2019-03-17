/*
 * unmarshal.go
 *
 * Copyright 2019 Dariusz Sikora <dev@isangeles.pl>
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

package parsexml

import (
	"github.com/isangeles/mural/core/object"
)

// UnmarshalAvatarAnim parses specified string
// to avatar animation type.
func UnmarshalAvatarAnim(s string) object.AvatarAnimType {
	switch s {
	case "idle":
		return object.AvatarIdle
	case "move":
		return object.AvatarMove
	case "cast":
		return object.AvatarCast
	case "melee":
		return object.AvatarMelee
	case "shoot":
		return object.AvatarShoot
	case "kneel":
		return object.AvatarKneel
	case "lie":
		return object.AvatarLie
	default:
		return -1
	}
}
