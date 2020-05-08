/*
 * res.go
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

// Package with pre-loaded game resources
// like item graphic data, textures, etc.
package res

var (
	avatarsData map[string]*AvatarData
	objectsData map[string]*ObjectGraphicData
	itemsData   map[string]*ItemGraphicData
	effectsData map[string]*EffectGraphicData
	skillsData  map[string]*SkillGraphicData
)

// On init.
func init() {
	avatarsData = make(map[string]*AvatarData)
	objectsData = make(map[string]*ObjectGraphicData)
	itemsData = make(map[string]*ItemGraphicData)
	effectsData = make(map[string]*EffectGraphicData)
	skillsData = make(map[string]*SkillGraphicData)
}

// Avatar returns avatar data for character
// with specified ID.
func Avatar(id string) *AvatarData {
	return avatarsData[id]
}

// Object returns object graphic data for
// object with specified ID.
func Object(id string) *ObjectGraphicData {
	return objectsData[id]
}

// Item returns graphic data for item
// with specified ID.
func Item(itemID string) *ItemGraphicData {
	return itemsData[itemID]
}

// Effect returns graphic data for effect
// with specified ID.
func Effect(id string) *EffectGraphicData {
	return effectsData[id]
}

// Skill returns graphic data for skill
// with specified ID.
func Skill(id string) *SkillGraphicData {
	return skillsData[id]
}

// AddObjectsData adds specified data to objects graphic
// recources data.
func AddObjectsData(data ...*ObjectGraphicData) {
	if objectsData == nil {
		objectsData = make(map[string]*ObjectGraphicData)
	}
	for _, d := range data {
		objectsData[d.ID] = d
	}
}

// AddItemsData adds specified data to items
// resources data.
func AddItemsData(data ...*ItemGraphicData) {
	if itemsData == nil {
		itemsData = make(map[string]*ItemGraphicData)
	}
	for _, d := range data {
		itemsData[d.ItemID] = d
	}
}

// AddEffectsData adds specified data to effects
// resources data.
func AddEffectsData(data ...*EffectGraphicData) {
	if effectsData == nil {
		effectsData = make(map[string]*EffectGraphicData)
	}
	for _, d := range data {
		effectsData[d.EffectID] = d
	}
}

// AddSkillsData adds specified data to skills
// resources data.
func AddSkillsData(data ...*SkillGraphicData) {
	if skillsData == nil {
		skillsData = make(map[string]*SkillGraphicData)
	}
	for _, d := range data {
		skillsData[d.SkillID] = d
	}
}

// AddAvatarData adds specified data to
// avatars resources.
func AddAvatarData(data ...*AvatarData) {
	if avatarsData == nil {
		avatarsData = make(map[string]*AvatarData)
	}
	for _, d := range data {
		avatarsData[d.ID] = d
	}
}

// SetObjectsData sets specified data as objects graphic
// recources data.
func SetObjectsData(data []*ObjectGraphicData) {
	if objectsData == nil {
		objectsData = make(map[string]*ObjectGraphicData)
	}
	for _, d := range data {
		objectsData[d.ID] = d
	}
}

// SetItemsData sets specified data as items
// resources data.
func SetItemsData(data []*ItemGraphicData) {
	if itemsData == nil {
		itemsData = make(map[string]*ItemGraphicData)
	}
	for _, d := range data {
		itemsData[d.ItemID] = d
	}
}

// SetEffectsData sets specified data as effects
// resources data.
func SetEffectsData(data []*EffectGraphicData) {
	if effectsData == nil {
		effectsData = make(map[string]*EffectGraphicData)
	}
	for _, d := range data {
		effectsData[d.EffectID] = d
	}
}

// SetSkillsData sets specified data as skills
// resources data.
func SetSkillsData(data []*SkillGraphicData) {
	if skillsData == nil {
		skillsData = make(map[string]*SkillGraphicData)
	}
	for _, d := range data {
		skillsData[d.SkillID] = d
	}
}

// SetAvatarData sets specified data as
// avatars resources.
func SetAvatarData(data []*AvatarData) {
	if avatarsData == nil {
		avatarsData = make(map[string]*AvatarData)
	}
	for _, d := range data {
		avatarsData[d.ID] = d
	}
}
