/*
 * infowindow.go
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

package mtk

import (
	"image/color"
	
	"github.com/faiface/pixel"
)

// InfoWindow struct for small text boxes with information about UI
// elements.
type InfoWindow struct {
	*Textbox
	drawArea pixel.Rect
}

// NewInfoWindow creates new information window.
func NewInfoWindow(size Size, color color.Color) *InfoWindow {
	iw := new(InfoWindow)
	iw.Textbox = NewTextbox(pixel.V(0, 0), size, color)
	return iw
}

// Draw draws info window.
func (iw *InfoWindow) Draw(t pixel.Target) {
	iw.Textbox.Draw(iw.drawArea, t)
}

// Update updates info window.
func (iw *InfoWindow) Update(win *Window) {
	iw.drawArea = pixel.R(win.MousePosition().X, win.MousePosition().Y,
		win.MousePosition().X + iw.Size().X,
		win.MousePosition().Y + iw.Size().Y * 1.5)
}
