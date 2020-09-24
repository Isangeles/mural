/*
 * loginmenu.go
 *
 * Copyright 2020 Dariusz Sikora <dev@isangeles.pl>
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
	"github.com/isangeles/flame/data/res/lang"

	"github.com/isangeles/fire/response"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/log"
)

// Struct for login menu.
type LoginMenu struct {
	mainmenu *MainMenu
	opened   bool
}

// newLoginMenu creates new login menu.
func newLoginMenu(mainmenu *MainMenu) *LoginMenu {
	lm := new(LoginMenu)
	lm.mainmenu = mainmenu
	return lm
}

// Draw draws login menu in specified window.
func (lm *LoginMenu) Draw(win *mtk.Window) {
}

// Update updates all menu elements.
func (lm *LoginMenu) Update(win *mtk.Window) {
	if lm.mainmenu.server != nil && lm.mainmenu.server.Authorized() {
		lm.mainmenu.server.SetOnResponseFunc(nil)
		lm.mainmenu.ShowMessage(lang.Text("login_logged_in_msg"))
		lm.mainmenu.OpenMenu()
	}
}

// Show shows menu.
func (lm *LoginMenu) Show() {
	lm.opened = true
	if lm.mainmenu.server == nil {
		return
	}
	lm.mainmenu.server.SetOnResponseFunc(lm.handleResponse)
	// Auto-login.
	if len(config.ServerLogin) > 0 && len(config.ServerPassword) > 0 {
		err := lm.mainmenu.server.Login(config.ServerLogin, config.ServerPassword)
		if err != nil {
			log.Err.Printf("Login menu: unable to send login request: %v", err)
		}
	}
}

// Hide hides menu.
func (lm *LoginMenu) Hide() {
	lm.opened = false
}

// Opened checks if menu is open.
func (lm *LoginMenu) Opened() bool {
	return lm.opened
}

// handleResponse handles specified response from Fire server.
func (lm *LoginMenu) handleResponse(resp response.Response) {
	for _, r := range resp.Error {
		log.Err.Printf("Login menu: server error: %v", r)
	}
}
