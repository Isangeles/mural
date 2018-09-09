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
	"github.com/faiface/pixel/pixelgl"
)

// MainMenu struct reperesents container with
// all menu screens(settings menu, new game menu, etc.)
type MainMenu struct {
	menu *Menu
}

// New returns new main menu
func New() (*MainMenu, error) {
	mm := new(MainMenu)
	m, err := newMenu()
	if err != nil {
		return nil, err
	}
	mm.menu = m
	return mm, nil
}

// Draw draws current menu screen.
func (mm *MainMenu) Draw(win *pixelgl.Window) {
	mm.menu.Draw(win)
}

// Update updates current menu screen.
func (mm *MainMenu) Update(win *pixelgl.Window) {
	mm.menu.Update(win)
}
