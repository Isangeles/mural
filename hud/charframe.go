/*
 * charframe.go
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
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// Struct for character frame with portrait
// and health/mana bars.
type CharFrame struct {
	hud      *HUD
	char     *object.Avatar
	bgSpr    *pixel.Sprite
	bgDraw   *imdraw.IMDraw
	drawArea pixel.Rect
	hpBar    *mtk.ProgressBar
}

// newCharFrame creates new HUD character frame for
// specified character avatar.
func newCharFrame(hud *HUD, char *object.Avatar) (*CharFrame, error) {
	cf := new(CharFrame)
	cf.char = char
	// Background.
	bg, err := data.PictureUI("charframe.png")
	if err == nil {
		cf.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	} else { // fallback
		cf.bgDraw = imdraw.New(nil)
		log.Err.Printf("hud_char_frame:bg_texture_not_found:%v", err)
	}
	// Bars.
	hpBarPic, err := data.PictureUI("bar_red.png")
	if err == nil {
		cf.hpBar = mtk.NewProgressBarSprite(hpBarPic, mtk.SIZE_MINI,
			lang.Text("gui", "char_frame_hp_bar_label"), 100)
	} else { // fallback
		cf.hpBar = mtk.NewProgressBar(mtk.SIZE_MINI, accent_color,
			lang.Text("gui", "char_frame_hp_bar_label"), 100)
		log.Err.Printf("hud_char_frame:hp_bar_texture_not_found:%v", err)
	}
	cf.hpBar.SetValue(50)
	return cf, nil
}

// Draw draws character frame.
func (cf *CharFrame) Draw(win *mtk.Window, matrix pixel.Matrix) {
	cf.drawArea = mtk.MatrixToDrawArea(matrix, cf.Bounds())
	// Background.
	if cf.bgSpr != nil {
		cf.bgSpr.Draw(win.Window, matrix)
	} else {
		cf.drawIMBackground(win)
	}
	// Portrait.
	portraitPos := pixel.V(mtk.ConvSize(-70), mtk.ConvSize(5))
	cf.char.Portrait().Draw(win, matrix.Scaled(cf.DrawArea().Center(),
		0.5).Moved(portraitPos))
	// Bars.
	hpBarPos := pixel.V(mtk.ConvSize(35), mtk.ConvSize(25))
	cf.hpBar.Draw(win.Window, matrix.Moved(hpBarPos))
}

// Update updates character frame.
func (cf *CharFrame) Update(win *mtk.Window) {
	cf.hpBar.Update(win)
}

// Bounds returns size bounds of character frame.
func (cf *CharFrame) Bounds() pixel.Rect {
	if cf.bgSpr == nil {
		// TODO: bounds for draw background.
		return pixel.R(0, 0, 0, 0)
	}
	return cf.bgSpr.Frame()
}

// DrawArea retruns current frame draw area.
func (cf *CharFrame) DrawArea() pixel.Rect {
	return cf.drawArea
}

// drawIMBackground draw character frame with pixel
// IMDraw.
func (cf *CharFrame) drawIMBackground(win *mtk.Window) {
	// TODO: draw background with IMDraw.
}
