/*
 * menubar.go
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

// Struct for HUD menu bar.
type MenuBar struct {
	hud        *HUD
	bgSpr      *pixel.Sprite
	bgDraw     *imdraw.IMDraw
	drawArea   pixel.Rect
	menuButton *mtk.Button
	invButton  *mtk.Button
}

// neMenuBar creates new menu bar for HUD.
func newMenuBar(hud *HUD) *MenuBar {
	mb := new(MenuBar)
	mb.hud = hud
	// Background.
	bg, err := data.PictureUI("menubar.png")
	if err != nil { // fallback
		mb.bgDraw = imdraw.New(nil)
	} else {
		mb.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	}
	// Buttons.
	menuButtonBG, err := data.PictureUI("menubutton.png")
	if err != nil { // fallback
		mb.menuButton = mtk.NewButton(mtk.SIZE_MINI, mtk.SHAPE_SQUARE, accent_color,
			"", lang.Text("gui", "hud_bar_menu_open_info"))
	} else {
		mb.menuButton = mtk.NewButtonSprite(menuButtonBG, mtk.SIZE_MINI, "",
			lang.Text("gui", "hud_bar_menu_open_info"))
	}
	mb.menuButton.SetOnClickFunc(mb.onMenuButtonClicked)
	invButtonBG, err := data.PictureUI("inventorybutton.png")
	if err != nil { // fallback
		mb.invButton = mtk.NewButton(mtk.SIZE_MINI, mtk.SHAPE_SQUARE, accent_color,
			"", lang.Text("gui", "hud_bar_menu_open_info"))
	} else {
		mb.invButton = mtk.NewButtonSprite(invButtonBG, mtk.SIZE_MINI, "",
			lang.Text("gui", "hud_bar_inv_open_info"))
	}
	mb.invButton.SetOnClickFunc(mb.onInvButtonClicked)
	return mb
}

// Draw draws menu bar.
func (mb *MenuBar) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	mb.drawArea = mtk.MatrixToDrawArea(matrix, mb.Bounds())
	// Background.
	if mb.bgSpr != nil {
		mb.bgSpr.Draw(win.Window, matrix)
	} else {
		mb.drawIMBackground(win.Window)
	}
	// Buttons.
	menuButtonPos := mtk.ConvVec(pixel.V(mb.Bounds().Max.X/2 - 30, 0))
	mb.menuButton.Draw(win.Window, matrix.Moved(menuButtonPos))
	invButtonPos := mtk.ConvVec(pixel.V(mb.Bounds().Max.X/2 - 65, 0))
	mb.invButton.Draw(win.Window, matrix.Moved(invButtonPos))
}

// Update updates menu bar.
func (mb *MenuBar) Update(win *mtk.Window) {
	// Buttons.
	mb.menuButton.Update(win)
	mb.invButton.Update(win)
}

// Bounds returns bounds of bar background.
func (mb *MenuBar) Bounds() pixel.Rect {
	if mb.bgSpr == nil {
		return pixel.R(0, 0, mtk.ConvSize(0), mtk.ConvSize(0))
	}
	return mb.bgSpr.Frame()
}

// DrawArea return current draw area of bar background.
func (mb *MenuBar) DrawArea() pixel.Rect {
	return mb.drawArea
}

// drawIMBackground draws menu bar background with IMDraw.
func (mb *MenuBar) drawIMBackground(t pixel.Target) {
	// TODO: draw background.
}

// Triggered after menu button clicked.
func (mb *MenuBar) onMenuButtonClicked(b *mtk.Button) {
	if mb.hud.menu.Opened() {
		mb.hud.menu.Show(false)
	} else {
		mb.hud.menu.Show(true)
	}
}

// Triggered after inventory button clicked.
func (mb *MenuBar) onInvButtonClicked(b *mtk.Button) {
	if mb.hud.inv.Opened() {
		mb.hud.inv.Show(false)
	} else {
		mb.hud.inv.Show(true)
	}
}
