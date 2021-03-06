/*
 * data.go
 *
 * Copyright 2018-2021 Dariusz Sikora <dev@isangeles.pl>
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

	"github.com/isangeles/flame"
	flamedata "github.com/isangeles/flame/data"
	flameres "github.com/isangeles/flame/data/res"

	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/data/res/audio"
	"github.com/isangeles/mural/data/res/graphic"
)

const (
	GUIModulePath   = "mural"
	SavesModulePath = GUIModulePath + "/saves"
	HUDFileExt      = ".hud"
	ErrorIcon       = "unknown.png"
)

var (
	// Paths.
	modAudioDirPath    string
	modGraphicDirPath  string
	modGraphicArchPath string
	modAudioArchPath   string
)

// LoadModuleData loads graphic data for specified module.
// Should be called by GUI before creating any
// in-game elements.
func LoadModuleData(mod *flame.Module) (err error) {
	// Load data resource paths.
	loadPaths(mod)
	// GUI textures.
	graphic.Textures, err = loadPicturesFromArch(modGraphicArchPath, "texture")
	if err != nil {
		return fmt.Errorf("unable to load textures: %v", err)
	}
	// Fonts.
	graphic.Fonts, err = loadFontsFromArch(modGraphicArchPath, "font")
	if err != nil {
		return fmt.Errorf("unable to load fonts: %v", err)
	}
	// Music.
	audio.Music, err = loadAudiosFromArch(modAudioArchPath, "music")
	if err != nil {
		return fmt.Errorf("unable to load music: %v", err)
	}
	// Audio effects.
	audio.Effects, err = loadAudiosFromArch(modAudioArchPath, "effect")
	if err != nil {
		return fmt.Errorf("unable to load audio effects: %v", err)
	}
	// Portraits.
	graphic.Portraits, err = loadPicturesFromArch(modGraphicArchPath, "portrait")
	if err != nil {
		return fmt.Errorf("unable to load portraits: %v", err)
	}
	// Avatar spritesheets.
	graphic.AvatarSpritesheets, err = loadPicturesFromArch(modGraphicArchPath, "spritesheet/avatar")
	if err != nil {
		return fmt.Errorf("unable to load avatars spritesheets: %v", err)
	}
	// Object spritesheets.
	graphic.ObjectSpritesheets, err = loadPicturesFromArch(modGraphicArchPath, "spritesheet/object")
	if err != nil {
		return fmt.Errorf("unable to load objects spritesheets: %v", err)
	}
	// Icons.
	graphic.Icons, err = loadPicturesFromArch(modGraphicArchPath, "icon")
	if err != nil {
		return fmt.Errorf("unable to load icons: %v", err)
	}
	// Avatars.
	path := filepath.Join(mod.Conf().Path, GUIModulePath, "avatars")
	avs, err := ImportAvatarsDir(path)
	if err != nil {
		return fmt.Errorf("unable to import avatars: %v", err)
	}
	res.SetAvatars(avs)
	// Objects graphics.
	path = filepath.Join(mod.Conf().Path, GUIModulePath, "objects")
	obGraphics, err := ImportObjectsGraphicsDir(path)
	if err != nil {
		return fmt.Errorf("unable to import objects graphics: %v", err)
	}
	res.SetObjects(obGraphics)
	// Items graphics.
	path = filepath.Join(mod.Conf().Path, GUIModulePath, "items")
	itGraphics, err := ImportItemsGraphicsDir(path)
	if err != nil {
		return fmt.Errorf("unable to import items graphics: %v", err)
	}
	res.SetItems(itGraphics)
	// Effects graphic.
	path = filepath.Join(mod.Conf().Path, GUIModulePath, "effects")
	effGraphics, err := ImportEffectsGraphicsDir(path)
	if err != nil {
		return fmt.Errorf("unable to import effects graphics: %v", err)
	}
	res.SetEffects(effGraphics)
	// Skills graphic.
	path = filepath.Join(mod.Conf().Path, GUIModulePath, "skills")
	skillGraphics, err := ImportSkillsGraphicsDir(path)
	if err != nil {
		return fmt.Errorf("unable to import skills graphics: %v", err)
	}
	res.SetSkills(skillGraphics)
	// Translations.
	path = filepath.Join(mod.Conf().Path, GUIModulePath, "lang")
	translations, err := flamedata.ImportLangDirs(path)
	if err != nil {
		return fmt.Errorf("unable to import translations: %v", err)
	}
	res.SetTranslationBases(translations)
	flameres.Add(flameres.ResourcesData{TranslationBases: translations})
	return nil
}

// LoadChapterData loads all graphical data for chapter.
func LoadChapterData(chapter *flame.Chapter) error {
	// Avatars.
	path := filepath.Join(chapter.Module().Conf().Path, GUIModulePath, "chapters",
		chapter.Conf().ID, "avatars")
	avs, err := ImportAvatarsDir(path)
	if err != nil {
		return fmt.Errorf("unable to import chapter avatars: %v", err)
	}
	res.SetAvatars(append(res.Avatars(), avs...))
	// Object graphics.
	path = filepath.Join(chapter.Module().Conf().Path, GUIModulePath, "chapters",
		chapter.Conf().ID, "objects")
	obGraphics, err := ImportObjectsGraphicsDir(path)
	if err != nil {
		return fmt.Errorf("unable to import objects graphics: %v", err)
	}
	res.SetObjects(append(res.Objects(), obGraphics...))
	return nil
}

// PlayablePortraits returns map with names of portraits as keys
// and portraits pictures as values avalible for player character.
func PlayablePortraits() (map[string]pixel.Picture, error) {
	path := filepath.Join(modGraphicDirPath, "portraits")
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir: %v", err)
	}
	portraits := make(map[string]pixel.Picture)
	for _, f := range files {
		if !f.IsDir() {
			img, err := loadPictureFromDir(filepath.FromSlash(
				path + "/" + f.Name()))
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

// Load loads grpahic directories.
func loadPaths(mod *flame.Module) {
	modGraphicDirPath = filepath.Join("data/modules", mod.Conf().ID, GUIModulePath)
	modAudioDirPath = filepath.Join("data/modules", mod.Conf().ID, GUIModulePath)
	modGraphicArchPath = filepath.Join("data/modules", mod.Conf().ID, GUIModulePath, "graphic.zip")
	modAudioArchPath = filepath.Join("data/modules", mod.Conf().ID, GUIModulePath, "audio.zip")
}
