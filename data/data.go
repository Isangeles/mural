/*
 * data.go
 *
 * Copyright 2018-2025 Dariusz Sikora <ds@isangeles.dev>
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
// graphic and audio data.
package data

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gopxl/pixel"

	flamedata "github.com/isangeles/flame/data"
	flameres "github.com/isangeles/flame/data/res"

	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/data/res/audio"
	"github.com/isangeles/mural/data/res/graphic"
	"github.com/isangeles/mural/log"
)

const (
	HUDDir     = "hud"
	HUDFileExt = ".json"
	ErrorIcon  = "unknown.png"
)

// LoadModuleData loads graphic data from specified path.
// Should be called by GUI before creating any
// in-game elements.
func LoadModuleData(path string) (err error) {
	// GUI textures.
	graphicArchPath := filepath.Join(path, "graphic.zip")
	graphic.Textures, err = loadPicturesFromArch(graphicArchPath, "texture")
	if err != nil {
		return fmt.Errorf("unable to load textures: %v", err)
	}
	// Fonts.
	graphic.Fonts, err = loadFontsFromArch(graphicArchPath, "font")
	if err != nil {
		return fmt.Errorf("unable to load fonts: %v", err)
	}
	audioArchPath := filepath.Join(path, "audio.zip")
	// Music.
	audio.Music, err = loadAudiosFromArch(audioArchPath, "music")
	if err != nil {
		return fmt.Errorf("unable to load music: %v", err)
	}
	// Audio effects.
	audio.Effects, err = loadAudiosFromArch(audioArchPath, "effect")
	if err != nil {
		return fmt.Errorf("unable to load audio effects: %v", err)
	}
	// Portraits.
	graphic.Portraits, err = loadPicturesFromArch(graphicArchPath, "portrait")
	if err != nil {
		return fmt.Errorf("unable to load portraits: %v", err)
	}
	// Avatar spritesheets.
	graphic.AvatarSpritesheets, err = loadPicturesFromArch(graphicArchPath, "spritesheet/avatar")
	if err != nil {
		return fmt.Errorf("unable to load avatars spritesheets: %v", err)
	}
	// Icons.
	graphic.Icons, err = loadPicturesFromArch(graphicArchPath, "icon")
	if err != nil {
		return fmt.Errorf("unable to load icons: %v", err)
	}
	// Avatars.
	res.Avatars, err = ImportAvatarsDir(filepath.Join(path, "avatars"))
	if err != nil {
		return fmt.Errorf("unable to import avatars: %v", err)
	}
	// Items graphics.
	res.Items, err = ImportItemsGraphicsDir(filepath.Join(path, "items"))
	if err != nil {
		return fmt.Errorf("unable to import items graphics: %v", err)
	}
	// Effects graphic.
	res.Effects, err = ImportEffectsGraphicsDir(filepath.Join(path, "effects"))
	if err != nil {
		return fmt.Errorf("unable to import effects graphics: %v", err)
	}
	// Skills graphic.
	res.Skills, err = ImportSkillsGraphicsDir(filepath.Join(path, "skills"))
	if err != nil {
		return fmt.Errorf("unable to import skills graphics: %v", err)
	}
	// Translations.
	translations, err := flamedata.ImportLangDirs(filepath.Join(path, "lang"))
	if err != nil {
		return fmt.Errorf("unable to import translations: %v", err)
	}
	res.AddTranslationBases(translations)
	flameres.Add(flameres.ResourcesData{TranslationBases: translations})
	return nil
}

// LoadChapterData loads all chapter graphical data from specified path.
func LoadChapterData(path string) error {
	// Avatars.
	avs, err := ImportAvatarsDir(filepath.Join(path, "avatars"))
	if err != nil {
		return fmt.Errorf("unable to import chapter avatars: %v", err)
	}
	res.Avatars = append(res.Avatars, avs...)
	return nil
}

// Pictures loads all pictures from specified path and returns
// map with file names as keys and pictures as values.
func Pictures(path string) (map[string]pixel.Picture, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir: %v", err)
	}
	portraits := make(map[string]pixel.Picture)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		img, err := loadPictureFromDir(filepath.Join(path, file.Name()))
		if err != nil {
			log.Err.Printf("Data pictures: unable to load picture: %v", err)
			continue
		}
		portraits[file.Name()] = img
	}
	return portraits, nil
}

// DirFiles returns names of all files matching specified
// file name pattern in directory with specified path.
func DirFiles(path, pattern string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir: %v", err)
	}
	names := make([]string, 0)
	for _, file := range files {
		match, err := regexp.MatchString(pattern, file.Name())
		if err != nil {
			return nil, fmt.Errorf("unable to execute pattern: %v", err)
		}
		if match {
			names = append(names, file.Name())
		}
	}
	return names, nil
}
