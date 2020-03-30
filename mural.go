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

	"golang.org/x/image/colornames"

	"github.com/faiface/beep"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame"
	flameconf "github.com/isangeles/flame/config"
	flamedata "github.com/isangeles/flame/data"
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/module"

	"github.com/isangeles/burn"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/ci"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/imp"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/hud"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/mainmenu"
)

var (
	mainMenu *mainmenu.MainMenu
	pcHUD    *hud.HUD
	mod      *module.Module
	game     *flame.Game
	inGame   bool
)

// On init.
func init() {
	// Load flame config.
	err := flameconf.LoadConfig()
	if err != nil {
		log.Err.Printf("unable to load flame config file: %v", err)
	}
	// Load UI translation files.
	err = flamedata.LoadTranslationData(flameconf.LangPath())
	if err != nil {
		log.Err.Printf("unable to load ui translation files: %v", err)
	}
	// Load GUI config.
	err = config.LoadConfig()
	if err != nil {
		log.Err.Printf("unable to load config file: %v", err)
	}
}

// Main function.
func main() {
	// Import module.
	m, err := flamedata.ImportModule(flameconf.ModulePath())
	if err != nil {
		panic(fmt.Errorf("unable to import module: %v", err))
	}
	setModule(m)
	// Load module translation data.
	err = flamedata.LoadModuleLang(mod, flameconf.Lang)
	if err != nil {
		panic(fmt.Errorf("unable to load module translation data: %v", err))
	}
	// Load UI graphic.
	err = data.LoadUIData(mod)
	if err != nil {
		panic(fmt.Errorf("unable to load gui data: %v", err))
	}
	// Load game graphic.
	err = data.LoadModuleData(mod)
	if err != nil {
		panic(fmt.Errorf("unable to load game graphic data: %v", err))
	}
	// Load module graphic data.
	err = imp.LoadModuleResources(mod)
	if err != nil {
		panic(fmt.Errorf("unable to load module resources: %v", err))
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
			log.Err.Printf("unable to load main theme audio data: %v", err)
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
	uiFont := data.Font(config.MainFont)
	if uiFont != nil {
		mtk.SetMainFont(uiFont)
	}
	// Audio effects.
	bClickSound, err := data.AudioEffect(config.ButtonClickSound)
	if err != nil {
		log.Err.Printf("init run: unable to retrieve button click audio data: %v",
			err)
	}
	mtk.SetButtonClickSound(bClickSound) // global button click sound
	// Create main menu.
	mainMenu = mainmenu.New(mod)
	mainMenu.SetOnGameCreatedFunc(EnterGame)
	mainMenu.SetOnSaveLoadFunc(LoadSavedGame)
	err = mainMenu.ImportPlayableChars(mod.Conf().CharactersPath())
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
		if config.Debug() {
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
		game.Update(win.Delta())
		if pcHUD.Exiting() {
			inGame = false
			// Reimport module.
			m, err := flamedata.ImportModule(flameconf.ModulePath())
			if err != nil {
				log.Err.Printf("unable to reimport module: %v", err)
			}
			setModule(m)
		}
	}
	// On exit.
	if win.Closed() {
		config.SaveConfig()
		flameconf.SaveConfig()
	}
}

// EnterGame creates HUD for specified game.
func EnterGame(g *flame.Game, pcs ...*object.Avatar) {
	mainMenu.OpenLoadingScreen(lang.Text("enter_game_info"))
	defer mainMenu.CloseLoadingScreen()
	game = g
	// Create HUD.
	hud := hud.New()
	// Set HUD.
	setHUD(hud)
	// Load GUI data.
	err := imp.LoadChapterResources(game.Module().Chapter())
	if err != nil {
		log.Err.Printf("enter game: unable to load chapter GUI data: %v", err)
		mainMenu.ShowMessage(lang.Text("load_game_err"))
		return
	}
	for _, pc := range pcs {
		hud.AddPlayer(pc)
	}
	// Set game for HUD.
	hud.SetGame(game)
	burn.Game = game
	inGame = true
	// Run module scripts.
	modpath := game.Module().Conf().Path
	scripts, err := data.ScriptsDir(modpath + "/gui/scripts/run")
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
	game, err := flamedata.ImportGame(mod, flameconf.ModuleSavegamesPath(), saveName)
	if err != nil {
		log.Err.Printf("load saved game: unable to import game: %v", err)
		mainMenu.ShowMessage(lang.Text("load_game_err"))
		return
	}
	// Import saved HUD state.
	guisav, err := imp.ImportGUISave(flameconf.ModuleSavegamesPath(), saveName)
	if err != nil {
		log.Err.Printf("load saved game: unable to load gui save: %v", err)
		mainMenu.ShowMessage(lang.Text("load_game_err"))
		return
	}
	pcs := make([]*object.Avatar, 0)
	for _, pcd := range guisav.PlayersData {
		char := game.Module().Chapter().Character(pcd.Avatar.ID, pcd.Avatar.Serial)
		if char == nil {
			log.Err.Printf("load saved game: unable to retrieve pc character: %s#%s",
				pcd.Avatar.ID, pcd.Avatar.Serial)
			continue
		}
		av := object.NewAvatar(char, pcd.Avatar)
		pcs = append(pcs, av)
	}
	// Enter game.
	EnterGame(game, pcs...)
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
func setModule(m *module.Module) {
	mod = m
	burn.Module = m
	if mainMenu != nil {
		mainMenu.SetModule(m)
	}
}
