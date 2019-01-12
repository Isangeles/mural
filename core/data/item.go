/*
 * item.go
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
	"fmt"
	"os"

	flamedata "github.com/isangeles/flame/core/data"
	"github.com/isangeles/flame/core/module"
	flameitem "github.com/isangeles/flame/core/module/object/item"
	
	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// ImportItemsGraphics import all items grpahics from
// base file with specified path.
func ImportItemsGraphics(mod *module.Module, path string) ([]*object.Item, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_base_file:%s",
			err)
	}
	xmlItems, err := parsexml.UnmarshalItemsGraphicsBase(f)
	if err != nil {
		return nil, fmt.Errorf("fail_to_parse_xml:%v", err)
	}
	items := make([]*object.Item, 0)
	for _,  xmlItem := range xmlItems {
		baseItem, err := flamedata.Item(mod, xmlItem.ID)
		if err != nil {
			log.Err.Printf("data_item_graphic_import:%s:fail_to_retrieve_base_item:%v",
				xmlItem.ID, err)
			continue 
		}
		it, err := buildXMLItemGraphic(baseItem, &xmlItem)
		if err != nil {
			log.Err.Printf("data_item_graphic_import:%s:fail_build_item_graphic:%v",
				xmlItem.ID, err)
			continue 
		}
		items = append(items, it)
	}
	return items, nil
}

// buildXMLItemGraphic creates item graphic object for
// specified item.
func buildXMLItemGraphic(it flameitem.Item, xmlItem *parsexml.ItemGraphicNodeXML) (*object.Item,
	error) {
	itSprite, err := ItemSpritesheet(xmlItem.Spritesheet)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_item_spritesheet:%v",
			err)
	}
	// TODO: retrieve item icon picture.
	item := object.NewItem(it, itSprite, nil)
	return item, nil
}
