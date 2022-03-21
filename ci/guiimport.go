/*
 * guiimport.go
 *
 * Copyright 2019-2022 Dariusz Sikora <dev@isangeles.pl>
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
	"path/filepath"
	
	"github.com/isangeles/burn"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/data"
)

// guiimport handles guiimport command.
func guiimport(cmd burn.Command) (int, string) {
	if len(cmd.OptionArgs()) < 1 {
		return 2, fmt.Sprintf("%s: no option args", GUIImport)
	}
	switch cmd.OptionArgs()[0] {
	case "hud", "hud-state":
		if len(cmd.Args()) < 1 {
			return 3, fmt.Sprintf("%s: no enought args for: %s",
				GUIImport, cmd.OptionArgs()[0])
		}
		if guiHUD == nil {
			return 3, fmt.Sprintf("%s: no HUD set", GUIImport)
		}
		name := cmd.Args()[0] + data.HUDFileExt
		path := filepath.Join(guiHUD.Game().Conf().Path, config.ModuleGUIDir,
			data.SavesDir, name)
		data, err := data.ImportHUD(path)
		if err != nil {
			return 3, fmt.Sprintf("%s: unable to import HUD: %v",
				GUIImport, err)
		}
		err = guiHUD.Apply(data)
		if err != nil {
			return 3, fmt.Sprintf("%s: unable to apply HUD data: %v",
				GUIImport, err)
		}
		return 0, ""
	default:
		return 2, fmt.Sprintf("%s: invalid option: '%s'", GUIImport,
			cmd.OptionArgs()[0])
	}
}
