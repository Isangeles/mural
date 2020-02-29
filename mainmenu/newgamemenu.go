/*
 * newgamemenu.go
 *
 * Copyright 2018-2020 Dariusz Sikora <dev@isangeles.pl>
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
	"path/filepath"

	"github.com/faiface/pixel"

	"github.com/isangeles/flame"
	flamedata "github.com/isangeles/flame/core/data"
	"github.com/isangeles/flame/core/data/res/lang"
	"github.com/isangeles/flame/core/module/character"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data/exp"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// NewGameMenu struct represents new game
// creation screen.
type NewGameMenu struct {
	mainmenu     *MainMenu
	title        *mtk.Text
	charSwitch   *mtk.Switch
	charInfo     *mtk.Textbox
	startButton  *mtk.Button
	exportButton *mtk.Button
	backButton   *mtk.Button
	opened       bool
}

// newNewGameMenu creates new game creation menu.
func newNewGameMenu(mainmenu *MainMenu) *NewGameMenu {
	ngm := new(NewGameMenu)
	ngm.mainmenu = mainmenu
	// Title.
	titleParams := mtk.Params{
		FontSize: mtk.SizeBig,
	}
	ngm.title = mtk.NewText(titleParams)
	ngm.title.SetText(lang.Text("newgame_menu_title"))
	// Swtches & text.
	charSwitchParams := mtk.Params{
		Size:      mtk.SizeBig,
		MainColor: mainColor,
	}
	ngm.charSwitch = mtk.NewSwitch(charSwitchParams)
	ngm.charSwitch.SetLabel(lang.Text("newgame_char_switch_label"))
	ngm.charSwitch.SetOnChangeFunc(ngm.onCharSwitchChanged)
	charInfoParams := mtk.Params{
		FontSize:    mtk.SizeBig,
		MainColor:   mainColor,
		AccentColor: accentColor,
	}
	ngm.charInfo = mtk.NewTextbox(charInfoParams)
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		FontSize:  mtk.SizeMedium,
		Shape:     mtk.ShapeRectangle,
		MainColor: accentColor,
	}
	ngm.startButton = mtk.NewButton(buttonParams)
	ngm.startButton.SetLabel(lang.Text("newgame_start_button_label"))
	ngm.startButton.SetOnClickFunc(ngm.onStartButtonClicked)
	ngm.exportButton = mtk.NewButton(buttonParams)
	ngm.exportButton.SetLabel(lang.Text("newgame_export_button_label"))
	ngm.exportButton.SetOnClickFunc(ngm.onExportButtonClicked)
	ngm.backButton = mtk.NewButton(buttonParams)
	ngm.backButton.SetLabel(lang.Text("back_b_label"))
	ngm.backButton.SetOnClickFunc(ngm.onBackButtonClicked)
	return ngm
}

// Draw draws all menu elements in specified window.
func (ngm *NewGameMenu) Draw(win *mtk.Window) {
	// Title.
	titlePos := pixel.V(win.Bounds().Center().X,
		win.Bounds().Max.Y-ngm.title.Size().Y)
	ngm.title.Draw(win, mtk.Matrix().Moved(titlePos))
	// Buttons.
	startButtonPos := mtk.DrawPosBR(win.Bounds(), ngm.startButton.Size())
	ngm.startButton.Draw(win.Window, mtk.Matrix().Moved(startButtonPos))
	exportButtonPos := mtk.LeftOf(ngm.startButton.DrawArea(),
		ngm.exportButton.Size(), 10)
	ngm.exportButton.Draw(win, mtk.Matrix().Moved(exportButtonPos))
	backButtonPos := mtk.DrawPosBL(win.Bounds(), ngm.backButton.Size())
	ngm.backButton.Draw(win.Window, mtk.Matrix().Moved(backButtonPos))
	// Portrait switch.
	charSwitchPos := mtk.BottomOf(ngm.title.DrawArea(), ngm.charSwitch.Size(), 10)
	ngm.charSwitch.Draw(win, mtk.Matrix().Moved(charSwitchPos))
	// Character info.
	charInfoSize := pixel.V(win.Bounds().W(), win.Bounds().H()/2)
	ngm.charInfo.SetSize(charInfoSize)
	charInfoPos := mtk.BottomOf(ngm.charSwitch.DrawArea(), ngm.charInfo.Size(), mtk.ConvSize(10))
	ngm.charInfo.Draw(win, mtk.Matrix().Moved(charInfoPos))
}

// Update updates all menu elements.
func (ngm *NewGameMenu) Update(win *mtk.Window) {
	ngm.charSwitch.Update(win)
	ngm.charInfo.Update(win)
	ngm.startButton.Update(win)
	ngm.exportButton.Update(win)
	ngm.backButton.Update(win)
}

// Show toggles menu visibility.
func (ngm *NewGameMenu) Show(show bool) {
	ngm.opened = show
	ngm.updateCharInfo()
	ngm.updateCharSwitch()
}

// Opened checks whether menu is open.
func (ngm *NewGameMenu) Opened() bool {
	return ngm.opened
}

// SetCharacters sets specified avatars as playable characters
// fo new game start.
func (ngm *NewGameMenu) SetCharacters(chars []PlayableCharData) {
	values := make([]mtk.SwitchValue, len(chars))
	for i, c := range chars {
		values[i] = mtk.SwitchValue{c.AvatarData.PortraitPic, c}
	}
	ngm.charSwitch.SetValues(values)
}

// updateCharInfo updates textbox with character informations.
func (ngm *NewGameMenu) updateCharInfo() {
	switchVal := ngm.charSwitch.Value()
	if switchVal == nil {
		return
	}
	c, ok := switchVal.Value.(PlayableCharData)
	if !ok {
		log.Err.Printf("unable to retrieve character data from switch")
		return
	}
	charData := c.CharData.BasicData
	charInfoForm := `Name:       %s
Level:      %d
Gender:     %s
Race:       %s
Alignment   %s
Attributes: %d, %d, %d, %d, %d`
	info := fmt.Sprintf(charInfoForm, charData.Name, charData.Level,
		lang.Text(character.Gender(charData.Sex).ID()),
		lang.Text(character.Race(charData.Race).ID()),
		lang.Text(character.Alignment(charData.Alignment).ID()),
		charData.Str, charData.Con, charData.Dex, charData.Int, charData.Wis)
	ngm.charInfo.SetText(info)
	return
}

// updateCharSwitch updates menu character switch.
func (ngm *NewGameMenu) updateCharSwitch() {
	ngm.charSwitch.SetValues(make([]mtk.SwitchValue, 0))
	ngm.SetCharacters(ngm.mainmenu.PlayableChars())
}

// exportChar exports currently selected character.
func (ngm *NewGameMenu) exportChar() error {
	switchVal := ngm.charSwitch.Value()
	if switchVal == nil {
		return nil
	}
	data, ok := switchVal.Value.(PlayableCharData)
	if !ok {
		return fmt.Errorf("unable to retrieve character data from switch")
	}
	c := character.New(*data.CharData)
	path := filepath.Join(flame.Mod().Conf().CharactersPath(), c.Name())
	err := flamedata.ExportCharacters(path, c)
	if err != nil {
		return fmt.Errorf("unable to export characters: %v", err)
	}
	av := object.NewAvatar(c, data.AvatarData)
	err = exp.ExportAvatar(av, flame.Mod().Conf().CharactersPath())
	if err != nil {
		return fmt.Errorf("unable to export avatar: %v", err)
	}
	return nil
}

// startGame starts new game.
func (ngm *NewGameMenu) startGame() {
	// Show loading screen.
	ngm.mainmenu.OpenLoadingScreen(lang.Text("newgame_start_info"))
	defer ngm.mainmenu.CloseLoadingScreen()
	// Retrive character from character switch.
	switchVal := ngm.charSwitch.Value()
	if switchVal == nil {
		log.Err.Printf("main menu: new game: no char switch value")
		return
	}
	pcd, ok := switchVal.Value.(PlayableCharData)
	if !ok {
		log.Err.Printf("main menu: new game: unable to retrieve avatar from switch")
		return
	}
	// Create game.
	c := character.New(*pcd.CharData)
	g, err := flame.StartGame(c)
	if err != nil {
		log.Err.Printf("main menu: new game: unable to start game: %v", err)
		return
	}
	// Create pc avatar.
	av := object.NewAvatar(c, pcd.AvatarData)
	// Pass new game.
	if ngm.mainmenu.onGameCreated == nil {
		return
	}
	ngm.mainmenu.onGameCreated(g, av)
}

// Triggered after start button clicked.
func (ngm *NewGameMenu) onStartButtonClicked(b *mtk.Button) {
	go ngm.startGame()
	ngm.mainmenu.OpenMenu()
}

// Triggered after back button clicked.
func (ngm *NewGameMenu) onBackButtonClicked(b *mtk.Button) {
	ngm.mainmenu.OpenMenu()
}

// Triggered after export button clicked.
func (ngm *NewGameMenu) onExportButtonClicked(b *mtk.Button) {
	err := ngm.exportChar()
	if err != nil {
		log.Err.Printf("main menu: new game: unable to export character: %v", err)
		return
	}
	msg := lang.Text("newgame_export_msg")
	ngm.mainmenu.ShowMessage(msg)
}

// Triggered after character switch change.
func (ngm *NewGameMenu) onCharSwitchChanged(s *mtk.Switch,
	old, new *mtk.SwitchValue) {
	ngm.updateCharInfo()
}
