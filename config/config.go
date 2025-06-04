/*
 * config.go
 *
 * Copyright 2018-2025 Dariusz Sikora <ds@isangeles.dev>
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

// Package with configuration values.
package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gopxl/pixel"

	"github.com/isangeles/mural/log"

	"github.com/isangeles/flame/data/text"
)

const (
	Name, Version = "Mural", "0.1.0-dev"
	ConfFileName  = ".mural"
)

var (
	Lang             = "english"
	Module           = ""
	ModulesPath      = "data/modules"
	GUIPath          = ""
	DefaultHUD       = "default.json"
	Debug            = true
	Fullscreen       = false
	MapFOW           = true
	MapFull          = true
	Resolution       pixel.Vec
	MaxFPS           = 60
	MainFont         = ""
	MenuMusic        = ""
	ButtonClickSound = ""
	EffectsVolume    = 0.0
	EffectsMute      = false
	MusicVolume      = 0.0
	MusicMute        = false
	ServerLogin      = ""
	ServerPassword   = ""
	ServerHost       = ""
	ServerPort       = ""
	ServerClose      = false
)

// Load loads configuration file.
func Load() error {
	file, err := os.Open(ConfFileName)
	if err != nil {
		return fmt.Errorf("Unable to open config file: %v", err)
	}
	conf, err := text.UnmarshalConfig(file)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal config: %v", err)
	}
	// Language.
	if len(conf["lang"]) > 0 {
		Lang = conf["lang"][0]
	}
	// Module.
	if len(conf["module"]) > 0 {
		Module = conf["module"][0]
		GUIPath = filepath.Join(ModulesPath, Module, "mural")
	}
	if len(conf["gui-path"]) > 0 {
		GUIPath = conf["gui-path"][0]
	}
	// Debug.
	if len(conf["debug"]) > 0 {
		Debug = conf["debug"][0] == "true"
	}
	// Fullscreen.
	if len(conf["fullscreen"]) > 0 {
		Fullscreen = conf["fullscreen"][0] == "true"
	}
	// Resolution.
	if len(conf["resolution"]) > 1 {
		Resolution.X, err = strconv.ParseFloat(conf["resolution"][0], 64)
		if err != nil {
			log.Err.Printf("Config: Unable to set resolution x: %v", err)
		}
		Resolution.Y, err = strconv.ParseFloat(conf["resolution"][1], 64)
		if err != nil {
			log.Err.Printf("Config: Unable to set resolution y: %v", err)
		}
	}
	// Max FPS.
	if len(conf["max-fps"]) > 0 {
		MaxFPS, err = strconv.Atoi(conf["max-fps"][0])
		if err != nil {
			log.Err.Printf("Config: Unable to set max FPS: %v", err)
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
	if len(conf["effects-volume"]) > 0 {
		EffectsVolume, err = strconv.ParseFloat(conf["effects-volume"][0], 64)
		if err != nil {
			log.Err.Printf("Config: Unable to set effects volume: %v", err)
		}
	}
	if len(conf["effects-mute"]) > 0 {
		EffectsMute = conf["effects-mute"][0] == "true"
	}
	if len(conf["music-volume"]) > 0 {
		MusicVolume, err = strconv.ParseFloat(conf["music-volume"][0], 64)
		if err != nil {
			log.Err.Printf("Config: Unable to set music volume: %v", err)
		}
	}
	if len(conf["music-mute"]) > 0 {
		MusicMute = conf["music-mute"][0] == "true"
	}
	// Server.
	if len(conf["server-user"]) > 1 {
		ServerLogin = conf["server-user"][0]
		ServerPassword = conf["server-user"][1]
	}
	if len(conf["server"]) > 1 {
		ServerHost = conf["server"][0]
		ServerPort = conf["server"][1]
	}
	if len(conf["server-close"]) > 0 {
		ServerClose = conf["server-close"][0] == "true"
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
		return fmt.Errorf("Unable to create config file: %v", err)
	}
	defer file.Close()
	// Marshal config.
	conf := make(map[string][]string)
	conf["lang"] = []string{Lang}
	conf["module"] = []string{Module}
	conf["gui-path"] = []string{GUIPath}
	conf["debug"] = []string{fmt.Sprintf("%v", Debug)}
	conf["fullscreen"] = []string{fmt.Sprintf("%v", Fullscreen)}
	conf["resolution"] = []string{
		fmt.Sprintf("%f", Resolution.X),
		fmt.Sprintf("%f", Resolution.Y),
	}
	conf["max-fps"] = []string{fmt.Sprintf("%v", MaxFPS)}
	conf["map-fow"] = []string{fmt.Sprintf("%v", MapFOW)}
	conf["map-full"] = []string{fmt.Sprintf("%v", MapFull)}
	conf["main-font"] = []string{MainFont}
	conf["menu-music"] = []string{MenuMusic}
	conf["button-click-sound"] = []string{ButtonClickSound}
	conf["effects-volume"] = []string{fmt.Sprintf("%f", EffectsVolume)}
	conf["effects-mute"] = []string{fmt.Sprintf("%v", EffectsMute)}
	conf["music-volume"] = []string{fmt.Sprintf("%f", MusicVolume)}
	conf["music-mute"] = []string{fmt.Sprintf("%v", MusicMute)}
	conf["server-user"] = []string{ServerLogin, ServerPassword}
	conf["server"] = []string{ServerHost, ServerPort}
	conf["server-close"] = []string{fmt.Sprintf("%v", ServerClose)}
	confText := text.MarshalConfig(conf)
	// Write config values.
	w := bufio.NewWriter(file)
	w.WriteString(confText)
	w.Flush()
	// Debug.
	log.Dbg.Print("Config file saved")
	return nil
}

// LangPath returns path to a UI lang directory.
func LangPath() string {
	return filepath.Join("data/lang", Lang)
}

// ModulePath returns path to directory of the current module.
func ModulePath() string {
	return filepath.Join(ModulesPath, Module)
}

// SupportedResolutions returns all resolutions
// supported by the UI.
func SupportedResolutions() []pixel.Vec {
	return []pixel.Vec{pixel.V(1920, 1080), pixel.V(1300, 720)}
}

// SuportedLangs retruns all languages supported by the UI.
func SupportedLangs() []string {
	return []string{"english"}
}
