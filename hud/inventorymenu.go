/*
 * inventorymenu.go
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

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/object/item"
	
	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// Struct for inventory menu.
type InventoryMenu struct {
	hud         *HUD
	bgSpr       *pixel.Sprite
	bgDraw      *imdraw.IMDraw
	drawArea    pixel.Rect
	titleText   *mtk.Text
	closeButton *mtk.Button
	slots       *mtk.SlotList
	opened      bool
	focused     bool
}

var (
	inv_slots         = 90
	inv_slot_size     = mtk.SIZE_BIG
	inv_slot_color    = pixel.RGBA{0.1, 0.1, 0.1, 0.5}
	inv_slot_eq_color = pixel.RGBA{0.3, 0.3, 0.3, 0.5}
	inv_special_key   = pixelgl.KeyLeftShift
)

// newInventoryMenu creates new inventory menu for HUD.
func newInventoryMenu(hud *HUD) *InventoryMenu {
	im := new(InventoryMenu)
	im.hud = hud
	// Background.
	im.bgDraw = imdraw.New(nil)
	bg, err := data.PictureUI("invbg.png")
	if err == nil {
		im.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	}
	// Title.
	im.titleText = mtk.NewText(mtk.SIZE_SMALL, 0)
	im.titleText.SetText(lang.Text("gui", "hud_inv_title"))
	// Buttons.
	im.closeButton = mtk.NewButton(mtk.SIZE_SMALL, mtk.SHAPE_SQUARE, accent_color)
	closeButtonBG, err := data.PictureUI("closebutton1.png")
	if err == nil {
		closeButtonSpr := pixel.NewSprite(closeButtonBG, closeButtonBG.Bounds())
		im.closeButton.SetBackground(closeButtonSpr)
	} else {
		log.Err.Printf("hud_inventory:fail_to_retrieve_background_tex:%v", err)
	}
	im.closeButton.SetOnClickFunc(im.onCloseButtonClicked)
	// Slots list.
	im.slots = mtk.NewSlotList(mtk.ConvVec(pixel.V(250, 300)),
		inv_slot_color, inv_slot_size)
	upButtonBG, err := data.PictureUI("scrollup.png")
	if err == nil {
		upBG := pixel.NewSprite(upButtonBG, upButtonBG.Bounds())
		im.slots.SetUpButtonBackground(upBG)
	} else {
		log.Err.Printf("hud_inv:fail_to_retrieve_slot_list_up_buttons_texture:%v",
			err)
	}
	downButtonBG, err := data.PictureUI("scrolldown.png")
	if err == nil {
		downBG := pixel.NewSprite(downButtonBG, downButtonBG.Bounds())
		im.slots.SetDownButtonBackground(downBG)
	} else {
		log.Err.Printf("hud_inv:fail_to_retrieve_slot_list_down_buttons_texture:%v",
			err)
	}
	// Create empty slots.
	for i := 0; i < inv_slots; i++ {
		s := im.createSlot()
		im.slots.Add(s)
	}
	return im
}

// Draw draws menu.
func (im *InventoryMenu) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	im.drawArea = mtk.MatrixToDrawArea(matrix, im.Size())
	// Background.
	if im.bgSpr != nil {
		im.bgSpr.Draw(win.Window, matrix)
	} else {
		mtk.DrawRectangle(win.Window, im.DrawArea(), nil)
	}
	// Title.
	titleTextPos := mtk.ConvVec(pixel.V(0, im.Size().Y/2-25))
	im.titleText.Draw(win.Window, matrix.Moved(titleTextPos))
	// Buttons.
	closeButtonPos := mtk.ConvVec(pixel.V(im.Size().X/2-20,
		im.Size().Y/2-15))
	im.closeButton.Draw(win.Window, matrix.Moved(closeButtonPos))
	// Slots.
	im.slots.Draw(win, matrix)
}

// Update updates menu.
func (im *InventoryMenu) Update(win *mtk.Window) {
	// Elements update.
	im.slots.Update(win)
	im.closeButton.Update(win)
}

// Opened checks whether menu is open.
func (im *InventoryMenu) Opened() bool {
	return im.opened
}

// Show toggles menu visibility.
func (im *InventoryMenu) Show(show bool) {
	im.opened = show
	if im.Opened() {
		im.hud.UserFocus().Focus(im)
		im.slots.Clear()
		im.insert(im.hud.ActivePlayer().Items()...)
		im.updateLayout()
	} else {
		im.hud.UserFocus().Focus(nil)
	}
}

// Focused checks whether menu is focused.
func (im *InventoryMenu) Focused() bool {
	return im.focused
}

// Focus toggles menu focus.
func (im *InventoryMenu) Focus(focus bool) {
	im.focused = focus
}

// DrawArea returns current menu draw area.
func (im *InventoryMenu) DrawArea() pixel.Rect {
	return im.drawArea
}

// Size returns size of menu background.
func (im *InventoryMenu) Size() pixel.Vec {
	if im.bgSpr == nil {
		// TODO: size for draw background.
		return pixel.V(mtk.ConvSize(0), mtk.ConvSize(0))
	}
	return im.bgSpr.Frame().Size()
}

// insert inserts specified items in inventory slots.
func (im *InventoryMenu) insert(items ...*object.ItemGraphic) {
	for _, it := range items {
		// Find proper slot.
		slot := im.slots.EmptySlot()
		layout := im.hud.Layout(im.hud.ActivePlayer().ID(), im.hud.ActivePlayer().Serial())
		slotID := layout.InvSlotID(it)
		if slotID > -1 { // insert item to slot from layout
			if slotID < len(im.slots.Slots())-1 {
				slot = im.slots.Slots()[slotID]
			}
		} else { // try to find slot with same content and available space
			for _, s := range im.slots.Slots() {
				if len(s.Values()) < 1 || len(s.Values()) >= it.MaxStack() {
					continue
				}
				slotIt, ok := s.Values()[0].(item.Item)
				if !ok {
					continue
				}
				if slotIt.ID() == it.ID() {
					slot = s
					break
				}
			}
		}
		if slot == nil {
			log.Err.Printf("hud_inv:no empty slots")
			return
		}
		// Insert item to slot.
		insertSlotItem(it, slot)
	}
}

// updateLayout updates inventory layout for active player.
func (im *InventoryMenu) updateLayout() {
	// Retrieve layout for current PC.
	layout := im.hud.layouts[im.hud.ActivePlayer().SerialID()]
	if layout == nil {
		layout = NewLayout()
	}
	// Clear layout.
	layout.SetInvSlots(make(map[string]int))
	// Set Layout.
	for i, s := range im.slots.Slots() {
		if len(s.Values()) < 1 {
			continue
		}
		for _, v := range s.Values() {
			it, ok := v.(*object.ItemGraphic)
			if !ok {
				log.Err.Printf("hud_inv:update_layout:fail to retrieve slot value")
				continue
			}
			layout.SaveInvSlot(it, i)
		}
	}
	im.hud.layouts[im.hud.ActivePlayer().SerialID()] = layout
}

// createSlot creates empty slot for inventory slots list.
func (im *InventoryMenu) createSlot() *mtk.Slot {
	s := mtk.NewSlot(inv_slot_size, mtk.SIZE_MINI)
	s.SetColor(inv_slot_color)
	s.SetSpecialKey(inv_special_key)
	s.SetOnRightClickFunc(im.onSlotRightClicked)
	s.SetOnLeftClickFunc(im.onSlotLeftClicked)
	s.SetOnSpecialLeftClickFunc(im.onSlotSpecialLeftClicked)
	return s
}

// draggedItems returns currently dragged slot
// with items.
func (im *InventoryMenu) draggedItems() *mtk.Slot {
	for _, s := range im.slots.Slots() {
		if s.Dragged() {
			return s
		}
	}
	return nil
}

// Triggered after close button clicked.
func (im *InventoryMenu) onCloseButtonClicked(b *mtk.Button) {
	im.Show(false)
}

// Triggered after one of item slots was clicked with
// right mosue button.
func (im *InventoryMenu) onSlotRightClicked(s *mtk.Slot) {
	if len(s.Values()) < 1 {
		return
	}
	it, ok := s.Values()[0].(*object.ItemGraphic)
	if !ok {
		log.Err.Printf("hud_inv_menu:%v:is not item", s.Values()[0])
		return
	}
	eit, ok := it.Item.(item.Equiper)
	if !ok {
		log.Err.Printf("hud_inv_menu:%s:is not equipable item", it.SerialID())
		return
	}
	if im.hud.ActivePlayer().Equipment().Equiped(eit) {
		err := im.hud.ActivePlayer().Equipment().Unequip(eit)
		if err != nil {
			log.Err.Printf("hud_inv_menu:%s:fail_to_unequip:%v", eit.SerialID(), err)
			return
		}
		s.SetColor(inv_slot_color)
	} else {
		err := im.hud.ActivePlayer().Equipment().Equip(eit)
		if err != nil {
			log.Err.Printf("hud_inv_menu:%s:fail_to_equip:%v", eit.SerialID(), err)
			return
		}
		s.SetColor(inv_slot_eq_color)
	}
}

// Triggered after one of item slots was clicked with
// laft mouse button.
func (im *InventoryMenu) onSlotLeftClicked(s *mtk.Slot) {
	for _, ds := range im.slots.Slots() {
		if !ds.Dragged() {
			continue
		}
		if s == ds {
			ds.Drag(false)
			return
		}
		mtk.SlotSwitch(s, ds)
		im.updateLayout()
		ds.Drag(false)
		return
	}
	if len(s.Values()) < 1 {
		return
	}
	s.Drag(true)
}

// Triggered after one of items slots was clicked with
// left mouse button and inv_slot_special_key pressed.
func (im *InventoryMenu) onSlotSpecialLeftClicked(s *mtk.Slot) {
	// Handle dragged slot.
	for _, ds := range im.slots.Slots() {
		if !ds.Dragged() {
			continue
		}
		if s == ds {
			ds.Drag(false)
			return
		}
		if len(s.Values()) < 1 {
			if dv := ds.Pop(); dv != nil {
				ig, ok := dv.(*object.ItemGraphic)
				if !ok {
					ds.AddValues(dv) // return value back to dragged slot
					return
				}
				insertSlotItem(ig, s)
			}
		} else {
			v, ok := s.Values()[0].(*object.ItemGraphic)
			if !ok {
				return
			}
			dv, ok := ds.Pop().(*object.ItemGraphic)
			if !ok {
				ds.AddValues(dv) // return value back to dragged slot
				return
			}
			if v.ID() != dv.ID() ||
				len(s.Values()) >= v.MaxStack() {
				return
			}
			s.AddValues(dv)
		}
		ds.Drag(false)
		im.updateLayout()
		return
	}
	if len(s.Values()) < 1 {
		return
	}
	s.Drag(true)
}

// insertSlotItem inserts specified item to specified slot.
func insertSlotItem(it *object.ItemGraphic, s *mtk.Slot) {
	s.AddValues(it)
	s.SetInfo(itemInfo(it.Item))
	s.SetIcon(it.Icon())
}

// itemInfo returns formated string with
// informations about specified item.
func itemInfo(it item.Item) string {
	switch i := it.(type) {
	case *item.Weapon:
		infoForm := "%s\n%d-%d"
		dmgMin, dmgMax := i.Damage()
		info := fmt.Sprintf(infoForm, i.Name(),
			dmgMin, dmgMax)
		if config.Debug() { // add serial ID info
			info = fmt.Sprintf("%s\n[%s_%s]", info,
				i.ID(), i.Serial())
		}
		return info
	default:
		return ""
	}
}
