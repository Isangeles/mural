/*
 * objectgraphicxml.go
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

package parsexml

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/log"
)

// Struct for object graphic base XML node.
type ObjectGraphics struct {
	XMLName        xml.Name        `xml:"object-graphics"`
	ObjectGraphics []ObjectGraphic `xml:"object-graphic"`
}

// Struct for object graphic XML node.
type ObjectGraphic struct {
	XMLName  xml.Name       `xml:"object-graphic"`
	ID       string         `xml:"id,attr"`
	Portrait ObjectPortrait `xml:"portrait"`
	Sprite   ObjectSprite   `xml:"sprite"`
}

// Struct for object sprite XML node.
type ObjectSprite struct {
	XMLName xml.Name `xml:"sprite"`
	Picture string   `xml:"picture,value"`
}

// Struct for object portrait XML node.
type ObjectPortrait struct {
	XMLName xml.Name `xml:"portrait"`
	Picture string   `xml:"picture,value"`
}

// UnmarshalkObjectsGraphics retrieves all object graphic data from
// specified XML data.
func UnmarshalObjectsGraphics(data io.Reader) ([]*res.ObjectGraphicData, error) {
	doc, _ := ioutil.ReadAll(data)
	xmlBase := new(ObjectGraphics)
	err := xml.Unmarshal(doc, xmlBase)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal xml data: %v", err)
	}
	objects := make([]*res.ObjectGraphicData, 0)
	for _, xmlData := range xmlBase.ObjectGraphics {
		ogd, err := buildObjectGraphicData(&xmlData)
		if err != nil {
			log.Err.Printf("xml: unmarshal object graphic: %s: unable to build data: %v",
				xmlData.ID, err)
			continue
		}
		objects = append(objects, ogd)
	}
	return objects, nil
}

// buildObjectGraphicData creates object graphic data from specified XML object
// node.
func buildObjectGraphicData(xmlObject *ObjectGraphic) (*res.ObjectGraphicData, error) {
	data := res.ObjectGraphicData{
		ID:       xmlObject.ID,
		Portrait: xmlObject.Portrait.Picture,
		Sprite:   xmlObject.Sprite.Picture,
	}
	return &data, nil
}
