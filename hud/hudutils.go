/*
 * hudutils.go
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

package hud

import (
	"fmt"
	
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/object/item"
	
	"github.com/isangeles/mtk"
	
	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/object"
)

// insertSlotItem inserts specified item to specified slot.
func (hud *HUD) insertSlotItem(it *object.ItemGraphic, s *mtk.Slot) {
	s.AddValues(it)
	s.SetInfo(hud.itemInfo(it.Item))
	s.SetIcon(it.Icon())
}

// itemInfo returns formated string with
// informations about specified item.
func (hud *HUD) itemInfo(it item.Item) string {
	// Retrieve translated item name and info.
	langPath := hud.game.Module().Conf().ItemsLangPath()
	nameInfo := lang.AllText(langPath, it.ID())
	// Compose info for item type.
	info := ""
	switch i := it.(type) {
	case *item.Weapon:
		infoForm := "%s\n%d-%d"
		dmgMin, dmgMax := i.Damage()
		info = fmt.Sprintf(infoForm, nameInfo[0],
			dmgMin, dmgMax)
	case *item.Misc:
		infoForm := "%s"
		info = fmt.Sprintf(infoForm, nameInfo[0])
	}
	if config.Debug() { // add serial ID info
		info = fmt.Sprintf("%s\n[%s_%s]", info,
			it.ID(), it.Serial())
	}
	return info
}
