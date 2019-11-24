/*
 * newcharactermenu.go
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
	"math/rand"
	"strings"
	"time"

	"github.com/faiface/pixel"

	flameconf "github.com/isangeles/flame/config"
	flamedata "github.com/isangeles/flame/core/data"
	flameres "github.com/isangeles/flame/core/data/res"
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/character"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

var (
	newCharIDFrom   = `player_%s` // player_[name]
	newCharAttrsMin = config.CharAttrsMin
	newCharAttrsMax = config.CharAttrsMax
)

// NewCharacterMenu struct represents new game character
// creation screen.
type NewCharacterMenu struct {
	mainmenu      *MainMenu
	title         *mtk.Text
	nameLabel     *mtk.Text
	nameEdit      *mtk.Textedit
	faceSwitch    *mtk.Switch
	pointsBox     *mtk.Text
	strSwitch     *mtk.Switch
	conSwitch     *mtk.Switch
	dexSwitch     *mtk.Switch
	intSwitch     *mtk.Switch
	wisSwitch     *mtk.Switch
	sexSwitch     *mtk.Switch
	raceSwitch    *mtk.Switch
	aliSwitch     *mtk.Switch
	doneButton    *mtk.Button
	backButton    *mtk.Button
	rollButton    *mtk.Button
	opened        bool
	rng           *rand.Rand
	attrPoints    int
	attrPointsMax int
}

// newNewCharacterMenu creates new character creation menu.
func newNewCharacterMenu(mainmenu *MainMenu) *NewCharacterMenu {
	ncm := new(NewCharacterMenu)
	ncm.mainmenu = mainmenu
	rngSrc := rand.NewSource(time.Now().UnixNano())
	ncm.rng = rand.New(rngSrc)
	// Title.
	titleParams := mtk.Params{
		FontSize: mtk.SizeBig,
	}
	ncm.title = mtk.NewText(titleParams)
	ncm.title.SetText(lang.Text("gui", "newchar_menu_title"))
	// Name Edit.
	labelParams := mtk.Params{
		FontSize: mtk.SizeMedium,
	}
	ncm.nameLabel = mtk.NewText(labelParams)
	nameLabelText := lang.TextDir(flameconf.LangPath(), "newchar_name_edit_label")
	ncm.nameLabel.SetText(fmt.Sprintf("%s:", nameLabelText))
	ncm.nameEdit = mtk.NewTextedit(mtk.SizeMedium, mainColor)
	// Points box.
	pointsBoxSize := mtk.SizeMedium.ButtonSize(mtk.ShapeRectangle)
	pointsBoxParams := mtk.Params{
		SizeRaw:     pointsBoxSize,
		FontSize:    mtk.SizeBig,
	}
	ncm.pointsBox = mtk.NewText(pointsBoxParams)
	// Portrait switch.
	faceSwitchParams := mtk.Params{
		Size:      mtk.SizeBig,
		MainColor: mainColor,
	}
	ncm.faceSwitch = mtk.NewSwitch(faceSwitchParams)
	ncm.faceSwitch.SetLabel(lang.Text("gui", "newchar_face_switch_label"))
	faces, err := data.PlayablePortraits()
	if err != nil {
		log.Err.Printf("new_char_menu:fail_to_retrieve_player_portraits:%v", err)
	}
	if faces != nil {
		ncm.faceSwitch.SetPictureValues(faces)
	}
	// Attributes switches.
	attrSwitchParams := mtk.Params{
		Size:      mtk.SizeMedium,
		MainColor: mainColor,
	}
	ncm.strSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.strSwitch.SetLabel(lang.Text("gui", "newchar_str_switch_label"))
	ncm.strSwitch.SetIntValues(0, 90)
	ncm.strSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.conSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.conSwitch.SetLabel(lang.Text("gui", "newchar_con_switch_label"))
	ncm.conSwitch.SetIntValues(0, 90)
	ncm.conSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.dexSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.dexSwitch.SetLabel(lang.Text("gui", "newchar_dex_switch_label"))
	ncm.dexSwitch.SetIntValues(0, 90)
	ncm.dexSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.intSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.intSwitch.SetLabel(lang.Text("gui", "newchar_int_switch_label"))
	ncm.intSwitch.SetIntValues(0, 90)
	ncm.intSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.wisSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.wisSwitch.SetLabel(lang.Text("gui", "newchar_wis_switch_label"))
	ncm.wisSwitch.SetIntValues(0, 90)
	ncm.wisSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	// Gender & alligment switches.
	maleSwitchVal := mtk.SwitchValue{lang.Text("ui", "gender_male"),
		character.Male}
	femaleSwitchVal := mtk.SwitchValue{lang.Text("ui", "gender_female"),
		character.Female}
	gens := []mtk.SwitchValue{maleSwitchVal, femaleSwitchVal}
	ncm.sexSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.sexSwitch.SetLabel(lang.Text("gui", "newchar_sex_switch_label"))
	ncm.sexSwitch.SetValues(gens)
	// Race switch.
	raceNames := lang.Texts("ui", "race_human", "race_elf", "race_dwarf",
		"race_gnome")
	races := []mtk.SwitchValue{
		mtk.SwitchValue{raceNames["race_human"], character.Human},
		mtk.SwitchValue{raceNames["race_elf"], character.Elf},
		mtk.SwitchValue{raceNames["race_dwarf"], character.Dwarf},
		mtk.SwitchValue{raceNames["race_gnome"], character.Gnome},
	}
	ncm.raceSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.raceSwitch.SetLabel(lang.Text("gui", "newchar_race_switch_label"))
	ncm.raceSwitch.SetValues(races)
	// Alignment switch.
	aliNames := lang.Texts("ui", "ali_law_good", "ali_neu_good", "ali_cha_good",
		"ali_law_neutral", "ali_tru_neutral", "ali_cha_neutral",
		"ali_law_evil", "ali_neu_evil", "ali_cha_evil")
	alis := []mtk.SwitchValue{
		mtk.SwitchValue{aliNames["ali_law_good"], character.Lawful_good},
		mtk.SwitchValue{aliNames["ali_neu_good"], character.Neutral_good},
		mtk.SwitchValue{aliNames["ali_cha_good"], character.Chaotic_good},
		mtk.SwitchValue{aliNames["ali_law_neutral"], character.Lawful_neutral},
		mtk.SwitchValue{aliNames["ali_tru_neutral"], character.True_neutral},
		mtk.SwitchValue{aliNames["ali_cha_neutral"], character.Chaotic_neutral},
		mtk.SwitchValue{aliNames["ali_law_evil"], character.Lawful_evil},
		mtk.SwitchValue{aliNames["ali_neu_evil"], character.Neutral_evil},
		mtk.SwitchValue{aliNames["ali_cha_evil"], character.Chaotic_evil},
	}
	ncm.aliSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.aliSwitch.SetLabel(lang.Text("gui", "newchar_ali_switch_label"))
	ncm.aliSwitch.SetValues(alis)
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		FontSize:  mtk.SizeMedium,
		Shape:     mtk.ShapeRectangle,
		MainColor: accentColor,
	}
	ncm.doneButton = mtk.NewButton(buttonParams)
	ncm.doneButton.SetLabel(lang.Text("gui", "done_b_label"))
	ncm.doneButton.SetOnClickFunc(ncm.onDoneButtonClicked)
	ncm.backButton = mtk.NewButton(buttonParams)
	ncm.backButton.SetLabel(lang.Text("gui", "back_b_label"))
	ncm.backButton.SetOnClickFunc(ncm.onBackButtonClicked)
	ncm.rollButton = mtk.NewButton(buttonParams)
	ncm.rollButton.SetLabel(lang.Text("gui", "newchar_roll_b_label"))
	ncm.rollButton.SetInfo(lang.Text("gui", "newchar_roll_b_info"))
	ncm.rollButton.SetOnClickFunc(ncm.onRollButtonClicked)
	return ncm
}

// Draw draws all menu elements in specified window.
func (ncm *NewCharacterMenu) Draw(win *mtk.Window) {
	// Title.
	titlePos := pixel.V(win.Bounds().Center().X,
		win.Bounds().H()-ncm.title.Size().Y)
	ncm.title.Draw(win, mtk.Matrix().Moved(titlePos))
	// Name edit.
	nameLabelPos := mtk.BottomOf(ncm.title.DrawArea(), ncm.nameEdit.Size(), 10)
	ncm.nameLabel.Draw(win, mtk.Matrix().Moved(nameLabelPos))
	nameEditPos := mtk.BottomOf(ncm.nameLabel.DrawArea(), ncm.nameEdit.Size(), 10)
	nameEditSize := ncm.title.DrawArea().Size()
	ncm.nameEdit.SetSize(nameEditSize)
	ncm.nameEdit.Draw(win.Window, mtk.Matrix().Moved(nameEditPos))
	// Points box.
	pointsBoxPos := mtk.DrawPosCL(win.Bounds(), ncm.pointsBox.Size())
	pointsBoxPos.X += mtk.ConvSize(100)
	ncm.pointsBox.Draw(win.Window, mtk.Matrix().Moved(pointsBoxPos))
	// Switches.
	ncm.faceSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.TopOf(
		ncm.pointsBox.DrawArea(), ncm.faceSwitch.Size(), 100)))
	ncm.strSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(
		ncm.pointsBox.DrawArea(), ncm.strSwitch.Size(), 5)))
	ncm.conSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(
		ncm.strSwitch.DrawArea(), ncm.conSwitch.Size(), 15)))
	ncm.dexSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(
		ncm.conSwitch.DrawArea(), ncm.dexSwitch.Size(), 15)))
	ncm.intSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(
		ncm.dexSwitch.DrawArea(), ncm.intSwitch.Size(), 15)))
	ncm.wisSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(
		ncm.intSwitch.DrawArea(), ncm.wisSwitch.Size(), 15)))
	ncm.sexSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.RightOf(
		ncm.wisSwitch.DrawArea(), ncm.sexSwitch.Size(), 30)))
	ncm.raceSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.BottomOf(
		ncm.sexSwitch.DrawArea(), ncm.raceSwitch.Size(), 10)))
	ncm.aliSwitch.Draw(win.Window, mtk.Matrix().Moved(mtk.BottomOf(
		ncm.raceSwitch.DrawArea(), ncm.aliSwitch.Size(), 10)))
	// Buttons.
	doneButtonPos := mtk.DrawPosBR(win.Bounds(), ncm.doneButton.Size())
	ncm.doneButton.Draw(win.Window, mtk.Matrix().Moved(doneButtonPos))
	backButtonPos := mtk.DrawPosBL(win.Bounds(), ncm.backButton.Size())
	ncm.backButton.Draw(win.Window, mtk.Matrix().Moved(backButtonPos))
	ncm.rollButton.Draw(win.Window, mtk.Matrix().Moved(mtk.BottomOf(
		ncm.pointsBox.DrawArea(), ncm.rollButton.Size(), 5)))
}

// Update updates all menu elements.
func (ncm *NewCharacterMenu) Update(win *mtk.Window) {
	ncm.nameEdit.Update(win)
	ncm.doneButton.Update(win)
	ncm.backButton.Update(win)
	ncm.rollButton.Update(win)
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
	ncm.attrPointsMax = ncm.rng.Intn(newCharAttrsMax-
		newCharAttrsMin) + newCharAttrsMin
	ncm.attrPoints = ncm.attrPointsMax
	ncm.updatePoints()
}

// Show toggles menu visibility.
func (ncm *NewCharacterMenu) Show(show bool) {
	ncm.opened = show
	if ncm.opened {
		ncm.rollPoints()
	}
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
	ncm.pointsBox.SetText(fmt.Sprintf("%d", ncm.attrPoints))
}

// createChar creates new game character.
func (ncm *NewCharacterMenu) createChar() (*character.Character, error) {
	// Name.
	name := ncm.nameEdit.Text()
	// Attributes.
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
	// Gender.
	gender, ok := ncm.sexSwitch.Value().Value.(character.Gender)
	if !ok {
		return nil, fmt.Errorf("fail_to_retrive_gender")
	}
	// Race.
	race, ok := ncm.raceSwitch.Value().Value.(character.Race)
	if !ok {
		return nil, fmt.Errorf("fail_to_retrieve_race")
	}
	// Alignment.
	alignment, ok := ncm.aliSwitch.Value().Value.(character.Alignment)
	if !ok {
		return nil, fmt.Errorf("fail_to_retrieve_alignment")
	}
	// ID.
	name = strings.ReplaceAll(name, " ", "_")
	id := fmt.Sprintf(newCharIDFrom, strings.ToLower(name))
	charData := flameres.CharacterBasicData{
		ID:        id,
		Name:      name,
		Level:     1,
		Sex:       int(gender),
		Race:      int(race),
		Alignment: int(alignment),
		Attitude:  int(character.Friendly),
		Str:       str,
		Con:       con,
		Dex:       dex,
		Int:       inte,
		Wis:       wis,
	}
	char := character.New(charData)
	// Player skills & items from interface config.
	for _, sid := range config.CharSkills {
		s, err := flamedata.Skill(sid)
		if err != nil {
			log.Err.Printf("newchar_menu:fail_to_retrieve_new_player_skill:%v", err)
			continue
		}
		char.AddSkill(s)
	}
	for _, iid := range config.CharItems {
		i, err := flamedata.Item(iid)
		if err != nil {
			log.Err.Printf("newchar_menu:fail_to_retireve_new_player_items:%v", err)
			continue
		}
		char.Inventory().AddItem(i)
	}
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
	avData := res.AvatarData{
		ID:           char.ID(),
		Serial:       char.Serial(),
		PortraitName: portraitName,
		SSHeadName:   ssHeadName,
		SSTorsoName:  ssTorsoName,
		PortraitPic:  portraitPic,
		SSHeadPic:    ssHeadPic,
		SSTorsoPic:   ssTorsoPic,
	}
	av := object.NewAvatar(char, &avData)
	ncm.mainmenu.AddPlayableChar(av)
	msg := lang.Text("gui", "newchar_create_msg")
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
	if pts > -1 && pts <= ncm.attrPointsMax {
		ncm.attrPoints = pts
	} else {
		s.SetIndex(s.Find(old.Value))
	}
	ncm.updatePoints()
}
