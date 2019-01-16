/*
 * chat.go
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

package hud

import (
	"fmt"
	"strings"
	
	"golang.org/x/image/colornames"
	
	"github.com/faiface/pixel"

	"github.com/isangeles/flame/core/enginelog"
	
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/log"
)

var (
	chat_command_prefix = "$"
)

// Chat represents HUD chat window.
type Chat struct {
	hud       *HUD
	textbox   *mtk.Textbox
	textedit  *mtk.Textedit
	activated bool
	onCommand func(line string) (int, string, error)

	textboxSize  pixel.Vec
	texteditSize pixel.Vec
}

// newChat creates new chat window for HUD.
func newChat(hud *HUD) *Chat {
	c := new(Chat)
	c.hud = hud
	c.textbox = mtk.NewTextbox(pixel.V(0, 0), mtk.SIZE_SMALL,
		colornames.Grey)
	c.textedit = mtk.NewTextedit(mtk.SIZE_MEDIUM, colornames.Grey, "")
	c.textedit.SetOnInputFunc(c.onTexteditInput)
	c.textboxSize = pixel.V(mtk.ConvSize(600), mtk.ConvSize(300))
	c.texteditSize = pixel.V(mtk.ConvSize(600), mtk.ConvSize(40))
	return c
}

// Draw draws chat window.
func (c *Chat) Draw(win *mtk.Window) {
	textboxDA := pixel.R(
		win.Bounds().Min.X + mtk.ConvSize(10),
		win.Bounds().Min.Y + c.texteditSize.Y + mtk.ConvSize(10),
		win.Bounds().Min.X + c.textboxSize.X,
		win.Bounds().Min.Y + c.texteditSize.Y + c.textboxSize.Y)
	texteditDA := pixel.R(
		win.Bounds().Min.X + mtk.ConvSize(10),
		win.Bounds().Min.Y + mtk.ConvSize(10),
		win.Bounds().Min.X + c.texteditSize.X,
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
	if c.Activated() {
		c.textedit.Update(win)
	}
}

// Activated checks whether chat input is
// active.
func (c *Chat) Activated() bool {
	return c.activated
}

// Active toggles chat intput activity.
func (c *Chat) Active(active bool) {
	c.activated = active
	c.textedit.Focus(c.Activated())
}

// SetOnCommandFunc sets specified function as
// function triggered on command input.
func (c *Chat) SetOnCommandFunc(f func(line string) (int, string, error)) {
	c.onCommand = f
}

// Echo displays specified text in chat log.
func (c *Chat) Echo(text string) {
	log.Inf.Printf("%s", text)
}

// Triggered after accepting input in text edit.
func (c *Chat) onTexteditInput(t *mtk.Textedit) {
	// Echo input to log.
	input := t.Text()
	c.Echo(input)
	defer t.Clear()
	// Execute command.
	if !strings.HasPrefix(input, chat_command_prefix) ||
		c.onCommand == nil {
		return
	}
	cmdInput := strings.TrimPrefix(input, chat_command_prefix)
	res, out, err := c.onCommand(cmdInput)
	if err != nil {
		log.Err.Printf("fail_to_execute_command:%v", err)
	}
	// Echo command result to log.
	log.Cli.Printf("[%d]:%s", res, out)
}
