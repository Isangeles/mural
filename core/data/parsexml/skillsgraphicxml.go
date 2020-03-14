/*
 * skillgraphicxml.go
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

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/log"
)

// Struct for skill graphics base XML node.
type SkillGraphics struct {
	XMLName  xml.Name       `xml:"skill-graphics"`
	Graphics []SkillGraphic `xml:"skill-graphic"`
}

// Struct for skill graphic XML node.
type SkillGraphic struct {
	XMLName    xml.Name        `xml:"skill-graphic"`
	ID         string          `xml:"id,attr"`
	Icon       string          `xml:"icon,attr"`
	Animations SkillAnimations `xml:"animations"`
	Audio      SkillAudio      `xml:"audio"`
}

// Struct for skill animations XML node.
type SkillAnimations struct {
	XMLName    xml.Name `xml:"animations"`
	Cast       string   `xml:"cast,value"`
	Activation string   `xml:"activation,value"`
}

// Struct for skill audio XML node.
type SkillAudio struct {
	XMLName    xml.Name `xml:"audio"`
	Cast       string   `xml:"cast,value"`
	Activation string   `xml:"activation,value"`
}

// UnmarshalSkillGraphics retrieves all skills graphic data
// for speicfied XML data.
func UnmarshalSkillGraphics(data io.Reader) ([]*res.SkillGraphicData, error) {
	doc, _ := ioutil.ReadAll(data)
	xmlBase := new(SkillGraphics)
	err := xml.Unmarshal(doc, xmlBase)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal xml data: %v", err)
	}
	skills := make([]*res.SkillGraphicData, 0)
	for _, xmlData := range xmlBase.Graphics {
		sgd, err := buildSkillGraphicData(&xmlData)
		if err != nil {
			log.Err.Printf("xml: unmarshal skill graphic: %s: unable to build data: %v",
				xmlData.ID, err)
			continue
		}
		skills = append(skills, sgd)
	}
	return skills, nil
}

// buildSkillGraphicData creates skill graphic data from specified
// skill XML data.
func buildSkillGraphicData(xmlSkill *SkillGraphic) (*res.SkillGraphicData, error) {
	icon := data.Icon(xmlSkill.Icon)
	if icon == nil {
		return nil, fmt.Errorf("unable to retrieve skill icon: %s",
			xmlSkill.Icon)
	}
	activationAnim := UnmarshalAvatarAnim(xmlSkill.Animations.Activation)
	skillData := res.SkillGraphicData{
		SkillID:        xmlSkill.ID,
		IconPic:        icon,
		ActivationAnim: int(activationAnim),
	}
	if xmlSkill.Audio.Activation != "" {
		activeAudio, err := data.AudioEffect(xmlSkill.Audio.Activation)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve skill audio: %v", err)
		}
		skillData.ActivationAudio = activeAudio
	}
	return &skillData, nil
}
