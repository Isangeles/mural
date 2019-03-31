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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/isangeles/flame/core/module/object/character"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/log"
)

var (
	AVATARS_FILE_EXT = ".avatars"
)

// ImportAvatarsData imports all avatars data for specified characters
// from avatar file with specified path.
func ImportAvatarsData(path string) ([]*res.AvatarData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("xml:%s:fail_to_open_avatars_file:%v",
			path, err)
	}
	avatars, err := parsexml.UnmarshalAvatarsBase(f)
	if err != nil {
		return nil, fmt.Errorf("xml:%s:fail_to_parse_XML:%v",
			path, err)
	}
	return avatars, nil
}

// ImportAvatarsDataDir imports all avatars data from avatars files
// in directory with specified path.
func ImportAvatarsDataDir(dirPath string) ([]*res.AvatarData, error) {
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
		impAvs, err := ImportAvatarsData(avFilePath)
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
		ID:           char.ID(),
		Serial:       char.Serial(),
		PortraitName: portraitName,
		SSHeadName:   ssHeadName,
		SSTorsoName:  ssTorsoName,
		PortraitPic:  portraitPic,
		SSHeadPic:    ssHeadPic,
		SSTorsoPic:   ssTorsoPic,
	}
	return &avData, nil
}

// defaultAvatarSpritesheet returns default spritesheet
// for specified character.
func defaultAvatarTorsoSpritesheet() string {
	return "m-cloth-1222211-80x90.png"
}

// defaultAvatarHeadSpritesheet retruns default spritesheet
// for specified character.
func defaultAvatarHeadSpritesheet() string {
	return "m-head-black-1222211-80x90.png"
}

// buildXMLAvatar builds avatar from specified XML data.
func buildXMLAvatarData(avXML *parsexml.AvatarXML) (*res.AvatarData, error) {
	ssHeadName := avXML.Spritesheet.Head
	ssTorsoName := avXML.Spritesheet.Torso
	portraitName := avXML.Portrait
	portraitPic, err := data.AvatarPortrait(portraitName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_portrait_picture:%v",
			avXML.ID, err)
	}
	if ssHeadName == "" {
		ssHeadName = defaultAvatarHeadSpritesheet()
	}
	if ssTorsoName == "" {
		ssTorsoName = defaultAvatarTorsoSpritesheet()
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
	avData := res.AvatarData{
		ID:           avXML.ID,
		Serial:       avXML.Serial,
		PortraitName: portraitName,
		SSHeadName:   ssHeadName,
		SSTorsoName:  ssTorsoName,
		PortraitPic:  portraitPic,
		SSHeadPic:    ssHeadPic,
		SSTorsoPic:   ssTorsoPic,
	}
	return &avData, nil
}
