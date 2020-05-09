/*
 * data.go
 *
 * Copyright 2018-2020 Dariusz Sikora <dev@isangeles.pl>
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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/salviati/go-tmx/tmx"

	"github.com/faiface/beep"

	"github.com/faiface/pixel"

	"github.com/isangeles/flame/module"
	
	"github.com/isangeles/burn/ash"

	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/data/res/graphic"
	"github.com/isangeles/mural/log"
)

const (
	SavesModulePath = "gui/saves"
	SaveFileExt     = ".savegui"
)

var (
	// Paths.
	modAudioDirPath     string
	modGraphicDirPath   string
	modGraphicArchPath  string
	modAudioArchPath    string
	// Scritps.
	ashScriptExt = ".ash"
)

// LoadModuleData loads graphic data for specified module.
// Should be called by GUI before creating any
// in-game elements.
func LoadModuleData(mod *module.Module) (err error) {
	// Load data resource paths.
	loadPaths(mod)
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
	return nil
}

// LoadUIData loads UI graphic data for specified module.
// Should be called by GUI before creating any
// GUI elements.
func LoadUIData(mod *module.Module) (err error) {
	// Load data sources paths.
	loadPaths(mod)
	// GUI elements textures.
	graphic.Textures, err = loadPicturesFromArch(modGraphicArchPath, "texture")
	if err != nil {
		return fmt.Errorf("unable to load textures: %v", err)
	}
	// Fonts.
	graphic.Fonts, err = loadFontsFromArch(modGraphicArchPath, "font")
	if err != nil {
		return fmt.Errorf("unable to load fonts: %v", err)
	}
	return nil
}

// PlayablePortraits returns map with names of portraits as keys
// and portraits pictures as values avalible for player character.
func PlayablePortraits() (map[string]pixel.Picture, error) {
	path := filepath.Join(modGraphicDirPath, "portraits")
	files, err := ioutil.ReadDir(path)
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

// Music returns audio stream data from file with specified name
// inside audio archive.
func Music(fileName string) (*beep.Buffer, error) {
	path := "music/" + fileName
	return loadAudioFromArch(modAudioArchPath, path)
}

// AudioEffect returns audio stream data from file with specified
// name inside audio archive.
func AudioEffect(fileName string) (*beep.Buffer, error) {
	path := "effect/" + fileName
	return loadAudioFromArch(modAudioArchPath, path)
}

// ErrorItemGraphic returns error graphic for item.
func ErrorItemGraphic() (*res.ItemGraphicData, error) {
	icon := graphic.Icons["unknown.png"]
	if icon == nil {
		return nil, fmt.Errorf("unable to retrieve error icon")
	}
	igd := res.ItemGraphicData{
		IconPic:  icon,
		MaxStack: 100,
	}
	return &igd, nil
}

// ScriptsDir returns all scripts from directory with
// specified path.
func ScriptsDir(path string) ([]*ash.Script, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir: %v", err)
	}
	scripts := make([]*ash.Script, 0)
	for _, info := range files {
		if !strings.HasSuffix(info.Name(), ashScriptExt) {
			continue
		}
		scriptPath := filepath.FromSlash(path + "/" + info.Name())
		s, err := Script(scriptPath)
		if err != nil {
			log.Err.Printf("data scripts dir: %s: unable to retrieve script: %v",
				path, err)
			continue
		}
		scripts = append(scripts, s)
	}
	return scripts, nil
}

// Script parses file with specified path to
// Ash scirpt.
func Script(path string) (*ash.Script, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	text, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %v", err)
	}
	scriptName := filepath.Base(path)
	script, err := ash.NewScript(scriptName, fmt.Sprintf("%s", text))
	if err != nil {
		return nil, fmt.Errorf("unable to parse script text: %v", err)
	}
	return script, nil
}

// Load loads grpahic directories.
func loadPaths(mod *module.Module) {
	modGraphicDirPath = filepath.Join("data/modules", mod.Conf().ID, "gui")
	modAudioDirPath = filepath.Join("data/modules", mod.Conf().ID, "gui")
	modGraphicArchPath = filepath.Join("data/modules", mod.Conf().ID, "gui/graphic.zip")
	modAudioArchPath = filepath.Join("data/modules", mod.Conf().ID, "gui/audio.zip")
}
