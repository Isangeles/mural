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

	"github.com/isangeles/mural/core"
)

// MainMenu struct reperesents container with
// all menu screens(settings menu, new game menu, etc.).
type MainMenu struct {
	menu    *Menu
	console *Console
	msgs    []*core.MessageWindow
}

// New returns new main menu
func New() (*MainMenu, error) {
	mm := new(MainMenu)
	
	m, err := newMenu()
	if err != nil {
		return nil, err
	}
	mm.menu = m

	c, err := newConsole()
	if err != nil {
		return nil, err
	}
	mm.console = c

	msg, err := core.NewMessageWindow("TEST TEST TEST")
	if err != nil {
		return nil, err
	}
	msg.Show(true)
	mm.msgs = append(mm.msgs, msg)

	return mm, nil
}

// Draw draws current menu screen.
func (mm *MainMenu) Draw(win *pixelgl.Window) {
	// Menu.
	mm.menu.Draw(win)
	// Console.
	if mm.console.Open() {
		conDrawMin := pixel.V(win.Bounds().Min.X, win.Bounds().Center().Y)
		conDrawMax := pixel.V(win.Bounds().Max.X, win.Bounds().Max.Y)
		mm.console.Draw(conDrawMin, conDrawMax, win)
	}
	// Messages.
	for _, msg := range mm.msgs {
		if msg.Open() {
			msgDrawMin := pixel.V(win.Bounds().Min.X + (win.Bounds().Max.X / 3),
				win.Bounds().Min.Y + (win.Bounds().Max.Y / 3))
			msgDrawMax := pixel.V(win.Bounds().Max.X - (win.Bounds().Max.X / 3),
				win.Bounds().Max.Y - (win.Bounds().Max.Y / 3))
			msg.Draw(msgDrawMin, msgDrawMax, win)
		}
	}
}

// Update updates current menu screen.
func (mm *MainMenu) Update(win *pixelgl.Window) {
	mm.menu.Update(win)
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
