/*
 * avatar.go
 *
 * Copyright 2019-2022 Dariusz Sikora <ds@isangeles.dev>
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
	"io"
	"os"
	"path/filepath"

	"github.com/isangeles/flame/character"

	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/log"
)

// ImportAvatars imports all avatars data for specified characters
// from avatar file with specified path.
func ImportAvatars(path string) ([]res.AvatarData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open data file: %v", err)
	}
	defer file.Close()
	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read data file: %v", err)
	}
	data := new(res.AvatarsData)
	err = xml.Unmarshal(buf, data)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal XML data: %v",
			err)
	}
	return data.Avatars, nil
}

// ImportAvatarsDir imports all avatars data from avatars files
// in directory with specified path.
func ImportAvatarsDir(dirPath string) ([]res.AvatarData, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir: %v", err)
	}
	avsData := make([]res.AvatarData, 0)
	for _, fInfo := range files {
		avFilePath := filepath.FromSlash(dirPath + "/" + fInfo.Name())
		impAvs, err := ImportAvatars(avFilePath)
		if err != nil {
			log.Err.Printf("data avatar import: %s: unable to parse char file: %v",
				avFilePath, err)
			continue
		}
		for _, av := range impAvs {
			avsData = append(avsData, av)
		}
	}
	return avsData, nil
}

// ExportAvatars exports specified avatars to the new avatar data file
// with specified path.
func ExportAvatars(path string, avatars ...res.AvatarData) error {
	// Marshal avatars.
	data := new(res.AvatarsData)
	for _, av := range avatars {
		data.Avatars = append(data.Avatars, av)
	}
	xml, err := xml.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal XML data: %v", err)
	}
	// Create base file.
	dirPath := filepath.Dir(path)
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("unable to create file directory: %v", err)
	}
	f, err := os.Create(filepath.FromSlash(path))
	if err != nil {
		return fmt.Errorf("unable to create avatars file: %v", err)
	}
	defer f.Close()
	// Write data to base file.
	w := bufio.NewWriter(f)
	w.Write(xml)
	w.Flush()
	return nil
}

// DefaultAvatarData creates default avatar data
// for specified character.
func DefaultAvatarData(char *character.Character) res.AvatarData {
	return res.AvatarData{
		ID:       char.ID(),
		Serial:   char.Serial(),
		FullBody: "default.png",
	}
}
