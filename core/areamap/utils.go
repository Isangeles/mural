/*
 * utils.go
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

package areamap

import (
	"fmt"
	"image"
	_ "image/png"
	"math"
	"os"
	
	"github.com/salviati/go-tmx/tmx"

	"github.com/faiface/pixel"
)

// tmxMap retieves tiled map from file with specified path.
func tmxMap(path string) (*tmx.Map, error) {
	tmxFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail to open tmx file: %v", err)
	}
	tmxMap, err := tmx.Read(tmxFile)
	if err != nil {
		return nil, fmt.Errorf("fail to read tmx file: %v", err)
	}
	return tmxMap, nil
}

// picture retieves picture from file with specified path.
func picture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail to open file: %v", err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("fail to decode img: %v", err)
	}
	return pixel.PictureDataFromImage(img), nil
}

// mapDrawPos translates real position to map draw position.
func mapDrawPos(pos pixel.Vec, drawMatrix pixel.Matrix) pixel.Vec {
	drawPos := pixel.V(drawMatrix[4], drawMatrix[5])
	drawScale := drawMatrix[0]
	posX := pos.X * drawScale
	posY := pos.Y * drawScale
	drawX := drawPos.X //* drawScale
	drawY := drawPos.Y //* drawScale
	return pixel.V(posX-drawX, posY-drawY)
}

// roundTilesetSize rounds tileset size to to value that can be divided
// by tile size without remainder.
func roundTilesetSize(tilesetSize pixel.Vec, tileSize pixel.Vec) pixel.Vec {
	size := pixel.V(0, 0)
	xCount := math.Floor(tilesetSize.X / tileSize.X)
	yCount := math.Floor(tilesetSize.Y / tileSize.Y)
	size.X = tileSize.X * xCount
	size.Y = tileSize.Y * yCount
	return size
}
