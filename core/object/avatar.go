/*
 * avatar.go
 *
 * Copyright 2018-2020 Dariusz Sikora <dev@isangeles.pl>
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

	flameres "github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/module/item"
	flameobject "github.com/isangeles/flame/module/objects"
	"github.com/isangeles/flame/module/character"
	"github.com/isangeles/flame/module/skill"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/object/internal"
	"github.com/isangeles/mural/log"
)

// Avatar struct for graphical representation of
// game character.
type Avatar struct {
	*character.Character
	data      *res.AvatarData
	sprite    *internal.AvatarSprite
	chat      *mtk.Text
	hovered   bool
	speaking  bool
	silenced  bool
	chatTimer int64
	items     map[string]*ItemGraphic
	eqItems   map[string]*ItemGraphic
	effects   map[string]*EffectGraphic
	skills    map[string]*SkillGraphic
}

// Type for avatar animations
// types.
type AvatarAnimType int

const (
	// Animation types.
	AvatarIdle AvatarAnimType = iota
	AvatarMove
	AvatarSpellCast
	AvatarCraftCast
	AvatarMelee
	AvatarShoot
	AvatarKneel
	AvatarLie
	// Chat popup visibility time.
	chatTimeMax = 2000
)

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
	chatParams := mtk.Params{
		FontSize: mtk.SizeSmall,
	}
	av.chat = mtk.NewText(chatParams)
	// Items, effects, skills.
	av.items = make(map[string]*ItemGraphic, 0)
	av.eqItems = make(map[string]*ItemGraphic, 0)
	av.effects = make(map[string]*EffectGraphic, 0)
	av.skills = make(map[string]*SkillGraphic, 0)
	// Events.
	av.SetOnSkillActivatedFunc(av.onSkillActivated)
	av.SetOnChatSentFunc(av.onChatSent)
	av.updateGraphic()
	return av
}

// Draw draws avatar.
func (av *Avatar) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Sprite.
	av.sprite.Draw(win, matrix)
	// Chat.
	chatPos := mtk.MoveTC(av.sprite.DrawArea().Size(), av.chat.Size())
	if av.speaking {
		av.chat.Draw(win, matrix.Moved(chatPos))
	}
}

// Update updates avatar.
func (av *Avatar) Update(win *mtk.Window) {
	// Animations
	switch {
	case av.castingSpell():
		av.sprite.SpellCast()
	case av.castingRecipe():
		av.sprite.CraftCast()
	case av.Moving():
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
	default:
		av.sprite.Idle()
	}
	// Sprite
	av.updateGraphic()
	av.sprite.Update(win)
	// Chat.
	if av.speaking {
		av.chatTimer += win.Delta()
		if av.chatTimer >= chatTimeMax {
			av.speaking = false
			av.chatTimer = 0
		}
	}
	av.hovered = av.sprite.DrawArea().Contains(win.MousePosition())
}

// DrawArea returns current draw area.
func (av *Avatar) DrawArea() pixel.Rect {
	return av.sprite.DrawArea()
}

// Portrait returns avatar portrait
// picture.
func (av *Avatar) Portrait() pixel.Picture {
	return av.data.PortraitPic
}

// Position returns current position of avatar.
func (av *Avatar) Position() pixel.Vec {
	x, y := av.Character.Position()
	return pixel.V(x, y)
}

// SetPosition sets current position of avatar.
func (av *Avatar) SetPosition(p pixel.Vec) {
	av.Character.SetPosition(p.X, p.Y)
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
	av.data.ID = av.ID()
	av.data.Serial = av.Serial()
	return av.data
}

// Silenced checks if audio effects are silenced.
func (av *Avatar) Silenced() bool {
	return av.silenced
}

// Silence toggles avatar audio effects.
func (av *Avatar) Silence(silence bool) {
	av.silenced = silence
}

// Hovered check if avatar is hovered
// by HUD user mouse cursor.
func (av *Avatar) Hovered() bool {
	return av.hovered
}

// equip adds graphic of specified item to avatar.
func (av *Avatar) equip(gItem *ItemGraphic) {
	switch gItem.Item.(type) {
	case *item.Weapon:
		av.sprite.SetWeapon(av.spritesheet(gItem.Spritesheets()))
		av.eqItems[gItem.ID()+gItem.Serial()] = gItem
	case *item.Armor:
		av.sprite.SetTorso(av.spritesheet(gItem.Spritesheets()))
		av.eqItems[gItem.ID()+gItem.Serial()] = gItem
	default:
		log.Dbg.Printf("avatar: %s#%s: equip: not equipable item type",
			av.ID(), av.Serial())
	}
}

// unequip removes graphic of specified item from
// avatar(if equiped).
func (av *Avatar) unequip(gItem *ItemGraphic) {
	switch gItem.Item.(type) {
	case *item.Weapon:
		av.sprite.SetWeapon(nil)
		delete(av.eqItems, gItem.ID()+gItem.Serial())
	case *item.Armor:
		av.sprite.SetTorso(nil)
		delete(av.eqItems, gItem.ID()+gItem.Serial())
	default:
		log.Dbg.Printf("avatar: %s#%s: equip: not equipable item type",
			av.ID(), av.Serial())
	}
}

// updateGraphic updates avatar grapphical
// content.
func (av *Avatar) updateGraphic() {
	// Clear items.
	for id, ig := range av.items {
		found := false
		for _, it := range av.Inventory().Items() {
			found = flameobject.Equals(it, ig)
		}
		if !found {
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
		found := false
		for _, eff := range av.Character.Effects() {
			found = flameobject.Equals(eg, eff)
		}
		if !found {
			delete(av.effects, id)
		}
	}
	// Clear skills.
	for id, sg := range av.skills {
		found := false
		for _, skill := range av.Character.Skills() {
			found = flameobject.Equals(sg, skill)
		}
		if !found {
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
		av.equip(itemGraphic)
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

// infoText returns info text about
// specified avatar.
func (av *Avatar) infoText() string {
	form := "%s"
	info := fmt.Sprintf(form, av.Name())
	if config.Debug {
		info = fmt.Sprintf("%s\n[%s_%s]", info, av.ID(), av.Serial())
	}
	return info
}

// Triggered after one of character skills was activated.
func (av *Avatar) onSkillActivated(s *skill.Skill) {
	sg := av.skills[s.ID()+s.Serial()]
	if sg == nil {
		log.Err.Printf("avatar: %s_%s: on skill activated: fail to find skill graphic: %s_%s",
			av.ID(), av.Serial(), s.ID(), s.Serial())
		return
	}
	// Animation.
	switch sg.ActivationAnim() {
	case AvatarMelee:
		av.sprite.MeleeOnce()
	}
	// Audio effect.
	if !av.Silenced() && mtk.Audio() != nil && sg.ActivationAudio() != nil {
		mtk.Audio().Play(sg.ActivationAudio())
	}
}

// Triggered after sending text to character
// chat channel.
func (av *Avatar) onChatSent(t string) {
	av.chat.SetText(t)
	av.speaking = true
}

// castingRecipe checks if avatar crafting
// any items right now.
func (av *Avatar) castingRecipe() bool {
	for _, r := range av.Crafting().Recipes() {
		if r.Casting() {
			return true
		}
	}
	return false
}

// castingSpell check if avatar casting
// any skills right now.
func (av *Avatar) castingSpell() bool {
	for _, s := range av.Skills() {
		if s.Casting() {
			return true
		}
	}
	return false
}

// spritesheet selects proper spritesheet for avatar from
// specified slice and returns its texture.
func (av *Avatar) spritesheet(sprs []*res.SpritesheetData) pixel.Picture {
	for _, s := range sprs {
		if s.Race != "*" {
			race := flameres.Race(s.Race)
			if race != nil && av.Race().ID() != race.ID {
				continue
			}
		}
		if s.Gender != "*" {
			gender := character.Gender(s.Gender)
			if av.Gender() != gender {
				continue
			}
		}
		return s.Texture
	}
	return nil
}
