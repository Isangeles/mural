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

	"github.com/isangeles/mural/objects"
)

// Struct for representation of avatars
// XML base.
type AvatarsBaseXML struct {
	XMLName string      `xml:"base"`
	Avatars []AvatarXML `xml:"avatar"`
}

// Struct for representation of XML avatar
// node.
type AvatarXML struct {
	XMLName     string `xml:"avatar"`
	ID          string `xml:"id,attr"`
	Portrait    string `xml:"portrait, value"`
	Spritesheet string `xml:"spritesheet, value"`
}

// MarshalAvatarsBase parses specified avatars to avatars
// base XML data.
func MarshalAvatarsBase(avs []*objects.Avatar) (string, error) {
	xmlAvatarsBase := new(AvatarsBaseXML)
	for _, av := range avs {
		xmlAvatar := AvatarXML{}
		xmlAvatar.ID = av.ID()
		xmlAvatar.Portrait = av.PortraitName()
		xmlAvatar.Spritesheet = av.SpritesheetName()
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
func MarshalAvatar(av *objects.Avatar) (string, error) {
	return MarshalAvatarsBase([]*objects.Avatar{av})
}

// UnmarshalAvatarsBase parses specified XML data to game
// characters avatars.
func UnmarshalAvatarsBase(data io.Reader) ([]AvatarXML, error) {
	doc, _ := ioutil.ReadAll(data)
	avatarsXML := new(AvatarsBaseXML)
	err := xml.Unmarshal(doc, avatarsXML)
	if err != nil {
		return nil, fmt.Errorf("fail_to_unmarshal_xml_data:%v", err)
	}
	return avatarsXML.Avatars, nil
}
