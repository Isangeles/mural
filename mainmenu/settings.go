/*
 * settings.go
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

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"

	"github.com/isangeles/flame/data/res/lang"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/log"
)

// Settings struct represents main menu
// settings screen.
type Settings struct {
	mainmenu          *MainMenu
	title             *mtk.Text
	backButton        *mtk.Button
	fullscrSwitch     *mtk.Switch
	resSwitch         *mtk.Switch
	langSwitch        *mtk.Switch
	musicVolumeSwitch *mtk.Switch
	musicMuteSwitch   *mtk.Switch
	opened            bool
	changed           bool
}

// newSettings returns new settings screen instance.
func newSettings(mainmenu *MainMenu) *Settings {
	s := new(Settings)
	s.mainmenu = mainmenu
	// Title.
	titleParams := mtk.Params{
		SizeRaw:  mtk.ConvVec(pixel.V(900, 0)),
		FontSize: mtk.SizeBig,
	}
	s.title = mtk.NewText(titleParams)
	s.title.SetText(lang.Text("settings_menu_title"))
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		FontSize:  mtk.SizeMedium,
		Shape:     mtk.ShapeRectangle,
		MainColor: accentColor,
	}
	s.backButton = mtk.NewButton(buttonParams)
	s.backButton.SetLabel(lang.Text("back_button_label"))
	s.backButton.SetOnClickFunc(s.onBackButtonClicked)
	// Switches.
	switchParams := mtk.Params{
		Size:      mtk.SizeMedium,
		MainColor: mainColor,
	}
	// Fullscreen.
	s.fullscrSwitch = mtk.NewSwitch(switchParams)
	s.fullscrSwitch.SetLabel(lang.Text("settings_fullscr_switch_label"))
	fullscrTrue := mtk.SwitchValue{lang.Text("com_yes"), true}
	fullscrFalse := mtk.SwitchValue{lang.Text("com_no"), false}
	fullscrValues := []mtk.SwitchValue{fullscrFalse, fullscrTrue}
	s.fullscrSwitch.SetValues(fullscrValues...)
	s.fullscrSwitch.SetOnChangeFunc(s.onSettingsSwitchChanged)
	// Resolution.
	s.resSwitch = mtk.NewSwitch(switchParams)
	s.resSwitch.SetLabel(lang.Text("resolution_switch_label"))
	var resValues []mtk.SwitchValue
	for _, res := range config.SupportedResolutions() {
		v := mtk.SwitchValue{fmt.Sprintf("%vx%v", res.X, res.Y), res}
		resValues = append(resValues, v)
	}
	s.resSwitch.SetValues(resValues...)
	s.resSwitch.SetOnChangeFunc(s.onSettingsSwitchChanged)
	// Language.
	s.langSwitch = mtk.NewSwitch(switchParams)
	s.langSwitch.SetLabel(lang.Text("lang_switch_label"))
	langValues := make([]mtk.SwitchValue, len(config.SupportedLangs()))
	for i, l := range config.SupportedLangs() {
		langValues[i] = mtk.SwitchValue{l, l}
	}
	s.langSwitch.SetValues(langValues...)
	s.langSwitch.SetOnChangeFunc(s.onSettingsSwitchChanged)
	// Music volume.
	s.musicVolumeSwitch = mtk.NewSwitch(switchParams)
	s.musicVolumeSwitch.SetLabel(lang.Text("settings_vol_switch_label"))
	volValues := []mtk.SwitchValue{
		mtk.SwitchValue{"-1", -1.0},
		mtk.SwitchValue{lang.Text("settings_vol_sys"), 0.0},
		mtk.SwitchValue{"+1", 1.0},
	}
	s.musicVolumeSwitch.SetValues(volValues...)
	s.musicVolumeSwitch.SetOnChangeFunc(s.onSettingsSwitchChanged)
	// Music mute.
	s.musicMuteSwitch = mtk.NewSwitch(switchParams)
	s.musicMuteSwitch.SetLabel(lang.Text("settings_mute_switch_label"))
	muteTrue := mtk.SwitchValue{lang.Text("com_yes"), true}
	muteFalse := mtk.SwitchValue{lang.Text("com_no"), false}
	muteValues := []mtk.SwitchValue{muteTrue, muteFalse}
	s.musicMuteSwitch.SetValues(muteValues...)
	s.musicMuteSwitch.SetOnChangeFunc(s.onSettingsSwitchChanged)
	return s
}

// Draw draws all menu elements.
func (s *Settings) Draw(win *pixelgl.Window) {
	// Title.
	titlePos := pixel.V(win.Bounds().Center().X, win.Bounds().Max.Y-s.title.Size().Y)
	s.title.Draw(win, mtk.Matrix().Moved(titlePos))
	// Switches.
	fullscrSwitchPos := mtk.BottomOf(s.title.DrawArea(), s.fullscrSwitch.Size(), 50)
	resSwitchPos := mtk.BottomOf(s.fullscrSwitch.DrawArea(), s.resSwitch.Size(), 30)
	langSwitchPos := mtk.BottomOf(s.resSwitch.DrawArea(), s.langSwitch.Size(), 30)
	mVolSwitchPos := mtk.BottomOf(s.resSwitch.DrawArea(), s.musicVolumeSwitch.Size(), 30)
	mMuteSwitchPos := mtk.BottomOf(s.musicVolumeSwitch.DrawArea(), s.musicMuteSwitch.Size(), 30)
	s.fullscrSwitch.Draw(win, mtk.Matrix().Moved(fullscrSwitchPos))
	s.resSwitch.Draw(win, mtk.Matrix().Moved(resSwitchPos))
	s.langSwitch.Draw(win, mtk.Matrix().Moved(langSwitchPos))
	s.musicVolumeSwitch.Draw(win, mtk.Matrix().Moved(mVolSwitchPos))
	s.musicMuteSwitch.Draw(win, mtk.Matrix().Moved(mMuteSwitchPos))
	// Buttons.
	backButtonPos := mtk.BottomOf(s.musicMuteSwitch.DrawArea(), s.backButton.Size(), 30)
	s.backButton.Draw(win, mtk.Matrix().Moved(backButtonPos))
}

// Update updates all menu elements.
func (s *Settings) Update(win *mtk.Window) {
	s.fullscrSwitch.Update(win)
	s.resSwitch.Update(win)
	s.langSwitch.Update(win)
	s.musicVolumeSwitch.Update(win)
	s.musicMuteSwitch.Update(win)
	s.backButton.Update(win)
}

// Opened checks whether menu should be drawn or not.
func (s *Settings) Opened() bool {
	return s.opened
}

// Show shows menu.
func (s *Settings) Show() {
	s.opened = true
	s.updateValues()
}

// Hide hides menu.
func (s *Settings) Hide() {
	s.opened = false
}

// Apply applies current settings values.
func (s *Settings) Apply() {
	// Fullscreen.
	fscr, ok := s.fullscrSwitch.Value().Value.(bool)
	if !ok {
		log.Err.Printf("settings menu: fail to retrive fullscreen switch value")
		return
	}
	config.Fullscreen = fscr
	// Resolution.
	res, ok := s.resSwitch.Value().Value.(pixel.Vec)
	if !ok {
		log.Err.Printf("settings menu: fail to retrive res switch value")
		return
	}
	config.Resolution = res
	// Language.
	lang, ok := s.langSwitch.Value().Value.(string)
	if !ok {
		log.Err.Printf("settings menu: fail to retrive lang switch value")
		return
	}
	config.Lang = lang
	// Music volume.
	mVol, ok := s.musicVolumeSwitch.Value().Value.(float64)
	if !ok {
		log.Err.Printf("settings menu: fail to retrive music volume switch value")
		return
	}
	config.MusicVolume = mVol
	// Music mute.
	mMute, ok := s.musicMuteSwitch.Value().Value.(bool)
	if !ok {
		log.Err.Printf("settings menu: fail to retrive music mute switch value")
		return
	}
	config.MusicMute = mMute
}

// Changed checks if any settings value was changed.
func (s *Settings) Changed() bool {
	return s.changed
}

// updateValues values of all settings elements.
func (s *Settings) updateValues() {
	fullscrIndex := s.fullscrSwitch.Find(config.Fullscreen)
	s.fullscrSwitch.SetIndex(fullscrIndex)
	resIndex := s.resSwitch.Find(config.Resolution)
	s.resSwitch.SetIndex(resIndex)
	langIndex := s.langSwitch.Find(config.Lang)
	s.langSwitch.SetIndex(langIndex)
	mVolIndex := s.musicVolumeSwitch.Find(config.MusicVolume)
	s.musicVolumeSwitch.SetIndex(mVolIndex)
	mMuteIndex := s.musicMuteSwitch.Find(config.MusicMute)
	s.musicMuteSwitch.SetIndex(mMuteIndex)
}

// close closes settings menu and displays message
// about required game restart if settings was changed.
func (s *Settings) close() {
	if s.Changed() {
		msg := lang.Text("settings_reset_msg")
		s.mainmenu.ShowMessage(msg)
		s.Apply()
	}
	s.mainmenu.OpenMenu()
}

// closeWithDialog creates settings apply dialog and puts it on
// main menu messages list.
func (s *Settings) closeWithDialog() {
	if s.Changed() {
		dlgParams := mtk.Params{
			Size:      mtk.SizeBig,
			FontSize:  mtk.SizeMedium,
			MainColor: mainColor,
			SecColor:  accentColor,
			Info:      lang.Text("settings_save_msg"),
		}
		dlg := mtk.NewDialogWindow(dlgParams)
		dlg.SetAcceptLabel(lang.Text("accept_button_label"))
		dlg.SetCancelLabel(lang.Text("cancel_button_label"))
		dlg.SetOnAcceptFunc(s.onSettingsApplyAccept)
		dlg.SetOnCancelFunc(s.onSettingsApplyCancel)
		s.mainmenu.ShowMessageWindow(dlg)
	} else {
		s.close()
	}
}

// Triggered after one of settings switch was changed.
func (s *Settings) onSettingsSwitchChanged(sw *mtk.Switch, old, new *mtk.SwitchValue) {
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
