/*
 * hudutils.go
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

package hud

import (
	"fmt"

	"github.com/isangeles/flame/area"
	"github.com/isangeles/flame/objects"
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/item"

	"github.com/isangeles/burn/ash"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/data"
	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/object"
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
	if eqIt, ok := it.Item.(item.Equiper); ok {
		if hud.PCAvatar() != nil && hud.PCAvatar().Equipment().Equiped(eqIt) {
			s.SetColor(invSlotEqColor)
		}
	}
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

// playerObject checks if object with specified id and serial
// is under player control.
func (hud *HUD) playerObject(id, serial string) bool {
	for _, c := range hud.game.PlayerChars() {
		if c.ID() == id && c.Serial() == serial {
			return true
		}
	}
	return false
}

// nearestObject returns the object nearest to the the active player character.
func (hud *HUD) nearestObject(obs []area.Object) (nearest area.Object) {
	pc := hud.game.ActivePlayerChar()
	for _, ob := range obs {
		if ob.ID() == pc.ID() && ob.Serial() == pc.Serial() {
			continue
		}
		if nearest == nil || objects.Range(ob, pc) < objects.Range(nearest, pc) {
			nearest = ob
		}
	}
	return nearest
}

// itemErrorGraphic returns error graphic data for specified item.
func itemErrorGraphic(it item.Item) *res.ItemGraphicData {
	return &res.ItemGraphicData{
		ItemID:   it.ID(),
		Icon:     data.ErrorIcon,
		MaxStack: 100,
	}
}

// itemGraphic returns item graphic for specified item.
func itemGraphic(it item.Item) *object.ItemGraphic {
	data := res.Item(it.ID())
	if data == nil {
		data = object.DefaultItemGraphic(it)
	}
	return object.NewItemGraphic(it, data)
}
