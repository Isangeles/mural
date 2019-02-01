/*
 * audioplayer.go
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

package mtk

import (
	"fmt"
	"time"
	
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

// Struct for audio player.
type AudioPlayer struct {
	playlist  []beep.Streamer
	playID    int
	mixer     *beep.Mixer
	ctrlMusic *beep.Ctrl
}

// NewAudioPlayer creates new audio player for specified
// stream format.
func NewAudioPlayer(format beep.Format) *AudioPlayer {
	p := new(AudioPlayer)
	p.playlist = make([]beep.Streamer, 0)
	p.mixer = new(beep.Mixer)
	p.ctrlMusic = new(beep.Ctrl)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(p.mixer)
	return p
}

// Add adds specified audio stream to playlist.
func (p *AudioPlayer) AddMusic(ab *beep.Buffer) error { 
	s := ab.Streamer(0, ab.Len())
	p.playlist = append(p.playlist, s)
	return nil
}

// SetPlaylist sets specified slice with audio streams
// as player playlist.
func (p *AudioPlayer) SetPlaylist(playlist []beep.Streamer) {
	p.playlist = playlist
}

// Play starts player.
func (p *AudioPlayer) PlayMusic() error {
	if p.playID < 0 || p.playID > len(p.playlist)-1 {
		return fmt.Errorf("audio_player:current playlist position nil")
	}
	m := p.playlist[p.playID]
	p.ctrlMusic.Streamer = m
	p.mixer.Play(p.ctrlMusic)
	return nil
}

// Play starts playing specified audio stream.
func (p *AudioPlayer) Play(ab *beep.Buffer) error {
	s := ab.Streamer(0, ab.Len())
	p.mixer.Play(s)
	return nil
}

// Stop stops player.
func (p *AudioPlayer) StopMusic() {
	if p.ctrlMusic.Streamer == nil {
		return
	}
	speaker.Lock()
	p.ctrlMusic.Streamer = nil
	speaker.Unlock()
}

// Reset stops player and moves play index to
// first music playlist index.
func (p *AudioPlayer) Reset() {
	p.StopMusic()
	p.SetPlayIndex(0)
}

// Next moves play index to next position
// on music playlist.
func (p *AudioPlayer) Next() {
	p.StopMusic()
	p.SetPlayIndex(p.playID+1)
	p.PlayMusic()
}

// Prev moves play index to previous position
// on music playlist.
func (p *AudioPlayer) Prev() {
	p.StopMusic()
	p.SetPlayIndex(p.playID-1)
	p.PlayMusic()
}

// Clear clears music playlist.
func (p *AudioPlayer) Clear() {
	p.playlist = make([]beep.Streamer, 0)
}

// SetPlayIndex sets specified index as current index
// on music playlist.
// If specified value is bigger than playlist lenght
// then first index is set, if is lower than 0 then
// last index is set.
func (p *AudioPlayer) SetPlayIndex(id int) {
	switch {
	case id > len(p.playlist)-1:
		p.playID = 0
	case id < 0:
		p.playID = len(p.playlist)-1
	default:
		p.playID = id
	}
}
