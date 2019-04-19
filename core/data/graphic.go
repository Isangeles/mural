/*
 * graphic.go
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

package data

import (
	"archive/zip"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"

	"github.com/faiface/pixel"
)

// loadPicturesFromArch loads all pictures from specified
// directory in ZIP archive with specified path.
func loadPicturesFromArch(archPath, dir string) (map[string]pixel.Picture, error) {
	r, err := zip.OpenReader(archPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_arch:%v", err)
	}
	defer r.Close()
	pics := make(map[string]pixel.Picture, 0)
	for _, f := range r.File {
		if isImage(f) && strings.HasPrefix(f.Name, dir) {
			fPath := strings.Split(f.Name, "/")
			fName := fPath[len(fPath)-1]
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("fail_to_open_arch_file:%v", err)
			}
			defer rc.Close()
			img, _, err := image.Decode(rc)
			if err != nil {
				return nil, fmt.Errorf("fail_to_decode_img:%v", err)
			}
			pics[fName] = pixel.PictureDataFromImage(img)
		}
	}
	return pics, nil
}

// loadFontsFromArch loads all fonts from specified
// directory in ZIP archive with specified path.
func loadFontsFromArch(archPath, dir string) (map[string]*truetype.Font, error) {
	r, err := zip.OpenReader(archPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_arch:%v", err)
	}
	defer r.Close()
	fonts := make(map[string]*truetype.Font, 0)
	for _, f := range r.File {
		if isFont(f) && strings.HasPrefix(f.Name, dir) {
			fPath := strings.Split(f.Name, "/")
			fName := fPath[len(fPath)-1]
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("fail_to_open_arch_file:%v", err)
			}
			defer rc.Close()
			bytes, err := ioutil.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("fail_to_read_arch_file:%v", err)
			}
			font, err := truetype.Parse(bytes)
			if err != nil {
				return nil, fmt.Errorf("fail_to_parse_font:%v", err)
			}
			fonts[fName] = font
		}
	}
	return fonts, nil
}

// loadPictureFromArch loads picture from ZIP archive from specified
// system path.
// Returns error if file with specified path inside archive was not found.
func loadPictureFromArch(archPath, filePath string) (pixel.Picture, error) {
	r, err := zip.OpenReader(archPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_arch:%v", err)
	}
	defer r.Close()
	for _, f := range r.File {
		if f.Name == filePath {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("fail_to_open_arch_file:%v", err)
			}
			defer rc.Close()
			img, _, err := image.Decode(rc)
			if err != nil {
				return nil, fmt.Errorf("fail_to_decode_img:%v", err)
			}
			return pixel.PictureDataFromImage(img), nil
		}
	}
	return nil, fmt.Errorf("file_not_found:%s", filePath)
}

// loadFontFromArch Returns font with specified path in archive in
// specified system path or nil if arch/font was not found.
func loadFontFromArch(archPath, filePath string) (*truetype.Font, error) {
	r, err := zip.OpenReader(archPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_arch:%v", err)
	}
	defer r.Close()
	for _, f := range r.File {
		if f.Name == filePath {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("fail_to_open_arch_file:%v", err)
			}
			defer rc.Close()
			bytes, err := ioutil.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("fail_to_read_arch_file:%v", err)
			}
			font, err := truetype.Parse(bytes)
			if err != nil {
				return nil, fmt.Errorf("fail_to_parse_font:%v", err)
			}
			return font, nil
		}
	}
	return nil, fmt.Errorf("file_not_found:%s", filePath)
}

// loadPictureFromDir loads picture from specified system path and
// returns picture object.
func loadPictureFromDir(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_file:%v", err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("fail_to_decode_img:%v", err)
	}
	return pixel.PictureDataFromImage(img), nil
}

// loadPicturesFromDir loads all pictures from speicified
// directory.x
func loadPicturesFromDir(path string) ([]pixel.Picture, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("fail_to_read_dir:%v", err)
	}
	var pics []pixel.Picture
	for _, f := range files {
		if !f.IsDir() {
			img, err := loadPictureFromDir(
				filepath.FromSlash(path + "/" + f.Name()))
			if err != nil {
				continue
			}
			pics = append(pics, img)
		}
	}
	return pics, nil
}

// loadFontFromDir loads font from specified system path.
func loadFontFromDir(path string) (*truetype.Font, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_file:%v", err)
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("fail_to_read_file:%v", err)
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, fmt.Errorf("fail_to_parse_font:%v", err)
	}
	return font, nil
}

// isImage checks whether specified ZIP file is a image.
func isImage(f *zip.File) bool {
	return strings.HasSuffix(f.Name, ".png") ||
		strings.HasSuffix(f.Name, ".jpg")
}

// isFont checks whether specified ZIP file is a
// font file.
func isFont(f *zip.File) bool {
	return strings.HasSuffix(f.Name, ".ttf")
}
