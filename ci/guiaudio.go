/*
 * guiaudio.go
 *
 * Copyright 2019-2024 Dariusz Sikora <ds@isangeles.dev>
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

	"github.com/isangeles/burn"
	"github.com/isangeles/mtk"
)

// guiaudio handles guiaudio command.
func guiaudio(cmd burn.Command) (int, string) {
	if len(cmd.OptionArgs()) < 1 {
		return 2, fmt.Sprintf("%s: no option args", GUIAudio)
	}
	if guiMenu == nil {
		return 2, fmt.Sprintf("%s: no main menu set", GUIAudio)
	}
	switch cmd.OptionArgs()[0] {
	case "play":
		guiMenu.Music().ResumePlaylist()
		return 0, ""
	case "stop":
		guiMenu.Music().Stop()
		return 0, ""
	case "next":
		guiMenu.Music().Stop()
		guiMenu.Music().SetPlayIndex(guiMenu.Music().PlayIndex() + 1)
		guiMenu.Music().ResumePlaylist()
		return 0, ""
	case "prev":
		guiMenu.Music().Stop()
		guiMenu.Music().SetPlayIndex(guiMenu.Music().PlayIndex() - 1)
		guiMenu.Music().ResumePlaylist()
		return 0, ""
	case "music-volume":
		out := fmt.Sprintf("%f", guiMenu.Music().Volume())
		return 0, out
	case "set-music-volume":
		if len(cmd.Args()) < 1 {
			return 3, fmt.Sprintf("%s: not enought args for: %s",
				GUIAudio, cmd.OptionArgs()[0])
		}
		vol, err := strconv.ParseFloat(cmd.Args()[0], 64)
		if err != nil {
			return 3, fmt.Sprintf("%s: invalid argument: '%s': %v",
				GUIAudio, cmd.Args()[0], err)
		}
		guiMenu.Music().SetVolume(vol)
		return 0, ""
	case "set-music-mute":
		if len(cmd.Args()) < 1 {
			return 3, fmt.Sprintf("%s: not enought args for: %s",
				GUIAudio, cmd.OptionArgs()[0])
		}
		mute := cmd.Args()[0] == "true"
		guiMenu.Music().SetMute(mute)
		return 0, ""
	case "effects-volume":
		out := fmt.Sprintf("%f", mtk.Audio().Volume())
		return 0, out
	case "set-effects-volume":
		if len(cmd.Args()) < 1 {
			return 3, fmt.Sprintf("%s: not enought args for: %s",
				GUIAudio, cmd.OptionArgs()[0])
		}
		vol, err := strconv.ParseFloat(cmd.Args()[0], 64)
		if err != nil {
			return 3, fmt.Sprintf("%s: invalid argument: '%s': %v",
				GUIAudio, cmd.Args()[0], err)
		}
		mtk.Audio().SetVolume(vol)
		return 0, ""
	case "set-effects-mute":
		if len(cmd.Args()) < 1 {
			return 3, fmt.Sprintf("%s: not enought args for: %s",
				GUIAudio, cmd.OptionArgs()[0])
		}
		mute := cmd.Args()[0] == "true"
		mtk.Audio().SetMute(mute)
		return 0, ""
	default:
		return 2, fmt.Sprintf("%s: invalid option: '%s'", GUIAudio,
			cmd.OptionArgs()[0])
	}
}
