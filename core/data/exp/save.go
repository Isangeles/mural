/*
 * save.go
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

package exp

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/log"
)

// ExportGUISave saves GUI state to file with specified name
// in directory with specified path.
func ExportGUISave(gui *res.GUISave, path string) error {
	gui.Name = filepath.Base(path)
	xml, err := parsexml.MarshalGUISave(gui)
	if err != nil {
		return fmt.Errorf("unable to marshal save: %v",
			err)
	}
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return fmt.Errorf("unable to create save directory: %v",
			err)
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create save file: %v",
			err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	w.WriteString(xml)
	w.Flush()
	log.Dbg.Printf("gui state saved in: %s", path)
	return nil
}
