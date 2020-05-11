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

package data

import (
	"bufio"
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/log"
)

// ImportGUISave imports GUI save file from specified path.
func ImportGUISave(path string) (*res.GUISave, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open save file: %v",
			err)
	}
	defer file.Close()
	save, err := parsexml.UnmarshalGUISave(file)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal save data: %v",
			err)
	}
	return save, nil
}

// ImportsGUISavesDir imports all saved GUIs from save files in
// directory with specified path.
func ImportGUISavesDir(dirPath string) ([]*res.GUISave, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir: %v", err)
	}
	saves := make([]*res.GUISave, 0)
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), SaveFileExt) {
			continue
		}
		savePath := filepath.Join(dirPath, f.Name())
		sav, err := ImportGUISave(savePath)
		if err != nil {
			log.Err.Printf("data saves import: unable to load save: %v",
				err)
			continue
		}
		saves = append(saves, sav)
	}
	return saves, nil
}

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
