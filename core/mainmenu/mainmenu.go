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

// mainmenu package contains main menu, settings,
// load/save and new game screens.
package mainmenu

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/mural/core/mtk"
)

// MainMenu struct reperesents container with
// all menu screens(settings menu, new game menu, etc.).
type MainMenu struct {
	menu     *Menu
	settings *Settings
	console  *Console
	msgs     []*mtk.MessageWindow
}

// New returns new main menu
func New() (*MainMenu, error) {
	mm := new(MainMenu)
	// Menu.
	m, err := newMenu()
	if err != nil {
		return nil, err
	}
	m.OnSettingsButtonClickedFunc(mm.OpenSettings)
	mm.menu = m
	// Settings.
	s, err := newSettings()
	if err != nil {
		return nil, err
	}
	s.OnBackButtonClickedFunc(mm.OpenMenu)
	mm.settings = s
	// Console.
	c, err := newConsole()
	if err != nil {
		return nil, err
	}
	mm.console = c
	
	msg, err := mtk.NewMessageWindow(
		"This is test UI message.\nClick 'Ok' to dismiss.") // test
	if err != nil {
		return nil, err
	}
	msg.Show(true)
	mm.msgs = append(mm.msgs, msg)

	mm.menu.Show(true)
	return mm, nil
}

// Draw draws current menu screen.
func (mm *MainMenu) Draw(win *pixelgl.Window) {
	// Menu.
	if mm.menu.Open() {
		mm.menu.Draw(win)
	}
	// Settings.
	if mm.settings.Open() {
		mm.settings.Draw(win)
	}
	// Messages.
	for _, msg := range mm.msgs {
		if msg.Open() {
			msg.Draw(mtk.DisBL(win.Bounds(), 0.3),
				mtk.DisTR(win.Bounds(), 0.3), win)
		}
	}
	// Console.
	if mm.console.Open() {
		conBottomLeft := pixel.V(win.Bounds().Min.X, win.Bounds().Center().Y)
		mm.console.Draw(conBottomLeft, mtk.DisTR(win.Bounds(), 0), win)
	}
}

// Update updates current menu screen.
func (mm *MainMenu) Update(win *pixelgl.Window) {
	mm.menu.Update(win)
	mm.settings.Update(win)
	mm.console.Update(win)
	for i, msg := range mm.msgs {
		if msg.Open() {
			msg.Update(win)
		}
		if msg.Dismissed() {
			mm.msgs = append(mm.msgs[:i], mm.msgs[i+1:]...) // remove dismissed message
		}
	}
}

// CloseMenus closes all menus.
func (mm *MainMenu) CloseMenus() {
	mm.menu.Show(false)
	mm.settings.Show(false)
	mm.settings.ApplySettings()
}

// OpenMenu closes all currently open
// menus and opens main menu.
func (mm *MainMenu) OpenMenu(b *mtk.Button) {
	mm.CloseMenus()
	mm.menu.Show(true)
}

// OpenSettings closes all currently open
// menus and opens settings menu.
func (mm *MainMenu) OpenSettings(b *mtk.Button) { 
	mm.CloseMenus()
	mm.settings.Show(true) 
}
