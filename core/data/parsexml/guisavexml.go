/*
 * guisavexml.go
 *
 * Copyright 2018-2020 Dariusz Sikora <dev@isangeles.pl>
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
type Save struct {
	XMLName xml.Name `xml:"save"`
	Name    string   `xml:"name,attr"`
	Players []Player `xml:"players>player"`
	Camera  Camera   `xml:"camera"`
}

// Struct for PC node.
type Player struct {
	XMLName   xml.Name  `xml:"player"`
	Avatar    Avatar    `xml:"avatar"`
	Inventory Inventory `xml:"inventory"`
	MenuBar   Bar       `xml:"bar"`
}

// Struct for GUI camera XML node.
type Camera struct {
	XMLName  xml.Name `xml:"camera"`
	X        float64  `xml:"x,attr"`
	Y        float64  `xml:"y,attr"`
}

// Struct for inventory node of avatar node.
type Inventory struct {
	XMLName xml.Name `xml:"inventory"`
	Slots   []Slot   `xml:"slot"`
}

// Struct for menu bar node of avatar node.
type Bar struct {
	XMLName xml.Name `xml:"bar"`
	Slots   []Slot   `xml:"slot"`
}

// Struct for slot node of inventory node.
type Slot struct {
	XMLName xml.Name `xml:"slot"`
	ID      int      `xml:"id,attr"`
	Content string   `xml:"content,attr"`
}

// MarshalGUISave parses specified game save to XML
// data.
func MarshalGUISave(save *res.GUISave) (string, error) {
	xmlGUI := new(Save)
	xmlGUI.Name = save.Name
	xmlGUI.Players = make([]Player, 0)
	// Players.
	for _, pcData := range save.PlayersData {
		xmlPC := new(Player)
		// Avatar.
		xmlPC.Avatar = buildAvatarDataXML(pcData.Avatar)
		// Layouts.
		for serial, slot := range pcData.InvSlots {
			xmlSlot := Slot{
				ID:      slot,
				Content: serial,
			}
			xmlPC.Inventory.Slots = append(xmlPC.Inventory.Slots, xmlSlot)
		}
		for serial, slot := range pcData.BarSlots {
			xmlSlot := Slot{
				ID:      slot,
				Content: serial,
			}
			xmlPC.MenuBar.Slots = append(xmlPC.MenuBar.Slots, xmlSlot)
		}
		xmlGUI.Players = append(xmlGUI.Players, *xmlPC)
	}
	xmlGUI.Camera.X, xmlGUI.Camera.Y = save.CameraPosX, save.CameraPosY
	out, err := xml.Marshal(xmlGUI)
	if err != nil {
		return "", fmt.Errorf("fail to marshal xml data: %v", err)
	}
	return string(out[:]), nil
}

// Unmarshal parses XML data to GUI save struct.
func UnmarshalGUISave(data io.Reader) (*res.GUISave, error) {
	doc, _ := ioutil.ReadAll(data)
	xmlGUISave := new(Save)
	err := xml.Unmarshal(doc, xmlGUISave)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal xml data: %v", err)
	}
	save, err := buildGUISave(xmlGUISave)
	if err != nil {
		return nil, fmt.Errorf("fail to build data: %v", err)
	}
	return save, nil
}

// buildGUISave builds GUI save from specified XML data.
func buildGUISave(xmlSave *Save) (*res.GUISave, error) {
	save := new(res.GUISave)
	// Save name.
	save.Name = xmlSave.Name
	// Players.
	for _, xmlPC := range xmlSave.Players {
		pcData := new(res.PlayerSave)
		// Avatar.
		avData, err := buildAvatarData(&xmlPC.Avatar)
		if err != nil {
			return nil, fmt.Errorf("player: %s#%s: fail to load player avatar: %v",
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
	save.CameraPosX, save.CameraPosY = xmlSave.Camera.X, xmlSave.Camera.Y
	return save, nil
}
