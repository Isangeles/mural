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

	"github.com/golang/freetype/truetype"

	"github.com/salviati/go-tmx/tmx"

	"github.com/faiface/beep"
	
	"github.com/faiface/pixel"

	"github.com/isangeles/flame"

	"github.com/isangeles/mural/log"
)

var (
	// Paths.
	a_dir_path  string
	g_dir_path  string
	g_arch_path string
	a_arch_path string
	// Textures & fonts.
	uiTexs      map[string]pixel.Picture
	avatarsTexs map[string]pixel.Picture
	itemsTexs   map[string]pixel.Picture
	icons       map[string]pixel.Picture
	portraits   map[string]pixel.Picture
	music       map[string]*beep.Buffer
	fonts       map[string]*truetype.Font
)

// LoadGameData loads game graphic data.
// Should be called by GUI before creating any
// in-game elements.
func LoadGameData() error {
	// Load data resources paths.
	err := loadPaths()
	if err != nil {
		return fmt.Errorf("fail_to_load_paths:%v", err)
	}
	// Portrait textures.
	portraitsTexs, err := loadPicturesFromArch(g_arch_path, "avatar/portrait")
	if err != nil {
		return fmt.Errorf("fail_to_load_portraits:%v", err)
	}
	portraits = portraitsTexs
	// Avatars spritesheets.
	avTexs, err := loadPicturesFromArch(g_arch_path, "avatar/spritesheet")
	if err != nil {
		return fmt.Errorf("fail_to_load_avatars_spritesheets:%v", err)
	}
	avatarsTexs = avTexs
	// Items spritesheets.
	itTexs, err := loadPicturesFromArch(g_arch_path, "item/spritesheet")
	if err != nil {
		return fmt.Errorf("fail_to_load_items_spritesheets:%v", err)
	}
	itemsTexs = itTexs
	// Icons.
	icons = make(map[string]pixel.Picture, 0)
	itemIcons, err := loadPicturesFromArch(g_arch_path, "item/icon")
	if err != nil {
		return fmt.Errorf("fail_to_load_items_icons:%v", err)
	}
	for name, i := range itemIcons {
		icons[name] = i
	}
	effectIcons, err := loadPicturesFromArch(g_arch_path, "effect/icon")
	if err != nil {
		return fmt.Errorf("fail_to_load_effects_icons:%v", err)
	}
	for name, i := range effectIcons {
		icons[name] = i
	}
	skillIcons, err := loadPicturesFromArch(g_arch_path, "skill/icon")
	if err != nil {
		return fmt.Errorf("fail_to_load_skills_icons:%v", err)
	}
	for name, i := range skillIcons {
		icons[name] = i
	}
	return nil
}

// LoadUIData loads UI graphic data.
// Should be called by GUI before creating any
// GUI elements.
func LoadUIData() error {
	// Load data sources paths.
	err := loadPaths()
	if err != nil {
		return fmt.Errorf("fail_to_load_paths:%v", err)
	}
	// GUI elements textures.
	texs, err := loadPicturesFromArch(g_arch_path, "ui")
	if err != nil {
		return fmt.Errorf("fail_to_load_UI_textures:%v", err)
	}
	uiTexs = texs
	// Fonts.
	ttfs, err := loadFontsFromArch(g_arch_path, "font")
	if err != nil {
		return fmt.Errorf("fail_to_load_fonts:%v", err)
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
	log.Dbg.Printf("data_picture_ui_fallback_load:%s", fileName)
	return loadPictureFromArch(g_arch_path, "ui/"+fileName)
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
	path := filepath.FromSlash(flame.Mod().FullPath() + "/gui/portraits/" +
		fileName)
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
	log.Dbg.Printf("data_avatar_spritesheet_fallback_load:%s",
		fileName)
	path := filepath.FromSlash("avatar/spritesheet/" + fileName)
	return loadPictureFromArch(g_arch_path, path)
}

// ItemSpritesheet returns picture with specified name
// for item sprite.
func ItemSpritesheet(fileName string) (pixel.Picture, error) {
	spritesheet := itemsTexs[fileName]
	if spritesheet != nil {
		return spritesheet, nil
	}
	// Fallback.
	log.Dbg.Printf("data_items_spritesheet_fallback_load:%s",
		fileName)
	path := filepath.FromSlash("item/spritesheet/" + fileName)
	return loadPictureFromArch(g_arch_path, path)
}

// Icon returns picture with specified name for icon.
func Icon(fileName string) (pixel.Picture, error) {
	icon := icons[fileName]
	if icon != nil {
		return icon, nil
	}
	// Fallback.
	log.Dbg.Printf("data_items_icon_fallback_load:%s",
		fileName)
	path := filepath.FromSlash("item/icon/" + fileName)
	return loadPictureFromArch(g_arch_path, path)
}

// PlayablePortraits returns map with names of portraits as keys
// and portraits pictures as values avalible for player character.
func PlayablePortraits() (map[string]pixel.Picture, error) {
	path := flame.Mod().FullPath() + "/gui/portraits"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
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
	log.Dbg.Printf("data_font_fallback_load:%s", fileName)
	fullpath := fmt.Sprintf("%s/%s/%s", g_dir_path, "font", fileName)
	return loadFontFromDir(filepath.FromSlash(fullpath))
}

// Map loads TMX map from file with specified
// directory and name.
func Map(mapDir, mapName string) (*tmx.Map, error) {
	mapPath := filepath.FromSlash(mapDir + "/" + mapName + ".tmx")
	tmxFile, err := os.Open(mapPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_tmx_file:%v", err)
	}
	tmxMap, err := tmx.Read(tmxFile)
	if err != nil {
		return nil, fmt.Errorf("fail_to_read_tmx_file:%v", err)
	}
	return tmxMap, nil
}

// Music returns audio stream data from file with specified name
// inside audio archive.
func Music(fileName string) (*beep.Buffer, error) {
	path := filepath.FromSlash("music/" + fileName)
	return loadAudioFromArch(a_arch_path, path)
}

// AudioEffect returns audio stream data from file with specified
// name inside audio archive.
func AudioEffect(fileName string) (*beep.Buffer, error) {
	path := filepath.FromSlash("effect/" + fileName)
	return loadAudioFromArch(a_arch_path, path)
}

// Load loads grpahic directories.
func loadPaths() error {
	if flame.Mod() == nil {
		return fmt.Errorf("no module loaded")
	}
	g_dir_path = filepath.FromSlash(fmt.Sprintf("data/modules/%s/gui",
		flame.Mod().Name()))
	a_dir_path = filepath.FromSlash(fmt.Sprintf("data/modules/%s/gui",
		flame.Mod().Name()))
	g_arch_path = filepath.FromSlash(fmt.Sprintf("data/modules/%s/gui/gdata.zip",
		flame.Mod().Name()))
	a_arch_path = filepath.FromSlash(fmt.Sprintf("data/modules/%s/gui/adata.zip",
		flame.Mod().Name()))
	return nil
}
