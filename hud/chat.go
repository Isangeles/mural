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
	"golang.org/x/image/colornames"
	
	"github.com/faiface/pixel"
	
	"github.com/isangeles/mural/core/mtk"
	//"github.com/isangeles/mural/log"
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
	c.texteditSize = pixel.V(mtk.ConvSize(600), mtk.ConvSize(20))
	return c
}

// Draw draws chat window.
func (c *Chat) Draw(win *mtk.Window) {
	textboxDA := pixel.R(win.Bounds().Max.X - c.textboxSize.X,
		win.Bounds().Min.Y + c.texteditSize.Y, win.Bounds().Max.X - mtk.ConvSize(10),
		win.Bounds().Min.Y + c.texteditSize.Y + c.textboxSize.Y)
	texteditDA := pixel.R(win.Bounds().Max.X - c.texteditSize.X,
		win.Bounds().Min.Y, win.Bounds().Max.X - mtk.ConvSize(10),
		win.Bounds().Min.Y + c.texteditSize.Y)
	c.textbox.Draw(textboxDA, win)
	c.textedit.Draw(texteditDA, win)
}

// Update updates chat window.
func (c *Chat) Update(win *mtk.Window) {
	c.textbox.Update(win)
	if c.Active() {
		c.textedit.Update(win)
	}
	// TODO: fill textbox with engine messages.
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

// Triggered after accepting input in text edit.
func (c *Chat) onTexteditInput(t *mtk.Textedit) {
	// TODO: send chat input to CI.
}
