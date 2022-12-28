/*
 * object.go
 *
 * Copyright 2019-2022 Dariusz Sikora <ds@isangeles.dev>
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

	"github.com/isangeles/flame/character"
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/effect"
	"github.com/isangeles/flame/objects"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/data/res/graphic"
	"github.com/isangeles/mural/log"
)

// Struct for graphical representation
// of area object.
type ObjectGraphic struct {
	*character.Character
	name         string
	sprite       *mtk.Animation
	portrait     pixel.Picture
	spriteName   string
	portraitName string
	hovered      bool
	silenced     bool
	effects      map[string]*EffectGraphic
	items        map[string]*ItemGraphic
	combatLog    *objects.Log
}

// NewObjectGraphic creates new object graphic for specified character.
func NewObjectGraphic(char *character.Character, data *res.ObjectGraphicData) *ObjectGraphic {
	og := new(ObjectGraphic)
	og.Character = char 
	og.name = lang.Text(og.ID())
	og.combatLog = objects.NewLog()
	// Sprite.
	spritePic := graphic.ObjectSpritesheets[data.Sprite]
	if spritePic != nil {
		og.sprite = mtk.NewAnimation(buildSpriteFrames(spritePic), 2)
		og.spriteName = data.Sprite
	} else {
		log.Err.Printf("object: %s#%s: sprite texture not found: %s", og.ID(),
			og.Serial(), data.Sprite)
	}
	// Portrait.
	og.portrait = graphic.Portraits[data.Portrait]
	if og.portrait != nil {
		og.portraitName = data.Portrait
	}
	// Effect, items.
	og.effects = make(map[string]*EffectGraphic)
	og.items = make(map[string]*ItemGraphic)
	// Events.
	og.SetOnModifierTakenFunc(og.onModifierTaken)
	return og
}

// Draw draws object sprite.
func (og *ObjectGraphic) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Sprite.
	og.sprite.Draw(t, matrix)
}

// Update updates object.
func (og *ObjectGraphic) Update(win *mtk.Window) {
	// Sprite.
	og.sprite.Update(win)
	og.hovered = og.sprite.DrawArea().Contains(win.MousePosition())
	// Graphic.
	og.updateGraphic()
}

// DrawArea returns current draw area of
// object sprite.
func (og *ObjectGraphic) DrawArea() pixel.Rect {
	return og.sprite.DrawArea()
}

// Portrait returns portrait picture.
func (og *ObjectGraphic) Portrait() pixel.Picture {
	return og.portrait
}

// Name returns object name.
func (og *ObjectGraphic) Name() string {
	return og.name
}

// Position return object position in form of
// pixel XY vector.
func (og *ObjectGraphic) Position() pixel.Vec {
	x, y := og.Character.Position()
	return pixel.V(x, y)
}

// Effects returns all object effects in form of
// graphical wrappers.
func (og *ObjectGraphic) Effects() []*EffectGraphic {
	effs := make([]*EffectGraphic, 0)
	for _, e := range og.effects {
		effs = append(effs, e)
	}
	return effs
}

// Items returns all object items(in form of
// graphical wrappers).
func (og *ObjectGraphic) Items() (items []*ItemGraphic) {
	for _, ig := range og.items {
		if og.Inventory().Item(ig.ID(), ig.Serial()) != nil {
			items = append(items, ig)
		}
	}
	return
}

// LootItems returns all 'lootable' items(in form of
// graphical wrappers).
func (og *ObjectGraphic) LootItems() (items []*ItemGraphic) {
	for _, ig := range og.items {
		if og.Inventory().LootItem(ig.ID(), ig.Serial()) != nil {
			items = append(items, ig)
		}
	}
	return
}

// MaxMana returns 0, object do not have mana.
// Function to satify frame target interface.
func (og *ObjectGraphic) MaxMana() int {
	return 0
}

// Silenced checks if audio effects are silenced.
func (og *ObjectGraphic) Silenced() bool {
	return og.silenced
}

// Silence toggles object audio effects.
func (og *ObjectGraphic) Silence(silence bool) {
	og.silenced = silence
}

// Hovered checks if object is hovered by
// HUD user mouse cursor.
func (og *ObjectGraphic) Hovered() bool {
	return og.hovered
}

// ComabatLog retruns object comabt log.
func (og *ObjectGraphic) CombatLog() *objects.Log {
	return og.combatLog
}

// updateGraphic updates object
// graphical content.
func (og *ObjectGraphic) updateGraphic() {
	// Clear items.
	for id, ig := range og.items {
		found := false
		for _, it := range og.Inventory().Items() {
			found = objects.Equals(it, ig)
		}
		if !found {
			delete(og.items, id)
		}
	}

	// Inventory.
	for _, it := range og.Inventory().Items() {
		if og.items[it.ID()+it.Serial()] != nil {
			continue
		}
		data := res.Item(it.ID())
		if data == nil {
			data = DefaultItemGraphic(it)
		}
		itemGraphic := NewItemGraphic(it, data)
		og.items[it.ID()+it.Serial()] = itemGraphic
	}
}

// Triggered after taking new modifier.
func (og *ObjectGraphic) onModifierTaken(m effect.Modifier) {
	switch m := m.(type) {
	case *effect.HealthMod:
		msg := objects.Message{
			Translated: true,
			Text: fmt.Sprintf("%s: %d", lang.Text("ob_health"),
				m.LastValue()),
		}
		og.CombatLog().Add(msg)
	}
}

// buildSpriteFrames creates animation frames from specified
// spritesheet.
func buildSpriteFrames(ss pixel.Picture) []*pixel.Sprite {
	frames := []*pixel.Sprite{
		pixel.NewSprite(ss, ss.Bounds()),
	}
	return frames
}
