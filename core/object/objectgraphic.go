/*
 * object.go
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

package object

import (
	"fmt"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"

	"github.com/isangeles/flame/core/module/objects"
	flameobject "github.com/isangeles/flame/core/module/object"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/data/res"
)

// Struct for graphical representation
// of area object.
type ObjectGraphic struct {
	*flameobject.Object
	data     *res.ObjectGraphicData
	sprite   *mtk.Animation
	info     *mtk.InfoWindow
	hovered  bool
	silenced bool
	effects  map[string]*EffectGraphic
	items    map[string]*ItemGraphic
}

// NewObject creates new graphical wrapper for specified object.
func NewObjectGraphic(ob *flameobject.Object, data *res.ObjectGraphicData) *ObjectGraphic {
	og := new(ObjectGraphic)
	og.Object = ob
	og.data = data
	// Sprite.
	og.sprite = mtk.NewAnimation(buildSpriteFrames(data.SpritePic), 2)
	// Info window.
	infoParams := mtk.Params{
		FontSize:  mtk.SizeSmall,
		MainColor: colornames.Grey,
	}
	og.info = mtk.NewInfoWindow(infoParams)
	og.info.SetText(og.infoText())
	// Effect, items.
	og.effects = make(map[string]*EffectGraphic)
	og.items = make(map[string]*ItemGraphic)
	return og
}

// Draw draws object sprite.
func (og *ObjectGraphic) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Sprite.
	og.sprite.Draw(t, matrix)
	// Info window.
	if og.hovered {
		og.info.Draw(t)
	}
}

// Update updates object.
func (og *ObjectGraphic) Update(win *mtk.Window) {
	// Sprite.
	og.sprite.Update(win)
	// Info window.
	og.info.Update(win)
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
	return og.data.PortraitPic
}

// Position return object position in form of
// pixel XY vector.
func (og *ObjectGraphic) Position() pixel.Vec {
	x, y := og.Object.Position()
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
		items = append(items, ig)
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
			continue
		}
		itemGraphic := NewItemGraphic(it, data)
		og.items[it.ID()+it.Serial()] = itemGraphic
	}
}

// infoText returns text for object info
// window.
func (og *ObjectGraphic) infoText() string {
	form := "%s"
	info := fmt.Sprintf(form, og.Name())
	if config.Debug() {
		info = fmt.Sprintf("%s\n[%s_%s]", info, og.ID(), og.Serial())
	}
	return info
}

// buildSpriteFrames creates animation frames from specified
// spritesheet.
func buildSpriteFrames(ss pixel.Picture) []*pixel.Sprite {
	frames := []*pixel.Sprite{
		pixel.NewSprite(ss, ss.Bounds()),
	}
	return frames
}
