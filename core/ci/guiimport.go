/*
 * guiimport.go
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
	
	"github.com/isangeles/burn"
	
	flameconf "github.com/isangeles/flame/config"
	
	"github.com/isangeles/mural/core/data/imp"
)

// guiimport handles guiimport command.
func guiimport(cmd burn.Command) (int, string) {
	if len(cmd.OptionArgs()) < 1 {
		return 2, fmt.Sprintf("%s: no option args", GUIImport)
	}
	switch cmd.OptionArgs()[0] {
	case "gui", "gui-state":
		if len(cmd.Args()) < 1 {
			return 3, fmt.Sprintf("%s: no enought args for: %s",
				GUIImport, cmd.OptionArgs()[0])
		}
		if guiHUD == nil {
			return 3, fmt.Sprintf("%s: no HUD set", GUIImport)
		}
		savName := cmd.Args()[0]
		savDir := flameconf.ModuleSavegamesPath()
		save, err := imp.ImportGUISave(savDir, savName)
		if err != nil {
			return 3, fmt.Sprintf("%s: fail to load save file: %v",
				GUIImport, err)
		}
		err = guiHUD.LoadGUISave(save)
		if err != nil {
			return 3, fmt.Sprintf("%s: fail to load gui state save: %v",
				GUIImport, err)
		}
		return 0, ""
	default:
		return 2, fmt.Sprintf("%s: invalid option: '%s'", GUIImport,
			cmd.OptionArgs()[0])
	}
}
