/*
 * hud.go
 *
 * Copyright 2018-2020 Dariusz Sikora <dev@isangeles.pl>
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
	"fmt"
	"path/filepath"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/module/area"
	flameobject "github.com/isangeles/flame/module/objects"

	"github.com/isangeles/burn/ash"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/game"
	"github.com/isangeles/mural/log"
)

var (
	// HUD colors.
	mainColor   = colornames.Grey
	secColor    = colornames.Blue
	accentColor = colornames.Red
	// Keys.
	pauseKey = pixelgl.KeySpace
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
	msgs          *mtk.MessagesQueue
	layouts       map[string]*Layout
	loading       bool
	exiting       bool
	loaderr       error
	onAreaChanged func(a *area.Area)
	areaScripts   []*ash.Script
}

// New creates new HUD instance.
func New() *HUD {
	hud := new(HUD)
	// Loading screen.
	hud.loadScreen = newLoadingScreen(hud)
	// Camera.
	hud.camera = newCamera(hud, config.Resolution)
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
	hud.msgs = mtk.NewMessagesQueue(hud.UserFocus())
	// Layouts.
	hud.layouts = make(map[string]*Layout)
	return hud
}

// Draw draws HUD elements.
func (hud *HUD) Draw(win *mtk.Window) {
	if hud.loading {
		hud.loadScreen.Draw(win)
		return
	}
	if hud.Game() == nil || hud.Game().ActivePlayer() == nil { // no game or active pc, don't draw
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
	if len(hud.Game().ActivePlayer().Targets()) > 0 {
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
	if hud.Game().ActivePlayer().Casted() != nil {
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
	if hud.Game() == nil || hud.Game().ActivePlayer() == nil { // no game or active pc, don't update
		return
	}
	// Handle area change.
	hud.updateCurrentArea()
	// Toggle game pause.
	if !hud.Chat().Activated() && win.JustPressed(pauseKey) {
		hud.Game().Pause(!hud.Game().Paused())
	}
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		hud.onMouseLeftPressed(win.MousePosition())
	}
	if win.JustPressed(pixelgl.MouseButtonRight) {
		hud.onMouseRightPressed(win.MousePosition())
	}
	// Put PC target into target frame.
	if len(hud.Game().ActivePlayer().Targets()) > 0 {
		for _, av := range hud.camera.Avatars() {
			if flameobject.Equals(hud.Game().ActivePlayer().Targets()[0], av.Character) {
				hud.tarFrame.SetObject(av)
			}
		}
		for _, ob := range hud.camera.AreaObjects() {
			if flameobject.Equals(hud.Game().ActivePlayer().Targets()[0], ob.Object) {
				hud.tarFrame.SetObject(ob)
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

// SetActivePlayer sets specified avatar as active
// player.
func (hud *HUD) SetActivePlayer(pc *game.Player) {
	hud.camera.CenterAt(hud.Game().ActivePlayer().Position())
	hud.pcFrame.SetObject(hud.Game().ActivePlayer())
	hud.Reload()
	// Setup active player area.
	if hud.game == nil {
		return
	}
	chapter := hud.game.Module().Chapter()
	pcArea := chapter.CharacterArea(hud.Game().ActivePlayer().Character)
	if pcArea == nil {
		log.Err.Printf("hud: set active pc: no pc area")
		return
	}
	err := hud.ChangeArea(pcArea)
	if err != nil {
		log.Err.Printf("hud: set active pc: fail to change area: %v", err)
	}
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
	l := hud.layouts[id+serial]
	if l == nil {
		l = NewLayout()
		hud.layouts[id+serial] = l
	}
	return l
}

// Reload reloads HUD layouts for
// active player.
func (hud *HUD) Reload() {
	if hud.Game().ActivePlayer() == nil {
		return
	}
	layout := hud.Layout(hud.Game().ActivePlayer().ID(), hud.Game().ActivePlayer().Serial())
	hud.bar.setLayout(layout)
}

// SetGame sets HUD game.
func (hud *HUD) SetGame(g *game.Game) {
	hud.game = g
	hud.game.SetOnActivePlayerChangeFunc(hud.SetActivePlayer)
	hud.game.SetActivePlayer(hud.game.ActivePlayer())
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
	// Reload HUD.
	hud.Reload()
	hud.runAreaScripts(area)
	return nil
}

// Saves GUI to save struct.
func (hud *HUD) NewGUISave() *res.GUISave {
	sav := new(res.GUISave)
	// Players.
	for _, pc := range hud.Game().Players() {
		pcData := res.PlayerSave{Avatar: pc.Data()}
		// Layout.
		layout := hud.layouts[pc.ID()+pc.Serial()]
		if layout != nil {
			for serialID, slot := range layout.InvSlots() {
				slotSave := res.SlotSave{slot, serialID}
				pcData.InvSlots = append(pcData.InvSlots, slotSave)
			}
			for serialID, slot := range layout.BarSlots() {
				slotSave := res.SlotSave{slot, serialID}
				pcData.BarSlots = append(pcData.BarSlots, slotSave)
			}
		}
		sav.Players = append(sav.Players, pcData)
	}
	// Camera XY position.
	sav.Camera.X = hud.Camera().Position().X
	sav.Camera.Y = hud.Camera().Position().Y
	return sav
}

// LoadGUISave load specified saved GUI state.
func (hud *HUD) LoadGUISave(save *res.GUISave) error {
	// Players.
	for _, pcd := range save.Players {
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
		layoutKey := pcd.Avatar.ID + pcd.Avatar.Serial
		hud.layouts[layoutKey] = layout
	}
	// Camera position.
	hud.camera.SetPosition(pixel.V(save.Camera.X, save.Camera.Y))
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
	mod := hud.Game().Module()
	path := filepath.Join(mod.Conf().Path, data.GUIModulePath, "chapters",
		mod.Chapter().ID(), "areas", a.ID(), "scripts")
	scripts, err := data.ScriptsDir(path)
	if err != nil {
		log.Err.Printf("hud: run area scripts: fail to retrieve scripts: %v", err)
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
	chapter := hud.Game().Module().Chapter()
	pcArea := chapter.CharacterArea(hud.Game().ActivePlayer().Character)
	if pcArea != hud.Camera().Area() {
		go hud.ChangeArea(pcArea)
	}
}

// containsPos checks is specified position is contained
// by any HUD element(except camera).
func (hud *HUD) containsPos(pos pixel.Vec) bool {
	if hud.bar.DrawArea().Contains(pos) ||
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
		(hud.charinfo.Opened() && hud.charinfo.DrawArea().Contains(pos)) {
		return true
	}
	return false
}

// Triggered after right mouse button was pressed.
func (hud *HUD) onMouseRightPressed(pos pixel.Vec) {
}

// Triggered after left mouse button was pressed.
func (hud *HUD) onMouseLeftPressed(pos pixel.Vec) {
}
