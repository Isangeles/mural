/*
 * craftingmenu.go
 *
 * Copyright 2019 Dariusz Sikora <dev@isangeles.pl>
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

	flameconf "github.com/isangeles/flame/config"
	"github.com/isangeles/flame/core/data/text/lang"
	"github.com/isangeles/flame/core/module/object/craft"

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
	opened      bool
	focused     bool
	recipeInfo  *mtk.Textbox
	recipesList *mtk.List
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
	cm.titleText = mtk.NewText(mtk.SIZE_SMALL, 0)
	cm.titleText.SetText(lang.TextDir(flameconf.LangPath(),
		"hud_crafting_title"))
	// Buttons.
	closeButtonParams := mtk.Params{
		Size: mtk.SIZE_MEDIUM,
		Shape: mtk.SHAPE_SQUARE,
		MainColor: accent_color,
	}
	makeButtonParams := mtk.Params{
		Size:      mtk.SIZE_MINI,
		FontSize:  mtk.SIZE_MINI,
		Shape:     mtk.SHAPE_RECTANGLE,
		MainColor: accent_color,
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
	cm.makeButton = mtk.NewButton(makeButtonParams)
	greenButtonBG, err := data.PictureUI("button_green.png")
	if err != nil {
		log.Err.Printf("hud_inventory:fail_to_retrieve_green_button_texture:%v", err)
	} else {
		bg := pixel.NewSprite(greenButtonBG, greenButtonBG.Bounds())
		cm.makeButton.SetBackground(bg)
	}
	cm.makeButton.SetOnClickFunc(cm.onMakeButtonClicked)
	cm.makeButton.SetLabel(lang.TextDir(flameconf.LangPath(), "hud_crafting_make"))
	// Recipe info.
	infoSize := pixel.V(cm.Size().X-mtk.ConvSize(20),
		cm.Size().Y/2-mtk.ConvSize(10))
	recipeInfoParams := mtk.Params{
		SizeRaw:     infoSize,
		FontSize:    mtk.SIZE_MINI,
		MainColor:   main_color,
		AccentColor: accent_color,
	}
	cm.recipeInfo = mtk.NewTextbox(recipeInfoParams)
	// Recipes list.
	recipesSize := pixel.V(cm.Size().X-mtk.ConvSize(20),
		cm.Size().Y/2-mtk.ConvSize(100))
	cm.recipesList = mtk.NewList(recipesSize, mtk.SIZE_MINI,
		main_color, sec_color, accent_color)
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
}

// Show toggles menu visibility.
func (cm *CraftingMenu) Show(show bool) {
	cm.opened = show
	if cm.Opened() {
		cm.recipesList.Clear()
		pc := cm.hud.ActivePlayer()
		cm.insertRecipes(pc.Recipes()...)
	} else {
		cm.recipeInfo.Clear()
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
	mod := cm.hud.game.Module()
	recipesLang := mod.Conf().RecipesLangPath()
	for _, r := range recipes {
		recipeText := lang.AllText(recipesLang, r.ID())
		cm.recipesList.AddItem(recipeText[0], r)
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
	// Show recipe info.
	mod := cm.hud.game.Module()
	recipesLang := mod.Conf().RecipesLangPath()
	recipeInfo := lang.AllText(recipesLang, recipe.ID())
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
	pc := cm.hud.ActivePlayer()
	if !pc.MeetReqs(recipe.Reqs()...) {
		log.Err.Printf(lang.TextDir(flameconf.LangPath(), "reqs_not_meet"))
	}
	pc.ChargeReqs(recipe.Reqs()...)
	res := recipe.Make()
	for _, i := range res {
		pc.Inventory().AddItem(i)
	}
	return
}
