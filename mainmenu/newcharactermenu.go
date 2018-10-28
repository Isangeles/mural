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
	"time"
	"math/rand"
	
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/core/data/text/lang"

	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/data"
)

// NewCharacterMenu struct represents new game character
// creation screen.
type NewCharacterMenu struct {
	mainmenu   *MainMenu  
	title      *mtk.Text
	nameEdit   *mtk.Textedit
	faceSwitch *mtk.Switch
	pointsBox  *mtk.Textbox
	strSwitch  *mtk.Switch
	conSwitch  *mtk.Switch
	dexSwitch  *mtk.Switch
	intSwitch  *mtk.Switch
	wisSwitch  *mtk.Switch
	doneButton *mtk.Button
	backButton *mtk.Button
	rollButton *mtk.Button
	opened     bool
	rng        *rand.Rand
	// Character.
	attrPoints, attrPointsMax int
}

// newNewCharacterMenu creates new character creation menu.
func newNewCharacterMenu(mainmenu *MainMenu) (*NewCharacterMenu, error) {
	ncm := new(NewCharacterMenu)
	ncm.mainmenu = mainmenu
	rngSrc := rand.NewSource(time.Now().UnixNano())
	ncm.rng = rand.New(rngSrc)
	// Character.
	ncm.attrPointsMax = ncm.rollPoints()
	ncm.attrPoints = ncm.attrPointsMax
	// Title.
	ncm.title = mtk.NewText(lang.Text("gui", "newchar_menu_title"),
		mtk.SIZE_BIG, 0)
	// Text fields.
	ncm.nameEdit = mtk.NewTextedit(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_name_edit_label"))
	ncm.pointsBox = mtk.NewTextbox(mtk.SIZE_MEDIUM, main_color)
	// Switches.
	faces, err := data.PlayablePortraits()
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_player_portraits:%v", err)
	}
	ncm.faceSwitch = mtk.NewPictureSwitch(mtk.SIZE_BIG, main_color,
		lang.Text("gui", "newchar_face_switch_label"), faces) 
	ncm.strSwitch = mtk.NewIntSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_str_switch_label"), 0, ncm.attrPointsMax)
	ncm.strSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.conSwitch = mtk.NewIntSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_con_switch_label"), 0, ncm.attrPointsMax)
	ncm.conSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.dexSwitch = mtk.NewIntSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_dex_switch_label"), 0, ncm.attrPointsMax)
	ncm.dexSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.intSwitch = mtk.NewIntSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_int_switch_label"), 0, ncm.attrPointsMax) 
	ncm.intSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.wisSwitch = mtk.NewIntSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_wis_switch_label"), 0, ncm.attrPointsMax)
	ncm.wisSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	// Buttons.
	ncm.doneButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		colornames.Red, lang.Text("gui", "done_b_label"), "")
	ncm.doneButton.SetOnClickFunc(ncm.onDoneButtonClicked)
	ncm.backButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		colornames.Red, lang.Text("gui", "back_b_label"), "")
	ncm.backButton.SetOnClickFunc(ncm.onBackButtonClicked)
	ncm.rollButton = mtk.NewButton(mtk.SIZE_SMALL, mtk.SHAPE_RECTANGLE,
		colornames.Red, lang.Text("gui", "newchar_roll_b_label"),
		lang.Text("gui", "newchar_roll_b_info"))
	ncm.rollButton.SetOnClickFunc(ncm.onRollButtonClicked)
	
	return ncm, nil
}

// Draw draws all menu elements in specified window.
func (ncm *NewCharacterMenu) Draw(win *mtk.Window) {
	// Title.
	titlePos := pixel.V(win.Bounds().Center().X,
		win.Bounds().Max.Y - ncm.title.Bounds().Size().Y)
	ncm.title.Draw(win, mtk.Matrix().Moved(titlePos))
	// Text fields.
	ncm.nameEdit.Draw(pixel.R(titlePos.X, titlePos.Y - mtk.ConvSize(30),
		titlePos.X + mtk.ConvSize(150), titlePos.Y - mtk.ConvSize(50)), win.Window)
	ncm.pointsBox.Draw(pixel.R(win.Bounds().Min.X + mtk.ConvSize(90),
		win.Bounds().Center().Y - mtk.ConvSize(40), win.Bounds().Min.X + mtk.ConvSize(140),
		win.Bounds().Center().Y + mtk.ConvSize(40)), win.Window)
	// Switches.
	ncm.faceSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.TopOf(ncm.pointsBox.DrawArea(),
		ncm.faceSwitch.Frame(), 100)))
	ncm.strSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(ncm.pointsBox.DrawArea(),
		ncm.strSwitch.Frame(), 5)))
	ncm.conSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(ncm.strSwitch.DrawArea(),
		ncm.conSwitch.Frame(), 15)))
	ncm.dexSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(ncm.conSwitch.DrawArea(),
		ncm.dexSwitch.Frame(), 15)))
	ncm.intSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(ncm.dexSwitch.DrawArea(),
		ncm.intSwitch.Frame(), 15)))
	ncm.wisSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(ncm.intSwitch.DrawArea(),
		ncm.wisSwitch.Frame(), 15)))
	// Buttons.
	ncm.doneButton.Draw(win.Window, mtk.Matrix().Moved(mtk.PosBR(ncm.doneButton.Frame(),
		pixel.V(win.Bounds().Max.X, win.Bounds().Min.Y))))
	ncm.backButton.Draw(win.Window, mtk.Matrix().Moved(mtk.PosBL(ncm.backButton.Frame(),
		win.Bounds().Min)))
	ncm.rollButton.Draw(win.Window, mtk.Matrix().Moved(mtk.BottomOf(ncm.pointsBox.DrawArea(),
		ncm.rollButton.Frame(), 5)))
}

// Update updates all menu elements.
func (ncm *NewCharacterMenu) Update(win *mtk.Window) {
	if ncm.Opened() {
		ncm.nameEdit.Update(win.Window)
		ncm.doneButton.Update(win.Window)
		ncm.backButton.Update(win.Window)
		ncm.rollButton.Update(win.Window)
		ncm.pointsBox.Update(win.Window)
		ncm.faceSwitch.Update(win.Window)
		ncm.strSwitch.Update(win.Window)
		ncm.conSwitch.Update(win.Window)
		ncm.dexSwitch.Update(win.Window)
		ncm.intSwitch.Update(win.Window)
		ncm.wisSwitch.Update(win.Window)
		ncm.updatePoints()
	}
}

// rollPoints draws random amount of attribute points for new character
// from range specified in Flame config.
func (ncm *NewCharacterMenu) rollPoints() int {
	return ncm.rng.Intn(flame.NewCharAttrMax() -
		flame.NewCharAttrMin()) + flame.NewCharAttrMin()
}

// updatePoints updates points box value.
func (ncm *NewCharacterMenu) updatePoints() {
	ncm.pointsBox.InsertText([]string{fmt.Sprintf("%d", ncm.attrPoints)})
}

// createChar creates new game character.
func (ncm *NewCharacterMenu) createChar() {
	// TODO: new character create.
}

// Show toggles menu visibility.
func (ncm *NewCharacterMenu) Show(show bool) {
	ncm.opened = show
}

// Opened checks whether menu is open.
func (ncm *NewCharacterMenu) Opened() bool {
	return ncm.opened
}

// Triggered after back button clicked.
func (ncm *NewCharacterMenu) onBackButtonClicked(b *mtk.Button) {
	ncm.mainmenu.OpenMenu()
}

// Triggered after done button clicked.
func (ncm *NewCharacterMenu) onDoneButtonClicked(b *mtk.Button) {
	ncm.createChar()
}

// Triggered after roll button clicked.
func (ncm *NewCharacterMenu) onRollButtonClicked(b *mtk.Button) {
	ncm.strSwitch.Reset()
	ncm.conSwitch.Reset()
	ncm.dexSwitch.Reset()
	ncm.intSwitch.Reset()
	ncm.wisSwitch.Reset()
	ncm.attrPointsMax = ncm.rollPoints()
	ncm.attrPoints = ncm.attrPointsMax
	ncm.updatePoints()
}

// Triggered after any attribute switch value changed.
func (ncm *NewCharacterMenu) onAttrSwitchChange(s *mtk.Switch,
	old, new *mtk.SwitchValue) {
	str, ok := ncm.strSwitch.Value().Value.(int)
	if !ok {
		log.Err.Print("new_char_menu:fail_to_retrieve_str_switch_value")
		return
	}
	con, ok := ncm.conSwitch.Value().Value.(int)
	if !ok {
		log.Err.Print("new_char_menu:fail_to_retrieve_con_switch_value")
		return
	}
	dex, ok := ncm.dexSwitch.Value().Value.(int)
	if !ok {
		log.Err.Print("new_char_menu:fail_to_retrieve_con_switch_value")
		return
	}
	inte, ok := ncm.intSwitch.Value().Value.(int)
	if !ok {
		log.Err.Print("new_char_menu:fail_to_retrieve_int_switch_value")
		return
	}
	wis, ok := ncm.wisSwitch.Value().Value.(int)
	if !ok {
		log.Err.Print("new_char_menu:fail_to_retrieve_wis_switch_value")
		return
	}
	pts := ncm.attrPointsMax
	pts -= str + con + dex + inte + wis
	if pts >= 0 && pts <= ncm.attrPointsMax {
		ncm.attrPoints = pts
	} else {
		s.SetIndex(s.Find(old.Value))
	}
	ncm.updatePoints()
}
