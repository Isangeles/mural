/*
 * avatar.go
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
	"io/ioutil"

	"github.com/isangeles/flame/core/module/object/character"

	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

var (
	AVATARS_FILE_EXT = ".avatars"
)

// CharacterAvatar imports and returns avatars for specified
// character from avatars file in directory with specified path.
func CharacterAvatar(importDir string, char *character.Character) (*object.Avatar,
	error) {
	// Search all avatars files in directory for character avatar.
	avs, err := ImportAvatarsDir([]*character.Character{char}, importDir)
	if err != nil {
		return nil, fmt.Errorf("fail_to_import_avatar:%v",
			err)
	}
	if len(avs) < 1 {
		return nil, fmt.Errorf("avatar_not_found_for:%s",
			char.ID())
	}
	return avs[0], nil
}

// ExportAvatars exports specified avatars to file
// with specified path.
func ExportAvatars(avs []*object.Avatar, basePath string) error {
	// Marshal avatars to base data.
	xml, err := parsexml.MarshalAvatarsBase(avs)
	if err != nil {
		return fmt.Errorf("fail_to_marshal_avatars:%v", err)
	}
	// Check whether file path ends with proper extension.
	if !strings.HasSuffix(basePath, AVATARS_FILE_EXT) {
		basePath = basePath + AVATARS_FILE_EXT
	}
	// Create base file.
	f, err := os.Create(filepath.FromSlash(basePath))
	if err != nil {
		return fmt.Errorf("fail_to_create_avatars_file:%v", err)
	}
	defer f.Close()
	// Write data to base file.
	w := bufio.NewWriter(f)
	w.WriteString(xml)
	w.Flush()
	return nil
}

// ExportAvatars exports specified avatar to directory
// with specified path.
func ExportAvatar(av *object.Avatar, dirPath string) error {
	filePath := filepath.FromSlash(dirPath + "/" + strings.ToLower(av.Name()) +
		AVATARS_FILE_EXT)
	return ExportAvatars([]*object.Avatar{av}, filePath)
}

// ImportAvatars imports all avatars for specified characters
// from avatar file with specified path.
func ImportAvatars(chars []*character.Character, path string) ([]*object.Avatar,
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
	avs := make([]*object.Avatar, 0)
	for _, avXML := range avsXML {
		for _, c := range chars {
			if avXML.ID != c.ID() {
				continue
			}
			var av *object.Avatar
			if avXML.Spritesheet.FullBody != "" {
				av, err = buildXMLStaticAvatar(c, &avXML)
			} else {
				av, err = buildXMLAvatar(c, &avXML)
			}
			if err != nil {
				log.Err.Printf("data_avatar_import:%s:parse_fail:%v",
					avXML.ID, err)
				continue
			}
			avs = append(avs, av)
		}
	}
	return avs, nil
}

// ImportAvatarsDir imports all avatars from avatars files
// in directory with specified path.
func ImportAvatarsDir(chars []*character.Character,
	dirPath string) ([]*object.Avatar, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_read_dir:%v", err)
	}
	avs := make([]*object.Avatar, 0)
	for _, fInfo := range files {
		if !strings.HasSuffix(fInfo.Name(), AVATARS_FILE_EXT) {
			continue
		}
		avFilePath := filepath.FromSlash(dirPath + "/" + fInfo.Name())
		impAvs, err := ImportAvatars(chars, avFilePath)
		if err != nil {
			log.Err.Printf("data_avatar_import:%s:fail_to_parse_char_file:%v",
				dirPath, err)
			continue
		}
		for _, av := range impAvs {
			avs = append(avs, av)
		}
	}
	return avs, nil
}

// DefaultAvatar creates default avatar for specified
// character.
func DefaultAvatar(char *character.Character) (*object.Avatar, error) {
	ssHeadName := "m-head-black-1222211-80x90.png"
	ssTorsoName := "m-cloth-1222211-80x90.png"
	portraitName := "male01.png"
	if char.Gender() == character.Female {
		ssHeadName = "f-head-black-1222211-80x90.png"
		ssTorsoName = "f-cloth-1222211-80x90.png"
		portraitName = "female01.png"
	}
	ssHeadPic, err := AvatarSpritesheet(ssHeadName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_head_spritesheet_picture:%v",
			err)
	}
	ssTorsoPic, err := AvatarSpritesheet(ssTorsoName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_torso_spritesheet_picture:%v",
			err)
	}
	portraitPic, err := AvatarPortrait(portraitName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_portrait_picture:%v\n",
			err)
	}
	av, err := object.NewAvatar(char, portraitPic, ssHeadPic, ssTorsoPic,
		portraitName, ssHeadName, ssTorsoName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_avatar:%v", err)
	}
	return av, nil	
}

// buildXMLAvatar builds avatar from specified XML data.
func buildXMLAvatar(char *character.Character, avXML *parsexml.AvatarXML) (*object.Avatar, error) {
	ssHeadName := avXML.Spritesheet.Head
	ssTorsoName := avXML.Spritesheet.Torso
	portraitName := avXML.Portrait
	portraitPic, err := AvatarPortrait(portraitName)
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
	ssHeadPic, err := AvatarSpritesheet(ssHeadName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_head_spritesheet_picture:%v", 
			avXML.ID, err)
	}
	ssTorsoPic, err := AvatarSpritesheet(ssTorsoName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_torso_spritesheet_picture:%v", 
			avXML.ID, err)
	}
	av, err := object.NewAvatar(char, portraitPic, ssHeadPic, ssTorsoPic,
		portraitName, ssHeadName, ssTorsoName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_avatar:%v", err)
	}
	return av, nil
}

// buildXMLStaticAvatar build new static avatar for specified
// character from specified XML data.
func buildXMLStaticAvatar(char *character.Character,
	avXML *parsexml.AvatarXML) (*object.Avatar, error) {
	ssFullBodyName := avXML.Spritesheet.FullBody
	portraitName := avXML.Portrait
	portraitPic, err := AvatarPortrait(portraitName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_portrait_picture:%v",
			err)
	}
	ssFullBodyPic, err := AvatarSpritesheet(ssFullBodyName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_head_spritesheet_picture:%v",
			err)
	}
	av, err := object.NewStaticAvatar(char, portraitPic, ssFullBodyPic,
		portraitName, ssFullBodyName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_create_avatar:%v", err)
	}
	return av, nil
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
