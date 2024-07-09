/*
 * res.go
 *
 * Copyright 2019-2024 Dariusz Sikora <ds@isangeles.dev>
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

// Package with pre-loaded game resources
// like item graphic data, textures, etc.
package res

import (
	flameres "github.com/isangeles/flame/data/res"
)

var (
	Avatars          []AvatarData
	Items            []ItemGraphicData
	Effects          []EffectGraphicData
	Skills           []SkillGraphicData
	TranslationBases []*flameres.TranslationBaseData
)

// Avatar returns avatar data for character
// with specified ID.
func Avatar(id string) *AvatarData {
	for _, d := range Avatars {
		if d.ID == id {
			return &d
		}
	}
	return nil
}

// Item returns graphic data for item
// with specified ID.
func Item(id string) *ItemGraphicData {
	for _, d := range Items {
		if d.ItemID == id {
			return &d
		}
	}
	return nil
}

// Effect returns graphic data for effect
// with specified ID.
func Effect(id string) *EffectGraphicData {
	for _, d := range Effects {
		if d.EffectID == id {
			return &d
		}
	}
	return nil
}

// Skill returns graphic data for skill
// with specified ID.
func Skill(id string) *SkillGraphicData {
	for _, d := range Skills {
		if d.SkillID == id {
			return &d
		}
	}
	return nil
}

// AddTranslationBases adds all translation bases
// to the translation resources.
func AddTranslationBases(bases []flameres.TranslationBaseData) {
	for _, b := range bases {
		TranslationBases = append(TranslationBases, &b)
	}
}

