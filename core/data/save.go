/*
 * save.go
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

package data

import (
	"bufio"
	"fmt"
	"path/filepath"
	"os"

	"github.com/isangeles/mural/core/data/save"
	"github.com/isangeles/mural/core/data/parsexml"
)

var (
	GUISAVE_FILE_EXT = ".savegui"
)

// SaveGUI saves GUI state to file with specified name
// in directory with specified path.
func SaveGUI(gui *save.GUISave, dirPath, saveName string) error {
	xml, err := parsexml.MarshalGUISave(gui)
	if err != nil {
		return fmt.Errorf("fail_to_marshal_gui_save:%v",
			err)
	}
	filePath := filepath.FromSlash(dirPath + "/" + saveName +
		GUISAVE_FILE_EXT)
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("fail_to_create_save_file:%v",
			err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(xml)
	w.Flush()
	return nil
}
