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
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame/core/data/text/lang"
	
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/core/mtk"
)

// MainMenu struct reperesents container with
// all menu screens(settings menu, new game menu, etc.).
// Handles switching betwen menus.
type MainMenu struct {
	menu      *Menu
	settings  *Settings
	console   *Console
	msgs      []*mtk.MessageWindow
	userFocus *mtk.Focus
}

// New returns new main menu
func New() (*MainMenu, error) {
	mm := new(MainMenu)
	// Menu.
	m, err := newMenu()
	if err != nil {
		return nil, err
	}
	m.SetOnSettingsButtonClickedFunc(mm.onSettingsButtonClicked)
	mm.menu = m
	// Settings.
	s, err := newSettings()
	if err != nil {
		return nil, err
	}
	s.SetOnBackButtonClickedFunc(mm.onCloseSettingsButtonClicked)
	mm.settings = s
	// Console.
	c, err := newConsole()
	if err != nil {
		return nil, err
	}
	mm.console = c
	// Messages & focus test.
	for i := 0; i < 2; i++ {
		msg, err := mtk.NewMessageWindow(mtk.SIZE_SMALL,
			fmt.Sprintf("This is test UI message.\n%d", i))
		if err != nil {
			return nil, err
		}
		msg.Show(true)
		mm.msgs = append(mm.msgs, msg)
	}

	mm.userFocus = new(mtk.Focus)
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
		if msg.Opened() {
			mm.userFocus.Focus(msg)
			msg.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
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
		if msg.Opened() {
			msg.Update(win)
		}
		if msg.Dismissed() {
			mm.msgs = append(mm.msgs[:i], mm.msgs[i+1:]...) // remove dismissed message
		}
	}
}

// OpenMenu opens menu.
func (mm *MainMenu) OpenMenu() {
	mm.HideMenus()
	mm.menu.Show(true)
}

// OpenSettings opens settings menu.
func (mm *MainMenu) OpenSettings() {
	mm.HideMenus()
	mm.settings.Show(true) 
}

// HideMenus hides all menus.
func (mm *MainMenu) HideMenus() {
	mm.menu.Show(false)
	mm.settings.Show(false)
}

// CloseSettings closes settings menu. Also displays message
// about required game restart if settings was changed.
func (mm *MainMenu) CloseSettings() {
	if mm.settings.Changed() {
		msg, err := mtk.NewMessageWindow(mtk.SIZE_SMALL,
			lang.Text("gui", "settings_reset_msg"))
		if err != nil {
			log.Err.Printf("mainmenu:fail_to_create_settings_change_message")
			mm.OpenMenu()
			return
		}
		msg.Show(true)
		mm.msgs = append(mm.msgs, msg)
		mm.settings.Apply()
	}
	mm.OpenMenu()
}

// CloseSettingsWithDialog creates settings apply dialog and puts it on
// main menu messages list.
func (mm *MainMenu) CloseSettingsWithDialog() {
	if mm.settings.Changed() {
		dlg, err := mtk.NewDialogWindow(mtk.SIZE_SMALL,
			lang.Text("gui", "settings_save_msg"))
		if err != nil {
			log.Err.Printf("mainmenu:fail_to_create_settings_confirm_dialog")
			mm.CloseSettings() 
			return
		}
		dlg.SetOnAcceptFunc(mm.onSettingsApplyAccept)
		dlg.SetOnCancelFunc(mm.onSettingsApplyCancel)
		dlg.Show(true)
		mm.msgs = append(mm.msgs, dlg)
	} else {
		mm.CloseSettings()
	}	
}

// closeSettingsConfirm displays dialog window with settings save
// confirmation.
func (mm *MainMenu) onCloseSettingsButtonClicked(b *mtk.Button) {
	mm.CloseSettingsWithDialog()
}

// onSettingsApplyAccept closes and applies settings. Triggered after
// accepting settings confirm dialog.
func (mm *MainMenu) onSettingsApplyAccept(m *mtk.MessageWindow) {
	mm.CloseSettings()
}

// onSettingsApplyCancel displays confirm dialog and closes settings
// without saving. Triggered after rejecting settings confirm dialog.
func (mm *MainMenu) onSettingsApplyCancel(m *mtk.MessageWindow) {
	mm.OpenMenu()
}

// onMenuButtonClicked closes all currently open
// menus and opens main menu.
func (mm *MainMenu) onMenuButtonClicked(b *mtk.Button) {
	mm.OpenMenu()
}

// onSettingsButtonClicked closes all currently open
// menus and opens settings menu.
func (mm *MainMenu) onSettingsButtonClicked(b *mtk.Button) {
	mm.OpenSettings()
}
