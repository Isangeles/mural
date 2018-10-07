/*
 * mtk.go
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

// Tool kid package for mural GUI.
package mtk

import (
	//"fmt"

	"golang.org/x/image/font"

	"github.com/faiface/pixel"

	"github.com/isangeles/mural/core/data"
)

// Type for sizes of UI elements, like buttons, switches, etc.
// Sizes: small(0), normal(1), big(2).
type Size int

// Type for shapes of UI elements.
// Shapes: rectangle(0), square(1).
type Shape int

const (
	SIZE_MINI Size = iota
	SIZE_SMALL 
	SIZE_MEDIUM
	SIZE_BIG

	SHAPE_RECTANGLE Shape = iota
	SHAPE_SQUARE 
)

// ButtonSize returns szie parameters for button with
// this size and with specifed shape.
func (s Size) ButtonSize(sh Shape) pixel.Rect {
	switch {
	case s <= SIZE_MINI && sh == SHAPE_SQUARE:
		return pixel.R(0, 0, ConvSize(30), ConvSize(30))
	case s == SIZE_SMALL && sh == SHAPE_SQUARE:
		return pixel.R(0, 0, ConvSize(45), ConvSize(45))
	case s == SIZE_MEDIUM && sh == SHAPE_SQUARE:
		return pixel.R(0, 0, ConvSize(100), ConvSize(100))
	case s <= SIZE_MINI && sh == SHAPE_RECTANGLE:
		return pixel.R(0, 0, ConvSize(30), ConvSize(15))
	case s == SIZE_SMALL && sh == SHAPE_RECTANGLE:
		return pixel.R(0, 0, ConvSize(70), ConvSize(35))
	case s == SIZE_MEDIUM && sh == SHAPE_RECTANGLE:
		return pixel.R(0, 0, ConvSize(100), ConvSize(50))
	default:
		return pixel.R(0, 0, ConvSize(70), ConvSize(35))
	}
}

// SwitchSize return rectangele parameters for switch
// with this size.
func (s Size) SwitchSize() pixel.Rect {
	switch {
	case s <= SIZE_SMALL:
		return pixel.R(0, 0, ConvSize(170), ConvSize(50))
	case s == SIZE_MEDIUM:
		return pixel.R(0, 0, ConvSize(200), ConvSize(70))
	default:
		return pixel.R(0, 0, ConvSize(70), ConvSize(35))
	}
}

// MessageWindowSize returns size parameters for message window. 
func (s Size) MessageWindowSize() pixel.Rect {
	switch {
	case s <= SIZE_SMALL:
		return pixel.R(0, 0, ConvSize(400), ConvSize(300))
	default:
		return pixel.R(0, 0, ConvSize(400), ConvSize(300))
	}
}

// Font returns main font in specified size from
// data package. 
func MainFont(s Size) font.Face {
	switch {
	case s == SIZE_SMALL:
		return data.MainFontSmall()
	case s == SIZE_MEDIUM:
		return data.MainFontNormal();
	case s >= SIZE_BIG:
		return data.MainFontBig();
	default:
		return data.MainFontSmall()
	}
}
