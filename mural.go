/*
 * mural.go
 *
 * Copyright 2018-2020 Dariusz Sikora <dev@isangeles.pl>
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

	flameconf "github.com/isangeles/flame/config"
	flamedata "github.com/isangeles/flame/data"
	"github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/module"

	"github.com/isangeles/burn"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/ci"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res/audio"
	"github.com/isangeles/mural/core/data/res/graphic"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/game"
	"github.com/isangeles/mural/hud"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/mainmenu"
)

var (
	mainMenu   *mainmenu.MainMenu
	pcHUD      *hud.HUD
	mod        *module.Module
	activeGame *game.Game
	inGame     bool
)

// On init.
func init() {
	// Load GUI config.
	err := config.Load()
	if err != nil {
		log.Err.Printf("unable to load config file: %v", err)
	}
	// Load UI translation files.
	err = flamedata.LoadTranslationData(config.LangPath())
	if err != nil {
		log.Err.Printf("unable to load ui translation files: %v", err)
	}
}

// Main function.
func main() {
	// Import module.
	modData, err := flamedata.ImportModule(config.ModulePath())
	if err != nil {
		panic(fmt.Errorf("unable to import module: %v", err))
	}
	setModule(modData)
	// Load module translation.
	err = flamedata.LoadModuleLang(mod, config.Lang)
	if err != nil {
		log.Err.Printf("unable to load module translation: %v", err)
	}
	// Load GUI graphic data.
	err = data.LoadModuleData(mod)
	if err != nil {
		panic(fmt.Errorf("unable to load game graphic data: %v", err))
	}
	// Music.
	mtk.InitAudio(beep.Format{44100, 2, 2})
	if mtk.Audio() != nil {
		ci.SetMusicPlayer(mtk.Audio())
		m := audio.Music[config.MenuMusic]
		if m != nil {
			pl := []beep.Streamer{m.Streamer(0, m.Len())}
			mtk.Audio().SetPlaylist(pl)
		} else {
			log.Err.Printf("main theme audio data not found: %s",
				config.MenuMusic)
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
	winConfig := pixelgl.WindowConfig{
		Title:  config.Name + " " + config.Version,
		Bounds: pixel.R(winPosX, winPosY, winRes.X, winRes.Y),
		VSync:  true,
	}
	if config.Fullscreen {
		winConfig.Monitor = pixelgl.PrimaryMonitor()
	}
	win, err := mtk.NewWindow(winConfig)
	if err != nil {
		panic(fmt.Errorf("unable to create mtk window: %v", err))
	}
	// UI Font.
	uiFont := graphic.Fonts[config.MainFont]
	if uiFont != nil {
		mtk.SetMainFont(uiFont)
	}
	// Audio effects.
	bClickSound := audio.Effects[config.ButtonClickSound]
	if bClickSound == nil {
		log.Err.Printf("init run: button click audio data not found: %s",
			config.ButtonClickSound)
	}
	mtk.SetButtonClickSound(bClickSound) // global button click sound
	// Fire mode.
	var server *game.Server
	if config.Fire {
		s, err := game.NewServer(config.ServerHost, config.ServerPort)
		if err != nil {
			log.Err.Printf("Init run: Unable to connect to the game server: %v",
				err)
		}
		server = s
	}
	// Create main menu.
	mainMenu = mainmenu.New(mod, server)
	mainMenu.SetOnGameCreatedFunc(EnterGame)
	mainMenu.SetOnSaveLoadFunc(LoadSavedGame)
	err = mainMenu.ImportPlayableChars()
	if err != nil {
		log.Err.Printf("init run: unable to import playable characters: %v",
			err)
	}
	ci.SetMainMenu(mainMenu)
	// Debug mode.
	textParams := mtk.Params{
		FontSize: mtk.SizeMedium,
	}
	fpsInfo := mtk.NewText(textParams)
	fpsInfo.Align(mtk.AlignRight)
	verInfo := mtk.NewText(textParams)
	verInfo.SetText(fmt.Sprintf("%s(%s)@%s(%s)", config.Name, config.Version,
		flameconf.Name, flameconf.Version))
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
		if config.Debug {
			fpsPos := mtk.DrawPosTR(win.Bounds(), fpsInfo.Size())
			fpsPos.Y -= mtk.ConvSize(10)
			fpsInfo.Draw(win, mtk.Matrix().Moved(fpsPos))
			verPos := mtk.LeftOf(fpsInfo.DrawArea(), verInfo.Size(), 5)
			verInfo.Draw(win, mtk.Matrix().Moved(verPos))
			fpsInfo.SetText(fmt.Sprintf("FPS: %d", win.FPS()))
		}
		// Update.
		win.Update()
		if !inGame {
			mainMenu.Update(win)
			continue
		}
		pcHUD.Update(win)
		activeGame.Update(win.Delta())
		if pcHUD.Exiting() || activeGame.Closing() {
			inGame = false
			// Reimport module.
			modData, err := flamedata.ImportModule(config.ModulePath())
			if err != nil {
				log.Err.Printf("unable to reimport module: %v", err)
			}
			setModule(modData)
		}
	}
	// On exit.
	if win.Closed() {
		config.Save()
	}
}

// EnterGame creates HUD for specified game.
func EnterGame(g *game.Game) {
	mainMenu.OpenLoadingScreen(lang.Text("enter_game_info"))
	defer mainMenu.CloseLoadingScreen()
	activeGame = g
	// Create HUD.
	hud := hud.New()
	// Set HUD.
	setHUD(hud)
	// Load GUI data.
	err := data.LoadChapterData(activeGame.Module().Chapter())
	if err != nil {
		log.Err.Printf("enter game: unable to load chapter GUI data: %v", err)
		mainMenu.ShowMessage(lang.Text("load_game_err"))
		return
	}
	// Set game for HUD.
	hud.SetGame(activeGame)
	burn.Game = activeGame.Game 
	inGame = true
	// Run module scripts.
	modpath := activeGame.Module().Conf().Path
	scriptsPath := filepath.Join(modpath, data.GUIModulePath, "scripts/run")
	scripts, err := data.ScriptsDir(scriptsPath)
	if err != nil {
		log.Err.Printf("enter game: unable to retrieve module scripts: %v", err)
	}
	for _, s := range scripts {
		go ci.RunScript(s)
	}
}

// LoadSavedGame creates game and HUD from saved data.
func LoadSavedGame(saveName string) {
	mainMenu.OpenLoadingScreen(lang.Text("loadgame_load_game_info"))
	defer mainMenu.CloseLoadingScreen()
	// Import saved game.
	savePath := filepath.Join(mod.Conf().SavesPath(),
		saveName+flamedata.SavegameFileExt)
	g, err := flamedata.ImportGame(mod, savePath)
	if err != nil {
		log.Err.Printf("load saved game: unable to import game: %v", err)
		mainMenu.ShowMessage(lang.Text("load_game_err"))
		return
	}
	// Create game wrapper.
	gameWrapper, err := game.New(g, nil)
	if err != nil {
		log.Err.Printf("load saved game: unable to create saved game: %v", err)
		mainMenu.ShowMessage(lang.Text("lang_game_err"))
		return
	}
	// Import saved HUD state.
	guiSavePath := filepath.Join(mod.Conf().Path, data.SavesModulePath,
		saveName+data.SaveFileExt)
	guisav, err := data.ImportGUISave(guiSavePath)
	if err != nil {
		log.Err.Printf("load saved game: unable to load gui save: %v", err)
		mainMenu.ShowMessage(lang.Text("load_game_err"))
		return
	}
	for _, pcd := range guisav.Players {
		char := gameWrapper.Module().Chapter().Character(pcd.Avatar.ID, pcd.Avatar.Serial)
		if char == nil {
			log.Err.Printf("load saved game: unable to retrieve pc character: %s#%s",
				pcd.Avatar.ID, pcd.Avatar.Serial)
			continue
		}
		av := object.NewAvatar(char, &pcd.Avatar)
		err := gameWrapper.AddPlayer(av)
		if err != nil {
			log.Err.Printf("load saved game: unable to add player to the game: %v",
				err)
			continue
		}
	}
	// Enter game.
	EnterGame(gameWrapper)
	// Load HUD state.
	err = pcHUD.LoadGUISave(guisav)
	if err != nil {
		log.Err.Printf("load saved game: unable to set hud layout: %v", err)
		mainMenu.ShowMessage(lang.Text("load_game_err"))
		return
	}
}

// setHUD sets specified HUD instance as current
// GUI player HUD.
func setHUD(h *hud.HUD) {
	pcHUD = h
	ci.SetHUD(pcHUD)
}

// setModule sets specified module for UI.
func setModule(data res.ModuleData) {
	mod = module.New(data)
	burn.Module = mod
	if mainMenu != nil {
		mainMenu.SetModule(mod)
	}
}
