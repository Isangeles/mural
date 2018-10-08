/*
 * messagewindow.go
 *
 * Copyright 2018 Dariusz Sikora <dev@isangeles.pl>
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

package mtk

import (
	"image/color"
	
	"golang.org/x/image/colornames"
	
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"

	"github.com/isangeles/flame/core/data/text/lang"
)

// MessageWindow struct represents UI message window.
type MessageWindow struct {
	bg           *imdraw.IMDraw
	drawArea     pixel.Rect
	size         Size
	textbox      *Textbox
	acceptButton *Button
	cancelButton *Button
	open         bool
	accepted     bool
	dismissed    bool
	onAccept     func(msg *MessageWindow)
	onCancel     func(msg *MessageWindow)
}

// NewMessageWindow creates new message window instance.
func NewMessageWindow(size Size, msg string) (*MessageWindow, error) {
	mw := new(MessageWindow)
	// Background.
	mw.bg = imdraw.New(nil)
	mw.size = size
	// Textbox.
	textbox, err := NewTextbox()
	if err != nil {
		return nil, err
	}
	mw.textbox = textbox
	tex := []string{msg}
	mw.textbox.InsertText(tex)
	// Buttons.
	acceptB := NewButton(SIZE_SMALL, SHAPE_RECTANGLE, colornames.Red,
		lang.Text("gui", "accept_b_label"))
	acceptB.SetOnClickFunc(mw.onAcceptButtonClicked)
	mw.acceptButton = acceptB

	return mw, nil
}

// NewDialogWindow creates new dialog window with message.
func NewDialogWindow(size Size, msg string) (*MessageWindow, error) {
	// Basic message window.
	mw, err := NewMessageWindow(size, msg)
	if err != nil {
		return nil, err
	}
	// Buttons.
	mw.cancelButton = NewButton(SIZE_SMALL, SHAPE_RECTANGLE, colornames.Red,
		lang.Text("gui", "cancel_b_label"))
	mw.cancelButton.SetOnClickFunc(mw.onCancelButtonClicked)
	
	return mw, nil
}

// Show toggles window visibility.
func (mw *MessageWindow) Show(open bool) {
	mw.open = open
}

// Open checks if window should be open.
func (mw *MessageWindow) Open() bool {
	return mw.open
}

// Dismissed checks if window was dismised.
func (mw *MessageWindow) Dismissed() bool {
	return mw.dismissed
}

// Accepted checks if message was accepted.
func (mw *MessageWindow) Accepted() bool {
	return mw.accepted
}

// Draw draws window.
func (mw *MessageWindow) Draw(t pixel.Target, matrix pixel.Matrix) {
	// Calculating draw area.
	mw.drawArea = MatrixToDrawArea(matrix, mw.Frame())
	// Background.
	mw.drawIMBackground(t, colornames.Grey)
	// Buttons.
	mw.acceptButton.Draw(t, pixel.IM.Moved(PosBR(mw.acceptButton.Frame(),
		pixel.V(mw.drawArea.Max.X, mw.drawArea.Min.Y))))
	if mw.cancelButton != nil {
		mw.cancelButton.Draw(t, pixel.IM.Moved(PosBL(mw.acceptButton.Frame(),
		mw.drawArea.Min)))
	}
	// Textbox.
	mw.textbox.Draw(pixel.V(mw.drawArea.Min.X, mw.acceptButton.DrawArea().Max.Y),
		mw.drawArea.Max, t)
}

// drawIMBackround Draws IMDraw background.
func (mw *MessageWindow) drawIMBackground(t pixel.Target, color color.Color) {
	mw.bg.Color = pixel.ToRGBA(color)
	mw.bg.Push(mw.drawArea.Min)
	mw.bg.Color = pixel.ToRGBA(color)
	mw.bg.Push(mw.drawArea.Max)
	mw.bg.Rectangle(0)
	mw.bg.Draw(t)
}

// Update handles key press events.
func (mw *MessageWindow) Update(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyEscape) {
		if mw.onCancel != nil {
			mw.cancel()
		}
	}

	mw.textbox.Update(win)
	mw.acceptButton.Update(win)
	if mw.cancelButton != nil {
		mw.cancelButton.Update(win)
	}
}

// Frame resturns message window size bounds.
func (mw *MessageWindow) Frame() pixel.Rect {
	return mw.size.MessageWindowSize()
}

// SetOnAcceptFunc sets specified function as function triggered after
// message was accepted.
func (mw *MessageWindow) SetOnAcceptFunc(f func(msg *MessageWindow)) {
	mw.onAccept = f;
}

// SetOnCancelFunc sets specified function as function triggered after
// message was canceled.
func (mw *MessageWindow) SetOnCancelFunc(f func(msg *MessageWindow)) {
	mw.onCancel = f
}

// Triggered after accept button clicked.
func (mw *MessageWindow) onAcceptButtonClicked(b *Button) {
	mw.accept()
}

// Triggered after cancel button clicked.
func (mw *MessageWindow) onCancelButtonClicked(b *Button) {
	mw.cancel()
}

// accept sets message as accepted.
func (mw *MessageWindow) accept() {
	mw.open = false
	mw.dismissed = true
	mw.accepted = true
	if mw.onAccept != nil {
		mw.onAccept(mw)
	}
}

// cancel sets message as canceld.
func (mw *MessageWindow) cancel() {
	mw.open = false
	mw.dismissed = true
	mw.accepted = false
	if mw.onCancel != nil {
		mw.onCancel(mw)
	}
}
