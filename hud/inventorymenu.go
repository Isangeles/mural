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
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/module/item"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data/res/graphic"
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
		mtk.DrawRectangle(win, im.DrawArea(), mainColor)
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
		im.refresh()
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
		return mtk.ConvVec(pixel.V(0, 0))
	}
	return mtk.ConvVec(im.bgSpr.Frame().Size())
}

// insertItems inserts specified items in inventory slots.
func (im *InventoryMenu) insertItems(items ...*object.ItemGraphic) {
	im.slots.Clear()
	// Insert items from layout first.
	pc := im.hud.Game().ActivePlayer()
	layout := im.hud.Layout(pc.ID(), pc.Serial())
	for _, it := range items {
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
	for _, it := range items {
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
	layout := im.hud.layouts[im.hud.Game().ActivePlayer().SerialID()]
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
	im.hud.layouts[im.hud.Game().ActivePlayer().SerialID()] = layout
}

// equip inserts specified equipable item to all
// compatible slots in active PC equipment.
func (im *InventoryMenu) equip(it item.Equiper) error {
	pc := im.hud.Game().ActivePlayer()
	if !pc.MeetReqs(it.EquipReqs()...) {
		return fmt.Errorf("requirements not meet")
	}
	for _, itSlot := range it.Slots() {
		equiped := false
		for _, eqSlot := range pc.Equipment().Slots() {
			if eqSlot.Item() != nil {
				continue
			}
			if eqSlot.Type() == itSlot {
				eqSlot.SetItem(it)
				equiped = true
				break
			}
		}
		if !equiped {
			pc.Equipment().Unequip(it)
			return fmt.Errorf("free slot not found: %s", itSlot)
		}
	}
	if !pc.Equipment().Equiped(it) {
		return fmt.Errorf("no compatible slots")
	}
	return nil
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
		im.hud.Game().ActivePlayer().Inventory().RemoveItem(it)
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

// refresh inserts player items to inventory
// slots and saves inventory layout.
func (im *InventoryMenu) refresh() {
	im.insertItems(im.hud.Game().ActivePlayer().Items()...)
	im.updateLayout()
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
		log.Err.Printf("inventory: inavlid slot value: %v", s.Values()[0])
		return
	}
	switch it := it.Item.(type) {
	case item.Equiper:
		if im.hud.Game().ActivePlayer().Equipment().Equiped(it) {
			im.hud.Game().ActivePlayer().Equipment().Unequip(it)
			s.SetColor(invSlotColor)
			break
		}
		err := im.equip(it)
		if err != nil {
			log.Err.Printf("inventory: item: %s %s: unable to equip: %v", it.ID(),
				it.Serial(), err)
			return
		}
		s.SetColor(invSlotEqColor)
	case *item.Misc:
		im.hud.Game().ActivePlayer().Use(it)
		if it.Consumable() {
			im.hud.Game().ActivePlayer().Inventory().RemoveItem(it)
			im.refresh()
		}
	}
}

// TRIGGERED after one of item slots was clicked with
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
