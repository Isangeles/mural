/*
 * guiset.go
 *
 * Copyright 2019 Dariusz Sikora <dev@isangeles.pl>
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
	
	"github.com/isangeles/burn"
	
	"github.com/isangeles/mural/config"
)

// guiset handles guiset command.
func guiset(cmd burn.Command) (int, string) {
	if len(cmd.OptionArgs()) < 1 {
		return 2, fmt.Sprintf("%s: no option args", GUISet)
	}
	switch cmd.OptionArgs()[0] {
	case "resolution":
		if len(cmd.Args()) < 1 {
                        return 3, fmt.Sprintf("%s: no enought args for: %s", GUISet,
				cmd.OptionArgs()[0])
                }
		resInput := cmd.Args()[0]
		resX, err := strconv.ParseFloat(strings.Split(resInput, "x")[0], 64)
		resY, err := strconv.ParseFloat(strings.Split(resInput, "x")[1], 64)
		if err != nil {
			return 3, fmt.Sprintf("%s: invalid input: '%s'", GUISet,
				cmd.Args()[0])
		}
		config.SetResolution(pixel.V(resX, resY))
		return 0, ""
	case "fow":
		if len(cmd.Args()) < 1 {
                        return 3, fmt.Sprintf("%s: no enought args for: %s", GUISet,
				cmd.OptionArgs()[0])
                }
		fow := cmd.Args()[0] == "on"
		config.SetMapFOW(fow)
		return 0, ""
	case "exit":
		if guiHUD != nil {
			guiHUD.Exit()
			return 0, ""
		}
		if guiMenu != nil {
			guiMenu.Exit()
			return 0, ""
		}
		return 3, fmt.Sprintf("no main menu or HUD set")
	default:
		return 2, fmt.Sprintf("%s: invalid option: '%s'", GUISet,
			cmd.OptionArgs()[0])
	}
}
