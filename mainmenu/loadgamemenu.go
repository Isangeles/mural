/*
 * loadgamemenu.go
 *
 * Copyright 2018-2021 Dariusz Sikora <dev@isangeles.pl>
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
	"path/filepath"
	"strings"

	"github.com/faiface/pixel"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/game"

	"github.com/isangeles/flame"
	flamedata "github.com/isangeles/flame/data"
	flameres "github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/serial"

	"github.com/isangeles/mtk"

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
	titleParams := mtk.Params{
		FontSize: mtk.SizeBig,
	}
	lgm.title = mtk.NewText(titleParams)
	lgm.title.SetText(lang.Text("loadgame_menu_title"))
	// Saves list.
	listSize := mtk.SizeBig.ListSize()
	listParams := mtk.Params{
		SizeRaw:     listSize,
		MainColor:   mainColor,
		SecColor:    secColor,
		AccentColor: accentColor,
	}
	lgm.savesList = mtk.NewList(listParams)
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		FontSize:  mtk.SizeMedium,
		Shape:     mtk.ShapeRectangle,
		MainColor: accentColor,
	}
	lgm.backButton = mtk.NewButton(buttonParams)
	lgm.backButton.SetLabel(lang.Text("back_button_label"))
	lgm.backButton.SetOnClickFunc(lgm.onBackButtonClicked)
	lgm.loadButton = mtk.NewButton(buttonParams)
	lgm.loadButton.SetLabel(lang.Text("load_button_label"))
	lgm.loadButton.SetOnClickFunc(lgm.onLoadButtonClicked)
	return lgm
}

// Draw draws all menu elements in specified window.
func (lgm *LoadGameMenu) Draw(win *mtk.Window) {
	// Title.
	titlePos := pixel.V(win.Bounds().Center().X,
		win.Bounds().H()-lgm.title.Size().Y)
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

// Show shows menu.
func (lgm *LoadGameMenu) Show() {
	lgm.opened = true
	lgm.mainmenu.userFocus.Focus(lgm.savesList)
	err := lgm.loadSaves()
	if err != nil {
		log.Err.Printf("load game menu: unable to load saves: %v", err)
	}
}

// Hide hides menu.
func (lgm *LoadGameMenu) Hide() {
	lgm.opened = false
	lgm.mainmenu.userFocus.Focus(nil)
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
	pattern := fmt.Sprintf(".*%s", flamedata.ModuleFileExt)
	saves, err := flamedata.DirFilesNames(lgm.mainmenu.mod.Conf().SavesPath(),
		pattern)
	if err != nil {
		return fmt.Errorf("unable to read saved games dir: %v", err)
	}
	for _, s := range saves {
		lgm.savesList.AddItem(s, s)
	}
	return nil
}

// loadSavedGame creates game and HUD from saved data.
func (lgm *LoadGameMenu) loadSavedGame(saveName string) error {
	// Import saved game.
	savePath := filepath.Join(lgm.mainmenu.mod.Conf().SavesPath(),
		saveName+flamedata.ModuleFileExt)
	modData, err := flamedata.ImportModuleFile(savePath)
	if err != nil {
		return fmt.Errorf("unable to import module: %v", err)
	}
	flameres.Clear()
	serial.Reset()
	flameres.TranslationBases = res.TranslationBases()
	m := flame.NewModule()
	m.Apply(modData)
	gameWrapper := game.New(m)
	// Import saved HUD state.
	guiSavePath := filepath.Join(lgm.mainmenu.mod.Conf().Path, data.SavesModulePath,
		saveName+data.SaveFileExt)
	layout, err := data.ImportGUISave(guiSavePath)
	if err != nil {
		return fmt.Errorf("unable to load GUI save: %v", err)
	}
	for _, pcd := range layout.Players {
		char := gameWrapper.Chapter().Character(pcd.Avatar.ID, pcd.Avatar.Serial)
		if char == nil {
			log.Err.Printf("Main menu: load game: unable to retrieve pc character: %s %s",
				pcd.Avatar.ID, pcd.Avatar.Serial)
			continue
		}
		av := object.NewAvatar(char, &pcd.Avatar)
		pc := game.NewPlayer(av, gameWrapper)
		gameWrapper.AddPlayer(pc)
	}
	// Enter game.
	if lgm.mainmenu.onGameCreated != nil {
		go lgm.mainmenu.onGameCreated(gameWrapper, layout)
	}
	return nil
}

// Triggered after back button clicked.
func (lgm *LoadGameMenu) onBackButtonClicked(b *mtk.Button) {
	lgm.mainmenu.OpenMenu()
}

// Triggered after load button clicked.
func (lgm *LoadGameMenu) onLoadButtonClicked(b *mtk.Button) {
	// Retrieve selected save name from list.
	if lgm.savesList.SelectedValue() == nil {
		return
	}
	selection := lgm.savesList.SelectedValue()
	fileName, ok := selection.(string)
	if !ok {
		log.Err.Printf("main menu: load game: unable to retrieve save name from list value")
		return
	}
	// Show loading screen.
	lgm.mainmenu.OpenLoadingScreen(lang.Text("loadgame_load_game_info"))
	defer lgm.mainmenu.CloseLoadingScreen()
	// Load saved game.
	saveName := strings.Replace(fileName, flamedata.ModuleFileExt, "", 1)
	err := lgm.loadSavedGame(saveName)
	if err != nil {
		log.Err.Printf("Main menu: load game: unable to load saved game: %v", err)
		lgm.mainmenu.ShowMessage(lang.Text("load_game_err"))
	}
	// Back to main menu.
	lgm.mainmenu.OpenMenu()
}
