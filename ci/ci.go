/*
 * ci.go
 *
 * Copyright 2018-2024 Dariusz Sikora <ds@isangeles.dev>
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
	"github.com/isangeles/burn"
	"github.com/isangeles/burn/ash"

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

// RunScript executes specified script, in case
// of error sends err message to Mural error log.
func RunScript(s *ash.Script) {
	err := ash.Run(s)
	if err != nil {
		log.Err.Printf("ci: fail to run script: %v", err)
		return
	}
}
