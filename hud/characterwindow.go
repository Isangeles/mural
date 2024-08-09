/*
 * characterwindow.go
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

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"

	"github.com/isangeles/flame/data/res/lang"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/data/res/graphic"
	"github.com/isangeles/mural/log"
)

var (
	charinfoKey = pixelgl.KeyC
)

// Struct for HUD character window.
type CharacterWindow struct {
	hud         *HUD
	bgSpr       *pixel.Sprite
	bgDraw      *imdraw.IMDraw
	drawArea    pixel.Rect
	titleText   *mtk.Text
	closeButton *mtk.Button
	charInfo    *mtk.Textbox
	opened      bool
	focused     bool
}

// newCharacterWindow creates new character
// window for HUD.
func newCharacterWindow(hud *HUD) *CharacterWindow {
	cw := new(CharacterWindow)
	cw.hud = hud
	// Background.
	cw.bgDraw = imdraw.New(nil)
	bg := graphic.Textures["menubg.png"]
	if bg != nil {
		cw.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	} else {
		log.Err.Printf("hud_char: unable to retrieve background texture")
	}
	// Title.
	titleParams := mtk.Params{
		FontSize: mtk.SizeSmall,
	}
	cw.titleText = mtk.NewText(titleParams)
	cw.titleText.SetText(lang.Text("hud_charwin_title"))
	// Close button.
	closeButtonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		Shape:     mtk.ShapeSquare,
		MainColor: accentColor,
	}
	cw.closeButton = mtk.NewButton(closeButtonParams)
	closeButtonBG := graphic.Textures["closebutton1.png"]
	if closeButtonBG != nil {
		closeBG := pixel.NewSprite(closeButtonBG, closeButtonBG.Bounds())
		cw.closeButton.SetBackground(closeBG)
	} else {
		log.Err.Printf("hud_char: unable to retrieve close button texture")
	}
	cw.closeButton.SetOnClickFunc(cw.onCloseButtonClicked)
	// Char info.
	infoSize := pixel.V(cw.Size().X-mtk.ConvSize(20),
		cw.Size().Y-mtk.ConvSize(70))
	charInfoParams := mtk.Params{
		SizeRaw: infoSize,
		FontSize: mtk.SizeSmall,
		MainColor: mainColor,
		AccentColor: accentColor,
	}
	cw.charInfo = mtk.NewTextbox(charInfoParams)
	upButtonTex := graphic.Textures["scrollup.png"]
	if upButtonTex != nil {
		upSprite := pixel.NewSprite(upButtonTex, upButtonTex.Bounds())
		cw.charInfo.SetUpButtonBackground(upSprite)
	}
	downButtonTex := graphic.Textures["scrolldown.png"]
	if downButtonTex != nil {
		downSprite := pixel.NewSprite(downButtonTex, downButtonTex.Bounds())
		cw.charInfo.SetDownButtonBackground(downSprite)
	}
	return cw
}

// Draw draws window.
func (cw *CharacterWindow) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	cw.drawArea = mtk.MatrixToDrawArea(matrix, cw.Size())
	// Background.
	if cw.bgSpr != nil {
		cw.bgSpr.Draw(win, matrix)
	} else {
		mtk.DrawRect(win, cw.DrawArea(), mainColor)
	}
	// Title.
	titleTextMove := pixel.V(0, cw.Size().Y/2-mtk.ConvSize(25))
	cw.titleText.Draw(win, matrix.Moved(titleTextMove))
	// Buttons.
	closeButtonMove := pixel.V(cw.Size().X/2-mtk.ConvSize(20),
		cw.Size().Y/2-mtk.ConvSize(15))
	cw.closeButton.Draw(win, matrix.Moved(closeButtonMove))
	// Char info.
	infoMove := mtk.ConvVec(pixel.V(0, -20))
	cw.charInfo.Draw(win, matrix.Moved(infoMove))
}

// Update updates window.
func (cw *CharacterWindow) Update(win *mtk.Window) {
	// Key events.
	if !cw.hud.Chat().Activated() && win.JustPressed(charinfoKey) {
		if cw.Opened() {
			cw.Hide()
		} else {
			cw.Show()
		}
	}
	if win.JustPressed(exitKey) {
		cw.Hide()
	}
	// Elements.
	if cw.Opened() {
		cw.closeButton.Update(win)
		cw.charInfo.Update(win)
	}
}

// Show shows window.
func (cw *CharacterWindow) Show() {
	cw.opened = true
	cw.updateInfo()
}

// Hide hides window.
func (cw *CharacterWindow) Hide() {
	cw.opened = false
}

// Opened checks if window is open.
func (cw *CharacterWindow) Opened() bool {
	return cw.opened
}

// DrawArea returns window draw area.
func (cw *CharacterWindow) DrawArea() pixel.Rect {
	return cw.drawArea
}

// Size returns window background size.
func (cw *CharacterWindow) Size() pixel.Vec {
	if cw.bgSpr == nil {
		return mtk.ConvVec(pixel.V(250, 350))
	}
	return mtk.ConvVec(cw.bgSpr.Frame().Size())
}

// updateInfo updates info textbox with
// information about active player.
func (cw *CharacterWindow) updateInfo() {
	infoForm := `
Name:       %s
Level:      %d
Experience: %d/%d
Gender:     %s
Race:       %s
Alignment   %s
Attributes: %s`
	pc := cw.hud.Game().ActivePlayerChar()
	race := lang.Text(pc.Race().ID())
	info := fmt.Sprintf(infoForm, pc.Name(), pc.Level(), pc.Experience(), pc.MaxExperience(),
		lang.Text(string(pc.Gender())), race, lang.Text(string(pc.Alignment())),
		pc.Attributes())
	cw.charInfo.SetText(info)
}

// Triggered on close button clicked.
func (cw *CharacterWindow) onCloseButtonClicked(b *mtk.Button) {
	cw.Hide()
}
