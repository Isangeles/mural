/*
 * import.go
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

package imp

import (
	"fmt"
	"path/filepath"

	"github.com/isangeles/flame/module"

	"github.com/isangeles/mural/core/data/res"
)

// LoadModuleResources loads all data(like items, skill, etc.) from module
// resources files.
func LoadModuleResources(mod *module.Module) error {
	// Paths.
	obPath := filepath.Join(mod.Conf().Path, "gui/objects")
	itPath := filepath.Join(mod.Conf().Path, "gui/items")
	effPath := filepath.Join(mod.Conf().Path, "gui/effects")
	skillPath := filepath.Join(mod.Conf().Path, "gui/skills")
	// Objects graphics.
	obGraphics, err := ImportObjectsGraphicsDir(obPath)
	if err != nil {
		return fmt.Errorf("unable to import objects graphics: %v", err)
	}
	res.SetObjects(obGraphics)
	// Items graphics.
	itGraphics, err := ImportItemsGraphicsDir(itPath)
	if err != nil {
		return fmt.Errorf("unable to import items graphics: %v", err)
	}
	res.SetItems(itGraphics)
	// Effects graphic.
	effGraphics, err := ImportEffectsGraphicsDir(effPath)
	if err != nil {
		return fmt.Errorf("unable to import effects graphics: %v", err)
	}
	res.SetEffects(effGraphics)
	// Skills graphic.
	skillGraphics, err := ImportSkillsGraphicsDir(skillPath)
	if err != nil {
		return fmt.Errorf("unable to import skills graphics: %v", err)
	}
	res.SetSkills(skillGraphics)
	return nil
}

// LoadChapterResources loads all data from chapter
// resources files.
func LoadChapterResources(chapter *module.Chapter) error {
	// Paths.
	avsPath := filepath.Join(chapter.Module().Conf().Path, "gui/chapters",
		chapter.Conf().ID, "npc")
	obsPath := filepath.Join(chapter.Module().Conf().Path, "gui/chapters",
		chapter.Conf().ID, "objects")
	// Avatars.
	avs, err := ImportAvatarsDataDir(avsPath)
	if err != nil {
		return fmt.Errorf("unable to import chapter avatars: %v", err)
	}
	res.SetAvatars(avs)
	// Objects graphics.
	obGraphics, err := ImportObjectsGraphicsDir(obsPath)
	if err != nil {
		return fmt.Errorf("unable to import objects graphics: %v", err)
	}
	res.SetObjects(append(res.Objects(), obGraphics...))
	return nil
}
