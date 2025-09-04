/*
 * effectgraphic.go
 *
 * Copyright 2019-2025 Dariusz Sikora <ds@isangeles.dev>
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

package object

import (
	"fmt"

	"golang.org/x/image/colornames"

	"github.com/gopxl/pixel"

	"github.com/isangeles/flame/effect"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/data/res/graphic"
	"github.com/isangeles/mural/log"
)

var (
	timeLabelColor    = colornames.Red
	defaultEffectIcon = "unknown.png"
)

// Graphical wrapper for effects.
type EffectGraphic struct {
	*effect.Effect
	icon     *pixel.Sprite
	timeText *mtk.Text
}

// NewEffectGraphic creates new graphical wrapper for specified effect.
func NewEffectGraphic(effect *effect.Effect, data *res.EffectGraphicData) *EffectGraphic {
	eg := new(EffectGraphic)
	eg.Effect = effect
	// Icon.
	iconPic := graphic.Icons[data.Icon]
	if iconPic != nil {
		eg.icon = pixel.NewSprite(iconPic, iconPic.Bounds())
	} else {
		log.Err.Printf("effect graphic: %s#%s: icon not found: %s", effect.ID(),
			effect.Serial(), data.Icon)
		iconPic = graphic.Icons[defaultEffectIcon]
	}
	if !eg.Infinite() {
		// Time text.
		textParams := mtk.Params{
			FontSize: mtk.SizeBig,
		}
		eg.timeText = mtk.NewText(textParams)
		eg.timeText.SetColor(timeLabelColor)
	}
	return eg
}

// DrawIcon draws effect icon and text label with
// remaining time(in seconds).
func (eg *EffectGraphic) DrawIcon(t pixel.Target, matrix pixel.Matrix) {
	eg.icon.Draw(t, matrix)
	if eg.timeText != nil {
		eg.timeText.SetText(fmt.Sprintf("%d", eg.Time()/1000))
		eg.timeText.Draw(t, matrix)
	}
}

// Icon returns effect icon.
func (eg *EffectGraphic) Icon() *pixel.Sprite {
	return eg.icon
}
