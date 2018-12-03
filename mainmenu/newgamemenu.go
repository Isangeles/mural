/*
 * newgamemenu.go
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

	"github.com/isangeles/flame"
	flamedata "github.com/isangeles/flame/core/data"
	"github.com/isangeles/flame/core/module/object/character"
	"github.com/isangeles/flame/core/data/text/lang"
	
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/objects"
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
func newNewGameMenu(mainmenu *MainMenu) (*NewGameMenu, error) {
	ngm := new(NewGameMenu)
	ngm.mainmenu = mainmenu
	// Title.
	ngm.title = mtk.NewText(lang.Text("gui", "newgame_menu_title"),
		mtk.SIZE_BIG, 0)
	// Swtches & text.
	ngm.charSwitch = mtk.NewSwitch(mtk.SIZE_BIG, main_color,
		lang.Text("gui", "newgame_char_switch_label"), "", nil)
	ngm.charSwitch.SetOnChangeFunc(ngm.onCharSwitchChanged)
	ngm.charInfo = mtk.NewTextbox(mtk.SIZE_BIG, main_color)
	// Buttons.
	ngm.startButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		accent_color, lang.Text("gui", "newgame_start_button_label"), "")
	ngm.startButton.SetOnClickFunc(ngm.onStartButtonClicked)
	ngm.exportButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		accent_color, lang.Text("gui", "newgame_export_button_label"), "")
	ngm.exportButton.SetOnClickFunc(ngm.onExportButtonClicked)
	ngm.backButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		accent_color, lang.Text("gui", "back_b_label"), "")
	ngm.backButton.SetOnClickFunc(ngm.onBackButtonClicked)
	ngm.updateCharSwitchValues()
	return ngm, nil
}

// Draw draws all menu elements in specified window.
func (ngm *NewGameMenu) Draw(win *mtk.Window) {
	// Title.
	titlePos := pixel.V(win.Bounds().Center().X,
		win.Bounds().Max.Y-ngm.title.Bounds().Size().Y)
	ngm.title.Draw(win, mtk.Matrix().Moved(titlePos))
	// Buttons.
	ngm.startButton.Draw(win.Window, mtk.Matrix().Moved(mtk.PosBR(
		ngm.startButton.Frame(), pixel.V(win.Bounds().Max.X,
			win.Bounds().Min.Y))))
	ngm.exportButton.Draw(win, mtk.Matrix().Moved(mtk.LeftOf(
		ngm.startButton.DrawArea(), ngm.exportButton.Frame(), 10)))
	ngm.backButton.Draw(win.Window, mtk.Matrix().Moved(mtk.PosBL(
		ngm.backButton.Frame(), win.Bounds().Min)))
	// Switches & text.
	ngm.charSwitch.Draw(win, mtk.Matrix().Moved(mtk.BottomOf(
		ngm.title.DrawArea(), ngm.charSwitch.Frame(), 10)))
	ngm.charInfo.Draw(pixel.R(win.Bounds().Min.X,
		ngm.backButton.DrawArea().Max.Y, win.Bounds().Max.X,
		ngm.charSwitch.DrawArea().Min.Y), win.Window)
}

// Update updates all menu elements.
func (ngm *NewGameMenu) Update(win *mtk.Window) {
	if ngm.Opened() {
		ngm.charSwitch.Update(win)
		ngm.charInfo.Update(win)
		ngm.startButton.Update(win)
		ngm.exportButton.Update(win)
		ngm.backButton.Update(win)
	}
}

// Show toggles menu visibility.
func (ngm *NewGameMenu) Show(show bool) {
	ngm.opened = show
	ngm.updateCharSwitchValues()
	ngm.updateCharInfo()
}

// Opened checks whether menu is open.
func (ngm *NewGameMenu) Opened() bool {
	return ngm.opened
}

// updateCharSwitchValues updates character switch values.
func (ngm *NewGameMenu) updateCharSwitchValues() {
	charSwitchValues := make([]mtk.SwitchValue, len(ngm.mainmenu.PlayableChars))
	for i, c := range ngm.mainmenu.PlayableChars {
		charSwitchValues[i] = mtk.SwitchValue{c.Portrait(), c}
	}
	ngm.charSwitch.SetValues(charSwitchValues)
}

// updateCharInfo updates textbox with character informations.
func (ngm *NewGameMenu) updateCharInfo() error {
	switchVal := ngm.charSwitch.Value()
	if switchVal == nil {
		return nil
	}
	c, ok := switchVal.Value.(*objects.Avatar)
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
	ngm.charInfo.Clear()
	ngm.charInfo.Add(fmt.Sprintf(charInfoForm, c.Name(), c.Level(),
		lang.Text("ui", c.Gender().Id()), lang.Text("ui", c.Race().Id()),
	        lang.Text("ui", c.Alignment().Id()), c.Attributes().String()))
	return nil
}

// exportChar exports currently selected character.
func (ngm *NewGameMenu) exportChar() error {
	switchVal := ngm.charSwitch.Value()
	if switchVal == nil {
		return nil
	}
	c, ok := switchVal.Value.(*objects.Avatar)
	if !ok {
		return fmt.Errorf("fail_to_retrieve_avatar_from_switch")
	}
	err := flamedata.ExportCharacter(c.Character, flame.Mod().CharactersPath())
	if err != nil {
		return err
	}
	err = data.ExportAvatar(c, flame.Mod().CharactersPath())
	if err != nil {
		return fmt.Errorf("fail_to_export_avatar:%v", err)
	}
	return nil
}

// startGame starts new game.
func (ngm *NewGameMenu) startGame() error {
	switchVal := ngm.charSwitch.Value()
	if switchVal == nil {
		return fmt.Errorf("no_char_switch_val")
	}
	c, ok := switchVal.Value.(*objects.Avatar)
	if !ok {
		return fmt.Errorf("fail_to_retrieve_avatar_from_switch")
	}
	g, err := flame.StartGame([]*character.Character{c.Character})
	if err != nil {
		return err
	}
	ngm.mainmenu.OnNewGameCreated(g, c)
	return nil
}

// Triggered after start button clicked.
func (ngm *NewGameMenu) onStartButtonClicked(b *mtk.Button) {
	err := ngm.startGame()
	if err != nil {
		log.Err.Printf("fail_to_start_new_game:%v\n", err)
	}
}

// Triggered after back button clicked.
func (ngm *NewGameMenu) onBackButtonClicked(b *mtk.Button) {
	ngm.mainmenu.OpenMenu()
}

// Triggered after export button clicked.
func (ngm *NewGameMenu) onExportButtonClicked(b *mtk.Button) {
	err := ngm.exportChar()
	if err != nil {
		log.Err.Printf("fail_to_export_character:%v", err)
		return
	}
	msg := mtk.NewMessageWindow(mtk.SIZE_SMALL, lang.Text("gui",
		"newgame_export_msg"))
	ngm.mainmenu.ShowMessage(msg)
}

// Triggered after character switch change.
func (ngm *NewGameMenu) onCharSwitchChanged(s *mtk.Switch,
		old, new *mtk.SwitchValue) {
	err := ngm.updateCharInfo()
	if err != nil {
		log.Err.Printf("fail_to_update_char_info:%v\n", err)
	}
}
