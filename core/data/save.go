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
	"io/ioutil"
	"path/filepath"
	"strings"

	flamecore "github.com/isangeles/flame/core"
	flameparsexml "github.com/isangeles/flame/core/data/parsexml"
	
	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/core/data/save"
	"github.com/isangeles/mural/log"
)

var (
	SAVEGUI_FILE_EXT = ".savegui"
)

// ExportGUISave saves GUI state to file with specified name
// in directory with specified path.
func ExportGUISave(gui *save.GUISave, dirPath, saveName string) error {
	gui.Name = saveName
	xml, err := parsexml.MarshalGUISave(gui)
	if err != nil {
		return fmt.Errorf("fail_to_marshal_gui_save:%v",
			err)
	}
	filePath := filepath.FromSlash(dirPath + "/" + saveName +
		SAVEGUI_FILE_EXT)
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("fail_to_create_save_file:%v",
			err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	w.WriteString(xml)
	w.Flush()
	log.Dbg.Printf("gui_state_saved_in:%s", filePath)
	return nil
}

// ImportGUISave imports GUI state from file with specified name
// in directory with specified path.
func ImportGUISave(game *flamecore.Game, dirPath,
	saveName string) (*save.GUISave, error) {
	if !strings.HasSuffix(saveName, SAVEGUI_FILE_EXT) {
		saveName = saveName + SAVEGUI_FILE_EXT
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
	save, err := buildXMLGUISave(game, xmlSave)
	if err != nil {
		return nil, fmt.Errorf("fail_to_build_save:%v",
			err)
	}
	return save, nil
}

// ImportsGUISavesDir imports all saved GUIs from save files in
// directory with specified path.
func ImportGUISavesDir(game *flamecore.Game, dirPath string) ([]*save.GUISave, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_read_dir:%v", err)
	}
	saves := make([]*save.GUISave, 0)
	for _, fInfo := range files {
		if !strings.HasSuffix(fInfo.Name(), SAVEGUI_FILE_EXT) {
			continue
		}
		sav, err := ImportGUISave(game, dirPath, fInfo.Name())
		if err != nil {
			log.Err.Printf("data_saves_import:fail_to_load_save_fail:%v",
				err)
			continue
		}
		saves = append(saves, sav)
	}
	return saves, nil
}

// buildGUISave builds GUI save from specified XML data.
func buildXMLGUISave(game *flamecore.Game, xmlSave *parsexml.GUISaveXML) (*save.GUISave, error) {
	save := new(save.GUISave)
	// Save name.
	save.Name = xmlSave.Name
	// Players.
	for _, xmlAvatar := range xmlSave.Players.Avatars {
		for _, pc := range game.Players() {
			if xmlAvatar.Serial == pc.Serial() {
				av, err := buildXMLAvatar(pc, &xmlAvatar)
				if err != nil {
					return nil, fmt.Errorf("player:%s:fail_to_load_player_avatar:%v",
						pc.SerialID, err)
				}
				save.Players = append(save.Players, av)
			}
		}
	}
	// Camera position.
	camX, camY, err := flameparsexml.UnmarshalPosition(xmlSave.Camera.Position)
	if err != nil {
		return nil, fmt.Errorf("fail_to_unmarshal_camera_position:%v",
			err)
	}
	save.CameraPosX, save.CameraPosY = camX, camY
	return save, nil
}
