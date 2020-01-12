/*
 * characterwindow.go
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
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame/core/data/res/lang"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data"
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
	bg, err := data.PictureUI("menubg.png")
	if err == nil {
		cw.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	} else {
		log.Err.Printf("hud_char: fail to retrieve background tex: %v", err)
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
	closeButtonBG, err := data.PictureUI("closebutton1.png")
	if err == nil {
		closeBG := pixel.NewSprite(closeButtonBG,
			closeButtonBG.Bounds())
		cw.closeButton.SetBackground(closeBG)
	} else {
		log.Err.Printf("hud_char: fail to retrieve close button tex: %v", err)
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
		mtk.DrawRectangle(win, cw.DrawArea(), mainColor)
	}
	// Title.
	titleTextMove := mtk.ConvVec(pixel.V(0, cw.Size().Y/2-25))
	cw.titleText.Draw(win, matrix.Moved(titleTextMove))
	// Buttons.
	closeButtonMove := mtk.ConvVec(pixel.V(cw.Size().X/2-20,
		cw.Size().Y/2-15))
	cw.closeButton.Draw(win, matrix.Moved(closeButtonMove))
	// Char info.
	infoMove := mtk.ConvVec(pixel.V(0, -20))
	cw.charInfo.Draw(win, matrix.Moved(infoMove))
}

// Update updates window.
func (cw *CharacterWindow) Update(win *mtk.Window) {
	// Key events.
	if !cw.hud.Chat().Activated() && win.JustPressed(charinfoKey) {
		cw.Show(!cw.Opened())
	}
	// Elements.
	if cw.Opened() {
		cw.closeButton.Update(win)
		cw.charInfo.Update(win)
	}
}

// Show toggles window visibility.
func (cw *CharacterWindow) Show(show bool) {
	cw.opened = show
	if cw.Opened() {
		cw.updateInfo()
	}
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
	return cw.bgSpr.Frame().Size()
}

// updateInfo updates info textbox with
// information about active player.
func (cw *CharacterWindow) updateInfo() {
	infoForm := `
Name:       %s
Level:      %d
Gender:     %s
Race:       %s
Alignment   %s
Attributes: %s`
	pc := cw.hud.ActivePlayer()
	info := fmt.Sprintf(infoForm, pc.Name(), pc.Level(), lang.Text(pc.Gender().ID()),
		lang.Text(pc.Race().ID()), lang.Text(pc.Alignment().ID()),
		pc.Attributes())
	cw.charInfo.SetText(info)
}

// Triggered on close button clicked.
func (cw *CharacterWindow) onCloseButtonClicked(b *mtk.Button) {
	cw.Show(false)
}
