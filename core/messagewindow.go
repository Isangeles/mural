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

package core

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"
)

// MessageWindow struct represents UI message window.
type MessageWindow struct {
	bg        *imdraw.IMDraw
	textbox   *Textbox
	open      bool
	dismissed bool
}

// NewMessageWindow return new message window instance.
func NewMessageWindow(msg string) (*MessageWindow, error) {
	mw := new(MessageWindow)
	mw.bg = imdraw.New(nil)
	textbox, err := NewTextbox()
	if err != nil {
		return nil, err
	}
	mw.textbox = textbox
	tex := []string{msg}
	mw.textbox.InsertText(tex)
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

// Draw draws window.
func (mw *MessageWindow) Draw(topLeft, bottomRight pixel.Vec, win *pixelgl.Window) {
	// Background.
	/*
	mw.bg.Color = pixel.RGB(0.6, 0.6, 0.6)
	mw.bg.Push(topLeft)
	mw.bg.Color = pixel.RGB(0.6, 0.6, 0.6)
	mw.bg.Push(bottomRight)
	mw.bg.Rectangle(0)
	mw.bg.Draw(win)
        */
	// Textbox.
	mw.textbox.Draw(topLeft, bottomRight, win)
}

// Update handles key events.
func (mw *MessageWindow) Update(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyEscape) {
		mw.open = false;
		mw.dismissed = true;
	}

	mw.textbox.Update(win)
}
