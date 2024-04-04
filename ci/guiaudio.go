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
)

// guiaudio handles guiaudio command.
func guiaudio(cmd burn.Command) (int, string) {
	if len(cmd.OptionArgs()) < 1 {
		return 2, fmt.Sprintf("%s: no option args", GUIAudio)
	}
	if guiMusic == nil {
		return 2, fmt.Sprintf("%s: no music player set", GUIAudio)
	}
	switch cmd.OptionArgs()[0] {
	case "play-music":
		guiMusic.ResumePlaylist()
		return 0, ""
	case "stop-music":
		guiMusic.StopPlaylist()
		return 0, ""
	case "next":
		guiMusic.StopPlaylist()
		guiMusic.SetPlayIndex(guiMusic.PlayIndex()+1)
		guiMusic.ResumePlaylist()
		return 0, ""
	case "prev":
		guiMusic.StopPlaylist()
		guiMusic.SetPlayIndex(guiMusic.PlayIndex()-1)
		guiMusic.ResumePlaylist()
		return 0, ""
	case "volume":
		out := fmt.Sprintf("%f", guiMusic.Volume())
		return 0, out
	case "set-volume":
		if len(cmd.Args()) < 1 {
			return 3, fmt.Sprintf("%s: not enought args for: %s",
				GUIAudio, cmd.OptionArgs()[0])
		}
		vol, err := strconv.ParseFloat(cmd.Args()[0], 64)
		if err != nil {
			return 3, fmt.Sprintf("%s: invalid argument: '%s': %v",
				GUIAudio, cmd.Args()[0], err)
		}
		guiMusic.SetVolume(vol)
		return 0, ""
	case "set-mute":
		if len(cmd.Args()) < 1 {
			return 3, fmt.Sprintf("%s: not enought args for: %s",
				GUIAudio, cmd.OptionArgs()[0])
		}
		mute := cmd.Args()[0] == "true"
		guiMusic.SetMute(mute)
		return 0, ""
	default:
		return 2, fmt.Sprintf("%s: invalid option: '%s'", GUIAudio,
			cmd.OptionArgs()[0])
	}
}
