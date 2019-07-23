/*
 * mural.go
 *
 * Copyright 2018-2019 Dariusz Sikora <dev@isangeles.pl>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either Version 2 of the License, or
 * (at your option) any later Version.
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
	"path/filepath"

	"golang.org/x/image/colornames"

	"github.com/faiface/beep"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame"
	flameconf "github.com/isangeles/flame/config"
	flamecore "github.com/isangeles/flame/core"
	flamedata "github.com/isangeles/flame/core/data"
	"github.com/isangeles/flame/core/data/text/lang"

	"github.com/isangeles/burn"
	"github.com/isangeles/burn/syntax"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/ci"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/imp"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/hud"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/mainmenu"
)

const (
	Name, Version = "Mural", "0.0.0"
)

var (
	mainMenu *mainmenu.MainMenu
	pcHUD    *hud.HUD
	game     *flamecore.Game
	inGame   bool
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
	// Check if Flame module is loaded.
	if flame.Mod() == nil {
		panic(fmt.Sprintf("%s\n", lang.Text("gui", "no_mod_loaded_err")))
	}
	// Load UI graphic.
	err := data.LoadUIData()
	if err != nil {
		panic(fmt.Errorf("fail_to_load_gui_data:%v", err))
	}
	// Load game graphic.
	err = data.LoadGameData()
	if err != nil {
		panic(fmt.Errorf("fail_to_load_game_graphic_data:%v", err))
	}
	// Load module data.
	err = flamedata.LoadModuleData(flame.Mod())
	if err != nil {
		panic(fmt.Errorf("fail_to_load_module_data:%v", err))
	}
	// Load module graphic data.
	err = imp.LoadModuleResources(flame.Mod())
	if err != nil {
		panic(fmt.Errorf("fail_to_load_module_resources:%v", err))
	}
	// Music.
	mtk.InitAudio(beep.Format{44100, 2, 2})
	if mtk.Audio() != nil {
		ci.SetMusicPlayer(mtk.Audio())
		m, err := data.Music(config.MenuMusicFile())
		if err != nil {
			log.Err.Printf("fail_to_load_main_theme_audio_data:%v", err)
		} else {
			mtk.Audio().AddAudio(m)
		}
		mtk.Audio().ResumePlaylist()
	}
	// Graphic.
	pixelgl.Run(run)
}

// All window code fired from there.
func run() {
	// Configure window.
	monitor := pixelgl.PrimaryMonitor()
	winPosX, winPosY := 0.0, 0.0
	winRes := config.Resolution()
	if winRes.X == 0 || winRes.Y == 0 {
		winRes.X, winRes.Y = monitor.Size()
	}
	cfg := pixelgl.WindowConfig{
		Title:  Name + " " + Version,
		Bounds: pixel.R(winPosX, winPosY, winRes.X, winRes.Y),
		VSync:  true,
	}
	if config.Fullscreen() {
		monitor := pixelgl.PrimaryMonitor()
		cfg.Monitor = monitor
	}
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("fail_to_create_mtk_window:%v", err))
	}
	// UI Font.
	uiFont, err := data.Font(config.MainFontName())
	if err == nil { // if font from config was found
		mtk.SetMainFont(uiFont)
	}
	// Audio effects.
	bClickSound, err := data.AudioEffect(config.ButtonClickSoundFile())
	if err != nil {
		log.Err.Printf("init_run:fail_to_retrieve_button_click_audio_data:%v",
			err)
	}
	mtk.SetButtonClickSound(bClickSound) // global button click sound
	// Create main menu.
	mainMenu = mainmenu.New()
	mainMenu.SetOnGameCreatedFunc(EnterGame)
	mainMenu.SetOnSaveLoadedFunc(EnterSavedGame)
	err = mainMenu.ImportPlayableChars(flame.Mod().Conf().CharactersPath())
	if err != nil {
		log.Err.Printf("init_run:fail_to_import_playable_characters:%v",
			err)
	}
	ci.SetMainMenu(mainMenu)
	mainMenu.Console().SetOnCommandFunc(ExecuteCommand)
	// Debug mode.
	fpsInfo := mtk.NewText(mtk.SizeMedium, 0)
	fpsInfo.JustRight()
	verInfo := mtk.NewText(mtk.SizeMedium, 0)
	verInfo.SetText(fmt.Sprintf("%s(%s)@%s(%s)", Name, Version,
		flame.NAME, flame.VERSION))
	verInfo.JustRight()
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
			fpsPos := mtk.DrawPosTR(win.Bounds(), fpsInfo.Size())
			fpsPos.Y -= mtk.ConvSize(10)
			fpsInfo.Draw(win, mtk.Matrix().Moved(fpsPos))
			verPos := mtk.LeftOf(fpsInfo.DrawArea(), verInfo.Size(), 5)
			verInfo.Draw(win, mtk.Matrix().Moved(verPos))
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
	// Create HUD.
	hud := hud.New()
	// Set game for HUD.
	err := imp.LoadChapterResources(game.Module().Chapter())
	if err != nil {
		log.Err.Printf("enter_game:fail_to_load_chapter_resources:%v", err)
		mainMenu.ShowMessage(lang.Text("gui", "load_game_err"))
		return
	}
	err = hud.SetGame(game)
	if err != nil {
		log.Err.Printf("enter_game:fail_to_set_hud_game:%v", err)
		mainMenu.ShowMessage(lang.Text("gui", "load_game_err"))
		return
	}
	// Set HUD.
	setHUD(hud)
	inGame = true
}

// EnterSavedGame creates game and HUD from saved data.
func EnterSavedGame(g *flamecore.Game, saveName string) {
	mainMenu.OpenLoadingScreen(lang.Text("gui", "loadgame_load_game_info"))
	defer mainMenu.CloseLoadingScreen()
	game = g
	// Import saved GUI state.
	guisav, err := imp.ImportGUISave(flameconf.ModuleSavegamesPath(), saveName)
	if err != nil {
		log.Err.Printf("enter_saved_game:fail_to_load_gui_save:%v", err)
		mainMenu.ShowMessage(lang.Text("gui", "load_game_err"))
		return
	}
	for _, pcd := range guisav.PlayersData {
		res.AddAvatarData(pcd.Avatar)
	}
	// Create HUD.
	hud := hud.New()
	err = imp.LoadChapterResources(game.Module().Chapter())
	err = hud.SetGame(game)
	if err != nil {
		log.Err.Printf("enter_saved_game:fail_to_set_hud_game:%v", err)
		mainMenu.ShowMessage(lang.Text("gui", "load_game_err"))
		return
	}
	// Load HUD state.
	err = hud.LoadGUISave(guisav)
	if err != nil {
		log.Err.Printf("enter_saved_game:fail_to_set_hud_layout:%v", err)
		mainMenu.ShowMessage(lang.Text("gui", "load_game_err"))
		return
	}
	// Set HUD.
	setHUD(hud)
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

// ExecuteScriptFile executes Ash script file
// with specified Name.
func ExecuteScriptFile(Name string, args ...string) error {
	modpath := game.Module().Conf().Path
	path := filepath.FromSlash(modpath + "/gui/scripts/" + Name + ".ash")
	return ci.RunScript(path, args...)
}

// setHUD sets specified HUD instance as current GUI player
// HUD.
func setHUD(h *hud.HUD) {
	pcHUD = h
	ci.SetHUD(pcHUD)
	pcHUD.Chat().SetOnCommandFunc(ExecuteCommand)
	pcHUD.Chat().SetOnScriptNameFunc(ExecuteScriptFile)
}
