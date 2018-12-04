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

	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/objects"
)

var (
	AVATAR_FILE_PREFIX = ".avatar"
)

// ExportAvatars exports specified avatar to '/characters'
// module directory.
func ExportAvatar(av *objects.Avatar, dirPath string) error {
	xml, err := parsexml.MarshalAvatar(av)
	if err != nil {
		return fmt.Errorf("fail_to_marshal_avatar:%v", err)
	}

	f, err := os.Create(filepath.FromSlash(dirPath + "/" + av.Name() +
		AVATAR_FILE_PREFIX))
	if err != nil {
		return fmt.Errorf("fail_to_create_avatar_file:%v", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(xml)
	w.Flush()
	return nil
}

// ImportAvatars imports all avatars from avatars files in
// directory with specified path.
func ImportAvatars(dirPath string) ([]*objects.Avatar, error) {
	avs := make([]*objects.Avatar, 0)
	// TODO: unmarshal XML avatars base.
	return avs, fmt.Errorf("unsupported_yet")
}
