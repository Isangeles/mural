/*
 * trainingwindow.go
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

	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/req"
	"github.com/isangeles/flame/training"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data/res/graphic"
	"github.com/isangeles/mural/log"
)

// Struct for HUD training window.
type TrainingWindow struct {
	hud           *HUD
	bgSpr         *pixel.Sprite
	bgDraw        *imdraw.IMDraw
	drawArea      pixel.Rect
	titleText     *mtk.Text
	closeButton   *mtk.Button
	trainButton   *mtk.Button
	trainingInfo  *mtk.Textbox
	trainingsList *mtk.List
	opened        bool
	focused       bool
	trainer       training.Trainer
}

// newTrainingWindow creates new training
// window for HUD.
func newTrainingWindow(hud *HUD) *TrainingWindow {
	tw := new(TrainingWindow)
	tw.hud = hud
	// Background.
	tw.bgDraw = imdraw.New(nil)
	bg := graphic.Textures["menubg.png"]
	if bg != nil {
		tw.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	} else {
		log.Err.Printf("hud training: unable to retrieve background texture")
	}
	// Title.
	titleParams := mtk.Params{
		FontSize: mtk.SizeSmall,
	}
	tw.titleText = mtk.NewText(titleParams)
	tw.titleText.SetText(lang.Text("hud_training_title"))
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
		log.Err.Printf("hud training: unable to retrieve close button texture")
	}
	tw.closeButton.SetOnClickFunc(tw.onCloseButtonClicked)
	// Train button.
	trainButtonParams := mtk.Params{
		Size:      mtk.SizeMini,
		FontSize:  mtk.SizeMini,
		Shape:     mtk.ShapeRectangle,
		MainColor: accentColor,
	}
	tw.trainButton = mtk.NewButton(trainButtonParams)
	trainButtonBG := graphic.Textures["button_green.png"]
	if trainButtonBG != nil {
		bg := pixel.NewSprite(trainButtonBG, trainButtonBG.Bounds())
		tw.trainButton.SetBackground(bg)
	} else {
		log.Err.Printf("hud training: unable to retrieve train button texture")
	}
	tw.trainButton.SetOnClickFunc(tw.onTrainButtonClicked)
	tw.trainButton.SetLabel(lang.Text("hud_training_train"))
	// Training info.
	infoSize := pixel.V(tw.Size().X-mtk.ConvSize(20),
		tw.Size().Y/2-mtk.ConvSize(10))
	trainingInfoParams := mtk.Params{
		SizeRaw:     infoSize,
		FontSize:    mtk.SizeSmall,
		MainColor:   mainColor,
		AccentColor: accentColor,
	}
	tw.trainingInfo = mtk.NewTextbox(trainingInfoParams)
	// Trainings list.
	trainingsSize := pixel.V(tw.Size().X-mtk.ConvSize(20),
		tw.Size().Y/2-mtk.ConvSize(100))
	trainingsParams := mtk.Params{
		SizeRaw:     trainingsSize,
		MainColor:   mainColor,
		SecColor:    secColor,
		AccentColor: accentColor,
	}
	tw.trainingsList = mtk.NewList(trainingsParams)
	upButtonBG := graphic.Textures["scrollup.png"]
	if upButtonBG != nil {
		upBG := pixel.NewSprite(upButtonBG,
			upButtonBG.Bounds())
		tw.trainingsList.SetUpButtonBackground(upBG)
	}
	downButtonBG := graphic.Textures["scrolldown.png"]
	if downButtonBG != nil {
		downBG := pixel.NewSprite(downButtonBG,
			downButtonBG.Bounds())
		tw.trainingsList.SetDownButtonBackground(downBG)
	}
	tw.trainingsList.SetOnItemSelectFunc(tw.onTrainingSelected)
	return tw
}

// Draw draws window.
func (tw *TrainingWindow) Draw(win *mtk.Window, matrix pixel.Matrix) {
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
	tw.titleText.Draw(win, matrix.Moved(titleTextMove))
	// Buttons.
	closeButtonMove := pixel.V(tw.Size().X/2-mtk.ConvSize(20),
		tw.Size().Y/2-mtk.ConvSize(15))
	trainButtonMove := mtk.ConvVec(pixel.V(0, -60))
	tw.closeButton.Draw(win, matrix.Moved(closeButtonMove))
	tw.trainButton.Draw(win, matrix.Moved(trainButtonMove))
	// Training info.
	trainingInfoMove := mtk.MoveTC(tw.Size(), tw.trainingInfo.Size())
	trainingInfoMove.Y -= mtk.ConvSize(50)
	tw.trainingInfo.Draw(win, matrix.Moved(trainingInfoMove))
	// Trainings list.
	trainingsMove := mtk.MoveBC(tw.Size(), tw.trainingsList.Size())
	trainingsMove.Y += mtk.ConvSize(20)
	tw.trainingsList.Draw(win, matrix.Moved(trainingsMove))
}

// Update updates window.
func (tw *TrainingWindow) Update(win *mtk.Window) {
	// Elements.
	if tw.Opened() {
		tw.closeButton.Update(win)
		tw.trainButton.Update(win)
		tw.trainingsList.Update(win)
	}
}

// Show shows window.
func (tw *TrainingWindow) Show() {
	tw.opened = true
	if tw.trainer != nil {
		tw.insertTrainings(tw.trainer.Trainings()...)
	}
}

// Hide hides window.
func (tw *TrainingWindow) Hide() {
	tw.opened = false
	tw.trainingsList.Clear()
}

// Opened checks if window is open.
func (tw *TrainingWindow) Opened() bool {
	return tw.opened
}

// DrawArea returns window draw area.
func (tw *TrainingWindow) DrawArea() pixel.Rect {
	return tw.drawArea
}

// Size returns window background size.
func (tw *TrainingWindow) Size() pixel.Vec {
	if tw.bgSpr == nil {
		return mtk.ConvVec(pixel.V(250, 350))
	}
	return mtk.ConvVec(tw.bgSpr.Frame().Size())
}

// SetTrainer sets trainer for window.
func (tw *TrainingWindow) SetTrainer(t training.Trainer) {
	tw.trainer = t
}

// insertTrainings adds all specified trainings to trainings list.
func (tw *TrainingWindow) insertTrainings(trainings ...*training.TrainerTraining) {
	for _, t := range trainings {
		tw.trainingsList.AddItem(lang.Text(t.ID()), t)
	}
}

// Triggered on close button clicked.
func (tw *TrainingWindow) onCloseButtonClicked(b *mtk.Button) {
	tw.Hide()
}

// Triggered after selecting training from list.
func (tw *TrainingWindow) onTrainingSelected(cs *mtk.CheckSlot) {
	// Retrieve training from slot.
	train, ok := cs.Value().(*training.TrainerTraining)
	if !ok {
		log.Err.Printf("hud training: unable to retrieve training from list")
		return
	}
	tw.trainButton.Active(true)
	// Show training info.
	nameInfo := lang.Texts(train.ID())
	if len(nameInfo) > 1 {
		trainingInfo := nameInfo[1]
		for _, r := range train.Requirements() {
			trainingInfo = fmt.Sprintf("%s\n%s", trainingInfo, reqInfo(r))
		}
		tw.trainingInfo.SetText(trainingInfo)
	}
}

// Triggered on train button clicked.
func (tw *TrainingWindow) onTrainButtonClicked(b *mtk.Button) {
	// Retrieve training from list.
	val := tw.trainingsList.SelectedValue()
	if val == nil {
		return
	}
	train, ok := val.(*training.TrainerTraining)
	if !ok {
		log.Err.Printf("hud training: unable to retrieve training from list")
		return
	}
	pc := tw.hud.Game().ActivePlayer()
	pc.Use(train)
}

// reqInfo returns information about specified
// requirement.
func reqInfo(r req.Requirement) string {
	info := ""
	switch r := r.(type) {
	case *req.Item:
		reqLabel := lang.Text("req_item")
		info = fmt.Sprintf("%s: %s x%d", reqLabel, r.ItemID(),
			r.ItemAmount())
	case *req.Currency:
		reqLabel := lang.Text("req_currency")
		info = fmt.Sprintf("%s: %d", reqLabel, r.Amount())
	default:
		return lang.Text("unknown")
	}
	return info
}
