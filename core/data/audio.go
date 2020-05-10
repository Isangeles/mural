/*
 * audio.go
 *
 * Copyright 2019-2020 Dariusz Sikora <dev@isangeles.pl>
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
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

// loadAudiosFromArch open ZIP archive iwth specified path and load all
// audio files from directory inside archive.
// Supported formats: vorbis, wav, mp3.
func loadAudiosFromArch(archPath, dir string) (map[string]*beep.Buffer, error) {
	r, err := zip.OpenReader(archPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open arch: %v", err)
	}
	defer r.Close()
	audio := make(map[string]*beep.Buffer)
	for _, f := range r.File {
		if !isAudio(f) || !strings.HasPrefix(f.Name, dir) {
			continue
		}
		data, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("unable to open arch file: %v", err)
		}
		buf, err := decodeAudio(data, filepath.Ext(f.Name))
		if err != nil {
			return nil, fmt.Errorf("unable to decode audio: %v",
				err)
		}
		audio[filepath.Base(f.Name)] = buf
	}
	return audio, nil
}

// loadAudioFromArch opens ZIP archive with specified path and loads audio
// file with specified path inside ZIP archive.
// Supported formats: vorbis, wav, mp3.
func loadAudioFromArch(archPath, filePath string) (*beep.Buffer, error) {
	r, err := zip.OpenReader(archPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open arch: %v", err)
	}
	defer r.Close()
	for _, f := range r.File {
		if f.Name != filePath {
			continue
		}
		data, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("unable to open arch file: %v", err)
		}
		buf, err := decodeAudio(data, filepath.Ext(f.Name))
		if err != nil {
			return nil, fmt.Errorf("unable to decode audio: %v",
				err)
		}
		return buf, nil
	}
	return nil, fmt.Errorf("file not found: %s", filePath)
}

// loadAudioFromDir load audio file with specified path.
// Supported formats: vorbis, wav, mp3.
func loadAudioFromDir(path string) (*beep.Buffer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()
	buf, err := decodeAudio(file, filepath.Ext(file.Name()))
	if err != nil {
		return nil, fmt.Errorf("unable to decode file: %v", err)
	}
	return buf, nil
}

// decodeAudio decodes specified audio data.
// Audio format need to be specified as a file extension.
// Supported extensions: .ogg, .wav, .mp3
func decodeAudio(buf io.ReadCloser, ext string) (*beep.Buffer, error) {
	switch ext {
	case ".ogg":
		//adata, err := ioutil.ReadAll(rc)
		//if err != nil {
		//	return nil, fmt.Errorf("fail_to_read_audio_data:%v",
		//		err)
		//}
		//ad := mtk.NewAudioData(adata, mtk.Vorbis_audio)
		s, f, err := vorbis.Decode(buf)
		if err != nil {
			return nil, fmt.Errorf("unable to decode vorbis data: %v",
				err)
		}
		ab := beep.NewBuffer(f)
		ab.Append(s)
		return ab, nil
	case ".wav":
		s, f, err := wav.Decode(buf)
		if err != nil {
			return nil, fmt.Errorf("unable to decode wav data: %v",
				err)
		}
		ab := beep.NewBuffer(f)
		ab.Append(s)
		return ab, nil
	case ".mp3":
		s, f, err := mp3.Decode(buf)
		if err != nil {
			return nil, fmt.Errorf("unable to decode mp3 data: %v",
				err)
		}
		ab := beep.NewBuffer(f)
		ab.Append(s)
		return ab, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s",
			ext)
	}
}

// isAudio checks if specified ZIP file is an audio file.
func isAudio(f *zip.File) bool {
	return strings.HasSuffix(f.Name, ".ogg") ||
		strings.HasSuffix(f.Name, ".wav") ||
		strings.HasSuffix(f.Name, "mp3")
}
