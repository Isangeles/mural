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
type ItemGraphics struct {
	XMLName      xml.Name      `xml:"item-graphics"`
	ItemGraphics []ItemGraphic `xml:"item-graphic"`
}

// Struct for XML item graphic node.
type ItemGraphic struct {
	XMLName     xml.Name `xml:"item-graphic"`
	ID          string   `xml:"id,attr"`
	Spritesheet string   `xml:"spritesheet,attr"`
	Icon        string   `xml:"icon,attr"`
	Stack       int      `xml:"stack,attr"`
}

// UnmarshalItemGraphics parses specified XML data
// to item grphics XML nodes.
func UnmarshalItemGraphics(data io.Reader) ([]*res.ItemGraphicData, error) {
	doc, _ := ioutil.ReadAll(data)
	xmlBase := new(ItemGraphics)
	err := xml.Unmarshal(doc, xmlBase)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal xml data: %v", err)
	}
	items := make([]*res.ItemGraphicData, 0)
	for _, xmlData := range xmlBase.ItemGraphics {
		igd, err := buildItemGraphicData(&xmlData)
		if err != nil {
			log.Err.Printf("xml: unmarshal item graphic: %s: fail to build data: %v",
				xmlData.ID, err)
			continue
		}
		items = append(items, igd)
	}
	return items, nil
}

// buildXMLItemGraphic creates item graphic object from
// specified item XML data.
func buildItemGraphicData(xmlItem *ItemGraphic) (*res.ItemGraphicData, error) {
	// Basic data.
	d := res.ItemGraphicData{
		ItemID:   xmlItem.ID,
		MaxStack: xmlItem.Stack,
	}
	// Icon.
	icon, err := data.Icon(xmlItem.Icon)
	if err != nil {
		return nil, fmt.Errorf("fail to retrieve item icon: %v", err)
	}
	d.IconPic = icon
	// Spritesheet.
	if len(xmlItem.Spritesheet) > 0 {
		sprite, err := data.ItemSpritesheet(xmlItem.Spritesheet)
		if err != nil {
			return nil, fmt.Errorf("fail to retrieve item spritesheet: %v", err)
		}
		d.SpritesheetPic = sprite
	}
	return &d, nil
}
