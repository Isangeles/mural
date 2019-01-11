/*
 * itemgraphicxml.go
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

package parsexml

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
)

// Struct for XML items grpahics base.
type ItemsGraphicsBaseXML struct {
	XMLName xml.Name             `xml:"base"`
	Items   []ItemGraphicNodeXML `xml:"item"`
}

// Struct for XML item graphic node.
type ItemGraphicNodeXML struct {
	XMLName      xml.Name `xml:"item"`
	ID           string   `xml:"id,attr"`
	Spritesheet  string   `xml:"spritesheet,attr"`
}

// UnmarshalItemsGraphicsBase parses specified XML data
// to item grphics XML nodes.
func UnmarshalItemsGraphicsBase(data io.Reader) ([]ItemGraphicNodeXML, error) {
	doc, _ := ioutil.ReadAll(data)
	itemsGraphicsXML := new(ItemsGraphicsBaseXML)
	err := xml.Unmarshal(doc, itemsGraphicsXML)
	if err != nil {
		return nil, fmt.Errorf("fail_to_unmarshal_xml_data:%v", err)
	}
	return itemsGraphicsXML.Items, nil
}
