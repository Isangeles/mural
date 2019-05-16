/*
 * castbar.go
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

	"github.com/isangeles/flame/core/module/object/skill"
	
	"github.com/isangeles/mtk"
)

// Struct for HUD cast bar.
// TODO: bar background.
type CastBar struct {
	hud   *HUD
	owner skill.SkillUser
	bar   *mtk.ProgressBar
}

// newCastBar creates new HUD cast bar.
func newCastBar(hud *HUD) *CastBar {
	cb := new(CastBar)
	cb.hud = hud
	cb.bar = mtk.NewProgressBar(mtk.SIZE_MEDIUM, accent_color)
	return cb
}

// Draw draws cast bar.
func (cb *CastBar) Draw(win *mtk.Window, matrix pixel.Matrix) {
	cb.bar.Draw(win.Window, matrix)
}

// Update updates cast bar.
func (cb *CastBar) Update(win *mtk.Window) {
	if cb.owner == nil {
		return
	}
	for _, s := range cb.owner.Skills() {
		if s.Casting() {
			cb.bar.SetMax(int(s.CastTimeMax()))
			cb.bar.SetValue(int(s.CastTime()))
		}
	}
	cb.bar.Update(win)
}

// SetOwner sets specified skill user as cast
// bar owner.
func (cb *CastBar) SetOwner(o skill.SkillUser) {
	cb.owner = o
}
