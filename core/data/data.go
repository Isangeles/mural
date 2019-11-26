/*
 * data.go
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

// data package contains functions for loading
// graphics/audio/text data.
package data

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"

	"github.com/salviati/go-tmx/tmx"

	"github.com/faiface/beep"

	"github.com/faiface/pixel"

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/core/module"
	
	"github.com/isangeles/burn/ash"

	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/log"
)

var (
	// Paths.
	modAudioDirPath     string
	modGraphicDirPath   string
	modGraphicArchPath  string
	modAudioArchPath    string
	mainGraphicArchPath = "data/gui/gdata.zip"
	mainAudioArchPath   = "data/gui/adata.zip"
	// Textures & fonts.
	uiTexs      map[string]pixel.Picture
	avatarsTexs map[string]pixel.Picture
	objectsTexs map[string]pixel.Picture
	itemsTexs   map[string]pixel.Picture
	icons       map[string]pixel.Picture
	portraits   map[string]pixel.Picture
	music       map[string]*beep.Buffer
	fonts       map[string]*truetype.Font
)

// Init function.
func init() {
	uiTexs = make(map[string]pixel.Picture)
	avatarsTexs = make(map[string]pixel.Picture)
	objectsTexs = make(map[string]pixel.Picture)
	itemsTexs = make(map[string]pixel.Picture)
	icons = make(map[string]pixel.Picture)
	portraits = make(map[string]pixel.Picture)
	music = make(map[string]*beep.Buffer)
	fonts = make(map[string]*truetype.Font)
}

// LoadModuleData loads graphic data for specified module.
// Should be called by GUI before creating any
// in-game elements.
func LoadModuleData(mod *module.Module) error {
	// Load data resource paths.
	loadPaths(mod)
	// Portraits.
	avsPortraits, err := loadPicturesFromArch(modGraphicArchPath, "avatar/portrait")
	if err != nil {
		return fmt.Errorf("fail to load avatars portraits: %v", err)
	}
	for n, p := range avsPortraits {
		portraits[n] = p
	}
	obPortraits, err := loadPicturesFromArch(modGraphicArchPath, "object/portrait")
	if err != nil {
		return fmt.Errorf("fail to load objects portraits: %v", err)
	}
	for n, p := range obPortraits {
		portraits[n] = p
	}
	// Avatars spritesheets.
	avTexs, err := loadPicturesFromArch(modGraphicArchPath, "avatar/spritesheet")
	if err != nil {
		return fmt.Errorf("fail to load avatars spritesheets: %v", err)
	}
	for n, t := range avTexs {
		avatarsTexs[n] = t
	}
	// Objects spritesheets.
	obTexs, err := loadPicturesFromArch(modGraphicArchPath, "object/spritesheet")
	if err != nil {
		return fmt.Errorf("fail to load objects spritesheets: %v", err)
	}
	for n, t := range obTexs {
		objectsTexs[n] = t
	}
	// Items spritesheets.
	itTexs, err := loadPicturesFromArch(modGraphicArchPath, "item/spritesheet")
	if err != nil {
		return fmt.Errorf("fail to load items spritesheets: %v", err)
	}
	for n, t := range itTexs {
		itemsTexs[n] = t
	}
	// Icons.
	itemIcons, err := loadPicturesFromArch(modGraphicArchPath, "item/icon")
	if err != nil {
		return fmt.Errorf("fail to load items icons: %v", err)
	}
	for name, i := range itemIcons {
		icons[name] = i
	}
	effectIcons, err := loadPicturesFromArch(modGraphicArchPath, "effect/icon")
	if err != nil {
		return fmt.Errorf("fail to load effects icons: %v", err)
	}
	for name, i := range effectIcons {
		icons[name] = i
	}
	skillIcons, err := loadPicturesFromArch(modGraphicArchPath, "skill/icon")
	if err != nil {
		return fmt.Errorf("fail to load skills icons: %v", err)
	}
	for name, i := range skillIcons {
		icons[name] = i
	}
	return nil
}

// LoadUIData loads UI graphic data for specified module.
// Should be called by GUI before creating any
// GUI elements.
func LoadUIData(mod *module.Module) error {
	// Load data sources paths.
	loadPaths(mod)
	// GUI elements textures.
	texs, err := loadPicturesFromArch(modGraphicArchPath, "ui")
	if err != nil {
		return fmt.Errorf("fail to load ui textures: %v", err)
	}
	uiTexs = texs
	// Fonts.
	ttfs, err := loadFontsFromArch(modGraphicArchPath, "font")
	if err != nil {
		return fmt.Errorf("fail to load fonts: %v", err)
	}
	fonts = ttfs
	return nil
}

// PictureUI loads image with specified name from UI data
// in gdata archive.
func PictureUI(fileName string) (pixel.Picture, error) {
	pic := uiTexs[fileName]
	if pic != nil {
		return pic, nil
	}
	// Fallback, load picture 'by hand'.
	log.Dbg.Printf("data picture ui fallback load: %s", fileName)
	return loadPictureFromArch(modGraphicArchPath, "ui/"+fileName)
}

// PictureFromDir loads image from specified system path.
func PictureFromDir(path string) (pixel.Picture, error) {
	return loadPictureFromDir(path)
}

// AvatarPortrait returns portrait with specified name.
func AvatarPortrait(fileName string) (pixel.Picture, error) {
	portrait := portraits[fileName]
	if portrait != nil {
		return portrait, nil
	}
	path := filepath.FromSlash(flame.Mod().Conf().Path +
		"/gui/portraits/" + fileName)
	return loadPictureFromDir(path)
}

// Portrait returns portrait with specified name.
func Portrait(fileName string) (pixel.Picture, error) {
	portrait := portraits[fileName]
	if portrait != nil {
		return portrait, nil
	}
	path := filepath.FromSlash(flame.Mod().Conf().Path +
		"/gui/portraits/" + fileName)
	return loadPictureFromDir(path)
}

// AvatarSpritesheet returns picture with specified name
// for avatar sprite.
func AvatarSpritesheet(fileName string) (pixel.Picture, error) {
	spritesheet := avatarsTexs[fileName]
	if spritesheet != nil {
		return spritesheet, nil
	}
	// Fallback.
	log.Dbg.Printf("data avatar spritesheet fallback load: %s",
		fileName)
	path := filepath.FromSlash("avatar/spritesheet/" + fileName)
	return loadPictureFromArch(modGraphicArchPath, path)
}

// ItemSpritesheet returns picture with specified name
// for item sprite.
func ItemSpritesheet(fileName string) (pixel.Picture, error) {
	spritesheet := itemsTexs[fileName]
	if spritesheet != nil {
		return spritesheet, nil
	}
	// Fallback.
	log.Dbg.Printf("data items spritesheet fallback load: %s",
		fileName)
	path := filepath.FromSlash("item/spritesheet/" + fileName)
	return loadPictureFromArch(modGraphicArchPath, path)
}

// ObjectSpritesheet returns picture with specified name
// for object sprite.
func ObjectSpritesheet(fileName string) (pixel.Picture, error) {
	spritesheet := objectsTexs[fileName]
	if spritesheet != nil {
		return spritesheet, nil
	}
	// Fallback.
	log.Dbg.Printf("data objects spritesheet fallback load: %s",
		fileName)
	path := filepath.FromSlash("object/spritesheet/" + fileName)
	return loadPictureFromArch(modGraphicArchPath, path)
}

// Icon returns picture with specified name for icon.
func Icon(fileName string) (pixel.Picture, error) {
	icon := icons[fileName]
	if icon != nil {
		return icon, nil
	}
	// Fallback.
	log.Dbg.Printf("data items icon fallback load: %s",
		fileName)
	path := filepath.FromSlash("item/icon/" + fileName)
	return loadPictureFromArch(modGraphicArchPath, path)
}

// PlayablePortraits returns map with names of portraits as keys
// and portraits pictures as values avalible for player character.
func PlayablePortraits() (map[string]pixel.Picture, error) {
	path := flame.Mod().Conf().Path + "/gui/portraits"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("fail to read dir: %v", err)
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

// Font loads font with specified name from gdata
// directory.
func Font(fileName string) (*truetype.Font, error) {
	font := fonts[fileName]
	if font != nil {
		return font, nil
	}
	// Fallback.
	log.Dbg.Printf("data font fallback load: %s", fileName)
	fullpath := fmt.Sprintf("%s/%s/%s", modGraphicDirPath, "font", fileName)
	return loadFontFromDir(filepath.FromSlash(fullpath))
}

// Map loads TMX map from file with specified
// directory and name.
func Map(areaDir string) (*tmx.Map, error) {
	tmxFile, err := os.Open(areaDir + "/map.tmx")
	if err != nil {
		return nil, fmt.Errorf("fail to open tmx file: %v", err)
	}
	tmxMap, err := tmx.Read(tmxFile)
	if err != nil {
		return nil, fmt.Errorf("fail to read tmx file: %v", err)
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
	icon, err := Icon("unknown.png")
	if err != nil {
		return nil, fmt.Errorf("fail to retrieve error icon: %v", err)
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
		return nil, fmt.Errorf("fail to read dir: %v", err)
	}
	scripts := make([]*ash.Script, 0)
	for _, info := range files {
		if !strings.HasSuffix(info.Name(), ash.SCRIPT_FILE_EXT) {
			continue
		}
		scriptPath := filepath.FromSlash(path + "/" + info.Name())
		s, err := Script(scriptPath)
		if err != nil {
			log.Err.Printf("data scripts dir: %s: fail to retrieve script: %v",
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
		return nil, fmt.Errorf("fail to open file: %v", err)
	}
	text, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("fail to read file: %v", err)
	}
	script, err := ash.NewScript(fmt.Sprintf("%s", text), file.Name())
	if err != nil {
		return nil, fmt.Errorf("fail to parse script text: %v", err)
	}
	return script, nil
}

// Load loads grpahic directories.
func loadPaths(mod *module.Module) {
	modGraphicDirPath = filepath.FromSlash(fmt.Sprintf("data/modules/%s/gui",
		mod.Conf().ID))
	modAudioDirPath = filepath.FromSlash(fmt.Sprintf("data/modules/%s/gui",
		mod.Conf().ID))
	modGraphicArchPath = filepath.FromSlash(fmt.Sprintf("data/modules/%s/gui/gdata.zip",
		mod.Conf().ID))
	modAudioArchPath = filepath.FromSlash(fmt.Sprintf("data/modules/%s/gui/adata.zip",
		mod.Conf().ID))
}
