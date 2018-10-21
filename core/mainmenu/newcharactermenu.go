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

	"github.com/isangeles/flame/core/data/text/lang"

	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/core/mtk"
)

// NewCharacterMenu struct represents new game character
// creation screen.
type NewCharacterMenu struct {
	mainmenu   *MainMenu  
	title      *mtk.Text
	nameEdit   *mtk.Textedit
	pointsBox  *mtk.Textbox
	strSwitch  *mtk.Switch
	backButton *mtk.Button
	opened     bool
	// Character.
	attrPoints, attrPointsMax int
}

// newNewCharacterMenu creates new character creation menu.
func newNewCharacterMenu(mainmenu *MainMenu) (*NewCharacterMenu, error) {
	ncm := new(NewCharacterMenu)
	ncm.mainmenu = mainmenu
	// Character.
	ncm.attrPointsMax = 5
	ncm.attrPoints = ncm.attrPointsMax
	// Title.
	ncm.title = mtk.NewText(lang.Text("gui", "newchar_menu_title"),
		mtk.SIZE_BIG, 0)
	// Text fields.
	ncm.nameEdit = mtk.NewTextedit(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_name_edit_label"))
	ncm.pointsBox = mtk.NewTextbox(mtk.SIZE_MEDIUM, main_color)
	// Buttons & switches.
	ncm.strSwitch = mtk.NewIntSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_str_switch_label"), 0, ncm.attrPointsMax)
	ncm.strSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.backButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		colornames.Red, lang.Text("gui", "back_b_label"), "")
	ncm.backButton.SetOnClickFunc(ncm.onBackButtonClicked)
	
	return ncm, nil
}

// Draw draws all menu elements in specified window.
func (ncm *NewCharacterMenu) Draw(win *pixelgl.Window) {
	// Title.
	titlePos := pixel.V(win.Bounds().Center().X,
		win.Bounds().Max.Y - ncm.title.Bounds().Size().Y)
	ncm.title.Draw(win, pixel.IM.Moved(titlePos))
	// Text fields.
	ncm.nameEdit.Draw(pixel.R(titlePos.X, titlePos.Y - mtk.ConvSize(30),
		titlePos.X + mtk.ConvSize(150), titlePos.Y - mtk.ConvSize(50)), win)
	ncm.pointsBox.Draw(pixel.R(win.Bounds().Min.X + mtk.ConvSize(10),
		win.Bounds().Center().Y, win.Bounds().Min.X + mtk.ConvSize(50),
		win.Bounds().Center().Y + mtk.ConvSize(40)), win)
	// Buttons && switches.
	ncm.strSwitch.Draw(win, pixel.IM.Moved(mtk.RightOf(ncm.pointsBox.DrawArea(),
		ncm.strSwitch.Frame(), 0)))
	ncm.backButton.Draw(win, pixel.IM.Moved(mtk.PosBL(ncm.backButton.Frame(),
		win.Bounds().Min)))
}

// Update updates all menu elements.
func (ncm *NewCharacterMenu) Update(win *pixelgl.Window) {
	if ncm.Opened() {
		ncm.nameEdit.Update(win)
		ncm.backButton.Update(win)
		ncm.pointsBox.Update(win)
		ncm.strSwitch.Update(win)
		ncm.pointsBox.InsertText([]string{fmt.Sprintf("%d", ncm.attrPoints)})
	}
}

// Show toggles menu visibility.
func (ncm *NewCharacterMenu) Show(show bool) {
	ncm.opened = show
	ncm.nameEdit.Focus(show)
}

// Opened checks whether menu is open.
func (ncm *NewCharacterMenu) Opened() bool {
	return ncm.opened
}

// Triggered after back button clicked.
func (ncm *NewCharacterMenu) onBackButtonClicked(b *mtk.Button) {
	ncm.mainmenu.OpenMenu()
}

// Triggered after strength value switch changed.
func (ncm *NewCharacterMenu) onAttrSwitchChange(s *mtk.Switch,
	old, new *mtk.SwitchValue) {
	str, ok := ncm.strSwitch.Value().Value.(int)
	if !ok {
		log.Err.Print("new_char_menu:fail_to_retrieve_str_switch_value")
		return
	}
	pts := ncm.attrPointsMax
	pts -= str
	if pts >= 0 && pts <= ncm.attrPointsMax {
		ncm.attrPoints = pts
	} else {
		s.SetIndex(s.Find(old.Value))
	}
	ncm.pointsBox.InsertText([]string{fmt.Sprint(ncm.attrPoints)})
}
