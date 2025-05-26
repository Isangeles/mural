/*
 * hud.go
 *
 * Copyright 2019-2025 Dariusz Sikora <ds@isangeles.dev>
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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/isangeles/mural/data/res"
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
	buf, err := io.ReadAll(file)
	if err != nil {
		return data, fmt.Errorf("unable to read data file: %v",
			err)
	}
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return data, fmt.Errorf("unable to unmarshal data: %v",
			err)
	}
	return data, nil
}

// ExportHUD exports HUD data to file with specified path.
func ExportHUD(hud res.HUDData, path string) error {
	// Marshal GUI.
	hud.Name = filepath.Base(path)
	data, err := json.Marshal(&hud)
	if err != nil {
		return fmt.Errorf("unable to marshal data: %v",
			err)
	}
	// Create save file.
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return fmt.Errorf("unable to create directory: %v",
			err)
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create file: %v",
			err)
	}
	defer file.Close()
	// Write save.
	writer := bufio.NewWriter(file)
	writer.Write(data)
	writer.Flush()
	return nil
}
