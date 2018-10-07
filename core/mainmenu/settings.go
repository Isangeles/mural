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
	"strings"
	"strconv"
	
	"golang.org/x/image/colornames"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"

	//"github.com/isangeles/flame"
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/log"
)

// Settings struct represents main menu
// settings screen.
type Settings struct {
	title      *text.Text
	backButton *mtk.Button
	resSwitch  *mtk.Switch
	open       bool
	changed    bool
}

// newSettings returns new settings screen
// instance.
func newSettings() (*Settings, error) {
	s := new(Settings)
	// Title.
	font := data.MainFontBig()
	atlas := text.NewAtlas(font, text.ASCII)
	s.title = text.New(pixel.V(0, 0), atlas)
	fmt.Fprintf(s.title, lang.Text("gui", "settings_menu_title"))
	// Buttons & switches.
	s.backButton = mtk.NewButton(mtk.SIZE_MEDIUM, mtk.SHAPE_RECTANGLE, colornames.Red,
		lang.Text("gui", "back_b_label"))
	var resSwitchValues []string
	var resSwitchIndex int
	for i, res := range config.SupportedResolutions() {
		resSwitchValues = append(resSwitchValues,
			fmt.Sprintf("%vx%v", res.X, res.Y))
		if res == config.Resolution() {
			resSwitchIndex = i
		}
	}
	s.resSwitch = mtk.NewSwitch(mtk.SIZE_MEDIUM, colornames.Blue,
		lang.Text("gui", "resolution_s_label"), resSwitchValues)
	s.resSwitch.SetIndex(resSwitchIndex)
	s.resSwitch.SetOnChangeFunc(s.onSettingsChanged)
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
	s.backButton.Draw(win, pixel.IM.Moved(pixel.V(titlePos.X,
		s.resSwitch.DrawArea().Min.Y - s.backButton.Frame().Size().Y)))
}

// Update updates all menu elements.
func (s *Settings) Update(win *pixelgl.Window) {
	if s.open {
		s.resSwitch.Update(win)
		s.backButton.Update(win)
	}
}

// Open checks whether menu should be drawn or not.
func (s *Settings) Open() bool {
	return s.open
}

// Show toggles menu visibility.
func (s *Settings) Show(show bool) {
	s.open = show
}

// Sets specified function as back button on-click
// callback function.
func (s *Settings) SetOnBackButtonClickedFunc(f func(b *mtk.Button)) {
	s.backButton.SetOnClickFunc(f)
}

// Apply applies current settings values.
func (s *Settings) Apply() {
	// Resolution.
	resText := strings.Split(s.resSwitch.Value(), "x")
	resValueX, err := strconv.ParseFloat(resText[0], 64)
	if err != nil {
		resValueX = 1920
		log.Err.Printf("settings_menu:fail_to_conv_res_switch_value_x:%v", err)
	}
	resValueY, err := strconv.ParseFloat(resText[1], 64)
	if err != nil {
		resValueY = 1080
		log.Err.Printf("settings_menu:fail_to_conv_res_switch_value_y:%v", err)
	}
	resValue := pixel.V(resValueX, resValueY)
	config.SetResolution(resValue)
}

// Changed checks if any settings value was changed.
func (s *Settings) Changed() bool {
	return s.changed
}

// Triggered after settings change.
func (s *Settings) onSettingsChanged(sw *mtk.Switch) {
	s.changed = true
}
