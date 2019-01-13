/*
 * avatar.go
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
	"path/filepath"
	"strings"
	"io/ioutil"

	"github.com/isangeles/flame/core/module/object/character"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/log"
)

var (
	AVATARS_FILE_EXT = ".avatars"
)

// CharacterAvatar imports and returns avatars for specified
// character from avatars file in directory with specified path.
func CharacterAvatarData(char *character.Character, importDir string) (*res.AvatarData,
	error) {
	// Search all avatars files in directory for character avatar.
	avsData, err := ImportAvatarsDataDir([]*character.Character{char}, importDir)
	if err != nil {
		return nil, fmt.Errorf("fail_to_import_avatar:%v",
			err)
	}
	if len(avsData) < 1 {
		return nil, fmt.Errorf("avatar_not_found_for:%s",
			char.ID())
	}
	return avsData[0], nil
}

// ImportAvatarsData imports all avatars data for specified characters
// from avatar file with specified path.
func ImportAvatarsData(chars []*character.Character, path string) ([]*res.AvatarData,
	error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("xml:%s:fail_to_open_avatars_file:%v",
			path, err)
	}
	avsXML, err := parsexml.UnmarshalAvatarsBase(f)
	if err != nil {
		return nil, fmt.Errorf("xml:%s:fail_to_parse_XML:%v",
			path, err)
	}
	avsData := make([]*res.AvatarData, 0)
	for _, avXML := range avsXML {
		for _, c := range chars {
			if avXML.ID != c.ID() {
				continue
			}
			var avData *res.AvatarData
			if avXML.Spritesheet.FullBody != "" {
				avData, err = buildXMLStaticAvatarData(c, &avXML)
			} else {
				avData, err = buildXMLAvatarData(c, &avXML)
			}
			if err != nil {
				log.Err.Printf("data_avatar_import:%s:parse_fail:%v",
					avXML.ID, err)
				continue
			}
			avsData = append(avsData, avData)
		}
	}
	return avsData, nil
}

// ImportAvatarsDataDir imports all avatars data from avatars files
// in directory with specified path.
func ImportAvatarsDataDir(chars []*character.Character,
	dirPath string) ([]*res.AvatarData, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_read_dir:%v", err)
	}
	avsData := make([]*res.AvatarData, 0)
	for _, fInfo := range files {
		if !strings.HasSuffix(fInfo.Name(), AVATARS_FILE_EXT) {
			continue
		}
		avFilePath := filepath.FromSlash(dirPath + "/" + fInfo.Name())
		impAvs, err := ImportAvatarsData(chars, avFilePath)
		if err != nil {
			log.Err.Printf("data_avatar_import:%s:fail_to_parse_char_file:%v",
				avFilePath, err)
			continue
		}
		for _, av := range impAvs {
			avsData = append(avsData, av)
		}
	}
	return avsData, nil
}

// DefaultAvatarData creates default avatar data
// for specified character.
func DefaultAvatarData(char *character.Character) (*res.AvatarData, error) {
	ssHeadName := "m-head-black-1222211-80x90.png"
	ssTorsoName := "m-cloth-1222211-80x90.png"
	portraitName := "male01.png"
	if char.Gender() == character.Female {
		ssHeadName = "f-head-black-1222211-80x90.png"
		ssTorsoName = "f-cloth-1222211-80x90.png"
		portraitName = "female01.png"
	}
	ssHeadPic, err := data.AvatarSpritesheet(ssHeadName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_head_spritesheet_picture:%v",
			err)
	}
	ssTorsoPic, err := data.AvatarSpritesheet(ssTorsoName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_torso_spritesheet_picture:%v",
			err)
	}
	portraitPic, err := data.AvatarPortrait(portraitName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_portrait_picture:%v\n",
			err)
	}
	avData := res.AvatarData{
		Character: char,
		PortraitName: portraitName,
		SSHeadName: ssHeadName,
		SSTorsoName: ssTorsoName,
		PortraitPic: portraitPic,
		SSHeadPic: ssHeadPic,
		SSTorsoPic: ssTorsoPic,
	}
	return &avData, nil	
}

// defaultAvatarSpritesheet returns default spritesheet
// for specified character.
func defaultAvatarTorsoSpritesheet(char *character.Character) string {
	switch char.Race() {
	default:
		if char.Gender() == character.Female {
			return  "f-cloth-1222211-80x90.png"
		}
		return  "m-cloth-1222211-80x90.png"
	}
}

// defaultAvatarHeadSpritesheet retruns default spritesheet
// for specified character.
func defaultAvatarHeadSpritesheet(char *character.Character) string {
	switch char.Race() {
	default:
		if char.Gender() == character.Female {
			return  "f-head-black-1222211-80x90.png"
		}
		return  "m-head-black-1222211-80x90.png"
	}
}

// buildXMLAvatar builds avatar from specified XML data.
func buildXMLAvatarData(char *character.Character, avXML *parsexml.AvatarXML) (*res.AvatarData, error) {
	ssHeadName := avXML.Spritesheet.Head
	ssTorsoName := avXML.Spritesheet.Torso
	portraitName := avXML.Portrait
	portraitPic, err := data.AvatarPortrait(portraitName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_portrait_picture:%v",
			avXML.ID, err)
	}
	if ssHeadName == "" {
		ssHeadName = defaultAvatarHeadSpritesheet(char)
	}
	if ssTorsoName == "" {
		ssTorsoName = defaultAvatarTorsoSpritesheet(char)
	}
	ssHeadPic, err := data.AvatarSpritesheet(ssHeadName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_head_spritesheet_picture:%v", 
			avXML.ID, err)
	}
	ssTorsoPic, err := data.AvatarSpritesheet(ssTorsoName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_torso_spritesheet_picture:%v", 
			avXML.ID, err)
	}
	eqItemsGraphics := make([]*res.ItemGraphicData, 0)
	for _, eqi := range char.Equipment().Items() {
		itemGData := res.ItemData(eqi.ID())
		if itemGData == nil {
			log.Err.Printf("data_buil_avatar_data:%s:item_graphic_not_found:%s",
				avXML.ID, eqi.ID())
			continue
		}
		eqItemsGraphics = append(eqItemsGraphics, itemGData)
	}
	avData := res.AvatarData{
		Character: char,
		PortraitName: portraitName,
		SSHeadName: ssHeadName,
		SSTorsoName: ssTorsoName,
		PortraitPic: portraitPic,
		SSHeadPic: ssHeadPic,
		SSTorsoPic: ssTorsoPic,
		EqItemsGraphics: eqItemsGraphics,
	}
	return &avData, nil
}

// buildXMLStaticAvatar build new static avatar for specified
// character from specified XML data.
func buildXMLStaticAvatarData(char *character.Character, avXML *parsexml.AvatarXML) (*res.AvatarData, error) {
	ssFullBodyName := avXML.Spritesheet.FullBody
	portraitName := avXML.Portrait
	portraitPic, err := data.AvatarPortrait(portraitName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_portrait_picture:%v",
			err)
	}
	ssFullBodyPic, err := data.AvatarSpritesheet(ssFullBodyName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_head_spritesheet_picture:%v",
			err)
	}
	eqItemsGraphics := make([]*res.ItemGraphicData, 0)
	for _, eqi := range char.Equipment().Items() {
		itemGData := res.ItemData(eqi.ID())
		if itemGData == nil {
			log.Err.Printf("data_buil_avatar_data:%s:item_graphic_not_found:%s",
				avXML.ID, eqi.ID())
			continue
		}
		eqItemsGraphics = append(eqItemsGraphics, itemGData)
	}
	avData := res.AvatarData{
		Character: char,
		PortraitName: portraitName,
		SSFullBodyName: ssFullBodyName,
		PortraitPic: portraitPic,
		SSFullBodyPic: ssFullBodyPic,
		EqItemsGraphics: eqItemsGraphics,
	}
	return &avData, nil
}
