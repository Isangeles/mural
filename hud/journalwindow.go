/*
 * journalwindow.go
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

	"github.com/isangeles/flame/core/data/res/lang"
	"github.com/isangeles/flame/core/module/quest"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/log"
)

var (
	journalKey = pixelgl.KeyL
)

// Struct for HUD journal window.
type JournalWindow struct {
	hud         *HUD
	bgSpr       *pixel.Sprite
	bgDraw      *imdraw.IMDraw
	drawArea    pixel.Rect
	titleText   *mtk.Text
	closeButton *mtk.Button
	opened      bool
	focused     bool
	questInfo   *mtk.Textbox
	questsList  *mtk.List
}

// newJournalWindow creates new journal window
// for HUD.
func newJournalWindow(hud *HUD) *JournalWindow {
	jw := new(JournalWindow)
	jw.hud = hud
	// Background.
	jw.bgDraw = imdraw.New(nil)
	bg, err := data.PictureUI("menubg.png")
	if err == nil {
		jw.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	} else {
		log.Err.Printf("hud journal: fail to retrieve bg tex: %v",
			err)
	}
	// Title.
	titleParams := mtk.Params{
		FontSize: mtk.SizeSmall,
	}
	jw.titleText = mtk.NewText(titleParams)
	jw.titleText.SetText(lang.Text("hud_journal_title"))
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		Shape:     mtk.ShapeSquare,
		MainColor: accentColor,
	}
	jw.closeButton = mtk.NewButton(buttonParams)
	closeButtonBG, err := data.PictureUI("closebutton1.png")
	if err == nil {
		closeBG := pixel.NewSprite(closeButtonBG,
			closeButtonBG.Bounds())
		jw.closeButton.SetBackground(closeBG)
	}
	jw.closeButton.SetOnClickFunc(jw.onCloseButtonClicked)
	// Quest info.
	questInfoSize := pixel.V(jw.Size().X-mtk.ConvSize(20),
		jw.Size().Y/2)
	questInfoParams := mtk.Params{
		SizeRaw:     questInfoSize,
		FontSize:    mtk.SizeSmall,
		MainColor:   mainColor,
		AccentColor: accentColor,
	}
	jw.questInfo = mtk.NewTextbox(questInfoParams)
	// Quests list.
	questsSize := pixel.V(jw.Size().X-mtk.ConvSize(20),
		jw.Size().Y/2-mtk.ConvSize(100))
	questsParams := mtk.Params{
		SizeRaw:     questsSize,
		MainColor:   mainColor,
		SecColor:    secColor,
		AccentColor: accentColor,
	}
	jw.questsList = mtk.NewList(questsParams)
	upButtonBG, err := data.PictureUI("scrollup.png")
	if err == nil {
		upBG := pixel.NewSprite(upButtonBG,
			upButtonBG.Bounds())
		jw.questsList.SetUpButtonBackground(upBG)
	}
	downButtonBG, err := data.PictureUI("scrolldown.png")
	if err == nil {
		downBG := pixel.NewSprite(downButtonBG,
			downButtonBG.Bounds())
		jw.questsList.SetDownButtonBackground(downBG)
	}
	jw.questsList.SetOnItemSelectFunc(jw.onQuestSelected)
	return jw
}

// Draw draws window.
func (jw *JournalWindow) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	jw.drawArea = mtk.MatrixToDrawArea(matrix, jw.Size())
	// Background.
	if jw.bgSpr != nil {
		jw.bgSpr.Draw(win, matrix)
	} else {
		mtk.DrawRectangle(win, jw.DrawArea(), nil)
	}
	// Title.
	titleTextMove := pixel.V(0, jw.Size().Y/2-mtk.ConvSize(25))
	jw.titleText.Draw(win, matrix.Moved(titleTextMove))
	// Buttons.
	closeButtonMove := mtk.ConvVec(pixel.V(jw.Size().X/2-20,
		jw.Size().Y/2-15))
	jw.closeButton.Draw(win, matrix.Moved(closeButtonMove))
	// Quest info.
	questInfoMove := mtk.MoveTC(jw.Size(), jw.questInfo.Size())
	questInfoMove.Y -= mtk.ConvSize(50)
	jw.questInfo.Draw(win, matrix.Moved(questInfoMove))
	// Quests list.
	questsMove := mtk.MoveBC(jw.Size(), jw.questsList.Size())
	questsMove.Y += mtk.ConvSize(20)
	jw.questsList.Draw(win, matrix.Moved(questsMove))
}

// Update updates window.
func (jw *JournalWindow) Update(win *mtk.Window) {
	// Key events.
	if !jw.hud.Chat().Activated() && win.JustPressed(journalKey) {
		jw.Show(!jw.Opened())
	}
	// Elements.
	if jw.Opened() {
		jw.closeButton.Update(win)
		jw.questsList.Update(win)
		jw.questInfo.Update(win)
	}
}

// Show toggles window visibility.
func (jw *JournalWindow) Show(show bool) {
	jw.opened = show
	if jw.Opened() {
		jw.questsList.Clear()
		pc := jw.hud.ActivePlayer()
		jw.insertQuests(pc.Journal().Quests()...)
	} else {
		jw.questInfo.Clear()
	}
}

// Opened checks if window is open.
func (jw *JournalWindow) Opened() bool {
	return jw.opened
}

// DrawArea returns window draw area.
func (jw *JournalWindow) DrawArea() pixel.Rect {
	return jw.drawArea
}

// Size returns window size.
func (jw *JournalWindow) Size() pixel.Vec {
	if jw.bgSpr == nil {
		return mtk.ConvVec(pixel.V(50, 200))
	}
	return jw.bgSpr.Frame().Size()
}

// insertQuests adds all specified quests to journal
// quests list.
func (jw *JournalWindow) insertQuests(quests ...*quest.Quest) {
	for _, q := range quests {
		questText := lang.Text(q.ID())
		jw.questsList.AddItem(questText, q)
	}
}

// Triggered after close button clicked.
func (jw *JournalWindow) onCloseButtonClicked(b *mtk.Button) {
	jw.Show(false)
}

// Triggered after selecting quest from quests list.
func (jw *JournalWindow) onQuestSelected(cs *mtk.CheckSlot) {
	// Retrive quest from slot.
	quest, ok := cs.Value().(*quest.Quest)
	if !ok {
		log.Err.Printf("hud journal: fail to retrive quest from list")
		return
	}
	// Show quest info.
	questInfo := lang.Texts(quest.ID())
	info := questInfo[0]
	if len(questInfo) > 1 {
		info = fmt.Sprintf("%s\n%s", info, questInfo[1])
	}
	stage := quest.ActiveStage()
	if stage != nil {
		if stage.Completed() {
			completeInfo := lang.Text("hud_journal_quest_complete")
			info = fmt.Sprintf("%s\n%s", info, completeInfo)
		} else {
			stageInfo := lang.Text(stage.ID())
			info = fmt.Sprintf("%s\n%s", info, stageInfo)
		}
	}
	jw.questInfo.SetText(info)
}
