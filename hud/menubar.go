/*
 * menubar.go
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
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"

	"github.com/isangeles/flame/core/data/text/lang"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// Struct for HUD menu bar.
type MenuBar struct {
	hud          *HUD
	bgSpr        *pixel.Sprite
	bgDraw       *imdraw.IMDraw
	drawArea     pixel.Rect
	menuButton   *mtk.Button
	invButton    *mtk.Button
	skillsButton *mtk.Button
	slots        []*mtk.Slot
}

var (
	bar_slots      = 10
	bar_slot_size  = mtk.SIZE_MEDIUM
	bar_slot_color = pixel.RGBA{0.1, 0.1, 0.1, 0.5}
)

// neMenuBar creates new menu bar for HUD.
func newMenuBar(hud *HUD) *MenuBar {
	mb := new(MenuBar)
	mb.hud = hud
	// Background.
	bg, err := data.PictureUI("menubar.png")
	if err != nil { // fallback
		mb.bgDraw = imdraw.New(nil)
	} else {
		mb.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	}
	// Menu Button.
	mb.menuButton = mtk.NewButton(mtk.SIZE_MINI, mtk.SHAPE_SQUARE, accent_color,
		"", lang.Text("gui", "hud_bar_menu_open_info"))
	menuButtonBG, err := data.PictureUI("menubutton.png")
	if err != nil {
		log.Err.Printf("hud_menu_bar:fail_to_retrieve_menu_button_texture:%v", err)
	} else {
		menuButtonSpr := pixel.NewSprite(menuButtonBG, menuButtonBG.Bounds())
		mb.menuButton.SetBackground(menuButtonSpr)
	}
	mb.menuButton.SetOnClickFunc(mb.onMenuButtonClicked)
	// Inventory button.
	mb.invButton = mtk.NewButton(mtk.SIZE_MINI, mtk.SHAPE_SQUARE, accent_color,
		"", lang.Text("gui", "hud_bar_inv_open_info"))
	invButtonBG, err := data.PictureUI("inventorybutton.png")
	if err != nil {
		log.Err.Printf("hud_menu_bar:fail_to_retrieve_inv_button_texture:%v", err)
	} else {
		invButtonSpr := pixel.NewSprite(invButtonBG, invButtonBG.Bounds())
		mb.invButton.SetBackground(invButtonSpr)
	}
	mb.invButton.SetOnClickFunc(mb.onInvButtonClicked)
	// Skills button.
	mb.skillsButton = mtk.NewButton(mtk.SIZE_MINI, mtk.SHAPE_SQUARE, accent_color,
		"", lang.Text("gui", "hud_bar_skills_open_info"))
	skillsButtonBG, err := data.PictureUI("skillsbutton.png")
	if err != nil {
		log.Err.Printf("hud_menu_bar:fail_to_retrieve_skills_button_texture:%v", err)
	} else {
		skillsButtonSpr := pixel.NewSprite(skillsButtonBG, skillsButtonBG.Bounds())
		mb.skillsButton.SetBackground(skillsButtonSpr)
	}
	mb.skillsButton.SetOnClickFunc(mb.onSkillsButtonClicked)
	// Slots.
	for i := 0; i < bar_slots; i++ {
		s := mb.createSlot()
		mb.slots = append(mb.slots, s)
	}
	return mb
}

// Draw draws menu bar.
func (mb *MenuBar) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	mb.drawArea = mtk.MatrixToDrawArea(matrix, mb.Bounds())
	// Background.
	if mb.bgSpr != nil {
		mb.bgSpr.Draw(win.Window, matrix)
	} else {
		mb.drawIMBackground(win.Window)
	}
	// Buttons.
	menuButtonPos := mtk.ConvVec(pixel.V(mb.Bounds().Max.X/2-30, 0))
	mb.menuButton.Draw(win.Window, matrix.Moved(menuButtonPos))
	invButtonPos := mtk.ConvVec(pixel.V(mb.Bounds().Max.X/2-65, 0))
	mb.invButton.Draw(win.Window, matrix.Moved(invButtonPos))
	skillsButtonPos := mtk.ConvVec(pixel.V(mb.Bounds().Max.X/2-100, 0))
	mb.skillsButton.Draw(win.Window, matrix.Moved(skillsButtonPos))
	// Slots.
	slotsStartPos := mtk.ConvVec(pixel.V(-163, 0))
	for _, s := range mb.slots {
		s.Draw(win.Window, matrix.Moved(slotsStartPos))
		slotsStartPos.X += s.Bounds().W() + mtk.ConvSize(6)
	}
}

// Update updates menu bar.
func (mb *MenuBar) Update(win *mtk.Window) {
	// Key events.
	if win.JustPressed(pixelgl.Key1) {
		mb.useSlot(mb.slots[0])
	}
	// Buttons.
	mb.menuButton.Update(win)
	mb.invButton.Update(win)
	mb.skillsButton.Update(win)
	// Slots.
	for _, s := range mb.slots {
		s.Update(win)
	}
}

// Bounds returns bounds of bar background.
func (mb *MenuBar) Bounds() pixel.Rect {
	if mb.bgSpr == nil {
		return pixel.R(0, 0, mtk.ConvSize(0), mtk.ConvSize(0))
	}
	return mb.bgSpr.Frame()
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
	s := mtk.NewSlot(bar_slot_size, mtk.SIZE_MINI)
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
	skill, ok := s.Values()[0].(*object.SkillGraphic)
	if ok {
		mb.hud.ActivePlayer().UseSkill(skill.Skill)
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
	skillSlot := mb.hud.skills.draggedSkill()
	if skillSlot != nil {
		mtk.SlotSwitch(s, skillSlot)
		skillSlot.Drag(false)
		return
	}
	mb.useSlot(s)
}
