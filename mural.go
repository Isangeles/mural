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
	"fmt"
	//"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame"
	flamecore "github.com/isangeles/flame/core"
	"github.com/isangeles/flame/core/module/object/character"
	"github.com/isangeles/flame/core/data/text/lang"

	"github.com/isangeles/mural/config" 
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/mainmenu"
	"github.com/isangeles/mural/hud"
)

const (
	NAME, VERSION = "Mural", "0.0.0"
)

var (
	mainMenu *mainmenu.MainMenu
	pcHUD    *hud.HUD
	game     *flamecore.Game

	focus= new(mtk.Focus)

	inGame bool
)

// On init.
func init() {
	err := flame.LoadConfig()
	if err != nil {
		log.Err.Printf("fail_to_load_flame_config_file:%v\n", err)
		flame.SaveConfig() // override 'corrupted' config file with default configuration
	}
	err = config.LoadConfig()
	if err != nil {
		log.Err.Printf("fail_to_load_config_file:%v\n", err)
	}
}

func main() {
	pixelgl.Run(run)
}

// All window code fired from there.
func run() {
	if flame.Mod() == nil {
		log.Err.Printf("%s\n", lang.Text("gui", "no_mod_loaded_err"))
		return
	}
	resolution := config.Resolution()
	if resolution.X == 0 || resolution.Y == 0 {
		monitor := pixelgl.PrimaryMonitor()
		resolution.X, resolution.Y = monitor.Size()
		//mPosX, mPosY := monitor.Position()
	}
	cfg := pixelgl.WindowConfig{
		Title:  NAME + " " + VERSION,
		Bounds: pixel.R(0, 0, resolution.X, resolution.Y),
		VSync:  true,
	}
	if config.Fullscreen() {
		monitor := pixelgl.PrimaryMonitor()
		cfg.Monitor = monitor
	}

	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	err = data.Load()
	if err != nil {
		panic(fmt.Errorf("data_load_fail:%v", err))
	}
	uiFont, err := data.Font("SIMSUN.ttf")
	if err != nil {
		panic(fmt.Errorf("fail_to_load_ui_font:%v", err))
	}
	mtk.SetMainFont(uiFont)
	
	mainMenu, err := mainmenu.New()
	if err != nil {
		panic(err)
	}
	mainMenu.SetOnGameCreatedFunc(EnterGame)
	fpsInfo := mtk.NewText("", mtk.SIZE_MEDIUM, 0)
	versionInfo := mtk.NewText(fmt.Sprintf("%s(%s)@%s(%s)", NAME, VERSION,
		flame.NAME, flame.VERSION), mtk.SIZE_MEDIUM, 0)
	versionInfo.JustLeft()
	
	// textbox test.
	/*
	for i := 0; i < 40; i ++ {
		log.Dbg.Printf("msg_%d", i)
	}
        */

	//last := time.Now()
	for !win.Closed() {
		// Delta.
		//dt := time.Since(last).Seconds()
		//last = time.Now()
		// Update.
		if inGame {
			// TODO: update HUD.
		} else {
			mainMenu.Update(win)
		}
		fpsInfo.SetText(fmt.Sprintf("FPS:%d", win.FPS()))
		
		// Draw.
		win.Clear(colornames.Black)
		if inGame {
			pcHUD.Draw(win) // is HUD nil check?
		} else {
			mainMenu.Draw(win)
		}
		if config.Debug() {
			fpsInfo.Draw(win, mtk.Matrix().Moved(mtk.PosTR(
				fpsInfo.Bounds(), win.Bounds().Max)))
			versionInfo.Draw(win, mtk.Matrix().Moved(mtk.PosBL(
				fpsInfo.Bounds(), win.Bounds().Min)))
		}

		win.Update()
	}

	// On exit.
	if win.Closed() {
		config.SaveConfig()
		flame.SaveConfig()
	}
}

// EnterGame opens HUD and area map for
// specified game.
func EnterGame(g *flamecore.Game, pc *character.Character) {
	game = g
	HUD, err := hud.NewHUD(g, pc)
	if err != nil {
		log.Err.Printf("fail_to_create_player_HUD:%v", err)
		return
	}
	pcHUD = HUD
	inGame = true
}
