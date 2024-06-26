/*
 * loginmenu.go
 *
 * Copyright 2020-2024 Dariusz Sikora <ds@isangeles.dev>
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
	"github.com/gopxl/pixel"

	"github.com/isangeles/flame/data/res/lang"

	"github.com/isangeles/fire/request"

	"github.com/isangeles/mtk"
)

// Struct for login menu.
type LoginMenu struct {
	mainmenu    *MainMenu
	title       *mtk.Text
	loginLabel  *mtk.Text
	passLabel   *mtk.Text
	loginEdit   *mtk.Textedit
	passEdit    *mtk.Textedit
	loginButton *mtk.Button
	backButton  *mtk.Button
	opened      bool
}

// newLoginMenu creates new login menu.
func newLoginMenu(mainmenu *MainMenu) *LoginMenu {
	lm := new(LoginMenu)
	lm.mainmenu = mainmenu
	// Labels.
	labelParams := mtk.Params{
		FontSize: mtk.SizeMedium,
	}
	lm.title = mtk.NewText(labelParams)
	lm.title.SetText(lang.Text("login_menu_title"))
	lm.loginLabel = mtk.NewText(labelParams)
	lm.loginLabel.SetText(lang.Text("login_menu_login_label"))
	lm.passLabel = mtk.NewText(labelParams)
	lm.passLabel.SetText(lang.Text("login_menu_pass_label"))
	// Text edit fields.
	texteditParams := mtk.Params{
		FontSize:  mtk.SizeMedium,
		MainColor: mainColor,
	}
	lm.loginEdit = mtk.NewTextedit(texteditParams)
	lm.passEdit = mtk.NewTextedit(texteditParams)
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		FontSize:  mtk.SizeMedium,
		Shape:     mtk.ShapeRectangle,
		MainColor: accentColor,
	}
	lm.loginButton = mtk.NewButton(buttonParams)
	lm.loginButton.SetLabel(lang.Text("login_menu_login_button_label"))
	lm.loginButton.SetOnClickFunc(lm.onLoginButtonClicked)
	lm.backButton = mtk.NewButton(buttonParams)
	lm.backButton.SetLabel(lang.Text("back_button_label"))
	lm.backButton.SetOnClickFunc(lm.onBackButtonClicked)
	return lm
}

// Draw draws login menu in specified window.
func (lm *LoginMenu) Draw(win *mtk.Window) {
	// Title.
	titlePos := pixel.V(win.Bounds().Center().X,
		win.Bounds().H()-lm.title.Size().Y)
	lm.title.Draw(win, mtk.Matrix().Moved(titlePos))
	// Login field.
	loginLabelPos := mtk.BottomOf(lm.title.DrawArea(), lm.loginEdit.Size(), 100)
	lm.loginLabel.Draw(win, mtk.Matrix().Moved(loginLabelPos))
	loginEditPos := mtk.BottomOf(lm.loginLabel.DrawArea(), lm.loginEdit.Size(), 10)
	loginEditSize := lm.title.DrawArea().Size()
	lm.loginEdit.SetSize(loginEditSize)
	lm.loginEdit.Draw(win.Window, mtk.Matrix().Moved(loginEditPos))
	// Password field.
	passLabelPos := mtk.BottomOf(lm.loginEdit.DrawArea(), lm.passEdit.Size(), 10)
	lm.passLabel.Draw(win, mtk.Matrix().Moved(passLabelPos))
	passEditPos := mtk.BottomOf(lm.passLabel.DrawArea(), lm.passEdit.Size(), 10)
	passEditSize := lm.title.DrawArea().Size()
	lm.passEdit.SetSize(passEditSize)
	lm.passEdit.Draw(win.Window, mtk.Matrix().Moved(passEditPos))
	// Buttons.
	loginButtonPos := mtk.BottomOf(lm.passEdit.DrawArea(), lm.loginButton.Size(), 10)
	lm.loginButton.Draw(win.Window, mtk.Matrix().Moved(loginButtonPos))
	backButtonPos := mtk.DrawPosBL(win.Bounds(), lm.backButton.Size())
	lm.backButton.Draw(win.Window, mtk.Matrix().Moved(backButtonPos))
}

// Update updates all menu elements.
func (lm *LoginMenu) Update(win *mtk.Window) {
	if lm.mainmenu.server != nil && lm.mainmenu.server.Authorized() {
		lm.mainmenu.ShowMessage(lang.Text("login_logged_in_msg"))
		lm.mainmenu.OpenMenu()
		return
	}
	lm.loginEdit.Update(win)
	lm.passEdit.Update(win)
	lm.loginButton.Update(win)
	lm.backButton.Update(win)
}

// Show shows menu.
func (lm *LoginMenu) Show() {
	lm.opened = true
}

// Hide hides menu.
func (lm *LoginMenu) Hide() {
	lm.opened = false
}

// Opened checks if menu is open.
func (lm *LoginMenu) Opened() bool {
	return lm.opened
}

// Triggered on login button click.
func (lm *LoginMenu) onLoginButtonClicked(b *mtk.Button) {
	loginReq := request.Login{lm.loginEdit.Text(), lm.passEdit.Text()}
	req := request.Request{Login: []request.Login{loginReq}}
	err := lm.mainmenu.server.Send(req)
	if err != nil {
		lm.mainmenu.ShowMessage(lang.Text("login_menu_create_login_req_err"))
	}
}

// Triggered on back button click.
func (lm *LoginMenu) onBackButtonClicked(b *mtk.Button) {
	lm.mainmenu.OpenMenu()
}
