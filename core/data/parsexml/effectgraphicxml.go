/*
 * effectgraphicxml.go
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

// Struct for XML effect graphic base.
type EffectGraphics struct {
	XMLName  xml.Name        `xml:"effect-graphics"`
	Graphics []EffectGraphic `xml:"effect-graphic"`
}

// Struct for XML effect graphic node.
type EffectGraphic struct {
	XMLName xml.Name `xml:"effect-graphic"`
	ID      string   `xml:"id,attr"`
	Icon    string   `xml:"icon,attr"`
}

// UnmarshalEffectGraphics retrieves all effect graphic
// data from specified XML data.
func UnmarshalEffectGraphics(data io.Reader) ([]*res.EffectGraphicData, error) {
	doc, _ := ioutil.ReadAll(data)
	xmlBase := new(EffectGraphics)
	err := xml.Unmarshal(doc, xmlBase)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal xml data: %v", err)
	}
	effects := make([]*res.EffectGraphicData, 0)
	for _, xmlData := range xmlBase.Graphics {
		egd, err := buildEffectGraphicData(&xmlData)
		if err != nil {
			log.Err.Printf("xml: unmarshal effect graphic: %s: fail to build data: %v",
				xmlData.ID, err)
			continue
		}
		effects = append(effects, egd)
	}
	return effects, nil
}

// buildEffectGraphicData creates effect graphic data from specified
// effect XML data.
func buildEffectGraphicData(xmlEffect *EffectGraphic) (*res.EffectGraphicData, error) {
	effIcon, err := data.Icon(xmlEffect.Icon)
	if err != nil {
		return nil, fmt.Errorf("fail to retrieve effect icon: %v", err)
	}
	effData := res.EffectGraphicData{
		EffectID: xmlEffect.ID,
		IconPic:  effIcon,
	}
	return &effData, nil
}
