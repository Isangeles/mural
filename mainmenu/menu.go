/*
 * menu.go
 *
 * Copyright 2018-2020 Dariusz Sikora <dev@isangeles.pl>
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
	"github.com/faiface/pixel"

	"github.com/isangeles/flame/data/res/lang"

	"github.com/isangeles/mtk"
)

// Menu struct represents main menu screen
// with buttons to other menus.
type Menu struct {
	mainmenu    *MainMenu
	title       *mtk.Text
	loginButton *mtk.Button
	newgameB    *mtk.Button
	newcharB    *mtk.Button
	loadgameB   *mtk.Button
	settingsB   *mtk.Button
	exitB       *mtk.Button
	opened      bool
}

// newMenu creates new menu.
func newMenu(mainmenu *MainMenu) *Menu {
	m := new(Menu)
	m.mainmenu = mainmenu
	// Title.
	titleParams := mtk.Params{
		FontSize: mtk.SizeBig,
		SizeRaw:  mtk.ConvVec(pixel.V(900, 0)),
	}
	m.title = mtk.NewText(titleParams)
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		FontSize:  mtk.SizeMedium,
		Shape:     mtk.ShapeRectangle,
		MainColor: accentColor,
	}
	m.loginButton = mtk.NewButton(buttonParams)
	m.loginButton.SetLabel(lang.Text("login_b_label"))
	m.loginButton.SetInfo(lang.Text("login_b_info"))
	m.loginButton.SetOnClickFunc(m.onLoginButtonClicked)
	m.newgameB = mtk.NewButton(buttonParams)
	m.newgameB.SetLabel(lang.Text("newgame_button_label"))
	m.newgameB.SetInfo(lang.Text("newgame_button_info"))
	m.newgameB.SetOnClickFunc(m.onNewGameButtonClicked)
	m.newcharB = mtk.NewButton(buttonParams)
	m.newcharB.SetLabel(lang.Text("newchar_button_label"))
	m.newcharB.SetInfo(lang.Text("newchar_button_info"))
	m.newcharB.SetOnClickFunc(m.onNewCharButtonClicked)
	m.loadgameB = mtk.NewButton(buttonParams)
	m.loadgameB.SetLabel(lang.Text("loadgame_button_label"))
	m.loadgameB.SetInfo(lang.Text("loadgame_button_info"))
	m.loadgameB.SetOnClickFunc(m.onLoadGameButtonClicked)
	m.settingsB = mtk.NewButton(buttonParams)
	m.settingsB.SetLabel(lang.Text("settings_button_label"))
	m.settingsB.SetInfo(lang.Text("settings_button_info"))
	m.settingsB.SetOnClickFunc(m.onSettingsButtonClicked)
	m.exitB = mtk.NewButton(buttonParams)
	m.exitB.SetLabel(lang.Text("exit_button_label"))
	m.exitB.SetInfo(lang.Text("exit_button_info"))
	m.exitB.SetOnClickFunc(m.onExitButtonClicked)
	return m
}

// Draw draws all menu elements in specified
// window.
func (m *Menu) Draw(win *mtk.Window) {
	// Title.
	titlePos := mtk.DrawPosTC(win.Bounds(), m.title.Size())
	titlePos.Y -= mtk.ConvSize(20)
	m.title.Draw(win.Window, mtk.Matrix().Moved(titlePos))
	// Buttons.
	loginPos := mtk.BottomOf(m.title.DrawArea(), m.loginButton.Size(), 10)
	m.loginButton.Draw(win.Window, mtk.Matrix().Moved(loginPos))
	newgamePos := mtk.BottomOf(m.loginButton.DrawArea(), m.newgameB.Size(), 5)
	m.newgameB.Draw(win.Window, mtk.Matrix().Moved(newgamePos))
	newcharPos := mtk.BottomOf(m.newgameB.DrawArea(), m.newcharB.Size(), 5)
	m.newcharB.Draw(win.Window, mtk.Matrix().Moved(newcharPos))
	loadgamePos := mtk.BottomOf(m.newcharB.DrawArea(), m.loadgameB.Size(), 5)
	m.loadgameB.Draw(win.Window, mtk.Matrix().Moved(loadgamePos))
	settingsPos := mtk.BottomOf(m.loadgameB.DrawArea(), m.settingsB.Size(), 5)
	m.settingsB.Draw(win.Window, mtk.Matrix().Moved(settingsPos))
	exitPos := mtk.BottomOf(m.settingsB.DrawArea(), m.exitB.Size(), 5)
	m.exitB.Draw(win.Window, mtk.Matrix().Moved(exitPos))
}

// Update updates all menu elements.
func (m *Menu) Update(win *mtk.Window) {
	if len(m.mainmenu.playableChars) < 1 {
		m.newgameB.Active(false)
	} else {
		m.newgameB.Active(true)
	}
	if m.mainmenu.server == nil || m.mainmenu.server.Authorized() {
		m.loginButton.Active(false)
		m.newgameB.Active(true)
		m.newcharB.Active(true)
		m.loadgameB.Active(true)
	} else {
		m.loginButton.Active(true)
		m.newgameB.Active(false)
		m.newcharB.Active(false)
		m.loadgameB.Active(false)
	}
	m.loginButton.Update(win)
	m.newgameB.Update(win)
	m.newcharB.Update(win)
	m.loadgameB.Update(win)
	m.settingsB.Update(win)
	m.exitB.Update(win)
}

// Opened checks whether menu is open.
func (m *Menu) Opened() bool {
	return m.opened
}

// Show shows menu.
func (m *Menu) Show() {
	m.opened = true
}

// Hide hides menu.
func (m *Menu) Hide() {
	m.opened = false
}

// Triggered after login button clicked.
func (m *Menu) onLoginButtonClicked(b *mtk.Button) {
	m.mainmenu.HideMenus()
	m.mainmenu.loginmenu.Show()
}

// Triggered after new game button clicked.
func (m *Menu) onNewGameButtonClicked(b *mtk.Button) {
	m.mainmenu.OpenNewGameMenu()
}

// onNewCharButtonClicked closes all currently open
// menus and opens new character creation  menu.
func (m *Menu) onNewCharButtonClicked(b *mtk.Button) {
	m.mainmenu.OpenNewCharMenu()
}

// Triggered after load game button clicked.
func (m *Menu) onLoadGameButtonClicked(b *mtk.Button) {
	m.mainmenu.OpenLoadGameMenu()
}

// onSettingsButtonClicked closes all currently open
// menus and opens settings menu.
func (m *Menu) onSettingsButtonClicked(b *mtk.Button) {
	m.mainmenu.OpenSettings()
}

// Triggered on exit button clicked.
func (m *Menu) onExitButtonClicked(b *mtk.Button) {
	m.mainmenu.Exit()
}
