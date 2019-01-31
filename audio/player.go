/*
 * player.go
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

package audio

import (
	"time"
	
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	
	"github.com/isangeles/mural/core/data"
	"github.com/isangeles/mural/log"
)

// Struct for music player.
type Player struct {
	playlist []*data.AudioData
	playID   int
	control  *beep.Ctrl
}

// NewPlayer creates new audio player for specified
// stream format.
func NewPlayer(format beep.Format) *Player {
	ap := new(Player)
	ap.playlist = make([]*data.AudioData, 0)
	ap.control = &beep.Ctrl{}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	return ap
}

// Add adds specified audio data to playlist.
func (p *Player) Add(a *data.AudioData) {
	p.playlist = append(p.playlist, a)
}

// SetPlaylist sets specified slice with audio data
// as player playlist.
func (p *Player) SetPlaylist(playlist []*data.AudioData) {
	p.playlist = playlist
}

// Play starts player.
func (p *Player) Play() {
	if p.playlist[p.playID] == nil {
		log.Err.Printf("audio_player:current playlist position nil")
		return
	}
	m := p.playlist[p.playID]
	p.control.Streamer = m.Stream
	speaker.Play(p.control)
}

// Stop stops player.
func (p *Player) Stop() {
	// TODO: don't work.
	if p.control.Streamer == nil {
		return
	}
	speaker.Lock()
	p.control.Streamer = nil
	speaker.Unlock()
}

// Reset stops player and moves play index to
// first playlist index.
func (p *Player) Reset() {
	p.Stop()
	p.SetPlayIndex(0)
}

// Next moves play index to next position
// on player playlist.
func (p *Player) Next() {
	p.Stop()
	p.SetPlayIndex(p.playID+1)
	p.Play()
}

// Prev moves play index to previous position
// on player playlist.
func (p *Player) Prev() {
	p.Stop()
	p.SetPlayIndex(p.playID-1)
	p.Play()
}

// Clear clears player playlist.
func (p *Player) Clear() {
	p.playlist = make([]*data.AudioData, 0)
}

// SetPlayIndex sets specified index as current index
// on player playlist.
// If specified value is bigger than playlist lenght
// then first index is set, if is lower than 0 then
// last index is set.
func (p *Player) SetPlayIndex(id int) {
	switch {
	case id > len(p.playlist)-1:
		p.playID = 0
	case id < 0:
		p.playID = len(p.playlist)-1
	default:
		p.playID = id
	}
}
