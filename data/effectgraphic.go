/*
 * effectgraphic.go
 *
 * Copyright 2019-2021 Dariusz Sikora <dev@isangeles.pl>
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
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/log"
)

var (
	EffectGraphicsFileExt = ".graphic"
)

// ImportEffectsGraphics imports all effects graphics from
// data file with specified path.
func ImportEffectsGraphics(path string) ([]res.EffectGraphicData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open data file: %v", err)
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read data file: %v", err)
	}
	data := new(res.EffectGraphicsData)
	err = xml.Unmarshal(buf, data)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal XML data: %v", err)
	}
	return data.Effects, nil
}

// ImportEffectsGraphicsDir imports all files with effects graphics from
// directory with specified path.
func ImportEffectsGraphicsDir(path string) ([]res.EffectGraphicData, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir: %v", err)
	}
	effects := make([]res.EffectGraphicData, 0)
	for _, finfo := range files {
		if !strings.HasSuffix(finfo.Name(), EffectGraphicsFileExt) {
			continue
		}
		basePath := filepath.FromSlash(path + "/" + finfo.Name())
		impEffects, err := ImportEffectsGraphics(basePath)
		if err != nil {
			log.Err.Printf("data effects graphic import: %s: unable to parse file: %v",
				basePath, err)
			continue
		}
		effects = append(effects, impEffects...)
	}
	return effects, nil
}
