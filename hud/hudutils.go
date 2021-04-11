/*
 * hudutils.go
 *
 * Copyright 2019-2021 Dariusz Sikora <dev@isangeles.pl>
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

	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/item"

	"github.com/isangeles/burn/ash"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// RunScript executes specified script, in case
// of error sends err message to Mural error log.
func (hud *HUD) RunScript(s *ash.Script) {
	err := ash.Run(s)
	if err != nil {
		log.Err.Printf("ci: unable to run script: %v", err)
		return
	}
}

// insertSlotItem inserts specified item to specified slot.
func (hud *HUD) insertSlotItem(it *object.ItemGraphic, s *mtk.Slot) {
	s.AddValues(it)
	s.SetInfo(hud.itemInfo(it.Item))
	s.SetIcon(it.Icon())
}

// itemInfo returns formated string with
// informations about specified item.
func (hud *HUD) itemInfo(it item.Item) string {
	// Compose info for item type.
	info := ""
	switch i := it.(type) {
	case *item.Weapon:
		infoForm := "%s\n%d-%d"
		dmgMin, dmgMax := i.Damage()
		info = fmt.Sprintf(infoForm, lang.Text(i.ID()),
			dmgMin, dmgMax)
	case *item.Misc:
		infoForm := "%s"
		info = fmt.Sprintf(infoForm, lang.Text(i.ID()))
	}
	if config.Debug { // add serial ID info
		info = fmt.Sprintf("%s\n[%s_%s]", info,
			it.ID(), it.Serial())
	}
	return info
}

// itemErrorGraphic returns error graphic data for specified item.
func itemErrorGraphic(it item.Item) *res.ItemGraphicData {
	return &res.ItemGraphicData{
		ItemID:   it.ID(),
		Icon:     data.ErrorIcon,
		MaxStack: 100,
	}
}
