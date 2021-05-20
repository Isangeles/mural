/*
 * avatar.go
 *
 * Copyright 2018-2021 Dariusz Sikora <dev@isangeles.pl>
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
	"time"

	"github.com/faiface/pixel"

	flameres "github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/character"
	"github.com/isangeles/flame/craft"
	"github.com/isangeles/flame/effect"
	"github.com/isangeles/flame/item"
	"github.com/isangeles/flame/objects"
	"github.com/isangeles/flame/skill"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/data/res/graphic"
	"github.com/isangeles/mural/object/internal"
	"github.com/isangeles/mural/log"
)

// Avatar struct for graphical representation of
// game character.
type Avatar struct {
	*character.Character
	name         string
	portrait     pixel.Picture
	sprite       *internal.AvatarSprite
	chat         *mtk.Text
	hovered      bool
	speaking     bool
	silenced     bool
	chatTimer    int64
	items        map[string]*ItemGraphic
	eqItems      map[string]*ItemGraphic
	effects      map[string]*EffectGraphic
	skills       map[string]*SkillGraphic
	portraitName string
	torsoName    string
	headName     string
	fullBodyName string
	combatLog    *objects.Log
	privateLog   *objects.Log
}

// Type for avatar animations
// types.
type AvatarAnimType string

const (
	// Animation types.
	AvatarIdle      AvatarAnimType = AvatarAnimType("idle")
	AvatarMove                     = AvatarAnimType("move")
	AvatarSpellCast                = AvatarAnimType("spell")
	AvatarCraftCast                = AvatarAnimType("craft")
	AvatarMelee                    = AvatarAnimType("melee")
	AvatarShoot                    = AvatarAnimType("shoot")
	AvatarKneel                    = AvatarAnimType("kneel")
	AvatarLie                      = AvatarAnimType("lie")
	// Chat popup visibility time.
	chatTimeMax = 2000
)

// NewAvatar creates new avatar for specified game character
// from specified avatar resources.
func NewAvatar(char *character.Character, data *res.AvatarData) *Avatar {
	av := new(Avatar)
	av.Character = char
	av.SetName(data.Name)
	// Translate name if empty.
	if len(av.Name()) < 1 {
		av.SetName(lang.Text(av.ID()))
	}
	// Portrait.
	av.portrait = graphic.Portraits[data.Portrait]
	if av.portrait != nil {
		av.portraitName = data.Portrait
	}
	// Sprite.
	fullBodyPic := graphic.AvatarSpritesheets[data.FullBody]
	if fullBodyPic != nil {
		av.sprite = internal.NewFullBodyAvatarSprite(fullBodyPic)
		av.fullBodyName = data.FullBody
	} else {
		torsoPic := graphic.AvatarSpritesheets[data.Torso]
		headPic := graphic.AvatarSpritesheets[data.Head]
		if torsoPic != nil && headPic != nil {
			av.sprite = internal.NewAvatarSprite(torsoPic, headPic)
			av.torsoName = data.Torso
			av.headName = data.Head
		} else {
			log.Err.Printf("avatar: %s#%s: sprite textures not found: fullbody: '%s' torso: '%s' head: '%s'",
				av.ID(), av.Serial(), data.FullBody, data.Torso, data.Head)
		}
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
	// Logs.
	av.combatLog = objects.NewLog()
	av.privateLog = objects.NewLog()
	// Events.
	av.SetOnSkillActivatedFunc(av.onSkillActivated)
	av.SetOnModifierTakenFunc(av.onModifierTaken)
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
	for _, m := range av.ChatLog().Messages() {
		duration := time.Since(m.Time())
		av.speaking = duration.Seconds() < 2
		if av.speaking {
			text := m.String()
			if !m.Translated {
				text = lang.Text(m.String())
			}
			av.chat.SetText(text)
			break
		}
	}
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

// Name returns avatar name.
func (av *Avatar) Name() string {
	return av.name
}

// SetName sets avatar name.
func (av *Avatar) SetName(name string) {
	av.name = name
}

// Portrait returns avatar portrait
// picture.
func (av *Avatar) Portrait() pixel.Picture {
	return av.portrait
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
		if av.Inventory().Item(ig.ID(), ig.Serial()) != nil {
			items = append(items, ig)
		}
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
func (av *Avatar) Data() res.AvatarData {
	data := res.AvatarData{
		ID:       av.ID(),
		Serial:   av.Serial(),
		Name:     av.Name(),
		Portrait: av.portraitName,
		Torso:    av.torsoName,
		Head:     av.headName,
		FullBody: av.fullBodyName,
	}
	return data
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

// CombatLog returns avatar combat log.
func (av *Avatar) CombatLog() *objects.Log {
	return av.combatLog
}

// PrivateLog returns avatar private log.
func (av *Avatar) PrivateLog() *objects.Log {
	return av.privateLog
}

// updateGraphic updates avatar grapphical
// content.
func (av *Avatar) updateGraphic() {
	// Clear items.
	for id, ig := range av.items {
		found := false
		for _, it := range av.Inventory().Items() {
			found = objects.Equals(it, ig)
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
			av.removeItemGraphic(ig)
			delete(av.eqItems, ig.ID()+ig.Serial())
		}
	}
	// Clear effects.
	for id, eg := range av.effects {
		found := false
		for _, eff := range av.Character.Effects() {
			found = objects.Equals(eg, eff)
		}
		if !found {
			delete(av.effects, id)
		}
	}
	// Clear skills.
	for id, sg := range av.skills {
		found := false
		for _, skill := range av.Character.Skills() {
			found = sg.ID() == skill.ID()
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
			data = DefaultItemGraphic(it)
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
		av.addItemGraphic(itemGraphic)
		av.eqItems[itemGraphic.ID()+itemGraphic.Serial()] = itemGraphic
	}
	// Effects.
	for _, e := range av.Character.Effects() {
		if av.effects[e.ID()+e.Serial()] != nil {
			continue
		}
		effectGData := res.Effect(e.ID())
		if effectGData == nil {
			effectGData = DefaultEffectGraphic(e)
		}
		effectGraphic := NewEffectGraphic(e, effectGData)
		av.effects[e.ID()+e.Serial()] = effectGraphic
	}
	// Skills.
	for _, s := range av.Character.Skills() {
		if av.skills[s.ID()] != nil {
			continue
		}
		data := res.Skill(s.ID())
		if data == nil {
			continue
		}
		skillGraphic := NewSkillGraphic(s, data)
		av.skills[s.ID()] = skillGraphic
	}
}

// addItemGraphic adds item graphic to avatar.
func (av *Avatar) addItemGraphic(gItem *ItemGraphic) {
	sprite := av.spritesheet(gItem.Spritesheets())
	if sprite == nil {
		return
	}
	tex := graphic.AvatarSpritesheets[sprite.Texture]
	if tex == nil {
		log.Err.Printf("avatar: %s#%s: item texture not found: %s",
			av.ID(), av.Serial(), sprite.Texture)
		return
	}
	switch gItem.Item.(type) {
	case *item.Weapon:
		av.sprite.SetWeapon(tex)
	case *item.Armor:
		av.sprite.SetTorso(tex)
	}
}

// removeItemGraphic removes item graphic from
// avatar.
func (av *Avatar) removeItemGraphic(gItem *ItemGraphic) {
	switch gItem.Item.(type) {
	case *item.Weapon:
		av.sprite.SetWeapon(nil)
	case *item.Armor:
		av.sprite.SetTorso(nil)
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
	sg := av.skills[s.ID()]
	if sg == nil {
		log.Err.Printf("avatar: %s %s: on skill activated: skill graphic not found: %s",
			av.ID(), av.Serial(), s.ID())
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

// Triggered on receiving new modifier.
func (av *Avatar) onModifierTaken(m effect.Modifier) {
	switch m := m.(type) {
	case *effect.HealthMod:
		msg := objects.Message{
			Translated: true,
			Text: fmt.Sprintf("%s: %d", lang.Text("ob_health"),
				m.LastValue()),
		}
		av.CombatLog().Add(msg)
	}
}

// castingRecipe checks if avatar crafting
// any items right now.
func (av *Avatar) castingRecipe() bool {
	_, ok := av.Casted().(*craft.Recipe)
	return ok
}

// castingSpell check if avatar casting
// any skills right now.
func (av *Avatar) castingSpell() bool {
	_, ok := av.Casted().(*skill.Skill)
	return ok
}

// spritesheet selects proper spritesheet for avatar from
// specified slice and returns its texture.
func (av *Avatar) spritesheet(sprs []*res.SpritesheetData) *res.SpritesheetData {
	for _, s := range sprs {
		if s.Race != "*" {
			race := flameres.Race(s.Race)
			if race == nil || av.Race().ID() != race.ID {
				continue
			}
		}
		if s.Gender != "*" {
			gender := character.Gender(s.Gender)
			if av.Gender() != gender {
				continue
			}
		}
		return s
	}
	return nil
}
