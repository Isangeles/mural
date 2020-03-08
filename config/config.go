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
	"github.com/isangeles/flame/data/parsetxt"
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
	CharArea         string
	CharPos          pixel.Vec
)

// LoadConfig loads configuration file.
func LoadConfig() error {
	file, err := os.Open(ConfFileName)
	if err != nil {
		return fmt.Errorf("unable to open config file: %v", err)
	}
	conf := parsetxt.UnmarshalConfig(file)
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
	// New char attributes points.
	if len(conf["newchar-attrs-min"]) > 0 {
		CharAttrsMin, err = strconv.Atoi(conf["newchar-attrs-min"][0])
		if err != nil {
			log.Err.Printf("config: unable to set char attrs min value: %v", err)
		}
	}
	if len(conf["newchar-attrs-max"]) > 0 {
		CharAttrsMax, err = strconv.Atoi(conf["newchar-attrs-max"][0])
		if err != nil {
			log.Err.Printf("config: unable to set char attrs max value: %v", err)
		}
	}
	// New char items & skills.
	CharSkills = conf["newchar-skills"]
	CharItems = conf["newchar-items"]
	// New char area & position.
	if len(conf["newchar-area"]) > 0 {
		CharArea = conf["newchar-area"][0]
	}
	if len(conf["newchar-pos"]) > 1 {
		CharPos.X, err = strconv.ParseFloat(conf["newchar-pos"][0], 64)
		if err != nil {
			log.Err.Printf("conf: unable to set new char position x: %v", err)
		}
		CharPos.Y, err = strconv.ParseFloat(conf["newchar-pos"][1], 64)
		if err != nil {
			log.Err.Printf("conf: unable to set new char position y: %v", err)
		}
	}
	log.Dbg.Print("config file loaded")
	return nil
}

// SaveConfig saves current configuration to file.
func SaveConfig() error {
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
	conf["main-font"] = []string{MainFont}
	conf["menu-music"] = []string{MenuMusic}
	conf["button-click-sound"] = []string{ButtonClickSound}
	conf["music-volume"] = []string{fmt.Sprintf("%f", MusicVolume)}
	conf["music-mute"] = []string{fmt.Sprintf("%v", MusicMute)}
	conf["newchar-attrs-min"] = []string{fmt.Sprintf("%d", CharAttrsMin)}
	conf["newchar-attrs-max"] = []string{fmt.Sprintf("%d", CharAttrsMax)}
	conf["newchar-skills"] = CharSkills
	conf["newchar-items"] = CharItems
	conf["newchar-area"] = []string{CharArea}
	conf["newchar-pos"] = []string{
		fmt.Sprintf("%f", CharPos.X),
		fmt.Sprintf("%f", CharPos.Y),
	}
	confText := parsetxt.MarshalConfig(conf)
	// Write config values.
	w := bufio.NewWriter(file)
	w.WriteString(confText)
	w.Flush()
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
