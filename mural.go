/*
 * mural.go
 *
 * Copyright 2018-2023 Dariusz Sikora <ds@isangeles.dev>
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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/image/colornames"

	"github.com/faiface/beep"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/character"
	flamedata "github.com/isangeles/flame/data"
	flameres "github.com/isangeles/flame/data/res"
	"github.com/isangeles/flame/data/res/lang"

	"github.com/isangeles/burn"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/ci"
	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/data"
	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/data/res/audio"
	"github.com/isangeles/mural/data/res/graphic"
	"github.com/isangeles/mural/game"
	"github.com/isangeles/mural/hud"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/mainmenu"
)

var (
	win        *mtk.Window
	mainMenu   *mainmenu.MainMenu
	pcHUD      *hud.HUD
	mod        *flame.Module
	activeGame *game.Game
	inGame     bool
)

// Main function.
func main() {
	// Load GUI config.
	err := config.Load()
	if err != nil {
		err = config.Save() // save default config
		if err != nil {
			log.Err.Printf("Unable to save default config: %v", err)
		}
	}
	log.PrintStdOut(config.Debug)
	// Import module.
	modData, err := flamedata.ImportModuleDir(config.ModulePath())
	if err != nil {
		panic(fmt.Errorf("unable to import module: %v", err))
	}
	setModule(modData)
	// Load GUI graphic data.
	err = data.LoadModuleData(config.GUIPath)
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
	var err error
	win, err = mtk.NewWindow(winConfig)
	if err != nil {
		panic(fmt.Errorf("unable to create mtk window: %v", err))
	}
	win.SetMaxFPS(config.MaxFPS)
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
	if len(config.ServerHost+config.ServerPort) > 1 {
		s, err := game.NewServer(config.ServerHost, config.ServerPort)
		if err != nil {
			log.Err.Printf("Init run: Unable to connect to the game server: %v",
				err)
		}
		server = s
	}
	// Create main menu.
	mainMenu = mainmenu.New()
	mainMenu.SetServer(server)
	if server == nil {
		mainMenu.SetModule(mod)
	}
	mainMenu.SetOnGameCreatedFunc(enterGame)
	ci.SetMainMenu(mainMenu)
	// Debug mode.
	textParams := mtk.Params{
		FontSize: mtk.SizeMedium,
	}
	fpsInfo := mtk.NewText(textParams)
	fpsInfo.Align(mtk.AlignRight)
	verInfo := mtk.NewText(textParams)
	verInfo.SetText(fmt.Sprintf("%s(%s)@%s(%s)", config.Name, config.Version,
		flame.Name, flame.Version))
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
		if pcHUD.Exiting() || activeGame.Closing() {
			inGame = false
			// Reimport module.
			modData, err := flamedata.ImportModuleDir(config.ModulePath())
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

// enterGame creates HUD and enters game.
// The second parameter is used to apply data on newly
// created HUD, if nil the HUD will be left in its
// default state after creation.
func enterGame(g *game.Game, hudData *res.HUDData) {
	mainMenu.OpenLoadingScreen(lang.Text("enter_game_info"))
	defer mainMenu.CloseLoadingScreen()
	activeGame = g
	activeGame.AddChapterChangeEvent(changeChapter)
	// Create HUD.
	hud := hud.New(win)
	// Set HUD.
	setHUD(hud)
	// Load GUI data.
	chapterGUIPath := filepath.Join(config.GUIPath, "chapters", activeGame.Chapter().Conf().ID)
	err := data.LoadChapterData(chapterGUIPath)
	if err != nil {
		log.Err.Printf("enter game: unable to load chapter GUI data: %v", err)
		mainMenu.ShowMessage(lang.Text("load_game_err"))
		return
	}
	// Set game for HUD.
	hud.SetGame(activeGame)
	if hudData != nil {
		err = pcHUD.Apply(*hudData)
		if err != nil {
			log.Err.Printf("enter game: unable to load HUD layout: %v", err)
		}
	}
	inGame = true
	// Run module scripts.
	err = runModuleScripts()
	if err != nil {
		log.Err.Printf("enter game: unable to run module scripts: %v", err)
	}
	// Run game update.
	go activeGame.Update()
}

// changeChater handles chaper change event triggered by specified character.
func changeChapter(ob *character.Character) {
	// Change game chapter.
	activeGame.Conf().Chapter = ob.ChapterID()
	chapterPath := filepath.Join(activeGame.Conf().ChaptersPath(), activeGame.Conf().Chapter)
	chapterData, err := flamedata.ImportChapterDir(chapterPath)
	if err != nil {
		log.Err.Printf("Chapter change: unable to load chapter data: %v", err)
		return
	}
	chapter := flame.NewChapter(activeGame.Module, chapterData)
	activeGame.SetChapter(chapter)
	// Load GUI data.
	chapterGUIPath := filepath.Join(config.GUIPath, "chapters", activeGame.Chapter().Conf().ID)
	err = data.LoadChapterData(chapterGUIPath)
	if err != nil {
		log.Err.Printf("Chapter change: unable to load chapter GUI data: %v", err)
		return
	}
	// Respawn the character.
	gameChar := activeGame.Char(ob.ID(), ob.Serial())
	if gameChar == nil {
		log.Err.Printf("Chapter change: game character not found: %s %s", ob.ID(),
			ob.Serial())
		return
	}
	err = activeGame.SpawnChar(gameChar)
	if err != nil {
		log.Err.Printf("Chapter change: unable to spawn character: %v", err)
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
func setModule(data flameres.ModuleData) {
	mod = flame.NewModule(data)
	burn.Module = mod
	if mainMenu != nil {
		mainMenu.SetModule(mod)
	}
}

// runModuleScripts starts all scripts from the module
// GUI directory(scripts/run).
func runModuleScripts() error {
	scriptsPath := filepath.Join(config.GUIPath, "scripts/run")
	if _, err := os.Stat(scriptsPath); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	scripts, err := data.ScriptsDir(scriptsPath)
	if err != nil {
		return fmt.Errorf("unable to retrieve scripts: %v", err)
	}
	for _, s := range scripts {
		go ci.RunScript(s)
	}
	return nil
}
