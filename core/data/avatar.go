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
	"github.com/isangeles/mural/objects"
	"github.com/isangeles/mural/log"
)

var (
	AVATARS_FILE_EXT = ".avatars"
)

// CharacterAvatar imports and returns avatars for specified
// character from avatars file in directory with specified path.
func CharacterAvatar(importDir string, char *character.Character) (*objects.Avatar,
	error) {
	// Search all avatars files in directory for character avatar.
	avs, err := ImportAvatarsDir([]*character.Character{char}, importDir)
	if err != nil {
		return nil, fmt.Errorf("fail_to_import_avatar:%v",
			err)
	}
	if len(avs) < 1 {
		// If avatar for character was not found,
		// then build and return default avatar.
		log.Dbg.Printf("avatar_not_found_for:%s",
			char.ID())
		return DefaultAvatar(char)
	}
	return avs[0], nil
}

// ExportAvatars exports specified avatars to file
// with specified path.
func ExportAvatars(avs []*objects.Avatar, basePath string) error {
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
func ExportAvatar(av *objects.Avatar, dirPath string) error {
	filePath := filepath.FromSlash(dirPath + "/" + strings.ToLower(av.Name()) +
		AVATARS_FILE_EXT)
	return ExportAvatars([]*objects.Avatar{av}, filePath)
}

// ImportAvatars imports all avatars for specified characters
// from avatar file with specified path.
func ImportAvatars(chars []*character.Character, path string) ([]*objects.Avatar,
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
	avs := make([]*objects.Avatar, 0)
	for _, avXML := range avsXML {
		for _, c := range chars {
			if avXML.ID != c.ID() {
				continue
			}
			av, err := buildXMLAvatar(c, &avXML)
			if err != nil {
				log.Err.Printf("data_avatars_import:parse_fail:%s:%v",
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
	dirPath string) ([]*objects.Avatar, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_read_dir:%v", err)
	}
	avs := make([]*objects.Avatar, 0)
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
func DefaultAvatar(char *character.Character) (*objects.Avatar, error) {
	spritesheetName := "test.png"
	spritesheetPic, err := AvatarSpritesheet(spritesheetName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_spritesheet_picture:%v",
			err)
	}
	portraitName := "male01.png"
	portraitPic, err := AvatarPortrait(portraitName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_portrait_picture:%v\n",
			err)
	}
	av := objects.NewAvatar(char, portraitPic, spritesheetPic, portraitName,
		spritesheetName)
	return av, nil	
}

// buildXMLAvatar builds avatar from specified XML data.
func buildXMLAvatar(char *character.Character, avXML *parsexml.AvatarXML) (*objects.Avatar, error) {
	portraitPic, err := AvatarPortrait(avXML.Portrait)
	if err != nil {
		return nil, fmt.Errorf("data:parse_fail:%s:fail_to_retrieve_portrait_picture:%v",
			avXML.ID, err)
	}
	spritesheetPic, err := AvatarSpritesheet(avXML.Spritesheet)
	if err != nil {
		return nil, fmt.Errorf("data:parse_fail:%s:fail_to_retrieve_spritesheet_picture:%v", 
			avXML.ID, err)
	}
	av := objects.NewAvatar(char, portraitPic, spritesheetPic,
		avXML.Portrait, avXML.Spritesheet)
	return av, nil
}