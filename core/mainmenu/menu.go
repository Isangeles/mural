/*
 * menu.go
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

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/mural/core"
	"github.com/isangeles/mural/core/data"
)

// Menu struct represents main menu screen
// with buttons to other menus.
type Menu struct {
	title    *text.Text
	exitB    *core.Button
	CloseReq bool
}

// newMenu returns new menu.
func newMenu() (*Menu, error) {
	m := new(Menu)
	// Title.
	font := data.MainFontBig()
	atlas := text.NewAtlas(font, text.ASCII)
	m.title = text.New(pixel.V(0, 0), atlas)
	fmt.Fprint(m.title, flame.Mod().Name())
	// Exit button.
	exitBBg, err := data.Picture("buttonS.png")
	if err != nil {
		return nil, err
	}
	m.exitB = core.NewButton(exitBBg, lang.Text("gui", "exit_b_label"))

	return m, nil
}

// Draw draws all menu elements in specified
// window.
func (m *Menu) Draw(win *pixelgl.Window) {
	titlePos := pixel.V(win.Bounds().Center().X, win.Bounds().Max.Y - m.title.Bounds().Size().Y)
	m.title.Draw(win, pixel.IM.Moved(titlePos))
	m.exitB.Draw(win, pixel.IM.Moved(pixel.V(titlePos.X, titlePos.Y - m.exitB.DrawArea().Size().Y)))
}

// Update updates all menu elements.
func (m *Menu) Update(win *pixelgl.Window) {
	m.exitB.Update(win)

	if win.JustReleased(pixelgl.MouseButtonLeft) {
		if m.exitB.ContainsPosition(win.MousePosition()) {
			win.SetClosed(true)
		}
	}
}
