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
	
	"golang.org/x/image/colornames"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"

	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/log"
)

// Settings struct represents main menu
// settings screen.
type Settings struct {
	mainmenu   *MainMenu
	title      *text.Text
	backButton *mtk.Button
	resSwitch  *mtk.Switch
	langSwitch *mtk.Switch
	opened     bool
	changed    bool
}

// newSettings returns new settings screen
// instance.
func newSettings(mainmenu *MainMenu) (*Settings, error) {
	s := new(Settings)
	s.mainmenu = mainmenu
	// Title.
	font := mtk.MainFont(mtk.SIZE_BIG)
	atlas := mtk.Atlas(&font)
	s.title = text.New(pixel.V(0, 0), atlas)
	fmt.Fprintf(s.title, lang.Text("gui", "settings_menu_title"))
	// Buttons & switches.
	s.backButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE, colornames.Red,
		lang.Text("gui", "back_b_label"), "")
	s.backButton.SetOnClickFunc(s.onBackButtonClicked)
	var resSwitchValues []mtk.SwitchValue
	for _, res := range config.SupportedResolutions() {
		resSwitchValues = append(resSwitchValues,
			mtk.SwitchValue{fmt.Sprintf("%vx%v", res.X, res.Y), res})
	}
	s.resSwitch = mtk.NewSwitch(mtk.SIZE_MEDIUM, colornames.Blue,
		lang.Text("gui", "resolution_s_label"), resSwitchValues)
	resSwitchIndex := s.resSwitch.Find(config.Resolution())
	s.resSwitch.SetIndex(resSwitchIndex)
	s.resSwitch.SetOnChangeFunc(s.onSettingsChanged)
	s.langSwitch = mtk.NewStringSwitch(mtk.SIZE_MEDIUM, colornames.Blue,
		lang.Text("gui", "lang_s_label"), config.SupportedLangs())
	langSwitchIndex := s.langSwitch.Find(config.Lang())
	s.langSwitch.SetIndex(langSwitchIndex)
	s.langSwitch.SetOnChangeFunc(s.onSettingsChanged)
	
	return s, nil
}

// Draw draws all menu elements.
func (s *Settings) Draw(win *pixelgl.Window) {
	// Title.
	titlePos :=pixel.V(win.Bounds().Center().X,
		win.Bounds().Max.Y - s.title.Bounds().Size().Y)
	s.title.Draw(win, pixel.IM.Moved(titlePos))
	// Buttons & switches.
	s.resSwitch.Draw(win, pixel.IM.Moved(pixel.V(titlePos.X,
		titlePos.Y - s.resSwitch.Frame().Size().Y)))
	s.langSwitch.Draw(win, pixel.IM.Moved(pixel.V(titlePos.X,
		s.resSwitch.DrawArea().Min.Y - s.langSwitch.Frame().Size().Y)))
	s.backButton.Draw(win, pixel.IM.Moved(pixel.V(titlePos.X,
		s.langSwitch.DrawArea().Min.Y - s.backButton.Frame().Size().Y)))
}

// Update updates all menu elements.
func (s *Settings) Update(win *pixelgl.Window) {
	if s.Opened() {
		s.resSwitch.Update(win)
		s.langSwitch.Update(win)
		s.backButton.Update(win)
	}
}

// Opened checks whether menu should be drawn or not.
func (s *Settings) Opened() bool {
	return s.opened
}

// Show toggles menu visibility.
func (s *Settings) Show(show bool) {
	s.opened = show
}

// Apply applies current settings values.
func (s *Settings) Apply() {
	// Resolution.
	res, ok := s.resSwitch.Value().Value.(pixel.Vec)
	if !ok {
		log.Err.Printf("settings_menu:fail_to_retrive_res_switch_value")
		return
	}
	lang, ok := s.langSwitch.Value().Value.(string)
	if !ok {
		log.Err.Printf("settings_menu:fail_to_retrive_lang_switch_value")
		return
	}
	
	config.SetLang(lang)
	config.SetResolution(res)
}

// Changed checks if any settings value was changed.
func (s *Settings) Changed() bool {
	return s.changed
}

// close closes settings menu and displays message
// about required game restart if settings was changed.
func (s *Settings) close() {
	if s.Changed() {
		msg, err := mtk.NewMessageWindow(mtk.SIZE_SMALL,
			lang.Text("gui", "settings_reset_msg"))
		if err != nil {
			log.Err.Printf("settings:fail_to_create_settings_change_message")
			s.mainmenu.OpenMenu()
			return
		}
		s.mainmenu.ShowMessage(msg)
		s.Apply()
	}
	s.mainmenu.OpenMenu()
}

// closeWithDialog creates settings apply dialog and puts it on
// main menu messages list.
func (s *Settings) closeWithDialog() {
	if s.Changed() {
		dlg, err := mtk.NewDialogWindow(mtk.SIZE_SMALL,
			lang.Text("gui", "settings_save_msg"))
		if err != nil {
			log.Err.Printf("settings:fail_to_create_settings_confirm_dialog")
			s.close() 
			return
		}
		dlg.SetOnAcceptFunc(s.onSettingsApplyAccept)
		dlg.SetOnCancelFunc(s.onSettingsApplyCancel)
		s.mainmenu.ShowMessage(dlg)
	} else {
		s.close()
	}	
}

// Triggered after settings change.
func (s *Settings) onSettingsChanged(sw *mtk.Switch) {
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
