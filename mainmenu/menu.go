/*
 * menu.go
 *
 * Copyright 2018-2019 Dariusz Sikora <dev@isangeles.pl>
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
	"github.com/isangeles/flame"
	"github.com/isangeles/flame/core/data/text/lang"

	"github.com/isangeles/mtk"
)

// Menu struct represents main menu screen
// with buttons to other menus.
type Menu struct {
	mainmenu  *MainMenu
	title     *mtk.Text
	newgameB  *mtk.Button
	newcharB  *mtk.Button
	loadgameB *mtk.Button
	settingsB *mtk.Button
	exitB     *mtk.Button
	opened    bool
}

// newMenu creates new menu.
func newMenu(mainmenu *MainMenu) *Menu {
	m := new(Menu)
	m.mainmenu = mainmenu
	// Title.
	m.title = mtk.NewText(mtk.SizeBig, mtk.ConvSize(900))
	m.title.SetText(flame.Mod().Conf().ID)
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		FontSize:  mtk.SizeMedium,
		Shape:     mtk.ShapeRectangle,
		MainColor: accentColor,
	}
	m.newgameB = mtk.NewButton(buttonParams)
	m.newgameB.SetLabel(lang.Text("gui", "newgame_b_label"))
	m.newgameB.SetInfo(lang.Text("gui", "newgame_b_info"))
	m.newgameB.SetOnClickFunc(m.onNewGameButtonClicked)
	m.newcharB = mtk.NewButton(buttonParams)
	m.newcharB.SetLabel(lang.Text("gui", "newchar_b_label"))
	m.newcharB.SetInfo(lang.Text("gui", "newchar_b_info"))
	m.newcharB.SetOnClickFunc(m.onNewCharButtonClicked)
	m.loadgameB = mtk.NewButton(buttonParams)
	m.loadgameB.SetLabel(lang.Text("gui", "loadgame_b_label"))
	m.loadgameB.SetInfo(lang.Text("gui", "loadgame_b_info"))
	m.loadgameB.SetOnClickFunc(m.onLoadGameButtonClicked)
	m.settingsB = mtk.NewButton(buttonParams)
	m.settingsB.SetLabel(lang.Text("gui", "settings_b_label"))
	m.settingsB.SetInfo(lang.Text("gui", "settings_b_info"))
	m.settingsB.SetOnClickFunc(m.onSettingsButtonClicked)
	m.exitB = mtk.NewButton(buttonParams)
	m.exitB.SetLabel(lang.Text("gui", "exit_b_label"))
	m.exitB.SetInfo(lang.Text("gui", "exit_b_info"))
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
	newgamePos := mtk.BottomOf(m.title.DrawArea(), m.newgameB.Size(), 10)
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

// Show toggles menu visibility.
func (m *Menu) Show(show bool) {
	m.opened = show
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
