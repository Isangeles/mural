/*
 * guiman.go
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

package ci

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/faiface/pixel"

	flameconf "github.com/isangeles/flame/config"
	"github.com/isangeles/flame/cmd/burn"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/data/exp"
	"github.com/isangeles/mural/core/data/imp"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/mainmenu"
	"github.com/isangeles/mural/hud"
)

var (
	gui_mmenu *mainmenu.MainMenu
	gui_hud   *hud.HUD
	gui_music *mtk.AudioPlayer
)

// SetMainMenu sets specified main menu as main
// menu for guiman to manage.
func SetMainMenu(menu *mainmenu.MainMenu) {
	gui_mmenu = menu
}

// SetHUD sets specified HUD as HUD for
// guiman to manage.
func SetHUD(h *hud.HUD) {
	gui_hud = h
}

// SetMusicPlayer sets specified audio player as
// player for guiman to manage.
func SetMusicPlayer(p *mtk.AudioPlayer) {
	gui_music = p
}

// handleGUICommand handles guiman tool commands.
// Returns response code and message.
func handleGUICommand(cmd burn.Command) (int, string) {
	if len(cmd.OptionArgs()) < 1 {
		return 3, fmt.Sprintf("%s:no_option_args", GUI_MAN)
	}

	switch cmd.OptionArgs()[0] {
	case "version":
		return 0, config.VERSION
	case "set":
		return setGUIOption(cmd)
	case "show":
		return showGUIOption(cmd)
	case "export", "save":
		return exportGUIOption(cmd)
	case "import", "load":
		return importGUIOption(cmd)
	case "start", "play":
		return startGUIOption(cmd)
	case "audio":
		return audioGUIOption(cmd)
	case "exit":
		if gui_hud != nil {
			gui_hud.Exit()
			return 0, ""
		}
		if gui_mmenu != nil {
			gui_mmenu.Exit()
			return 0, ""
		}
		return 5, fmt.Sprintf("no main menu or HUD set")
	default:
		return 4, fmt.Sprintf("%s:no_such_option:%s", GUI_MAN,
			cmd.OptionArgs()[0])
	}
}

// setGUIOption Handles set coptions for guiman commands.
func setGUIOption(cmd burn.Command) (int, string) {
	if len(cmd.TargetArgs()) < 1 {
		return 5, fmt.Sprintf("%s:no_enought_target_args_for:%s", GUI_MAN,
			cmd.OptionArgs()[0])
	}

	switch cmd.TargetArgs()[0] {
	case "resolution":
		if len(cmd.Args()) < 1 {
                        return 7, fmt.Sprintf("%s:no_enought_args_for:%s", GUI_MAN,
				cmd.TargetArgs()[0])
                }
		
		resInput := cmd.Args()[0]
		resX, err := strconv.ParseFloat(strings.Split(resInput, "x")[0], 64)
		resY, err := strconv.ParseFloat(strings.Split(resInput, "x")[1], 64)
		if err != nil {
			return 8, fmt.Sprintf("%s:invalid_input:%s", GUI_MAN,
				cmd.OptionArgs()[0])
		}
		config.SetResolution(pixel.V(resX, resY))
		return 0, ""
	case "fow":
		if len(cmd.Args()) < 1 {
                        return 7, fmt.Sprintf("%s:no_enought_args_for:%s", GUI_MAN,
				cmd.TargetArgs()[0])
                }
		fow := cmd.Args()[0] == "on"
		config.SetMapFOW(fow)
		return 0, ""
	default:
		return 6, fmt.Sprintf("%s:no_vaild_target_for_%s:'%s'", GUI_MAN,
			cmd.OptionArgs()[0], cmd.TargetArgs()[0])
	}
}

// showGUIOption handles 'show' option for guiman tool.
func showGUIOption(cmd burn.Command) (int, string) {
	if len(cmd.TargetArgs()) < 1 {
		return 5, fmt.Sprintf("%s:no_enought_target_args_for:%s", GUI_MAN,
			cmd.OptionArgs()[0])
	}

	switch cmd.TargetArgs()[0] {
	case "playable-chars":
		if gui_mmenu == nil {
			return 7, fmt.Sprintf("%s:no gui main menu set", GUI_MAN)
		}
		out := ""
		for _, c := range gui_mmenu.PlayableChars {
			out = out + fmt.Sprintf("%s[%s]", c.Name(), c.SerialID()) +
				","
		}
		return 0, out
	default:
		return 6, fmt.Sprintf("%s:no_vaild_target_for_%s:'%s'", GUI_MAN,
			cmd.OptionArgs()[0], cmd.TargetArgs()[0])
	}
}

// exportGUIOption handles 'export' option for guiman tool.
func exportGUIOption(cmd burn.Command) (int, string) {
	if len(cmd.TargetArgs()) < 1 {
		return 5, fmt.Sprintf("%s:no_enought_target_args_for:%s", GUI_MAN,
			cmd.OptionArgs()[0])
	}

	switch cmd.TargetArgs()[0] {
	case "avatar":
		if len(cmd.TargetArgs()) < 2 {
			return 7, fmt.Sprintf("%s:no_enought_target_args_for:%s",
				GUI_MAN, cmd.TargetArgs()[0])
		}
		if gui_hud == nil {
			return 7, fmt.Sprintf("%s:no HUD set", GUI_MAN)
		}
		for _, av := range gui_hud.AreaAvatars() {
			if av.SerialID() == cmd.TargetArgs()[1] {
				err := exp.ExportAvatar(av,
					gui_hud.Game().Module().Conf().CharactersPath())
				if err != nil {
					return 8, fmt.Sprintf("%s:fail_to_export_avatar:%v",
						GUI_MAN, err)
				}
				return 0, ""
			}
		}
		return 7, fmt.Sprintf("%s:avatar_not_found:%s", GUI_MAN,
			cmd.TargetArgs()[1])
	case "gui-state":
		if len(cmd.Args()) < 1 {
			return 7, fmt.Sprintf("%s:no_enought_args_for:%s",
				GUI_MAN, cmd.TargetArgs()[0])
		}
		if gui_hud == nil {
			return 7, fmt.Sprintf("%s:no HUD set", GUI_MAN)
		}
		savName := cmd.Args()[0]
		savDir := flameconf.ModuleSavegamesPath()
		sav := gui_hud.NewGUISave()
		err := exp.ExportGUISave(sav, savDir, savName)
		if err != nil {
			return 8, fmt.Sprintf("%s:fail_to_save_gui_state:%v",
				GUI_MAN, err)
		}
		return 0, ""
	default:
		return 6, fmt.Sprintf("%s:no_vaild_target_for_%s:'%s'", GUI_MAN,
			cmd.OptionArgs()[0], cmd.TargetArgs()[0])
	}
}

// importGUIOption handles import option for guiman.
func importGUIOption(cmd burn.Command) (int, string) {
	if len(cmd.TargetArgs()) < 1 {
		return 5, fmt.Sprintf("%s:no_enought_target_args_for:%s", GUI_MAN,
			cmd.OptionArgs()[0])
	}
	switch cmd.TargetArgs()[0] {
	case "gui-state":
		if len(cmd.Args()) < 1 {
			return 7, fmt.Sprintf("%s:no_enought_args_for:%s",
				GUI_MAN, cmd.TargetArgs()[0])
		}
		if gui_hud == nil {
			return 7, fmt.Sprintf("%s:no HUD set", GUI_MAN)
		}
		savName := cmd.Args()[0]
		savDir := flameconf.ModuleSavegamesPath()
		save, err := imp.ImportGUISave(gui_hud.Game(), savDir, savName)
		if err != nil {
			return 9, fmt.Sprintf("%s:fail_to_load_save_file:%v",
				GUI_MAN, err)
		}
		err = gui_hud.LoadGUISave(save)
		if err != nil {
			return 9, fmt.Sprintf("%s:fail_to_load_gui_state_save:%v",
				GUI_MAN, err)
		}
		return 0, ""
	default:
		return 6, fmt.Sprintf("%s:no_vaild_target_for_%s:'%s'", GUI_MAN,
			cmd.OptionArgs()[0], cmd.TargetArgs()[0])
	}
}

// startGUIOption handles import option for guiman.
func startGUIOption(cmd burn.Command) (int, string) {
	if len(cmd.TargetArgs()) < 1 {
		return 5, fmt.Sprintf("%s:no_enought_target_args_for:%s", GUI_MAN,
			cmd.OptionArgs()[0])
	}
	switch cmd.TargetArgs()[0] {
	default:
		return 6, fmt.Sprintf("%s:no_vaild_target_for_%s:'%s'", GUI_MAN,
			cmd.OptionArgs()[0], cmd.TargetArgs()[0])
	}
}

// audioGUIOption handles audio option for guiman.
func audioGUIOption(cmd burn.Command) (int, string) {
	if gui_music == nil {
		return 7, fmt.Sprintf("%s:%s:no audio player set",
			GUI_MAN, cmd.TargetArgs()[0])
	}
	if len(cmd.Args()) < 1 {
		return 7, fmt.Sprintf("%s:no_enought_args_for:%s",
			GUI_MAN, cmd.TargetArgs()[0])
	}
	switch cmd.Args()[0] {
	case "play-music":
		gui_music.PlayMusic()
		return 0, ""
	case "stop-music":
		gui_music.StopMusic()
		return 0, ""
	case "next":
		gui_music.Next()
		return 0, ""
	case "prev":
		gui_music.Prev()
		return 0, ""
	default:
		return 6, fmt.Sprintf("%s:no_vaild_target_for_%s:'%s'", GUI_MAN,
			cmd.OptionArgs()[0], cmd.TargetArgs()[0])
	}
}
