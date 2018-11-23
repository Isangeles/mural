/*
 * chat.go
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

package hud

import (
	"fmt"
	"strings"
	
	"golang.org/x/image/colornames"
	
	"github.com/faiface/pixel"

	"github.com/isangeles/flame/core/enginelog"
	"github.com/isangeles/flame/cmd/command"

	"github.com/isangeles/mural/core/ci"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/log"
)

var (
	chat_command_prefix = "$"
)

// Chat represents HUD chat window.
type Chat struct {
	hud      *HUD
	textbox  *mtk.Textbox
	textedit *mtk.Textedit
	active   bool

	textboxSize  pixel.Vec
	texteditSize pixel.Vec
}

// newChat creates new chat window for HUD.
func newChat(hud *HUD) *Chat {
	c := new(Chat)
	c.hud = hud
	c.textbox = mtk.NewTextbox(mtk.SIZE_MEDIUM, colornames.Grey)
	c.textedit = mtk.NewTextedit(mtk.SIZE_MEDIUM, colornames.Grey, "")
	c.textedit.SetOnInputFunc(c.onTexteditInput)
	c.textboxSize = pixel.V(mtk.ConvSize(600), mtk.ConvSize(300))
	c.texteditSize = pixel.V(mtk.ConvSize(600), mtk.ConvSize(40))
	return c
}

// Draw draws chat window.
func (c *Chat) Draw(win *mtk.Window) {
	textboxDA := pixel.R(win.Bounds().Max.X - c.textboxSize.X,
		win.Bounds().Min.Y + c.texteditSize.Y + mtk.ConvSize(10),
		win.Bounds().Max.X - mtk.ConvSize(10),
		win.Bounds().Min.Y + c.texteditSize.Y + c.textboxSize.Y)
	texteditDA := pixel.R(win.Bounds().Max.X - c.texteditSize.X,
		win.Bounds().Min.Y + mtk.ConvSize(10),
		win.Bounds().Max.X - mtk.ConvSize(10),
		win.Bounds().Min.Y + c.texteditSize.Y)
	c.textbox.Draw(textboxDA, win)
	c.textedit.Draw(texteditDA, win)
}

// Update updates chat window.
func (c *Chat) Update(win *mtk.Window) {
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
	if c.Active() {
		c.textedit.Update(win)
	}
}

// Active checks whether chat input is
// active.
func (c *Chat) Active() bool {
	return c.active
}

// SetActive toggles chat intput activity.
func (c *Chat) SetActive(active bool) {
	c.active = active
	c.textedit.Focus(active)
}

// Execute executes specified text command.
func (c *Chat) ExecuteCommand(line string) {
	log.Cli.Printf(">%s", line)
	cmd, err := command.NewStdCommand(line)
	if err != nil {
		log.Err.Printf("invalid_input:%s", line)
		return
	}
	res, out := ci.HandleCommand(cmd)
	log.Cli.Printf("[%d]:%s", res, out)
}

// Echo displays specified text in chat log.
func (c *Chat) Echo(text string) {
	log.Inf.Printf("%s", text)
}

// Triggered after accepting input in text edit.
func (c *Chat) onTexteditInput(t *mtk.Textedit) {
	if strings.HasPrefix(t.Text(), chat_command_prefix) {
		c.ExecuteCommand(strings.TrimPrefix(t.Text(), chat_command_prefix))
	} else {
		c.Echo(t.Text())
	}
	t.Clear()
}
