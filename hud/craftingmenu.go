/*
 * craftingmenu.go
 *
 * Copyright 2019-2020 Dariusz Sikora <dev@isangeles.pl>
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

	"github.com/isangeles/flame/core/data/res/lang"
	"github.com/isangeles/flame/core/module/craft"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/core/data" 
	"github.com/isangeles/mural/log"
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
	bg, err := data.PictureUI("menubg.png")
	if err == nil {
		cm.bgSpr = pixel.NewSprite(bg, bg.Bounds())
	} else {
		log.Err.Printf("hud_crafting:fail_to_retrieve_bg_tex:%v",
			err)
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
	closeButtonBG, err := data.PictureUI("closebutton1.png")
	if err == nil {
		closeBG := pixel.NewSprite(closeButtonBG,
			closeButtonBG.Bounds())
		cm.closeButton.SetBackground(closeBG)
	} else {
		log.Err.Printf("hud_crafting:fail_to_retrieve_close_button_tex:%v",
			err)
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
	makeButtonBG, err := data.PictureUI("button_green.png")
	if err == nil {
		bg := pixel.NewSprite(makeButtonBG, makeButtonBG.Bounds())
		cm.makeButton.SetBackground(bg)
	} else {
		log.Err.Printf("hud_inventory:fail_to_retrieve_make_button_texture:%v", err)
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
	upButtonBG, err := data.PictureUI("scrollup.png")
	if err == nil {
		upBG := pixel.NewSprite(upButtonBG,
			upButtonBG.Bounds())
		cm.recipesList.SetUpButtonBackground(upBG)
	}
	downButtonBG, err := data.PictureUI("scrolldown.png")
	if err == nil {
		downBG := pixel.NewSprite(downButtonBG,
			downButtonBG.Bounds())
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
	closeButtonMove := mtk.ConvVec(pixel.V(cm.Size().X/2 - 20,
		cm.Size().Y/2 - 15))
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
	// Elements.
	cm.closeButton.Update(win)
	cm.makeButton.Update(win)
	cm.recipesList.Update(win)
	cm.recipeInfo.Update(win)
}

// Show toggles menu visibility.
func (cm *CraftingMenu) Show(show bool) {
	cm.opened = show
	if cm.Opened() {
		cm.recipesList.Clear()
		pc := cm.hud.ActivePlayer()
		cm.insertRecipes(pc.Crafting().Recipes()...)
	} else {
		cm.recipeInfo.Clear()
		cm.makeButton.Active(false)
	}
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
	return cm.bgSpr.Frame().Size()
}

// insertRecipes adds all specified recipes to crafting
// recipes list.
func (cm *CraftingMenu) insertRecipes(recipes ...*craft.Recipe) {
	for _, r := range recipes {
		recipeText := lang.Text(r.ID())
		cm.recipesList.AddItem(recipeText, r)
	}
}

// Triggered after close button clicked.
func (cm *CraftingMenu) onCloseButtonClicked(b *mtk.Button) {
	cm.Show(false)
}

// Triggered after selecting recipe from recipes list.
func (cm *CraftingMenu) onRecipeSelected(cs *mtk.CheckSlot) {
	// Retrieve recipe from slot.
	recipe, ok := cs.Value().(*craft.Recipe)
	if !ok {
		log.Err.Printf("hud_crafting:fail to retrieve recipe from list")
		return
	}
	cm.makeButton.Active(true)
	// Show recipe info.
	recipeInfo := lang.Texts(recipe.ID())
	info := recipeInfo[0]
	if len(recipeInfo) > 1 {
		info = fmt.Sprintf("%s\n%s\n", info, recipeInfo[1])
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
	cm.hud.ActivePlayer().Craft(recipe)
	return
}
