/*
 * avatarxml.go
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

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// Struct for representation of avatars
// XML base.
type AvatarsBaseXML struct {
	XMLName xml.Name    `xml:"base"`
	Avatars   []AvatarXML `xml:"avatar"`
}

// Struct for representation of XML avatar
// node.
type AvatarXML struct {
	XMLName     xml.Name            `xml:"avatar"`
	ID          string              `xml:"id,attr"`
	Serial      string              `xml:"serial,attr"`
	Portrait    string              `xml:"portrait,value"`
	Spritesheet AvatarSpriteNodeXML `xml:"sprite"`
}

// Struct for spritesheet node of avatar node.
type AvatarSpriteNodeXML struct {
	XMLName  xml.Name `xml:"sprite"`
	Head     string   `xml:"head,value"`
	Torso    string   `xml:"torso,value"`
	FullBody string   `xml:"fullbody,value"`
}

// MarshalAvatarsBase parses specified avatars to avatars
// base XML data.
func MarshalAvatarsBase(avs []*object.Avatar) (string, error) {
	xmlAvatarsBase := new(AvatarsBaseXML)
	for _, av := range avs {
		xmlAvatar := buildAvatarXML(av)
		xmlAvatarsBase.Avatars = append(xmlAvatarsBase.Avatars, xmlAvatar)
	}
	out, err := xml.Marshal(xmlAvatarsBase)
	if err != nil {
		return "", fmt.Errorf("fail_to_marshal_avatars_base:%v", err)
	}
	return string(out[:]), nil
}

// MarshalAvatar parses specified character avatar to
// XML data.
func MarshalAvatar(av *object.Avatar) (string, error) {
	return MarshalAvatarsBase([]*object.Avatar{av})
}

// UnmarshalAvatarsBase retrieves all avatar data from specified
// XML data.
func UnmarshalAvatarsBase(data io.Reader) ([]*res.AvatarData, error) {
	doc, _ := ioutil.ReadAll(data)
	xmlBase := new(AvatarsBaseXML)
	err := xml.Unmarshal(doc, xmlBase)
	if err != nil {
		return nil, fmt.Errorf("fail_to_unmarshal_xml_data:%v", err)
	}
	avatars := make([]*res.AvatarData, 0)
	for _, xmlData := range xmlBase.Avatars {
		ad, err := buildAvatarData(&xmlData)
		if err != nil {
			log.Err.Printf("xml:unmarshal_avatar:%s:fail_to_build_data:%v",
				xmlData.ID, err)
			continue
		}
		avatars = append(avatars, ad)
	}
	return avatars, nil
}

// buildAvatarXML build XML node struct for specified
// avatar.
func buildAvatarXML(av *object.Avatar) AvatarXML {
	xmlAvatar := AvatarXML{}
	xmlAvatar.ID = av.ID()
	xmlAvatar.Serial = av.Serial()
	xmlAvatar.Portrait = av.Data().PortraitName
	xmlAvatar.Spritesheet.Head = av.Data().SSHeadName
	xmlAvatar.Spritesheet.Torso = av.Data().SSTorsoName
	xmlAvatar.Spritesheet.FullBody = av.Data().SSFullBodyName
	return xmlAvatar
}

// buildAvatarXML build XML node struct for specified
// avatar.
func buildAvatarDataXML(avData *res.AvatarData) AvatarXML {
	xmlAvatar := AvatarXML{}
	xmlAvatar.ID = avData.CharID
	xmlAvatar.Serial = avData.CharSerial
	xmlAvatar.Portrait = avData.PortraitName
	xmlAvatar.Spritesheet.Head = avData.SSHeadName
	xmlAvatar.Spritesheet.Torso = avData.SSTorsoName
	xmlAvatar.Spritesheet.FullBody = avData.SSFullBodyName
	return xmlAvatar
}

// buildAvatarData build data avatar from specified XML data.
func buildAvatarData(avXML *AvatarXML) (*res.AvatarData, error) {
	ssHeadName := avXML.Spritesheet.Head
	ssTorsoName := avXML.Spritesheet.Torso
	portraitName := avXML.Portrait
	portraitPic, err := data.AvatarPortrait(portraitName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_portrait_picture:%v", err)
	}
	ssHeadPic, err := data.AvatarSpritesheet(ssHeadName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_head_spritesheet_picture:%v", err)
	}
	ssTorsoPic, err := data.AvatarSpritesheet(ssTorsoName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_torso_spritesheet_picture:%v", err)
	}
	avData := res.AvatarData{
		CharID:       avXML.ID,
		CharSerial:   avXML.Serial,
		PortraitName: portraitName,
		SSHeadName:   ssHeadName,
		SSTorsoName:  ssTorsoName,
		PortraitPic:  portraitPic,
		SSHeadPic:    ssHeadPic,
		SSTorsoPic:   ssTorsoPic,
	}
	return &avData, nil
}

// buildStaticAvatarData build new static avatar data for specified
// character from specified XML data.
func buildStaticAvatarData(avXML *AvatarXML) (*res.AvatarData, error) {
	ssFullBodyName := avXML.Spritesheet.FullBody
	portraitName := avXML.Portrait
	portraitPic, err := data.AvatarPortrait(portraitName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_portrait_picture:%v",
			err)
	}
	ssFullBodyPic, err := data.AvatarSpritesheet(ssFullBodyName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_head_spritesheet_picture:%v",
			err)
	}
	avData := res.AvatarData{
		CharID: avXML.ID,
		CharSerial: avXML.Serial,
		PortraitName: portraitName,
		SSFullBodyName: ssFullBodyName,
		PortraitPic: portraitPic,
		SSFullBodyPic: ssFullBodyPic,
	}
	return &avData, nil
}
