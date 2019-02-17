/*
 * save.go
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

package imp

import (
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"strings"

	flamecore "github.com/isangeles/flame/core"
	flamexml "github.com/isangeles/flame/core/data/parsexml"
	
	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/log"
)

var (
	SAVEGUI_FILE_EXT = ".savegui"
)

// ImportGUISave imports GUI state from file with specified name
// in directory with specified path.
func ImportGUISave(game *flamecore.Game, dirPath, saveName string) (*res.GUISave,
	error) {
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
func ImportGUISavesDir(game *flamecore.Game, dirPath string) ([]*res.GUISave, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_read_dir:%v", err)
	}
	saves := make([]*res.GUISave, 0)
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
func buildXMLGUISave(game *flamecore.Game, xmlSave *parsexml.GUISaveXML) (*res.GUISave, error) {
	save := new(res.GUISave)
	// Save name.
	save.Name = xmlSave.Name
	// Players.
	for _, xmlPC := range xmlSave.PlayersNode.Players {
		for _, pc := range game.Players() {
			if xmlPC.Avatar.Serial == pc.Serial() {
				pcData := new(res.PlayerSave)
				// Character.
				pcData.Character = pc
				// Avatar.
				avData, err := buildXMLAvatarData(&xmlPC.Avatar)
				if err != nil {
					return nil, fmt.Errorf("player:%s:fail_to_load_player_avatar:%v",
						pc.SerialID, err)
				}
				pcData.Avatar = avData
				// Inventory layout.
				pcData.InvSlots = make(map[string]int)
				for _, xmlSlot := range xmlPC.Inventory.Slots {
					pcData.InvSlots[xmlSlot.Content] = xmlSlot.ID
				}
				save.PlayersData = append(save.PlayersData, pcData)
			}
		}
	}
	// Camera position.
	camX, camY, err := flamexml.UnmarshalPosition(xmlSave.CameraNode.Position)
	if err != nil {
		return nil, fmt.Errorf("fail_to_unmarshal_camera_position:%v",
			err)
	}
	save.CameraPosX, save.CameraPosY = camX, camY
	return save, nil
}
