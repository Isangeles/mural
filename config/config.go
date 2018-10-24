/*
 * config.go
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
)

const (
	NAME, VERSION = "Mural", "0.0.0"
	CONF_FILE_NAME string = ".mural"
)

var (
	fullscreen bool
	resolution pixel.Vec
	lang = flame.LangID()
)

// LoadConfig loads configuration file.
func LoadConfig() error {
	confValues, err := text.ReadConfigValue(CONF_FILE_NAME, "fullscreen",
		"resolution")
	if err != nil {
		return err
	}

	fullscreen = confValues[0] == "true"
	
	resValue := confValues[1]
	resolution.X, err = strconv.ParseFloat(strings.Split(resValue,
		"x")[0], 64)
	resolution.Y, err = strconv.ParseFloat(strings.Split(resValue,
		"x")[1], 64)
	if err != nil {
		log.Err.Printf("fail_to_set_custom_resolution:%s", resValue)
	}
	
	log.Dbg.Print("config_file_loaded")
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
	w.WriteString(fmt.Sprintf("%s\n", "#Mural GUI configuration file.")) // default header
	w.WriteString(fmt.Sprintf("fullscreen:%v;\n", fullscreen))
	w.WriteString(fmt.Sprintf("resolution:%fx%f;\n", resolution.X,
		resolution.Y))
	w.Flush()

	log.Dbg.Print("config_file_saved")
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

// SetFullscreen toggles fullscreen mode.
func SetFullscreen(fs bool) {
	fullscreen = fs
}

// SetResolution sets specified XY size as current
// resolution.
func SetResolution(res pixel.Vec) {
	resolution = res
}

// SetLang sets language with specified ID as current language.
func SetLang(langID string) {
	// TODO: set flame lang ID.
}

// SupportedResolutions returns all resolutions supported by UI.
func SupportedResolutions() []pixel.Vec {
	return []pixel.Vec{pixel.V(1920, 1080), pixel.V(1300, 720)}
}

// SuportedLangs retruns all languages supported by UI.
func SupportedLangs() []string {
	return []string{"english"}
}
