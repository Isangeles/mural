/*
 * mainmenu.go
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

// mainmenu package contains main menu and also settings,
// load/save and new game menus.
package mainmenu

import (
	"fmt"
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/isangeles/flame"
	flamecore "github.com/isangeles/flame/core"
	flamedata "github.com/isangeles/flame/core/data"
	flamesave "github.com/isangeles/flame/core/data/save" 

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/objects"
)

var (
	main_color   color.Color = colornames.Grey
	sec_color    color.Color = colornames.Blue
	accent_color color.Color = colornames.Red
)

// MainMenu struct reperesents container with
// all menu screens(settings menu, new game menu, etc.).
// Wraps all main menu screens.
type MainMenu struct {
	menu          *Menu
	newgamemenu   *NewGameMenu
	newcharmenu   *NewCharacterMenu
	loadgamemenu  *LoadGameMenu
	settings      *Settings
	console       *Console
	loadscreen    *LoadingScreen
	userFocus     *mtk.Focus
	msgs          *mtk.MessagesQueue
	PlayableChars []*objects.Avatar
	onGameCreated func(g *flamecore.Game, player *objects.Avatar)
	onGameLoaded  func(gameSav *flamesave.SaveGame)
	loading       bool
	exiting       bool
}

// New returns new main menu
func New() (*MainMenu, error) {
	mm := new(MainMenu)
	// Menu.
	m, err := newMenu(mm)
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_main_menu:%v",
			err)
	}
	mm.menu = m
	// New game menu.
	ngm, err := newNewGameMenu(mm)
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_new_game_menu:%v",
			err)
	}
	mm.newgamemenu = ngm
	// New character menu.
	ncm, err := newNewCharacterMenu(mm)
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_new_character_menu:%v",
			err)
	}
	mm.newcharmenu = ncm
	// Load game menu.
	lgm, err := newLoadGameMenu(mm)
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_load_game_menu:%v",
			err)
	}
	mm.loadgamemenu = lgm
	// Settings.
	s, err := newSettings(mm)
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_settings_menu:%v",
			err)
	}
	mm.settings = s
	// Console.
	c, err := newConsole()
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_console:%v",
			err)
	}
	mm.console = c
	// Loading screen.
	ls := newLoadingScreen(mm)
	mm.loadscreen = ls
	// Messages & focus.
	mm.userFocus = new(mtk.Focus)
	mm.msgs = mtk.NewMessagesQueue(mm.userFocus)

	mm.menu.Show(true)
	return mm, nil
}

// Draw draws current menu screen.
func (mm *MainMenu) Draw(win *mtk.Window) {
	if mm.loading {
		mm.loadscreen.Draw(win)
		return
	}
	// Menu.
	if mm.menu.Opened() {
		mm.menu.Draw(win)
	}
	// New game menu.
	if mm.newgamemenu.Opened() {
		mm.newgamemenu.Draw(win)
	}
	// New character menu.
	if mm.newcharmenu.Opened() {
		mm.newcharmenu.Draw(win)
	}
	// Load game menu.
	if mm.loadgamemenu.Opened() {
		mm.loadgamemenu.Draw(win)
	}
	// Settings.
	if mm.settings.Opened() {
		mm.settings.Draw(win.Window)
	}
	// Messages.
	mm.msgs.Draw(win.Window, mtk.Matrix().Moved(win.Bounds().Center()))
	// Console.
	if mm.console.Opened() {
		mm.console.Draw(win)
	}
}

// Update updates current menu screen.
func (mm *MainMenu) Update(win *mtk.Window) {
	if mm.exiting {
		win.SetClosed(true)
		return
	}
	if mm.loading {
		mm.loadscreen.Update(win)
	}
	if mm.menu.Opened() {
		mm.menu.Update(win)
	}
	if mm.newgamemenu.Opened() {
		mm.newgamemenu.Update(win)
	}
	if mm.newcharmenu.Opened() {
		mm.newcharmenu.Update(win)
	}
	if mm.loadgamemenu.Opened() {
		mm.loadgamemenu.Update(win)
	}
	if mm.settings.Opened() {
		mm.settings.Update(win)
	}
	mm.console.Update(win)
	mm.msgs.Update(win)
}

// Exit sends exit request to main menu.
func (mm *MainMenu) Exit() {
	mm.exiting = true
}

// SetOnGameCreatedFunc sets specified function as function
// triggered after new game created.
func (mm *MainMenu) SetOnGameCreatedFunc(f func(g *flamecore.Game,
	player *objects.Avatar)) {
	mm.onGameCreated = f
}

// SetOnGameLoadedFunc sets specified function as function
// triggered after game loaded.
func (mm *MainMenu) SetOnGameLoadedFunc(f func(gameSav *flamesave.SaveGame)) {
	mm.onGameLoaded = f
}

// OpenMenu opens menu.
func (mm *MainMenu) OpenMenu() {
	mm.HideMenus()
	mm.menu.Show(true)
}

// OpenNewGameMenu opens new game creation menu.
func (mm *MainMenu) OpenNewGameMenu() {
	mm.HideMenus()
	mm.newgamemenu.Show(true)
}

// OpenNewCharMenu opens new character creation menu.
func (mm *MainMenu) OpenNewCharMenu() {
	mm.HideMenus()
	mm.newcharmenu.Show(true)
}

// OpenLoadGameMenu opens load game menu.
func (mm *MainMenu) OpenLoadGameMenu() {
	mm.HideMenus()
	mm.loadgamemenu.Show(true)
}

// OpenSettings opens settings menu.
func (mm *MainMenu) OpenSettings() {
	mm.HideMenus()
	mm.settings.Show(true) 
}

// OpenLoadingScreen opens loading screen
// with specified loading information.
func (mm *MainMenu) OpenLoadingScreen(loadInfo string) {
	mm.loading = true
	mm.loadscreen.SetLoadInfo(loadInfo)
}

// CloseLoadingScreen closes loading screen.
func (mm *MainMenu) CloseLoadingScreen() {
	mm.loading = false
}

// HideMenus hides all menus.
func (mm *MainMenu) HideMenus() {
	mm.menu.Show(false)
	mm.newgamemenu.Show(false)
	mm.newcharmenu.Show(false)
	mm.loadgamemenu.Show(false)
	mm.settings.Show(false)
}

// ShowMessage adds specified message to messages queue
// and turns message visible(if not visible already).
func (mm *MainMenu) ShowMessage(m *mtk.MessageWindow) {
	m.Show(true)
	mm.msgs.Append(m)
}

// Console returns main menu console.
func (mm *MainMenu) Console() *Console {
	return mm.console
}

// AddPlaybaleChar adds new playable character to playable
// characters list.
func (mm *MainMenu) AddPlayableChar(c *objects.Avatar) {
	mm.PlayableChars = append(mm.PlayableChars, c)
}

// Triggered after new game was created(by new game creation menu).
func (mm *MainMenu) OnNewGameCreated(g *flamecore.Game,
	player *objects.Avatar) {
	if mm.onGameCreated == nil {
		return
	}
	mm.onGameCreated(g, player)
}

// Triggered after saved game was loaded.
func (mm *MainMenu) OnGameLoaded(gameSav *flamesave.SaveGame) {
	if mm.onGameLoaded == nil {
		return
	}
	mm.onGameLoaded(gameSav)
}

// ImportPlayableChars import all characters from specified
// path.
func (mm *MainMenu) ImportPlayableChars(path string) error {
	chars, err := flamedata.ImportCharactersDir(flame.Mod(), path)
	if err != nil {
		return fmt.Errorf("fail_to_import_characters:%v", err)
	}
	avs, err := data.ImportAvatarsDir(chars, path)
	if err != nil {
		return fmt.Errorf("fail_to_import_avatars:%v", err)
	}
	for _, av := range avs {
		mm.PlayableChars = append(mm.PlayableChars, av)
	}
	return nil
}

