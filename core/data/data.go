/*
 * data.go
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

// data package contains functions for loading
// graphics/audi/text data
package data

import (
	"fmt"
	"path/filepath"

	"github.com/faiface/pixel"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

var (
	g_dir_path     = filepath.FromSlash("data/gdata")
	g_arch_path    = filepath.FromSlash("data/gdata.zip")
	mainFontSmall  font.Face
	mainFontNormal font.Face
	mainFontBig    font.Face
)

func init() {
	var err error
	mainFontSmall, err = Font("SIMSUN.ttd", 10)
	if err != nil {
		// TODO: log somewhere information about main font load failure.
		mainFontSmall = basicfont.Face7x13
	}
	mainFontNormal, err = Font("SIMSUN.ttf", 20)
	if err != nil {
		// TODO: log somewhere information about main font load failure.
		mainFontNormal = basicfont.Face7x13
	}
	mainFontBig, err = Font("SIMSUN.ttf", 40)
	if err != nil {
		// TODO: log somewhere information about main font load failure.
		mainFontBig = basicfont.Face7x13
	}
}

// Sprite loads image with specified name from gdata
// archive.
func Picture(filePath string) (pixel.Picture, error) {
	return loadPictureFromArch(g_arch_path, filePath)
}

// Font loads font with specified name from gdata
// directory.
func Font(fileName string, size float64) (font.Face, error) {
	fullpath := fmt.Sprintf("%s/%s/%s", g_dir_path, "font", fileName)
	return loadFontFromDir(filepath.FromSlash(fullpath), size)
}

// MainFontSmall returns standard font in small
// size.
func MainFontSmall() font.Face {
	return mainFontSmall
}

// MainFontNormal returns standard font in normal
// size.
func MainFontNormal() font.Face {
	return mainFontNormal
}

// MainFontBig returns standard in big size.
func MainFontBig() font.Face {
	return mainFontBig
}
