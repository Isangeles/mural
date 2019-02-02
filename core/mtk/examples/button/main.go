/*
 * main.go
 *
 * Copyright 2019 Dariusz Sikora <dev@isangeles.pl>
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

// Example for creating simple MTK button with draw background
// and custom on-click function.
package main

import (
	"fmt"
	
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/mural/core/mtk"
)

var (
	exitreq bool // menu exit request flag
)

// Main function.
func main() {
	// Run Pixel graphic.
	pixelgl.Run(run)
}

// All window code fired from there.
func run() {
	// Create Pixel window configuration.
	cfg := pixelgl.WindowConfig{
		Title:  "MTK menu example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("fail_to_create_mtk_window:%v", err))
	}
	// Create MTK button for exit.
	exitB := mtk.NewButton(mtk.SIZE_BIG, mtk.SHAPE_RECTANGLE, colornames.Red,
		"Exit", "Exit menu")
	// Set function for exit button click event.
	exitB.SetOnClickFunc(onExitButtonClicked)
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Create exit button position.
		exitBPos := mtk.DrawPosC(win.Bounds(), exitB.Frame())
		// Draw exit button.
		exitB.Draw(win, mtk.Matrix().Moved(exitBPos))
		// Update.
		win.Update()
		exitB.Update(win)
		// On exit request.
		if exitreq {
			win.SetClosed(true)
		}
	}
}

// onExitButtonClicked handles exit button click
// event.
func onExitButtonClicked(b *mtk.Button) {
	exitreq = true
}
