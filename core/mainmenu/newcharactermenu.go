/*
 * newcharactermenu.go
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

package mainmenu

import (
	"fmt"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"

	"github.com/isangeles/flame/core/data/text/lang"

	"github.com/isangeles/mural/core/mtk"
)

// NewCharacterMenu struct represents new game character
// creation screen.
type NewCharacterMenu struct {
	title      *text.Text
	backButton *mtk.Button
	opened    bool
}

// newNewCharacterMenu creates new character creation menu.
func newNewCharacterMenu() (*NewCharacterMenu, error) {
	ncm := new(NewCharacterMenu)
	// Title.
	font := mtk.MainFont(mtk.SIZE_BIG)
	atlas := mtk.Atlas(&font)
	ncm.title = text.New(pixel.V(0, 0), atlas)
	fmt.Fprint(ncm.title, lang.Text("gui", "newchar_menu_title"))
	// Button.
	ncm.backButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		colornames.Red, lang.Text("gui", "back_b_label"), "")
	
	return ncm, nil
}

// Draw draws all menu elements in specified window.
func (ncm *NewCharacterMenu) Draw(win *pixelgl.Window) {
	// Title.
	titlePos := pixel.V(win.Bounds().Center().X,
		win.Bounds().Max.Y - ncm.title.Bounds().Size().Y)
	ncm.title.Draw(win, pixel.IM.Moved(titlePos))
	// Buttons.
	ncm.backButton.Draw(win, pixel.IM.Moved(mtk.PosBL(ncm.backButton.Frame(),
		win.Bounds().Min)))
}

// Update updates all menu elements.
func (ncm *NewCharacterMenu) Update(win *pixelgl.Window) {
	if ncm.Opened() {
		ncm.backButton.Update(win)
	}
}

// Show toggles menu visibility.
func (ncm *NewCharacterMenu) Show(show bool) {
	ncm.opened = show
}

// Opened checks whether menu is open.
func (ncm *NewCharacterMenu) Opened() bool {
	return ncm.opened
}

// Sets scpecified function as back button on-click callback
// function.
func (ncm *NewCharacterMenu) SetOnBackFunc(f func(b *mtk.Button)) {
	ncm.backButton.SetOnClickFunc(f)
}
