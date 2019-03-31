/*
 * import.go
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

package imp

import (
	"fmt"

	"github.com/isangeles/flame/core/module"

	"github.com/isangeles/mural/core/data/res"
)

// LoadModuleResources loads all data(like items, skill, etc.) from module
// resources files.
func LoadModuleResources(mod *module.Module) error {
	// Objectsd graphics.
	objectsData := make(map[string]*res.ObjectGraphicData)
	obGraphics, err := ImportObjectsGraphicsDir(mod.Conf().ObjectsPath())
	if err != nil {
		return fmt.Errorf("fail_to_import_objects_graphics:%v", err)
	}
	for _, obGraphic := range obGraphics {
		objectsData[obGraphic.ID] = obGraphic
	}
	res.SetObjectsData(objectsData)
	// Items graphics.
	itemsData := make(map[string]*res.ItemGraphicData)
	itGraphics, err := ImportItemsGraphicsDir(mod.Conf().ItemsPath())
	if err != nil {
		return fmt.Errorf("fail_to_import_items_graphics:%v", err)
	}
	for _, itGraphic := range itGraphics {
		itemsData[itGraphic.ItemID] = itGraphic
	}
	res.SetItemsData(itemsData)
	// Effects graphic.
	effectsData := make(map[string]*res.EffectGraphicData)
	effGraphics, err := ImportEffectsGraphicsDir(mod.Conf().EffectsPath())
	if err != nil {
		return fmt.Errorf("fail_to_import_effects_graphics:%v", err)
	}
	for _, effGraphic := range effGraphics {
		effectsData[effGraphic.EffectID] = effGraphic
	}
	res.SetEffectsData(effectsData)
	// Skills graphic.
	skillsData := make(map[string]*res.SkillGraphicData)
	skillGraphics, err := ImportSkillsGraphicsDir(mod.Conf().SkillsPath())
	if err != nil {
		return fmt.Errorf("fail_to_import_skills_graphics:%v", err)
	}
	for _, skillGraphic := range skillGraphics {
		skillsData[skillGraphic.SkillID] = skillGraphic
	}
	res.SetSkillsData(skillsData)
	return nil
}

// LoadChapterResources loads all data from chapter
// resources files.
func LoadChapterResources(chapter *module.Chapter) error {
	// Avatars.
	avatarsData := make(map[string]*res.AvatarData)
	avs, err := ImportAvatarsDataDir(chapter.Conf().NPCPath())
	if err != nil {
		return fmt.Errorf("fail_to_import_chapter_avatars:%v", err)
	}
	for _, av := range avs {
		avatarsData[av.ID] = av
	}
	res.SetAvatarsData(avatarsData)
	return nil
}
