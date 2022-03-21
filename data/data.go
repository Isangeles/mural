/*
 * data.go
 *
 * Copyright 2018-2022 Dariusz Sikora <dev@isangeles.pl>
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

	"github.com/salviati/go-tmx/tmx"

	"github.com/faiface/pixel"

	flamedata "github.com/isangeles/flame/data"
	flameres "github.com/isangeles/flame/data/res"

	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/data/res/audio"
	"github.com/isangeles/mural/data/res/graphic"
)

const (
	SavesDir   = "saves"
	HUDFileExt = ".hud"
	ErrorIcon  = "unknown.png"
)

var guiPath string

// LoadModuleData loads graphic data from specified path.
// Should be called by GUI before creating any
// in-game elements.
func LoadModuleData(path string) (err error) {
	guiPath = path
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
	// Object spritesheets.
	graphic.ObjectSpritesheets, err = loadPicturesFromArch(graphicArchPath, "spritesheet/object")
	if err != nil {
		return fmt.Errorf("unable to load objects spritesheets: %v", err)
	}
	// Icons.
	graphic.Icons, err = loadPicturesFromArch(graphicArchPath, "icon")
	if err != nil {
		return fmt.Errorf("unable to load icons: %v", err)
	}
	// Avatars.
	avs, err := ImportAvatarsDir(filepath.Join(path, "avatars"))
	if err != nil {
		return fmt.Errorf("unable to import avatars: %v", err)
	}
	res.SetAvatars(avs)
	// Objects graphics.
	obGraphics, err := ImportObjectsGraphicsDir(filepath.Join(path, "objects"))
	if err != nil {
		return fmt.Errorf("unable to import objects graphics: %v", err)
	}
	res.SetObjects(obGraphics)
	// Items graphics.
	itGraphics, err := ImportItemsGraphicsDir(filepath.Join(path, "items"))
	if err != nil {
		return fmt.Errorf("unable to import items graphics: %v", err)
	}
	res.SetItems(itGraphics)
	// Effects graphic.
	effGraphics, err := ImportEffectsGraphicsDir(filepath.Join(path, "effects"))
	if err != nil {
		return fmt.Errorf("unable to import effects graphics: %v", err)
	}
	res.SetEffects(effGraphics)
	// Skills graphic.
	skillGraphics, err := ImportSkillsGraphicsDir(filepath.Join(path, "skills"))
	if err != nil {
		return fmt.Errorf("unable to import skills graphics: %v", err)
	}
	res.SetSkills(skillGraphics)
	// Translations.
	translations, err := flamedata.ImportLangDirs(filepath.Join(path, "lang"))
	if err != nil {
		return fmt.Errorf("unable to import translations: %v", err)
	}
	res.SetTranslationBases(translations)
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
	res.SetAvatars(append(res.Avatars(), avs...))
	// Object graphics.
	obGraphics, err := ImportObjectsGraphicsDir(filepath.Join(path, "objects"))
	if err != nil {
		return fmt.Errorf("unable to import objects graphics: %v", err)
	}
	res.SetObjects(append(res.Objects(), obGraphics...))
	return nil
}

// PlayablePortraits returns map with names of portraits as keys
// and portraits pictures as values avalible for player character.
func PlayablePortraits() (map[string]pixel.Picture, error) {
	files, err := os.ReadDir(filepath.Join(guiPath, "portraits"))
	if err != nil {
		return nil, fmt.Errorf("unable to read dir: %v", err)
	}
	portraits := make(map[string]pixel.Picture)
	for _, f := range files {
		if !f.IsDir() {
			img, err := loadPictureFromDir(filepath.FromSlash(
				guiPath + "/" + f.Name()))
			if err != nil {
				continue
			}
			portraits[f.Name()] = img
		}
	}
	return portraits, nil
}

// Map loads TMX map from file with specified
// directory and name.
func Map(areaDir string) (*tmx.Map, error) {
	tmxFile, err := os.Open(areaDir + "/map.tmx")
	if err != nil {
		return nil, fmt.Errorf("unable to open tmx file: %v", err)
	}
	tmxMap, err := tmx.Read(tmxFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read tmx file: %v", err)
	}
	return tmxMap, nil
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
