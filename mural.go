/*
 * mural.go
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

// Mural is 2D graphical frontend for Flame engine.
package main

import (
	"log"
	//"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/isangeles/mural/core/mainmenu"

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/core/data/text/lang"
)

const (
	NAME, VERSION = "Mural", "0.0.0"
	log_prefix    = "mural-gui"
)

var (
	inflog *log.Logger = log.New(flame.InfLog, "mural>", 0)
	errlog *log.Logger = log.New(flame.ErrLog, "mural>", 0)
	dbglog *log.Logger = log.New(flame.DbgLog, "mural-debug>", 0)
)

// On init.
func init() {
	err := LoadConfig()
	if err != nil {
		errlog.Printf("fail_to_load_config_file:%v\n", err)
	}
}

func main() {
	pixelgl.Run(run)
}

// All window code fired from there.
func run() {
	if flame.Mod() == nil {
		errlog.Printf("%s\n", lang.Text("gui", "no_mod_loaded_err"))
		return
	}
	monitor := pixelgl.PrimaryMonitor()
	mResW, mResH := monitor.Size()
	//mPosX, mPosY := monitor.Position()
	cfg := pixelgl.WindowConfig{
		Title:  "Mural GUI",
		Bounds: pixel.R(0, 0, mResW, mResH),
		VSync:  true,
	}
	if fullscreen {
		cfg.Monitor = monitor
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)

	mainMenu, err := mainmenu.New()
	if err != nil {
		panic(err)
	}

	//last := time.Now()
	for !win.Closed() {
		//dt := time.Since(last).Seconds()
		//last = time.Now()

		mainMenu.Update(win)

		win.Clear(colornames.Black)
		mainMenu.Draw(win)

		win.Update()
	}

	// On exit.
	if win.Closed() {
		SaveConfig()
		flame.SaveConfig()
	}
}
