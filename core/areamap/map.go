/*
 * map.go
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

// Package for tiled area map.
package areamap

import (
	"fmt"
	"path/filepath"
	"math"
	
	"github.com/salviati/go-tmx/tmx"

	"github.com/faiface/pixel"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/log"
)

// Struct for graphical representation of area map.
type Map struct {
	tmxMap       *tmx.Map
	tilesets     map[string]pixel.Picture
	tilesBatches map[pixel.Picture]*pixel.Batch
	tilesize     pixel.Vec
	mapsize      pixel.Vec
	tilescount   pixel.Vec
	// Layers.
	ground      []*tile
}

// NewMap creates new map for specified scenario area.
func NewMap(tmxDir, tmxName string) (*Map, error) {
	m := new(Map)
	tmxMap, err := data.Map(tmxDir, tmxName)
	if err != nil {
		return nil, fmt.Errorf("fail_to_retrieve_tmx_map:%v", err)
	}
	m.tmxMap = tmxMap
	m.tilesize = pixel.V(float64(m.tmxMap.TileWidth),
		float64(m.tmxMap.TileHeight))
	m.tilescount = pixel.V(float64(m.tmxMap.Width),
		float64(m.tmxMap.Height))
	m.mapsize = pixel.V(float64(int(m.tilesize.X) * m.tmxMap.Width),
		float64(int(m.tilesize.Y) * m.tmxMap.Height))
	m.tilesets = make(map[string]pixel.Picture)
	m.tilesBatches = make(map[pixel.Picture]*pixel.Batch)
	// Tilesets.
	for _, ts := range m.tmxMap.Tilesets {
		tsPath := filepath.FromSlash(tmxDir + "/" + ts.Image.Source)
		tsPic, err := data.PictureFromDir(tsPath)
		if err != nil {
			return nil, fmt.Errorf("fail_to_retrieve_tilset_source:%v",
				ts.Name)
		}
		m.tilesets[ts.Name] = tsPic
		m.tilesBatches[tsPic] = pixel.NewBatch(&pixel.TrianglesData{}, tsPic)
	}
	// Map layers.
	for _, l := range m.tmxMap.Layers {
		switch l.Name {
		case "ground":
			l, err := mapLayer(m, l)
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
func (m *Map) Draw(win *mtk.Window, startPoint pixel.Vec, size pixel.Vec) {
	// Map draw area.
	drawArea := pixel.R(startPoint.X, startPoint.Y, size.X, size.Y)
	// Clear all tilesets draw batches.
	for _, batch := range m.tilesBatches {
		batch.Clear()
	}
	// Draw layers tiles to tilesets batechs.
	for _, t := range m.ground {
		if drawArea.Contains(t.Position()) {
			batch := m.tilesBatches[t.Picture()]
			if batch == nil {
				continue
			}
			tileDrawPos := mapDrawPos(t.Position(), startPoint)
			t.Draw(batch, mtk.Matrix().Moved(tileDrawPos))
		}
	}
	// Draw bateches with layers tiles.
	for _, batch := range m.tilesBatches {
		batch.Draw(win)
	}
}

// DrawForChar draws only part of map visible for specified game character.
func (m *Map) DrawWithFOW(win *mtk.Window, startPoint pixel.Vec, size pixel.Vec,
	sightStart pixel.Vec, sightRange float64) {
	drawArea := pixel.R(startPoint.X, startPoint.Y, size.X, size.Y)
	for _, t := range m.ground {
		if drawArea.Contains(t.Position()) &&
			mtk.Range(sightStart, t.Position()) <= sightRange {
			tileDrawPos := mapDrawPos(t.Position(), startPoint)
			t.Draw(win.Window, mtk.Matrix().Moved(tileDrawPos))
		}
	}
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
	for h := tilesetSize.Y - m.tilesize.Y; h - m.tilesize.Y >= 0;
	h -= m.tilesize.Y {
		for w := 0.0; w + m.tilesize.X <= tilesetSize.X; w += m.tilesize.X {
			//fmt.Printf("tile_id:%d\n", tileCount)
			if tileCount == int(tileId) {
				tileBounds := pixel.R(w, h, w + m.tilesize.X,
					h + m.tilesize.Y)
				//fmt.Printf("tile_id:%d, tile_bounds:%v\n", tileCount, tileBounds)
				return tileBounds
			}
			tileCount++
		}
	}
	return pixel.R(0, 0, 0, 0)
}

// mapLayer parses specified TMX layer data to slice
// with tile sprites.
func mapLayer(m *Map, layer tmx.Layer) ([]*tile, error) {
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
	}
	return tiles, nil
}

// mapDrawPos translates real position to map draw position.
func mapDrawPos(pos, drawPos pixel.Vec) pixel.Vec {
	return pixel.V(mtk.ConvSize(pos.X - drawPos.X),
		mtk.ConvSize(pos.Y - drawPos.Y))
}

// roundTilesetSize rounds tileset size to to value that can be divided
// by tile size without rest.
func roundTilesetSize(tilesetSize pixel.Vec, tileSize pixel.Vec) pixel.Vec {
	size := pixel.V(0, 0)
	xCount := math.Floor(tilesetSize.X / tileSize.X)
	yCount := math.Floor(tilesetSize.Y / tileSize.Y)
	size.X = tileSize.X * xCount
	size.Y = tileSize.Y * yCount
	return size
}
