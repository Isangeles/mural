/*
 * settings.go
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

	//"github.com/isangeles/flame"
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/mural/core"
	"github.com/isangeles/mural/core/data"
)

// Settings struct represents main menu
// settings screen.
type Settings struct {
	title *text.Text
	backB *core.Button
	open  bool
}

// newSettings returns new settings screen
// instance.
func newSettings() (*Settings, error) {
	s := new(Settings)
	// Title.
	font := data.MainFontBig()
	atlas := text.NewAtlas(font, text.ASCII)
	s.title = text.New(pixel.V(0, 0), atlas)
	fmt.Fprintf(s.title, lang.Text("gui", "settings_menu_title"))
	// Buttons.
	s.backB = core.NewButtonDraw(core.SIZE_SMALL, lang.Text("gui", "back_b_label"))

	return s, nil
}

// Open checks whether menu should be drawn or not.
func (s *Settings) Open() bool {
	return s.open
}

// Show toggles menu visibility.
func (s *Settings) Show(show bool) {
	s.open = show
}

// Draw draws all menu elements.
func (s *Settings) Draw(win *pixelgl.Window) {
	titlePos :=pixel.V(win.Bounds().Center().X,
		win.Bounds().Max.Y - s.title.Bounds().Size().Y)
	s.title.Draw(win, pixel.IM.Moved(titlePos))
	s.backB.Draw(win, pixel.IM.Moved(pixel.V(titlePos.X,
		titlePos.Y - s.backB.Frame().Size().Y)))
}

// Update updates all menu elements.
func (s *Settings) Update(win *pixelgl.Window) {
	s.backB.Update(win)

	if s.backB.ContainsPosition(win.MousePosition()) {
		s.Show(false)
	}
}
