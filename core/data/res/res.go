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
	avatars        map[string]*AvatarData
	objects        map[string]*ObjectGraphicData
	items          map[string]*ItemGraphicData
	effects        map[string]*EffectGraphicData
	skills         map[string]*SkillGraphicData
)

// On init.
func init() {
	avatars = make(map[string]*AvatarData)
	objects = make(map[string]*ObjectGraphicData)
	items = make(map[string]*ItemGraphicData)
	effects = make(map[string]*EffectGraphicData)
	skills = make(map[string]*SkillGraphicData)
}

// Avatar returns avatar data for character
// with specified ID.
func Avatar(id string) *AvatarData {
	return avatars[id]
}

// Object returns object graphic data for
// object with specified ID.
func Object(id string) *ObjectGraphicData {
	return objects[id]
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
func Avatars() (data []*AvatarData) {
	for _, ad := range avatars {
		data = append(data, ad)
	}
	return
}

// Objects returns all object graphic resources.
func Objects() (data []*ObjectGraphicData) {
	for _, od := range objects {
		data = append(data, od)
	}
	return
}

// Items returns all item resources.
func Items() (data []*ItemGraphicData) {
	for _, id := range items {
		data = append(data, id)
	}
	return
}

// Effects returns all effect resources.
func Effects() (data []*EffectGraphicData) {
	for _, ed := range effects {
		data = append(data, ed)
	}
	return
}

// Skills returns all skill resources.
func Skills() (data []*SkillGraphicData) {
	for _, sd := range skills {
		data = append(data, sd)
	}
	return
}

// SetAvatars sets specified data as
// avatars resources.
func SetAvatars(data []AvatarData) {
	for i, _ := range data {
		avatars[data[i].ID] = &data[i]
	}
}

// SetObjects sets specified data as objects graphic
// recources data.
func SetObjects(data []ObjectGraphicData) {
	for i, _ := range data {
		objects[data[i].ID] = &data[i]
	}
}

// SetItems sets specified data as items
// resources data.
func SetItems(data []*ItemGraphicData) {
	for _, d := range data {
		items[d.ItemID] = d
	}
}

// SetEffects sets specified data as effects
// resources data.
func SetEffects(data []*EffectGraphicData) {
	for _, d := range data {
		effects[d.EffectID] = d
	}
}

// SetSkills sets specified data as skills
// resources data.
func SetSkills(data []*SkillGraphicData) {
	for _, d := range data {
		skills[d.SkillID] = d
	}
}
