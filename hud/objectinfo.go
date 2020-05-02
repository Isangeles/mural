/*
 * objectinfo.go
 *
 * Copyright 2020 Dariusz Sikora <dev@isangeles.pl>
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

package hud

import (
	"fmt"
	
	"github.com/faiface/pixel"
	
	"github.com/isangeles/mtk"

	"github.com/isangeles/mural/config"
)

// Struct for HUD object info window.
type ObjectInfo struct {
	hud    *HUD
	info   *mtk.InfoWindow
	object InfoTarget
	opened bool
}

// Interface for object info targets.
type InfoTarget interface {
	ID() string
	Serial() string
	Name() string
}

// newObjectInfo creates object info for HUD.
func newObjectInfo(hud *HUD) *ObjectInfo {
	oi := new(ObjectInfo)
	oi.hud = hud
	infoParams := mtk.Params{
		FontSize:  mtk.SizeSmall,
		MainColor: pixel.RGBA{0.1, 0.1, 0.1, 0.5},
	}
	oi.info = mtk.NewInfoWindow(infoParams)
	return oi
}

// Draw draws info window.
func (oi *ObjectInfo) Draw(win *mtk.Window) {
	oi.info.Draw(win)
}

// Update updates object info.
func (oi *ObjectInfo) Update(win *mtk.Window) {
	oi.info.Update(win)
	oi.object = nil
	for _, av := range oi.hud.Camera().Avatars() {
		if !av.Hovered() {
			continue
		}
		oi.SetObject(av)
		oi.Open(true)
		return
	}
	for _, ob := range oi.hud.Camera().AreaObjects() {
		if !ob.Hovered() {
			continue
		}
		oi.SetObject(ob)
		oi.Open(true)
		return
	}
	if oi.object == nil {
		oi.Open(false)
	}
}

// Opened checks if object info is open.
func (oi *ObjectInfo) Opened() bool {
	return oi.opened
}

// SetOpen toggles info window visibility.
func (oi *ObjectInfo) Open(open bool) {
	oi.opened = open
}

// SetObject sets specified object for object info.
func (oi *ObjectInfo) SetObject(ob InfoTarget) {
	oi.object = ob
	oi.info.SetText(objectInfo(oi.object))
}

// objectInfo returns info text for specified
// object.
func objectInfo(o InfoTarget) string {
	info := fmt.Sprintf("%s", o.Name())
	if config.Debug {
		info = fmt.Sprintf("%s\n%s#%s", info, o.ID(),
			o.Serial())
	}
	return info
}


