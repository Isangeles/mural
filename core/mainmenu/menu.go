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

	"golang.org/x/image/colornames"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/core/data/text/lang"
	
	//"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/mtk"
)

// Menu struct represents main menu screen
// with buttons to other menus.
type Menu struct {
	mainmenu  *MainMenu
	title     *text.Text
	newcharB  *mtk.Button
	settingsB *mtk.Button
	exitB     *mtk.Button
	opened    bool
	exitReq   bool
}

// newMenu creates new menu.
func newMenu(mainmenu *MainMenu) (*Menu, error) {
	m := new(Menu)
	m.mainmenu = mainmenu
	// Title.
	font := mtk.MainFont(mtk.SIZE_BIG)
	atlas := mtk.Atlas(&font)
	m.title = text.New(pixel.V(0, 0), atlas)
	fmt.Fprint(m.title, flame.Mod().Name())
	// Buttons.
	m.newcharB = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		colornames.Red, lang.Text("gui", "newchar_b_label"),
		lang.Text("gui", "newchar_b_info"))
	m.newcharB.SetOnClickFunc(m.onNewCharButtonClicked)
	m.settingsB = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		colornames.Red, lang.Text("gui", "settings_b_label"),
		lang.Text("gui", "settings_b_info"))
	m.settingsB.SetOnClickFunc(m.onSettingsButtonClicked)
	//buttonExitBG, err := data.Picture("buttonS.png")
	//if err != nil {
	//	return nil, err
	//}
	//m.exitB = mtk.NewButtonSprite(buttonExitBG, lang.Text("gui", "exit_b_label"))
	m.exitB = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		colornames.Red, lang.Text("gui", "exit_b_label"),
		lang.Text("gui", "exit_b_info"))
	m.exitB.SetOnClickFunc(m.onExitButtonClicked)

	return m, nil
}

// Draw draws all menu elements in specified
// window.
func (m *Menu) Draw(win *pixelgl.Window) {
	// Title.
	titlePos := pixel.V(win.Bounds().Center().X,
		win.Bounds().Max.Y - m.title.Bounds().Size().Y)
	m.title.Draw(win, pixel.IM.Moved(titlePos))
	// Buttons.
	m.newcharB.Draw(win, pixel.IM.Moved(pixel.V(titlePos.X,
		titlePos.Y - m.exitB.Frame().Size().Y)))
	m.settingsB.Draw(win, pixel.IM.Moved(pixel.V(titlePos.X,
		m.newcharB.DrawArea().Min.Y - m.newcharB.Frame().Size().Y)))
	m.exitB.Draw(win, pixel.IM.Moved(pixel.V(titlePos.X,
		m.settingsB.DrawArea().Min.Y - m.settingsB.Frame().Size().Y)))
}

// Update updates all menu elements.
func (m *Menu) Update(win *pixelgl.Window) {
	if m.Opened() {
		m.newcharB.Update(win)
		m.settingsB.Update(win)
		m.exitB.Update(win)
	}
	if m.exitReq {
		win.SetClosed(true)
	}
}

// Opened checks whether menu is open.
func (m *Menu) Opened() bool {
	return m.opened
}

// Show toggles menu visibility.
func (m *Menu) Show(show bool) {
	m.opened = show
}

// onNewCharButtonClicked closes all currently open
// menus and opens new character creation  menu.
func (m *Menu) onNewCharButtonClicked(b *mtk.Button) {
	m.mainmenu.OpenNewCharMenu()
}

// onSettingsButtonClicked closes all currently open
// menus and opens settings menu.
func (m *Menu) onSettingsButtonClicked(b *mtk.Button) {
	m.mainmenu.OpenSettings()
}

// Triggered on exit button clicked.
func (m *Menu) onExitButtonClicked(b *mtk.Button) {
	m.exitReq = true
}
