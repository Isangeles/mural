/*
 * effectgraphic.go
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
	"io/ioutil"
	"strings"
	"path/filepath"

	"github.com/isangeles/mural/core/data"	
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/log"
)

var (
	EFFECTS_GRAPHIC_FILE_EXT = ".graphic"
)

// ImportEffectsGraphics imports all effects graphics from
// base file with specified path.
func ImportEffectsGraphics(path string) ([]*res.EffectGraphicData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_base_file:%s", err)
	}
	xmlEffects, err := parsexml.UnmarshalEffectsGraphicsBase(f)
	if err != nil {
		return nil, fmt.Errorf("fail_to_parse_xml:%v", err)
	}
	effects := make([]*res.EffectGraphicData, 0)
	for _, xmlEffect := range xmlEffects {
		data, err := buildXMLEffectGraphicData(&xmlEffect)
		if err != nil {
			log.Err.Printf("data_imp_effect_graphic:%s:fail_to_build_data:%v",
				xmlEffect.ID, err)
		}
		effects = append(effects, data)
	}
	return effects, nil
}

// ImportEffectsGraphicsDir imports all files with effects graphics from
// directory with specified path.
func ImportEffectsGraphicsDir(dirPath string) ([]*res.EffectGraphicData, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_read_dir:%v", err)
	}
	effects := make([]*res.EffectGraphicData, 0)
	for _, fInfo := range files {
		if !strings.HasSuffix(fInfo.Name(), EFFECTS_GRAPHIC_FILE_EXT) {
			continue
		}
		effectsGraphicFilePath := filepath.FromSlash(dirPath + "/" + fInfo.Name())
		impEffects, err := ImportEffectsGraphics(effectsGraphicFilePath)
		if err != nil {
			log.Err.Printf("data_effects_graphic_import:%s:fail_to_parse_file:%v",
				effectsGraphicFilePath, err)
			continue
		}
		for _, eff := range impEffects {
			effects = append(effects, eff)
		}
	}
	return effects, nil
}

// buildXMLEffectGraphicData creates effect graphic object from specified
// effect XML data.
func buildXMLEffectGraphicData(xmlEffect *parsexml.EffectGraphicNodeXML) (*res.EffectGraphicData,
	error) {
	effIcon, err := data.Icon(xmlEffect.Icon)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_effect_icon:%v", err)
	}
	effData := res.EffectGraphicData {
		EffectID: xmlEffect.ID,
		IconPic: effIcon,
	}
	return &effData, nil
}