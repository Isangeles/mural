/*
 * menu.go
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

package mainmenu

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/mural/core"
	"github.com/isangeles/mural/core/data"
)

// Menu struct represents main menu screen
// with buttons to other menus.
type Menu struct {
	title     *text.Text
	settingsB *core.Button
	exitB     *core.Button
	open      bool
	exitReq   bool
}

// newMenu returns new menu.
func newMenu() (*Menu, error) {
	m := new(Menu)
	// Title.
	font := data.MainFontBig()
	atlas := text.NewAtlas(font, text.ASCII)
	m.title = text.New(pixel.V(0, 0), atlas)
	fmt.Fprint(m.title, flame.Mod().Name())
	// Buttons.
	m.settingsB = core.NewButtonDraw(core.SIZE_SMALL,
		lang.Text("gui", "settings_b_label"))
	//buttonExitBG, err := data.Picture("buttonS.png")
	//if err != nil {
	//	return nil, err
	//}
	//m.exitB = core.NewButton(buttonExitBG, lang.Text("gui", "exit_b_label"))
	m.exitB = core.NewButtonDraw(core.SIZE_SMALL, lang.Text("gui", "exit_b_label"))
	m.exitB.OnClick(m.onExitButtonClicked)

	return m, nil
}

// Draw draws all menu elements in specified
// window.
func (m *Menu) Draw(win *pixelgl.Window) {
	titlePos := pixel.V(win.Bounds().Center().X, win.Bounds().Max.Y - m.title.Bounds().Size().Y)
	m.title.Draw(win, pixel.IM.Moved(titlePos))
	m.settingsB.Draw(win, pixel.IM.Moved(pixel.V(titlePos.X,
		titlePos.Y - m.exitB.Frame().Size().Y)))
	m.exitB.Draw(win, pixel.IM.Moved(pixel.V(titlePos.X,
		m.settingsB.DrawArea().Min.Y - m.settingsB.Frame().Size().Y)))
}

// Update updates all menu elements.
func (m *Menu) Update(win *pixelgl.Window) {
	m.settingsB.Update(win)
	m.exitB.Update(win)

	if m.exitReq {
		win.SetClosed(true)
	}
}

// Open checks if menu should be displayed.
func (m *Menu) Open() bool {
	return m.open
}

// Toggles menu visibility.
func (m *Menu) Show(show bool) {
	m.open = show
}

// Triggered on settings button clicked.
func (m *Menu) onSettingsButtonClicked(b *core.Button) {
	// TODO: settings toggle.
}

// Triggered on exit button clicked.
func (m *Menu) onExitButtonClicked(b *core.Button) {
	m.exitReq = true
}
