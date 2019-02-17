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
	"io/ioutil"
	"strings"
	
	"github.com/faiface/pixel"
	
	"github.com/isangeles/flame"
	flameconf "github.com/isangeles/flame/config"
	flamedata "github.com/isangeles/flame/core/data"
	"github.com/isangeles/flame/core/data/text/lang"

	"github.com/isangeles/mural/core/data/imp"
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
func newLoadGameMenu(mainmenu *MainMenu) (*LoadGameMenu, error) {
	lgm := new(LoadGameMenu)
	lgm.mainmenu = mainmenu
	// Title.
	lgm.title = mtk.NewText(mtk.SIZE_BIG, 0)
	lgm.title.SetText(lang.Text("gui", "loadgame_menu_title"))
	// Saves list.
	lgm.savesList = mtk.NewList(mtk.SIZE_BIG, main_color, sec_color,
		accent_color)
	gameSaves, err := saveGamesFiles(flameconf.ModuleSavegamesPath())
	if err != nil {
		return nil, fmt.Errorf("fail_to_read_saved_games_dir:%v",
			err)
	}
	for _, sav := range gameSaves {
		lgm.savesList.AddItem(sav, sav)
	}
	// Buttons.
	lgm.backButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		accent_color, lang.Text("gui", "back_b_label"), "")
	lgm.backButton.SetOnClickFunc(lgm.onBackButtonClicked)
	lgm.loadButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		accent_color, lang.Text("gui", "load_b_label"), "")
	lgm.loadButton.SetOnClickFunc(lgm.onLoadButtonClicked)
	return lgm, nil
}

// Draw draws all menu elements in specified window.
func (lgm *LoadGameMenu) Draw(win *mtk.Window) {
	// Title.
	titlePos := pixel.V(win.Bounds().Center().X,
		win.Bounds().Max.Y - lgm.title.Bounds().Size().Y)
	lgm.title.Draw(win.Window, mtk.Matrix().Moved(titlePos))
	// Saves list.
	lgm.savesList.Draw(win.Window, mtk.Matrix().Moved(mtk.BottomOf(
		lgm.title.DrawArea(), lgm.savesList.Bounds(), 10)))
	// Buttons.
	backButtonPos := mtk.DrawPosBL(win.Bounds(), lgm.backButton.Frame())
	loadButtonPos := mtk.DrawPosBR(win.Bounds(), lgm.loadButton.Frame())
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
	if show {
		lgm.mainmenu.userFocus.Focus(lgm.savesList)
	} else {
		lgm.mainmenu.userFocus.Focus(nil)
	}
}

// Opened checks whether menu is open.
func (lgm *LoadGameMenu) Opened() bool {
	return lgm.opened
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
	savName, ok := selection.(string)
	if !ok {
		log.Err.Printf("load_game_menu:fail to retrieve save name from list value")
		return
	}
	lgm.mainmenu.OpenLoadingScreen(lang.Text("gui", "loadgame_load_info"))
	go lgm.loadGame(savName)
}

// saveGamesFiles returns names of all save files
// in directory with specified path.
func saveGamesFiles(dirPath string) ([]string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_read_dir:%v",
			err)
	}
	filesNames := make([]string, 0)
	for _, fInfo := range files {
		if strings.HasSuffix(fInfo.Name(), flamedata.SAVEGAME_FILE_EXT) {
			filesNames = append(filesNames, fInfo.Name())
		}
	}
	return filesNames, nil
}

// loadGame loads saved game file with
// specified name.
func (lgm *LoadGameMenu) loadGame(savName string) {
	err := imp.LoadModuleResources(flame.Mod())
	if err != nil {
		log.Err.Printf("fail_to_load_resources:%v", err)
		return
	}
	sav, err := flamedata.ImportSavedGame(flame.Mod(), flameconf.ModuleSavegamesPath(),
		savName)
	if err != nil {
		log.Err.Printf("load_game_menu:fail_to_load_saved_game:%v", err)
		return
	}
	lgm.mainmenu.CloseLoadingScreen()
	lgm.mainmenu.OnGameLoaded(sav)
}
