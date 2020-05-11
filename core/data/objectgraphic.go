/*
 * objectgraphic.go
 *
 * Copyright 2019-2020 Dariusz Sikora <dev@isangeles.pl>
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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/log"
)

var (
	ObjectGraphicsFileExt = ".graphic"
)

// ImportObjectsGraphics imports all objects grpahics from
// base file with specified path.
func ImportObjectsGraphics(path string) ([]*res.ObjectGraphicData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open base file: %s", err)
	}
	objects, err := parsexml.UnmarshalObjectsGraphics(f)
	if err != nil {
		return nil, fmt.Errorf("unable to parse xml: %v", err)
	}
	return objects, nil
}

// ImportObjectsGraphicsDir imports all files with objects graphics from
// directory with specified path.
func ImportObjectsGraphicsDir(dirPath string) ([]*res.ObjectGraphicData, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir: %v", err)
	}
	objects := make([]*res.ObjectGraphicData, 0)
	for _, fInfo := range files {
		if !strings.HasSuffix(fInfo.Name(), ObjectGraphicsFileExt) {
			continue
		}
		basePath := filepath.FromSlash(dirPath + "/" + fInfo.Name())
		impObjects, err := ImportObjectsGraphics(basePath)
		if err != nil {
			log.Err.Printf("data items graphic import: %s: unable to parse file: %v",
				basePath, err)
			continue
		}
		for _, ob := range impObjects {
			objects = append(objects, ob)
		}
	}
	return objects, nil
}
