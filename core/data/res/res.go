/*
 * loadgamemenu.go
 *
 * Copyright 2019 Dariusz Sikora <dev@isangeles.pl>
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
// like items data, saves, etc.
package res

var (
	avatarsData map[string]*AvatarData
	itemsData   map[string]*ItemGraphicData
	effectsData map[string]*EffectGraphicData
)

// Avatar returns avatar data for character
// with specified ID.
func Avatar(id string) *AvatarData {
	return avatarsData[id]
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

// SetAvatarsData sets specified map with avatars
// data as loaded resources data.
func SetAvatarsData(data map[string]*AvatarData) {
	avatarsData = data
}

// SetItemsData sets specified map with items
// graphic data as loaded resources data.
func SetItemsData(data map[string]*ItemGraphicData) {
	itemsData = data
}

// SetEffectsData sets specified map with effects
// data as loaded resources data.
func SetEffectsData(data map[string]*EffectGraphicData) {
	effectsData = data
}
