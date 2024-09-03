/*
 * hud.go
 *
 * Copyright 2018-2024 Dariusz Sikora <ds@isangeles.dev>
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

// Package with 'head-up display' elements.
package hud

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/image/colornames"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"

	"github.com/isangeles/flame/area"
	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/objects"

	"github.com/isangeles/burn/ash"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/data"
	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/game"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/object"
)

var (
	// HUD colors.
	mainColor   = colornames.Grey
	secColor    = colornames.Blue
	accentColor = colornames.Red
	// Keys.
	pauseKey = pixelgl.KeySpace
	exitKey  = pixelgl.KeyEscape
)

// Struct for 'head-up display'.
type HUD struct {
	loadScreen    *LoadingScreen
	camera        *Camera
	bar           *MenuBar
	menu          *Menu
	savemenu      *SaveMenu
	pcFrame       *ObjectFrame
	tarFrame      *ObjectFrame
	objectInfo    *ObjectInfo
	castBar       *CastBar
	chat          *Chat
	inv           *InventoryMenu
	skills        *SkillMenu
	loot          *LootWindow
	dialog        *DialogWindow
	journal       *JournalWindow
	crafting      *CraftingMenu
	charinfo      *CharacterWindow
	trade         *TradeWindow
	training      *TrainingWindow
	game          *game.Game
	userFocus     *mtk.Focus
	msgs          *mtk.MessageQueue
	layouts       map[string]*Layout
	defaultLayout *Layout
	loading       bool
	exiting       bool
	loaderr       error
	onAreaChanged func(a *area.Area)
	areaScripts   []*ash.Script
}

// New creates new HUD instance.
func New(win *mtk.Window) *HUD {
	hud := new(HUD)
	// Loading screen.
	hud.loadScreen = newLoadingScreen(hud)
	// Camera.
	hud.camera = newCamera(hud, win.Bounds().Size())
	// Active player & target frames.
	hud.pcFrame = newObjectFrame(hud)
	hud.tarFrame = newObjectFrame(hud)
	// Hovered object info window.
	hud.objectInfo = newObjectInfo(hud)
	// Cast bar.
	hud.castBar = newCastBar(hud)
	// Windows & menus.
	hud.bar = newMenuBar(hud)
	hud.menu = newMenu(hud)
	hud.savemenu = newSaveMenu(hud)
	hud.chat = newChat(hud)
	hud.inv = newInventoryMenu(hud)
	hud.skills = newSkillMenu(hud)
	hud.loot = newLootWindow(hud)
	hud.dialog = newDialogWindow(hud)
	hud.journal = newJournalWindow(hud)
	hud.crafting = newCraftingMenu(hud)
	hud.charinfo = newCharacterWindow(hud)
	hud.trade = newTradeWindow(hud)
	hud.training = newTrainingWindow(hud)
	// Messages & focus.
	hud.userFocus = new(mtk.Focus)
	hud.msgs = mtk.NewMessageQueue(hud.UserFocus())
	// Layouts.
	hud.layouts = make(map[string]*Layout)
	hud.defaultLayout = NewLayout()
	return hud
}

// Draw draws HUD elements.
func (hud *HUD) Draw(win *mtk.Window) {
	if hud.loading {
		hud.loadScreen.Draw(win)
		return
	}
	if hud.Game() == nil || hud.Game().ActivePlayerChar() == nil { // no game or active pc, don't draw
		return
	}
	// Elements positions.
	pcFramePos := mtk.DrawPosTL(win.Bounds(), hud.pcFrame.Size())
	tarFramePos := mtk.RightOf(hud.pcFrame.DrawArea(), hud.tarFrame.Size(), 0)
	castBarPos := win.Bounds().Center()
	barPos := mtk.DrawPosBC(win.Bounds(), hud.bar.Size())
	chatPos := mtk.DrawPosBL(win.Bounds(), hud.chat.Size())
	menuPos := win.Bounds().Center()
	saveMenuPos := win.Bounds().Center()
	invPos := win.Bounds().Center()
	skillsPos := win.Bounds().Center()
	lootPos := win.Bounds().Center()
	dialogPos := win.Bounds().Center()
	journalPos := win.Bounds().Center()
	craftingPos := win.Bounds().Center()
	charinfoPos := win.Bounds().Center()
	tradePos := win.Bounds().Center()
	trainPos := win.Bounds().Center()
	// Draw elements.
	hud.camera.Draw(win)
	hud.bar.Draw(win, mtk.Matrix().Moved(barPos))
	hud.chat.Draw(win, mtk.Matrix().Moved(chatPos))
	hud.pcFrame.Draw(win, mtk.Matrix().Moved(pcFramePos))
	if len(hud.Game().ActivePlayerChar().Targets()) > 0 {
		hud.tarFrame.Draw(win, mtk.Matrix().Moved(tarFramePos))
	}
	if hud.savemenu.Opened() {
		hud.savemenu.Draw(win, mtk.Matrix().Moved(saveMenuPos))
	}
	if hud.inv.Opened() {
		hud.inv.Draw(win, mtk.Matrix().Moved(invPos))
	}
	if hud.skills.Opened() {
		hud.skills.Draw(win, mtk.Matrix().Moved(skillsPos))
	}
	if hud.loot.Opened() {
		hud.loot.Draw(win, mtk.Matrix().Moved(lootPos))
	}
	if hud.dialog.Opened() {
		hud.dialog.Draw(win, mtk.Matrix().Moved(dialogPos))
	}
	if hud.journal.Opened() {
		hud.journal.Draw(win, mtk.Matrix().Moved(journalPos))
	}
	if hud.crafting.Opened() {
		hud.crafting.Draw(win, mtk.Matrix().Moved(craftingPos))
	}
	if hud.charinfo.Opened() {
		hud.charinfo.Draw(win, mtk.Matrix().Moved(charinfoPos))
	}
	if hud.trade.Opened() {
		hud.trade.Draw(win, mtk.Matrix().Moved(tradePos))
	}
	if hud.training.Opened() {
		hud.training.Draw(win, mtk.Matrix().Moved(trainPos))
	}
	if hud.objectInfo.Opened() {
		hud.objectInfo.Draw(win)
	}
	if hud.Game().ActivePlayerChar().Casted() != nil {
		hud.castBar.Draw(win, mtk.Matrix().Moved(castBarPos))
	}
	if hud.menu.Opened() {
		hud.menu.Draw(win, mtk.Matrix().Moved(menuPos))
	}
	// Messages.
	msgPos := win.Bounds().Center()
	hud.msgs.Draw(win, mtk.Matrix().Moved(msgPos))
}

// Update updated HUD elements.
func (hud *HUD) Update(win *mtk.Window) {
	// HUD state.
	if hud.loading && hud.loaderr != nil { // on loading error
		log.Err.Printf("hud loading fail: %v", hud.loaderr)
		hud.Exit()
	}
	if hud.Game() == nil || hud.Game().ActivePlayerChar() == nil { // no game or active pc, don't update
		return
	}
	// Handle area change.
	hud.updateCurrentArea()
	// Toggle game pause.
	if !hud.Chat().Activated() && win.JustPressed(pauseKey) {
		hud.Game().SetPause(!hud.Game().Pause())
	}
	// Put PC target into target frame.
	if hud.camera.area != nil && len(hud.Game().ActivePlayerChar().Targets()) > 0 {
		for _, av := range hud.camera.area.Avatars() {
			if objects.Equals(hud.Game().ActivePlayerChar().Targets()[0], av.Character) {
				hud.tarFrame.SetObject(av)
			}
		}
	}
	// Elements update.
	hud.loadScreen.Update(win)
	hud.camera.Update(win)
	hud.bar.Update(win)
	hud.chat.Update(win)
	hud.pcFrame.Update(win)
	hud.tarFrame.Update(win)
	hud.castBar.Update(win)
	hud.objectInfo.Update(win)
	hud.menu.Update(win)
	hud.savemenu.Update(win)
	hud.inv.Update(win)
	hud.skills.Update(win)
	hud.loot.Update(win)
	hud.dialog.Update(win)
	hud.journal.Update(win)
	hud.crafting.Update(win)
	hud.charinfo.Update(win)
	hud.trade.Update(win)
	hud.training.Update(win)
	hud.msgs.Update(win)
}

// SetActiveChar sets specified game character as active
// character.
func (hud *HUD) SetActiveChar(char *game.Character) {
	pcAvatar := hud.PCAvatar()
	if pcAvatar != nil {
		hud.camera.CenterAt(pcAvatar.Position())
		hud.pcFrame.SetObject(pcAvatar)
	}
	hud.Reload()
}

// Exit sends exit request to HUD.
func (hud *HUD) Exit() {
	hud.exiting = true
}

// Checks whether exit was requested.
func (hud *HUD) Exiting() bool {
	return hud.exiting
}

// Chat returns HUD chat.
func (hud *HUD) Chat() *Chat {
	return hud.chat
}

// Camera returns HUD camera.
func (hud *HUD) Camera() *Camera {
	return hud.camera
}

// Game returns HUD current game.
func (hud *HUD) Game() *game.Game {
	return hud.game
}

// UserFocus returns HUD user focus.
func (hud *HUD) UserFocus() *mtk.Focus {
	return hud.userFocus
}

// ShowMessage displays specified message window in HUD
// messages queue.
func (hud *HUD) ShowMessage(msg *mtk.MessageWindow) {
	msg.Show(true)
	hud.msgs.Append(msg)
}

// OpenLoadingScreen opens loading screen with
// specified loading information.
func (hud *HUD) OpenLoadingScreen(info string) {
	hud.loading = true
	hud.loadScreen.SetLoadInfo(info)
}

// Close loading screen closes loading screen.
func (hud *HUD) CloseLoadingScreen() {
	hud.loading = false
}

// Layout returns layout for player with specified ID
// serial value(creates new layout if there is no saved
// layout for such player).
func (hud *HUD) Layout(id, serial string) *Layout {
	layout := hud.layouts[id+serial]
	if layout == nil {
		layout = hud.defaultLayout
		hud.layouts[id+serial] = layout
	}
	return layout
}

// Reload reloads HUD layouts for
// active player.
func (hud *HUD) Reload() {
	pc := hud.Game().ActivePlayerChar()
	if pc == nil {
		return
	}
	layout := hud.Layout(pc.ID(), pc.Serial())
	hud.bar.setLayout(layout)
}

// SetGame sets HUD game.
func (hud *HUD) SetGame(g *game.Game) {
	hud.game = g
	hud.game.SetOnPlayerCharChangeFunc(hud.SetActiveChar)
}

// PCAvatar return avatar for player current character.
func (hud *HUD) PCAvatar() *object.Avatar {
	if hud.camera.area == nil {
		return nil
	}
	pc := hud.game.ActivePlayerChar()
	for _, av := range hud.camera.area.Avatars() {
		if av.ID() == pc.ID() && av.Serial() == pc.Serial() {
			return av
		}
	}
	return nil
}

// ChangeArea changes current HUD area.
func (hud *HUD) ChangeArea(area *area.Area) error {
	// Stop previous area scripts.
	for _, s := range hud.areaScripts {
		s.Stop(true)
	}
	hud.areaScripts = make([]*ash.Script, 0)
	// Map.
	hud.OpenLoadingScreen(lang.Text("load_map_info"))
	defer hud.CloseLoadingScreen()
	err := hud.camera.SetArea(area)
	if err != nil {
		hud.loaderr = fmt.Errorf("unable to set camera area: %v",
			err)
		return hud.loaderr
	}
	pcAvatar := hud.PCAvatar()
	if pcAvatar != nil {
		hud.camera.CenterAt(pcAvatar.Position())
		hud.pcFrame.SetObject(pcAvatar)
	}
	// Reload HUD.
	hud.Reload()
	hud.runAreaScripts(area)
	return nil
}

// Data returns data struct for HUD.
func (hud *HUD) Data() res.HUDData {
	var data res.HUDData
	// Players.
	for _, pc := range hud.game.PlayerChars() {
		pcData := res.Player{ID: pc.ID(), Serial: pc.Serial()}
		// Layout.
		layout := hud.layouts[pc.ID()+pc.Serial()]
		if layout != nil {
			for serialID, slot := range layout.InvSlots() {
				slot := res.Slot{slot, serialID}
				pcData.InvSlots = append(pcData.InvSlots, slot)
			}
			for serialID, slot := range layout.BarSlots() {
				slot := res.Slot{slot, serialID}
				pcData.BarSlots = append(pcData.BarSlots, slot)
			}
		}
		data.Players = append(data.Players, pcData)
	}
	// Camera XY position.
	data.Camera.X = hud.Camera().Position().X
	data.Camera.Y = hud.Camera().Position().Y
	return data
}

// Apply applies specified data on the HUD.
func (hud *HUD) Apply(data res.HUDData) error {
	// Players.
	for _, pcd := range data.Players {
		layout := NewLayout()
		slotsLayout := make(map[string]int)
		for _, s := range pcd.InvSlots {
			slotsLayout[s.Content] = s.ID
		}
		layout.SetInvSlots(slotsLayout)
		slotsLayout = make(map[string]int)
		for _, s := range pcd.BarSlots {
			slotsLayout[s.Content] = s.ID
		}
		layout.SetBarSlots(slotsLayout)
		if pcd.ID == "*" {
			hud.defaultLayout = layout
			continue
		}
		layoutKey := pcd.ID + pcd.Serial
		hud.layouts[layoutKey] = layout
	}
	// Camera position.
	hud.camera.SetPosition(pixel.V(data.Camera.X, data.Camera.Y))
	// Reload UI.
	hud.Reload()
	return nil
}

// SetOnAreaChangedFunc sets function triggered on area change.
func (hud *HUD) SetOnAreaChangedFunc(f func(a *area.Area)) {
	hud.onAreaChanged = f
}

// runAreaScripts executes all scripts for specified area
// placed in ui/mural/chapters/[chapter]/areas/scripts/[area].
func (hud *HUD) runAreaScripts(a *area.Area) {
	// Retrive scripts.
	mod := hud.Game().Module
	path := filepath.Join(config.GUIPath, "chapters", mod.Chapter().ID(), "areas", a.ID(),
		"scripts")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return
	}
	scripts, err := data.ScriptsDir(path)
	if err != nil {
		log.Err.Printf("hud: run area scripts: unable to retrieve scripts: %v", err)
		return
	}
	// Run scripts in background.
	for _, s := range scripts {
		go hud.RunScript(s)
		hud.areaScripts = append(hud.areaScripts, s)
		log.Dbg.Printf("script started: %s\n", s.Name())
	}
}

// updateCurrentArea updates HUD area to active player area.
func (hud *HUD) updateCurrentArea() {
	chapter := hud.Game().Chapter()
	pcArea := chapter.ObjectArea(hud.Game().ActivePlayerChar())
	if hud.game == nil {
		return
	}
	if pcArea == nil {
		return
	}
	if hud.camera.Area() == nil || pcArea.ID() != hud.camera.Area().ID() {
		go hud.ChangeArea(pcArea)
	}
}

// containsPos checks if specified position is contained
// by any HUD element(except camera).
func (hud *HUD) containsPos(pos pixel.Vec) bool {
	return hud.bar.DrawArea().Contains(pos) ||
		hud.chat.DrawArea().Contains(pos) ||
		hud.pcFrame.DrawArea().Contains(pos) ||
		(hud.inv.Opened() && hud.inv.DrawArea().Contains(pos)) ||
		(hud.menu.Opened() && hud.menu.DrawArea().Contains(pos)) ||
		(hud.savemenu.Opened() && hud.savemenu.DrawArea().Contains(pos)) ||
		(hud.skills.Opened() && hud.skills.DrawArea().Contains(pos)) ||
		(hud.loot.Opened() && hud.loot.DrawArea().Contains(pos)) ||
		(hud.dialog.Opened() && hud.dialog.DrawArea().Contains(pos)) ||
		(hud.journal.Opened() && hud.journal.DrawArea().Contains(pos)) ||
		(hud.crafting.Opened() && hud.crafting.DrawArea().Contains(pos)) ||
		(hud.trade.Opened() && hud.trade.DrawArea().Contains(pos)) ||
		(hud.training.Opened() && hud.training.DrawArea().Contains(pos)) ||
		(hud.charinfo.Opened() && hud.charinfo.DrawArea().Contains(pos))
}

// menuOpen checks if any HUD menu is open.
func (hud *HUD) menuOpen() bool {
	return hud.charinfo.Opened() || hud.crafting.Opened() ||
		hud.dialog.Opened() || hud.inv.Opened() ||
		hud.journal.Opened() || hud.loot.Opened() ||
		hud.menu.Opened() || hud.savemenu.Opened() ||
		hud.trade.Opened() || hud.training.Opened()
}
