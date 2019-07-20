/*
 * dialogwindow.go
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

	flameconf "github.com/isangeles/flame/config"
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/object/dialog"
	"github.com/isangeles/flame/core/module/object/effect"
	"github.com/isangeles/flame/core/module/object/item"
	"github.com/isangeles/flame/core/module/train"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/log"
)

// Struct for HUD dialog window.
type DialogWindow struct {
	hud         *HUD
	bgSpr       *pixel.Sprite
	bgDraw      *imdraw.IMDraw
	drawArea    pixel.Rect
	titleText   *mtk.Text
	closeButton *mtk.Button
	chatBox     *mtk.Textbox
	answersList *mtk.List
	opened      bool
	focused     bool
	dialog      *dialog.Dialog
}

// newDialogWindow creates new dialog
// window for HUD.
func newDialogWindow(hud *HUD) *DialogWindow {
	dw := new(DialogWindow)
	dw.hud = hud
	// Background.
	dw.bgDraw = imdraw.New(nil)
	bg, err := data.PictureUI("menubg.png")
	if err == nil {
		dw.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	}
	// Title.
	dw.titleText = mtk.NewText(mtk.SizeSmall, 0)
	dw.titleText.SetText(lang.TextDir(flameconf.LangPath(), "hud_dialog_title"))
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		Shape:     mtk.ShapeSquare,
		MainColor: accentColor,
	}
	dw.closeButton = mtk.NewButton(buttonParams)
	closeButtonBG, err := data.PictureUI("closebutton1.png")
	if err == nil {
		closeBG := pixel.NewSprite(closeButtonBG, closeButtonBG.Bounds())
		dw.closeButton.SetBackground(closeBG)
	}
	dw.closeButton.SetOnClickFunc(dw.onCloseButtonClicked)
	// Chat.
	chatSize := pixel.V(dw.Size().X-mtk.ConvSize(20),
		dw.Size().Y/2)
	chatParams := mtk.Params{
		SizeRaw:     chatSize,
		FontSize:    mtk.SizeMedium,
		MainColor:   mainColor,
		AccentColor: accentColor,
	}
	dw.chatBox = mtk.NewTextbox(chatParams)
	// Answers list.
	answersSize := pixel.V(dw.Size().X-mtk.ConvSize(20),
		dw.Size().Y/2-mtk.ConvSize(100))
	dw.answersList = mtk.NewList(answersSize, mtk.SizeSmall, mainColor, secColor,
		accentColor)
	upButtonBG, err := data.PictureUI("scrollup.png")
	if err == nil {
		upBG := pixel.NewSprite(upButtonBG, upButtonBG.Bounds())
		dw.answersList.SetUpButtonBackground(upBG)
	}
	downButtonBG, err := data.PictureUI("scrolldown.png")
	if err == nil {
		downBG := pixel.NewSprite(downButtonBG, downButtonBG.Bounds())
		dw.answersList.SetDownButtonBackground(downBG)
	}
	dw.answersList.SetOnItemSelectFunc(dw.onAnswerSelected)
	return dw
}

// Draw draws window.
func (dw *DialogWindow) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	dw.drawArea = mtk.MatrixToDrawArea(matrix, dw.Size())
	// Background.
	if dw.bgSpr != nil {
		dw.bgSpr.Draw(win.Window, matrix)
	} else {
		mtk.DrawRectangle(win.Window, dw.DrawArea(), mainColor)
	}
	// Title.
	titleTextMove := pixel.V(0, dw.Size().Y/2-mtk.ConvSize(20))
	dw.titleText.Draw(win.Window, matrix.Moved(titleTextMove))
	// Buttons.
	closeButtonPos := mtk.ConvVec(pixel.V(dw.Size().X/2-20,
		dw.Size().Y/2-15))
	dw.closeButton.Draw(win.Window, matrix.Moved(closeButtonPos))
	// Chat & answers list.
	chatMove := mtk.MoveTC(dw.Size(), dw.chatBox.Size())
	chatMove.Y -= mtk.ConvSize(50)
	answersMove := mtk.MoveBC(dw.Size(), dw.answersList.Size())
	answersMove.Y += mtk.ConvSize(20)
	dw.chatBox.Draw(win, matrix.Moved(chatMove))
	dw.answersList.Draw(win, matrix.Moved(answersMove))
}

// Updates updates window.
func (dw *DialogWindow) Update(win *mtk.Window) {
	// Elements.
	dw.closeButton.Update(win)
	dw.chatBox.Update(win)
	dw.answersList.Update(win)
	// Dialog.
	if dw.dialog != nil {
		if dw.dialog.Finished() {
			dw.Show(false)
		}
	}
}

// Show toggles window visibility.
func (dw *DialogWindow) Show(show bool) {
	dw.opened = show
	if !dw.Opened() {
		dw.chatBox.Clear()
		dw.answersList.Clear()
	}
}

// Opened checks if window is open.
func (dw *DialogWindow) Opened() bool {
	return dw.opened
}

// Size returns window size.
func (dw *DialogWindow) Size() pixel.Vec {
	if dw.bgSpr == nil {
		return mtk.ConvVec(pixel.V(0, 0)) // TODO: draw bg size
	}
	return dw.bgSpr.Frame().Size()
}

// DrawArea returns current draw area.
func (dw *DialogWindow) DrawArea() pixel.Rect {
	return dw.drawArea
}

// SetDialog sets dialog for window.
func (dw *DialogWindow) SetDialog(d *dialog.Dialog) {
	dw.dialog = d
	dw.dialog.Restart()
	dw.dialogUpdate()
}

// dialogUpdate updates window components to
// current dialog phase.
func (dw *DialogWindow) dialogUpdate() {
	if dw.dialog == nil || dw.dialog.Finished() {
		return
	}
	// Search for proper dialog phase.
	var phase *dialog.Phase
	for _, p := range dw.dialog.Phases() {
		if dw.hud.ActivePlayer().MeetReqs(p.Requirements()...) {
			phase = p
		}
	}
	if phase == nil {
		log.Err.Printf("hud_dialog:no suitable dialog phase found")
		return
	}
	// Print phase text to chat box.
	chapter := dw.hud.game.Module().Chapter()
	dialogLine := lang.AllText(chapter.Conf().DialogsLangPath(), phase.ID())
	text := fmt.Sprintf("[%s]:%s\n", dw.dialog.Owner().Name(), dialogLine[0])
	dw.chatBox.AddText(text)
	dw.chatBox.ScrollBottom()
	// Apply phase modifiers.
	if tar, ok := dw.dialog.Owner().(effect.Target); ok {
		for _, mod := range phase.OwnerModifiers() {
			mod.Affect(dw.hud.ActivePlayer().Character, tar)
		}
		for _, mod := range phase.TalkerModifiers() {
			mod.Affect(tar, dw.hud.ActivePlayer().Character)
		}
	} else {
		for _, mod := range phase.TalkerModifiers() {
			mod.Affect(nil, dw.hud.ActivePlayer().Character)
		}
	}
	// Select answers.
	answers := make([]*dialog.Answer, 0)
	for _, a := range phase.Answers() {
		if dw.hud.ActivePlayer().MeetReqs(a.Requirements()...) {
			answers = append(answers, a)
		}
	}
	// Insert answers to answers list.
	dw.answersList.Clear()
	for i, a := range answers {
		answerText := lang.AllText(chapter.Conf().DialogsLangPath(), a.ID())[0]
		answerText = fmt.Sprintf("%d)%s", i, answerText)
		dw.answersList.AddItem(answerText, a)
	}
}

// Triggered after clicking close button.
func (dw *DialogWindow) onCloseButtonClicked(b *mtk.Button) {
	dw.Show(false)
}

// Triggered after selecting answer from answers list.
func (dw *DialogWindow) onAnswerSelected(cs *mtk.CheckSlot) {
	if dw.dialog == nil {
		return
	}
	// Retrieve answer from slot.
	answer, ok := cs.Value().(*dialog.Answer)
	if !ok {
		log.Err.Printf("hud_dialog:fail to retrieve answer from list")
		return
	}
	// Print answer to chat box.
	chapter := dw.hud.game.Module().Chapter()
	answerText := lang.AllText(chapter.Conf().DialogsLangPath(), answer.ID())[0]
	dw.chatBox.AddText(fmt.Sprintf("[%s]:%s\n", dw.hud.ActivePlayer().Name(), answerText))
	dw.chatBox.ScrollBottom()
	// Move dialog forward.
	dw.dialog.Next(answer)
	// Apply answer modifiers.
	if tar, ok := dw.dialog.Owner().(effect.Target); ok {
		for _, mod := range answer.OwnerModifiers() {
			mod.Affect(dw.hud.ActivePlayer().Character, tar)
		}
		for _, mod := range answer.TalkerModifiers() {
			mod.Affect(tar, dw.hud.ActivePlayer().Character)
		}
	} else {
		for _, mod := range answer.TalkerModifiers() {
			mod.Affect(nil, dw.hud.ActivePlayer().Character)
		}
	}
	// On trade.
	if dw.dialog.Trading() {
		con, ok := dw.dialog.Owner().(item.Container)
		if !ok {
			log.Err.Printf("hud_dialog:dialog onwer has no inventory")
			return
		}
		dw.Show(false)
		dw.hud.trade.SetSeller(con)
		dw.hud.trade.Show(true)
		return
	}
	// On training.
	if dw.dialog.Training() {
		tra, ok := dw.dialog.Owner().(train.Trainer)
		if !ok {
			log.Err.Printf("hud_dialog:dialog onwer is not a trainer")
			return
		}
		dw.Show(false)
		dw.hud.training.SetTrainer(tra)
		dw.hud.training.Show(true)
		return
	}
	// Update dialog view.
	dw.dialogUpdate()
}
