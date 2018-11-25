/*
 * map.go
 *
 * Copyright 2018 Dariusz Sikora <dev@isangeles.pl>
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

// Package for tiled area map.
package areamap

import (
	"fmt"
	"path/filepath"
	"math"
	
	"github.com/salviati/go-tmx/tmx"

	"github.com/faiface/pixel"
	
	"github.com/isangeles/flame/core/module/scenario"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/log"
	"github.com/isangeles/mural/objects"
)

// Struct for graphical representation of area map.
// TODO: implementation is bulky as f**k, slow and don't work well.
type Map struct {
	area       *scenario.Area
	tmxMap     *tmx.Map
	tilesets   map[string]pixel.Picture
	tilesize   pixel.Vec
	mapsize    pixel.Vec
	tilescount pixel.Vec
	// Layers.
	ground []*tile 
}

// NewMap creates new map for specified scenario area.
func NewMap(area *scenario.Area, areasPath string) (*Map, error) {
	m := new(Map)
	m.area = area
	mapsPath := filepath.FromSlash(areasPath + "/maps")
	tm, err := data.Map(area.Id(), areasPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_load_TMX_file:%v", err)
	}
	m.tmxMap = tm
	m.tilesize = pixel.V(float64(m.tmxMap.TileWidth),
		float64(m.tmxMap.TileHeight))
	m.tilescount = pixel.V(float64(m.tmxMap.Width),
		float64(m.tmxMap.Height))
	m.mapsize = pixel.V(float64(int(m.tilesize.X) * m.tmxMap.Width),
		float64(int(m.tilesize.Y) * m.tmxMap.Height))
	m.tilesets = make(map[string]pixel.Picture)
	for _, ts := range m.tmxMap.Tilesets {
		tsPath := filepath.FromSlash(mapsPath + "/" + ts.Image.Source)
		tsPic, err := data.PictureFromDir(tsPath)
		if err != nil {
			return nil, fmt.Errorf("fail_to_retrieve_tilset_source:%v",
				ts.Name)
		}
		m.tilesets[ts.Name] = tsPic
	}
	for _, l := range m.tmxMap.Layers {
		switch l.Name {
		case "ground":
			l, err := m.mapLayer(l, mapsPath)
			if err != nil {
				return nil,
				fmt.Errorf("fail_to_create_ground_layer:%v", err)
			}
			m.ground = l
		default:
			log.Err.Printf("map_builder:unknown_layer:%s", l.Name) 
		}
	}
	return m, nil
}

// Draw draws specified part of map(specific amount of tile starting
// from specific point on map)
// TODO: slow as f**k.
func (m *Map) Draw(win *mtk.Window, startPoint pixel.Vec, size pixel.Vec) {
	drawArea := pixel.R(startPoint.X, startPoint.Y, size.X, size.Y)
	for _, t := range m.ground {
		if drawArea.Contains(t.Position()) {
			t.Draw(win.Window, mtk.Matrix().Moved(pixel.V(
				mtk.ConvSize(t.Position().X),
				mtk.ConvSize(t.Position().Y))))
		}
	}
}

// DrawCircle draws map tiles in circular form(all tiles in specified
// radius from specified position).
func (m *Map) DrawCircle(win *mtk.Window, startPoint pixel.Vec, radius float64) {
	for _, t := range m.ground {
		if mtk.Range(startPoint, t.Position()) <= radius {
			t.Draw(win.Window, pixel.IM.Moved(pixel.V(
				mtk.ConvSize(t.Position().X),
				mtk.ConvSize(t.Position().Y))))
		}
	}
}

// DrawForChar draws only part of map visible for specified game character.
func (m *Map) DrawForChar(win *mtk.Window, startPoint pixel.Vec, size pixel.Vec,
	av *objects.Avatar) {
	drawArea := pixel.R(startPoint.X, startPoint.Y, size.X, size.Y)
	for _, t := range m.ground {
		if drawArea.Contains(t.Position()) &&
			mtk.Range(av.Position(), t.Position()) <= av.SightRange() {
			tilePos := mapDrawPos(t.Position(), startPoint)
			t.Draw(win.Window, mtk.Matrix().Moved(pixel.V(
				mtk.ConvSize(tilePos.X),
				mtk.ConvSize(tilePos.Y))))
		}
	}
}

// mapLayer parses specified TMX layer data to slice
// with tile sprites.
func (m *Map) mapLayer(layer tmx.Layer, mapsPath string) ([]*tile, error) {
	tiles := make([]*tile, 0)
	tileIdX := 1
	tileIdY := 1
	for _, dt := range layer.DecodedTiles {
		tileset := dt.Tileset
		if tileset != nil {
			tilesetPic := m.tilesets[tileset.Name]
			if tilesetPic == nil {
				return nil, fmt.Errorf(
					"fail_to_found_tileset_source:%s",
					tileset.Name)
			}
			tileBounds := m.tileBounds(tilesetPic, dt.ID)
			pic := pixel.NewSprite(tilesetPic, tileBounds)
			tilePos := pixel.V(float64(int(m.tilesize.X)*tileIdX),
				float64(int(m.tilesize.Y)*tileIdY))
			tile := newTile(pic, tilePos)
			tiles = append(tiles, tile)	
			tileIdX++
			if tileIdX > int(m.tilescount.X) {
				tileIdX = 0
				tileIdY++
			}		
		}
		//log.Dbg.Printf("tile[%d]:%v, tileset:%s", i, dt.ID,
		//	dt.Tileset.Name)
	}
	return tiles, nil
}

// TileSize returns size of map tile.
func (m *Map) TileSize() pixel.Vec {
	return m.tilesize
}

// Size returns size of the map.
func (m *Map) Size() pixel.Vec {
	return m.mapsize
}

// tileBounds returns bounds for tile with specified size and ID
// from specified tileset picture.
func (m *Map) tileBounds(tileset pixel.Picture, tileId tmx.ID) pixel.Rect {
	tileCount := 0
	tilesetSize := roundTilesetSize(tileset.Bounds().Max, m.tilesize)
	//fmt.Printf("tileset_size:%v\n", tilesetSize)
	for h := tilesetSize.Y - m.tilesize.Y; h - m.tilesize.Y >= 0; h -= m.tilesize.Y {
		for w := 0.0; w + m.tilesize.X <= tilesetSize.X; w += m.tilesize.X {
			//fmt.Printf("tile_id:%d\n", tileCount)
			if tileCount == int(tileId) {
				tileBounds := pixel.R(w, h, w + m.tilesize.X, h + m.tilesize.Y)
				//fmt.Printf("tile_id:%d, tile_bounds:%v\n", tileCount, tileBounds)
				return tileBounds
			}
			tileCount++
		}
	}
	return pixel.R(0, 0, 0, 0)
}

// mapDrawPos translates real position to map draw position.
func mapDrawPos(pos, drawPos pixel.Vec) pixel.Vec {
	return pixel.V(pos.X - drawPos.X, pos.Y - drawPos.Y)
}

// roundMapSize rounds map size 
func roundTilesetSize(tilesetSize pixel.Vec, tileSize pixel.Vec) pixel.Vec {
	size := pixel.V(0, 0)
	xCount := math.Floor(tilesetSize.X / 32)
	yCount := math.Floor(tilesetSize.Y / 32)
	size.X = tileSize.X * xCount
	size.Y = tileSize.Y * yCount
	return size
}
