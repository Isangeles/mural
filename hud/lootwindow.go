/*
 * lootwindow.go
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
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"

	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/object/item"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// Struct for HUD loot window.
type LootWindow struct {
	hud         *HUD
	bgSpr       *pixel.Sprite
	bgDraw      *imdraw.IMDraw
	drawArea    pixel.Rect
	titleText   *mtk.Text
	closeButton *mtk.Button
	slots       *mtk.SlotList
	opened      bool
	focused     bool
	target      LootTarget
}

// Interface for 'lootable' objects.
type LootTarget interface {
	Inventory() *item.Inventory
	Items() []*object.ItemGraphic
}

var (
	lootSlots     = 90
	lootSlotSize  = mtk.SizeBig
	lootSlotColor = pixel.RGBA{0.1, 0.1, 0.1, 0.5}
)

// newLootWindow creates new loot window for HUD.
func newLootWindow(hud *HUD) *LootWindow {
	lw := new(LootWindow)
	lw.hud = hud
	// Background.
	lw.bgDraw = imdraw.New(nil)
	bg, err := data.PictureUI("invbg.png")
	if err == nil {
		lw.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	}
	// Title.
	titleParams := mtk.Params{
		FontSize: mtk.SizeSmall,
	}
	lw.titleText = mtk.NewText(titleParams)
	lw.titleText.SetText(lang.Text("gui", "hud_loot_title"))
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		Shape:     mtk.ShapeSquare,
		MainColor: accentColor,
	}
	lw.closeButton = mtk.NewButton(buttonParams)
	closeButtonBG, err := data.PictureUI("closebutton1.png")
	if err == nil {
		spr := pixel.NewSprite(closeButtonBG, closeButtonBG.Bounds())
		lw.closeButton.SetBackground(spr)
	}
	lw.closeButton.SetOnClickFunc(lw.onCloseButtonClicked)
	// Slots list.
	lw.slots = mtk.NewSlotList(mtk.ConvVec(pixel.V(250, 300)), lootSlotColor,
		lootSlotSize)
	for i := 0; i < lootSlots; i++ {
		s := lw.createSlot()
		lw.slots.Add(s)
	}
	// Slot list scroll buttons.
	upButtonBG, err := data.PictureUI("scrollup.png")
	if err == nil {
		spr := pixel.NewSprite(upButtonBG, upButtonBG.Bounds())
		lw.slots.SetUpButtonBackground(spr)
	}
	downButtonBG, err := data.PictureUI("scrolldown.png")
	if err == nil {
		spr := pixel.NewSprite(downButtonBG, downButtonBG.Bounds())
		lw.slots.SetDownButtonBackground(spr)
	}
	return lw
}

// Draw draws window.
func (lw *LootWindow) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	lw.drawArea = mtk.MatrixToDrawArea(matrix, lw.Size())
	// Background.
	if lw.bgSpr != nil {
		lw.bgSpr.Draw(win.Window, matrix)
	} else {
		mtk.DrawRectangle(win.Window, lw.DrawArea(), nil)
	}
	// Title.
	titleTextPos := mtk.ConvVec(pixel.V(0, lw.Size().Y/2-25))
	lw.titleText.Draw(win.Window, matrix.Moved(titleTextPos))
	// Buttons.
	closeButtonPos := mtk.ConvVec(pixel.V(lw.Size().X/2-20,
		lw.Size().Y/2-15))
	lw.closeButton.Draw(win.Window, matrix.Moved(closeButtonPos))
	// Slots.
	lw.slots.Draw(win, matrix)
}

// Update updates window.
func (lw *LootWindow) Update(win *mtk.Window) {
	// Elements.
	lw.slots.Update(win)
	lw.closeButton.Update(win)
}

// Show toggles window visibility.
func (lw *LootWindow) Show(show bool) {
	lw.opened = show
}

// Opened checks whether window is open.
func (lw *LootWindow) Opened() bool {
	return lw.opened
}

// Focus toggles window focus.
func (lw *LootWindow) Focus(f bool) {
	lw.focused = f
}

// Focused checks whether window is focused.
func (lw *LootWindow) Focused() bool {
	return lw.focused
}

// Size returns size of loot window background.
func (lw *LootWindow) Size() pixel.Vec {
	if lw.bgSpr == nil {
		// TODO: size for draw background.
		return pixel.V(mtk.ConvSize(0), mtk.ConvSize(0))
	}
	return lw.bgSpr.Frame().Size()
}

// DrawArea returns current draw area of window
// background.
func (lw *LootWindow) DrawArea() pixel.Rect {
	return lw.drawArea
}

// SetTarget sets object with inventory to loot.
func (lw *LootWindow) SetTarget(t LootTarget) {
	lw.target = t
	lw.slots.Clear()
	for _, it := range lw.target.Inventory().Items() {
		if !it.Loot() {
			continue
		}
		lw.insertItems(it)
	}
}

// insert inserts specified items in window slots.
func (lw *LootWindow) insert(items ...*object.ItemGraphic) {
	for _, it := range items {
		slot := lw.slots.EmptySlot()
		if slot == nil {
			slot = lw.createSlot()
			lw.slots.Add(slot)
		}
		lw.hud.insertSlotItem(it, slot)
	}
}

// insertItems inserts specified items in window slots.
func (lw *LootWindow) insertItems(items ...item.Item) {
	for _, it := range items {
		// Retrieve item graphic.
		data := res.Item(it.ID())
		if data == nil {
			log.Err.Printf("hud_loot:item_graphic_not_found:%s\n", it.ID())
			continue
		}
		ig := object.NewItemGraphic(it, data)
		// Find empty slot.
		slot := lw.slots.EmptySlot()
		if slot == nil {
			slot = lw.createSlot()
			lw.slots.Add(slot)
		}
		lw.hud.insertSlotItem(ig, slot)
	}
}

// createSlot creates empty slot for loot slots list.
func (lw *LootWindow) createSlot() *mtk.Slot {
	params := mtk.Params{
		Size:      lootSlotSize,
		FontSize:  mtk.SizeMini,
		MainColor: lootSlotColor,
	}
	s := mtk.NewSlot(params)
	s.SetOnLeftClickFunc(lw.onSlotLeftClicked)
	return s
}

// Triggered after close button was clicked.
func (lw *LootWindow) onCloseButtonClicked(b *mtk.Button) {
	lw.Show(false)
}

// Triggered after one of items slots was clicked with
// left mouse button.
func (lw *LootWindow) onSlotLeftClicked(s *mtk.Slot) {
	if len(s.Values()) < 1 {
		return
	}
	valuesLen := len(s.Values())
	for i := 0; i < valuesLen; i++ {
		v := s.Pop()
		ig, ok := v.(*object.ItemGraphic)
		if !ok {
			continue
		}
		err := lw.hud.ActivePlayer().Inventory().AddItem(ig.Item)
		if err != nil {
			log.Err.Printf("hud_loot:fail_to_transfer_item:%v", err)
			continue
		}
		lw.target.Inventory().RemoveItem(ig.Item)
	}
}
