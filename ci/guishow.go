/*
 * guishow.go
 *
 * Copyright 2019-2020 Dariusz Sikora <dev@isangeles.pl>
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
	"strings"
	
	"github.com/isangeles/burn"
	
	"github.com/isangeles/mural/config"
)

// guishow handles guishow command.
func guishow(cmd burn.Command) (int, string) {
	if len(cmd.OptionArgs()) < 1 {
		return 2, fmt.Sprintf("%s: no option args", GUIShow)
	}
	switch cmd.OptionArgs()[0] {
	case "version":
		return 0, config.Version
	case "playable-chars":
		if guiMenu == nil {
			return 3, fmt.Sprintf("%s: no gui main menu set", GUIShow)
		}
		out := ""
		for _, c := range guiMenu.PlayableChars() {
			out = fmt.Sprintf("%s%s#%s ", out, c.ID, c.Serial)
		}
		out = strings.TrimSpace(out)
		return 0, out
	default:
		return 2, fmt.Sprintf("%s: invalid option: '%s'", GUIShow,
			cmd.OptionArgs()[0])
	}
}


