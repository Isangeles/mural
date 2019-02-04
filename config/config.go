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

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/core/data/text"
	"github.com/isangeles/flame/core/enginelog"
)

const (
	NAME, VERSION = "Mural", "0.0.0"
	CONF_FILE_NAME = ".mural"
)

var (
	fullscreen bool
	mapfow     bool = true
	resolution pixel.Vec
	lang = flame.LangID()
	mainFontName = ""
	menuMusic = ""
	bClickSound = ""
)

// LoadConfig loads configuration file.
func LoadConfig() error {
	confValues, err := text.ReadConfigValue(CONF_FILE_NAME, "fullscreen",
		"resolution", "map_fow", "main_font", "menu_music", "button_click_sound")
	if err != nil {
		return err
	}
	// Fullscreen.
	fullscreen = confValues[0] == "true"
	// Resolution.
	resValue := confValues[1]
	resolution.X, err = strconv.ParseFloat(strings.Split(resValue,
		"x")[0], 64)
	resolution.Y, err = strconv.ParseFloat(strings.Split(resValue,
		"x")[1], 64)
	if err != nil {
		log.Err.Printf("fail_to_set_custom_resolution:%s", resValue)
	}
	// Map FOW effect.
	mapfow = confValues[2] == "true"
	mainFontName = confValues[3]
	menuMusic = confValues[4]
	bClickSound = confValues[5]
	
	log.Dbg.Print("config file loaded")
	return nil
}

// SaveConfig saves current configuration to file.
func SaveConfig() error {
	f, err := os.Create(CONF_FILE_NAME)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(fmt.Sprintf("%s\n", "# Mural GUI configuration file.")) // default header
	w.WriteString(fmt.Sprintf("fullscreen:%v;\n", fullscreen))
	w.WriteString(fmt.Sprintf("resolution:%fx%f;\n", resolution.X, resolution.Y))
	w.WriteString(fmt.Sprintf("map_fow:%v;\n", mapfow))
	w.WriteString(fmt.Sprintf("main_font:%s;\n", mainFontName))
	w.WriteString(fmt.Sprintf("menu_music:%s;\n", menuMusic))
	w.WriteString(fmt.Sprintf("button_click_sound:%s;\n", bClickSound))
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
	return enginelog.IsDebug()
}

// MapFOW checks whether map 'Fog Of War' effect
// in enabled.
func MapFOW() bool {
	return mapfow
}

// MainFontName returns name of main font
// for UI.
func MainFontName() string  {
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
	_ = flame.SetLang(langID)
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
