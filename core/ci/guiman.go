/*
 * guiman.go
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

package ci

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/faiface/pixel"

	flameci "github.com/isangeles/flame/cmd/ci"

	"github.com/isangeles/mural/config"
)

// handleGUICommand handles guiman tool commands.
// Returns response code and message.
func handleGUICommand(cmd flameci.Command) (int, string) {
	if len(cmd.OptionArgs()) < 1 {
		return 3, fmt.Sprintf("%s:no_option_args", GUI_MAN)
	}

	switch cmd.OptionArgs()[0] {
	case "version":
		return 0, config.VERSION
	case "set":
		return setGUIOption(cmd)
	default:
		return 4, fmt.Sprintf("%s:no_such_option:%s", GUI_MAN,
			cmd.OptionArgs()[0])
	}
}

// setGUIOption Handles set coptions for guiman commands.
func setGUIOption(cmd flameci.Command) (int, string) {
	if len(cmd.TargetArgs()) < 1 {
		return 5, fmt.Sprintf("%s:no_enought_target_args_for:%s", GUI_MAN,
			cmd.OptionArgs()[0])
	}

	switch cmd.TargetArgs()[0] {
	case "resolution":
		if len(cmd.Args()) < 1 {
                        return 7, fmt.Sprintf("%s:no_enought_args_for:%s", GUI_MAN,
				cmd.OptionArgs()[1])
                }
		
		resInput := cmd.Args()[0]
		resX, err := strconv.ParseFloat(strings.Split(resInput,
			"x")[0], 64)
		resY, err := strconv.ParseFloat(strings.Split(resInput,
			"x")[1], 64)
		if err != nil {
			return 8, fmt.Sprintf("%s:invalid_input:%s", GUI_MAN,
				cmd.OptionArgs()[0])
		}
		config.SetResolution(pixel.V(resX, resY))
		return 0, ""
	default:
		return 6, fmt.Sprintf("%s:no_vaild_target_for_%s:'%s'", GUI_MAN,
			cmd.OptionArgs()[0], cmd.TargetArgs()[0])
	}
}

