/*
 * parse.go
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
	AVATAR_FILE_EXT = ".avatars"
)

// ExportAvatars exports specified avatars to file
// with specified path.
func ExportAvatars(avs []*objects.Avatar, basePath string) error {
	xml, err := parsexml.MarshalAvatarsBase(avs)
	if err != nil {
		return fmt.Errorf("fail_to_marshal_avatars:%v", err)
	}

	if !strings.HasSuffix(basePath, AVATAR_FILE_EXT) {
		basePath = basePath + AVATAR_FILE_EXT
	}
	f, err := os.Create(filepath.FromSlash(basePath))
	if err != nil {
		return fmt.Errorf("fail_to_create_avatars_file:%v", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(xml)
	w.Flush()
	return nil
}

// ExportAvatars exports specified avatar to '/characters'
// module directory.
func ExportAvatar(av *objects.Avatar, dirPath string) error {
	filePath := filepath.FromSlash(dirPath + "/" + av.Name() +
		AVATAR_FILE_EXT)
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
			portraitPic, err := AvatarPortrait(avXML.Portrait)
			if err != nil {
				log.Err.Printf("data:parse_fail:%s:fail_to_retrieve_portrait_picture:%v",
					avXML.ID, err)
				continue
			}
			spritesheetPic, err := AvatarSpritesheet(avXML.Spritesheet)
			if err != nil {
				log.Err.Printf("data:parse_fail:%s:fail_to_retrieve_spritesheet_picture:%v",
					avXML.ID, err)
				continue
			}
			av := objects.NewAvatar(c, portraitPic, spritesheetPic,
				avXML.Portrait, avXML.Spritesheet)
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
		return nil, fmt.Errorf("xml_dir:%s:fail_to_read_dir:%v",
			dirPath, err)
	}
	avs := make([]*objects.Avatar, 0)
	for _, fInfo := range files {
		if !strings.HasSuffix(fInfo.Name(), AVATAR_FILE_EXT) {
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
