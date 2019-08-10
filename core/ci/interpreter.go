/*
 * interpreter.go
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

// ci package provides GUI-specific command line tools
// for Burn CI.
package ci

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/isangeles/burn"
	"github.com/isangeles/burn/ash"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/hud"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/mainmenu"
)

const (
	GUIShow   = "guishow"
	GUISet    = "guiset"
	GUIExport = "guiexport"
	GUIImport = "guiimport"
	GUIAudio  = "guiaudio"
)

var (
	guiMenu  *mainmenu.MainMenu
	guiHUD   *hud.HUD
	guiMusic *mtk.AudioPlayer
)

// On init.
func init() {
	burn.AddToolHandler(GUIAudio, guiaudio)
	burn.AddToolHandler(GUIShow, guishow)
	burn.AddToolHandler(GUISet, guiset)
	burn.AddToolHandler(GUISet, guiset)
	burn.AddToolHandler(GUIExport, guiexport)
	burn.AddToolHandler(GUIImport, guiimport)
}

// RunScriptsDir runs in background all Ash scripts
// in directory with specified path.
func RunScriptsDir(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("fail to read dir: %v", err)
	}
	for _, finfo := range files {
		if !strings.HasSuffix(finfo.Name(), ash.SCRIPT_FILE_EXT) {
			continue
		}
		filepath := filepath.FromSlash(path + "/" + finfo.Name())
		err := RunScript(filepath)
		if err != nil {
			log.Err.Printf("ci: script: %s: %v", err)
		}
	}
	return nil
}

// RunScript runs in background Asg script
// from specified path.
func RunScript(path string, args ...string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("fail to open file: %v", err)
	}
	text, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("fail to read file: %v", err)
	}
	scriptPath := strings.Split(path, "/")
	scriptName := scriptPath[len(scriptPath)-1]
	scriptArgs := []string{scriptName}
	for _, a := range args {
		scriptArgs = append(scriptArgs, a)
	}
	script, err := ash.NewScript(fmt.Sprintf("%s", text), scriptArgs...)
	if err != nil {
		return fmt.Errorf("fail to create ash script: %v", err)
	}
	go runScript(script)
	return nil
}

// SetMainMenu sets specified main menu as main
// menu for guiman to manage.
func SetMainMenu(menu *mainmenu.MainMenu) {
	guiMenu = menu
}

// SetHUD sets specified HUD as HUD for
// guiman to manage.
func SetHUD(h *hud.HUD) {
	guiHUD = h
}

// SetMusicPlayer sets specified audio player as
// player for guiman to manage.
func SetMusicPlayer(p *mtk.AudioPlayer) {
	guiMusic = p
}

// runScript executes specified script,
// in case of error sends err message to
// Mural log.
func runScript(s *ash.Script) {
	err := ash.Run(s)
	if err != nil {
		log.Err.Printf("ci: fail to run script: %v", err)
		return
	}
}
