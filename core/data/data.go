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
// graphics/audio/text data.
package data

import (
	"fmt"
	"path/filepath"

	"golang.org/x/image/font"

	"github.com/golang/freetype/truetype"

	"github.com/faiface/pixel"

	"github.com/isangeles/flame"
)

var (
	g_dir_path  string
	g_arch_path string
	mainFont    *truetype.Font
)

// Called by GUI before creating any GUI elements.
func Load() error {
	if flame.Mod() == nil {
		return fmt.Errorf("no module loaded")
	}
	g_dir_path = filepath.FromSlash(fmt.Sprintf("data/modules/%s/gui",
		flame.Mod().Name()))
	g_arch_path = filepath.FromSlash(fmt.Sprintf("data/modules/%s/gui/gdata.zip",
		flame.Mod().Name()))
	var err error
	mainFont, err = Font("SIMSUN.ttf")
	if err != nil {
		return fmt.Errorf("fail to load main font")
		//mainFont = basicfont.Face7x13
	}
	return nil
}

// Sprite loads image with specified name from gdata
// archive.
func Picture(filePath string) (pixel.Picture, error) {
	return loadPictureFromArch(g_arch_path, filePath)
}

// Font loads font with specified name from gdata
// directory.
func Font(fileName string) (*truetype.Font, error) {
	fullpath := fmt.Sprintf("%s/%s/%s", g_dir_path, "font", fileName)
	return loadFontFromDir(filepath.FromSlash(fullpath))
}

// MainFont returns standard font in specified size.
func MainFont(size float64) font.Face {
	return truetype.NewFace(mainFont, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})
}
