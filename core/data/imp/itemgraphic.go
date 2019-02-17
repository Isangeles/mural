/*
 * itemgraphic.go
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

	"github.com/isangeles/flame/core/module"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/log"
)

var (
	ITEMS_GRAPHIC_FILE_EXT = ".graphic"
)

// ImportItemsGraphics import all items grpahics from
// base file with specified path.
func ImportItemsGraphics(mod *module.Module, path string) ([]*res.ItemGraphicData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_base_file:%s",
			err)
	}
	xmlItems, err := parsexml.UnmarshalItemsGraphicsBase(f)
	if err != nil {
		return nil, fmt.Errorf("fail_to_parse_xml:%v", err)
	}
	items := make([]*res.ItemGraphicData, 0)
	for _, xmlItem := range xmlItems {
		data, err := buildXMLItemGraphicData(&xmlItem)
		if err != nil {
			log.Err.Printf("data_imp_item_graphic:%s:fail_to_build_data:%v",
				xmlItem.ID, err)
			continue
		}
		items = append(items, data)
	}
	return items, nil
}

// ImportItemsGraphicsDir imports all files with items graphics from
// directory with specified path.
func ImportItemsGraphicsDir(mod *module.Module, dirPath string) ([]*res.ItemGraphicData, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_read_dir:%v", err)
	}
	items := make([]*res.ItemGraphicData, 0)
	for _, fInfo := range files {
		if !strings.HasSuffix(fInfo.Name(), ITEMS_GRAPHIC_FILE_EXT) {
			continue
		}
		itemsGraphicFilePath := filepath.FromSlash(dirPath + "/" + fInfo.Name())
		impItems, err := ImportItemsGraphics(mod, itemsGraphicFilePath)
		if err != nil {
			log.Err.Printf("data_items_graphic_import:%s:fail_to_parse_file:%v",
				itemsGraphicFilePath, err)
			continue
		}
		for _, it := range impItems {
			items = append(items, it)
		}
	}
	return items, nil
}

// buildXMLItemGraphic creates item graphic object for
// specified item.
func buildXMLItemGraphicData(xmlItem *parsexml.ItemGraphicNodeXML) (*res.ItemGraphicData, error) {
	itSprite, err := data.ItemSpritesheet(xmlItem.Spritesheet)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_item_spritesheet:%v",
			err)
	}
	itIcon, err := data.ItemIcon(xmlItem.Icon)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_item_icon:%v", err)
	}
	itData := res.ItemGraphicData {
		ItemID: xmlItem.ID,
		IconPic: itIcon,
		SpritesheetPic: itSprite,
		MaxStack: xmlItem.MaxStack,
	}
	return &itData, nil
}
