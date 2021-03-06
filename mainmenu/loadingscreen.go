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

package mainmenu

import (
	"fmt"
	
	"github.com/isangeles/mtk"
)

// Struct for main menu loading screen.
type LoadingScreen struct {
	mainmenu *MainMenu
	info     *mtk.Text
}

// newLoadingScreen creates new main menu
// loading screen.
func newLoadingScreen(mainmenu *MainMenu) *LoadingScreen {
	ls := new(LoadingScreen)
	ls.mainmenu = mainmenu
	infoParams := mtk.Params{
		SizeRaw:     mtk.SizeMedium.MessageWindowSize(),
		FontSize:    mtk.SizeMedium,
		MainColor:   mainColor,
		AccentColor: accentColor,
	}
	ls.info = mtk.NewText(infoParams)
	return ls
}

// Draw draws loading screen.
func (ls *LoadingScreen) Draw(win *mtk.Window) {
	infoPos := win.Bounds().Center()
	ls.info.Draw(win, mtk.Matrix().Moved(infoPos))
}

// Update updates loading screen.
func (ls *LoadingScreen) Update(win *mtk.Window) {
}

// SetLoadInfo sets specified text as current load info text.
func (ls *LoadingScreen) SetLoadInfo(text string) {
	fmt.Sprintf("load screen: set text: %s\n", text)
	ls.info.SetText(text)
}
