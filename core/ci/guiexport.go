/*
 * guiexport.go
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
	
	flameconf "github.com/isangeles/flame/config"
	
	"github.com/isangeles/burn"
	
	"github.com/isangeles/mural/core/data/exp"
)

// guiexport handles guiexport command.
func guiexport(cmd burn.Command) (int, string) {
	if len(cmd.OptionArgs()) < 1 {
		return 2, fmt.Sprintf("%s: no option args", GUIExport)
	}
	switch cmd.OptionArgs()[0] {
	case "avatar":
		if len(cmd.TargetArgs()) < 1 {
			return 3, fmt.Sprintf("%s: no enought target args for: %s",
				GUIExport, cmd.OptionArgs()[0])
		}
		if guiHUD == nil {
			return 3, fmt.Sprintf("%s: no HUD set", GUIExport)
		}
		for _, av := range guiHUD.Camera().Avatars() {
			if av.ID()+"#"+av.Serial() == cmd.TargetArgs()[0] {
				expPath := guiHUD.Game().Module().Conf().CharactersPath()
				err := exp.ExportAvatar(av, expPath)
				if err != nil {
					return 8, fmt.Sprintf("%s: fail to export avatar: %v",
						GUIExport, err)
				}
				return 0, ""
			}
		}
		return 3, fmt.Sprintf("%s: avatar not found: %s", GUIExport,
			cmd.TargetArgs()[0])
	case "gui", "gui-state":
		if len(cmd.Args()) < 1 {
			return 3, fmt.Sprintf("%s: no enought args for: %s",
				GUIExport, cmd.OptionArgs()[0])
		}
		if guiHUD == nil {
			return 3, fmt.Sprintf("%s:no HUD set", GUIExport)
		}
		savName := cmd.Args()[0]
		savDir := flameconf.ModuleSavegamesPath()
		sav := guiHUD.NewGUISave()
		err := exp.ExportGUISave(sav, savDir, savName)
		if err != nil {
			return 3, fmt.Sprintf("%s: fail to save gui state: %v",
				GUIExport, err)
		}
		return 0, ""
	default:
		return 2, fmt.Sprintf("%s: invalid option: '%s'", GUIExport,
			cmd.OptionArgs()[0])
	}
}