/*
 * objectframe.go
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

	"github.com/isangeles/flame/data/res/lang"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/data/res/graphic"
	"github.com/isangeles/mural/object"
	"github.com/isangeles/mural/log"
)

// Struct for object frame with portrait,
// health/mana bars and effects icons.
type ObjectFrame struct {
	hud      *HUD
	object   FrameTarget
	portrait *pixel.Sprite
	bgSpr    *pixel.Sprite
	bgDraw   *imdraw.IMDraw
	drawArea pixel.Rect
	hpBar    *mtk.ProgressBar
	manaBar  *mtk.ProgressBar
}

// Interface for HUD frame object.
type FrameTarget interface {
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
	bg := graphic.Textures["charframe.png"]
	if bg != nil {
		of.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	} else { // fallback
		of.bgDraw = imdraw.New(nil)
		log.Err.Printf("hud object frame: bg texture not found")
	}
	// Bars.
	of.hpBar = mtk.NewProgressBar(mtk.SizeMini, accentColor)
	of.hpBar.SetLabel(lang.Text("char_frame_hp_bar_label"))
	hpBarPic := graphic.Textures["bar_red.png"]
	if hpBarPic != nil {
		of.hpBar.SetBackground(hpBarPic)
	} else {
		log.Err.Printf("hud object frame: hp bar texture not found")
	}
	of.manaBar = mtk.NewProgressBar(mtk.SizeMini, accentColor)
	of.manaBar.SetLabel(lang.Text("char_frame_mana_bar_label"))
	manaBarPic := graphic.Textures["bar_blue.png"]
	if manaBarPic != nil {
		of.manaBar.SetBackground(manaBarPic)
	} else {
		log.Err.Printf("hud object frame: mana bar texture not found")
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
	return mtk.ConvVec(of.bgSpr.Frame().Size())
}

// DrawArea retruns current frame draw area.
func (of *ObjectFrame) DrawArea() pixel.Rect {
	return of.drawArea
}

// SetObject sets specified object as object to
// display in frame.
func (of *ObjectFrame) SetObject(ob FrameTarget) {
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
	of.bgDraw.Color = pixel.ToRGBA(mainColor)
	of.bgDraw.Push(of.DrawArea().Min)
	of.bgDraw.Push(of.DrawArea().Max)
	of.bgDraw.Rectangle(0)
	of.bgDraw.Draw(t)
}
