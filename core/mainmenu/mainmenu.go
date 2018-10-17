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
	menu        *Menu
	newcharmenu *NewCharacterMenu
	settings    *Settings
	console     *Console
	msgs        *mtk.MessagesQueue
	userFocus   *mtk.Focus
}

// New returns new main menu
func New() (*MainMenu, error) {
	mm := new(MainMenu)
	// Menu.
	m, err := newMenu()
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_main_menu:%v",
			err)
	}
	m.SetOnNewCharButtonClickedFunc(mm.onNewCharButtonClicked)
	m.SetOnSettingsButtonClickedFunc(mm.onSettingsButtonClicked)
	mm.menu = m
	// New character menu.
	ncm, err := newNewCharacterMenu()
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_new_character_menu:%v",
			err)
	}
	ncm.SetOnBackFunc(mm.onNewCharBackButtonClicked)
	mm.newcharmenu = ncm
	// Settings.
	s, err := newSettings()
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_settings_menu:%v",
			err)
	}
	s.SetOnBackFunc(mm.onCloseSettingsButtonClicked)
	mm.settings = s
	// Console.
	c, err := newConsole()
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_console:%v",
			err)
	}
	mm.console = c
	// Messages & focus test.
	mm.userFocus = new(mtk.Focus)
	mm.msgs = mtk.NewMessagesQueue(mm.userFocus)
	/*
	for i := 0; i < 2; i++ {
		msg, err := mtk.NewMessageWindow(mtk.SIZE_SMALL,
			fmt.Sprintf("This is test UI message.\n%d", i))
		if err != nil {
			return nil, err
		}
		msg.Show(true)
		mm.msgs.Append(msg)
	}
        */

	mm.menu.Show(true)
	return mm, nil
}

// Draw draws current menu screen.
func (mm *MainMenu) Draw(win *pixelgl.Window) {
	// Menu.
	if mm.menu.Opened() {
		mm.menu.Draw(win)
	}
	// New character menu.
	if mm.newcharmenu.Opened() {
		mm.newcharmenu.Draw(win)
	}
	// Settings.
	if mm.settings.Opened() {
		mm.settings.Draw(win)
	}
	// Messages.
	mm.msgs.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	// Console.
	if mm.console.Opened() {
		conBottomLeft := pixel.V(win.Bounds().Min.X, win.Bounds().Center().Y)
		mm.console.Draw(conBottomLeft, mtk.DisTR(win.Bounds(), 0), win)
	}
}

// Update updates current menu screen.
func (mm *MainMenu) Update(win *pixelgl.Window) {
	mm.menu.Update(win)
	mm.newcharmenu.Update(win)
	mm.settings.Update(win)
	mm.console.Update(win)
	mm.msgs.Update(win)
}

// OpenMenu opens menu.
func (mm *MainMenu) OpenMenu() {
	mm.HideMenus()
	mm.menu.Show(true)
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
	mm.newcharmenu.Show(false)
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
		mm.msgs.Append(msg)
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
		mm.msgs.Append(dlg)
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

// onNewCharBackButtonClicked closes all currently open
// menus and opens main menu.
func (mm *MainMenu) onNewCharBackButtonClicked(b *mtk.Button) {
	mm.OpenMenu()
}

// onNewCharButtonClicked closes all currently open
// menus and opens new character creation  menu.
func (mm *MainMenu) onNewCharButtonClicked(b *mtk.Button) {
	mm.OpenNewCharMenu()
}

// onSettingsButtonClicked closes all currently open
// menus and opens settings menu.
func (mm *MainMenu) onSettingsButtonClicked(b *mtk.Button) {
	mm.OpenSettings()
}
