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

// Package for area map.
package areamap

import (
	"fmt"
	"path/filepath"
	
	"github.com/salviati/go-tmx/tmx"

	"github.com/faiface/pixel"
	
	"github.com/isangeles/flame/core/module/scenario"

	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/core/mtk"
	"github.com/isangeles/mural/log"
)

// Struct for graphical representation of area map.
// TODO: implementation is bulky as f**k, slow and don't work well.
type Map struct {
	area     *scenario.Area
	tmxMap   *tmx.Map
	tilesets map[string]pixel.Picture
	tilesize pixel.Vec
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
			l, err := m.mapLayer(l, mapsPath, m.tmxMap.TileWidth,
				m.tmxMap.TileHeight, m.tmxMap.Width, m.tmxMap.Height)
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
			t.Draw(win.Window, mtk.Matrix())
		}
	}
}

// mapLayer parses specified TMX layer data to slice
// with tile sprites.
func (m *Map) mapLayer(layer tmx.Layer, mapsPath string,
	tileWidth, tileHeight, mapWidth, mapHeight int) ([]*tile, error) {
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
			tileWidth := float64(tileWidth)
			tileHeight := float64(tileHeight)
			tileBounds := tileBounds(tilesetPic, pixel.V(
				tileWidth, tileHeight), dt.ID)
			pic := pixel.NewSprite(tilesetPic, tileBounds)
			tilePos := pixel.V(float64(int(tileWidth)*tileIdX),
				float64(int(tileHeight)*tileIdY))
			tile := newTile(pic, tilePos)
			tiles = append(tiles, tile)	
			tileIdX++
			if tileIdX > mapWidth {
				tileIdX = 0
				tileIdY++
			}		
		}
		//log.Dbg.Printf("tile[%d]:%v, tileset:%s", i, dt.ID,
		//	dt.Tileset.Name)
	}
	return tiles, nil
}

// TileSize return size of map tile.
func (m *Map) TileSize() pixel.Vec {
	return m.tilesize
}

// tileBounds returns bounds for tile with specified size and ID
// from specified tileset picture.
func tileBounds(tileset pixel.Picture, tileSize pixel.Vec,
	tileId tmx.ID) pixel.Rect {
	tileCount := 0
	for h := 0.0; h <= tileset.Bounds().H(); h += tileSize.Y {
		for w := 0.0; w <= tileset.Bounds().W(); w += tileSize.X {
			if tileCount == int(tileId) {
				return pixel.R(w, h, tileSize.X, tileSize.Y)
			}
			tileCount++
		}
	}
	return pixel.R(0, 0, 0, 0)
}
