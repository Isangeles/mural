/*
 * loadingscreen.go
 *
 * Copyright 2018 Dariusz Sikora <dev@isangeles.pl>
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
	
	"github.com/isangeles/mural/core/mtk"
)

// Struct for HUD loading screen.
type LoadingScreen struct {
	info *mtk.Textbox
}

// newLoadingScreen returns new HUD loadin screen.
func newLoadingScreen(hud *HUD) *LoadingScreen {
	ls := new(LoadingScreen)
	ls.info = mtk.NewTextbox(pixel.V(0, 0),
		mtk.SIZE_MEDIUM, main_color)
	return ls
}

// Draw draws loading screen.
func (ls *LoadingScreen) Draw(win *mtk.Window) {
	infoBounds := mtk.SIZE_MEDIUM.MessageWindowSize()
	dw := mtk.MatrixToDrawArea(mtk.Matrix().Moved(win.Bounds().Center()),
		infoBounds)
	ls.info.Draw(dw, win)
}

// Update updates loading screen.
func (ls *LoadingScreen) Update(win *mtk.Window) {
	ls.info.Update(win)
}

// SetLoadInfo sets specified text as current load info text.
func (ls *LoadingScreen) SetLoadInfo(text string) {
	ls.info.Clear()
	ls.info.Add(text)
}
