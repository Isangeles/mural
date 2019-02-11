/*
 * main.go
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

// Example for creating simple MTK animation.
package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/isangeles/mural/core/mtk"
)

// Main function.
func main() {
	// Run Pixel graphic.
	pixelgl.Run(run)
}

// All window code fired from there.
func run() {
	// Create Pixel window configuration.
	cfg := pixelgl.WindowConfig{
		Title:  "MTK animation example",
		Bounds: pixel.R(0, 0, 1600, 900),
	}
	// Create MTK warpper for Pixel window.
	win, err := mtk.NewWindow(cfg)
	if err != nil {
		panic(fmt.Errorf("fail_to_create_mtk_window:%v", err))
	}
	// Load spritesheet image.
	ss, err := loadPicture("spritesheet.png")
	if err != nil {
		panic(fmt.Errorf("fail_to_load_spritesheet:%v", err))
	}
	// Retrieve frames from spritesheet.
	frames := cutFrames(ss)
	// Create animation from frames,
	// 2 frames per second.
	anim := mtk.NewAnimation(frames, 2)
	// Main loop.
	for !win.Closed() {
		// Clear window.
		win.Clear(colornames.Black)
		// Draw.
		animPos := win.Bounds().Center()
		anim.Draw(win.Window, mtk.Matrix().Moved(animPos))
		// Update.
		win.Update()
		anim.Update(win)
	}
}

// loadPicture loads picture from file with specified path.
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_image_file:%v", err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("fail_to_decode_image:%v", err)
	}
	return pixel.PictureDataFromImage(img), nil
}

// cutFrames retrieves all frames from 100x100
// spritesheet(4 50x50 frames).
func cutFrames(ss pixel.Picture) []*pixel.Sprite {
	frame1 := pixel.NewSprite(ss, pixel.R(0, 0, 50, 50))
	frame2 := pixel.NewSprite(ss, pixel.R(50, 0, 100, 50))
	frame3 := pixel.NewSprite(ss, pixel.R(0, 50, 50, 100))
	frame4 := pixel.NewSprite(ss, pixel.R(50, 50, 100, 100))
	return []*pixel.Sprite{frame1, frame2, frame3, frame4}
}
