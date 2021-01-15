/*
 * chat.go
 *
 * Copyright 2018-2021 Dariusz Sikora <dev@isangeles.pl>
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
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame/data/res/lang"
	flamelog "github.com/isangeles/flame/log"
	"github.com/isangeles/flame/module/objects"

	"github.com/isangeles/burn"
	"github.com/isangeles/burn/syntax"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res/graphic"
	"github.com/isangeles/mural/log"
)

var (
	chatKey           = pixelgl.KeyGraveAccent
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
	msgs         map[string]*flamelog.Message
	activated    bool
	lastInput    string
	onScriptName func(name string, args ...string) error
}

// Interface for objects with combat log.
type CombatLogger interface {
	objects.Logger
	CombatLog() *objects.Log
}

// Struct for log message.
type Message struct {
	author string
	time   time.Time
	text   string
}

// Struct for sorting messages by the messsage time.
type MessagesByTime []Message

func (mbt MessagesByTime) Len() int           { return len(mbt) }
func (mbt MessagesByTime) Swap(i, j int)      { mbt[i], mbt[j] = mbt[j], mbt[i] }
func (mbt MessagesByTime) Less(i, j int) bool { return mbt[i].time.UnixNano() < mbt[j].time.UnixNano() }

// newChat creates new chat window for HUD.
func newChat(hud *HUD) *Chat {
	c := new(Chat)
	c.hud = hud
	c.msgs = make(map[string]*flamelog.Message)
	// Background.
	c.bgDraw = imdraw.New(nil)
	bg := graphic.Textures["chatbg.png"]
	if bg != nil {
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
	// Key events.
	if win.JustPressed(chatKey) {
		// Toggle chat activity.
		c.Activate(!c.Activated())
	}
	// Clear textbox.
	scrollBottom := c.textbox.AtBottom()
	c.textbox.Clear()
	// Add messages from players and nearby objects.
	messages := make([]Message, 0)
	for _, pc := range c.hud.Game().Players() {
		// PC's private messages.
		for _, m := range pc.PrivateLog().Messages() {
			m := Message{
				author: pc.ID(),
				time:   m.Time(),
				text:   fmt.Sprintf("%s\n", m.String()),
			}
			messages = append(messages, m)
		}
		// Near objects chat & combat.
		area := c.hud.Game().Module().Chapter().CharacterArea(pc.Character)
		if area == nil {
			continue
		}
		for _, tar := range area.NearTargets(pc.Character, pc.SightRange()) {
			log, ok := tar.(objects.Logger)
			if !ok {
				continue
			}
			for _, m := range log.ChatLog().Messages() {
				m := Message{
					author: log.ID(),
					time:   m.Time(),
					text:   fmt.Sprintf("%s\n", m.String()),
				}
				messages = append(messages, m)
			}
			cmbLog := c.combatLogger(log)
			if cmbLog == nil {
				continue
			}
			for _, m := range cmbLog.CombatLog().Messages() {
				m := Message{
					author: log.ID(),
					time:   m.Time(),
					text:   fmt.Sprintf("%s\n", m.String()),
				}
				messages = append(messages, m)
			}
		}
	}
	// Add engine log messages.
	for _, m := range flamelog.Messages() {
		m := Message{
			author: "system",
			time:   m.Date(),
			text:   m.String(),
		}
		messages = append(messages, m)
	}
	// Sort and print messages.
	sort.Sort(MessagesByTime(messages))
	for _, m := range messages {
		c.textbox.AddText(fmt.Sprintf("%s: %s", lang.Text(m.author),
			m.text))
	}
	if scrollBottom {
		c.textbox.ScrollBottom()
	}
	// Elements update.
	c.textbox.Update(win)
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
	c.hud.Camera().Lock(c.Activated())
}

// Echo displays specified text in chat log.
func (c *Chat) Echo(text string) {
	log.Inf.Printf("%s", text)
}

// combatLogger retruns returns object with combat
// log for specified logger, or nil if such object
// does not exists.
func (c *Chat) combatLogger(l objects.Logger) CombatLogger {
	for _, a := range c.hud.Camera().Avatars() {
		if a.ID() == l.ID() && a.Serial() == l.Serial() {
			return a
		}
	}
	return nil
}

// Triggered after accepting input in text edit.
func (c *Chat) onTexteditInput(t *mtk.Textedit) {
	// Save last input.
	input := t.Text()
	c.lastInput = input
	defer t.Clear()
	// Execute command.
	if strings.HasPrefix(input, chatCommandPrefix) {
		cmdInput := strings.TrimPrefix(input, chatCommandPrefix)
		res, out, err := executeCommand(cmdInput)
		if err != nil {
			log.Err.Printf("unable to execute command: %v", err)
		}
		// Echo command result to log.
		log.Cli.Printf("[%d]:%s", res, out)
		return
	}
	// Execute script file.
	if strings.HasPrefix(input, chatScriptPrefix) && c.onScriptName != nil {
		input = strings.TrimPrefix(input, chatScriptPrefix)
		args := strings.Split(input, " ")
		err := c.executeScriptFile(args[0], args...)
		if err != nil {
			log.Err.Printf("unable to execute script: %v", err)
		}
		return
	}
	// Echo chat.
	c.hud.Game().ActivePlayer().AddChatMessage(input)
}

// executeScriptFile executes Ash script from file
// with specified name in background.
func (c *Chat) executeScriptFile(name string, args ...string) error {
	modpath := c.hud.Game().Module().Conf().Path
	path := filepath.FromSlash(modpath + "/gui/scripts/" + name + ".ash")
	script, err := data.Script(path)
	if err != nil {
		return fmt.Errorf("unable to retrieve script: %v", err)
	}
	go c.hud.RunScript(script)
	return nil
}

// executeCommand handles specified text line
// as CI command.
// Returns result code and output text, or error if
// specified line is not valid command.
func executeCommand(line string) (int, string, error) {
	cmd, err := syntax.NewSTDExpression(line)
	if err != nil {
		return -1, "", fmt.Errorf("invalid input: %s", line)
	}
	res, out := burn.HandleExpression(cmd)
	return res, out, nil
}
