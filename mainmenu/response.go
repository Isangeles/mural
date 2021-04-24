/*
 * response.go
 *
 * Copyright 2020-2021 Dariusz Sikora <dev@isangeles.pl>
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

package mainmenu

import (
	"path/filepath"

	"github.com/isangeles/flame"
	flameres "github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/serial"

	"github.com/isangeles/fire/response"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/game"
	"github.com/isangeles/mural/log"
)

// handleResponse handles specified response from Fire server.
func (mm *MainMenu) handleResponse(resp response.Response) {
	if !resp.Logon {
		if len(resp.Load.Save) > 0 {
			mm.handleLoadResponse(resp.Load)
		}
		mm.handleUpdateResponse(resp.Update)
		for _, r := range resp.Character {
			mm.handleCharacterResponse(r)
		}
	}
	for _, r := range resp.Error {
		log.Err.Printf("Main menu: server error: %v", r)
	}
}

// handleNewCharResponse handles new char response.
func (mm *MainMenu) handleCharacterResponse(resp response.Character) {
	if mm.mod == nil {
		return
	}
	for _, c := range mm.continueChars {
		if c.ID() == resp.ID && c.Serial() == resp.Serial {
			return
		}
	}
	for _, c := range mm.mod.Chapter().Characters() {
		if c.ID() == resp.ID && c.Serial() == resp.Serial {
			mm.continueChars = append(mm.continueChars, c)
		}
	}
}

// handleLoadResponse handles load response.
func (mm *MainMenu) handleLoadResponse(resp response.Load) {
	// Recreate saved game.
	flameres.Clear()
	serial.Reset()
	flameres.TranslationBases = res.TranslationBases()
	m := flame.NewModule()
	m.Apply(resp.Module)
	gameWrapper := game.New(m)
	gameWrapper.SetServer(mm.server)
	// Import saved HUD state.
	hudPath := filepath.Join(mm.mod.Conf().Path, data.SavesModulePath,
		resp.Save+data.HUDFileExt)
	hud, err := data.ImportHUD(hudPath)
	if err != nil {
		log.Err.Printf("Main menu: handle load response: unable to import HUD: %v", err)
		return
	}
	// Run on game created function.
	if mm.onGameCreated != nil {
		mm.onGameCreated(gameWrapper, &hud)
	}
}

// handleUpdateResponse handles update response.
func (mm *MainMenu) handleUpdateResponse(resp response.Update) {
	flameres.Clear()
	flameres.TranslationBases = res.TranslationBases()
	if mm.mod == nil {
		serial.Reset()
		mm.mod = flame.NewModule()
		mm.mod.Apply(resp.Module)
		return
	}
	mm.mod.Apply(resp.Module)
}
