/*
 * coords.go
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
	"math"
	
	"github.com/faiface/pixel"

	"github.com/isangeles/mural/config"
)

const (
	def_res_x, def_res_y float64 = 1920, 1080
)


// Scale return scale value for current resolution.
func Scale() float64 {
	res := config.Resolution()
	scaleX := res.X / def_res_x;
	scaleY := res.Y / def_res_y;
	scale := math.Round(math.Min(scaleX, scaleY) * 10) / 10;
	return scale
}

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

// Size converts specified default size value(for 1080p)
// to value for current resolution.
func ConvSize(size1080p float64) float64 {
	return size1080p * Scale()
}


