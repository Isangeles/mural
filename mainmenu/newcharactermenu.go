/*
 * newcharactermenu.go
 *
 * Copyright 2018-2024 Dariusz Sikora <ds@isangeles.dev>
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
	"path/filepath"
	"strings"
	"time"

	"github.com/gopxl/pixel"

	"github.com/isangeles/flame/character"
	flameres "github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/data/res/lang"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/data"
	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/log"
)

var (
	playerIDPrefix = "player_"
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
	ncm.title.SetText(lang.Text("newchar_menu_title"))
	// Name Edit.
	labelParams := mtk.Params{
		FontSize: mtk.SizeMedium,
	}
	ncm.nameLabel = mtk.NewText(labelParams)
	nameLabelText := lang.Text("newchar_name_edit_label")
	ncm.nameLabel.SetText(fmt.Sprintf("%s", nameLabelText))
	texteditParams := mtk.Params{
		FontSize:  mtk.SizeMedium,
		MainColor: mainColor,
	}
	ncm.nameEdit = mtk.NewTextedit(texteditParams)
	// Points box.
	pointsBoxSize := mtk.SizeMedium.ButtonSize(mtk.ShapeRectangle)
	pointsBoxParams := mtk.Params{
		SizeRaw:  pointsBoxSize,
		FontSize: mtk.SizeBig,
	}
	ncm.pointsBox = mtk.NewText(pointsBoxParams)
	// Portrait switch.
	faceSwitchParams := mtk.Params{
		Size:      mtk.SizeBig,
		MainColor: mainColor,
	}
	ncm.faceSwitch = mtk.NewSwitch(faceSwitchParams)
	ncm.faceSwitch.SetLabel(lang.Text("newchar_face_switch_label"))
	portraitsPath := filepath.Join(config.GUIPath, "portraits")
	faces, err := data.Pictures(portraitsPath)
	if err != nil {
		log.Err.Printf("new char menu: unable to retrieve player portraits: %v", err)

	}
	faceValues := make([]mtk.SwitchValue, 0)
	for n, p := range faces {
		value := mtk.SwitchValue{p, n}
		faceValues = append(faceValues, value)
	}
	ncm.faceSwitch.SetValues(faceValues...)
	// Attributes switches.
	attrSwitchParams := mtk.Params{
		Size:      mtk.SizeMedium,
		MainColor: mainColor,
	}
	ncm.strSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.strSwitch.SetLabel(lang.Text("newchar_str_switch_label"))
	ncm.strSwitch.SetIntValues(0, 90)
	ncm.strSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.conSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.conSwitch.SetLabel(lang.Text("newchar_con_switch_label"))
	ncm.conSwitch.SetIntValues(0, 90)
	ncm.conSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.dexSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.dexSwitch.SetLabel(lang.Text("newchar_dex_switch_label"))
	ncm.dexSwitch.SetIntValues(0, 90)
	ncm.dexSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.intSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.intSwitch.SetLabel(lang.Text("newchar_int_switch_label"))
	ncm.intSwitch.SetIntValues(0, 90)
	ncm.intSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	ncm.wisSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.wisSwitch.SetLabel(lang.Text("newchar_wis_switch_label"))
	ncm.wisSwitch.SetIntValues(0, 90)
	ncm.wisSwitch.SetOnChangeFunc(ncm.onAttrSwitchChange)
	// Gender & alligment switches.
	maleSwitchVal := mtk.SwitchValue{lang.Text(string(character.Male)),
		string(character.Male)}
	femaleSwitchVal := mtk.SwitchValue{lang.Text(string(character.Female)),
		string(character.Female)}
	gens := []mtk.SwitchValue{maleSwitchVal, femaleSwitchVal}
	ncm.sexSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.sexSwitch.SetLabel(lang.Text("newchar_sex_switch_label"))
	ncm.sexSwitch.SetValues(gens...)
	// Race switch.
	races := []mtk.SwitchValue{}
	for _, r := range flameres.Races {
		if !r.Playable {
			continue
		}
		val := mtk.SwitchValue{lang.Text(r.ID), r.ID}
		races = append(races, val)
	}
	ncm.raceSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.raceSwitch.SetLabel(lang.Text("newchar_race_switch_label"))
	ncm.raceSwitch.SetValues(races...)
	// Alignment switch.
	alis := []mtk.SwitchValue{
		mtk.SwitchValue{lang.Text(string(character.LawfulGood)), string(character.LawfulGood)},
		mtk.SwitchValue{lang.Text(string(character.NeutralGood)), string(character.NeutralGood)},
		mtk.SwitchValue{lang.Text(string(character.ChaoticGood)), string(character.ChaoticGood)},
		mtk.SwitchValue{lang.Text(string(character.LawfulNeutral)), string(character.LawfulNeutral)},
		mtk.SwitchValue{lang.Text(string(character.TrueNeutral)), string(character.TrueNeutral)},
		mtk.SwitchValue{lang.Text(string(character.ChaoticNeutral)), string(character.ChaoticNeutral)},
		mtk.SwitchValue{lang.Text(string(character.LawfulEvil)), string(character.LawfulEvil)},
		mtk.SwitchValue{lang.Text(string(character.NeutralEvil)), string(character.NeutralEvil)},
		mtk.SwitchValue{lang.Text(string(character.ChaoticEvil)), string(character.ChaoticEvil)},
	}
	ncm.aliSwitch = mtk.NewSwitch(attrSwitchParams)
	ncm.aliSwitch.SetLabel(lang.Text("newchar_ali_switch_label"))
	ncm.aliSwitch.SetValues(alis...)
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		FontSize:  mtk.SizeMedium,
		Shape:     mtk.ShapeRectangle,
		MainColor: accentColor,
	}
	ncm.doneButton = mtk.NewButton(buttonParams)
	ncm.doneButton.SetLabel(lang.Text("done_button_label"))
	ncm.doneButton.SetOnClickFunc(ncm.onDoneButtonClicked)
	ncm.backButton = mtk.NewButton(buttonParams)
	ncm.backButton.SetLabel(lang.Text("back_button_label"))
	ncm.backButton.SetOnClickFunc(ncm.onBackButtonClicked)
	ncm.rollButton = mtk.NewButton(buttonParams)
	ncm.rollButton.SetLabel(lang.Text("newchar_roll_button_label"))
	ncm.rollButton.SetInfo(lang.Text("newchar_roll_button_info"))
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
	ncm.doneButton.Active(ncm.canCreate())
}

// rollPoints rolls random amount of attribute points for new character
// from range specified in the chapter config of active module.
func (ncm *NewCharacterMenu) rollPoints() {
	ncm.strSwitch.Reset()
	ncm.conSwitch.Reset()
	ncm.dexSwitch.Reset()
	ncm.intSwitch.Reset()
	ncm.wisSwitch.Reset()
	attrsMax := ncm.mainmenu.mod.Chapter().Conf().StartAttrs
	if attrsMax < 1 {
		return
	}
	ncm.attrPointsMax = ncm.rng.Intn(attrsMax-1) + 1
	ncm.attrPoints = ncm.attrPointsMax
	ncm.updatePoints()
}

// Show shows menu.
func (ncm *NewCharacterMenu) Show() {
	ncm.opened = true
	ncm.rollPoints()
}

// Hide hides menu.
func (ncm *NewCharacterMenu) Hide() {
	ncm.opened = false
}

// Opened checks whether menu is open.
func (ncm *NewCharacterMenu) Opened() bool {
	return ncm.opened
}

// canCreate checks whether its possible to create new character.
func (ncm *NewCharacterMenu) canCreate() bool {
	return len(ncm.nameEdit.Text()) > 0 && ncm.attrPoints < 1
}

// updatePoints updates points box value.
func (ncm *NewCharacterMenu) updatePoints() {
	ncm.pointsBox.SetText(fmt.Sprintf("%d", ncm.attrPoints))
}

// createChar creates new game character.
func (ncm *NewCharacterMenu) createCharData() (flameres.CharacterData, error) {
	charData := flameres.CharacterData{
		Level:     1,
		Attitude:  string(character.Friendly),
	}
	// ID.
	name := ncm.nameEdit.Text()
	id := strings.ReplaceAll(name, " ", "_")
	charData.ID = fmt.Sprintf("%s%s", playerIDPrefix, strings.ToLower(id))
	// Add name translation.
	nameTrans := flameres.TranslationData{charData.ID, []string{name}}
	lang.AddTranslation(nameTrans)
	// Attributes.
	var ok bool
	charData.Attributes.Str, ok = ncm.strSwitch.Value().Value.(int)
	if !ok {
		return charData, fmt.Errorf("unable to retrieve strenght switch value")
	}
	charData.Attributes.Con, ok = ncm.conSwitch.Value().Value.(int)
	if !ok {
		return charData, fmt.Errorf("unable to retrieve constitution switch value")
	}
	charData.Attributes.Dex, ok = ncm.dexSwitch.Value().Value.(int)
	if !ok {
		return charData, fmt.Errorf("unable to retireve dexterity switch value")
	}
	charData.Attributes.Con, ok = ncm.intSwitch.Value().Value.(int)
	if !ok {
		return charData, fmt.Errorf("unable to retrieve inteligence switch value")
	}
	charData.Attributes.Wis, ok = ncm.wisSwitch.Value().Value.(int)
	if !ok {
		return charData, fmt.Errorf("unable to retrieve wisdom switch value")
	}
	// Gender.
	charData.Sex, ok = ncm.sexSwitch.Value().Value.(string)
	if !ok {
		return charData, fmt.Errorf("unable to retrive gender")
	}
	// Race.
	charData.Race, ok = ncm.raceSwitch.Value().Value.(string)
	if !ok {
		return charData, fmt.Errorf("unable to retrieve race")
	}
	// Alignment.
	charData.Alignment, ok = ncm.aliSwitch.Value().Value.(string)
	if !ok {
		return charData, fmt.Errorf("unable to retrieve alignment")
	}
	// Player skills & items from the chapter config.
	for _, sid := range ncm.mainmenu.mod.Chapter().Conf().StartSkills {
		skill := flameres.ObjectSkillData{ID: sid}
		charData.Skills = append(charData.Skills, skill)
	}
	for _, iid := range ncm.mainmenu.mod.Chapter().Conf().StartItems {
		item := flameres.InventoryItemData{ID: iid}
		charData.Inventory.Items = append(charData.Inventory.Items, item)
	}
	return charData, nil
}

// createAvatarData Creates avatar data for specified character data.
func (ncm *NewCharacterMenu) createAvatarData(charData flameres.CharacterData) (res.AvatarData, error) {
	avData := res.AvatarData{
		ID:       charData.ID,
		Serial:   charData.Serial,
	}
	avData.Head = "m-head-black-1222211-80x90.png"
	avData.Torso = "m-cloth-1222211-80x90.png"
	if character.Gender(charData.Sex) == character.Female {
		avData.Head = "f-cloth-1222211-80x90.png"
		avData.Torso = "f-head-black-1222211-80x90.png"
	}
	var ok bool
	avData.Portrait, ok = ncm.faceSwitch.Value().Value.(string)
	if !ok {
		return avData, fmt.Errorf("unable to retrieve portrait name from switch")
	}
	res.Avatars = append(res.Avatars, avData)
	return avData, nil
}

// Triggered after back button clicked.
func (ncm *NewCharacterMenu) onBackButtonClicked(b *mtk.Button) {
	ncm.mainmenu.OpenMenu()
}

// Triggered after done button clicked.
func (ncm *NewCharacterMenu) onDoneButtonClicked(b *mtk.Button) {
	// Create character.
	charData, err := ncm.createCharData()
	if err != nil {
		log.Err.Printf("new char menu: unable to create character: %v", err)
		return
	}
	// Create avatar.
	avData, err := ncm.createAvatarData(charData)
	if err != nil {
		log.Err.Printf("new char menu: unable to create avatar: %v", err)
		return
	}
	// Add as playable char.
	pc := PlayableCharData{charData, avData}
	ncm.mainmenu.AddPlayableChar(pc)
	// Show message & go back to main menu.
	msg := lang.Text("newchar_create_msg")
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
		log.Err.Print("new char menu: unable to retrieve str switch value")
		return
	}
	con, ok := ncm.conSwitch.Value().Value.(int)
	if !ok {
		log.Err.Print("new char menu: unable to retrieve con switch value")
		return
	}
	dex, ok := ncm.dexSwitch.Value().Value.(int)
	if !ok {
		log.Err.Print("new char menu: unable to retrieve con switch value")
		return
	}
	inte, ok := ncm.intSwitch.Value().Value.(int)
	if !ok {
		log.Err.Print("new char menu: unable to retrieve int switch value")
		return
	}
	wis, ok := ncm.wisSwitch.Value().Value.(int)
	if !ok {
		log.Err.Print("new char menu: unable to retrieve wis switch value")
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
