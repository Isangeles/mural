/*
 * savemenu.go
 *
 * Copyright 2019 Dariusz Sikora <dev@isangeles.pl>
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

package hud

import (
	"fmt"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"

	flameconf "github.com/isangeles/flame/config"
	flamedata "github.com/isangeles/flame/core/data"
	"github.com/isangeles/flame/core/data/text/lang"
	
	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/exp"
	"github.com/isangeles/mural/log"
)

// Struct for HUD save game menu.
type SaveMenu struct {
	hud          *HUD
	bgSpr        *pixel.Sprite
	bgDraw       *imdraw.IMDraw
	drawArea     pixel.Rect
	titleText    *mtk.Text
	savesList    *mtk.List
	saveNameEdit *mtk.Textedit
	closeButton  *mtk.Button
	saveButton   *mtk.Button
	opened       bool
	focused      bool
}

// newSaveMenu creates new save menu for HUD.
func newSaveMenu(hud *HUD) *SaveMenu {
	sm := new(SaveMenu)
	sm.hud = hud
	// Background.
	sm.bgDraw = imdraw.New(nil)
	bg, err := data.PictureUI("menubg.png")
	if err == nil {
		sm.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	}
	// Title.
	sm.titleText = mtk.NewText(mtk.SIZE_SMALL, 0)
	sm.titleText.SetText(lang.Text("gui", "hud_save_menu_title"))
	// Saves list.
	bgSize := sm.Size()
	savesListSize := pixel.V(bgSize.X-mtk.ConvSize(50),
		bgSize.Y-mtk.ConvSize(200))
	sm.savesList = mtk.NewList(savesListSize, mtk.SIZE_MINI, main_color,
		sec_color, accent_color)
	sm.savesList.SetOnItemSelectFunc(sm.onSaveSelected)
	// Text field.
	sm.saveNameEdit = mtk.NewTextedit(mtk.SIZE_SMALL, main_color)
	saveNameSize := pixel.V(savesListSize.X, mtk.ConvSize(20))
	sm.saveNameEdit.SetSize(saveNameSize)
	// Buttons.
	sm.closeButton = mtk.NewButton(mtk.SIZE_SMALL, mtk.SHAPE_SQUARE,
		accent_color)
	closeButtonBG, err := data.PictureUI("closebutton1.png")
	if err == nil {
		closeBG := pixel.NewSprite(closeButtonBG, closeButtonBG.Bounds())
		sm.closeButton.SetBackground(closeBG)
	}
	sm.closeButton.SetOnClickFunc(sm.onCloseButtonClicked)
	sm.saveButton = mtk.NewButton(mtk.SIZE_SMALL, mtk.SHAPE_RECTANGLE,
		accent_color)
	sm.saveButton.SetLabel(lang.Text("gui", "save_b_label"))
	saveButtonBG, err := data.PictureUI("button_green.png")
	if err == nil {
		bg := pixel.NewSprite(saveButtonBG, saveButtonBG.Bounds())
		sm.saveButton.SetBackground(bg)
	}
	sm.saveButton.SetOnClickFunc(sm.onSaveButtonClicked)
	return sm
}

// Draw draws menu.
func (sm *SaveMenu) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	sm.drawArea = mtk.MatrixToDrawArea(matrix, sm.Size())
	// Background.
	if sm.bgSpr != nil {
		sm.bgSpr.Draw(win.Window, matrix)
	} else {
		mtk.DrawRectangle(win.Window, sm.DrawArea(), nil)
	}
	// Title.
	titleTextMove := pixel.V(0, sm.Size().Y/2-mtk.ConvSize(20))
	sm.titleText.Draw(win.Window, matrix.Moved(titleTextMove))
	// Saves.
	savesListMove := pixel.V(0, 0)
	sm.savesList.Draw(win, matrix.Moved(savesListMove))
	// Save name filed.
	saveNameMove := pixel.V(0, -mtk.ConvSize(150))
	sm.saveNameEdit.Draw(win, matrix.Moved(saveNameMove))
	// Buttons.
	closeButtonMove := mtk.ConvVec(pixel.V(sm.Size().X/2-20,
		sm.Size().Y/2-15))
	sm.closeButton.Draw(win.Window, matrix.Moved(closeButtonMove))
	saveButtonMove := mtk.ConvVec(pixel.V(0, -sm.Size().Y/2+30))
	sm.saveButton.Draw(win, matrix.Moved(saveButtonMove))
}

// Update updates menu.
func (sm *SaveMenu) Update(win *mtk.Window) {
	// Elements.
	sm.closeButton.Update(win)
	sm.saveButton.Update(win)
	sm.saveNameEdit.Update(win)
	sm.savesList.Update(win)
}

// Size returns menu background size.
func (sm *SaveMenu) Size() pixel.Vec {
	if sm.bgSpr == nil {
		// TODO: menu draw background size.
		return pixel.V(mtk.ConvSize(0), mtk.ConvSize(0))
	}
	return sm.bgSpr.Frame().Size()
}

// DrawArea return menu background draw area.
func (sm *SaveMenu) DrawArea() pixel.Rect {
	return sm.drawArea
}

// Opened check if menu is open.
func (sm *SaveMenu) Opened() bool {
	return sm.opened
}

// Show toggles menu visibility.
func (sm *SaveMenu) Show(show bool) {
	sm.opened = show
	sm.hud.Camera().Lock(sm.Opened())
	err := sm.loadSaves()
	if err != nil {
		log.Err.Printf("hud_save_menu:fail_to_load_saves:%e", err)
	}
}

// Focused checks if menu us focused.
func (sm *SaveMenu) Focused() bool {
	return sm.focused
}

// Focus toggles menu focus.
func (sm *SaveMenu) Focus(focus bool) {
	sm.focused = focus
}

// loadSaves updates saves list with
// current saves from saves dir.
func (sm *SaveMenu) loadSaves() error {
	// Clear list.
	sm.savesList.Clear()
	// Insert save names.
	pattern := fmt.Sprintf(".*%s", flamedata.SAVEGAME_FILE_EXT)
	saves, err := flamedata.DirFilesNames(flameconf.ModuleSavegamesPath(),
		pattern)
	if err != nil {
		return fmt.Errorf("fail_to_read_saved_games_dir:%v", err)
	}
	for _, s := range saves {
		sm.savesList.AddItem(s, s)
	}
	return nil
}

// Triggered after selecting one of save list items.
func (sm *SaveMenu) onSaveSelected(cs *mtk.CheckSlot) {
	saveName, ok := cs.Value().(string)
	if !ok {
		log.Err.Printf("hud_savegame:fail to retireve save name")
		return
	}
	sm.saveNameEdit.SetText(saveName)
}

// Triggered after close button clicked.
func (sm *SaveMenu) onCloseButtonClicked(b *mtk.Button) {
	sm.Show(false)
}

// Triggered after save button clicked.
func (sm *SaveMenu) onSaveButtonClicked(b *mtk.Button) {
	saveFileName := sm.saveNameEdit.Text()
	saveName := strings.Split(saveFileName, ".")[0]
	err := sm.save(saveName)
	if err != nil {
		log.Err.Printf("hud_savegame:fail_to_save:%v")
	}
}

// Save saves GUI and game state to
// savegames directory.
func (sm *SaveMenu) save(saveName string) error {
	// Retrieve saves path.
	savesPath := flameconf.ModuleSavegamesPath()
	// Save current game.
	err := flamedata.SaveGame(sm.hud.Game(), savesPath, saveName)
	if err != nil {
		return fmt.Errorf("fail_to_save_game:%v",
			err)
	}
	// Save GUI state.
	guisav := sm.hud.NewGUISave()
	err = exp.ExportGUISave(guisav, savesPath, saveName)
	if err != nil {
		return fmt.Errorf("fail_to_save_gui:%v", err)
	}
	return nil
}
