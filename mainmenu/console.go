/*
 * console.go
 *
 * Copyright 2018-2019 Dariusz Sikora <dev@isangeles.pl>
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

package mainmenu

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/isangeles/flame/core/enginelog"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/log"
)

// Struct for game console.
type Console struct {
	textbox   *mtk.Textbox
	textedit  *mtk.Textedit
	drawArea  pixel.Rect
	opened    bool
	lastInput string
	onCommand func(cmd string) (int, string, error)
}

// newConsole creates game console.
func newConsole() *Console {
	c := new(Console)
	// Text box.
	c.textbox = mtk.NewTextbox(pixel.V(0, 0), mtk.SIZE_MEDIUM, colornames.Grey)
	// Text input.
	c.textedit = mtk.NewTextedit(mtk.SIZE_MEDIUM, colornames.Grey)
	c.textedit.SetOnInputFunc(c.onTexteditInput)
	return c
}

// Draw draws console.
func (c *Console) Draw(win *mtk.Window) {
	// Text box.
	boxSize := pixel.V(win.Bounds().W(), win.Bounds().H()/2)
	c.textbox.SetSize(boxSize)
	boxMove := mtk.DrawPosTC(win.Bounds(), c.textbox.Size())
	c.textbox.Draw(win, mtk.Matrix().Moved(boxMove))
	// Text edit.
	editSize := pixel.V(c.textbox.Size().X, mtk.ConvSize(30))
	c.textedit.SetSize(editSize)
	editMove := mtk.BottomOf(c.textbox.DrawArea(), c.textedit.Size(), 0)
	c.textedit.Draw(win.Window, mtk.Matrix().Moved(editMove))
}

// Update handles key events and updates console.
func (c *Console) Update(win *mtk.Window) {
	// Key events.
	if win.JustPressed(pixelgl.KeyGraveAccent) {
		if !c.opened {
			c.Show(true)
		} else {
			c.Show(false)
		}
		defer c.textedit.Clear()
	}
	if win.JustPressed(pixelgl.KeyUp) {
		c.textedit.SetText(c.lastInput)
	}
	// Messages.
	c.textbox.Clear()
	for _, msg := range enginelog.Messages() {
		c.textbox.AddText(msg.String())
	}
	c.textbox.SetMaxTextWidth(win.Bounds().W())
	// Elements.
	c.textbox.Update(win)
	c.textedit.Update(win)
}

// Show toggles console visibility.
func (c *Console) Show(show bool) {
	c.opened = show
	c.textedit.Focus(show)
	c.textedit.Clear()
}

// Checks if console is open.
func (c *Console) Opened() bool {
	return c.opened
}

// SetOnCommandFunc sets specified function as
// function triggered on command input.
func (c *Console) SetOnCommandFunc(f func(cmd string) (int, string, error)) {
	c.onCommand = f
}

// Echo prints specified text to console.
func (c *Console) Echo(text string) {
	log.Cli.Printf(">%s", text)
}

// Triggered after accepting input in text edit.
func (c *Console) onTexteditInput(t *mtk.Textedit) {
	// Echo input to log.
	input := t.Text()
	c.Echo(input)
	c.lastInput = input
	defer t.Clear()
	// Execute command.
	if c.onCommand == nil {
		return
	}
	res, out, err := c.onCommand(input)
	if err != nil {
		log.Err.Printf("fail_to_execute_command:%s", input)
	}
	// Echo command result to log.
	log.Cli.Printf("[%d]:%s", res, out)
}
