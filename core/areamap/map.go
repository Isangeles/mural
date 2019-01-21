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
func NewMap(mapData *tmx.Map, mapDir string) (*Map, error) {
	m := new(Map)
	m.tmxMap = mapData
	m.tilesize = pixel.V(float64(m.tmxMap.TileWidth),
		float64(m.tmxMap.TileHeight))
	m.tilescount = pixel.V(float64(m.tmxMap.Width),
		float64(m.tmxMap.Height))
	m.mapsize = pixel.V(float64(int(m.tilesize.X * m.tilescount.X)),
		float64(int(m.tilesize.Y * m.tilescount.Y)))
	m.tilesets = make(map[string]pixel.Picture)
	m.tilesBatches = make(map[pixel.Picture]*pixel.Batch)
	// Tilesets.
	for _, ts := range m.tmxMap.Tilesets {
		tsPath := filepath.FromSlash(mapDir + "/" + ts.Image.Source)
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
				return nil, fmt.Errorf("fail_to_create_ground_layer:%v",
					err)
			}
			m.ground = l
		default:
			fmt.Printf("map_builder:unknown_layer:%s\n", l.Name)
		}
	}
	return m, nil
}

// Draw draws map tiles with positions within specified
// draw area.
// TODO: don't work well.
func (m *Map) Draw(win pixel.Target, matrix pixel.Matrix, size pixel.Vec) {
	drawArea := pixel.R(matrix[4], matrix[5], size.X, size.Y)
	// Clear all tilesets draw batches.
	for _, batch := range m.tilesBatches {
		batch.Clear()
	}
	// Draw layers tiles to tilesets batechs.
	for _, t := range m.ground {
		tileDrawPos := MapDrawPos(t.Position(), matrix)
		if drawArea.Contains(tileDrawPos) {
			batch := m.tilesBatches[t.Picture()]
			if batch == nil {
				continue
			}
			t.Draw(batch, pixel.IM.Scaled(pixel.V(0, 0),
				matrix[0]).Moved(tileDrawPos))
		}
	}
	// Draw bateches with layers tiles.
	for _, batch := range m.tilesBatches {
		batch.Draw(win)
	}
}

// DrawFull draws whole map starting from specified position.
func (m *Map) DrawFull(win pixel.Target, matrix pixel.Matrix) {
	// Clear all tilesets draw batches.
	for _, batch := range m.tilesBatches {
		batch.Clear()
	}
	// Draw layers tiles to tilesets batechs.
	for _, t := range m.ground {
		tileDrawPos := MapDrawPos(t.Position(), matrix)
		batch := m.tilesBatches[t.Picture()]
		if batch == nil {
			continue
		}
		t.Draw(batch, pixel.IM.Scaled(pixel.V(0, 0),
			matrix[0]).Moved(tileDrawPos))
	}
	// Draw bateches with layers tiles.
	for _, batch := range m.tilesBatches {
		batch.Draw(win)
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
	for h := tilesetSize.Y - m.tilesize.Y; h - m.tilesize.Y >= 0;
	h -= m.tilesize.Y {
		for w := 0.0; w + m.tilesize.X <= tilesetSize.X; w += m.tilesize.X {
			if tileCount == int(tileId) {
				tileBounds := pixel.R(w, h, w + m.tilesize.X,
					h + m.tilesize.Y)
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
	tileIdX := 0
	tileIdY := 0
	for _, dt := range layer.DecodedTiles {
		tileset := dt.Tileset
		if tileset != nil {
			tilesetPic := m.tilesets[tileset.Name]
			if tilesetPic == nil {
				return nil, fmt.Errorf("fail_to_found_tileset_source:%s",
					tileset.Name)
			}
			tileBounds := m.tileBounds(tilesetPic, dt.ID)
			pic := pixel.NewSprite(tilesetPic, tileBounds)
			tilePos := pixel.V(float64(int(m.tilesize.X)*tileIdX),
				float64(int(m.tilesize.Y)*tileIdY))
			tile := newTile(pic, tilePos)
			tiles = append(tiles, tile)
			tileIdX++
			if tileIdX > int(m.tilescount.X)-1 {
				tileIdX = 0
				tileIdY++
			}		
		}
	}
	return tiles, nil
}

// MapDrawPos translates real position to map draw position.
func MapDrawPos(pos pixel.Vec, drawMatrix pixel.Matrix) pixel.Vec {
	drawPos := pixel.V(drawMatrix[4], drawMatrix[5]) // 
	drawScale := drawMatrix[0]
	posX := pos.X * drawScale
	posY := pos.Y * drawScale
	drawX := drawPos.X * drawScale
	drawY := drawPos.Y * drawScale
	return pixel.V(posX - drawX, posY - drawY)
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
