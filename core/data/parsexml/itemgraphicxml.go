/*
 * itemgraphicxml.go
 *
 * Copyright 2018-2019 Dariusz Sikora <dev@isangeles.pl>
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

package parsexml

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/log"
)

// Struct for XML items grpahics base.
type ItemsGraphicBaseXML struct {
	XMLName xml.Name         `xml:"base"`
	Nodes   []ItemGraphicXML `xml:"item"`
}

// Struct for XML item graphic node.
type ItemGraphicXML struct {
	XMLName     xml.Name `xml:"item"`
	ID          string   `xml:"id,attr"`
	Spritesheet string   `xml:"spritesheet,attr"`
	Icon        string   `xml:"icon,attr"`
	MaxStack    int      `xml:"stack,attr"`
}

// UnmarshalItemsGraphicsBase parses specified XML data
// to item grphics XML nodes.
func UnmarshalItemsGraphicBase(data io.Reader) ([]*res.ItemGraphicData, error) {
	doc, _ := ioutil.ReadAll(data)
	xmlBase := new(ItemsGraphicBaseXML)
	err := xml.Unmarshal(doc, xmlBase)
	if err != nil {
		return nil, fmt.Errorf("fail_to_unmarshal_xml_data:%v", err)
	}
	items := make([]*res.ItemGraphicData, 0)
	for _, xmlData := range xmlBase.Nodes {
		igd, err := buildItemGraphicData(&xmlData)
		if err != nil {
			log.Err.Printf("xml:unmarshal_item_graphic:%s:fail_to_build_data:%v",
				xmlData.ID, err)
			continue
		}
		items = append(items, igd)
	}
	return items, nil
}

// buildXMLItemGraphic creates item graphic object from
// specified item XML data.
func buildItemGraphicData(xmlItem *ItemGraphicXML) (*res.ItemGraphicData, error) {
	// Basic data.
	d := res.ItemGraphicData{
		ItemID:         xmlItem.ID,
		MaxStack:       xmlItem.MaxStack,
	}
	// Icon.
	icon, err := data.Icon(xmlItem.Icon)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_item_icon:%v", err)
	}
	d.IconPic = icon
	// Spritesheet.
	if len(xmlItem.Spritesheet) > 0 {
		sprite, err := data.ItemSpritesheet(xmlItem.Spritesheet)
		if err != nil {
			return nil, fmt.Errorf("fail_to_retrieve_item_spritesheet:%v", err)
		}
		d.SpritesheetPic = sprite
	}
	return &d, nil
}
