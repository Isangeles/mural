/*
 * menu.go
 *
 * Copyright 2019-2020 Dariusz Sikora <dev@isangeles.pl>
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
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame/core/data/res/lang"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data"
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
	bg, err := data.PictureUI("menubg.png")
	if err == nil {
		m.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	} else {
		log.Err.Printf("hud menu: bg texture not found: %v", err)
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
	closeButtonBG, err := data.PictureUI("closebutton1.png")
	if err == nil {
		closeBG := pixel.NewSprite(closeButtonBG, closeButtonBG.Bounds())
		m.closeButton.SetBackground(closeBG)
	} else {
		log.Err.Printf("hud menu: fail to retrieve exit button texture: %v", err)
	}
	m.closeButton.SetOnClickFunc(m.onCloseButtonClicked)
	// Menu buttons.
	menuButtonParams := mtk.Params{
		Size:      mtk.SizeMini,
		FontSize:  mtk.SizeMini,
		Shape:     mtk.ShapeRectangle,
		MainColor: accentColor,
	}
	greenButtonBG, err := data.PictureUI("button_green.png")
	if err != nil {
		log.Err.Printf("hud menu: fail to retrieve green button texture: %v", err)
	}
	m.saveButton = mtk.NewButton(menuButtonParams)
	m.saveButton.SetLabel(lang.Text("savegame_b_label"))
	m.saveButton.SetInfo(lang.Text("savegame_b_info"))
	if greenButtonBG != nil {
		bg := pixel.NewSprite(greenButtonBG, greenButtonBG.Bounds())
		m.saveButton.SetBackground(bg)
	}
	m.saveButton.SetOnClickFunc(m.onSaveButtonClicked)
	m.exitButton = mtk.NewButton(menuButtonParams)
	m.exitButton.SetLabel(lang.Text("exit_b_label"))
	m.exitButton.SetInfo(lang.Text("exit_hud_b_info"))
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
		mtk.DrawRectangle(win.Window, m.DrawArea(), nil)
	}
	// Title.
	titleTextPos := mtk.ConvVec(pixel.V(0, m.Size().Y/2-25))
	m.titleText.Draw(win.Window, matrix.Moved(titleTextPos))
	// Buttons.
	closeButtonPos := mtk.ConvVec(pixel.V(m.Size().X/2-20,
		m.Size().Y/2-15))
	saveButtonPos := mtk.ConvVec(pixel.V(0, -m.Size().X/2+20))
	exitButtonPos := mtk.ConvVec(pixel.V(0, -m.Size().X/2-20))
	m.closeButton.Draw(win.Window, matrix.Moved(closeButtonPos))
	m.saveButton.Draw(win.Window, matrix.Moved(saveButtonPos))
	m.exitButton.Draw(win.Window, matrix.Moved(exitButtonPos))
}

// Update updates menu.
func (m *Menu) Update(win *mtk.Window) {
	// Key events.
	if win.JustPressed(menuKey) {
		// Show menu.
		m.Show(!m.Opened())
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
	return m.bgSpr.Frame().Size()
}

// Opened checks wheter menu is open.
func (m *Menu) Opened() bool {
	return m.opened
}

// Show toggles menu visibility.
func (m *Menu) Show(show bool) {
	m.opened = show
	m.hud.camera.Lock(m.Opened())
	m.hud.game.Pause(m.Opened())
	if m.Opened() {
		m.hud.UserFocus().Focus(m)
	} else {
		m.hud.UserFocus().Focus(nil)
	}
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
	m.Show(false)
}

// Triggered after save button clicked.
func (m *Menu) onSaveButtonClicked(b *mtk.Button) {
	m.Show(false)
	m.hud.savemenu.Show(true)
}

// Triggered after exit button clicked.
func (m *Menu) onExitButtonClicked(b *mtk.Button) {
	m.hud.Exit()
}
