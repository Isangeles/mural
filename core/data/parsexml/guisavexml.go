/*
 * guisavexml.go
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
	
	"github.com/isangeles/mural/core/data/res"
)

// Struct for XML GUI save.
type GUISaveXML struct {
	XMLName     xml.Name   `xml:"save"`
	Name        string     `xml:"name,attr"`
	PlayersNode PlayersXML `xml:"players"`
        CameraNode  CameraXML  `xml:"camera"`
}

// Struct for PCs avatars node.
type PlayersXML struct {
	XMLName xml.Name    `xml:"players"`
	Players []PlayerXML `xml:"player"`
}

// Struct for PC node.
type PlayerXML struct {
	XMLName   xml.Name     `xml:"player"`
	Avatar    AvatarXML    `xml:"avatar"`
	Inventory InventoryXML `xml:"inventory"`
}

// Struct for GUI camera XML node.
type CameraXML struct {
	XMLName  xml.Name `xml:"camera"`
	Position string   `xml:"position,attr"`
}

// Struct for inventory node of avatar node.
type InventoryXML struct {
	XMLName xml.Name  `xml:"inventory"`
	Slots   []SlotXML `xml:"slot"`
}

// Struct for slot node of inventory node.
type SlotXML struct {
	XMLName xml.Name `xml:"slot"`
	ID      int      `xml:"id,attr"`
	Content string   `xml:"content,attr"`
}

// MarshalGUISave parses specified game save to XML
// data.
func MarshalGUISave(save *res.GUISave) (string, error) {
	xmlGUI := new(GUISaveXML)
	xmlGUI.Name = save.Name
	xmlGUI.PlayersNode.Players = make([]PlayerXML, 0)
	for _, pcData := range save.PlayersData {
		xmlPC := new(PlayerXML)
		xmlPC.Avatar = buildAvatarDataXML(pcData.Avatar)
		for serial, slot := range pcData.InvSlots {
			xmlSlot := new(SlotXML)
			xmlSlot.ID = slot
			xmlSlot.Content = serial
			xmlPC.Inventory.Slots = append(xmlPC.Inventory.Slots, *xmlSlot)
		}
		xmlGUI.PlayersNode.Players = append(xmlGUI.PlayersNode.Players, *xmlPC)
	}
	xmlGUI.CameraNode.Position = fmt.Sprintf("%fx%f", save.CameraPosX,
		save.CameraPosY)
	out, err := xml.Marshal(xmlGUI)
	if err != nil {
		return "", fmt.Errorf("fail_to_marshal_xml_data:%v",
			err)
	}
	return string(out[:]), nil
}

// Unmarshal parses XML data to GUI save struct.
func UnmarshalGUISave(data io.Reader) (*GUISaveXML, error) {
	doc, _ := ioutil.ReadAll(data)
	xmlGUISave := new(GUISaveXML)
	err := xml.Unmarshal(doc, xmlGUISave)
	if err != nil {
		return nil, fmt.Errorf("fail_to_unmarshal_xml_data:%v",
			err)
	}
	return xmlGUISave, nil
}
