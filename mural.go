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
	"github.com/isangeles/flame/cmd/burn"
	"github.com/isangeles/flame/cmd/burn/syntax"
	"github.com/isangeles/flame/core/data/text/lang"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/core/ci"
	"github.com/isangeles/mural/hud"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/mainmenu"
	"github.com/isangeles/mural/objects"
)

const (
	NAME, VERSION = "Mural", "0.0.0"
)

var (
	mainMenu *mainmenu.MainMenu
	pcHUD    *hud.HUD
	game     *flamecore.Game
	inGame   bool

	focus = new(mtk.Focus)
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
	// Check whether Flame module is loaded.
	if flame.Mod() == nil {
		log.Err.Printf("%s\n", lang.Text("gui", "no_mod_loaded_err"))
		return
	}
	// Configure window.
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
	// Load UI data.
	err = data.LoadUIData()
	if err != nil {
		panic(fmt.Errorf("data_load_fail:%v", err))
	}
	uiFont, err := data.Font("SIMSUN.ttf")
	if err != nil {
		panic(fmt.Errorf("fail_to_load_ui_font:%v", err))
	}
	mtk.SetMainFont(uiFont)
	// Create main menu.
	mainMenu, err := mainmenu.New()
	if err != nil {
		panic(err)
	}
	mainMenu.SetOnGameCreatedFunc(EnterGame)
 	err = mainMenu.ImportPlayableChars(flame.Mod().CharactersPath())
	if err != nil {
		log.Err.Printf("init_run:fail_to_import_playable_characters:%v",
			err)
	}
	ci.SetMainMenu(mainMenu)
	mainMenu.Console().SetOnCommandFunc(ExecuteCommand)
	// Debug mode.
	fpsInfo := mtk.NewText("", mtk.SIZE_MEDIUM, 0)
	versionInfo := mtk.NewText(fmt.Sprintf("%s(%s)@%s(%s)", NAME, VERSION,
		flame.NAME, flame.VERSION), mtk.SIZE_MEDIUM, 0)
	versionInfo.JustLeft()
	// Main loop.
	//last := time.Now()
	for !win.Closed() {
		// Delta.
		//dt := time.Since(last).Seconds()
		//last = time.Now()
		// Update.
		if inGame {
			pcHUD.Update(win)
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
				versionInfo.Bounds(), win.Bounds().Min)))
		}
		win.Update()
	}
	// On exit.
	if win.Closed() {
		config.SaveConfig()
		flame.SaveConfig()
	}
}

// EnterGame creates HUD for specified game.
func EnterGame(g *flamecore.Game, pc *objects.Avatar) {
	game = g
	HUD, err := hud.NewHUD(g, pc)
	if err != nil {
		log.Err.Printf("fail_to_create_player_HUD:%v", err)
		return
	}
	pcHUD = HUD
	ci.SetHUD(pcHUD)
	pcHUD.Chat().SetOnCommandFunc(ExecuteCommand)
	inGame = true
}

// ExecuteCommand handles specified text line
// as CI command.
// Returns result code and output text, or error if
// specified line is not valid command.
func ExecuteCommand(line string) (int, string, error) {
	cmd, err := syntax.NewSTDExpression(line)
	if err != nil {
		return -1, "", fmt.Errorf("invalid_input:%s", line)
	}
	res, out := burn.HandleExpression(cmd)
	return res, out, nil
}
