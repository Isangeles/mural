/*
 * hud.go
 *
 * Copyright 2019-2021 Dariusz Sikora <dev@isangeles.pl>
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
	"encoding/xml"
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/log"
)

// ImportHUD imports HUD data file from specified path.
func ImportHUD(path string) (res.HUDData, error) {
	var data res.HUDData
	file, err := os.Open(path)
	if err != nil {
		return data, fmt.Errorf("unable to open data file: %v",
			err)
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return data, fmt.Errorf("unable to read data file: %v",
			err)
	}
	err = xml.Unmarshal(buf, &data)
	if err != nil {
		return data, fmt.Errorf("unable to unmarshal data: %v",
			err)
	}
	return data, nil
}

// ImportHUDDir imports all HUD data files in directory with
// specified path.
func ImportHUDDir(dirPath string) ([]res.HUDData, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir: %v", err)
	}
	data := make([]res.HUDData, 0)
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), HUDFileExt) {
			continue
		}
		path := filepath.Join(dirPath, f.Name())
		hud, err := ImportHUD(path)
		if err != nil {
			log.Err.Printf("data: import hud dir: unable to import data file: %v",
				err)
			continue
		}
		data = append(data, hud)
	}
	return data, nil
}

// ExportHUD exports HUD data to file with specified name
// in directory with specified path.
func ExportHUD(hud res.HUDData, path string) error {
	// Marshal GUI.
	hud.Name = filepath.Base(path)
	xml, err := xml.Marshal(&hud)
	if err != nil {
		return fmt.Errorf("unable to marshal hud: %v",
			err)
	}
	// Create save file.
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return fmt.Errorf("unable to create hud directory: %v",
			err)
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create hud file: %v",
			err)
	}
	defer file.Close()
	// Write save.
	w := bufio.NewWriter(file)
	w.Write(xml)
	w.Flush()
	log.Dbg.Printf("HUD exported to: %s", path)
	return nil
}
