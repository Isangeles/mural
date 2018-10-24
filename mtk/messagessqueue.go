/*
 * messagesqueue.go
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

package mtk

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// MessagesQueue struct for list with messages to display.
type MessagesQueue struct{
	queue []*MessageWindow
	focus *Focus
}

// NewMessagesQueue creates new messages queue.
func NewMessagesQueue(focus *Focus) *MessagesQueue {
	mq := new(MessagesQueue)
	mq.queue = make([]*MessageWindow, 0)
	mq.focus = focus
	return mq
}

// Draw draws all messages 
func (mq *MessagesQueue) Draw(t pixel.Target, matrix pixel.Matrix) {
	for _, m := range mq.queue {
		if m.Opened() {
			m.Draw(t, matrix)
		}
	}
}

// Update updates all messages in queue.
func (mq *MessagesQueue) Update(win *pixelgl.Window) {
	for i, m := range mq.queue {
		if m.Opened() {
			if i == len(mq.queue)-1 {
				m.Active(true)
				mq.focus.Focus(m)
			} else {
				m.Active(false)
			}
			m.Update(win)
		}
		if m.Dismissed() {
			mq.Remove(i)
		}
	}
}

// Append adds specified message to the front of queue.
func (mq *MessagesQueue) Append(m *MessageWindow) {
	mq.queue = append(mq.queue, m)
}

// Remove removes message with specified index from queue.
func (mq *MessagesQueue) Remove(i int) {
	mq.queue = append(mq.queue[:i], mq.queue[i+1:]...)
}

