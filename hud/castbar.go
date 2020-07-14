/*
 * castbar.go
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

	"github.com/isangeles/mtk"
)

// Struct for HUD cast bar.
// TODO: bar background.
type CastBar struct {
	hud   *HUD
	bar   *mtk.ProgressBar
}

// newCastBar creates new HUD cast bar.
func newCastBar(hud *HUD) *CastBar {
	cb := new(CastBar)
	cb.hud = hud
	cb.bar = mtk.NewProgressBar(mtk.SizeMedium, accentColor)
	return cb
}

// Draw draws cast bar.
func (cb *CastBar) Draw(win *mtk.Window, matrix pixel.Matrix) {
	cb.bar.Draw(win.Window, matrix)
}

// Update updates cast bar.
func (cb *CastBar) Update(win *mtk.Window) {
	pc := cb.hud.ActivePlayer()
	if pc == nil {
		return
	}
	if pc.Casted() != nil {
		cb.bar.SetMax(int(pc.Casted().UseAction().CastMax()))
		cb.bar.SetValue(int(pc.Casted().UseAction().Cast()))
	}
	cb.bar.Update(win)
}
