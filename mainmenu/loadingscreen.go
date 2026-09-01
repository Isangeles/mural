/*
 * loadingscreen.go
 *
 * Copyright 2018-2026 Dariusz Sikora <ds@isangeles.dev>
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
	"sync"

	"github.com/isangeles/mtk"
)

// Struct for main menu loading screen.
type LoadingScreen struct {
	mainmenu *MainMenu
	info     *mtk.Text
	loadInfo string
	infoText string
	mutex    sync.Mutex
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
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	if ls.infoText != ls.loadInfo {
		ls.info.SetText(ls.loadInfo)
		ls.infoText = ls.loadInfo
	}
}

// SetLoadInfo sets specified text as current load info text.
// Text is set on the info display during the next update,
// so this function is safe to call outside the main thread.
func (ls *LoadingScreen) SetLoadInfo(text string) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	ls.loadInfo = text
}
