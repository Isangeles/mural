/*
 * mainmenu.go
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
	"github.com/faiface/pixel"
)

// DisBL returns bottom left position of specified rect
// multiplied by specified value.
func DisBL(rect pixel.Rect, scale float64) pixel.Vec {
	return pixel.V(rect.Min.X + (rect.Max.X * scale),
		rect.Min.Y + (rect.Max.Y * scale))
}

// DisTR returns top right position of specified rect
// multiplied by specified value.
func DisTR(rect pixel.Rect, scale float64) pixel.Vec {
	return pixel.V(rect.Max.X - (rect.Max.X * scale),
		rect.Max.Y - (rect.Max.Y * scale))
}

// TopLeftDis3 returns top left point of specified window
// divided by 3.
func BottomLeftDis3(rect pixel.Rect) pixel.Vec {
	return pixel.V(rect.Min.X + (rect.Max.X / 3), rect.Min.Y + (rect.Max.Y / 3))
}

// BottomRightDis3 returns bottom right point of specified window,
// divided by 3.
func TopRightDis3(rect pixel.Rect) pixel.Vec {
	return pixel.V(rect.Max.X - (rect.Max.X / 3), rect.Max.Y - (rect.Max.Y / 3))
}

// PosBL return bottom left point for specified position
// in specified rectangle.
func PosBL(size pixel.Rect, pos pixel.Vec) pixel.Vec {
	return pixel.V(pos.X + (size.Size().X / 2), pos.Y + (size.Size().Y / 2))
}

// PosBR returns bottom right point for specified position
// in specified rectangle.
func PosBR(size pixel.Rect, pos pixel.Vec) pixel.Vec {
	return pixel.V(pos.X - (size.Size().X / 2), pos.Y + (size.Size().Y / 2))
}


