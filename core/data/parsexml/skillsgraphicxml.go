/*
 * skillgraphicxml.go
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

// Struct for skill graphics base XML node.
type SkillsGraphicsBaseXML struct {
	XMLName xml.Name          `xml:"base"`
	Nodes   []SkillGraphicXML `xml:"skill"`
}

// Struct for skill graphic XML node.
type SkillGraphicXML struct {
	XMLName    xml.Name           `xml:"skill"`
	ID         string             `xml:"id,attr"`
	Icon       string             `xml:"icon,attr"`
	Animations SkillAnimationsXML `xml:"animations"`
}

// Struct for skill animations XML node.
type SkillAnimationsXML struct {
	XMLName    xml.Name `xml:"animations"`
	Activation string   `xml:"activation,value"`
	Cast       string   `xml:"cast,value"`
}

// UnmarshalSkillsGraphicsBase retrieves all skills graphic data
// for speicfied XML data.
func UnmarshalSkillsGraphicsBase(data io.Reader) ([]*res.SkillGraphicData, error) {
	doc, _ := ioutil.ReadAll(data)
	xmlBase := new(SkillsGraphicsBaseXML)
	err := xml.Unmarshal(doc, xmlBase)
	if err != nil {
		return nil, fmt.Errorf("fail_to_unmarshal_xml_data:%v", err)
	}
	skills := make([]*res.SkillGraphicData, 0)
	for _, xmlData := range xmlBase.Nodes {
		sgd, err := buildSkillGraphicData(&xmlData)
		if err != nil {
			log.Err.Printf("xml:unmarshal_skill_graphic:%s:fail_to_build_data:%v",
				xmlData.ID, err)
			continue
		}
		skills = append(skills, sgd)
	}
	return skills, nil
}

// buildSkillGraphicData creates skill graphic data from specified
// skill XML data.
func buildSkillGraphicData(xmlSkill *SkillGraphicXML) (*res.SkillGraphicData, error) {
	skillIcon, err := data.Icon(xmlSkill.Icon)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_skill_icon:%v", err)
	}
	activationAnim := UnmarshalAvatarAnim(xmlSkill.Animations.Activation)
	skillData := res.SkillGraphicData{
		SkillID:        xmlSkill.ID,
		IconPic:        skillIcon,
		ActivationAnim: int(activationAnim),
	}
	return &skillData, nil
}
