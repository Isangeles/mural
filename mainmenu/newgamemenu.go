/*
 * newgamemenu.go
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

	"github.com/faiface/pixel"

	"github.com/isangeles/flame"
	flamedata "github.com/isangeles/flame/core/data"
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/object/character"

	"github.com/isangeles/mural/core/data/exp"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/mtk"
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
	ngm.title = mtk.NewText(mtk.SIZE_BIG, 0)
	ngm.title.SetText(lang.Text("gui", "newgame_menu_title"))
	// Swtches & text.
	ngm.charSwitch = mtk.NewSwitch(mtk.SIZE_BIG, main_color)
	ngm.charSwitch.SetLabel(lang.Text("gui", "newgame_char_switch_label"))
	ngm.charSwitch.SetOnChangeFunc(ngm.onCharSwitchChanged)
	ngm.charInfo = mtk.NewTextbox(pixel.V(0, 0), mtk.SIZE_MINI, mtk.SIZE_BIG,
		accent_color, main_color)
	// Buttons.
	ngm.startButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		accent_color)
	ngm.startButton.SetLabel(lang.Text("gui", "newgame_start_button_label"))
	ngm.startButton.SetOnClickFunc(ngm.onStartButtonClicked)
	ngm.exportButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		accent_color)
	ngm.exportButton.SetLabel(lang.Text("gui", "newgame_export_button_label"))
	ngm.exportButton.SetOnClickFunc(ngm.onExportButtonClicked)
	ngm.backButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		accent_color)
	ngm.backButton.SetLabel(lang.Text("gui", "back_b_label"))
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
func (ngm *NewGameMenu) SetCharacters(chars []*object.Avatar) {
	values := make([]mtk.SwitchValue, len(chars))
	for i, c := range chars {
		values[i] = mtk.SwitchValue{c.Portrait(), c}
	}
	ngm.charSwitch.SetValues(values)
}

// updateCharInfo updates textbox with character informations.
func (ngm *NewGameMenu) updateCharInfo() error {
	switchVal := ngm.charSwitch.Value()
	if switchVal == nil {
		return nil
	}
	c, ok := switchVal.Value.(*object.Avatar)
	if !ok {
		return fmt.Errorf("fail_to_retrieve_avatar_from_switch")
	}
	charInfoForm := `
                         Name:       %s
                         Level:      %d
                         Gender:     %s
                         Race:       %s
                         Alignment   %s
                         Attributes: %s`
	info := fmt.Sprintf(charInfoForm, c.Name(), c.Level(),
		lang.Text("ui", c.Gender().ID()), lang.Text("ui", c.Race().ID()),
		lang.Text("ui", c.Alignment().ID()), c.Attributes().String())
	ngm.charInfo.SetText(info)
	return nil
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
	c, ok := switchVal.Value.(*object.Avatar)
	if !ok {
		return fmt.Errorf("fail_to_retrieve_avatar_from_switch")
	}
	err := flamedata.ExportCharacter(c.Character, flame.Mod().Conf().CharactersPath())
	if err != nil {
		return err
	}
	err = exp.ExportAvatar(c, flame.Mod().Conf().CharactersPath())
	if err != nil {
		return fmt.Errorf("fail_to_export_avatar:%v", err)
	}
	return nil
}

// startGame starts new game.
func (ngm *NewGameMenu) startGame() {
	ngm.mainmenu.OpenLoadingScreen(lang.Text("gui", "newgame_start_info"))
	defer ngm.mainmenu.CloseLoadingScreen()
	switchVal := ngm.charSwitch.Value()
	if switchVal == nil {
		log.Err.Printf("main_menu:new_game:no char switch value")
		return
	}
	// Character from avatar switch.
	c, ok := switchVal.Value.(*object.Avatar)
	if !ok {
		log.Err.Printf("main_menu:new_game:fail to retrieve avatar from switch")
		return
	}
	// Add avatar data to resources base.
	res.AddAvatarData(c.Data())
	// Create game.
	g, err := flame.StartGame([]*character.Character{c.Character})
	if err != nil {
		log.Err.Printf("main_menu:new_game:fail_to_start_game:%v", err)
		return
	}
	// Pass new game.
	if ngm.mainmenu.onGameCreated == nil {
		return
	}
	ngm.mainmenu.onGameCreated(g)
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
		log.Err.Printf("main_menu:new_game:fail_to_export_character:%v", err)
		return
	}
	msg := lang.Text("gui", "newgame_export_msg")
	ngm.mainmenu.ShowMessage(msg)
}

// Triggered after character switch change.
func (ngm *NewGameMenu) onCharSwitchChanged(s *mtk.Switch,
	old, new *mtk.SwitchValue) {
	err := ngm.updateCharInfo()
	if err != nil {
		log.Err.Printf("main_menu:new_game:fail_to_update_char_info:%v\n", err)
	}
}
