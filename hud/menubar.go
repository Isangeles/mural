/*
 * menubar.go
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

// Struct for HUD menu bar.
type MenuBar struct {
	hud           *HUD
	bgSpr         *pixel.Sprite
	bgDraw        *imdraw.IMDraw
	drawArea      pixel.Rect
	menuButton    *mtk.Button
	invButton     *mtk.Button
	skillsButton  *mtk.Button
	journalButton *mtk.Button
	charButton    *mtk.Button
	slots         []*mtk.Slot
}

var (
	barSlots     = 10
	barSlotSize  = mtk.SizeMedium
	barSlotColor = pixel.RGBA{0.1, 0.1, 0.1, 0.5}
)

// newMenuBar creates new menu bar for HUD.
func newMenuBar(hud *HUD) *MenuBar {
	mb := new(MenuBar)
	mb.hud = hud
	// Background.
	mb.bgDraw = imdraw.New(nil)
	bg := graphic.Textures["menubar.png"]
	if bg != nil {
		mb.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	}
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		Shape:     mtk.ShapeSquare,
		MainColor: accentColor,
	}
	// Menu Button.
	mb.menuButton = mtk.NewButton(buttonParams)
	mb.menuButton.SetInfo(lang.Text("hud_bar_menu_open_info"))
	menuButtonBG := graphic.Textures["menubutton.png"]
	if menuButtonBG != nil {
		menuButtonSpr := pixel.NewSprite(menuButtonBG, menuButtonBG.Bounds())
		mb.menuButton.SetBackground(menuButtonSpr)
	} else {
		log.Err.Printf("hud: menu bar: unable to retrieve menu button texture")
	}
	mb.menuButton.SetOnClickFunc(mb.onMenuButtonClicked)
	// Inventory button.
	mb.invButton = mtk.NewButton(buttonParams)
	mb.invButton.SetInfo(lang.Text("hud_bar_inv_open_info"))
	invButtonBG := graphic.Textures["inventorybutton.png"]
	if invButtonBG != nil {
		invButtonSpr := pixel.NewSprite(invButtonBG, invButtonBG.Bounds())
		mb.invButton.SetBackground(invButtonSpr)
	} else {
		log.Err.Printf("hud: menu bar: unable to retrieve inv button texture")
	}
	mb.invButton.SetOnClickFunc(mb.onInvButtonClicked)
	// Skills button.
	mb.skillsButton = mtk.NewButton(buttonParams)
	mb.skillsButton.SetInfo(lang.Text("hud_bar_skills_open_info"))
	skillsButtonBG := graphic.Textures["skillsbutton.png"]
	if skillsButtonBG != nil {
		skillsButtonSpr := pixel.NewSprite(skillsButtonBG, skillsButtonBG.Bounds())
		mb.skillsButton.SetBackground(skillsButtonSpr)
	} else {
		log.Err.Printf("hud: menu bar: unable to retrieve skills button texture")
	}
	mb.skillsButton.SetOnClickFunc(mb.onSkillsButtonClicked)
	// Journal button.
	mb.journalButton = mtk.NewButton(buttonParams)
	journalInfo := lang.Text("hud_bar_journal_open_info")
	mb.journalButton.SetInfo(journalInfo)
	journalButtonBG := graphic.Textures["questsbutton.png"]
	if journalButtonBG != nil {
		journalButtonSpr := pixel.NewSprite(journalButtonBG, journalButtonBG.Bounds())
		mb.journalButton.SetBackground(journalButtonSpr)
	} else {
		log.Err.Printf("hud: menu bar: unable to retrieve quests button texture")
	}
	mb.journalButton.SetOnClickFunc(mb.onJournalButtonClicked)
	// Character button.
	mb.charButton = mtk.NewButton(buttonParams)
	charInfo := lang.Text("hud_bar_char_open_info")
	mb.charButton.SetInfo(charInfo)
	charButtonBG := graphic.Textures["charbutton.png"]
	if charButtonBG != nil {
		charButtonSpr := pixel.NewSprite(charButtonBG, charButtonBG.Bounds())
		mb.charButton.SetBackground(charButtonSpr)
	} else {
		log.Err.Printf("hud: menu bar: unable to retrieve char button texture")
	}
	mb.charButton.SetOnClickFunc(mb.onCharButtonClicked)
	// Slots.
	for i := 0; i < barSlots; i++ {
		s := mb.createSlot()
		s.SetLabel(fmt.Sprintf("%d", i+1))
		mb.slots = append(mb.slots, s)
	}
	return mb
}

// Draw draws menu bar.
func (mb *MenuBar) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	mb.drawArea = mtk.MatrixToDrawArea(matrix, mb.Size())
	// Background.
	if mb.bgSpr != nil {
		mb.bgSpr.Draw(win, matrix)
	} else {
		mb.drawIMBackground(win)
	}
	// Slots.
	slotsStartPos := mtk.ConvVec(pixel.V(-163, 0))
	for _, s := range mb.slots {
		s.Draw(win, matrix.Moved(slotsStartPos))
		slotsStartPos.X += s.Size().X + mtk.ConvSize(6)
	}
	// Buttons.
	menuButtonPos := pixel.V(mb.Size().X/2-mtk.ConvSize(30), mtk.ConvSize(0))
	invButtonPos := pixel.V(mb.Size().X/2-mtk.ConvSize(65), mtk.ConvSize(0))
	skillsButtonPos := pixel.V(mb.Size().X/2-mtk.ConvSize(100), mtk.ConvSize(0))
	journalButtonPos := pixel.V(-mb.Size().X/2+mtk.ConvSize(100), mtk.ConvSize(0))
	charButtonPos := pixel.V(-mb.Size().X/2+mtk.ConvSize(65), mtk.ConvSize(0))
	mb.menuButton.Draw(win, matrix.Moved(menuButtonPos))
	mb.invButton.Draw(win, matrix.Moved(invButtonPos))
	mb.skillsButton.Draw(win, matrix.Moved(skillsButtonPos))
	mb.journalButton.Draw(win, matrix.Moved(journalButtonPos))
	mb.charButton.Draw(win, matrix.Moved(charButtonPos))
}

// Update updates menu bar.
func (mb *MenuBar) Update(win *mtk.Window) {
	// Key events.
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		if !mb.DrawArea().Contains(win.MousePosition()) {
			for _, s := range mb.slots {
				if !s.Dragged() {
					continue
				}
				s.Clear()
				mb.updateLayout()
			}
		}
	}
	if win.JustPressed(pixelgl.Key1) {
		mb.useSlot(mb.slots[0])
	}
	if win.JustPressed(pixelgl.Key2) {
		mb.useSlot(mb.slots[1])
	}
	if win.JustPressed(pixelgl.Key3) {
		mb.useSlot(mb.slots[2])
	}
	if win.JustPressed(pixelgl.Key4) {
		mb.useSlot(mb.slots[3])
	}
	if win.JustPressed(pixelgl.Key5) {
		mb.useSlot(mb.slots[4])
	}
	if win.JustPressed(pixelgl.Key6) {
		mb.useSlot(mb.slots[5])
	}
	if win.JustPressed(pixelgl.Key7) {
		mb.useSlot(mb.slots[6])
	}
	if win.JustPressed(pixelgl.Key8) {
		mb.useSlot(mb.slots[7])
	}
	if win.JustPressed(pixelgl.Key9) {
		mb.useSlot(mb.slots[8])
	}
	if win.JustPressed(pixelgl.Key0) {
		mb.useSlot(mb.slots[9])
	}
	// Buttons.
	mb.menuButton.Update(win)
	mb.invButton.Update(win)
	mb.skillsButton.Update(win)
	mb.journalButton.Update(win)
	mb.charButton.Update(win)
	// Slots.
	for _, s := range mb.slots {
		s.Update(win)
	}
}

// Size returns size of bar background.
func (mb *MenuBar) Size() pixel.Vec {
	if mb.bgSpr == nil {
		return pixel.V(mtk.ConvSize(0), mtk.ConvSize(0))
	}
	return mtk.ConvVec(mb.bgSpr.Frame().Size())
}

// DrawArea return current draw area of bar background.
func (mb *MenuBar) DrawArea() pixel.Rect {
	return mb.drawArea
}

// drawIMBackground draws menu bar background with IMDraw.
func (mb *MenuBar) drawIMBackground(t pixel.Target) {
	// TODO: draw background.
}

// createSlot creates new slot for bar.
func (mb *MenuBar) createSlot() *mtk.Slot {
	params := mtk.Params{
		Size:     barSlotSize,
		FontSize: mtk.SizeMini,
	}
	s := mtk.NewSlot(params)
	s.SetOnRightClickFunc(mb.onSlotRightClicked)
	s.SetOnLeftClickFunc(mb.onSlotLeftClicked)
	return s
}

// useSlot starts action specific to current
// slot content.
func (mb *MenuBar) useSlot(s *mtk.Slot) {
	if len(s.Values()) < 1 {
		return
	}
	// Skill.
	val := s.Values()[0]
	skill, ok := val.(*object.SkillGraphic)
	if ok {
		mb.hud.Game().ActivePlayer().Use(skill.Skill)
		return
	}
	// Item.
	it, ok := val.(*object.ItemGraphic)
	if ok {
		eqit, ok := it.Item.(item.Equiper)
		if ok {
			pc := mb.hud.Game().ActivePlayer()
			if pc.Equipment().Equiped(eqit) {
				pc.Equipment().Unequip(eqit)
				return
			}
			mb.hud.inv.equip(eqit)
			return
		}
	}
}

// updateLayout updates menu bar layout for
// active player.
func (mb *MenuBar) updateLayout() {
	// Retrieve layout for current PC.
	pc := mb.hud.Game().ActivePlayer()
	layout := mb.hud.Layout(pc.ID(), pc.Serial())
	// Clear layout.
	layout.SetBarSlots(make(map[string]int))
	// Set layout.
	for i, s := range mb.slots {
		if len(s.Values()) < 1 {
			continue
		}
		for _, v := range s.Values() {
			layout.SaveBarSlot(v, i)
		}
	}
	mb.hud.layouts[pc.ID()+pc.Serial()] = layout
}

// setLayout sets specified layout as
// current bar layout.
func (mb *MenuBar) setLayout(l *Layout) {
	// Skills.
	for _, s := range mb.hud.Game().ActivePlayer().Skills() {
		slotID := l.BarSlotID(s)
		if slotID < 0 {
			continue
		}
		slot := mb.slots[slotID]
		if slot == nil {
			log.Err.Printf("hud: menu bar: set layout: unable to find slot: %d",
				slotID)
			continue
		}
		insertSlotSkill(s, slot)
	}
	// Items.
	for _, i := range mb.hud.Game().ActivePlayer().Items() {
		slotID := l.BarSlotID(i)
		if slotID < 0 {
			continue
		}
		slot := mb.slots[slotID]
		if slot == nil {
			log.Err.Printf("hud: menu bar: set layout: unable to find slot: %d",
				slotID)
			continue
		}
		mb.hud.insertSlotItem(i, slot)
	}
}

// Triggered after menu button clicked.
func (mb *MenuBar) onMenuButtonClicked(b *mtk.Button) {
	if mb.hud.menu.Opened() {
		mb.hud.menu.Show(false)
	} else {
		mb.hud.menu.Show(true)
	}
}

// Triggered after inventory button clicked.
func (mb *MenuBar) onInvButtonClicked(b *mtk.Button) {
	if mb.hud.inv.Opened() {
		mb.hud.inv.Show(false)
	} else {
		mb.hud.inv.Show(true)
	}
}

// Triggered after skills button clicked.
func (mb *MenuBar) onSkillsButtonClicked(b *mtk.Button) {
	if mb.hud.skills.Opened() {
		mb.hud.skills.Show(false)
	} else {
		mb.hud.skills.Show(true)
	}
}

// Triggered after journal button clicked.
func (mb *MenuBar) onJournalButtonClicked(b *mtk.Button) {
	if mb.hud.journal.Opened() {
		mb.hud.journal.Show(false)
	} else {
		mb.hud.journal.Show(true)
	}
}

// Triggered after character button clicked.
func (mb *MenuBar) onCharButtonClicked(b *mtk.Button) {
	if mb.hud.charinfo.Opened() {
		mb.hud.charinfo.Show(false)
	} else {
		mb.hud.charinfo.Show(true)
	}
}

// Triggered after one of bar slots was clicked with right
// mouse button.
func (mb *MenuBar) onSlotRightClicked(s *mtk.Slot) {
	for _, s := range mb.slots {
		s.Drag(false)
	}
	if len(s.Values()) < 1 {
		return
	}
	s.Drag(true)
}

// Triggered after one of bar slots was clicked with
// left mouse button.
func (mb *MenuBar) onSlotLeftClicked(s *mtk.Slot) {
	// Insert dragged skill from skill menu.
	dragSlot := mb.hud.skills.draggedSkill()
	if dragSlot != nil {
		copyMenuSlot(dragSlot, s)
		dragSlot.Drag(false)
		mb.updateLayout()
		return
	}
	// Insert dragged item from inventory menu.
	dragSlot = mb.hud.inv.draggedSlot()
	if dragSlot != nil {
		copyMenuSlot(dragSlot, s)
		dragSlot.Drag(false)
		mb.updateLayout()
		return
	}
	// Move dragged content from another bar slot.
	for _, dragSlot := range mb.slots {
		if !dragSlot.Dragged() {
			continue
		}
		switchMenuSlot(dragSlot, s)
		dragSlot.Drag(false)
	}
	// Use slot content.
	mb.useSlot(s)
}

// copyMenuSlot copies content(without label)
// from slot a to slot b.
func copyMenuSlot(a, b *mtk.Slot) {
	lab := b.Label()
	mtk.SlotCopy(a, b)
	b.SetLabel(lab)
}

// switchMenuSlot switches content between
// slots a and b(leaves slots labels unchanged).
func switchMenuSlot(a, b *mtk.Slot) {
	labA := a.Label()
	labB := b.Label()
	mtk.SlotSwitch(a, b)
	a.SetLabel(labA)
	b.SetLabel(labB)
}
