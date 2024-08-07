/*
 * menu.go
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
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"

	"github.com/isangeles/flame/data/res/lang"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/data/res/graphic"
	"github.com/isangeles/mural/log"
)

var (
	menuKey = pixelgl.KeyEscape
)

// Struct for HUD menu.
type Menu struct {
	hud         *HUD
	bgSpr       *pixel.Sprite
	bgDraw      *imdraw.IMDraw
	drawArea    pixel.Rect
	titleText   *mtk.Text
	closeButton *mtk.Button
	saveButton  *mtk.Button
	exitButton  *mtk.Button
	opened      bool
	focused     bool
}

// newMenu creates menu for HUD.
func newMenu(hud *HUD) *Menu {
	m := new(Menu)
	m.hud = hud
	// Background.
	m.bgDraw = imdraw.New(nil)
	bg := graphic.Textures["menubg.png"]
	if bg != nil {
		m.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	} else {
		log.Err.Printf("hud menu: bg texture not found")
	}
	// Title.
	titleParams := mtk.Params{
		FontSize: mtk.SizeSmall,
	}
	m.titleText = mtk.NewText(titleParams)
	m.titleText.SetText(lang.Text("hud_menu_title"))
	// Close button.
	closeButtonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		Shape:     mtk.ShapeSquare,
		MainColor: accentColor,
	}
	m.closeButton = mtk.NewButton(closeButtonParams)
	closeButtonBG := graphic.Textures["closebutton1.png"]
	if closeButtonBG != nil {
		closeBG := pixel.NewSprite(closeButtonBG, closeButtonBG.Bounds())
		m.closeButton.SetBackground(closeBG)
	} else {
		log.Err.Printf("hud menu: unable to retrieve exit button texture")
	}
	m.closeButton.SetOnClickFunc(m.onCloseButtonClicked)
	// Menu buttons.
	menuButtonParams := mtk.Params{
		Size:      mtk.SizeMini,
		FontSize:  mtk.SizeMini,
		Shape:     mtk.ShapeRectangle,
		MainColor: accentColor,
	}
	greenButtonBG := graphic.Textures["button_green.png"]
	if greenButtonBG == nil {
		log.Err.Printf("hud menu: unable to retrieve green button texture")
	}
	m.saveButton = mtk.NewButton(menuButtonParams)
	m.saveButton.SetLabel(lang.Text("savegame_button_label"))
	m.saveButton.SetInfo(lang.Text("savegame_button_info"))
	if greenButtonBG != nil {
		bg := pixel.NewSprite(greenButtonBG, greenButtonBG.Bounds())
		m.saveButton.SetBackground(bg)
	}
	m.saveButton.SetOnClickFunc(m.onSaveButtonClicked)
	m.exitButton = mtk.NewButton(menuButtonParams)
	m.exitButton.SetLabel(lang.Text("exit_button_label"))
	m.exitButton.SetInfo(lang.Text("exit_hud_button_info"))
	if greenButtonBG != nil {
		bg := pixel.NewSprite(greenButtonBG, greenButtonBG.Bounds())
		m.exitButton.SetBackground(bg)
	}
	m.exitButton.SetOnClickFunc(m.onExitButtonClicked)
	return m
}

// Draw draws menu.
func (m *Menu) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	m.drawArea = mtk.MatrixToDrawArea(matrix, m.Size())
	// Background.
	if m.bgSpr != nil {
		m.bgSpr.Draw(win.Window, matrix)
	} else {
		mtk.DrawRect(win.Window, m.DrawArea(), nil)
	}
	// Title.
	titleTextPos := pixel.V(mtk.ConvSize(0),
		m.Size().Y/2-mtk.ConvSize(25))
	m.titleText.Draw(win.Window, matrix.Moved(titleTextPos))
	// Buttons.
	closeButtonPos := pixel.V(m.Size().X/2-mtk.ConvSize(20),
		m.Size().Y/2-mtk.ConvSize(15))
	saveButtonPos := pixel.V(mtk.ConvSize(0), -m.Size().X/2+mtk.ConvSize(20))
	exitButtonPos := pixel.V(mtk.ConvSize(0), -m.Size().X/2-mtk.ConvSize(20))
	m.closeButton.Draw(win.Window, matrix.Moved(closeButtonPos))
	m.saveButton.Draw(win.Window, matrix.Moved(saveButtonPos))
	m.exitButton.Draw(win.Window, matrix.Moved(exitButtonPos))
}

// Update updates menu.
func (m *Menu) Update(win *mtk.Window) {
	// Key events.
	if win.JustPressed(menuKey) {
		// Show menu.
		if m.Opened() {
			m.Hide()
		} else if !m.hud.menuOpen() {
			m.Show()
		}
	}
	// Elements.
	if m.Opened() {
		m.closeButton.Update(win)
		m.saveButton.Update(win)
		m.exitButton.Update(win)
	}
}

// DrawArea returns current draw area of
// menu background.
func (m *Menu) DrawArea() pixel.Rect {
	return m.drawArea
}

// Size return size of menu background.
func (m *Menu) Size() pixel.Vec {
	if m.bgSpr == nil {
		// TODO: menu draw background size.
		return pixel.V(mtk.ConvSize(0), mtk.ConvSize(0))
	}
	return mtk.ConvVec(m.bgSpr.Frame().Size())
}

// Opened checks wheter menu is open.
func (m *Menu) Opened() bool {
	return m.opened
}

// Show shows menu.
func (m *Menu) Show() {
	m.opened = true
	m.hud.camera.Lock(true)
	m.hud.game.SetPause(true)
	m.hud.UserFocus().Focus(m)
}

// Hide hides menu.
func (m *Menu) Hide() {
	m.opened = false
	m.hud.camera.Lock(false)
	m.hud.game.SetPause(false)
	m.hud.UserFocus().Focus(nil)
}

// Focused checks whether menu is focused.
func (m *Menu) Focused() bool {
	return m.focused
}

// Focus toggles menu focus.
func (m *Menu) Focus(focus bool) {
	m.focused = focus
}

// Triggered after close button clicked.
func (m *Menu) onCloseButtonClicked(b *mtk.Button) {
	m.Hide()
}

// Triggered after save button clicked.
func (m *Menu) onSaveButtonClicked(b *mtk.Button) {
	m.Hide()
	m.hud.savemenu.Show()
}

// Triggered after exit button clicked.
func (m *Menu) onExitButtonClicked(b *mtk.Button) {
	m.hud.Exit()
}
