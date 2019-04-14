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
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"

	"github.com/isangeles/flame/core/data/text/lang"
	
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/mtk"
)

// Struct for HUD save game menu.
type SaveMenu struct {
	hud         *HUD
	bgSpr       *pixel.Sprite
	bgDraw      *imdraw.IMDraw
	drawArea    pixel.Rect
	titleText   *mtk.Text
	savesList   *mtk.List
	closeButton *mtk.Button
	opened      bool
	focused     bool
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
	sm.savesList = mtk.NewList(sm.Bounds().Size(), mtk.SIZE_MINI, main_color,
		sec_color, accent_color)
	// Buttons.
	sm.closeButton = mtk.NewButton(mtk.SIZE_SMALL, mtk.SHAPE_SQUARE,
		accent_color, "", "")
	closeButtonBG, err := data.PictureUI("closebutton1.png")
	if err == nil {
		closeBG := pixel.NewSprite(closeButtonBG, closeButtonBG.Bounds())
		sm.closeButton.SetBackground(closeBG)
	}
	sm.closeButton.SetOnClickFunc(sm.onCloseButtonClicked)
	return sm
}

// Draw draws menu.
func (sm *SaveMenu) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	sm.drawArea = mtk.MatrixToDrawArea(matrix, sm.Bounds())
	// Background.
	if sm.bgSpr != nil {
		sm.bgSpr.Draw(win.Window, matrix)
	} else {
		mtk.DrawRectangle(win.Window, sm.Bounds(), nil)
	}
	// Title.
	titleTextPos := mtk.ConvVec(pixel.V(0, sm.Bounds().Max.Y/2 - 25))
	sm.titleText.Draw(win.Window, matrix.Moved(titleTextPos))
	// Buttons.
	closeButtonPos := mtk.ConvVec(pixel.V(sm.Bounds().Max.X/2 - 20,
		sm.Bounds().Max.Y/2 - 15))
	sm.closeButton.Draw(win.Window, matrix.Moved(closeButtonPos))
}

// Update updates menu.
func (sm *SaveMenu) Update(win *mtk.Window) {
	// Elements.
	sm.closeButton.Update(win)
}

// Bounds returns menu background size.
func (sm *SaveMenu) Bounds() pixel.Rect {
	if sm.bgSpr == nil {
		// TODO: menu draw background size.
		return pixel.R(0, 0, mtk.ConvSize(0), mtk.ConvSize(0))
	}
	return sm.bgSpr.Frame()
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
}

// Focused checks if menu us focused.
func (sm *SaveMenu) Focused() bool {
	return sm.focused
}

// Focus toggles menu focus.
func (sm *SaveMenu) Focus(focus bool) {
	sm.focused = focus
}

// Triggered after close button clicked.
func (sm *SaveMenu) onCloseButtonClicked(b *mtk.Button) {
	sm.Show(false)
}
