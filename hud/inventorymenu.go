/*
 * inventorymenu.go
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
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"

	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/item"

	"github.com/isangeles/fire/request"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/data/res/graphic"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/object"
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
	bg := graphic.Textures["invbg.png"]
	if bg != nil {
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
	closeButtonBG := graphic.Textures["closebutton1.png"]
	if closeButtonBG != nil {
		closeButtonSpr := pixel.NewSprite(closeButtonBG, closeButtonBG.Bounds())
		im.closeButton.SetBackground(closeButtonSpr)
	} else {
		log.Err.Printf("hud: inventory menu: unable to retrieve background texture")
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
	upButtonBG := graphic.Textures["scrollup.png"]
	if upButtonBG != nil {
		upBG := pixel.NewSprite(upButtonBG, upButtonBG.Bounds())
		im.slots.SetUpButtonBackground(upBG)
	} else {
		log.Err.Printf("hud: inventory menu: unable to retrieve slot list up button texture")
	}
	downButtonBG := graphic.Textures["scrolldown.png"]
	if downButtonBG != nil {
		downBG := pixel.NewSprite(downButtonBG, downButtonBG.Bounds())
		im.slots.SetDownButtonBackground(downBG)
	} else {
		log.Err.Printf("hud: inventory menu: unable to retrieve slot list down button texture")
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
		mtk.DrawRect(win, im.DrawArea(), mainColor)
	}
	// Title.
	titleTextPos := pixel.V(0, im.Size().Y/2-mtk.ConvSize(25))
	im.titleText.Draw(win, matrix.Moved(titleTextPos))
	// Buttons.
	closeButtonPos := pixel.V(im.Size().X/2-mtk.ConvSize(20),
		im.Size().Y/2-mtk.ConvSize(15))
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
		if im.Opened() {
			im.Hide()
		} else {
			im.Show()
		}
	}
	if win.JustPressed(exitKey) {
		im.Hide()
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

// Show shows menu.
func (im *InventoryMenu) Show() {
	im.opened = true
	im.hud.UserFocus().Focus(im)
	im.refresh()
}

// Hide hides menu.
func (im *InventoryMenu) Hide() {
	im.opened = false
	im.hud.UserFocus().Focus(nil)
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
		return mtk.ConvVec(pixel.V(0, 0))
	}
	return mtk.ConvVec(im.bgSpr.Frame().Size())
}

// insertItems inserts specified items in inventory slots.
func (im *InventoryMenu) insertItems(items ...*item.InventoryItem) {
	im.slots.Clear()
	// Insert items from layout first.
	pc := im.hud.Game().ActivePlayerChar()
	layout := im.hud.Layout(pc.ID(), pc.Serial())
	for _, i := range items {
		it := itemGraphic(i.Item)
		slotID := layout.InvSlotID(it)
		if slotID < 0 {
			continue
		}
		if slotID < len(im.slots.Slots())-1 {
			slot := im.slots.Slots()[slotID]
			im.hud.insertSlotItem(it, slot)
		}
	}
	// Insert new items.
	for _, i := range items {
		it := itemGraphic(i.Item)
		// Skip items from layout.
		slotID := layout.InvSlotID(it)
		if slotID > -1 {
			continue
		}
		// Find proper slot.
		slot := im.slots.EmptySlot()
		// Try to find slot with same content and available space.
		for _, s := range im.slots.Slots() {
			if len(s.Values()) < 1 || len(s.Values()) >= it.MaxStack() {
				continue
			}
			slotIt, ok := s.Values()[0].(*object.ItemGraphic)
			if !ok {
				continue
			}
			if slotIt.ID() == it.ID() {
				slot = s
				break
			}
		}
		if slot == nil {
			log.Err.Printf("hud: inventory menu: no empty slots")
			return
		}
		// Insert item to slot.
		im.hud.insertSlotItem(it, slot)
	}
}

// updateLayout updates inventory layout for active player.
func (im *InventoryMenu) updateLayout() {
	// Retrieve layout for current PC.
	pc := im.hud.Game().ActivePlayerChar()
	layout := im.hud.layouts[pc.ID()+pc.Serial()]
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
				log.Err.Printf("hud: inventory menu: update layout: unable to retrieve slot value")
				continue
			}
			layout.SaveInvSlot(it, i)
		}
	}
	im.hud.layouts[pc.ID()+pc.Serial()] = layout
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

// removeSlotItems removes all items from specified slot and
// from PC inventory.
func (im *InventoryMenu) removeSlotItems(s *mtk.Slot) {
	items := make(map[string][]string, 0)
	for _, v := range s.Values() {
		it, ok := v.(item.Item)
		if !ok {
			continue
		}
		items[it.ID()] = append(items[it.ID()], it.Serial())
		im.hud.Game().ActivePlayerChar().Inventory().RemoveItem(it)
	}
	s.Clear()
	// Server request to remove items.
	if im.hud.Game().Server() != nil {
		pc := im.hud.Game().ActivePlayerChar()
		throwItemsReq := request.ThrowItems{
			ObjectID:     pc.ID(),
			ObjectSerial: pc.Serial(),
			Items:        items,
		}
		req := request.Request{ThrowItems: []request.ThrowItems{throwItemsReq}}
		err := im.hud.Game().Server().Send(req)
		if err != nil {
			log.Err.Printf("HUD: Inventory Menu: unable to send throw items request: %v",
				err)
		}
	}
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
	dlg.SetAcceptLabel(lang.Text("accept_button_label"))
	dlg.SetCancelLabel(lang.Text("cancel_button_label"))
	rmFunc := func(mw *mtk.MessageWindow) {
		im.removeSlotItems(s)
	}
	dlg.SetOnAcceptFunc(rmFunc)
	im.hud.ShowMessage(dlg)
}

// refresh inserts player items to inventory
// slots and saves inventory layout.
func (im *InventoryMenu) refresh() {
	pcAvatar := im.hud.PCAvatar()
	if pcAvatar != nil {
		im.insertItems(pcAvatar.Inventory().Items()...)
	}
	im.updateLayout()
}

// Triggered after close button clicked.
func (im *InventoryMenu) onCloseButtonClicked(b *mtk.Button) {
	im.Hide()
}

// Triggered after one of item slots was clicked with
// right mosue button.
func (im *InventoryMenu) onSlotRightClicked(s *mtk.Slot) {
	if len(s.Values()) < 1 {
		return
	}
	it, ok := s.Values()[0].(*object.ItemGraphic)
	if !ok {
		log.Err.Printf("inventory: inavlid slot value: %v", s.Values()[0])
		return
	}
	switch it := it.Item.(type) {
	case item.Equiper:
		if im.hud.Game().ActivePlayerChar().Equipment().Equiped(it) {
			im.hud.Game().ActivePlayerChar().Unequip(it)
			s.SetColor(invSlotColor)
			break
		}
		err := im.hud.Game().ActivePlayerChar().Equip(it)
		if err != nil {
			log.Err.Printf("inventory: item: %s %s: unable to equip: %v", it.ID(),
				it.Serial(), err)
			return
		}
		s.SetColor(invSlotEqColor)
	case *item.Misc:
		im.hud.Game().ActivePlayerChar().Use(it)
		if it.Consumable() {
			im.hud.Game().ActivePlayerChar().Inventory().RemoveItem(it)
			im.refresh()
		}
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
