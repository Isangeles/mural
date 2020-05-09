/*
 * avatarxml.go
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
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// Struct for representation of avatars
// XML base.
type Avatars struct {
	XMLName xml.Name `xml:"avatars"`
	Avatars []Avatar `xml:"avatar"`
}

// Struct for representation of XML avatar
// node.
type Avatar struct {
	XMLName     xml.Name     `xml:"avatar"`
	ID          string       `xml:"id,attr"`
	Serial      string       `xml:"serial,attr"`
	Portrait    string       `xml:"portrait,value"`
	Spritesheet AvatarSprite `xml:"sprite"`
}

// Struct for spritesheet node of avatar node.
type AvatarSprite struct {
	XMLName  xml.Name `xml:"sprite"`
	Head     string   `xml:"head,value"`
	Torso    string   `xml:"torso,value"`
	FullBody string   `xml:"fullbody,value"`
}

// MarshalAvatars parses specified avatars to avatars
// base XML data.
func MarshalAvatars(avs []*object.Avatar) (string, error) {
	xmlAvatarsBase := new(Avatars)
	for _, av := range avs {
		xmlAvatar := buildAvatarXML(av)
		xmlAvatarsBase.Avatars = append(xmlAvatarsBase.Avatars, xmlAvatar)
	}
	out, err := xml.Marshal(xmlAvatarsBase)
	if err != nil {
		return "", fmt.Errorf("unable to marshal avatars base: %v", err)
	}
	return string(out[:]), nil
}

// MarshalAvatar parses specified character avatar to
// XML data.
func MarshalAvatar(av *object.Avatar) (string, error) {
	return MarshalAvatars([]*object.Avatar{av})
}

// UnmarshalAvatars retrieves all avatar data from specified
// XML data.
func UnmarshalAvatars(data io.Reader) ([]*res.AvatarData, error) {
	doc, _ := ioutil.ReadAll(data)
	xmlBase := new(Avatars)
	err := xml.Unmarshal(doc, xmlBase)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal xml data: %v", err)
	}
	avatars := make([]*res.AvatarData, 0)
	for _, xmlData := range xmlBase.Avatars {
		ad, err := buildAvatarData(&xmlData)
		if err != nil {
			log.Err.Printf("xml: unmarshal avatar: %s: unable to build data: %v",
				xmlData.ID, err)
			continue
		}
		avatars = append(avatars, ad)
	}
	return avatars, nil
}

// buildAvatarXML build XML node struct for specified
// avatar.
func buildAvatarXML(av *object.Avatar) Avatar {
	xmlAvatar := Avatar{}
	xmlAvatar.ID = av.ID()
	xmlAvatar.Serial = av.Serial()
	xmlAvatar.Portrait = av.Data().Portrait
	xmlAvatar.Spritesheet.Head = av.Data().Head
	xmlAvatar.Spritesheet.Torso = av.Data().Torso
	xmlAvatar.Spritesheet.FullBody = av.Data().FullBody
	return xmlAvatar
}

// buildAvatarXML build XML node struct for specified
// avatar.
func buildAvatarDataXML(avData *res.AvatarData) Avatar {
	xmlAvatar := Avatar{}
	xmlAvatar.ID = avData.ID
	xmlAvatar.Serial = avData.Serial
	xmlAvatar.Portrait = avData.Portrait
	xmlAvatar.Spritesheet.Head = avData.Head
	xmlAvatar.Spritesheet.Torso = avData.Torso
	xmlAvatar.Spritesheet.FullBody = avData.FullBody
	return xmlAvatar
}

// buildAvatarData build data avatar from specified XML data.
func buildAvatarData(avXML *Avatar) (*res.AvatarData, error) {
	avData := res.AvatarData{
		ID:       avXML.ID,
		Serial:   avXML.Serial,
		Portrait: avXML.Portrait,
		Head:     avXML.Spritesheet.Head,
		Torso:    avXML.Spritesheet.Torso,
	}
	return &avData, nil
}

// buildStaticAvatarData build new static avatar data for specified
// character from specified XML data.
func buildStaticAvatarData(avXML *Avatar) (*res.AvatarData, error) {
	avData := res.AvatarData{
		ID:       avXML.ID,
		Serial:   avXML.Serial,
		Portrait: avXML.Portrait,
		FullBody: avXML.Spritesheet.FullBody,
	}
	return &avData, nil
}
