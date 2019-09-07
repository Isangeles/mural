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

package exp

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/core/object"
)

var (
	AvatarsFileExt = ".avatars"
)

// ExportAvatars exports specified avatars to file
// with specified path.
func ExportAvatars(avs []*object.Avatar, basePath string) error {
	// Marshal avatars to base data.
	xml, err := parsexml.MarshalAvatars(avs)
	if err != nil {
		return fmt.Errorf("fail to marshal avatars: %v", err)
	}
	// Check whether file path ends with proper extension.
	if !strings.HasSuffix(basePath, AvatarsFileExt) {
		basePath = basePath + AvatarsFileExt
	}
	// Create base file.
	f, err := os.Create(filepath.FromSlash(basePath))
	if err != nil {
		return fmt.Errorf("fail to create avatars file: %v", err)
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
		AvatarsFileExt)
	return ExportAvatars([]*object.Avatar{av}, filePath)
}
