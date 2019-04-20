/*
 * coords.go
 *
 * Copyright 2018-2019 Dariusz Sikora <dev@isangeles.pl>
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

package mtk

import (
	"math"

	"github.com/faiface/pixel"
)

const (
	def_res_x, def_res_y float64 = 1920, 1080
)

var (
	scale float64 = 0.0
	res   pixel.Vec
)

// Scale return scale value for current resolution.
func Scale() float64 {
	/*
	res := config.Resolution()
	scaleX := res.X / def_res_x;
	scaleY := res.Y / def_res_y;
	s := math.Round(math.Min(scaleX, scaleY) * 10) / 10;
	return s
	*/
	return scale
}

// DisBR returns bottom right position of specified rectangle
// multiplied by specified value.
func DisBR(rect pixel.Rect, scale float64) pixel.Vec {
	return pixel.V(rect.Min.X-(rect.Max.X*scale),
		rect.Min.Y+(rect.Max.Y*scale))
}

// DisBL returns bottom left position of specified rectangle
// multiplied by specified value.
func DisBL(rect pixel.Rect, scale float64) pixel.Vec {
	return pixel.V(rect.Min.X+(rect.Max.X*scale),
		rect.Min.Y+(rect.Max.Y*scale))
}

// DisTR returns top right position of specified rectangle
// multiplied by specified value.
func DisTR(rect pixel.Rect, scale float64) pixel.Vec {
	return pixel.V(rect.Max.X-(rect.Max.X*scale),
		rect.Max.Y-(rect.Max.Y*scale))
}

// DisTL returns top left position of specified rectangle
// multiplied by specified value.
func DisTL(rect pixel.Rect, scale float64) pixel.Vec {
	return pixel.V(rect.Min.X+(rect.Max.X*scale),
		rect.Max.Y-(rect.Max.Y*scale))
}

// PosTR returns top right point for specified position
// of specified rectangle.
func PosTR(size, pos pixel.Vec) pixel.Vec {
	return pixel.V(pos.X-(size.X/2), pos.Y-(size.Y))
}

// PosTL returns top left point for specified position
// of specified rectangle.
func PosTL(size pixel.Rect, pos pixel.Vec) pixel.Vec {
	return pixel.V(pos.X+(size.Size().X/2), pos.Y-(size.Size().Y))
}

// PosBL returns bottom left point for specified position
// of specified rectangle.
func PosBL(size pixel.Rect, pos pixel.Vec) pixel.Vec {
	return pixel.V(pos.X+(size.Size().X/2), pos.Y+(size.Size().Y/2))
}

// PosBR returns bottom right point for specified position
// of specified rectangle.
func PosBR(size pixel.Rect, pos pixel.Vec) pixel.Vec {
	return pixel.V(pos.X-(size.Size().X/2), pos.Y+(size.Size().Y/2))
}

// DrawPosTR returns top left position for draw specifed
// object with specified size.
func DrawPosTL(bg pixel.Rect, size pixel.Vec) pixel.Vec {
	return pixel.V(bg.Min.X+(size.X)/2, bg.Max.Y-(size.Y)/2)
}

// DrawPosTR returns top right draw position(center) on specified
// background for object with specified size.
func DrawPosTR(bg pixel.Rect, size pixel.Vec) pixel.Vec {
	return pixel.V(bg.Max.X-(size.X)/2, bg.Max.Y-(size.Y)/2)
}

// DrawPosTC returns top center draw position(center) on specified
// background for object with specified size.
func DrawPosTC(bg pixel.Rect, size pixel.Vec) pixel.Vec {
	return pixel.V(bg.Center().X, bg.Max.Y-(size.Y)/2)
}

// DrawPosBL returns bottom left draw position(center) on specified
// background for object with specified size.
func DrawPosBL(bg pixel.Rect, size pixel.Vec) pixel.Vec {
	return pixel.V(bg.Min.X+(size.X/2), bg.Min.Y+(size.Y/2))
}

// DrawPosBR returns bottom right draw position(center) on specified
// background for object with specified size.
func DrawPosBR(bg pixel.Rect, size pixel.Vec) pixel.Vec {
	return pixel.V(bg.Max.X-(size.X/2), bg.Min.Y+(size.Y)/2)
}

// DrawPosBC returns bottom center draw position(center) on specified
// background for object with specified size.
func DrawPosBC(bg pixel.Rect, size pixel.Vec) pixel.Vec {
	return pixel.V(bg.Center().X, bg.Min.Y+(size.Y)/2)
}

// DrawPosCR returns center right draw position(center) on specified
// background for object with specified size.
func DrawPosCR(bg pixel.Rect, size pixel.Vec) pixel.Vec {
	return pixel.V(bg.Max.X-(size.X)/2, bg.Center().Y)
}

// DrawPosCL returns center left draw position(center) on specified
// background for object with specified size.
func DrawPosCL(bg pixel.Rect, size pixel.Vec) pixel.Vec {
	return pixel.V(bg.Min.X+(size.X)/2, bg.Center().Y)
}

// MoveTR returns move vector from center of background with
// specified size to top right point draw position(center) for
// specified size.
func MoveTR(bgSize, obSize pixel.Vec) pixel.Vec {
	return pixel.V(bgSize.X/2-obSize.X/2, bgSize.Y/2-obSize.Y/2)
}

// MoveTL returns move vector from center of background with specified
// size to top left point draw position(center) for specified size.
func MoveTL(bgSize, obSize pixel.Vec) pixel.Vec {
	return pixel.V(-bgSize.X/2-obSize.X/2, bgSize.Y/2-obSize.Y/2)
}

// MoveBR returns move vector from center of background with specified
// size to bottom right point draw position(center) for specified size.
func MoveBR(bgSize, obSize pixel.Vec) pixel.Vec {
	return pixel.V(bgSize.X/2-obSize.X/2, -bgSize.Y/2+obSize.Y/2)
}

// MoveBL returns move vector from center of background with specified
// size to bottom left point draw position(center) for specified size.
func MoveBL(bgSize, obSize pixel.Vec) pixel.Vec {
	return pixel.V(-bgSize.X/2+obSize.X/2, -bgSize.Y/2+obSize.Y/2)
}

// MoveBC returns move vector from center of background with specified
// size to bottom center point draw position(center) for specified size.
func MoveBC(bgSize, obSize pixel.Vec) pixel.Vec {
	return pixel.V(0, -bgSize.Y/2+obSize.Y/2)
}

// MoveTC returns move vector from center of background with specified
// size to top center point draw position(center) for specified size.
func MoveTC(bgSize, obSize pixel.Vec) pixel.Vec {
	return pixel.V(0, bgSize.Y/2-obSize.Y/2)
}

// TopOf returns position for rect with specified size at the top of
// specified draw area, with specified offset value.
func TopOf(drawArea pixel.Rect, size pixel.Vec, offset float64) pixel.Vec {
	return pixel.V(drawArea.Center().X, drawArea.Max.Y+
		(size.Y/2)+ConvSize(offset))
}

// ReightOf returns position for rect with specified size at the right side
// of specified draw area, with specified offset value.
func RightOf(drawArea pixel.Rect, size pixel.Vec, offset float64) pixel.Vec {
	return pixel.V(drawArea.Max.X+(size.X/2)+ConvSize(offset),
		drawArea.Min.Y+(size.Y/2))
}

// BottomOf returns position of rect with specified size at the bottom of
// speicified draw area, width specified offset value.
func BottomOf(drawArea pixel.Rect, size pixel.Vec, offset float64) pixel.Vec {
	return pixel.V(drawArea.Center().X, drawArea.Min.Y-
		(size.Y/2)-ConvSize(offset))
}

// LeftOf returns position for rect with specified size at the left side of
// specified draw area, with specified offset value.
func LeftOf(drawArea pixel.Rect, size pixel.Vec, offset float64) pixel.Vec {
	return pixel.V(drawArea.Min.X-(size.X/2)-ConvSize(offset),
		drawArea.Min.Y+(size.Y/2))
}

// Range returns range between two specified positions.
func Range(from, to pixel.Vec) float64 {
	return math.Hypot(to.X-from.X, to.Y-from.Y)
}

// Size converts specified default size value(for 1080p)
// to value for current resolution.
func ConvSize(size1080p float64) float64 {
	return size1080p * Scale()
}

// ConvVec converts specified default Pixel XY vector values(for 1080p)
// to vector with values for current resolution.
func ConvVec(vec1080p pixel.Vec) pixel.Vec {
	return pixel.V(ConvSize(vec1080p.X), ConvSize(vec1080p.Y))
}

// MatrixToDrawArea calculates draw area based on specified
// matrix and rectangle.
func MatrixToDrawArea(matrix pixel.Matrix, rectSize pixel.Vec) (drawArea pixel.Rect) {
	bgBottomX := matrix[4] - (rectSize.X / 2)
	bgBottomY := matrix[5] - (rectSize.Y / 2)
	drawArea.Min = pixel.V(bgBottomX, bgBottomY)
	drawArea.Max = drawArea.Min.Add(rectSize)
	return
}

// initScale calculates global scale for MTK elements for specified
// resolution.
// Called on new MTK window create.
func initScale(r pixel.Vec) {
	res = r
	scaleX := res.X / def_res_x
	scaleY := res.Y / def_res_y
	scale = math.Round(math.Min(scaleX, scaleY)*10) / 10
}
