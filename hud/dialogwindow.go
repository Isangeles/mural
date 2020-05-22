/*
 * dialogwindow.go
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

	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/module/dialog"
	"github.com/isangeles/flame/module/item"
	"github.com/isangeles/flame/module/train"

	"github.com/isangeles/mtk"
 
	"github.com/isangeles/mural/core/data/res/graphic"
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
	bg := graphic.Textures["menubg.png"]
	if bg != nil {
		dw.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	}
	// Title.
	titleParams := mtk.Params{
		FontSize: mtk.SizeSmall,
	}
	dw.titleText = mtk.NewText(titleParams)
	dw.titleText.SetText(lang.Text("hud_dialog_title"))
	// Buttons.
	buttonParams := mtk.Params{
		Size:      mtk.SizeMedium,
		Shape:     mtk.ShapeSquare,
		MainColor: accentColor,
	}
	dw.closeButton = mtk.NewButton(buttonParams)
	closeButtonBG := graphic.Textures["closebutton1.png"]
	if closeButtonBG != nil {
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
	answersParams := mtk.Params{
		SizeRaw:     answersSize,
		MainColor:   mainColor,
		SecColor:    secColor,
		AccentColor: accentColor,
	}
	dw.answersList = mtk.NewList(answersParams)
	upButtonBG := graphic.Textures["scrollup.png"]
	if upButtonBG != nil {
		upBG := pixel.NewSprite(upButtonBG, upButtonBG.Bounds())
		dw.answersList.SetUpButtonBackground(upBG)
	}
	downButtonBG := graphic.Textures["scrolldown.png"]
	if downButtonBG != nil {
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
	if dw.Opened() {
		dw.closeButton.Update(win)
		dw.chatBox.Update(win)
		dw.answersList.Update(win)
	}
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
		if dw.dialog != nil {
			dw.dialog.Restart()
		}
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
	dw.dialog.SetTarget(dw.hud.ActivePlayer())
	dw.dialogUpdate()
}

// dialogUpdate updates window components to
// current dialog stage.
func (dw *DialogWindow) dialogUpdate() {
	if dw.dialog == nil || dw.dialog.Finished() {
		return
	}
	if dw.dialog.Stage() == nil {
		log.Err.Printf("hud: dialog window: no suitable dialog phase found")
		return
	}
	// Print stage text to chat box.
	text := fmt.Sprintf("[%s]: %s\n", dw.dialog.Owner().Name(), dw.dialog.Stage())
	dw.chatBox.AddText(text)
	dw.chatBox.ScrollBottom()
	// Select answers.
	answers := make([]*dialog.Answer, 0)
	for _, a := range dw.dialog.Stage().Answers() {
		if dw.hud.ActivePlayer().MeetReqs(a.Requirements()...) {
			answers = append(answers, a)
		}
	}
	// Insert answers to answers list.
	dw.answersList.Clear()
	for i, a := range answers {
		answerText := fmt.Sprintf("%d) %s", i, a)
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
		log.Err.Printf("hud: dialog window: unable to retrieve answer from list")
		return
	}
	// Print answer to chat box.
	dw.chatBox.AddText(fmt.Sprintf("[%s]: %s\n", dw.hud.ActivePlayer().Name(), answer))
	dw.chatBox.ScrollBottom()
	// Move dialog forward.
	dw.dialog.Next(answer)
	// On trade.
	if dw.dialog.Trading() {
		con, ok := dw.dialog.Owner().(item.Container)
		if !ok {
			log.Err.Printf("hud: dialog window: dialog onwer has no inventory")
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
			log.Err.Printf("hud: dialog window: dialog onwer is not a trainer")
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
