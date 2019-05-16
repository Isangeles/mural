/*
 * objectframe.go
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
	
	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// Struct for object frame with portrait,
// health/mana bars and effects icons.
type ObjectFrame struct {
	hud      *HUD
	object   Target
	portrait *pixel.Sprite
	bgSpr    *pixel.Sprite
	bgDraw   *imdraw.IMDraw
	drawArea pixel.Rect
	hpBar    *mtk.ProgressBar
	manaBar  *mtk.ProgressBar
}

// Interface for HUD frame object.
type Target interface {
	Name() string
	Health() int
	MaxHealth() int
	Mana() int
	MaxMana() int
	Portrait() pixel.Picture
	Effects() []*object.EffectGraphic
}

// newCharFrame creates new HUD character frame for
// specified character avatar.
func newObjectFrame(hud *HUD) *ObjectFrame {
	of := new(ObjectFrame)
	of.hud = hud
	// Background.
	bg, err := data.PictureUI("charframe.png")
	if err != nil { // fallback
		of.bgDraw = imdraw.New(nil)
		log.Err.Printf("hud_char_frame:bg_texture_not_found:%v", err)
	} else {
		of.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	}
	// Bars.
	of.hpBar = mtk.NewProgressBar(mtk.SIZE_MINI, accent_color)
	of.hpBar.SetLabel(lang.Text("gui", "char_frame_hp_bar_label"))
	hpBarPic, err := data.PictureUI("bar_red.png")
	if err != nil {
		log.Err.Printf("hud_char_frame:hp_bar_texture_not_found:%v", err)
	} else {
		of.hpBar.SetBackground(hpBarPic)
	}
	of.manaBar = mtk.NewProgressBar(mtk.SIZE_MINI, accent_color)
	of.manaBar.SetLabel(lang.Text("gui", "char_frame_mana_bar_label"))
	manaBarPic, err := data.PictureUI("bar_blue.png")
	if err != nil {
		log.Err.Printf("hud_char_frame:mana_bar_texture_not_found:%v", err)
	} else {
		of.manaBar.SetBackground(manaBarPic)
	}
	return of
}

// Draw draws character frame.
func (of *ObjectFrame) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	of.drawArea = mtk.MatrixToDrawArea(matrix, of.Size())
	// Portrait.
	portraitPos := pixel.V(mtk.ConvSize(-70), mtk.ConvSize(0))
	if of.object != nil && of.portrait != nil {
		of.portrait.Draw(win, matrix.Scaled(of.DrawArea().Center(),
			0.6).Moved(portraitPos))
	}
	// Background.
	if of.bgSpr != nil {
		of.bgSpr.Draw(win.Window, matrix)
	} else {
		of.drawIMBackground(win.Window)
	}
	// Bars.
	hpBarPos := pixel.V(mtk.ConvSize(35), mtk.ConvSize(25))
	of.hpBar.Draw(win.Window, matrix.Moved(hpBarPos))
	manaBarPos := pixel.V(mtk.ConvSize(35), mtk.ConvSize(10))
	of.manaBar.Draw(win.Window, matrix.Moved(manaBarPos))
	// Effects icons.
	if of.object != nil {
		iconsStartPos := pixel.V(mtk.ConvSize(0), mtk.ConvSize(-30))
		for i, e := range of.object.Effects() {
			// Icon.
			iconMove := iconsStartPos
			iconMove.X += mtk.ConvSize(e.Icon().Frame().W()) * float64(i)
			if i % 3 == 0 {
				iconMove.X = iconsStartPos.X
				iconMove.Y -= mtk.ConvSize(e.Icon().Frame().H()) * float64(i/3)
			}
			e.DrawIcon(win.Window, matrix.Moved(iconMove))
		}
	}
}

// Update updates character frame.
func (of *ObjectFrame) Update(win *mtk.Window) {
	if of.object == nil {
		return
	}
	// Bars values.
	of.hpBar.SetMax(of.object.MaxHealth())
	of.hpBar.SetValue(of.object.Health())
	of.manaBar.SetMax(of.object.MaxMana())
	of.manaBar.SetValue(of.object.Mana())
	// Bars.
	of.hpBar.Update(win)
	of.manaBar.Update(win)
}

// Size returns size of character frame background.
func (of *ObjectFrame) Size() pixel.Vec {
	if of.bgSpr == nil {
		return pixel.V(mtk.ConvSize(200), mtk.ConvSize(50))
	}
	return of.bgSpr.Frame().Size()
}

// DrawArea retruns current frame draw area.
func (of *ObjectFrame) DrawArea() pixel.Rect {
	return of.drawArea
}

// SetObject sets specified object as object to
// display in frame.
func (of *ObjectFrame) SetObject(ob Target) {
	of.object = ob
	if ob.Portrait() != nil {
		of.portrait = pixel.NewSprite(ob.Portrait(),
			ob.Portrait().Bounds())
	} else {
		of.portrait = nil
	}
}

// drawIMBackground draw character frame with pixel
// IMDraw.
func (of *ObjectFrame) drawIMBackground(t pixel.Target) {
	of.bgDraw.Clear()
	of.bgDraw.Color = pixel.ToRGBA(main_color)
	of.bgDraw.Push(of.DrawArea().Min)
	of.bgDraw.Push(of.DrawArea().Max)
	of.bgDraw.Rectangle(0)
	of.bgDraw.Draw(t)
}
