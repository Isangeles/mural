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
	
	"github.com/faiface/pixel"

	flamecore "github.com/isangeles/flame/core"
	
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/objects"
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
	settings      *Settings
	console       *Console
	userFocus     *mtk.Focus
	msgs          *mtk.MessagesQueue
	PlayableChars []*objects.Avatar
	onGameCreated func(g *flamecore.Game, player *objects.Avatar)
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
	// Messages & focus.
	mm.userFocus = new(mtk.Focus)
	mm.msgs = mtk.NewMessagesQueue(mm.userFocus)

	mm.menu.Show(true)
	return mm, nil
}

// Draw draws current menu screen.
func (mm *MainMenu) Draw(win *mtk.Window) {
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
	// Settings.
	if mm.settings.Opened() {
		mm.settings.Draw(win.Window)
	}
	// Messages.
	mm.msgs.Draw(win.Window, mtk.Matrix().Moved(win.Bounds().Center()))
	// Console.
	if mm.console.Opened() {
		conBottomLeft := pixel.V(win.Bounds().Min.X,
			win.Bounds().Center().Y)
		mm.console.Draw(conBottomLeft, mtk.DisTR(win.Bounds(), 0),
			win.Window)
	}
}

// Update updates current menu screen.
func (mm *MainMenu) Update(win *mtk.Window) {
	mm.menu.Update(win)
	mm.newgamemenu.Update(win)
	mm.newcharmenu.Update(win)
	mm.settings.Update(win)
	mm.console.Update(win)
	mm.msgs.Update(win)
}

// SetOnGameCreatedFunc sets specified function as function
// triggered after new game created.
func (mm *MainMenu) SetOnGameCreatedFunc(f func(g *flamecore.Game,
	player *objects.Avatar)) {
	mm.onGameCreated = f
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


// OpenSettings opens settings menu.
func (mm *MainMenu) OpenSettings() {
	mm.HideMenus()
	mm.settings.Show(true) 
}

// HideMenus hides all menus.
func (mm *MainMenu) HideMenus() {
	mm.menu.Show(false)
	mm.newgamemenu.Show(false)
	mm.newcharmenu.Show(false)
	mm.settings.Show(false)
}

// ShowMessage adds specified message to messages queue
// and turns message visible(if not visible already).
func (mm *MainMenu) ShowMessage(m *mtk.MessageWindow) {
	m.Show(true)
	mm.msgs.Append(m)
}

// AddPlaybaleChar adds new playable character to playable
// characters list.
func (mm *MainMenu) AddPlayableChar(c *objects.Avatar) {
	mm.PlayableChars = append(mm.PlayableChars, c)
}

// Triggered after new game was created(by new game creation menu).
func (mm *MainMenu) OnNewGameCreated(g *flamecore.Game,
	player *objects.Avatar) {
	mm.onGameCreated(g, player)
}

