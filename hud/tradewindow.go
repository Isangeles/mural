/*
 * tradewindow.go
 *
 * Copyright 2019-2021 Dariusz Sikora <dev@isangeles.pl>
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
	"github.com/isangeles/flame/module/objects"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/data/res/graphic"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// Struct for HUD trade window.
type TradeWindow struct {
	hud         *HUD
	bgSpr       *pixel.Sprite
	bgDraw      *imdraw.IMDraw
	drawArea    pixel.Rect
	titleText   *mtk.Text
	valueText   *mtk.Text
	closeButton *mtk.Button
	tradeButton *mtk.Button
	buySlots    *mtk.SlotList
	sellSlots   *mtk.SlotList
	seller      item.Container
	sellItems   map[string]item.Item
	buyItems    map[string]item.Item
	opened      bool
	focused     bool
}

var (
	tradeBuySlots        = 90
	tradeSellSlots       = 90
	tradeSlotSize        = mtk.SizeMedium
	tradeSlotColor       = pixel.RGBA{0.1, 0.1, 0.1, 0.5}
	tradeSelectSlotColor = pixel.RGBA{0.3, 0.3, 0.3, 0.5}
	tradeSpecialKey      = pixelgl.KeyLeftShift
)

// newTradeWindow creates new trade
// window for HUD.
func newTradeWindow(hud *HUD) *TradeWindow {
	tw := new(TradeWindow)
	tw.hud = hud
	tw.sellItems = make(map[string]item.Item)
	tw.buyItems = make(map[string]item.Item)
	// Background.
	tw.bgDraw = imdraw.New(nil)
	bg := graphic.Textures["menubg.png"]
	if bg != nil {
		tw.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	} else {
		log.Err.Printf("hud trade: unable to retrieve background texure")
	}
	// Title.
	titleParams := mtk.Params{
		FontSize: mtk.SizeSmall,
	}
	tw.titleText = mtk.NewText(titleParams)
	tw.titleText.SetText(lang.Text("hud_trade_title"))
	// Trade value text.
	valueTextParams := mtk.Params{
		FontSize: mtk.SizeMini,
	}
	tw.valueText = mtk.NewText(valueTextParams)
	// Close button.
	closeButtonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		Shape:     mtk.ShapeSquare,
		MainColor: accentColor,
	}
	tw.closeButton = mtk.NewButton(closeButtonParams)
	closeButtonBG := graphic.Textures["closebutton1.png"]
	if closeButtonBG != nil {
		closeBG := pixel.NewSprite(closeButtonBG,
			closeButtonBG.Bounds())
		tw.closeButton.SetBackground(closeBG)
	} else {
		log.Err.Printf("hud trade: unable to retrieve close button texture")
	}
	tw.closeButton.SetOnClickFunc(tw.onCloseButtonClicked)
	// Trade button.
	tradeButtonParams := mtk.Params{
		Size:      mtk.SizeMini,
		FontSize:  mtk.SizeMini,
		Shape:     mtk.ShapeRectangle,
		MainColor: accentColor,
	}
	tw.tradeButton = mtk.NewButton(tradeButtonParams)
	tradeButtonBG := graphic.Textures["button_green.png"]
	if tradeButtonBG != nil {
		bg := pixel.NewSprite(tradeButtonBG, tradeButtonBG.Bounds())
		tw.tradeButton.SetBackground(bg)
	} else {
		log.Err.Printf("hud trade: unable to retrieve trade button texture")
	}
	tw.tradeButton.SetOnClickFunc(tw.onTradeButtonClicked)
	tw.tradeButton.SetLabel(lang.Text("hud_trade_accept"))
	// Buy slot list.
	tw.buySlots = mtk.NewSlotList(mtk.ConvVec(pixel.V(250, 150)),
		tradeSlotColor, tradeSlotSize)
	for i := 0; i < tradeBuySlots; i++ {
		s := tw.createBuySlot()
		tw.buySlots.Add(s)
	}
	// Sell slot list.
	tw.sellSlots = mtk.NewSlotList(mtk.ConvVec(pixel.V(250, 150)),
		tradeSlotColor, tradeSlotSize)
	for i := 0; i < tradeSellSlots; i++ {
		s := tw.createSellSlot()
		tw.sellSlots.Add(s)
	}
	// Slot lists scroll buttons.
	upButtonBG := graphic.Textures["scrollup.png"]
	if upButtonBG != nil {
		upBG := pixel.NewSprite(upButtonBG, upButtonBG.Bounds())
		tw.buySlots.SetUpButtonBackground(upBG)
		tw.sellSlots.SetUpButtonBackground(upBG)
	} else {
		log.Err.Printf("hud trade: unable to retrieve slot list up button texture")
	}
	downButtonBG := graphic.Textures["scrolldown.png"]
	if downButtonBG != nil {
		downBG := pixel.NewSprite(downButtonBG, downButtonBG.Bounds())
		tw.buySlots.SetDownButtonBackground(downBG)
		tw.sellSlots.SetDownButtonBackground(downBG)
	} else {
		log.Err.Printf("hud trade: unable to retrieve slot list down button texture")
	}
	tw.updateTradeValue()
	return tw
}

// Draw draws window.
func (tw *TradeWindow) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	tw.drawArea = mtk.MatrixToDrawArea(matrix, tw.Size())
	// Background.
	if tw.bgSpr != nil {
		tw.bgSpr.Draw(win, matrix)
	} else {
		mtk.DrawRectangle(win, tw.DrawArea(), mainColor)
	}
	// Title & trade value.
	titleTextMove := pixel.V(mtk.ConvSize(0), tw.Size().Y/2-mtk.ConvSize(25))
	valueTextMove := pixel.V(mtk.ConvSize(-80), -tw.Size().Y/2+mtk.ConvSize(30))
	tw.titleText.Draw(win, matrix.Moved(titleTextMove))
	tw.valueText.Draw(win, matrix.Moved(valueTextMove))
	// Buttons.
	closeButtonMove := pixel.V(tw.Size().X/2-mtk.ConvSize(20),
		tw.Size().Y/2-mtk.ConvSize(15))
	tradeButtonMove := pixel.V(mtk.ConvSize(50), -tw.Size().Y/2+mtk.ConvSize(30))
	tw.closeButton.Draw(win, matrix.Moved(closeButtonMove))
	tw.tradeButton.Draw(win, matrix.Moved(tradeButtonMove))
	// Slot lists.
	buySlotsMove := mtk.MoveTC(tw.Size(), tw.buySlots.Size())
	buySlotsMove.Y -= mtk.ConvSize(50)
	sellSlotsMove := mtk.MoveBC(tw.Size(), tw.sellSlots.Size())
	sellSlotsMove.Y += mtk.ConvSize(60)
	tw.buySlots.Draw(win, matrix.Moved(buySlotsMove))
	tw.sellSlots.Draw(win, matrix.Moved(sellSlotsMove))
}

// Update updates window.
func (tw *TradeWindow) Update(win *mtk.Window) {
	// Elements.
	if tw.Opened() {
		tw.closeButton.Update(win)
		tw.tradeButton.Update(win)
		tw.buySlots.Update(win)
		tw.sellSlots.Update(win)
	}
}

// Show shows window.
func (tw *TradeWindow) Show() {
	tw.opened = true
	if tw.seller != nil {
		tw.insertBuyItems(tw.seller.Inventory().TradeItems()...)
	}
	tw.insertSellItems(tw.hud.Game().ActivePlayer().Inventory().Items()...)
}

// Hide hides window.
func (tw *TradeWindow) Hide() {
	tw.opened = false
	tw.reset()
}

// Opened checks if window is open.
func (tw *TradeWindow) Opened() bool {
	return tw.opened
}

// DrawArea returns window draw area.
func (tw *TradeWindow) DrawArea() pixel.Rect {
	return tw.drawArea
}

// Size returns window background size.
func (tw *TradeWindow) Size() pixel.Vec {
	if tw.bgSpr == nil {
		return mtk.ConvVec(pixel.V(250, 350))
	}
	return mtk.ConvVec(tw.bgSpr.Frame().Size())
}

// SetSeller sets c as seller.
func (tw *TradeWindow) SetSeller(c item.Container) {
	tw.seller = c
}

// reset resets all window elements to
// default state.
func (tw *TradeWindow) reset() {
	tw.sellItems = make(map[string]item.Item)
	tw.buyItems = make(map[string]item.Item)
	tw.sellSlots.Clear()
	tw.buySlots.Clear()
	// Trade items highlight.
	for _, s := range tw.buySlots.Slots() {
		s.SetColor(tradeSlotColor)
	}
	for _, s := range tw.sellSlots.Slots() {
		s.SetColor(tradeSlotColor)
	}
}

// tradeValue returns current trade value.
func (tw *TradeWindow) tradeValue() (v int) {
	for _, it := range tw.buyItems {
		ti, ok := it.(*item.TradeItem)
		if !ok {
			log.Err.Printf("hud trade: item a trade item: %s#%s",
				it.ID(), it.Serial())
			continue
		}
		v -= ti.Price
	}
	for _, it := range tw.sellItems {
		v += it.Value()
	}
	return
}

// updateTradeValue updates trade value.
func (tw *TradeWindow) updateTradeValue() {
	// Trade value label.
	value := tw.tradeValue()
	label := lang.Text("hud_trade_value")
	tw.valueText.SetText(fmt.Sprintf("%s:%d", label, value))
}

// insertBuyItems inserts specified items in buy slots.
func (tw *TradeWindow) insertBuyItems(items ...*item.TradeItem) {
	tw.buySlots.Clear()
	for _, it := range items {
		// Retrieve item graphic.
		igd := res.Item(it.ID())
		if igd == nil { // if icon was found
			log.Err.Printf("hud trade: item graphic not found: %s", it.ID())
			// Get fallback graphic.
			igd = itemErrorGraphic(it)
		}
		ig := object.NewItemGraphic(it, igd)
		// Find proper slot.
		slot := tw.buySlots.EmptySlot()
		// Try to find slot with same content and available space.
		for _, s := range tw.buySlots.Slots() {
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
		if slot == nil {
			log.Err.Printf("hud trade: no empty buy slots")
			return
		}
		// Insert item to slot.
		tw.hud.insertSlotItem(ig, slot)
	}
}

// insertSellItem inserts specified items in sell slots.
func (tw *TradeWindow) insertSellItems(items ...item.Item) {
	tw.sellSlots.Clear()
	for _, it := range items {
		// Retrieve item graphic.
		igd := res.Item(it.ID())
		if igd == nil { // if icon was found
			log.Err.Printf("hud trade: item graphic not found: %s", it.ID())
			// Get fallback graphic.
			igd = itemErrorGraphic(it)
		}
		ig := object.NewItemGraphic(it, igd)
		// Find proper slot.
		slot := tw.sellSlots.EmptySlot()
		// Try to find slot with same content and available space.
		for _, s := range tw.sellSlots.Slots() {
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
		if slot == nil {
			log.Err.Printf("hud trade: no empty sell slots")
			return
		}
		// Insert item to slot.
		tw.hud.insertSlotItem(ig, slot)
	}
}

// createBuySlot creates slot for buy list.
func (tw *TradeWindow) createBuySlot() *mtk.Slot {
	params := mtk.Params{
		Size:      tradeSlotSize,
		FontSize:  mtk.SizeMini,
		MainColor: tradeSlotColor,
	}
	s := mtk.NewSlot(params)
	s.SetSpecialKey(tradeSpecialKey)
	s.SetOnRightClickFunc(tw.onBuySlotRightClicked)
	s.SetOnLeftClickFunc(tw.onBuySlotLeftClicked)
	s.SetOnSpecialLeftClickFunc(tw.onBuySlotSpecialLeftClicked)
	s.SetOnSpecialRightClickFunc(tw.onBuySlotSpecialRightClicked)
	return s
}

// createSellSlot creates slot for sell list.
func (tw *TradeWindow) createSellSlot() *mtk.Slot {
	params := mtk.Params{
		Size:      tradeSlotSize,
		FontSize:  mtk.SizeMini,
		MainColor: tradeSlotColor,
	}
	s := mtk.NewSlot(params)
	s.SetSpecialKey(tradeSpecialKey)
	s.SetOnRightClickFunc(tw.onSellSlotRightClicked)
	s.SetOnLeftClickFunc(tw.onSellSlotLeftClicked)
	s.SetOnSpecialLeftClickFunc(tw.onSellSlotSpecialLeftClicked)
	s.SetOnSpecialRightClickFunc(tw.onSellSlotSpecialRightClicked)
	return s
}

// Triggered after close button clicked.
func (tw *TradeWindow) onCloseButtonClicked(b *mtk.Button) {
	tw.Hide()
}

// triggered after trade button clicked.
func (tw *TradeWindow) onTradeButtonClicked(b *mtk.Button) {
	// Check trade value.
	if tw.tradeValue() < 0 {
		msg := objects.Message{Text: lang.Text("hud_trade_low_value_msg")}
		tw.hud.Game().ActivePlayer().PrivateLog().Add(msg)
		return
	}
	// Trade.
	for _, it := range tw.buyItems {
		tw.seller.Inventory().RemoveItem(it)
		tw.hud.Game().ActivePlayer().Inventory().AddItem(it)
	}
	for _, it := range tw.sellItems {
		tw.hud.Game().ActivePlayer().Inventory().RemoveItem(it)
		tw.seller.Inventory().AddItem(it)
	}
	tw.Hide()
}

// Triggered after one of buy slots was clicked
// with right mouse button.
func (tw *TradeWindow) onBuySlotRightClicked(s *mtk.Slot) {
	if len(s.Values()) < 1 {
		return
	}
	for _, v := range s.Values() {
		itg, ok := v.(*object.ItemGraphic)
		if !ok {
			log.Err.Printf("hud trade: invalid slot value: %v", v)
			return
		}
		delete(tw.buyItems, itg.ID()+itg.Serial())
	}
	s.SetColor(tradeSlotColor)
	tw.updateTradeValue()
}

// Triggered after one of buy slots was clicked
// with left mouse button.
func (tw *TradeWindow) onBuySlotLeftClicked(s *mtk.Slot) {
	if len(s.Values()) < 1 {
		return
	}
	for _, v := range s.Values() {
		itg, ok := v.(*object.ItemGraphic)
		if !ok {
			log.Err.Printf("hud trade: invalid slot value: %v", v)
			return
		}
		tw.buyItems[itg.ID()+itg.Serial()] = itg.Item
	}
	s.SetColor(tradeSelectSlotColor)
	tw.updateTradeValue()
}

// Triggered after one of buy slots was clicked with
// left mouse button and special key.
func (tw *TradeWindow) onBuySlotSpecialLeftClicked(s *mtk.Slot) {
	if len(s.Values()) < 1 {
		return
	}
	for _, v := range s.Values() {
		itg, ok := v.(*object.ItemGraphic)
		if !ok {
			log.Err.Printf("hud trade: invalid slot value: %v", v)
			return
		}
		if tw.buyItems[itg.ID()+itg.Serial()] == nil {
			tw.buyItems[itg.ID()+itg.Serial()] = itg.Item
			break
		}
	}
	s.SetColor(tradeSelectSlotColor)
	tw.updateTradeValue()
}

// Triggered after one of buy slots was clicked with
// right mouse button and special key.
func (tw *TradeWindow) onBuySlotSpecialRightClicked(s *mtk.Slot) {

}

// Triggered after one of sell slots was clicked
// with right mouse button.
func (tw *TradeWindow) onSellSlotRightClicked(s *mtk.Slot) {
	if len(s.Values()) < 1 {
		return
	}
	for _, v := range s.Values() {
		itg, ok := v.(*object.ItemGraphic)
		if !ok {
			log.Err.Printf("hud trade: invalid slot value: %v", v)
			return
		}
		delete(tw.sellItems, itg.ID()+itg.Serial())
	}
	s.SetColor(tradeSlotColor)
	tw.updateTradeValue()
}

// Triggered after one of sell slots was clicked
// with left mosue button.
func (tw *TradeWindow) onSellSlotLeftClicked(s *mtk.Slot) {
	if len(s.Values()) < 1 {
		return
	}
	for _, v := range s.Values() {
		itg, ok := v.(*object.ItemGraphic)
		if !ok {
			log.Err.Printf("hud trade: invalid slot value: %v", v)
			return
		}
		tw.sellItems[itg.ID()+itg.Serial()] = itg.Item
	}
	s.SetColor(tradeSelectSlotColor)
	tw.updateTradeValue()
}

// Triggered after one of sell slots was clicked with
// left mouse button and special key.
func (tw *TradeWindow) onSellSlotSpecialLeftClicked(s *mtk.Slot) {
	if len(s.Values()) < 1 {
		return
	}
	for _, v := range s.Values() {
		itg, ok := v.(*object.ItemGraphic)
		if !ok {
			log.Err.Printf("hud trade: invalid slot value: %v", v)
			return
		}
		if tw.sellItems[itg.ID()+itg.Serial()] == nil {
			tw.sellItems[itg.ID()+itg.Serial()] = itg.Item
			break
		}
	}
	s.SetColor(tradeSelectSlotColor)
	tw.updateTradeValue()
}

// Triggered after one of sell slots was clicked with
// right mouse button and special key.
func (tw *TradeWindow) onSellSlotSpecialRightClicked(s *mtk.Slot) {

}
