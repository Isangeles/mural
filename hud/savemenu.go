/*
 * savemenu.go
 *
 * Copyright 2019-2024 Dariusz Sikora <ds@isangeles.dev>
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
	"os"
	"path/filepath"
	"strings"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"

	flamedata "github.com/isangeles/flame/data"
	"github.com/isangeles/flame/data/res/lang"

	"github.com/isangeles/fire/request"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/data"
	"github.com/isangeles/mural/data/res/graphic"
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
	bg := graphic.Textures["menubg.png"]
	if bg != nil {
		sm.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	}
	// Title.
	titleParams := mtk.Params{
		FontSize: mtk.SizeSmall,
	}
	sm.titleText = mtk.NewText(titleParams)
	sm.titleText.SetText(lang.Text("hud_save_menu_title"))
	// Saves list.
	bgSize := sm.Size()
	savesListSize := pixel.V(bgSize.X-mtk.ConvSize(50),
		bgSize.Y-mtk.ConvSize(200))
	savesListParams := mtk.Params{
		SizeRaw:     savesListSize,
		MainColor:   mainColor,
		SecColor:    secColor,
		AccentColor: accentColor,
	}
	sm.savesList = mtk.NewList(savesListParams)
	sm.savesList.SetOnItemSelectFunc(sm.onSaveSelected)
	upButtonTex := graphic.Textures["scrollup.png"]
	if upButtonTex != nil {
		upSprite := pixel.NewSprite(upButtonTex, upButtonTex.Bounds())
		sm.savesList.SetUpButtonBackground(upSprite)
	}
	downButtonTex := graphic.Textures["scrolldown.png"]
	if downButtonTex != nil {
		downSprite := pixel.NewSprite(downButtonTex, downButtonTex.Bounds())
		sm.savesList.SetDownButtonBackground(downSprite)
	}
	// Text field.
	texteditParams := mtk.Params{
		FontSize:  mtk.SizeSmall,
		MainColor: mainColor,
	}
	sm.saveNameEdit = mtk.NewTextedit(texteditParams)
	saveNameSize := pixel.V(savesListSize.X, mtk.ConvSize(20))
	sm.saveNameEdit.SetSize(saveNameSize)
	// Buttons.
	closeButtonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		Shape:     mtk.ShapeSquare,
		MainColor: accentColor,
	}
	saveButtonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		Shape:     mtk.ShapeRectangle,
		MainColor: accentColor,
	}
	sm.closeButton = mtk.NewButton(closeButtonParams)
	closeButtonBG := graphic.Textures["closebutton1.png"]
	if closeButtonBG != nil {
		closeBG := pixel.NewSprite(closeButtonBG, closeButtonBG.Bounds())
		sm.closeButton.SetBackground(closeBG)
	}
	sm.closeButton.SetOnClickFunc(sm.onCloseButtonClicked)
	sm.saveButton = mtk.NewButton(saveButtonParams)
	sm.saveButton.SetLabel(lang.Text("save_button_label"))
	saveButtonBG := graphic.Textures["button_green.png"]
	if saveButtonBG != nil {
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
	closeButtonMove := pixel.V(sm.Size().X/2-mtk.ConvSize(20),
		sm.Size().Y/2-mtk.ConvSize(15))
	sm.closeButton.Draw(win.Window, matrix.Moved(closeButtonMove))
	saveButtonMove := pixel.V(mtk.ConvSize(0), -sm.Size().Y/2+mtk.ConvSize(30))
	sm.saveButton.Draw(win, matrix.Moved(saveButtonMove))
}

// Update updates menu.
func (sm *SaveMenu) Update(win *mtk.Window) {
	if !sm.Opened() {
		return
	}
	// Key events.
	if win.JustPressed(exitKey) {
		sm.Hide()
	}
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
	return mtk.ConvVec(sm.bgSpr.Frame().Size())
}

// DrawArea return menu background draw area.
func (sm *SaveMenu) DrawArea() pixel.Rect {
	return sm.drawArea
}

// Opened check if menu is open.
func (sm *SaveMenu) Opened() bool {
	return sm.opened
}

// Show shows menu.
func (sm *SaveMenu) Show() {
	sm.opened = true
	sm.hud.Camera().Lock(true)
	sm.hud.bar.Lock(true)
	err := sm.loadSaves()
	if err != nil {
		log.Err.Printf("hud: savegame menu: unable to load saves: %v", err)
	}
}

// Hide hides menu.
func (sm *SaveMenu) Hide() {
	sm.opened = false
	sm.hud.Camera().Lock(false)
	sm.hud.bar.Lock(false)
}

// Focused checks if menu us focused.
func (sm *SaveMenu) Focused() bool {
	return sm.focused
}

// Focus toggles menu focus.
func (sm *SaveMenu) Focus(focus bool) {
	sm.focused = focus
}

// loadSaves updates saves list with current
// saves from the saves dir.
func (sm *SaveMenu) loadSaves() error {
	// Clear list.
	sm.savesList.Clear()
	// Check if saves dir exists.
	path := filepath.Join(config.GUIPath, data.HUDDir)
	_, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	// Retrive save names.
	pattern := fmt.Sprintf(".*%s", data.HUDFileExt)
	saves, err := data.DirFiles(path, pattern)
	if err != nil {
		return fmt.Errorf("unable to retrieve save files: %v", err)
	}
	// Add save names to the list.
	for _, s := range saves {
		if s != config.DefaultHUD {
			sm.savesList.AddItem(s, s)
		}
	}
	return nil
}

// Triggered after selecting one of save list items.
func (sm *SaveMenu) onSaveSelected(cs *mtk.CheckSlot) {
	saveName, ok := cs.Value().(string)
	if !ok {
		log.Err.Printf("hud savegame menu: unable to retireve save name")
		return
	}
	sm.saveNameEdit.SetText(saveName)
}

// Triggered after close button clicked.
func (sm *SaveMenu) onCloseButtonClicked(b *mtk.Button) {
	sm.Hide()
}

// Triggered after save button clicked.
func (sm *SaveMenu) onSaveButtonClicked(b *mtk.Button) {
	saveFileName := sm.saveNameEdit.Text()
	saveName := strings.Split(saveFileName, ".")[0]
	if len(saveName) < 1 {
		return
	}
	err := sm.save(saveName)
	if err != nil {
		log.Err.Printf("hud: savegame menu: unable to save: %v", err)
	}
}

// Save saves GUI and module state.
func (sm *SaveMenu) save(saveName string) error {
	// Save HUD.
	hudPath := filepath.Join(config.GUIPath, data.HUDDir,
		saveName+data.HUDFileExt)
	err := data.ExportHUD(sm.hud.Data(), hudPath)
	if err != nil {
		return fmt.Errorf("unable to export hud: %v", err)
	}
	// Save current game.
	if sm.hud.game.Server() != nil {
		req := request.Request{Save: []string{saveName}}
		err = sm.hud.game.Server().Send(req)
		if err != nil {
			return fmt.Errorf("unable to send save request: %v",
				err)
		}
		return nil
	}
	modPath := filepath.Join(config.ModulesPath, saveName)
	err = flamedata.ExportModule(modPath, sm.hud.Game().Data())
	if err != nil {
		return fmt.Errorf("unable to export module: %v", err)
	}
	return nil
}
