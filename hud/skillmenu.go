/*
 * skillmenu.go
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

	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/object/skill"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

// Struct for skills menu.
type SkillMenu struct {
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
	skills_slots      = 50
	skills_slot_size  = mtk.SIZE_BIG
	skills_slot_color = pixel.RGBA{0.1, 0.1, 0.1, 0.5}
)

// newSkillsMenu creates new skills menu for HUD.
func newSkillMenu(hud *HUD) *SkillMenu {
	sm := new(SkillMenu)
	sm.hud = hud
	// Background.
	sm.bgDraw = imdraw.New(nil)
	bg, err := data.PictureUI("skillsbg.png")
	if err == nil {
		sm.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	}
	// Title.
	sm.titleText = mtk.NewText(mtk.SIZE_SMALL, 0)
	sm.titleText.SetText(lang.Text("gui", "hud_skills_title"))
	// Buttons.
	sm.closeButton = mtk.NewButton(mtk.SIZE_SMALL, mtk.SHAPE_SQUARE, accent_color)
	closeButtonBG, err := data.PictureUI("closebutton1.png")
	if err != nil {
		log.Err.Printf("hud_skills:fail_to_retrieve_background_tex:%v", err)
	} else {
		closeButtonSpr := pixel.NewSprite(closeButtonBG, closeButtonBG.Bounds())
		sm.closeButton.SetBackground(closeButtonSpr)
	}
	sm.closeButton.SetOnClickFunc(sm.onCloseButtonClicked)
	// Slots.
	sm.slots = mtk.NewSlotList(mtk.ConvVec(pixel.V(250, 350)),
		inv_slot_color, inv_slot_size)
	upButtonBG, err := data.PictureUI("scrollup.png")
	if err != nil {
		log.Err.Printf("hud_inv:fail_to_retrieve_slot_list_up_buttons_texture:%v",
			err)
	} else {
		upBG := pixel.NewSprite(upButtonBG, upButtonBG.Bounds())
		sm.slots.SetUpButtonBackground(upBG)
	}
	downButtonBG, err := data.PictureUI("scrolldown.png")
	if err != nil {
		log.Err.Printf("hud_inv:fail_to_retrieve_slot_list_down_buttons_texture:%v",
			err)
	} else {
		downBG := pixel.NewSprite(downButtonBG, downButtonBG.Bounds())
		sm.slots.SetDownButtonBackground(downBG)
	}
	// Create empty slots.
	for i := 0; i < skills_slots; i++ {
		s := sm.createSlot()
		sm.slots.Add(s)
	}
	return sm
}

// Draw draws menu.
func (sm *SkillMenu) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	sm.drawArea = mtk.MatrixToDrawArea(matrix, sm.Bounds())
	// Background.
	if sm.bgSpr != nil {
		sm.bgSpr.Draw(win.Window, matrix)
	} else {
		mtk.DrawRectangle(win.Window, sm.DrawArea(), nil)
	}
	// Title.
	titleTextPos := mtk.ConvVec(pixel.V(0, sm.Bounds().Max.Y/2-25))
	sm.titleText.Draw(win.Window, matrix.Moved(titleTextPos))
	// Buttons.
	closeButtonPos := mtk.ConvVec(pixel.V(sm.Bounds().Max.X/2-20,
		sm.Bounds().Max.Y/2-15))
	sm.closeButton.Draw(win.Window, matrix.Moved(closeButtonPos))
	// Slots.
	slotsPos := pixel.V(mtk.ConvSize(0), mtk.ConvSize(-10))
	sm.slots.Draw(win, matrix.Moved(slotsPos))
}

// Update updates window.
func (sm *SkillMenu) Update(win *mtk.Window) {
	// Elements update.
	sm.closeButton.Update(win)
	sm.slots.Update(win)
}

// Opened checks whether menu is open.
func (sm *SkillMenu) Opened() bool {
	return sm.opened
}

// Show toggles menu visibility.
func (sm *SkillMenu) Show(show bool) {
	sm.opened = show
	if sm.Opened() {
		sm.slots.Clear()
		sm.insert(sm.hud.ActivePlayer().Skills()...)
		sm.hud.UserFocus().Focus(sm)
	} else {
		sm.hud.UserFocus().Focus(nil)
	}
}

// Focused checks whether menu us focused.
func (sm *SkillMenu) Focused() bool {
	return sm.focused
}

// Focus toggles menu focus.
func (sm *SkillMenu) Focus(focus bool) {
	sm.focused = focus
}

// DrawArea returns menu draw area.
func (sm *SkillMenu) DrawArea() pixel.Rect {
	return sm.drawArea
}

// Bounds returns size bounds of menu background.
func (sm *SkillMenu) Bounds() pixel.Rect {
	if sm.bgSpr == nil {
		// TODO: bounds for draw background.
		return pixel.R(0, 0, mtk.ConvSize(0), mtk.ConvSize(0))
	}
	return sm.bgSpr.Frame()
}

// insert inserts specified skills in menu slots.
func (sm *SkillMenu) insert(skills ...*object.SkillGraphic) {
	for _, s := range skills {
		slot := sm.slots.EmptySlot()
		if slot == nil {
			log.Err.Printf("hud_skills:no empty slots")
			return
		}
		// Insert skill to slot.
		insertSlotSkill(s, slot)
	}
}

// createSlot creates empty slot for skills slot list.
func (sm *SkillMenu) createSlot() *mtk.Slot {
	s := mtk.NewSlot(skills_slot_size, mtk.SIZE_MINI)
	s.SetColor(skills_slot_color)
	s.SetOnRightClickFunc(sm.onSlotRightClicked)
	s.SetOnLeftClickFunc(sm.onSlotLeftClicked)
	return s
}

// draggedSkill returns currently dragged slot
// with skill.
func (sm *SkillMenu) draggedSkill() *mtk.Slot {
	for _, s := range sm.slots.Slots() {
		if s.Dragged() {
			return s
		}
	}
	return nil
}

// Triggered after close button clicked.
func (sm *SkillMenu) onCloseButtonClicked(b *mtk.Button) {
	sm.Show(false)
}

// Triggered after one of skill slots was clicked with
// right mouse button.
func (sm *SkillMenu) onSlotRightClicked(s *mtk.Slot) {
	if len(s.Values()) < 1 {
		return
	}
	skill, ok := s.Values()[0].(*object.SkillGraphic)
	if !ok {
		log.Err.Printf("hud_skills_menu:%v:is not skill", s.Values()[0])
	}
	sm.hud.ActivePlayer().UseSkill(skill.Skill)
}

// Triggered after one of skill slots was clicked with
// left mouse button.
func (sm *SkillMenu) onSlotLeftClicked(s *mtk.Slot) {
	for _, s := range sm.slots.Slots() {
		s.Drag(false)
	}
	if len(s.Values()) < 1 {
		return
	}
	s.Drag(true)
}

// insertSlotSkill inserts specified skill to specified slot.
func insertSlotSkill(skill *object.SkillGraphic, slot *mtk.Slot) {
	slot.AddValues(skill)
	slot.SetIcon(skill.Icon())
	slot.SetInfo(skillInfo(skill.Skill))
}

// skillInfo returns formated string with
// informations about specified skill.
func skillInfo(s *skill.Skill) string {
	infoForm := "%s"
	info := fmt.Sprintf(infoForm, s.Name())
	if config.Debug() {
		info = fmt.Sprintf("%s\n[%s_%s]", info,
			s.ID(), s.Serial())
	}
	return info
}
