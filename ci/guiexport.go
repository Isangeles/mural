/*
 * guiexport.go
 *
 * Copyright 2019-2023 Dariusz Sikora <ds@isangeles.dev>
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
		for _, av := range guiHUD.Camera().Area().Avatars() {
			if av.ID()+"#"+av.Serial() == cmd.TargetArgs()[0] {
				avatarsPath := filepath.Join(config.GUIPath, "avatars", av.ID())
				err := data.ExportAvatars(avatarsPath, av.Data())
				if err != nil {
					return 3, fmt.Sprintf("unable to export avatar: %v", err)
				}
				return 0, ""
			}
		}
		return 3, fmt.Sprintf("%s: avatar not found: %s", GUIExport,
			cmd.TargetArgs()[0])
	case "hud", "hud-state":
		if len(cmd.Args()) < 1 {
			return 3, fmt.Sprintf("%s: no enought args for: %s",
				GUIExport, cmd.OptionArgs()[0])
		}
		if guiHUD == nil {
			return 3, fmt.Sprintf("%s:no HUD set", GUIExport)
		}
		name := cmd.Args()[0]
		path := filepath.Join(config.GUIPath, data.SavesDir, name + data.HUDFileExt)
		hudData := guiHUD.Data()
		err := data.ExportHUD(hudData, path)
		if err != nil {
			return 3, fmt.Sprintf("%s: unable to export HUD state: %v",
				GUIExport, err)
		}
		return 0, ""
	default:
		return 2, fmt.Sprintf("%s: invalid option: '%s'", GUIExport,
			cmd.OptionArgs()[0])
	}
}
