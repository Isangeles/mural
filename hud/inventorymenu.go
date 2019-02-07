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
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"

	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/object/item"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/log"
)

const (
	inv_slots = 90
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
	im.titleText = mtk.NewText(lang.Text("gui", "hud_inv_title"), mtk.SIZE_SMALL, 0)
	// Buttons.
	closeButtonBG, err := data.PictureUI("closebutton1.png")
	if err != nil { // fallback
		log.Err.Printf("hud_inventory:fail_to_retrieve_background_tex:%v", err)
		im.closeButton = mtk.NewButton(mtk.SIZE_SMALL, mtk.SHAPE_SQUARE, accent_color,
			"", "")
	} else {
		im.closeButton = mtk.NewButtonSprite(closeButtonBG, mtk.SIZE_SMALL, "", "")
	}
	im.closeButton.SetOnClickFunc(im.onCloseButtonClicked)
	// Slots.
	slotsBGColor := pixel.RGBA{0.1, 0.1, 0.1, 0.5}
	im.slots = mtk.NewSlotList(mtk.ConvVec(pixel.V(250, 300)), slotsBGColor, mtk.SIZE_MINI)
	for i := 0; i < inv_slots; i ++ { // create empty slots
		s := mtk.NewSlot(mtk.SIZE_MINI, mtk.SIZE_MINI, mtk.SHAPE_SQUARE)
		im.slots.Add(s)
	}
	return im
}

// Draw draws menu.
func (im *InventoryMenu) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	im.drawArea = mtk.MatrixToDrawArea(matrix, im.Bounds())
	// Background.
	if im.bgSpr != nil {
		im.bgSpr.Draw(win.Window, matrix)
	} else {
		im.drawIMBackground(win.Window)
	}
	// Title.
	titleTextPos := mtk.ConvVec(pixel.V(0, im.Bounds().Max.Y/2 - 25))
	im.titleText.Draw(win.Window, matrix.Moved(titleTextPos))
	// Buttons.
	closeButtonPos := mtk.ConvVec(pixel.V(im.Bounds().Max.X/2 - 20,
		im.Bounds().Max.Y/2 - 15))
	im.closeButton.Draw(win.Window, matrix.Moved(closeButtonPos))
	// Slots.
	//firstSlotPos := mtk.ConvVec(pixel.V(-im.Bounds().W()/2 + 30, im.Bounds().H()/2 - 70))
	im.slots.Draw(win, matrix)
}

// Update updates menu.
func (im *InventoryMenu) Update(win *mtk.Window) {
	// Elements update.
	im.slots.Update(win)
	im.closeButton.Update(win)
}

// Opened checks wheter menu is open.
func (im *InventoryMenu) Opened() bool {
	return im.opened
}

// Show toggles menu visibility.
func (im *InventoryMenu) Show(show bool) {
	im.opened = show
	if im.Opened() {
		im.hud.UserFocus().Focus(im)
		im.insert(im.hud.ActivePlayer().Inventory().Items()...)
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

// Bounds returns size bounds of menu background.
func (im *InventoryMenu) Bounds() pixel.Rect {
	if im.bgSpr == nil {
		// TODO: bounds for draw background.
		return pixel.R(0, 0, mtk.ConvSize(0), mtk.ConvSize(0))
	}
	return im.bgSpr.Frame()
}

// drawIMBackground draw menu background with IMDraw.
func (im *InventoryMenu) drawIMBackground(t pixel.Target) {
	// TODO: draw background with IMDraw.
}

// insert inserts specified items in inventory slots.
func (im *InventoryMenu) insert(items ...item.Item) {
	for _, i := range items {
		slot := im.slots.Slots()[0]
		itemGraphic := res.ItemData(i.ID())
		if itemGraphic == nil {
			log.Err.Printf("hud_inv_menu:fail_to_find_item_graphic:%s",
				i.ID())
			slot.SetValue(i)
			continue
		}
		slot.SetValue(i)
		slot.SetIcon(itemGraphic.IconPic)
		slot.SetInfo(i.Name())
	}
}

// Triggered after close button clicked.
func (im *InventoryMenu) onCloseButtonClicked(b *mtk.Button) {
	im.Show(false)
}

// Triggered after one of item slots was clicked.
func (im *InventoryMenu) onSlotClicked(s *mtk.Slot) {

}
