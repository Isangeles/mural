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
	"fmt"

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
func newConsole() (*Console, error) {
	c := new(Console)
	// Text box.
	c.textbox = mtk.NewTextbox(pixel.V(0, 0), mtk.SIZE_MEDIUM, colornames.Grey)
	// Text input.
	c.textedit = mtk.NewTextedit(mtk.SIZE_MEDIUM, colornames.Grey)
	c.textedit.SetOnInputFunc(c.onTexteditInput)
	return c, nil
}

// Draw draws console.
func (c *Console) Draw(win *mtk.Window) {
	drawMin := pixel.V(win.Bounds().Min.X, win.Bounds().Center().Y)
	drawMax := win.Bounds().Max
	// Text box.
	c.textbox.Draw(pixel.R(drawMin.X, drawMin.Y, drawMax.X, drawMax.Y), win)
	// Text edit.
	editSize := pixel.V(c.textbox.DrawArea().Size().X, mtk.ConvSize(30))
	editMove := pixel.V(win.Bounds().Center().X, drawMin.Y-mtk.ConvSize(20))
	c.textedit.SetSize(editSize)
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
	if win.JustPressed(pixelgl.KeyDown) {
		c.textedit.SetText(c.lastInput)
	}
	// Messages.
	var msgs []fmt.Stringer
	engineMsgs := enginelog.Messages()
	/* 
        for i := len(engineMsgs)-1; i >= 0; i-- {
		msgs = append(msgs, engineMsgs[i])
	}
	*/
	for _, msg := range engineMsgs {
		msgs = append(msgs, msg)
	}
	c.textbox.SetMaxTextWidth(win.Bounds().Max.X)
	c.textbox.Insert(msgs)
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
