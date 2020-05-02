/*
 * config.go
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

// Config package with configuration values.
package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/faiface/pixel"

	"github.com/isangeles/mural/log"

	flameconf "github.com/isangeles/flame/config"
	"github.com/isangeles/flame/data/text"
)

const (
	Name, Version = "Mural", "0.0.0"
	ConfFileName  = ".mural"
)

var (
	Fullscreen       = false
	MapFOW           = true
	MapFull          = true
	Resolution       pixel.Vec
	MainFont         = ""
	MenuMusic        = ""
	ButtonClickSound = ""
	MusicVolume      = 0.0
	MusicMute        = false
)

// Load loads configuration file.
func Load() error {
	file, err := os.Open(ConfFileName)
	if err != nil {
		return fmt.Errorf("unable to open config file: %v", err)
	}
	conf, err := text.UnmarshalConfig(file)
	if err != nil {
		return fmt.Errorf("unable to unmarshal config: %v", err)
	}
	// Fullscreen.
	if len(conf["fullscreen"]) > 0 {
		Fullscreen = conf["fullscreen"][0] == "true"
	}
	// Resolution.
	if len(conf["resolution"]) > 1 {
		Resolution.X, err = strconv.ParseFloat(conf["resolution"][0], 64)
		if err != nil {
			log.Err.Printf("config: unable to set resolution x: %v", err)
		}
		Resolution.Y, err = strconv.ParseFloat(conf["resolution"][1], 64)
		if err != nil {
			log.Err.Printf("config: unable to set resolution y: %v", err)
		}
	}
	// Graphic effects.
	if len(conf["map-fow"]) > 0 {
		MapFOW = conf["map-fow"][0] == "true"
	}
	if len(conf["map-full"]) > 0 {
		MapFull = conf["map-full"][0] == "true"
	}
	if len(conf["main-font"]) > 0 {
		MainFont = conf["main-font"][0]
	}
	// Audio effects.
	if len(conf["menu-music"]) > 0 {
		MenuMusic = conf["menu-music"][0]
	}
	if len(conf["button-click-sound"]) > 0 {
		ButtonClickSound = conf["button-click-sound"][0]
	}
	if len(conf["music-volume"]) > 0 {
		MusicVolume, err = strconv.ParseFloat(conf["music-volume"][0], 64)
		if err != nil {
			log.Err.Printf("config: unable to set music volume: %v", err)
		}
	}
	if len(conf["music-mute"]) > 0 {
		MusicMute = conf["music-mute"][0] == "true"
	}
	// Debug.
	log.Dbg.Print("Config file loaded")
	return nil
}

// Save saves current configuration to file.
func Save() error {
	// Create file.
	file, err := os.Create(ConfFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	// Create config text.
	conf := make(map[string][]string)
	conf["fullscreen"] = []string{fmt.Sprintf("%v", Fullscreen)}
	conf["resolution"] = []string{
		fmt.Sprintf("%f", Resolution.X),
		fmt.Sprintf("%f", Resolution.Y),
	}
	conf["map-fow"] = []string{fmt.Sprintf("%v", MapFOW)}
	conf["map-full"] = []string{fmt.Sprintf("%v", MapFull)}
	conf["main-font"] = []string{MainFont}
	conf["menu-music"] = []string{MenuMusic}
	conf["button-click-sound"] = []string{ButtonClickSound}
	conf["music-volume"] = []string{fmt.Sprintf("%f", MusicVolume)}
	conf["music-mute"] = []string{fmt.Sprintf("%v", MusicMute)}
	confText := text.MarshalConfig(conf)
	// Write config values.
	w := bufio.NewWriter(file)
	w.WriteString(confText)
	w.Flush()
	// Debug.
	log.Dbg.Print("Config file saved")
	return nil
}

// Debug checks whether debug mode is enabled.
func Debug() bool {
	return flameconf.Debug
}

// Lang returns ID of current language.
func Lang() string {
	return flameconf.Lang
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
