/*
 * config.go
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

// Config package with configuration values.
package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/faiface/pixel"

	"github.com/isangeles/mural/log"

	flameconfig "github.com/isangeles/flame/config"
	"github.com/isangeles/flame/core/data/text"
)

const (
	NAME, VERSION  = "Mural", "0.0.0"
	CONF_FILE_NAME = ".mural"
)

var (
	fullscreen   = false
	mapfow       = true
	resolution   pixel.Vec
	lang         = flameconfig.LangID()
	mainFontName = ""
	menuMusic    = ""
	bClickSound  = ""
	attrsPtsMin  = 1
	attrsPtsMax  = 10
)

// LoadConfig loads configuration file.
func LoadConfig() error {
	confValues, err := text.ReadValue(CONF_FILE_NAME, "fullscreen",
		"resolution", "map-fow", "main-font", "menu-music", "button-click-sound")
	if err != nil {
		return fmt.Errorf("fail_to_retrieve_config_value:%v", err)
	}
	confInts, err := text.ReadInt(CONF_FILE_NAME, "newchar-attrs-min", "newchar-attrs-max")
	if err != nil {
		return fmt.Errorf("fail_to_retrieve_config_value:%v", err)
	}
	// Fullscreen.
	fullscreen = confValues["fullscreen"] == "true"
	// Resolution.
	resValue := confValues["resolution"]
	resolution.X, err = strconv.ParseFloat(strings.Split(resValue, "x")[0], 64)
	resolution.Y, err = strconv.ParseFloat(strings.Split(resValue, "x")[1], 64)
	if err != nil {
		log.Err.Printf("fail_to_set_custom_resolution:%s", resValue)
	}
	// Graphic effects.
	mapfow = confValues["map-fow"] == "true"
	mainFontName = confValues["main-font"]
	// Audio effects.
	menuMusic = confValues["menu-music"]
	bClickSound = confValues["button-click-sound"]
	// New chars attributes points.
	attrsPtsMin = confInts["newchar-attrs-min"]
	attrsPtsMax = confInts["newchar-attrs-max"]

	log.Dbg.Print("config file loaded")
	return nil
}

// SaveConfig saves current configuration to file.
func SaveConfig() error {
	// Create file.
	f, err := os.Create(CONF_FILE_NAME)
	if err != nil {
		return err
	}
	defer f.Close()
	// Write config values.
	w := bufio.NewWriter(f)
	w.WriteString(fmt.Sprintf("%s\n", "# Mural GUI configuration file.")) // default header
	w.WriteString(fmt.Sprintf("fullscreen:%v;\n", fullscreen))
	w.WriteString(fmt.Sprintf("resolution:%fx%f;\n", resolution.X, resolution.Y))
	w.WriteString(fmt.Sprintf("map-fow:%v;\n", mapfow))
	w.WriteString(fmt.Sprintf("main-font:%s;\n", mainFontName))
	w.WriteString(fmt.Sprintf("menu-music:%s;\n", menuMusic))
	w.WriteString(fmt.Sprintf("button-click-sound:%s;\n", bClickSound))
	w.WriteString(fmt.Sprintf("newchar-attrs-min:%d;\n", attrsPtsMin))
	w.WriteString(fmt.Sprintf("newchar-attrs-max:%d;\n", attrsPtsMax))
	w.Flush()
	log.Dbg.Print("config file saved")
	return nil
}

// Fullscreen returns fullscreen config value.
func Fullscreen() bool {
	return fullscreen
}

// Returns current resolution width and height.
func Resolution() pixel.Vec {
	return resolution
}

// Lang returns current language ID.
func Lang() string {
	return lang
}

// Debug checks whether debug mode is enabled.
func Debug() bool {
	return flameconfig.Debug()
}

// MapFOW checks whether map 'Fog Of War' effect
// in enabled.
func MapFOW() bool {
	return mapfow
}

// MainFontName returns name of main font
// for UI.
func MainFontName() string {
	return mainFontName
}

// MenuMusicFile returns name of audio file
// with main menu music theme.
func MenuMusicFile() string {
	return menuMusic
}

// ButtonClickSoundFile returns name of audio file
// with button click audio effect.
func ButtonClickSoundFile() string {
	return bClickSound
}

// NewCharAttrsMin returns minimal
// amount of attributes points for
// new character.
func NewCharAttrsMin() int {
	return attrsPtsMin
}

// NewCharAttrsMax returns maximal
// amount of attributes points for
// new character.
func NewCharAttrsMax() int {
	return attrsPtsMax
}

// SetFullscreen toggles fullscreen mode.
func SetFullscreen(fs bool) {
	fullscreen = fs
}

// SetResolution sets specified XY size as current
// resolution.
func SetResolution(res pixel.Vec) {
	resolution = res
}

// SetLang sets language with specified ID as current
// language.
func SetLang(langID string) {
	_ = flameconfig.SetLang(langID)
}

// SetMapFOW toggles map 'Fog Of War' graphical
// effect.
func SetMapFOW(fow bool) {
	mapfow = fow
}

// SetMainFotName set specified font file name
// as name for main UI font.
func SetMainFontName(font string) {
	mainFontName = font
}

// SetMenuMusicFile sets name of audio file with
// main menu music theme.
func SetMenuMusicFile(fileName string) {
	menuMusic = fileName
}

// SetButtonClickSoundFile sets specified file name
// as name of audio file with global button click
// effect.
func SetButtonClickSoundFile(fileName string) {
	bClickSound = fileName
}

// SupportedResolutions returns all resolutions
// supported by UI.
func SupportedResolutions() []pixel.Vec {
	return []pixel.Vec{pixel.V(1920, 1080), pixel.V(1300, 720)}
}

// SuportedLangs retruns all languages supported by UI.
func SupportedLangs() []string {
	return []string{"english"}
}
