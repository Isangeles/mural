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
	"github.com/isangeles/flame/core/module/scenario"

	"github.com/isangeles/burn"
	"github.com/isangeles/burn/ash"
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

var (
	mainMenu    *mainmenu.MainMenu
	pcHUD       *hud.HUD
	game        *flamecore.Game
	inGame      bool
	areaScripts []*ash.Script
)

// On init.
func init() {
	// Load flame config.
	err := flameconf.LoadConfig()
	if err != nil {
		log.Err.Printf("fail to load flame config file: %v", err)
		flameconf.SaveConfig() // override 'corrupted' config file with default configuration
	}
	// Load module.
	m, err := flamedata.Module(flameconf.ModulePath(), flameconf.LangID())
	if err != nil {
		log.Err.Printf("fail to load config module: %v", err)
	}
	flame.SetModule(m)
	// Load GUI config.
	err = config.LoadConfig()
	if err != nil {
		log.Err.Printf("fail to load config file: %v", err)
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
		panic(fmt.Errorf("fail to load gui data: %v", err))
	}
	// Load game graphic.
	err = data.LoadGameData()
	if err != nil {
		panic(fmt.Errorf("fail to load game graphic data: %v", err))
	}
	// Load module data.
	err = flamedata.LoadModuleData(flame.Mod())
	if err != nil {
		panic(fmt.Errorf("fail to load module data: %v", err))
	}
	// Load module graphic data.
	err = imp.LoadModuleResources(flame.Mod())
	if err != nil {
		panic(fmt.Errorf("fail to load module resources: %v", err))
	}
	// Music.
	mtk.InitAudio(beep.Format{44100, 2, 2})
	if mtk.Audio() != nil {
		ci.SetMusicPlayer(mtk.Audio())
		m, err := data.Music(config.MenuMusic)
		if err == nil {
			pl := []beep.Streamer{m.Streamer(0, m.Len())}
			mtk.Audio().SetPlaylist(pl)
		} else {
			log.Err.Printf("fail to load main theme audio data: %v", err)
		}
		mtk.Audio().SetVolume(config.MusicVolume)
		mtk.Audio().SetMute(config.MusicMute)
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
	winRes := config.Resolution
	if winRes.X == 0 || winRes.Y == 0 {
		winRes.X, winRes.Y = monitor.Size()
	}
	cfg := pixelgl.WindowConfig{
		Title:  config.Name + " " + config.Version,
		Bounds: pixel.R(winPosX, winPosY, winRes.X, winRes.Y),
		VSync:  true,
	}
	if config.Fullscreen {
		monitor := pixelgl.PrimaryMonitor()
		cfg.Monitor = monitor
	}
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("fail to create mtk window: %v", err))
	}
	// UI Font.
	uiFont, err := data.Font(config.MainFont)
	if err == nil { // if font from config was found
		mtk.SetMainFont(uiFont)
	}
	// Audio effects.
	bClickSound, err := data.AudioEffect(config.ButtonClickSound)
	if err != nil {
		log.Err.Printf("init run: fail to retrieve button click audio data: %v",
			err)
	}
	mtk.SetButtonClickSound(bClickSound) // global button click sound
	// Create main menu.
	mainMenu = mainmenu.New()
	mainMenu.SetOnGameCreatedFunc(EnterGame)
	mainMenu.SetOnSaveLoadedFunc(EnterSavedGame)
	err = mainMenu.ImportPlayableChars(flame.Mod().Conf().CharactersPath())
	if err != nil {
		log.Err.Printf("init run: fail to import playable characters: %v",
			err)
	}
	ci.SetMainMenu(mainMenu)
	mainMenu.Console().SetOnCommandFunc(ExecuteCommand)
	// Debug mode.
	textParams := mtk.Params{
		FontSize: mtk.SizeMedium,
	}
	fpsInfo := mtk.NewText(textParams)
	fpsInfo.Align(mtk.AlignRight)
	verInfo := mtk.NewText(textParams)
	verInfo.SetText(fmt.Sprintf("%s(%s)@%s(%s)", config.Name, config.Version,
		flame.NAME, flame.VERSION))
	verInfo.Align(mtk.AlignRight)
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
	// Set HUD.
	setHUD(hud)
	// Set game for HUD.
	err := imp.LoadChapterResources(game.Module().Chapter())
	if err != nil {
		log.Err.Printf("enter game: fail to load chapter resources: %v", err)
		mainMenu.ShowMessage(lang.Text("gui", "load_game_err"))
		return
	}
	err = hud.SetGame(game)
	if err != nil {
		log.Err.Printf("enter game: fail to set hud game: %v", err)
		mainMenu.ShowMessage(lang.Text("gui", "load_game_err"))
		return
	}
	inGame = true
	// Run module scripts.
	modpath := game.Module().Conf().Path
	scripts, err := data.ScriptsDir(modpath + "/gui/scripts/run")
	if err != nil {
		log.Err.Printf("enter saved game: fail to retrieve module scripts: %v", err)
	}
	for _, s := range scripts {
		go ci.RunScript(s)
	}
}

// EnterSavedGame creates game and HUD from saved data.
func EnterSavedGame(g *flamecore.Game, saveName string) {
	mainMenu.OpenLoadingScreen(lang.Text("gui", "loadgame_load_game_info"))
	defer mainMenu.CloseLoadingScreen()
	game = g
	// Import saved GUI state.
	guisav, err := imp.ImportGUISave(flameconf.ModuleSavegamesPath(), saveName)
	if err != nil {
		log.Err.Printf("enter saved game: fail  load gui save: %v", err)
		mainMenu.ShowMessage(lang.Text("gui", "load_game_err"))
		return
	}
	for _, pcd := range guisav.PlayersData {
		res.AddAvatarData(pcd.Avatar)
	}
	// Create HUD.
	hud := hud.New()
	// Set HUD.
	setHUD(hud)
	err = imp.LoadChapterResources(game.Module().Chapter())
	err = hud.SetGame(game)
	if err != nil {
		log.Err.Printf("enter saved game: fail to set hud game: %v", err)
		mainMenu.ShowMessage(lang.Text("gui", "load_game_err"))
		return
	}
	// Load HUD state.
	err = hud.LoadGUISave(guisav)
	if err != nil {
		log.Err.Printf("enter saved game: fail to set hud layout: %v", err)
		mainMenu.ShowMessage(lang.Text("gui", "load_game_err"))
		return
	}
	inGame = true
	// Run module scripts.
	modpath := game.Module().Conf().Path
	scripts, err := data.ScriptsDir(modpath + "/gui/scripts/run")
	if err != nil {
		log.Err.Printf("enter saved game: fail to retrieve module scripts: %v", err)
	}
	for _, s := range scripts {
		go ci.RunScript(s)
	}
}

// ExecuteCommand handles specified text line
// as CI command.
// Returns result code and output text, or error if
// specified line is not valid command.
func ExecuteCommand(line string) (int, string, error) {
	cmd, err := syntax.NewSTDExpression(line)
	if err != nil {
		return -1, "", fmt.Errorf("invalid input: %s", line)
	}
	res, out := burn.HandleExpression(cmd)
	return res, out, nil
}

// ExecuteScriptFile executes Ash script from file
// with specified name in background.
func ExecuteScriptFile(name string, args ...string) error {
	modpath := game.Module().Conf().Path
	path := filepath.FromSlash(modpath + "/gui/scripts/" + name + ".ash")
	script, err := data.Script(path)
	if err != nil {
		return fmt.Errorf("fail to retrieve script: %v", err)
	}
	go ci.RunScript(script)
	return nil
}

// RunAreaScripts executes all scripts for specified area
// placed in gui/chapters/[chapter]/areas/scripts/[area].
func RunAreaScripts(a *scenario.Area) {
	// Stop previous area scripts.
	for _, s := range areaScripts {
		s.Stop(true)
	}
	// Retrive scripts.
	mod := game.Module()
	path := filepath.FromSlash(mod.Conf().Path + "/gui/chapters/" +
		mod.Chapter().ID() + "/areas/scripts/" + a.ID())
	scripts, err := data.ScriptsDir(path)
	if err != nil {
		log.Err.Printf("run area scripts: fail to retrieve scripts: %v", err)
		return
	}
	// Run scripts in background.
	for _, s := range scripts {
		go ci.RunScript(s)
		areaScripts = append(areaScripts, s)
	}
}

// setHUD sets specified HUD instance as current
// GUI player HUD.
func setHUD(h *hud.HUD) {
	pcHUD = h
	ci.SetHUD(pcHUD)
	pcHUD.SetOnAreaChangedFunc(RunAreaScripts)
	pcHUD.Chat().SetOnCommandFunc(ExecuteCommand)
	pcHUD.Chat().SetOnScriptNameFunc(ExecuteScriptFile)
}
