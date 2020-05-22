/*
 * skillgraphic.go
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

package res

import (
	"encoding/xml"
)

type SkillGraphicsData struct {
	XMLName xml.Name           `xml:"skill-graphics" json:"-"`
	Skills  []SkillGraphicData `xml:"skill-graphic" json:"skill-graphics"`
}

// Struct for skill graphical data.
type SkillGraphicData struct {
	SkillID         string `xml:"id,attr" json:"id"`
	Icon            string `xml:"icon,attr" json:"icon"`
	ActivationAudio string `xml:"activation-audio,attr" json:"activation-audio"`
	ActivationAnim  string `xml:"activation-anim,attr" json:"actvation-anim"`
}
