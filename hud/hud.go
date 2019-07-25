/*
 * hud.go
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

// Package with 'head-up display' elements.
package hud

import (
	"fmt"
	"path/filepath"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	flamecore "github.com/isangeles/flame/core"
	"github.com/isangeles/flame/core/data/text/lang"
	flameobject "github.com/isangeles/flame/core/module/object"
	"github.com/isangeles/flame/core/module/object/character"
	"github.com/isangeles/flame/core/module/scenario"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/core/areamap"
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/object"
	"github.com/isangeles/mural/log"
)

var (
	// HUD colors.
	mainColor   = colornames.Grey
	secColor    = colornames.Blue
	accentColor = colornames.Red
	// Keys.
	pauseKey    = pixelgl.KeySpace
	menuKey     = pixelgl.KeyEscape
	chatKey     = pixelgl.KeyGraveAccent
	invKey      = pixelgl.KeyB
	skillsKey   = pixelgl.KeyK
	journalKey  = pixelgl.KeyL
	craftingKey = pixelgl.KeyV
)

// Struct for 'head-up display'.
type HUD struct {
	loadScreen *LoadingScreen
	camera     *Camera
	bar        *MenuBar
	menu       *Menu
	savemenu   *SaveMenu
	pcFrame    *ObjectFrame
	tarFrame   *ObjectFrame
	castBar    *CastBar
	chat       *Chat
	inv        *InventoryMenu
	skills     *SkillMenu
	loot       *LootWindow
	dialog     *DialogWindow
	journal    *JournalWindow
	crafting   *CraftingMenu
	trade      *TradeWindow
	training   *TrainingWindow
	game       *flamecore.Game
	pcs        []*object.Avatar
	activePC   *object.Avatar
	userFocus  *mtk.Focus
	msgs       *mtk.MessagesQueue
	layouts    map[string]*Layout
	loading    bool
	exiting    bool
	loaderr    error
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
	if hud.ActivePlayer() == nil { // no active pc, don't draw
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
	tradePos := win.Bounds().Center()
	trainPos := win.Bounds().Center()
	// Draw elements.
	hud.camera.Draw(win)
	hud.bar.Draw(win, mtk.Matrix().Moved(barPos))
	hud.chat.Draw(win, mtk.Matrix().Moved(chatPos))
	hud.pcFrame.Draw(win, mtk.Matrix().Moved(pcFramePos))
	if hud.ActivePlayer().Targets()[0] != nil {
		hud.tarFrame.Draw(win, mtk.Matrix().Moved(tarFramePos))
	}
	if hud.menu.Opened() {
		hud.menu.Draw(win, mtk.Matrix().Moved(menuPos))
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
	if hud.trade.Opened() {
		hud.trade.Draw(win, mtk.Matrix().Moved(tradePos))
	}
	if hud.training.Opened() {
		hud.training.Draw(win, mtk.Matrix().Moved(trainPos))
	}
	if hud.ActivePlayer().Casting() {
		hud.castBar.Draw(win, mtk.Matrix().Moved(castBarPos))
	}
	// Messages.
	msgPos := win.Bounds().Center()
	hud.msgs.Draw(win, mtk.Matrix().Moved(msgPos))
}

// Update updated HUD elements.
func (hud *HUD) Update(win *mtk.Window) {
	// HUD state.
	if hud.loading {
		if hud.loaderr != nil { // on loading error
			log.Err.Printf("hud_loading_fail:%v", hud.loaderr)
			hud.Exit()
		}
	}
	if hud.ActivePlayer() == nil { // no active pc, don't update
		return
	}
	// Key events.
	if win.JustPressed(chatKey) {
		// Toggle chat activity.
		if !hud.chat.Activated() {
			hud.chat.Active(true)
			hud.camera.Lock(true)
		} else {
			hud.chat.Active(false)
			hud.camera.Lock(false)
		}
	}
	if !hud.chat.Activated() { // block rest of key events if chat is active
		if win.JustPressed(pauseKey) {
			// Pause game.
			if !hud.game.Paused() {
				hud.game.Pause(true)
			} else {
				hud.game.Pause(false)
			}
		}
		if win.JustPressed(menuKey) {
			// Show menu.
			if !hud.menu.Opened() {
				hud.menu.Show(true)
			} else {
				hud.menu.Show(false)
			}
		}
		if win.JustPressed(invKey) {
			// Show inventory.
			if !hud.inv.Opened() {
				hud.inv.Show(true)
			} else {
				hud.inv.Show(false)
			}
		}
		if win.JustPressed(skillsKey) {
			// Show skills.
			if !hud.skills.Opened() {
				hud.skills.Show(true)
			} else {
				hud.skills.Show(false)
			}
		}
		if win.JustPressed(journalKey) {
			// Show journal.
			if !hud.journal.Opened() {
				hud.journal.Show(true)
			} else {
				hud.journal.Show(false)
			}
		}
		if win.JustPressed(craftingKey) {
			// Show crafting menu.
			if !hud.crafting.Opened() {
				hud.crafting.Show(true)
			} else {
				hud.crafting.Show(false)
			}
		}
	}
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		hud.onMouseLeftPressed(win.MousePosition())
	}
	if win.JustPressed(pixelgl.MouseButtonRight) {
		hud.onMouseRightPressed(win.MousePosition())
	}
	// PC target.
	if hud.ActivePlayer().Targets()[0] != nil {
		for _, av := range hud.camera.Avatars() {
			if flameobject.Equals(hud.ActivePlayer().Targets()[0], av.Character) {
				hud.tarFrame.SetObject(av)
			}
		}
		for _, ob := range hud.camera.AreaObjects() {
			if flameobject.Equals(hud.ActivePlayer().Targets()[0], ob.Object) {
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
	if hud.menu.Opened() {
		hud.menu.Update(win)
	}
	if hud.savemenu.Opened() {
		hud.savemenu.Update(win)
	}
	if hud.inv.Opened() {
		hud.inv.Update(win)
	}
	if hud.skills.Opened() {
		hud.skills.Update(win)
	}
	if hud.loot.Opened() {
		hud.loot.Update(win)
	}
	if hud.dialog.Opened() {
		hud.dialog.Update(win)
	}
	if hud.journal.Opened() {
		hud.journal.Update(win)
	}
	if hud.crafting.Opened() {
		hud.crafting.Update(win)
	}
	if hud.trade.Opened() {
		hud.trade.Update(win)
	}
	if hud.training.Opened() {
		hud.training.Update(win)
	}
	hud.msgs.Update(win)
}

// Players returns player characters
// (game characters of HUD user).
func (hud *HUD) Players() []*object.Avatar {
	return hud.pcs
}

// ActivePlayers retruns currently active PC.
func (hud *HUD) ActivePlayer() *object.Avatar {
	return hud.activePC
}

// AddAvatar adds specified character to HUD
// player characters list.
func (hud *HUD) AddPlayer(char *character.Character) error {
	avData := res.Avatar(char.ID())
	if avData == nil {
		return fmt.Errorf("fail_to_find_avatar_data:%s", char.ID())
	}
	av := object.NewAvatar(char, avData)
	hud.pcs = append(hud.pcs, av)
	if hud.ActivePlayer() == nil {
		hud.SetActivePlayer(av)
	}
	return nil
}

// SetActivePlayer sets specified avatar as active
// player.
func (hud *HUD) SetActivePlayer(pc *object.Avatar) {
	hud.activePC = pc
	hud.camera.CenterAt(hud.ActivePlayer().Position())
	hud.pcFrame.SetObject(hud.ActivePlayer())
	hud.castBar.SetOwner(hud.ActivePlayer().Character)
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
func (hud *HUD) Game() *flamecore.Game {
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
	// TODO: sometimes SetLoadInfo causes panic on text draw.
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
	if hud.ActivePlayer() == nil {
		return
	}
	layout := hud.Layout(hud.ActivePlayer().ID(), hud.ActivePlayer().Serial())
	hud.bar.setLayout(layout)
}

// SetGame sets HUD game.
func (hud *HUD) SetGame(g *flamecore.Game) error {
	hud.game = g
	// Players.
	if len(hud.game.Players()) < 1 {
		return fmt.Errorf("no player characters")
	}
	for _, pc := range hud.game.Players() {
		err := hud.AddPlayer(pc)
		if err != nil {
			return fmt.Errorf("fail_to_add_player:%v", err)
		}
	}
	// Setup active player area.
	chapter := hud.game.Module().Chapter()
	pcArea := chapter.CharacterArea(hud.ActivePlayer().Character)
	if pcArea == nil {
		hud.loaderr = fmt.Errorf("no pc area")
		return hud.loaderr
	}
	err := hud.ChangeArea(pcArea)
	if err != nil {
		hud.loaderr = fmt.Errorf("fail_to_change_area:%v", err)
		return hud.loaderr
	}
	return nil
}

// ChangeArea changes current HUD area.
func (hud *HUD) ChangeArea(area *scenario.Area) error {
	// Map.
	hud.OpenLoadingScreen(lang.Text("gui", "load_map_info"))
	defer hud.CloseLoadingScreen()
	chapter := hud.game.Module().Chapter()
	mapsPath := filepath.FromSlash(chapter.Conf().ModulePath +
		"/gui/chapters/" + chapter.Conf().ID + "/areas/maps")
	tmxMap, err := data.Map(mapsPath, area.ID())
	if err != nil {
		hud.loaderr = fmt.Errorf("fail_to_retrieve_tmx_map:%v", err)
		return hud.loaderr
	}
	areaMap, err := areamap.NewMap(tmxMap, mapsPath)
	if err != nil {
		hud.loaderr = fmt.Errorf("fail_to_create_pc_area_map:%v", err)
		return hud.loaderr
	}
	hud.camera.SetMap(areaMap)
	// Avatars.
	hud.OpenLoadingScreen(lang.Text("gui", "load_avatars_info"))
	avatars := make([]*object.Avatar, 0)
	for _, c := range area.Characters() {
		var pcAvatar *object.Avatar
		for _, pc := range hud.Players() {
			if c == pc.Character {
				pcAvatar = pc
				break
			}
		}
		if pcAvatar != nil { // skip players, PCs already have avatars
			avatars = append(avatars, pcAvatar)
			continue
		}
		avData := res.Avatar(c.ID())
		if avData == nil {
			log.Err.Printf("hud:area_change:avatar_data_not_found:%s", c.ID())
			continue
		}
		av := object.NewAvatar(c, avData)
		avatars = append(avatars, av)
	}
	hud.camera.SetAvatars(avatars)
	// Objects.
	hud.OpenLoadingScreen(lang.Text("gui", "load_objects_info"))
	objects := make([]*object.ObjectGraphic, 0)
	for _, o := range area.Objects() {
		ogData := res.Object(o.ID())
		if ogData == nil {
			log.Err.Printf("hud:area_change:object_data_not_found:%s", o.ID())
			continue
		}
		og := object.NewObjectGraphic(o, ogData)
		objects = append(objects, og)
	}
	hud.camera.SetObjects(objects)
	// Reload HUD.
	hud.Reload()
	return nil
}

// Saves GUI to save struct.
func (hud *HUD) NewGUISave() *res.GUISave {
	sav := new(res.GUISave)
	// Players.
	for _, pc := range hud.Players() {
		pcData := new(res.PlayerSave)
		pcData.Avatar = pc.Data()
		// Layout.
		layout := hud.layouts[pc.SerialID()]
		if layout != nil {
			pcData.InvSlots = layout.InvSlots()
			pcData.BarSlots = layout.BarSlots()
		}
		sav.PlayersData = append(sav.PlayersData, pcData)
	}
	// Camera XY position.
	sav.CameraPosX = hud.Camera().Position().X
	sav.CameraPosY = hud.Camera().Position().Y
	return sav
}

// LoadGUISave load specified saved GUI state.
func (hud *HUD) LoadGUISave(save *res.GUISave) error {
	// Players.
	for _, pcd := range save.PlayersData {
		layout := NewLayout()
		layout.SetInvSlots(pcd.InvSlots)
		layout.SetBarSlots(pcd.BarSlots)
		layoutKey := pcd.Avatar.ID + pcd.Avatar.Serial
		hud.layouts[layoutKey] = layout
	}
	// Camera position.
	hud.camera.SetPosition(pixel.V(save.CameraPosX, save.CameraPosY))
	// Reload UI.
	hud.Reload()
	return nil
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
		(hud.training.Opened() && hud.training.DrawArea().Contains(pos)) {
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
