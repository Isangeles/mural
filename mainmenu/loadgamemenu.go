/*
 * loadgamemenu.go
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
	"fmt"
	"strings"
	
	"github.com/faiface/pixel"
	
	"github.com/isangeles/flame"
	flameconf "github.com/isangeles/flame/config"
	flamedata "github.com/isangeles/flame/core/data"
	"github.com/isangeles/flame/core/data/text/lang"

	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/log"
)

// LoadGameMenu struct represents load game
// menu.
type LoadGameMenu struct {
	mainmenu   *MainMenu
	title      *mtk.Text
	savesList  *mtk.List
	backButton *mtk.Button
	loadButton *mtk.Button
	opened     bool
}

// newLoadGameMenu creates load game menu.
func newLoadGameMenu(mainmenu *MainMenu) *LoadGameMenu {
	lgm := new(LoadGameMenu)
	lgm.mainmenu = mainmenu
	// Title.
	lgm.title = mtk.NewText(mtk.SIZE_BIG, 0)
	lgm.title.SetText(lang.Text("gui", "loadgame_menu_title"))
	// Saves list.
	listSize := mtk.SIZE_BIG.ListSize().Size()
	lgm.savesList = mtk.NewList(listSize, mtk.SIZE_BIG, main_color,
		sec_color, accent_color)
	// Buttons.
	lgm.backButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		accent_color)
	lgm.backButton.SetLabel(lang.Text("gui", "back_b_label"))
	lgm.backButton.SetOnClickFunc(lgm.onBackButtonClicked)
	lgm.loadButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		accent_color)
	lgm.loadButton.SetLabel(lang.Text("gui", "load_b_label"))
	lgm.loadButton.SetOnClickFunc(lgm.onLoadButtonClicked)
	return lgm
}

// Draw draws all menu elements in specified window.
func (lgm *LoadGameMenu) Draw(win *mtk.Window) {
	// Title.
	titlePos := pixel.V(win.Bounds().Center().X,
		win.Bounds().H() - lgm.title.Size().Y)
	lgm.title.Draw(win.Window, mtk.Matrix().Moved(titlePos))
	// Saves list.
	savesListPos := win.Bounds().Center()
	lgm.savesList.Draw(win.Window, mtk.Matrix().Moved(savesListPos))
	// Buttons.
	backButtonPos := mtk.DrawPosBL(win.Bounds(), lgm.backButton.Size())
	loadButtonPos := mtk.DrawPosBR(win.Bounds(), lgm.loadButton.Size())
	lgm.backButton.Draw(win.Window, mtk.Matrix().Moved(backButtonPos))
	lgm.loadButton.Draw(win.Window, mtk.Matrix().Moved(loadButtonPos))
}

// Update updates all menu elements.
func (lgm *LoadGameMenu) Update(win *mtk.Window) {
	lgm.backButton.Update(win)
	lgm.loadButton.Update(win)
	lgm.savesList.Update(win)
}

// Show toggles menu visibility.
func (lgm *LoadGameMenu) Show(show bool) {
	lgm.opened = show
	if lgm.Opened() {
		lgm.mainmenu.userFocus.Focus(lgm.savesList)
		err := lgm.loadSaves()
		if err != nil {
			log.Err.Printf("load_game_menu:fail_to_load_saves:%v", err)
		}
	} else {
		lgm.mainmenu.userFocus.Focus(nil)
	}
}

// Opened checks whether menu is open.
func (lgm *LoadGameMenu) Opened() bool {
	return lgm.opened
}

// loadSaves updates saves list with currrent
// saves from saves dir.
func (lgm *LoadGameMenu) loadSaves() error {
	// Clear list.
	lgm.savesList.Clear()
	// Insert save names.
	pattern := fmt.Sprintf(".*%s", flamedata.SAVEGAME_FILE_EXT)
	saves, err := flamedata.DirFilesNames(flameconf.ModuleSavegamesPath(),
		pattern)
	if err != nil {
		return fmt.Errorf("fail_to_read_saved_games_dir:%v", err)
	}
	for _, s := range saves {
		lgm.savesList.AddItem(s, s)
	}
	return nil
}

// importSave imports saved game from file with
// specified name.
func (lgm *LoadGameMenu) loadSave(savName string) {
	// Import saved game from file.
	lgm.mainmenu.OpenLoadingScreen(lang.Text("gui", "loadgame_import_save_info"))
	defer lgm.mainmenu.CloseLoadingScreen()
	// Load game.
	g, err := flame.LoadGame(savName)
	if err != nil {
		log.Err.Printf("load_game_menu:fail_to_load_game_save:%v", err)
		lgm.mainmenu.ShowMessage(lang.Text("gui", "load_game_err"))
		return
	}
	// Pass imported save.
	if lgm.mainmenu.onSaveLoaded == nil {
		return
	}
	lgm.mainmenu.onSaveLoaded(g, savName)
}

// Triggered after back button clicked.
func (lgm *LoadGameMenu) onBackButtonClicked(b *mtk.Button) {
	lgm.mainmenu.OpenMenu()
}

// Triggered after load button clicked.
func (lgm *LoadGameMenu) onLoadButtonClicked(b *mtk.Button) {
	if lgm.savesList.SelectedValue() == nil {
		return
	}
	selection := lgm.savesList.SelectedValue()
	filename, ok := selection.(string)
	if !ok {
		log.Err.Printf("main_menu:load_game:fail to retrieve save name from list value")
		return
	}
	savename := strings.Replace(filename, ".savegame", "", 1)
	go lgm.loadSave(savename)
	lgm.mainmenu.OpenMenu()
}
