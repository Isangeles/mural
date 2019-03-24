/*
 * objectgraphicxml.go
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

// Struct for object graphic base XML node.
type ObjectGraphicBaseXML struct {
	XMLName xml.Name           `xml:"base"`
	Nodes   []ObjectGraphicXML `xml:"object"`
}

// Struct for object graphic XML node.
type ObjectGraphicXML struct {
	XMLName  xml.Name          `xml:"object"`
	ID       string            `xml:"id,attr"`
	Portrait ObjectPortraitXML `xml:"portrait"`
	Sprite   ObjectSpriteXML   `xml:"sprite"`
}

// Struct for object sprite XML node.
type ObjectSpriteXML struct {
	XMLName xml.Name `xml:"sprite"`
	Picture string   `xml:"picture,value"`
}

// Struct for object portrait XML node.
type ObjectPortraitXML struct {
	XMLName xml.Name `xml:"portrait"`
	Picture string   `xml:"picture,value"`
}

// UnmarshalkObjectsGraphicBase retrieves all object graphic data from
// specified XML data.
func UnmarshalObjectsGraphicBase(data io.Reader) ([]*res.ObjectGraphicData, error) {
	doc, _ := ioutil.ReadAll(data)
	xmlBase := new(ObjectGraphicBaseXML)
	err := xml.Unmarshal(doc, xmlBase)
	if err != nil {
		return nil, fmt.Errorf("fail_to_unmarshal_xml_data:%v", err)
	}
	objects := make([]*res.ObjectGraphicData, 0)
	for _, xmlData := range xmlBase.Nodes {
		ogd, err := buildObjectGraphicData(&xmlData)
		if err != nil {
			log.Err.Printf("xml:unmarshal_object_graphic:%s:fail_to_build_data:%v",
				xmlData.ID, err)
			continue
		}
		objects = append(objects, ogd)
	}
	return objects, nil
}

// buildObjectGraphicData creates object graphic data from specified XML object
// node.
func buildObjectGraphicData(xmlObject *ObjectGraphicXML) (*res.ObjectGraphicData, error) {
	sprite, err := data.ObjectSpritesheet(xmlObject.Sprite.Picture)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retireve_object_spritesheet:%v", err)
	}
	portrait, err := data.Portrait(xmlObject.Portrait.Picture)
	if err != nil {
		log.Err.Printf("xml:build_object_graphc_data:%s:fail_to_retrieve_object_portrait:%v",
			xmlObject.ID, err)
	}
	data := res.ObjectGraphicData{
		ID:          xmlObject.ID,
		PortraitPic: portrait,
		SpritePic:   sprite,
	}
	return &data, nil
}
