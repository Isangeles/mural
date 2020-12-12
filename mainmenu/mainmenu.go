/*
 * mainmenu.go
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

// mainmenu package contains main menu and also settings,
// load/save and new game menus.
package mainmenu

import (
	"fmt"
	"image/color"
	"path/filepath"

	"golang.org/x/image/colornames"

	flameres "github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/module"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/game"
)

var (
	// Main menu elements colors.
	mainColor   color.Color = colornames.Grey
	secColor    color.Color = colornames.Blue
	accentColor color.Color = colornames.Red
)

// MainMenu struct reperesents container with
// all menu screens(settings menu, new game menu, etc.).
// Wraps all main menu screens.
type MainMenu struct {
	menu          *Menu
	loginmenu     *LoginMenu
	newgamemenu   *NewGameMenu
	newcharmenu   *NewCharacterMenu
	loadgamemenu  *LoadGameMenu
	settings      *Settings
	console       *Console
	loadscreen    *LoadingScreen
	userFocus     *mtk.Focus
	msgs          *mtk.MessagesQueue
	server        *game.Server
	mod           *module.Module
	playableChars []PlayableCharData
	onGameCreated func(g *game.Game)
	onSaveLoad    func(savename string)
	loading       bool
	exiting       bool
}

// Struct with character and avatar data
// for playable characters.
type PlayableCharData struct {
	CharData   *flameres.CharacterData
	AvatarData *res.AvatarData
}

// New creates new main menu
func New() *MainMenu {
	mm := new(MainMenu)
	// Menus.
	mm.menu = newMenu(mm)
	mm.loginmenu = newLoginMenu(mm)
	mm.newgamemenu = newNewGameMenu(mm)
	mm.newcharmenu = newNewCharacterMenu(mm)
	mm.loadgamemenu = newLoadGameMenu(mm)
	mm.settings = newSettings(mm)
	// Console.
	mm.console = newConsole()
	// Loading screen.
	mm.loadscreen = newLoadingScreen(mm)
	// Messages & focus.
	mm.userFocus = new(mtk.Focus)
	mm.msgs = mtk.NewMessagesQueue(mm.userFocus)
	mm.OpenMenu()
	return mm
}

// Draw draws current menu screen.
func (mm *MainMenu) Draw(win *mtk.Window) {
	if mm.loading {
		mm.loadscreen.Draw(win)
		return
	}
	// Menus
	if mm.menu.Opened() {
		mm.menu.Draw(win)
	}
	if mm.loginmenu.Opened() {
		mm.loginmenu.Draw(win)
	}
	if mm.newgamemenu.Opened() {
		mm.newgamemenu.Draw(win)
	}
	if mm.newcharmenu.Opened() {
		mm.newcharmenu.Draw(win)
	}
	if mm.loadgamemenu.Opened() {
		mm.loadgamemenu.Draw(win)
	}
	if mm.settings.Opened() {
		mm.settings.Draw(win.Window)
	}
	// Messages.
	mm.msgs.Draw(win.Window, mtk.Matrix().Moved(win.Bounds().Center()))
	// Console.
	if mm.console.Opened() {
		mm.console.Draw(win)
	}
}

// Update updates current menu screen.
func (mm *MainMenu) Update(win *mtk.Window) {
	if mm.exiting {
		win.SetClosed(true)
		return
	}
	if mm.loading {
		mm.loadscreen.Update(win)
	}
	if mm.menu.Opened() {
		mm.menu.Update(win)
	}
	if mm.loginmenu.Opened() {
		mm.loginmenu.Update(win)
	}
	if mm.newgamemenu.Opened() {
		mm.newgamemenu.Update(win)
	}
	if mm.newcharmenu.Opened() {
		mm.newcharmenu.Update(win)
	}
	if mm.loadgamemenu.Opened() {
		mm.loadgamemenu.Update(win)
	}
	if mm.settings.Opened() {
		mm.settings.Update(win)
	}
	mm.console.Update(win)
	mm.msgs.Update(win)
}

// SetMod sets module for main menu.
func (mm *MainMenu) SetModule(mod *module.Module) {
	mm.mod = mod
	mm.menu.title.SetText(mod.Conf().ID)
}

// SetServer sets game server for main menu.
func (mm *MainMenu) SetServer(server *game.Server) {
	mm.server = server
}

// Exit sends exit request to main menu.
func (mm *MainMenu) Exit() {
	mm.exiting = true
}

// SetOnGameCreatedFunc sets specified function as function
// triggered after new game created.
func (mm *MainMenu) SetOnGameCreatedFunc(f func(g *game.Game)) {
	mm.onGameCreated = f
}

// SetOnSaveImportedFunc sets specified function as function
// triggered after save game imported.
func (mm *MainMenu) SetOnSaveLoadFunc(f func(savename string)) {
	mm.onSaveLoad = f
}

// OpenMenu opens menu.
func (mm *MainMenu) OpenMenu() {
	mm.HideMenus()
	mm.menu.Show()
}

// OpenNewGameMenu opens new game creation menu.
func (mm *MainMenu) OpenNewGameMenu() {
	mm.HideMenus()
	mm.newgamemenu.Show()
}

// OpenNewCharMenu opens new character creation menu.
func (mm *MainMenu) OpenNewCharMenu() {
	mm.HideMenus()
	mm.newcharmenu.Show()
}

// OpenLoadGameMenu opens load game menu.
func (mm *MainMenu) OpenLoadGameMenu() {
	mm.HideMenus()
	mm.loadgamemenu.Show()
}

// OpenSettings opens settings menu.
func (mm *MainMenu) OpenSettings() {
	mm.HideMenus()
	mm.settings.Show()
}

// OpenLoadingScreen opens loading screen
// with specified loading information.
func (mm *MainMenu) OpenLoadingScreen(loadInfo string) {
	mm.loadscreen.SetLoadInfo(loadInfo)
	mm.loading = true
}

// CloseLoadingScreen closes loading screen.
func (mm *MainMenu) CloseLoadingScreen() {
	mm.loading = false
}

// HideMenus hides all menus.
func (mm *MainMenu) HideMenus() {
	mm.menu.Hide()
	mm.loginmenu.Hide()
	mm.newgamemenu.Hide()
	mm.newcharmenu.Hide()
	mm.loadgamemenu.Hide()
	mm.settings.Hide()
}

// ShowMessageWindow adds specified message to messages queue
// and turns message visible(if not visible already).
func (mm *MainMenu) ShowMessageWindow(m *mtk.MessageWindow) {
	m.Show(true)
	mm.msgs.Append(m)
}

// ShowMessage creates new message window with specified message
// and adds it to messages queue.
func (mm *MainMenu) ShowMessage(msg string) {
	params := mtk.Params{
		Size:      mtk.SizeBig,
		FontSize:  mtk.SizeMedium,
		MainColor: mainColor,
		SecColor:  accentColor,
		Info:      msg,
	}
	mw := mtk.NewMessageWindow(params)
	mw.SetAcceptLabel(lang.Text("accept_button_label"))
	mw.Show(true)
	mm.msgs.Append(mw)
}

// Console returns main menu console.
func (mm *MainMenu) Console() *Console {
	return mm.console
}

// PlayableChars returns all playable characters.
func (mm *MainMenu) PlayableChars() []PlayableCharData {
	return mm.playableChars
}

// AddPlaybaleChar adds new playable character to playable
// characters list.
func (mm *MainMenu) AddPlayableChar(c PlayableCharData) {
	mm.playableChars = append(mm.playableChars, c)
}

// ImportPlayableChars import all characters from current module.
func (mm *MainMenu) ImportPlayableChars() error {
	avatarsPath := filepath.Join("data/modules", mm.mod.Conf().ID, data.GUIModulePath, "avatars")
	avsData, err := data.ImportAvatarsDir(avatarsPath)
	if err != nil {
		return fmt.Errorf("unable to import avatars: %v", err)
	}
	for _, avData := range avsData {
		for _, charData := range mm.mod.Res.Characters {
			if avData.ID != charData.ID {
				continue
			}
			pc := PlayableCharData{&charData, &avData}
			mm.playableChars = append(mm.playableChars, pc)
			// Add translation for character name.
			nameTrans := flameres.TranslationData{charData.ID, []string{avData.Name}}
			lang.AddTranslation(nameTrans)
		}
	}
	mm.newgamemenu.SetCharacters(mm.playableChars)
	return nil
}
