/*
 * settings.go
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
	
	"golang.org/x/image/colornames"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/log"
)

// Settings struct represents main menu
// settings screen.
type Settings struct {
	mainmenu      *MainMenu
	title         *mtk.Text
	backButton    *mtk.Button
	fullscrSwitch *mtk.Switch
	resSwitch     *mtk.Switch
	langSwitch    *mtk.Switch
	opened        bool
	changed       bool
}

// newSettings returns new settings screen
// instance.
func newSettings(mainmenu *MainMenu) (*Settings, error) {
	s := new(Settings)
	s.mainmenu = mainmenu
	// Title.
	s.title = mtk.NewText(mtk.SIZE_BIG, 900)
	s.title.SetText(lang.Text("gui", "settings_menu_title"))
	// Buttons.
	s.backButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE,
		colornames.Red)
	s.backButton.SetLabel(lang.Text("gui", "back_b_label"))
	s.backButton.SetOnClickFunc(s.onBackButtonClicked)
	// Switches.
	fullscrTrue := mtk.SwitchValue{lang.Text("ui", "com_yes"), true}
	fullscrFalse := mtk.SwitchValue{lang.Text("ui", "com_no"), false}
	fullscrValues := []mtk.SwitchValue{fullscrFalse, fullscrTrue}
	s.fullscrSwitch = mtk.NewSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "settings_fullscr_switch_label"), "", fullscrValues)
	s.fullscrSwitch.SetOnChangeFunc(s.onSettingsChanged)
	var resSwitchValues []mtk.SwitchValue
	for _, res := range config.SupportedResolutions() {
		resSwitchValues = append(resSwitchValues,
			mtk.SwitchValue{fmt.Sprintf("%vx%v", res.X, res.Y), res})
	}
	s.resSwitch = mtk.NewSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "resolution_s_label"), "", resSwitchValues)
	s.resSwitch.SetOnChangeFunc(s.onSettingsChanged)
	s.langSwitch = mtk.NewSwitch(mtk.SIZE_MEDIUM, main_color,
		lang.Text("gui", "lang_s_label"), "", nil)
	s.langSwitch.SetTextValues(config.SupportedLangs())
	s.langSwitch.SetOnChangeFunc(s.onSettingsChanged)
	
	return s, nil
}

// Draw draws all menu elements.
func (s *Settings) Draw(win *pixelgl.Window) {
	// Title.
	titlePos :=pixel.V(win.Bounds().Center().X,
		win.Bounds().Max.Y - s.title.Bounds().Size().Y)
	s.title.Draw(win, mtk.Matrix().Moved(titlePos))
	// Buttons & switches.
	s.fullscrSwitch.Draw(win, mtk.Matrix().Moved(mtk.BottomOf(
		s.title.DrawArea(), s.fullscrSwitch.Bounds(), 50)))
	s.resSwitch.Draw(win, mtk.Matrix().Moved(mtk.BottomOf(
		s.fullscrSwitch.DrawArea(), s.resSwitch.Bounds(), 30)))
	s.langSwitch.Draw(win, mtk.Matrix().Moved(mtk.BottomOf(
		s.resSwitch.DrawArea(), s.langSwitch.Bounds(), 30)))
	s.backButton.Draw(win, mtk.Matrix().Moved(mtk.BottomOf(
		s.langSwitch.DrawArea(), s.backButton.Frame(), 30)))
}

// Update updates all menu elements.
func (s *Settings) Update(win *mtk.Window) {
	s.fullscrSwitch.Update(win)
	s.resSwitch.Update(win)
	s.langSwitch.Update(win)
	s.backButton.Update(win)
}

// Opened checks whether menu should be drawn or not.
func (s *Settings) Opened() bool {
	return s.opened
}

// Show toggles menu visibility.
func (s *Settings) Show(show bool) {
	s.opened = show
	s.updateValues()
}

// Apply applies current settings values.
func (s *Settings) Apply() {
	// Fullscreen.
	fscr, ok := s.fullscrSwitch.Value().Value.(bool)
	if !ok {
		log.Err.Printf(
			"settings_menu:fail_to_retrive_fullscreen_switch_value")
		return
	}
	// Resolution.
	res, ok := s.resSwitch.Value().Value.(pixel.Vec)
	if !ok {
		log.Err.Printf("settings_menu:fail_to_retrive_res_switch_value")
		return
	}
	// Language.
	lang, ok := s.langSwitch.Value().Value.(string)
	if !ok {
		log.Err.Printf("settings_menu:fail_to_retrive_lang_switch_value")
		return
	}

	config.SetFullscreen(fscr)
	config.SetResolution(res)
	config.SetLang(lang)
}

// Changed checks if any settings value was changed.
func (s *Settings) Changed() bool {
	return s.changed
}

// updateValues values of all settings elements.
func (s *Settings) updateValues() {
	fullscrSwitchIndex := s.fullscrSwitch.Find(config.Fullscreen())
	s.fullscrSwitch.SetIndex(fullscrSwitchIndex)
	resSwitchIndex := s.resSwitch.Find(config.Resolution())
	s.resSwitch.SetIndex(resSwitchIndex)
	langSwitchIndex := s.langSwitch.Find(config.Lang())
	s.langSwitch.SetIndex(langSwitchIndex)
}

// close closes settings menu and displays message
// about required game restart if settings was changed.
func (s *Settings) close() {
	if s.Changed() {
		msg := lang.Text("gui", "settings_reset_msg")
		s.mainmenu.ShowMessage(msg)
		s.Apply()
	}
	s.mainmenu.OpenMenu()
}

// closeWithDialog creates settings apply dialog and puts it on
// main menu messages list.
func (s *Settings) closeWithDialog() {
	if s.Changed() {
		dlg := mtk.NewDialogWindow(mtk.SIZE_SMALL,
			lang.Text("gui", "settings_save_msg"))
		dlg.SetOnAcceptFunc(s.onSettingsApplyAccept)
		dlg.SetOnCancelFunc(s.onSettingsApplyCancel)
		s.mainmenu.ShowMessageWindow(dlg)
	} else {
		s.close()
	}	
}

// Triggered after settings change.
func (s *Settings) onSettingsChanged(sw *mtk.Switch, old, new *mtk.SwitchValue) {
	s.changed = true
}

// Triggered after back button clicked.
func (s *Settings) onBackButtonClicked(b *mtk.Button) {
	s.closeWithDialog()
}

// Triggered after settings apply dialog accepted.
func (s *Settings) onSettingsApplyAccept(m *mtk.MessageWindow) {
	s.close()
}

// Triggered after settings apply dialog dismissed.
func (s *Settings) onSettingsApplyCancel(m *mtk.MessageWindow) {
	s.close()
}
