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
	"os"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

// loadAudiosFromArch loads all audio files data from specified
// directory insied ZIP archive with specified path.x
func loadAudiosFromArch(archPath, path string) (map[string]*beep.Buffer, error) {
	r, err := zip.OpenReader(archPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_archive_reader:%v", err)
	}
	defer r.Close()
	return nil, fmt.Errorf("unsupported yet")
}

// loadAudioFromArch loads audio file from specified directory
// in ZIP archive from specified path.
// Supported formats: vorbis, wav, mp3.
func loadAudioFromArch(archPath, filePath string) (*beep.Buffer, error) {
	r, err := zip.OpenReader(archPath)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_arch:%v", err)
	}
	defer r.Close()
	for _, f := range r.File {
		if f.Name == filePath {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("fail_to_open_file_inside_arch:%v",
					err)
			}
			defer rc.Close()
			switch {
			case strings.HasSuffix(f.Name, ".ogg"): // vorbis
				//adata, err := ioutil.ReadAll(rc)
				//if err != nil {
				//	return nil, fmt.Errorf("fail_to_read_audio_data:%v",
				//		err)
				//}
				//ad := mtk.NewAudioData(adata, mtk.Vorbis_audio)
				s, f, err := vorbis.Decode(rc)
				if err != nil {
					return nil, fmt.Errorf("fail_to_decode_vorbis_data:%v",
						err)
				}
				ab := beep.NewBuffer(f)
				ab.Append(s)
				return ab, nil
			case strings.HasSuffix(f.Name, ".wav"): // wav
				s, f, err := wav.Decode(rc)
				if err != nil {
					return nil, fmt.Errorf("fail_to_decode_wav_data:%v",
						err)
				}
				ab := beep.NewBuffer(f)
				ab.Append(s)
				return ab, nil
			case strings.HasSuffix(f.Name, ".mp3"): // mp3
				s, f, err := mp3.Decode(rc)
				if err != nil {
					return nil, fmt.Errorf("fail_to_decode_mp3_data:%v",
						err)
				}
				ab := beep.NewBuffer(f)
				ab.Append(s)
				return ab, nil
			default:
				return nil, fmt.Errorf("unsupported format:%s",
					f.Name)
			}
		}
	}
	return nil, fmt.Errorf("file_not_found:%s", filePath)
}

// loadAudioFromDir load audio file with specified path.
// Supported formats: vorbis, wav, mp3.
func loadAudioFromDir(path string) (*beep.Buffer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_file:%v", err)
	}
	defer file.Close()
	switch {
	case strings.HasSuffix(path, ".ogg"): // vorbis
		s, f, err := vorbis.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("fail_to_decode_vorbis_data:%v",
				err)
		}
		ab := beep.NewBuffer(f)
		ab.Append(s)
		return ab, nil
	case strings.HasSuffix(path, ".wav"): // wav
		s, f, err := wav.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("fail_to_decode_wav_data:%v",
				err)
		}
		ab := beep.NewBuffer(f)
		ab.Append(s)
		return ab, nil
	case strings.HasSuffix(path, ".mp3"): // mp3
		s, f, err := mp3.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("fail_to_decode_mp3_data:%v",
				err)
		}
		ab := beep.NewBuffer(f)
		ab.Append(s)
		return ab, nil
	default:
		return nil, fmt.Errorf("unsupported_format:%s",
			path)
	}
}
