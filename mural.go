/*
 * mural.go
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

// Mural is 2D graphical frontend for Flame engine.
package main

import (
	"fmt"

	"golang.org/x/image/colornames"

	"github.com/faiface/beep"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/cmd/burn"
	"github.com/isangeles/flame/cmd/burn/syntax"
	flameconf "github.com/isangeles/flame/config"
	flamecore "github.com/isangeles/flame/core"
	flamedata "github.com/isangeles/flame/core/data"
	flamesave "github.com/isangeles/flame/core/data/save"
	"github.com/isangeles/flame/core/data/text/lang"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/ci"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/imp"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/hud"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/mainmenu"
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
	// Load flame config.
	err := flameconf.LoadConfig()
	if err != nil {
		log.Err.Printf("fail_to_load_flame_config_file:%v\n", err)
		flameconf.SaveConfig() // override 'corrupted' config file with default configuration
	}
	// Load module.
	m, err := flamedata.Module(flameconf.ModulePath(), flameconf.LangID())
	if err != nil {
		log.Err.Printf("fail_to_load_config_module:%v", err)
	}
	flame.SetModule(m)
	// Load GUI config.
	err = config.LoadConfig()
	if err != nil {
		log.Err.Printf("fail_to_load_config_file:%v\n", err)
	}
}

// Main function.
func main() {
	// Check whether Flame module is loaded.
	if flame.Mod() == nil {
		panic(fmt.Sprintf("%s\n", lang.Text("gui", "no_mod_loaded_err")))
	}
	// Load UI graphic.
	err := data.LoadUIData()
	if err != nil {
		panic(fmt.Errorf("data_load_fail:%v", err))
	}
	// Load game graphic.
	data.LoadGameData()
	// Load module data.
	err = flamedata.LoadModuleData(flame.Mod())
	if err != nil {
		panic(fmt.Errorf("fail_to_load_module_data:%v", err))
	}
	// Load module graphic data.
	err = imp.LoadModuleResources(flame.Mod())
	if err != nil {
		panic(fmt.Errorf("main_menu:load_game:fail_to_load_resources:%v", err))
	}
	// Music.
	mtk.InitAudio(beep.Format{44100, 2, 2})
	if mtk.Audio() != nil {
		ci.SetMusicPlayer(mtk.Audio())
		m, err := data.Music(config.MenuMusicFile())
		if err != nil {
			log.Err.Printf("fail_to_load_main_theme_audio_data:%v", err)
		} else {
			mtk.Audio().AddMusic(m)
		}
		mtk.Audio().PlayMusic()
	}
	// Graphic.
	pixelgl.Run(run)
}

// All window code fired from there.
func run() {
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
	// UI Font.
	uiFont, err := data.Font(config.MainFontName())
	if err == nil { // if font from config was found
		mtk.SetMainFont(uiFont)
	}
	// Audio effects.
	bClickSound1, err := data.AudioEffect(config.ButtonClickSoundFile())
	if err != nil {
		log.Err.Printf("init_run:fail_to_retrieve_button_click_audio_data:%v",
			err)
	}
	mtk.SetButtonClickSound(bClickSound1) // global button click sound
	// Create main menu.
	mainMenu, err = mainmenu.New()
	if err != nil {
		panic(err)
	}
	mainMenu.SetOnGameCreatedFunc(EnterGame)
	mainMenu.SetOnSaveImportedFunc(EnterSavedGame)
	err = mainMenu.ImportPlayableChars(flame.Mod().Conf().CharactersPath())
	if err != nil {
		log.Err.Printf("init_run:fail_to_import_playable_characters:%v",
			err)
	}
	ci.SetMainMenu(mainMenu)
	mainMenu.Console().SetOnCommandFunc(ExecuteCommand)
	// Debug mode.
	fpsInfo := mtk.NewText(mtk.SIZE_MEDIUM, 0)
	fpsInfo.JustRight()
	versionInfo := mtk.NewText(mtk.SIZE_MEDIUM, 0)
	versionInfo.SetText(fmt.Sprintf("%s(%s)@%s(%s)", NAME, VERSION,
		flame.NAME, flame.VERSION))
	versionInfo.JustRight()
	// Main loop.
	for !win.Closed() {
		// Draw.
		win.Clear(colornames.Black)
		if inGame {
			pcHUD.Draw(win)
		} else {
			mainMenu.Draw(win)
		}
		if config.Debug() {
			fpsInfo.Draw(win, mtk.Matrix().Moved(mtk.PosTR(
				fpsInfo.Bounds(), win.Bounds().Max)))
			versionInfo.Draw(win, mtk.Matrix().Moved(mtk.LeftOf(
				fpsInfo.DrawArea(), versionInfo.Bounds(), 5)))
		}
		// Update.
		win.Update()
		if inGame {
			pcHUD.Update(win)
			game.Update(win.Delta()) // game update
			if pcHUD.Exiting() {
				inGame = false
			}
		} else {
			mainMenu.Update(win)
		}
		fpsInfo.SetText(fmt.Sprintf("FPS:%d", win.FPS()))
	}
	// On exit.
	if win.Closed() {
		config.SaveConfig()
		flameconf.SaveConfig()
	}
}

// EnterGame creates HUD for specified game.
func EnterGame(g *flamecore.Game) {
	mainMenu.OpenLoadingScreen(lang.Text("gui", "enter_game_info"))
	defer mainMenu.CloseLoadingScreen()
	game = g
	HUD, err := hud.NewHUD(game)
	if err != nil {
		log.Err.Printf("fail_to_create_player_HUD:%v", err)
		return
	}
	setHUD(HUD)
	inGame = true
}

// EnterSavedGame creates game and HUD from saved data.
func EnterSavedGame(save *flamesave.SaveGame) {
	mainMenu.OpenLoadingScreen(lang.Text("gui", "loadgame_load_game_info"))
	defer mainMenu.CloseLoadingScreen()
	// Load game.
	game = flamecore.LoadGame(save)
	flame.SetGame(game)
	// Import saved GUI state.
	guisav, err := imp.ImportGUISave(flameconf.ModuleSavegamesPath(), save.Name)
	if err == nil {
		// Add avatars from saved GUI to global avatars data resources.
		for _, pcd := range guisav.PlayersData {
			res.AddAvatarData(pcd.Avatar)
		}
	} else {
		log.Err.Printf("fail_to_load_gui_save:%v", err)
	}
	// Create HUD.
	HUD, err := hud.NewHUD(game)
	if err != nil {
		log.Err.Printf("fail_to_create_player_HUD:%v", err)
		msg := lang.Text("gui", "load_game_err")
		mainMenu.ShowMessage(msg)
		return
	}
	// Load HUD state.
	HUD.LoadGUISave(guisav)
	// Show HUD.
	setHUD(HUD)
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

// setHUD sets specified HUD instance as current GUI player
// HUD.
func setHUD(h *hud.HUD) {
	pcHUD = h
	ci.SetHUD(pcHUD)
	pcHUD.Chat().SetOnCommandFunc(ExecuteCommand)
}
