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

	"github.com/isangeles/flame/core/module/object/character"
	flameitem "github.com/isangeles/flame/core/module/object/item"

	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/object/internal"
	"github.com/isangeles/mural/log"
)

// Avatar struct for graphical representation of
// game character.
type Avatar struct {
	*character.Character

	sprite         *internal.AvatarSprite
	portrait       *pixel.Sprite
	portraitName   string
	ssHeadName     string
	ssTorsoName    string
	ssFullBodyName string
	eqItems        map[string]*ItemGraphic
}

// NewAvatar creates new avatar for specified game character.
// Portrait and spritesheet names are required for saving and
// loading avatar file.
func NewAvatar(data *res.AvatarData) *Avatar {
	av := new(Avatar)
	av.Character = data.Character
	// Portrait & spritesheets names.
	av.portraitName = data.PortraitName
	av.ssHeadName = data.SSHeadName
	av.ssTorsoName = data.SSTorsoName
	// Sprite.
	av.sprite = internal.NewAvatarSprite(data.SSTorsoPic, data.SSHeadPic)
	// Portrait.
	av.portrait = pixel.NewSprite(data.PortraitPic, data.PortraitPic.Bounds())
	// Items.
	av.eqItems = make(map[string]*ItemGraphic, 0)
	for _, eqItemData := range data.EqItemsGraphics {
		eqItem := NewItemGraphic(eqItemData)
		av.Equip(eqItem)
	}
	return av
}

// NewStaticAvatar creates new avatar with static(not affected by
// equipped items) body sprite.
// Portrait and spritesheet names are required for saving and
// loading avatar file.
func NewStaticAvatar(data *res.AvatarData) (*Avatar, error) {
	av := new(Avatar)
	av.Character = data.Character
	// Portrait & spritesheet names.
	av.portraitName = data.PortraitName
	av.ssFullBodyName = data.SSFullBodyName
	// Sprite.
	av.sprite = internal.NewFullBodyAvatarSprite(data.SSFullBodyPic)
	// Portrait.
	av.portrait = pixel.NewSprite(data.PortraitPic, data.PortraitPic.Bounds())
	// Items.
	av.eqItems = make(map[string]*ItemGraphic, 0)
	for _, eqItemData := range data.EqItemsGraphics {
		eqItem := NewItemGraphic(eqItemData)
		av.Equip(eqItem)
	}
	return av, nil
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
	av.updateApperance()
	av.sprite.Update(win)
}

// Portrait returns avatar portrait.
func (av *Avatar) Portrait() *pixel.Sprite {
	return av.portrait
}

// PortraitName returns name of portrait picture
// file.
func (av *Avatar) PortraitName() string {
	return av.portraitName
}

// TorsoSpritesheetName returns name of base torso
// spritesheet picture file.
func (av *Avatar) TorsoSpritesheetName() string {
	return av.ssTorsoName
}

// HeadSpritesheetName returns name of base head
// spritesheet picture file.
func (av *Avatar) HeadSpritesheetName() string {
	return av.ssHeadName
}

// FullBodySpritesheetName returns name of base
// full body spritesheet picture file.
func (av *Avatar) FullBodySpritesheetName() string {
	return av.ssFullBodyName
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

// Equip equips specified graphical item.
func (av *Avatar) Equip(gItem *ItemGraphic) error {
	switch it := gItem.Item.(type) {
	case *flameitem.Weapon:
		err := av.Equipment().EquipHandRight(it)
		if err != nil {
			return err
		}
		av.eqItems[it.SerialID()] = gItem
		av.sprite.SetWeapon(gItem.Spritesheet())
		av.eqItems[gItem.SerialID()] = gItem
		return nil
	default:
		return fmt.Errorf("unequipable_item_type")
	}
}

// updateApperance updates avatar sprite apperance.
func (av *Avatar) updateApperance() {
	// Update graphical items list.
	for _, eqi := range av.Equipment().Items() {
		if av.eqItems[eqi.SerialID()] != nil {
			continue
		}
		itemGData := res.ItemData(eqi.ID())
		if itemGData == nil {
			continue
		}
		itemGraphic := NewItemGraphic(itemGData)
		av.eqItems[eqi.SerialID()] = itemGraphic
	}
	// Equipped items.
	for _, eqItemGraphic := range av.eqItems {
		if av.Equipment().Equiped(eqItemGraphic.Item) {
			continue
		}
		err := av.Equip(eqItemGraphic)
		if err != nil {
			log.Err.Printf("new_avatar:%s:fail_to_equip_graphical_item:%v",
				av.SerialID(), err)
		}
	}
}
