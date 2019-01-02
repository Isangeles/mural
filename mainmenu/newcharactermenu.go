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
	"math/rand"
	"strings"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/object/character"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/core/objects"
)

var (
	new_char_base_id = "player_"
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
	sexSwitch  *mtk.Switch
	raceSwitch *mtk.Switch
	aliSwitch  *mtk.Switch
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
	// Title.
	ncm.title = mtk.NewText(lang.Text("gui", "newchar_menu_title"),
		mtk.SIZE_BIG, 0)
	// Text fields.
	ncm.nameEdit = mtk.NewTextedit(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_name_edit_label"))
	ncm.pointsBox = mtk.NewTextbox(pixel.V(0, 0), mtk.SIZE_MEDIUM,
		main_color)
	// Portrait switch.
	faces, err := data.PlayablePortraits()
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_player_portraits:%v", err)
	}
	ncm.faceSwitch = mtk.NewSwitch(mtk.SIZE_BIG, main_color,
		lang.Text("gui", "newchar_face_switch_label"), "", nil)
	ncm.faceSwitch.SetPictureValues(faces)
	// Attributes switches.
	ncm.strSwitch = mtk.NewSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_str_switch_label"), "", nil)
	ncm.strSwitch.SetIntValues(0, 90)
	ncm.strSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.conSwitch = mtk.NewSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_con_switch_label"), "", nil)
	ncm.conSwitch.SetIntValues(0, 90)
	ncm.conSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.dexSwitch = mtk.NewSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_dex_switch_label"), "", nil)
	ncm.dexSwitch.SetIntValues(0, 90)
	ncm.dexSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.intSwitch = mtk.NewSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_int_switch_label"), "", nil)
	ncm.intSwitch.SetIntValues(0, 90)
	ncm.intSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.wisSwitch = mtk.NewSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_wis_switch_label"), "", nil)
	ncm.wisSwitch.SetIntValues(0, 90)
	ncm.wisSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	// Gender & alligment switches.
	maleSwitchVal := mtk.SwitchValue{lang.Text("ui", "gender_male"),
		character.Male}
	femaleSwitchVal := mtk.SwitchValue{lang.Text("ui", "gender_female"),
		character.Female}
	gens := []mtk.SwitchValue{maleSwitchVal, femaleSwitchVal}
	ncm.sexSwitch = mtk.NewSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_sex_switch_label"), "", gens)
	// Race switch.
	raceNames := lang.Texts("ui", "race_human", "race_elf", "race_dwarf",
		"race_gnome")
	races := []mtk.SwitchValue{
		mtk.SwitchValue{raceNames[0], character.Human},
		mtk.SwitchValue{raceNames[1], character.Elf},
		mtk.SwitchValue{raceNames[2], character.Dwarf},
		mtk.SwitchValue{raceNames[3], character.Gnome},
	}
	ncm.raceSwitch = mtk.NewSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_race_switch_label"), "", races)
	// Alignment switch.
	aliNames := lang.Texts("ui", "ali_law_good", "ali_neu_good", "ali_cha_good",
		"ali_law_neutral", "ali_tru_neutral", "ali_cha_neutral",
		"ali_law_evil", "ali_neu_evil", "ali_cha_evil")
	alis := []mtk.SwitchValue{
		mtk.SwitchValue{aliNames[0], character.Lawful_good},
		mtk.SwitchValue{aliNames[1], character.Neutral_good},
		mtk.SwitchValue{aliNames[2], character.Chaotic_good},
		mtk.SwitchValue{aliNames[3], character.Lawful_neutral},
		mtk.SwitchValue{aliNames[4], character.True_neutral},
		mtk.SwitchValue{aliNames[5], character.Chaotic_neutral},
		mtk.SwitchValue{aliNames[6], character.Lawful_evil},
		mtk.SwitchValue{aliNames[7], character.Neutral_evil},
		mtk.SwitchValue{aliNames[8], character.Chaotic_evil},
	}
	ncm.aliSwitch = mtk.NewSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "newchar_ali_switch_label"), "", alis)
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
	// Character.
	ncm.rollPoints()

	return ncm, nil
}

// Draw draws all menu elements in specified window.
func (ncm *NewCharacterMenu) Draw(win *mtk.Window) {
	// Title.
	titlePos := pixel.V(win.Bounds().Center().X,
		win.Bounds().Max.Y-ncm.title.Bounds().Size().Y)
	ncm.title.Draw(win, mtk.Matrix().Moved(titlePos))
	// Text fields.
	ncm.nameEdit.Draw(pixel.R(titlePos.X, titlePos.Y-mtk.ConvSize(30),
		titlePos.X+mtk.ConvSize(150), titlePos.Y-mtk.ConvSize(50)),
		win.Window)
	ncm.pointsBox.Draw(pixel.R(win.Bounds().Min.X+mtk.ConvSize(90),
		win.Bounds().Center().Y-mtk.ConvSize(40),
		win.Bounds().Min.X+mtk.ConvSize(140),
		win.Bounds().Center().Y+mtk.ConvSize(40)), win.Window)
	// Switches.
	ncm.faceSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.TopOf(
		ncm.pointsBox.DrawArea(), ncm.faceSwitch.Frame(), 100)))
	ncm.strSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(
		ncm.pointsBox.DrawArea(), ncm.strSwitch.Frame(), 5)))
	ncm.conSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(
		ncm.strSwitch.DrawArea(), ncm.conSwitch.Frame(), 15)))
	ncm.dexSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(
		ncm.conSwitch.DrawArea(), ncm.dexSwitch.Frame(), 15)))
	ncm.intSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(
		ncm.dexSwitch.DrawArea(), ncm.intSwitch.Frame(), 15)))
	ncm.wisSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(
		ncm.intSwitch.DrawArea(), ncm.wisSwitch.Frame(), 15)))
	ncm.sexSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(
		ncm.wisSwitch.DrawArea(), ncm.sexSwitch.Frame(), 30)))
	ncm.raceSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.BottomOf(
		ncm.sexSwitch.DrawArea(), ncm.raceSwitch.Frame(), 10)))
	ncm.aliSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.BottomOf(
		ncm.raceSwitch.DrawArea(), ncm.aliSwitch.Frame(), 10)))
	// Buttons.
	ncm.doneButton.Draw(win.Window, mtk.Matrix().Moved(mtk.PosBR(
		ncm.doneButton.Frame(), pixel.V(win.Bounds().Max.X,
			win.Bounds().Min.Y))))
	ncm.backButton.Draw(win.Window, mtk.Matrix().Moved(mtk.PosBL(
		ncm.backButton.Frame(), win.Bounds().Min)))
	ncm.rollButton.Draw(win.Window, mtk.Matrix().Moved(mtk.BottomOf(
		ncm.pointsBox.DrawArea(), ncm.rollButton.Frame(), 5)))
}

// Update updates all menu elements.
func (ncm *NewCharacterMenu) Update(win *mtk.Window) {
	ncm.nameEdit.Update(win)
	ncm.doneButton.Update(win)
	ncm.backButton.Update(win)
	ncm.rollButton.Update(win)
	ncm.pointsBox.Update(win)
	ncm.faceSwitch.Update(win)
	ncm.strSwitch.Update(win)
	ncm.conSwitch.Update(win)
	ncm.dexSwitch.Update(win)
	ncm.intSwitch.Update(win)
	ncm.wisSwitch.Update(win)
	ncm.sexSwitch.Update(win)
	ncm.raceSwitch.Update(win)
	ncm.aliSwitch.Update(win)
	ncm.updatePoints()
	if ncm.canCreate() {
		ncm.doneButton.Active(false)
	} else {
		ncm.doneButton.Active(true)
	}
}

// rollPoints draws random amount of attribute points for new character
// from range specified in Flame config.
func (ncm *NewCharacterMenu) rollPoints() {
	ncm.strSwitch.Reset()
	ncm.conSwitch.Reset()
	ncm.dexSwitch.Reset()
	ncm.intSwitch.Reset()
	ncm.wisSwitch.Reset()
	ncm.attrPointsMax = ncm.rng.Intn(flame.Mod().NewcharAttrsMax()-
		flame.Mod().NewcharAttrsMin()) + flame.Mod().NewcharAttrsMin()
	ncm.attrPoints = ncm.attrPointsMax
	ncm.updatePoints()
}

// Show toggles menu visibility.
func (ncm *NewCharacterMenu) Show(show bool) {
	ncm.opened = show
}

// Opened checks whether menu is open.
func (ncm *NewCharacterMenu) Opened() bool {
	return ncm.opened
}

// canCreate checks whether its possible to create new character.
func (ncm *NewCharacterMenu) canCreate() bool {
	return ncm.nameEdit.Text() == "" || ncm.attrPoints > 0
}

// updatePoints updates points box value.
func (ncm *NewCharacterMenu) updatePoints() {
	ncm.pointsBox.InsertText([]string{fmt.Sprintf("%d", ncm.attrPoints)})
}

// createChar creates new game character.
func (ncm *NewCharacterMenu) createChar() (*character.Character, error) {
	name := ncm.nameEdit.Text()
	str, err := ncm.strSwitch.Value().IntValue()
	if err != nil {
		return nil, err
	}
	con, err := ncm.conSwitch.Value().IntValue()
	if err != nil {
		return nil, err
	}
	dex, err := ncm.dexSwitch.Value().IntValue()
	if err != nil {
		return nil, err
	}
	inte, err := ncm.intSwitch.Value().IntValue()
	if err != nil {
		return nil, err
	}
	wis, err := ncm.wisSwitch.Value().IntValue()
	if err != nil {
		return nil, err
	}
	attrs := character.Attributes{str, con, dex, inte, wis}
	gender, ok := ncm.sexSwitch.Value().Value.(character.Gender)
	if !ok {
		return nil, fmt.Errorf("fail_to_retrive_gender")
	}
	race, ok := ncm.raceSwitch.Value().Value.(character.Race)
	if !ok {
		return nil, fmt.Errorf("fail_to_retrieve_race")
	}
	alignment, ok := ncm.aliSwitch.Value().Value.(character.Alignment)
	if !ok {
		return nil, fmt.Errorf("fail_to_retrieve_alignment")
	}
	id := new_char_base_id + strings.ToLower(name)
	char := character.NewCharacter(id, name, 1, gender, race,
		character.Friendly, character.NewGuild("none"), attrs, alignment)
	return char, nil
}

// Triggered after back button clicked.
func (ncm *NewCharacterMenu) onBackButtonClicked(b *mtk.Button) {
	ncm.mainmenu.OpenMenu()
}

// Triggered after done button clicked.
func (ncm *NewCharacterMenu) onDoneButtonClicked(b *mtk.Button) {
	char, err := ncm.createChar()
	if err != nil {
		log.Err.Printf("newchar_menu:fail_to_create_character:%v", err)
		return
	}
	ssHeadName := "m-head-black-1222211-80x90.png"
	ssTorsoName := "m-cloth-1222211-80x90.png"
	if char.Gender() == character.Female {
		ssTorsoName = "f-cloth-1222211-80x90.png"
		ssHeadName = "f-head-black-1222211-80x90.png"
	}
	ssHeadPic, err := data.AvatarSpritesheet(ssHeadName)
	if err != nil {
		log.Err.Printf("newchar_menu:fail_to_retrieve_head_spritesheet_picture:%v",
			err)
		return
	}
	ssTorsoPic, err := data.AvatarSpritesheet(ssTorsoName)
	if err != nil {
		log.Err.Printf("newchar_menu:fail_to_retrieve_torso_spritesheet_picture:%v",
			err)
		return
	}
	portraitName, err := ncm.faceSwitch.Value().TextValue()
	portraitPic, err := data.AvatarPortrait(portraitName)
	if err != nil {
		log.Err.Printf("newchar_menu:fail_to_retrieve_portrait_picture:%v",
			err)
		return
	}
	av, err := objects.NewAvatar(char, portraitPic,
		ssHeadPic, ssTorsoPic, portraitName,
		ssHeadName, ssTorsoName)
	if err != nil {
		log.Err.Printf("newchar_menu:fail_to_create_avatar:%v",
			err)
		return
	}
	ncm.mainmenu.AddPlayableChar(av)
	msg := mtk.NewMessageWindow(mtk.SIZE_SMALL,
		lang.Text("gui", "newchar_create_msg"))
	ncm.mainmenu.ShowMessage(msg)
	ncm.mainmenu.OpenMenu()
}

// Triggered after roll button clicked.
func (ncm *NewCharacterMenu) onRollButtonClicked(b *mtk.Button) {
	ncm.rollPoints()
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
