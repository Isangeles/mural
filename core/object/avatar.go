/*
 * avatar.go
 *
 * Copyright 2018-2019 Dariusz Sikora <dev@isangeles.pl>
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

package object

import (
	"fmt"

	"github.com/faiface/pixel"

	flameobject "github.com/isangeles/flame/core/module/object"
	"github.com/isangeles/flame/core/module/object/character"
	"github.com/isangeles/flame/core/module/object/item"

	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/object/internal"
)

// Avatar struct for graphical representation of
// game character.
type Avatar struct {
	*character.Character
	data     *res.AvatarData
	sprite   *internal.AvatarSprite
	portrait *pixel.Sprite
	items    map[string]*ItemGraphic
	eqItems  map[string]*ItemGraphic
	effects  map[string]*EffectGraphic
	skills   map[string]*SkillGraphic
}

// NewAvatar creates new avatar for specified game character
// from specified avatar resources.
func NewAvatar(char *character.Character, data *res.AvatarData) *Avatar {
	av := new(Avatar)
	av.Character = char
	av.data = data
	// Sprite.
	if data.SSFullBodyPic != nil {
		av.sprite = internal.NewFullBodyAvatarSprite(data.SSFullBodyPic)
	}
	if data.SSTorsoPic != nil && data.SSHeadPic != nil {
		av.sprite = internal.NewAvatarSprite(data.SSTorsoPic, data.SSHeadPic)
	}
	// Portrait.
	av.portrait = pixel.NewSprite(data.PortraitPic, data.PortraitPic.Bounds())
	// Items, effects, skills.
	av.items = make(map[string]*ItemGraphic, 0)
	av.eqItems = make(map[string]*ItemGraphic, 0)
	av.effects = make(map[string]*EffectGraphic, 0)
	av.skills = make(map[string]*SkillGraphic, 0)
	av.updateGraphic()
	return av
}

// Draw draws avatar.
func (av *Avatar) Draw(win *mtk.Window, matrix pixel.Matrix) {
	av.sprite.Draw(win, matrix)
}

// Update updates avatar.
func (av *Avatar) Update(win *mtk.Window) {
	if av.InMove() {
		av.sprite.Move()
		pos := av.Position()
		dest := av.DestPoint()
		switch {
		case pos.X < dest.X:
			av.sprite.Right()
		case pos.Y < dest.Y:
			av.sprite.Up()
		case pos.X > dest.X:
			av.sprite.Left()
		case pos.Y > dest.Y:
			av.sprite.Down()
		}
	} else {
		av.sprite.Idle()
	}
	av.updateGraphic()
	av.sprite.Update(win)
}

// DrawArea returns current draw area.
func (av *Avatar) DrawArea() pixel.Rect {
	return av.sprite.DrawArea()
}

// Portrait returns avatar portrait.
func (av *Avatar) Portrait() *pixel.Sprite {
	return av.portrait
}

// Position return current position of avatar.
func (av *Avatar) Position() pixel.Vec {
	x, y := av.Character.Position()
	return pixel.V(x, y)
}

// DestPoint returns current destination point of
// avatar.
func (av *Avatar) DestPoint() pixel.Vec {
	x, y := av.Character.DestPoint()
	return pixel.V(x, y)
}

// Items returns all avatar items(in form of
// graphical wrappers).
func (av *Avatar) Items() (items []*ItemGraphic) {
	for _, ig := range av.items {
		items = append(items, ig)
	}
	return
}

// Effects returns all visible effects active on
// avatar character.
func (av *Avatar) Effects() (effects []*EffectGraphic) {
	for _, eg := range av.effects {
		effects = append(effects, eg)
	}
	return
}

// Skills retruns all avatar skills(in form of
// graphical wrappers).
func (av *Avatar) Skills() (skills []*SkillGraphic) {
	for _, sg := range av.skills {
		skills = append(skills, sg)
	}
	return
}

// Data returns avatar graphical data.
func (av *Avatar) Data() *res.AvatarData {
	return av.data
}

// equip equips specified graphical item.
func (av *Avatar) equip(gItem *ItemGraphic) error {
	switch gItem.Item.(type) {
	case *item.Weapon:
		av.sprite.SetWeapon(gItem.Spritesheet())
		av.eqItems[gItem.SerialID()] = gItem
		return nil
	default:
		return fmt.Errorf("not_equipable_item_type")
	}
}

// unequip removes graphic of specified item from
// avatar(if equiped).
func (av *Avatar) unequip(gItem *ItemGraphic) {
	switch gItem.Item.(type) {
	case *item.Weapon:
		av.sprite.SetWeapon(nil)
		delete(av.eqItems, gItem.SerialID())
	}
}

// updateApperance updates avatar grapphical
// content.
func (av *Avatar) updateGraphic() {
	// Clear items.
	for id, ig := range av.items {
		for _, it := range av.Inventory().Items() {
			if flameobject.Equals(it, ig) {
				continue
			}
			delete(av.items, id)
		}
	}
	// Clear unequipped items.
	for _, ig := range av.eqItems {
		eit, ok := ig.Item.(item.Equiper)
		if !ok {
			continue
		}
		if !av.Equipment().Equiped(eit) {
			av.unequip(ig)
		}
	}
	// Clear effects.
	for id, eg := range av.effects {
		for _, eff := range av.Character.Effects() {
			if flameobject.Equals(eg, eff) {
				continue
			}
			delete(av.effects, id)
		}
	}
	// Clear skills.
	for id, sg := range av.skills {
		for _, skill := range av.Character.Skills() {
			if flameobject.Equals(sg, skill) {
				continue
			}
			delete(av.skills, id)
		}
	}
	// Inventory.
	for _, it := range av.Inventory().Items() {
		if av.items[it.ID()+it.Serial()] != nil {
			continue
		}
		data := res.Item(it.ID())
		if data == nil {
			continue
		}
		itemGraphic := NewItemGraphic(it, data)
		av.items[it.ID()+it.Serial()] = itemGraphic
	}
	// Equipment.
	for _, eqi := range av.Equipment().Items() {
		it, ok := eqi.(item.Item)
		if !ok {
			continue
		}
		if av.eqItems[it.ID()+it.Serial()] != nil {
			continue
		}
		itemGData := res.Item(eqi.ID())
		if itemGData == nil {
			continue
		}
		itemGraphic := NewItemGraphic(it, itemGData)
		err := av.equip(itemGraphic)
		if err != nil {
			av.eqItems[it.ID()+it.Serial()] = itemGraphic
		}
	}
	// Effects.
	for _, e := range av.Character.Effects() {
		if av.effects[e.ID()+e.Serial()] != nil {
			continue
		}
		effectGData := res.Effect(e.ID())
		if effectGData == nil {
			continue
		}
		effectGraphic := NewEffectGraphic(e, effectGData)
		av.effects[e.ID()+e.Serial()] = effectGraphic
	}
	// Skills.
	for _, s := range av.Character.Skills() {
		if av.skills[s.ID()+s.Serial()] != nil {
			continue
		}
		data := res.Skill(s.ID())
		if data == nil {
			continue
		}
		skillGraphic := NewSkillGraphic(s, data)
		av.skills[s.ID()+s.Serial()] = skillGraphic
	}
}
