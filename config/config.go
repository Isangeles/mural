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

	flameconf "github.com/isangeles/flame/config"
	"github.com/isangeles/flame/core/data/text"
)

const (
	Name, Version = "Mural", "0.0.0"
	ConfFileName  = ".mural"
)

var (
	Fullscreen       = false
	MapFOW           = true
	Resolution       pixel.Vec
	MainFont         = ""
	MenuMusic        = ""
	ButtonClickSound = ""
	MusicVolume      = 0.0
	MusicMute        = false
	CharAttrsMin     = 1
	CharAttrsMax     = 10
	CharSkills       []string
	CharItems        []string
)

// LoadConfig loads configuration file.
func LoadConfig() error {
	values, err := text.ReadValue(ConfFileName, "fullscreen", "resolution",
		"map-fow", "main-font", "menu-music", "button-click-sound",
		"newchar-skills", "newchar-items", "music-volume", "music-mute")
	if err != nil {
		return fmt.Errorf("config_load: fail to retrieve config value: %v", err)
	}
	intValues, err := text.ReadInt(ConfFileName, "newchar-attrs-min", "newchar-attrs-max")
	if err != nil {
		return fmt.Errorf("config_load: fail to retrieve config value: %v", err)
	}
	// Fullscreen.
	Fullscreen = values["fullscreen"] == "true"
	// Resolution.
	resValue := values["resolution"]
	Resolution.X, err = strconv.ParseFloat(strings.Split(resValue, "x")[0], 64)
	Resolution.Y, err = strconv.ParseFloat(strings.Split(resValue, "x")[1], 64)
	if err != nil {
		log.Err.Printf("config_load: fail to set resolution: %v", err)
	}
	// Graphic effects.
	MapFOW = values["map-fow"] == "true"
	MainFont = values["main-font"]
	// Audio effects.
	MenuMusic = values["menu-music"]
	ButtonClickSound = values["button-click-sound"]
	MusicVolume, err = strconv.ParseFloat(values["music-volume"], 64)
	if err != nil {
		log.Err.Printf("config_load: fail to set music volume: %v", err)
	}
	MusicMute = values["music-mute"] == "true"
	// New char attributes points.
	CharAttrsMin = intValues["newchar-attrs-min"]
	CharAttrsMax = intValues["newchar-attrs-max"]
	// New char items & skills.
	CharSkills = strings.Split(values["newchar-skills"], ";")
	CharItems = strings.Split(values["newchar-items"], ";")
	log.Dbg.Print("config file loaded")
	return nil
}

// SaveConfig saves current configuration to file.
func SaveConfig() error {
	// Create file.
	f, err := os.Create(ConfFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	// Write config values.
	w := bufio.NewWriter(f)
	w.WriteString(fmt.Sprintf("%s\n", "# Mural GUI configuration file.")) // default header
	w.WriteString(fmt.Sprintf("fullscreen:%v\n", Fullscreen))
	w.WriteString(fmt.Sprintf("resolution:%fx%f\n", Resolution.X, Resolution.Y))
	w.WriteString(fmt.Sprintf("map-fow:%v\n", MapFOW))
	w.WriteString(fmt.Sprintf("main-font:%s\n", MainFont))
	w.WriteString(fmt.Sprintf("menu-music:%s\n", MenuMusic))
	w.WriteString(fmt.Sprintf("button-click-sound:%s\n", ButtonClickSound))
	w.WriteString(fmt.Sprintf("music-volume:%f\n", MusicVolume))
	w.WriteString(fmt.Sprintf("music-mute:%v\n", MusicMute))
	w.WriteString(fmt.Sprintf("newchar-attrs-min:%d\n", CharAttrsMin))
	w.WriteString(fmt.Sprintf("newchar-attrs-max:%d\n", CharAttrsMax))
	w.WriteString("newchar-skills:")
	for _, sid := range CharSkills {
		w.WriteString(sid + ";")
	}
	w.WriteString("\n")
	w.WriteString("newchar-items:")
	for _, iid := range CharItems {
		w.WriteString(iid + ";")
	}
	w.WriteString("\n")
	w.Flush()
	log.Dbg.Print("config file saved")
	return nil
}

// Debug checks whether debug mode is enabled.
func Debug() bool {
	return flameconf.Debug()
}

// Lang returns ID of current language.
func Lang() string {
	return flameconf.LangID()
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
