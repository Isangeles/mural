/*
 * chat.go
 *
 * Copyright 2018-2020 Dariusz Sikora <dev@isangeles.pl>
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

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"

	"github.com/isangeles/flame/core/enginelog"
	"github.com/isangeles/flame/core/module/objects"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/log"
)

var (
	chatCommandPrefix = "$"
	chatScriptPrefix  = "%"
)

// Chat represents HUD chat window.
type Chat struct {
	hud          *HUD
	bgSpr        *pixel.Sprite
	bgDraw       *imdraw.IMDraw
	drawArea     pixel.Rect
	textbox      *mtk.Textbox
	textedit     *mtk.Textedit
	msgs         map[string]*enginelog.Message
	activated    bool
	lastInput    string
	onCommand    func(line string) (int, string, error)
	onScriptName func(name string, args ...string) error
}

// newChat creates new chat window for HUD.
func newChat(hud *HUD) *Chat {
	c := new(Chat)
	c.hud = hud
	c.msgs = make(map[string]*enginelog.Message)
	// Background.
	c.bgDraw = imdraw.New(nil)
	bg, err := data.PictureUI("chatbg.png")
	if err == nil { // fallback
		c.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	}
	// Textbox.
	boxSize := c.Size()
	boxSize.Y -= mtk.ConvSize(70)
	textboxParams := mtk.Params{
		SizeRaw:     boxSize,
		FontSize:    mtk.SizeMedium,
		MainColor:   mainColor,
		AccentColor: accentColor,
	}
	c.textbox = mtk.NewTextbox(textboxParams)
	// Textedit.
	c.textedit = mtk.NewTextedit(mtk.SizeMedium, mainColor)
	c.textedit.SetOnInputFunc(c.onTexteditInput)
	return c
}

// Draw draws chat window.
func (c *Chat) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Background.
	if c.bgSpr != nil {
		c.bgSpr.Draw(win, matrix)
	}
	// Textbox.
	c.textbox.Draw(win, matrix)
	// Textedit.
	editSize := pixel.V(c.Size().X, mtk.ConvSize(30))
	c.textedit.SetSize(editSize)
	editMove := pixel.V(0, -c.Size().Y/2+mtk.ConvSize(10))
	c.textedit.Draw(win, matrix.Moved(editMove))
}

// Update updates chat window.
func (c *Chat) Update(win *mtk.Window) {
	// Print log messages.
	for _, msg := range enginelog.Messages() {
		if c.msgs[msg.ID()] != nil {
			continue
		}
		c.msgs[msg.ID()] = &msg
		c.textbox.AddText(msg.String())
	}
	// Print messages from players and nearby objects.
	game := c.hud.Game()
	for _, pc := range game.Players() {
		select {
		case msg := <-pc.PrivateLog():
			c.textbox.AddText(fmt.Sprintf("%s\n", msg))
		default:
		}
		// Near objects.
		area := game.Module().Chapter().CharacterArea(pc)
		if area == nil {
			continue
		}
		for _, tar := range area.NearTargets(pc, pc.SightRange()) {
			tar, ok := tar.(objects.Logger)
			if !ok {
				continue
			}
			select {
			case msg := <-tar.CombatLog():
				c.textbox.AddText(fmt.Sprintf("%s\n", msg))
			case msg := <-tar.ChatLog():
				c.textbox.AddText(fmt.Sprintf("%s: %s\n", tar.Name(), msg))
			case msg := <-tar.PrivateLog():
				if tar == pc {
					c.textbox.AddText(fmt.Sprintf("%s\n", msg))
				}
			default:
			}
		}
	}
	// Elements update.
	c.textbox.Update(win)
	c.Activate(c.textedit.Focused())
	if c.Activated() {
		c.textedit.Update(win)
	}
}

// DrawArea returns current chat draw area.
func (c *Chat) DrawArea() pixel.Rect {
	return c.drawArea
}

// Size returns chat background size.
func (c *Chat) Size() pixel.Vec {
	if c.bgSpr == nil {
		// TODO: return draw background bounds.
		return mtk.ConvVec(pixel.V(0, 0))
	}
	return mtk.ConvVec(c.bgSpr.Frame().Size())
}

// Activated checks whether chat input is
// active.
func (c *Chat) Activated() bool {
	return c.activated
}

// Active toggles chat intput activity.
func (c *Chat) Activate(active bool) {
	c.activated = active
	c.textedit.Focus(c.Activated())
}

// SetOnCommandFunc sets specified function as
// function triggered on command input.
func (c *Chat) SetOnCommandFunc(f func(line string) (int, string, error)) {
	c.onCommand = f
}

// SetOnScriptNameFunc sets specified function as
// function triggered after script name input.
func (c *Chat) SetOnScriptNameFunc(f func(name string, args ...string) error) {
	c.onScriptName = f
}

// Echo displays specified text in chat log.
func (c *Chat) Echo(text string) {
	log.Inf.Printf("%s", text)
}

// Triggered after accepting input in text edit.
func (c *Chat) onTexteditInput(t *mtk.Textedit) {
	// Save last input.
	input := t.Text()
	c.lastInput = input
	defer t.Clear()
	// Execute command.
	if strings.HasPrefix(input, chatCommandPrefix) && c.onCommand != nil {
		cmdInput := strings.TrimPrefix(input, chatCommandPrefix)
		res, out, err := c.onCommand(cmdInput)
		if err != nil {
			log.Err.Printf("fail_to_execute_command:%v", err)
		}
		// Echo command result to log.
		log.Cli.Printf("[%d]:%s", res, out)
		return
	}
	// Execute script file.
	if strings.HasPrefix(input, chatScriptPrefix) && c.onScriptName != nil {
		input = strings.TrimPrefix(input, chatScriptPrefix)
		args := strings.Split(input, " ")
		err := c.onScriptName(args[0], args...)
		if err != nil {
			log.Err.Printf("fail_to_execute_script:%v", err)
		}
		return
	}
	// Echo chat.
	c.hud.ActivePlayer().SendChat(input)
}
