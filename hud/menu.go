/*
 * menu.go
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
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"

	"github.com/isangeles/flame/core/data/text/lang"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/log"
)

// Struct for HUD menu.
type Menu struct {
	hud         *HUD
	bgSpr       *pixel.Sprite
	bgDraw      *imdraw.IMDraw
	drawArea    pixel.Rect
	titleText   *mtk.Text
	closeButton *mtk.Button
	exitButton  *mtk.Button
	opened      bool
}

// newMenu creates menu for HUD.
func newMenu(hud *HUD) *Menu {
	m := new(Menu)
	m.hud = hud
	// Background.
	bg, err := data.PictureUI("menubg.png")
	if err != nil { // fallback
		m.bgDraw = imdraw.New(nil)
		log.Err.Printf("hud_menu:bg_texture_not_found:%v", err)
	} else {
		m.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	}
	// Title.
	m.titleText = mtk.NewText(lang.Text("gui", "hud_menu_title"), mtk.SIZE_SMALL, 0)
	// Buttons.
	closeButtonBG, err := data.PictureUI("closebutton1.png")
	if err != nil { // fallback
		m.closeButton = mtk.NewButton(mtk.SIZE_SMALL, mtk.SHAPE_SQUARE, accent_color,
			"", "")
	} else {
		m.closeButton = mtk.NewButtonSprite(closeButtonBG, mtk.SIZE_SMALL, "", "")
	}
	m.closeButton.SetOnClickFunc(m.onCloseButtonClicked)
	exitButtonBG, err := data.PictureUI("button_green.png")
	if err != nil { // fallback
		m.exitButton = mtk.NewButton(mtk.SIZE_SMALL, mtk.SHAPE_RECTANGLE, accent_color,
			lang.Text("gui", ""), lang.Text("gui", ""))
	} else {
		m.exitButton = mtk.NewButtonSprite(exitButtonBG, mtk.SIZE_SMALL, lang.Text("gui", "exit_b_label"),
			lang.Text("gui", "exit_b_info"))
	}
	m.exitButton.SetOnClickFunc(m.onExitButtonClicked)
	return m
}

// Draw draws menu.
func (m *Menu) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	m.drawArea = mtk.MatrixToDrawArea(matrix, m.Bounds())
	// Background.
	if m.bgSpr != nil {
		m.bgSpr.Draw(win.Window, matrix)
	} else {
		m.drawIMBackground(win.Window)
	}
	// Title.
	titleTextPos := mtk.ConvVec(pixel.V(0, m.Bounds().Max.Y/2 - 25))
	m.titleText.Draw(win.Window, matrix.Moved(titleTextPos))
	// Buttons.
	closeButtonPos := mtk.ConvVec(pixel.V(m.Bounds().Max.X/2 - 20, m.Bounds().Max.Y/2 - 15))
	m.closeButton.Draw(win.Window, matrix.Moved(closeButtonPos))
	exitButtonPos := mtk.ConvVec(pixel.V(0, -m.Bounds().Max.X/2 - 20))
	m.exitButton.Draw(win.Window, matrix.Moved(exitButtonPos))
}

// Update updates menu.
func (m *Menu) Update(win *mtk.Window) {
	if !m.hud.chat.Activated() { // block key events if chat is active
		if win.JustPressed(pixelgl.KeyEscape) {
			// Show menu.
			if !m.hud.menu.Opened() {
				m.Show(true)
			} else {
				m.Show(false)
			}
		}
	}
	m.closeButton.Update(win)
	m.exitButton.Update(win)
}

// DrawArea returns current draw area of
// menu background.
func (m *Menu) DrawArea() pixel.Rect {
	return m.drawArea
}

// Bounds return size parameter of menu background.
func (m *Menu) Bounds() pixel.Rect {
	if m.bgSpr == nil {
		return pixel.R(0, 0, mtk.ConvSize(0), mtk.ConvSize(0))
	}
	return m.bgSpr.Frame()
}

// Opened checks wheter menu is open.
func (m *Menu) Opened() bool {
	return m.opened
}

// Show toggles menu visibility.
func (m *Menu) Show(show bool) {
	m.opened = show
	m.hud.camera.Lock(show)
	m.hud.game.Pause(show)
}

// drawIMBackground draws menu background with
// Pixel IMDraw.
func (m *Menu) drawIMBackground(t pixel.Target) {
	// TODO: draw background with IMDraw.
}

// Triggered after close button clicked.
func (m *Menu) onCloseButtonClicked(b *mtk.Button) {
	m.Show(false)
}

// Triggered after exit button clicked.
func (m *Menu) onExitButtonClicked(b *mtk.Button) {
	m.hud.Exit()
}
