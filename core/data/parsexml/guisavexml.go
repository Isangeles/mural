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
	
	flamexml "github.com/isangeles/flame/core/data/parsexml"
	
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
	MenuBar   BarXML       `xml:"bar"`
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

// Struct for menu bar node of avatar node.
type BarXML struct {
	XMLName xml.Name  `xml:"bar"`
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
	// Players.
	for _, pcData := range save.PlayersData {
		xmlPC := new(PlayerXML)
		// Avatar.
		xmlPC.Avatar = buildAvatarDataXML(pcData.Avatar)
		// Layouts.
		for serial, slot := range pcData.InvSlots {
			xmlSlot := SlotXML{
				ID:      slot,
				Content: serial,
			}
			xmlPC.Inventory.Slots = append(xmlPC.Inventory.Slots, xmlSlot)
		}
		for serial, slot := range pcData.BarSlots {
			xmlSlot := SlotXML{
				ID:      slot,
				Content: serial,
			}
			xmlPC.MenuBar.Slots = append(xmlPC.MenuBar.Slots, xmlSlot)
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
func UnmarshalGUISave(data io.Reader) (*res.GUISave, error) {
	doc, _ := ioutil.ReadAll(data)
	xmlGUISave := new(GUISaveXML)
	err := xml.Unmarshal(doc, xmlGUISave)
	if err != nil {
		return nil, fmt.Errorf("fail_to_unmarshal_xml_data:%v",
			err)
	}
	save, err := buildGUISave(xmlGUISave)
	if err != nil {
		return nil, fmt.Errorf("fail_to_build_data:%v", err)
	}
	return save, nil
}

// buildGUISave builds GUI save from specified XML data.
func buildGUISave(xmlSave *GUISaveXML) (*res.GUISave, error) {
	save := new(res.GUISave)
	// Save name.
	save.Name = xmlSave.Name
	// Players.
	for _, xmlPC := range xmlSave.PlayersNode.Players {
		pcData := new(res.PlayerSave)
		// Avatar.
		avData, err := buildAvatarData(&xmlPC.Avatar)
		if err != nil {
			return nil, fmt.Errorf("player:%s_%s:fail_to_load_player_avatar:%v",
				pcData.Avatar.ID, pcData.Avatar.Serial, err)
		}
		pcData.Avatar = avData
		// Inventory layout.
		pcData.InvSlots = make(map[string]int)
		for _, xmlSlot := range xmlPC.Inventory.Slots {
			pcData.InvSlots[xmlSlot.Content] = xmlSlot.ID
		}
		save.PlayersData = append(save.PlayersData, pcData)
		// Menu bar layout.
		pcData.BarSlots = make(map[string]int)
		for _, xmlSlot := range xmlPC.MenuBar.Slots {
			pcData.BarSlots[xmlSlot.Content] = xmlSlot.ID
		}
	}
	// Camera position.
	camX, camY, err := flamexml.UnmarshalPosition(xmlSave.CameraNode.Position)
	if err != nil {
		return nil, fmt.Errorf("fail_to_unmarshal_camera_position:%v",
			err)
	}
	save.CameraPosX, save.CameraPosY = camX, camY
	return save, nil
}
