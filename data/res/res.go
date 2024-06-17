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
	"sync"

	flameres "github.com/isangeles/flame/data/res"
)

var (
	avatars          *sync.Map
	items            map[string]*ItemGraphicData
	effects          map[string]*EffectGraphicData
	skills           map[string]*SkillGraphicData
	translationBases map[string]*flameres.TranslationBaseData
)

// On init.
func init() {
	avatars = new(sync.Map)
	items = make(map[string]*ItemGraphicData)
	effects = make(map[string]*EffectGraphicData)
	skills = make(map[string]*SkillGraphicData)
	translationBases = make(map[string]*flameres.TranslationBaseData)
}

// Avatar returns avatar data for character
// with specified ID.
func Avatar(id string) *AvatarData {
	ad, ok := avatars.Load(id)
	if !ok {
		return nil
	}
	return ad.(*AvatarData)
}

// Item returns graphic data for item
// with specified ID.
func Item(itemID string) *ItemGraphicData {
	return items[itemID]
}

// Effect returns graphic data for effect
// with specified ID.
func Effect(id string) *EffectGraphicData {
	return effects[id]
}

// Skill returns graphic data for skill
// with specified ID.
func Skill(id string) *SkillGraphicData {
	return skills[id]
}

// Avatars returns all avatars resources.
func Avatars() (data []AvatarData) {
	addAvatar := func(k, v interface{}) bool {
		ad, ok := v.(*AvatarData)
		if ok {
			data = append(data, *ad)
		}
		return true
	}
	avatars.Range(addAvatar)
	return
}

// Items returns all item resources.
func Items() (data []ItemGraphicData) {
	for _, id := range items {
		data = append(data, *id)
	}
	return
}

// Effects returns all effect resources.
func Effects() (data []EffectGraphicData) {
	for _, ed := range effects {
		data = append(data, *ed)
	}
	return
}

// Skills returns all skill resources.
func Skills() (data []SkillGraphicData) {
	for _, sd := range skills {
		data = append(data, *sd)
	}
	return
}

// TranslationBases returns all translation bases.
func TranslationBases() (data []*flameres.TranslationBaseData) {
	for _, tbd := range translationBases {
		data = append(data, tbd)
	}
	return
}

// SetAvatars sets specified data as
// avatars resources.
func SetAvatars(data []AvatarData) {
	for i, _ := range data {
		avatars.Store(data[i].ID, &data[i])
	}
}

// SetItems sets specified data as items
// resources data.
func SetItems(data []ItemGraphicData) {
	for i, _ := range data {
		items[data[i].ItemID] = &data[i]
	}
}

// SetEffects sets specified data as effects
// resources data.
func SetEffects(data []EffectGraphicData) {
	for i, _ := range data {
		effects[data[i].EffectID] = &data[i]
	}
}

// SetSkills sets specified data as skills
// resources data.
func SetSkills(data []SkillGraphicData) {
	for i, _ := range data {
		skills[data[i].SkillID] = &data[i]
	}
}

// SetTranslationBases sets specified data as translation
// resources.
func SetTranslationBases(data []flameres.TranslationBaseData) {
	for i, _ := range data {
		translationBases[data[i].ID] = &data[i]
	}
}
