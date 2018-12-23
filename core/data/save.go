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
	"os"
	"path/filepath"
	"strings"

	flameparsexml "github.com/isangeles/flame/core/data/parsexml"
	
	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/core/data/save"
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

// LoadGUISave loads GUI state from file with specified name
// in directory with specified path.
func LoadGUISave(dirPath, saveName string) (*save.GUISave, error) {
	if !strings.HasSuffix(saveName, GUISAVE_FILE_EXT) {
		saveName = saveName + GUISAVE_FILE_EXT
	}
	filePath := filepath.FromSlash(dirPath + "/" + saveName)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_save_file:%v",
			err)
	}
	xmlSave, err := parsexml.UnmarshalGUISave(f)
	if err != nil {
		return nil, fmt.Errorf("fail_to_unmarshal_save_data:%v",
			err)
	}
	save, err := buildXMLGUISave(xmlSave)
	if err != nil {
		return nil, fmt.Errorf("fail_to_build_save:%v",
			err)
	}
	return save, nil
}

// buildGUISave builds GUI save from specified XML data.
func buildXMLGUISave(xmlSave *parsexml.GUISaveXML) (*save.GUISave, error) {
	save := new(save.GUISave)
	// Camera position.
	camX, camY, err := flameparsexml.UnmarshalPosition(xmlSave.Camera.Position)
	if err != nil {
		return nil, fmt.Errorf("fail_to_unmarshal_camera_position:%v",
			err)
	}
	save.CameraPosX, save.CameraPosY = camX, camY
	return save, nil
}
