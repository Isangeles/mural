/*
 * console.go
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

package mainmenu

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/isangeles/flame/cmd/command"
	"github.com/isangeles/flame/core/enginelog"

	"github.com/isangeles/mural/core/ci"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/log"
)

// Struct for game console.
type Console struct {
	textbox   *mtk.Textbox
	textedit  *mtk.Textedit
	opened    bool
	lastInput string
}

// newConsole creates game console.
func newConsole() (*Console, error) {
	c := new(Console)
	// Text box.
	c.textbox = mtk.NewTextbox(mtk.SIZE_MEDIUM, colornames.Grey)
	// Text input.
	c.textedit = mtk.NewTextedit(mtk.SIZE_MEDIUM, colornames.Grey, "")
	c.textedit.SetOnInputFunc(c.onTexteditInput)

	return c, nil
}

// Draw draws console.
func (c *Console) Draw(drawMin, drawMax pixel.Vec, win *pixelgl.Window) {
	c.textbox.Draw(pixel.R(drawMin.X, drawMin.Y, drawMax.X, drawMax.Y), win)
	c.textedit.Draw(pixel.R(drawMin.X, drawMin.Y-mtk.ConvSize(20), drawMax.X,
		drawMin.Y), win)
}

// Update handles key events and updates console.
func (c *Console) Update(win *mtk.Window) {
	if win.JustPressed(pixelgl.KeyGraveAccent) {
		if !c.opened {
			c.Show(true)
		} else {
			c.Show(false)
		}
	}
	if win.JustPressed(pixelgl.KeyDown) {
		c.textedit.SetText(c.lastInput)
	}
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
	c.textbox.Insert(msgs)
	
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

// Execute executes specified text command.
func (c *Console) Execute(line string) {
	log.Cli.Printf(">%s", line)
	c.lastInput = line
	cmd, err := command.NewStdCommand(line)
	if err != nil {
		log.Err.Printf("invalid_input:%s", line)
		return
	}
	res, out := ci.HandleCommand(cmd)
	log.Cli.Printf("[%d]:%s", res, out)
}

// Triggered after accept input in text edit.
func (c *Console) onTexteditInput(t *mtk.Textedit) {
	c.Execute(t.Text())
	t.Clear()
}
