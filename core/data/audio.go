/*
 * audio.go
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

package data

import (
	"archive/zip"
	"fmt"
	"strings"
	"os"
	
	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/faiface/beep/vorbis"
)

// loadAudioFromArch loads audio file from specified directory
// in ZIP archive from specified path.
// Supported formats: vorbis, wav.
func loadAudioFromArch(archPath, filePath string) (beep.Streamer, beep.Format, error) {
	r, err := zip.OpenReader(archPath)
	if err != nil {
		return nil, beep.Format{}, err
	}
	//defer r.Close()
	for _, f := range r.File {
		if f.Name == filePath {
			rc, err := f.Open()
			if err != nil {
				return nil, beep.Format{}, err
			}
			//defer rc.Close()
			if strings.HasSuffix(f.Name, ".ogg") { // vorbis
				s, format, err := vorbis.Decode(rc)
				if err != nil {
					return nil, beep.Format{}, err
				}
				return s, format, nil
			} else if strings.HasSuffix(f.Name, ".wav") { // wav
				s, format, err := wav.Decode(rc)
				if err != nil {
					return nil, beep.Format{}, err
				}
				return s, format, nil
			} else {
				return nil, beep.Format{}, fmt.Errorf("unsupported format:%s",
					f.Name)
			}
		}
	}
	return nil, beep.Format{}, fmt.Errorf("file not found:%s", filePath)
}

// loadAudioFromDir load audio file with specified path.
// Supported formats: vorbis, wav.
func loadAudioFromDir(path string) (beep.Streamer, beep.Format, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, beep.Format{}, err
	}
	//defer file.Close()
	if strings.HasSuffix(path, ".ogg") { // vorbis
		s, format, err := vorbis.Decode(file)
		if err != nil {
			return nil, beep.Format{}, err
		}
		return s, format, nil
	} else if strings.HasSuffix(path, ".wav") { // wav
		s, format, err := wav.Decode(file)
		if err != nil {
			return nil, beep.Format{}, err
		}
		return s, format, nil
	} else {
		return nil, beep.Format{}, fmt.Errorf("unsupported format:%s",
			path)
	}
}
