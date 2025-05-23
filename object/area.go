/*
 * area.go
 *
 * Copyright 2023-2025 Dariusz Sikora <ds@isangeles.dev>
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
	"sync"

	"github.com/isangeles/flame/area"
	"github.com/isangeles/flame/character"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"

	"github.com/isangeles/stone"

	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
	"github.com/isangeles/mural/data"
	"github.com/isangeles/mural/data/res"
	"github.com/isangeles/mural/game"
	"github.com/isangeles/mural/log"
)

const (
	dayColorAlpha     = 0.0
	eveningColorAlpha = 0.5
	nightColorAlpha   = 0.9
)

var (
	fowColor pixel.RGBA = pixel.RGBA{0.1, 0.1, 0.1, 0.7}
	dayColor pixel.RGBA = pixel.RGBA{0.1, 0.1, 0.1, dayColorAlpha}
)

// Graphical wrapper for area.
type Area struct {
	*area.Area
	game    *game.Game
	areaMap *stone.Map
	fow     *imdraw.IMDraw
	avatars *sync.Map
}

// NewArea returns new graphical wrapper for specified area.
func NewArea(game *game.Game, area *area.Area, areaMap *stone.Map) *Area {
	a := &Area{
		Area:    area,
		game:    game,
		areaMap: areaMap,
		fow:     imdraw.New(nil),
		avatars: new(sync.Map),
	}
	a.updateObjects()
	return a
}

// Draw draws area.
func (a *Area) Draw(win *mtk.Window, matrix pixel.Matrix, size pixel.Vec) {
	// Map.
	if config.MapFull {
		a.areaMap.Draw(win.Window, matrix)
	} else {
		a.areaMap.DrawPart(win.Window, matrix, size)
	}
	// Avatars.
	for _, av := range a.Avatars() {
		if !a.game.VisibleForPlayer(av.Position().X, av.Position().Y) {
			continue
		}
		avPos := a.convAreaPos(av.Position(), matrix)
		av.Draw(win, mtk.Matrix().Moved(avPos))
	}
	// FOW effect.
	if a.areaMap != nil && config.MapFOW {
		a.drawMapFOW(win.Window, matrix)
	}
	// Day.
	dayColor.A = a.dayTransparency()
	mtk.DrawRect(win, win.Bounds(), dayColor)
}

// Update updates area.
func (a *Area) Update(win *mtk.Window) {
	a.updateObjects()
	// Avatars.
	for _, av := range a.Avatars() {
		if a.game.VisibleForPlayer(av.Position().X, av.Position().Y) {
			av.Silence(false)
			av.Update(win)
		} else {
			av.Silence(true)
		}
	}
}

// Avatars returns all avatars from current area.
func (a *Area) Avatars() (avatars []*Avatar) {
	addAvatar := func(k, v interface{}) bool {
		av, ok := v.(*Avatar)
		if !ok {
			return true
		}
		avatars = append(avatars, av)
		return true
	}
	a.avatars.Range(addAvatar)
	return
}

// Map returns area map.
func (a *Area) Map() *stone.Map {
	return a.areaMap
}

// PassablePosition checks if specified position is 'passable',
// i.e. map there is visible layer on this position where player
// is allowed to move(like 'ground' layer').
func (a *Area) PassablePosition(pos pixel.Vec) bool {
	layer := a.areaMap.PositionLayer(pos)
	if layer == nil {
		return false
	}
	return layer.Name() == "ground"
}

// updateObjects checks and creates new avatars for area
// objects if needed.
func (a *Area) updateObjects() {
	a.clearObjects()
	// Add new objects & characters.
	for _, ob := range a.Objects() {
		char, isChar := ob.(*character.Character)
		if !isChar {
			continue
		}
		_, avExists := a.avatars.Load(char.ID() + char.Serial())
		if avExists {
			continue
		}
		gameChar := a.game.Char(char.ID(), char.Serial())
		if gameChar == nil {
			continue
		}
		// Avatar.
		avData := res.Avatar(char.ID())
		if avData == nil {
			defData := data.DefaultAvatarData(char)
			res.Avatars = append(res.Avatars, defData)
			og, err := NewAvatar(gameChar, &defData)
			if err != nil {
				log.Err.Printf("Area objects update: unable to create default avatar: %s %s: %v",
					char.ID(), char.Serial(), err)
				continue
			}
			a.avatars.Store(char.ID()+char.Serial(), og)
			continue
		}
		av, err := NewAvatar(gameChar, avData)
		if err != nil {
			log.Err.Printf("Area objects update: unable to create avatar: %s %s: %v",
				char.ID(), char.Serial(), err)
			continue
		}
		a.avatars.Store(char.ID()+char.Serial(), av)
	}
}

// clearObjecs removes avatars & graphics without
// corresponding objects in the area.
func (a *Area) clearObjects() {
	for _, av := range a.Avatars() {
		found := false
		for _, ob := range a.Objects() {
			_, isChar := ob.(*character.Character)
			if isChar && ob.ID() == av.ID() && ob.Serial() == av.Serial() {
				found = true
				break
			}
		}
		if found {
			continue
		}
		a.avatars.Delete(av.ID() + av.Serial())
	}
}

// drawMapFOW draws 'Fog Of War' effect on current area map.
func (a *Area) drawMapFOW(t pixel.Target, matrix pixel.Matrix) {
	a.fow.Clear()
	w, h := 0.0, 0.0
	for h < a.areaMap.Size().Y {
		if !a.game.VisibleForPlayer(w, h) {
			// Draw FOW tile.
			tileSizeX := mtk.ConvSize(a.areaMap.TileSize().X)
			tileSizeY := mtk.ConvSize(a.areaMap.TileSize().Y)
			tileDrawMin := a.convAreaPos(pixel.V(w, h), matrix)
			tileDrawMax := pixel.V(tileDrawMin.X+tileSizeX,
				tileDrawMin.Y+tileSizeY)
			a.fow.Color = fowColor
			a.fow.Push(tileDrawMin)
			a.fow.Push(tileDrawMax)
			a.fow.Rectangle(0)
		}
		// Next tile.
		w += a.areaMap.TileSize().X
		if w > a.areaMap.Size().X {
			w = 0.0
			h += a.areaMap.TileSize().Y
		}
	}
	a.fow.Draw(t)
}

// convAreaPos translates specified area
// position to camera position.
func (a *Area) convAreaPos(pos pixel.Vec, matrix pixel.Matrix) pixel.Vec {
	drawPos := pixel.V(matrix[4], matrix[5])
	drawScale := matrix[0]
	posX := pos.X * drawScale
	posY := pos.Y * drawScale
	drawX := drawPos.X //* drawScale
	drawY := drawPos.Y //* drawScale
	return pixel.V(posX-drawX, posY-drawY)
}

// dayTransparency return transparency level for
// current phase of the day in the area.
func (a *Area) dayTransparency() float64 {
	hour := a.Time.Hour()
	switch {
	case hour > 17:
		return eveningColorAlpha
	case hour > 21 || hour < 5:
		return nightColorAlpha
	default:
		return dayColorAlpha
	}
}
