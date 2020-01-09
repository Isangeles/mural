/*
 * loadgamemenu.go
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
	"fmt"
	"strings"

	"github.com/faiface/pixel"

	flameconf "github.com/isangeles/flame/config"
	flamedata "github.com/isangeles/flame/core/data"
	"github.com/isangeles/flame/core/data/res/lang"

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
	lgm.backButton.SetLabel(lang.Text("back_b_label"))
	lgm.backButton.SetOnClickFunc(lgm.onBackButtonClicked)
	lgm.loadButton = mtk.NewButton(buttonParams)
	lgm.loadButton.SetLabel(lang.Text("load_b_label"))
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

// Show toggles menu visibility.
func (lgm *LoadGameMenu) Show(show bool) {
	lgm.opened = show
	if lgm.Opened() {
		lgm.mainmenu.userFocus.Focus(lgm.savesList)
		err := lgm.loadSaves()
		if err != nil {
			log.Err.Printf("load game menu: fail to load saves: %v", err)
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
	pattern := fmt.Sprintf(".*%s", flamedata.SavegameFileExt)
	saves, err := flamedata.DirFilesNames(flameconf.ModuleSavegamesPath(),
		pattern)
	if err != nil {
		return fmt.Errorf("fail to read saved games dir: %v", err)
	}
	for _, s := range saves {
		lgm.savesList.AddItem(s, s)
	}
	return nil
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
		log.Err.Printf("main menu: load game: fail to retrieve save name from list value")
		return
	}
	savename := strings.Replace(filename, ".savegame", "", 1)
	if lgm.mainmenu.onSaveLoad != nil {
		go lgm.mainmenu.onSaveLoad(savename)
	}
	lgm.mainmenu.OpenMenu()
}
