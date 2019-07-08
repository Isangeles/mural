/*
 * loadingscreen.go
 *
 * Copyright 2018-2019 Dariusz Sikora <dev@isangeles.pl>
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
	"github.com/isangeles/mtk"
)

// Struct for HUD loading screen.
type LoadingScreen struct {
	hud  *HUD
	info *mtk.Textbox
}

// newLoadingScreen returns new HUD loading
// screen.
func newLoadingScreen(hud *HUD) *LoadingScreen {
	ls := new(LoadingScreen)
	ls.hud = hud
	infoParams := mtk.Params{
		FontSize:    mtk.SizeMini,
		MainColor:   main_color,
		AccentColor: accent_color,
	}
	ls.info = mtk.NewTextbox(infoParams)
	return ls
}

// Draw draws loading screen.
func (ls *LoadingScreen) Draw(win *mtk.Window) {
	infoSize := mtk.SizeMedium.MessageWindowSize()
	infoPos := win.Bounds().Center()
	ls.info.SetSize(infoSize)
	ls.info.Draw(win, mtk.Matrix().Moved(infoPos))
}

// Update updates loading screen.
func (ls *LoadingScreen) Update(win *mtk.Window) {
	ls.info.Update(win)
}

// SetLoadInfo sets specified text as current load info text.
func (ls *LoadingScreen) SetLoadInfo(text string) {
	ls.info.SetText(text)
}
