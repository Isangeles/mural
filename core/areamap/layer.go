/*
 * layer.go
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

package areamap

import (
	"fmt"
	
	"github.com/salviati/go-tmx/tmx"

	"github.com/faiface/pixel"
)

// Struct for map layer.
type Layer struct {
	name    string
	tiles   []*Tile
}

// newLayer creates new layer with tiles for specified map.
func newLayer(m *Map, tmxLayer tmx.Layer) (*Layer, error) {
	l := new(Layer)
	l.name = tmxLayer.Name
	l.tiles = make([]*Tile, 0)
	tileX := 0
	tileY := 0
	for _, dt := range tmxLayer.DecodedTiles {
		tileset := dt.Tileset
		if tileset != nil {
			tilesetPic := m.tilesets[tileset.Name]
			if tilesetPic == nil {
				return nil, fmt.Errorf("fail to found tileset source: %s",
					tileset.Name)
			}
			tileBounds := m.tileBounds(tilesetPic, dt.ID)
			pic := pixel.NewSprite(tilesetPic, tileBounds)
			tilePos := pixel.V(float64(int(m.tilesize.X)*tileX),
				float64(int(m.tilesize.Y)*tileY))
			tilePos.Y = m.Size().Y - tilePos.Y
			tile := newTile(pic, tilePos)
			l.tiles = append(l.tiles, tile)		
		}
		tileX++
		if tileX > int(m.tilescount.X)-1 {
			tileX = 0
			tileY++
		}
	}
	return l, nil
}

// Name returns layer name from tmx data.
func (l *Layer) Name() string {
	return l.name
}
