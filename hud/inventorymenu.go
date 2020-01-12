/*
 * inventorymenu.go
 *
 * Copyright 2019-2020 Dariusz Sikora <dev@isangeles.pl>
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
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame/core/data/res/lang"
	"github.com/isangeles/flame/core/module/item"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

var (
	invKey         = pixelgl.KeyB
	invSlots       = 90
	invSlotSize    = mtk.SizeBig
	invSlotColor   = pixel.RGBA{0.1, 0.1, 0.1, 0.5}
	invSlotEqColor = pixel.RGBA{0.3, 0.3, 0.3, 0.5}
	invSpecialKey  = pixelgl.KeyLeftShift
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
	titleParams := mtk.Params{
		FontSize: mtk.SizeSmall,
	}
	im.titleText = mtk.NewText(titleParams)
	im.titleText.SetText(lang.Text("hud_inv_title"))
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		Shape:     mtk.ShapeSquare,
		MainColor: accentColor,
	}
	im.closeButton = mtk.NewButton(buttonParams)
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
		invSlotColor, invSlotSize)
	// Create empty slots.
	for i := 0; i < invSlots; i++ {
		s := im.createSlot()
		im.slots.Add(s)
	}
	// Slot list scroll buttons.
	upButtonBG, err := data.PictureUI("scrollup.png")
	if err == nil {
		upBG := pixel.NewSprite(upButtonBG, upButtonBG.Bounds())
		im.slots.SetUpButtonBackground(upBG)
	} else {
		log.Err.Printf("hud_inv:fail_to_retrieve_slot_list_up_button_texture:%v",
			err)
	}
	downButtonBG, err := data.PictureUI("scrolldown.png")
	if err == nil {
		downBG := pixel.NewSprite(downButtonBG, downButtonBG.Bounds())
		im.slots.SetDownButtonBackground(downBG)
	} else {
		log.Err.Printf("hud_inv:fail_to_retrieve_slot_list_down_button_texture:%v",
			err)
	}
	return im
}

// Draw draws menu.
func (im *InventoryMenu) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	im.drawArea = mtk.MatrixToDrawArea(matrix, im.Size())
	// Background.
	if im.bgSpr != nil {
		im.bgSpr.Draw(win, matrix)
	} else {
		mtk.DrawRectangle(win, im.DrawArea(), mainColor)
	}
	// Title.
	titleTextPos := mtk.ConvVec(pixel.V(0, im.Size().Y/2-25))
	im.titleText.Draw(win, matrix.Moved(titleTextPos))
	// Buttons.
	closeButtonPos := mtk.ConvVec(pixel.V(im.Size().X/2-20,
		im.Size().Y/2-15))
	im.closeButton.Draw(win, matrix.Moved(closeButtonPos))
	// Slots.
	im.slots.Draw(win, matrix)
}

// Update updates menu.
func (im *InventoryMenu) Update(win *mtk.Window) {
	// Ket events.
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		dragSlot := im.draggedSlot()
		if dragSlot != nil && !im.DrawArea().Contains(win.MousePosition()) {
			im.confirmRemove(dragSlot)
		}
	}
	if !im.hud.Chat().Activated() && win.JustPressed(invKey) {
		im.Show(!im.Opened())
	}
	// Elements update.
	if im.Opened() {
		im.slots.Update(win)
		im.closeButton.Update(win)
	}
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
		im.insertItems(im.hud.ActivePlayer().Inventory().Items()...)
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

// insertItems inserts specified items in inventory slots.
func (im *InventoryMenu) insertItems(items ...item.Item) {
	im.slots.Clear()
	for _, it := range items {
		// Retrieve item graphic.
		igd := res.Item(it.ID())
		if igd == nil { // if icon was found
			log.Err.Printf("hud_inv:item_graphic_not_found:%s", it.ID())
			// Get error icon.
			errData, err := data.ErrorItemGraphic()
			if err != nil {
				log.Err.Printf("hud_inv:fail_to_retrieve_error_graphic:%v", err)
				continue
			}
			errData.ItemID = it.ID()
			igd = errData
		}
		ig := object.NewItemGraphic(it, igd)
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
				if len(s.Values()) < 1 || len(s.Values()) >= ig.MaxStack() {
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
		im.hud.insertSlotItem(ig, slot)
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
	params := mtk.Params{
		Size:      invSlotSize,
		FontSize:  mtk.SizeMini,
		MainColor: invSlotColor,
	}
	s := mtk.NewSlot(params)
	s.SetSpecialKey(invSpecialKey)
	s.SetOnRightClickFunc(im.onSlotRightClicked)
	s.SetOnLeftClickFunc(im.onSlotLeftClicked)
	s.SetOnSpecialLeftClickFunc(im.onSlotSpecialLeftClicked)
	return s
}

// draggedSlot returns currently dragged inventory
// slot or nil.
func (im *InventoryMenu) draggedSlot() *mtk.Slot {
	for _, s := range im.slots.Slots() {
		if s.Dragged() {
			return s
		}
	}
	return nil
}

// removeSlotItem removes item from specified slot and
// from PC inventory.
func (im *InventoryMenu) removeSlotItem(s *mtk.Slot) {
	for _, v := range s.Values() {
		it, ok := v.(item.Item)
		if !ok {
			continue
		}
		im.hud.ActivePlayer().Inventory().RemoveItem(it)
	}
	s.Clear()
}

// confirmRemove shows warning message and removes content
// from inventory slot after warning dialog accepted.
func (im *InventoryMenu) confirmRemove(s *mtk.Slot) {
	s.Drag(false)
	msg := lang.Text("hud_inv_remove_item_warn")
	dlgParams := mtk.Params{
		Size:      mtk.SizeMedium,
		FontSize:  mtk.SizeMedium,
		MainColor: mainColor,
		SecColor:  accentColor,
		Info:      msg,
	}
	dlg := mtk.NewDialogWindow(dlgParams)
	dlg.SetAcceptLabel(lang.Text("accept_b_label"))
	dlg.SetCancelLabel(lang.Text("cancel_b_label"))
	rmFunc := func(mw *mtk.MessageWindow) {
		im.removeSlotItem(s)
	}
	dlg.SetOnAcceptFunc(rmFunc)
	im.hud.ShowMessage(dlg)
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
		log.Err.Printf("hud_inv_menu:not_item:%v", s.Values()[0])
		return
	}
	eit, ok := it.Item.(item.Equiper)
	if !ok {
		log.Err.Printf("hud_inv_menu:not_equipable_item:%s", it.ID())
		return
	}
	if im.hud.ActivePlayer().Equipment().Equiped(eit) {
		im.hud.ActivePlayer().Equipment().Unequip(eit)
		s.SetColor(invSlotColor)
	} else {
		err := im.hud.ActivePlayer().Equipment().Equip(eit)
		if err != nil {
			log.Err.Printf("hud_inv_menu:item:%s_%s:fail_to_equip:%v", eit.ID(),
				eit.Serial(), err)
			return
		}
		s.SetColor(invSlotEqColor)
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
				im.hud.insertSlotItem(ig, s)
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
