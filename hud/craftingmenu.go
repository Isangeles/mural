/*
 * craftingmenu.go
 *
 * Copyright 2019-2021 Dariusz Sikora <dev@isangeles.pl>
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

package hud

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/flame/data/res/lang"
	"github.com/isangeles/flame/craft"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/data/res/graphic"
	"github.com/isangeles/mural/log"
)

var (
	craftingKey = pixelgl.KeyV
)

// Struct for HUD crafting menu.
type CraftingMenu struct {
	hud         *HUD
	bgSpr       *pixel.Sprite
	bgDraw      *imdraw.IMDraw
	drawArea    pixel.Rect
	titleText   *mtk.Text
	closeButton *mtk.Button
	makeButton  *mtk.Button
	recipeInfo  *mtk.Textbox
	recipesList *mtk.List
	opened      bool
	focused     bool
}

// newCraftingMenu creates new crafting
// menu for HUD.
func newCraftingMenu(hud *HUD) *CraftingMenu {
	cm := new(CraftingMenu)
	cm.hud = hud
	// Background.
	cm.bgDraw = imdraw.New(nil)
	bg := graphic.Textures["menubg.png"]
	if bg != nil {
		cm.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	} else {
		log.Err.Printf("hud: crafting menu: unable to retrieve bg texture")
	}
	// Title.
	titleParams := mtk.Params{
		FontSize: mtk.SizeSmall,
	}
	cm.titleText = mtk.NewText(titleParams)
	cm.titleText.SetText(lang.Text("hud_crafting_title"))
	// Close button.
	closeButtonParams := mtk.Params{
		Size: mtk.SizeMedium,
		Shape: mtk.ShapeSquare,
		MainColor: accentColor,
	}
	cm.closeButton = mtk.NewButton(closeButtonParams)
	closeButtonBG := graphic.Textures["closebutton1.png"]
	if closeButtonBG != nil {
		closeBG := pixel.NewSprite(closeButtonBG, closeButtonBG.Bounds())
		cm.closeButton.SetBackground(closeBG)
	} else {
		log.Err.Printf("hud: crafting menu: unable to retrieve close button texture")
	}
	cm.closeButton.SetOnClickFunc(cm.onCloseButtonClicked)
	// Make button.
	makeButtonParams := mtk.Params{
		Size:      mtk.SizeMini,
		FontSize:  mtk.SizeMini,
		Shape:     mtk.ShapeRectangle,
		MainColor: accentColor,
	}
	cm.makeButton = mtk.NewButton(makeButtonParams)
	makeButtonBG := graphic.Textures["button_green.png"]
	if makeButtonBG != nil {
		bg := pixel.NewSprite(makeButtonBG, makeButtonBG.Bounds())
		cm.makeButton.SetBackground(bg)
	} else {
		log.Err.Printf("hud: crafting menu: unable to retrieve make button texture")
	}
	cm.makeButton.SetOnClickFunc(cm.onMakeButtonClicked)
	cm.makeButton.SetLabel(lang.Text("hud_crafting_make"))
	// Recipe info.
	infoSize := pixel.V(cm.Size().X-mtk.ConvSize(20),
		cm.Size().Y/2-mtk.ConvSize(10))
	recipeInfoParams := mtk.Params{
		SizeRaw:     infoSize,
		FontSize:    mtk.SizeSmall,
		MainColor:   mainColor,
		AccentColor: accentColor,
	}
	cm.recipeInfo = mtk.NewTextbox(recipeInfoParams)
	// Recipes list.
	recipesSize := pixel.V(cm.Size().X-mtk.ConvSize(20),
		cm.Size().Y/2-mtk.ConvSize(100))
	recipesParams := mtk.Params{
		SizeRaw:     recipesSize,
		MainColor:   mainColor,
		SecColor:    secColor,
		AccentColor: accentColor,
	}
	cm.recipesList = mtk.NewList(recipesParams)
	upButtonBG := graphic.Textures["scrollup.png"]
	if upButtonBG != nil {
		upBG := pixel.NewSprite(upButtonBG, upButtonBG.Bounds())
		cm.recipesList.SetUpButtonBackground(upBG)
	}
	downButtonBG := graphic.Textures["scrolldown.png"]
	if downButtonBG != nil {
		downBG := pixel.NewSprite(downButtonBG, downButtonBG.Bounds())
		cm.recipesList.SetDownButtonBackground(downBG)
	}
	cm.recipesList.SetOnItemSelectFunc(cm.onRecipeSelected)
	return cm
}

// Draw draws menu.
func (cm *CraftingMenu) Draw(win *mtk.Window, matrix pixel.Matrix) {
	// Draw area.
	cm.drawArea = mtk.MatrixToDrawArea(matrix, cm.Size())
	// Background.
	if cm.bgSpr != nil {
		cm.bgSpr.Draw(win, matrix)
	} else {
		mtk.DrawRectangle(win, cm.DrawArea(), nil)
	}
	// Title.
	titleTextMove := pixel.V(0, cm.Size().Y/2-mtk.ConvSize(25))
	cm.titleText.Draw(win, matrix.Moved(titleTextMove))
	// Buttons.
	closeButtonMove := pixel.V(cm.Size().X/2 - mtk.ConvSize(20),
		cm.Size().Y/2 - mtk.ConvSize(15))
	makeButtonMove := mtk.ConvVec(pixel.V(0, -60))
	cm.closeButton.Draw(win, matrix.Moved(closeButtonMove))
	cm.makeButton.Draw(win, matrix.Moved(makeButtonMove))
	// Recipe info.
	recipeInfoMove := mtk.MoveTC(cm.Size(), cm.recipeInfo.Size())
	recipeInfoMove.Y -= mtk.ConvSize(50)
	cm.recipeInfo.Draw(win, matrix.Moved(recipeInfoMove))
	// Recipes list.
	recipesMove := mtk.MoveBC(cm.Size(), cm.recipesList.Size())
	recipesMove.Y += mtk.ConvSize(20)
	cm.recipesList.Draw(win, matrix.Moved(recipesMove))
}

// Update updates menu.
func (cm *CraftingMenu) Update(win *mtk.Window) {
	// Key events.
	if !cm.hud.Chat().Activated() && win.JustPressed(craftingKey) {
		if cm.Opened() {
			cm.Hide()
		} else {
			cm.Show()
		}
	}
	// Elements.
	if cm.Opened() {
		cm.closeButton.Update(win)
		cm.makeButton.Update(win)
		cm.recipesList.Update(win)
		cm.recipeInfo.Update(win)
	}
}

// Show shows menu.
func (cm *CraftingMenu) Show() {
	cm.opened = true
	cm.recipesList.Clear()
	pc := cm.hud.Game().ActivePlayer()
	cm.insertRecipes(pc.Crafting().Recipes()...)
}

// Hide hides menu.
func (cm *CraftingMenu) Hide() {
	cm.opened = false
	cm.recipeInfo.Clear()
	cm.makeButton.Active(false)
}

// Opened checks if menu is open.
func (cm *CraftingMenu) Opened() bool {
	return cm.opened
}

// DrawArea returns current draw area.
func (cm *CraftingMenu) DrawArea() pixel.Rect {
	return cm.drawArea
}

// Size returns menu background size.
func (cm *CraftingMenu) Size() pixel.Vec {
	if cm.bgSpr == nil {
		return mtk.ConvVec(pixel.V(50, 100))
	}
	return mtk.ConvVec(cm.bgSpr.Frame().Size())
}

// insertRecipes adds all specified recipes to crafting
// recipes list.
func (cm *CraftingMenu) insertRecipes(recipes ...*craft.Recipe) {
	for _, r := range recipes {
		cm.recipesList.AddItem(lang.Text(r.ID()), r)
	}
}

// Triggered after close button clicked.
func (cm *CraftingMenu) onCloseButtonClicked(b *mtk.Button) {
	cm.Hide()
}

// Triggered after selecting recipe from recipes list.
func (cm *CraftingMenu) onRecipeSelected(cs *mtk.CheckSlot) {
	// Retrieve recipe from slot.
	recipe, ok := cs.Value().(*craft.Recipe)
	if !ok {
		log.Err.Printf("hud: crafting menu: unable to retrieve recipe from list")
		return
	}
	cm.makeButton.Active(true)
	// Show recipe info.
	nameInfo := lang.Texts(recipe.ID())
	info := fmt.Sprintf("%s", nameInfo[0])
	if len(nameInfo) > 1 {
		info = fmt.Sprintf("%s\n%s\n", info, nameInfo[1])
	}
	cm.recipeInfo.SetText(info)
}

// Triggered after make button clicked.
func (cm *CraftingMenu) onMakeButtonClicked(b *mtk.Button) {
	val := cm.recipesList.SelectedValue()
	if val == nil {
		return
	}
	recipe, ok := val.(*craft.Recipe)
	if !ok {
		return
	}
	cm.hud.Game().ActivePlayer().Use(recipe)
	return
}
