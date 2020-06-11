/*
 * skillmenu.go
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
	"github.com/isangeles/flame/module/skill"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/data/res/graphic"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

var (
	skillsKey       = pixelgl.KeyK
	skillsSlots     = 50
	skillsSlotSize  = mtk.SizeBig
	skillsSlotColor = pixel.RGBA{0.1, 0.1, 0.1, 0.5}
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

// newSkillsMenu creates new skills menu for HUD.
func newSkillMenu(hud *HUD) *SkillMenu {
	sm := new(SkillMenu)
	sm.hud = hud
	// Background.
	sm.bgDraw = imdraw.New(nil)
	bg := graphic.Textures["skillsbg.png"]
	if bg != nil {
		sm.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	}
	// Title.
	titleParams := mtk.Params{
		FontSize: mtk.SizeSmall,
	}
	sm.titleText = mtk.NewText(titleParams)
	sm.titleText.SetText(lang.Text("hud_skills_title"))
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		Shape:     mtk.ShapeSquare,
		MainColor: accentColor,
	}
	sm.closeButton = mtk.NewButton(buttonParams)
	closeButtonBG := graphic.Textures["closebutton1.png"]
	if closeButtonBG != nil {
		closeButtonSpr := pixel.NewSprite(closeButtonBG, closeButtonBG.Bounds())
		sm.closeButton.SetBackground(closeButtonSpr)
	} else {
		log.Err.Printf("hud: skills menu: unable to retrieve background texture")
	}
	sm.closeButton.SetOnClickFunc(sm.onCloseButtonClicked)
	// Slots.
	sm.slots = mtk.NewSlotList(mtk.ConvVec(pixel.V(250, 350)),
		skillsSlotColor, skillsSlotSize)
	upButtonBG := graphic.Textures["scrollup.png"]
	if upButtonBG != nil {
		upBG := pixel.NewSprite(upButtonBG, upButtonBG.Bounds())
		sm.slots.SetUpButtonBackground(upBG)
	} else {
		log.Err.Printf("hud: skills menu: unable to retrieve slot list up button texture")
	}
	downButtonBG := graphic.Textures["scrolldown.png"]
	if downButtonBG != nil {
		downBG := pixel.NewSprite(downButtonBG, downButtonBG.Bounds())
		sm.slots.SetDownButtonBackground(downBG)
	} else {
		log.Err.Printf("hud: skills menu: unable to retrieve slot list down buttons texture")
	}
	// Create empty slots.
	for i := 0; i < skillsSlots; i++ {
		s := sm.createSlot()
		sm.slots.Add(s)
	}
	return sm
}

// Draw draws menu.
func (sm *SkillMenu) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	sm.drawArea = mtk.MatrixToDrawArea(matrix, sm.Size())
	// Background.
	if sm.bgSpr != nil {
		sm.bgSpr.Draw(win.Window, matrix)
	} else {
		mtk.DrawRectangle(win.Window, sm.DrawArea(), nil)
	}
	// Title.
	titleTextPos := pixel.V(mtk.ConvSize(0), sm.Size().Y/2-mtk.ConvSize(25))
	sm.titleText.Draw(win.Window, matrix.Moved(titleTextPos))
	// Buttons.
	closeButtonPos := pixel.V(sm.Size().X/2-mtk.ConvSize(20),
		sm.Size().Y/2-mtk.ConvSize(15))
	sm.closeButton.Draw(win.Window, matrix.Moved(closeButtonPos))
	// Slots.
	slotsPos := pixel.V(mtk.ConvSize(0), mtk.ConvSize(-10))
	sm.slots.Draw(win, matrix.Moved(slotsPos))
}

// Update updates window.
func (sm *SkillMenu) Update(win *mtk.Window) {
	// Key events.
	if !sm.hud.Chat().Activated() && win.JustPressed(skillsKey) {
		sm.Show(!sm.Opened())
	}
	// Elements update.
	if sm.Opened() {
		sm.closeButton.Update(win)
		sm.slots.Update(win)
	}
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

// Size returns size of menu background.
func (sm *SkillMenu) Size() pixel.Vec {
	if sm.bgSpr == nil {
		// TODO: size for draw background.
		return pixel.V(mtk.ConvSize(0), mtk.ConvSize(0))
	}
	return mtk.ConvVec(sm.bgSpr.Frame().Size())
}

// insert inserts specified skills in menu slots.
func (sm *SkillMenu) insert(skills ...*object.SkillGraphic) {
	for _, s := range skills {
		slot := sm.slots.EmptySlot()
		if slot == nil {
			log.Err.Printf("hud: skills menu: no empty slots")
			return
		}
		// Insert skill to slot.
		insertSlotSkill(s, slot)
	}
}

// createSlot creates empty slot for skills slot list.
func (sm *SkillMenu) createSlot() *mtk.Slot {
	params := mtk.Params{
		Size:      skillsSlotSize,
		FontSize:  mtk.SizeMini,
		MainColor: skillsSlotColor,
	}
	s := mtk.NewSlot(params)
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
		log.Err.Printf("hud: skills menu: %v: is not skill", s.Values()[0])
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
	if config.Debug {
		info = fmt.Sprintf("%s\n[%s]", info,
			s.ID())
	}
	return info
}
