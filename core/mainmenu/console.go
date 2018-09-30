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
	//"golang.org/x/image/colornames"

	"github.com/isangeles/flame/core/enginelog"

	"github.com/isangeles/mural/core/mtk"
)

// Struct for game console.
// TODO: very slow while open.
type Console struct {
	textbox *mtk.Textbox
	open    bool
}

// newConsole creates game console.
func newConsole() (*Console, error) {
	c := new(Console)
	
	// Text box.
	textbox, err := mtk.NewTextbox()
	if err != nil {
		return nil, err
	}
	c.textbox = textbox
	
	return c, nil
}

// Draw draws console.
func (c *Console) Draw(drawMin, drawMax pixel.Vec, win *pixelgl.Window) {
	c.textbox.Draw(drawMin, drawMax, win)
}

// Update handles key events and updates console.
func (c *Console) Update(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyGraveAccent) {
		if !c.open {
			c.open = true
		} else {
			c.open = false
		}
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
}

// Checks if console is open(should be open).
func (c *Console) Open() bool {
	return c.open
}
