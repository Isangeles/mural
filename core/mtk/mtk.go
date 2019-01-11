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

// Tool kit package for mural GUI.
package mtk

import (
	//"fmt"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"

	"github.com/golang/freetype/truetype"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

const (
	SIZE_MINI Size = iota
	SIZE_SMALL
	SIZE_MEDIUM
	SIZE_BIG
	SIZE_HUGE

	SHAPE_RECTANGLE Shape = iota
	SHAPE_SQUARE
)

var (
	fallback_font  font.Face = basicfont.Face7x13
	main_font_base *truetype.Font
	sec_timer = time.Tick(time.Second)
)

// Type for shapes of UI elements.
// Shapes: rectangle(0), square(1).
type Shape int

// Type for sizes of UI elements, like buttons, switches, etc.
// Sizes: small(0), normal(1), big(2).
type Size int

// Interface for all 'focusable' UI elements, like buttons,
// switches, etc.
type Focuser interface {
	Focus(focus bool)
	Focused() bool
}

// Focus represents user focus on UI element.
type Focus struct {
	element Focuser
}

// Focus sets focus on specified focusable element.
// Previously focused element(if exists) is unfocused before
// focusing specified one.
func (f *Focus) Focus(e Focuser) {
	if f.element != nil {
		f.element.Focus(false)
	}
	f.element = e
	f.element.Focus(true)
}

// ButtonSize returns szie parameters for button with
// this size and with specifed shape.
func (s Size) ButtonSize(sh Shape) pixel.Rect {
	switch {
	case s <= SIZE_MINI && sh == SHAPE_SQUARE:
		return pixel.R(0, 0, ConvSize(30), ConvSize(30))
	case s == SIZE_SMALL && sh == SHAPE_SQUARE:
		return pixel.R(0, 0, ConvSize(45), ConvSize(45))
	case s == SIZE_MEDIUM && sh == SHAPE_SQUARE:
		return pixel.R(0, 0, ConvSize(60), ConvSize(60))
	case s >= SIZE_BIG && sh == SHAPE_SQUARE:
		return pixel.R(0, 0, ConvSize(70), ConvSize(70))
	case s <= SIZE_MINI && sh == SHAPE_RECTANGLE:
		return pixel.R(0, 0, ConvSize(30), ConvSize(15))
	case s == SIZE_SMALL && sh == SHAPE_RECTANGLE:
		return pixel.R(0, 0, ConvSize(70), ConvSize(35))
	case s == SIZE_MEDIUM && sh == SHAPE_RECTANGLE:
		return pixel.R(0, 0, ConvSize(100), ConvSize(50))
	case s >= SIZE_BIG && sh == SHAPE_RECTANGLE:
		return pixel.R(0, 0, ConvSize(120), ConvSize(70))
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
		return pixel.R(0, 0, ConvSize(210), ConvSize(80))
	case s == SIZE_BIG:
		return pixel.R(0, 0, ConvSize(260), ConvSize(140))
	case s >= SIZE_HUGE:
		return pixel.R(0, 0, ConvSize(270), ConvSize(130))
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

// ListSize returns size parameters for list.
func (s Size) ListSize() pixel.Rect {
	switch {
	case s <= SIZE_MEDIUM:
		return pixel.R(0, 0, ConvSize(600), ConvSize(500))
	default:
		return pixel.R(0, 0, ConvSize(600), ConvSize(500))
	}
}

// Sets specified truetype font as current main font of
// the interface.
func SetMainFont(font *truetype.Font) {
	main_font_base = font
}

// MainFont returns main font in specified size from
// data package.
func MainFont(s Size) font.Face {
	switch {
	case s <= SIZE_MINI:
		return createMainFont(10)
	case s <= SIZE_SMALL:
		return createMainFont(15)
	case s == SIZE_MEDIUM:
		return createMainFont(20)
	case s >= SIZE_BIG:
		return createMainFont(30)
	default:
		return createMainFont(10)
	}
}

// Atlas returns atlas for UI text with specified
// font.
func Atlas(f *font.Face) *text.Atlas {
	return text.NewAtlas(*f, text.ASCII)
}

// Matrix return scaled identity matrix.
func Matrix() pixel.Matrix {
	return pixel.IM.Scaled(pixel.V(0, 0), Scale())
}

// createMainFont creates new main font face with
// specified size.
func createMainFont(size float64) font.Face {
	if main_font_base == nil {
		return fallback_font
	}
	return truetype.NewFace(main_font_base, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})
}
