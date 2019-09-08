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

	flamexml "github.com/isangeles/flame/core/data/parsexml"

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
	XMLName      xml.Name      `xml:"item-graphic"`
	ID           string        `xml:"id,attr"`
	Spritesheet  string        `xml:"spritesheet,attr"`
	Icon         string        `xml:"icon,attr"`
	Stack        int           `xml:"stack,attr"`
	Spritesheets []Spritesheet `xml:"spritesheets>spritesheet"`
}

// Struct for XML graphic spritesheet node.
type Spritesheet struct {
	XMLName xml.Name `xml:"spritesheet"`
	Texture string   `xml:"texture,attr"`
	Race    string   `xml:"race,attr"`
	Gender  string   `xml:"gender,attr"`
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

// buildXMLItemGraphic creates item graphic data from specified item graphic node.
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
	// Spritesheets.
	d.Spritesheets = make([]*res.SpritesheetData, 0)
	for _, xmlSpritesheet := range xmlItem.Spritesheets {
		s, err := buildSpritesheetData(&xmlSpritesheet)
		if err != nil {
			return nil, fmt.Errorf("fail to build spritesheet data: %v", err)
		}
		d.Spritesheets = append(d.Spritesheets, s)
	}
	return &d, nil
}

// buildSpritesheetData creates spriteseheet data from specified XML
// spritesheet node.
func buildSpritesheetData(xmlSpritesheet *Spritesheet) (*res.SpritesheetData, error) {
	tex, err := data.ItemSpritesheet(xmlSpritesheet.Texture)
	if err != nil {
		return nil, fmt.Errorf("fail to retrieve texture: %v", err)
	}
	race, err := flamexml.UnmarshalRace(xmlSpritesheet.Race)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal race: %v", err)
	}
	gender, err := flamexml.UnmarshalGender(xmlSpritesheet.Gender)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal gender: %v", err)
	}
	d := res.SpritesheetData{
		Texture: tex,
		Race:    int(race),
		Gender:  int(gender),
	}
	return &d, nil
}
