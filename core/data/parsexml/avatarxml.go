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

// marshalAvatarXML parses specified character avatar to
// XML in form of bytes.
func MarshalAvatar(av *objects.Avatar) (string, error) {
	xmlAvatarsBase := new(AvatarsBaseXML)
	xmlAvatar := AvatarXML{}
	xmlAvatar.ID = av.ID()
	xmlAvatar.Portrait = av.PortraitName()
	xmlAvatar.Spritesheet = av.SpritesheetName()
	xmlAvatarsBase.Avatars = append(xmlAvatarsBase.Avatars, xmlAvatar)
	out, err := xml.Marshal(xmlAvatarsBase)
	if err != nil {
		return "", fmt.Errorf("fail_to_marshal_avatars_base:%v", err)
	}
	return string(out[:]), nil
}

// unmarshalAvatarsBaseXML parses specified XML data to game
// characters avatars.
func UnmarshalAvatarsBase(data io.Reader) ([]*objects.Avatar, error) {
	// TODO: unmarshal XML avatars base.
	return nil, fmt.Errorf("unsuported_yet")
}

