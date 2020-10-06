/*
 * response.go
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

package mainmenu

import (
	"github.com/isangeles/flame/module"
	"github.com/isangeles/flame/module/serial"

	"github.com/isangeles/fire/response"

	"github.com/isangeles/mural/log"
)

// handleResponse handles specified response from Fire server.
func (mm *MainMenu) handleResponse(resp response.Response) {
	if !resp.Logon {
		mm.handleUpdateResponse(resp.Update)
	}
	for _, r := range resp.Error {
		log.Err.Printf("Login menu: server error: %v", r)
	}
}

// handleUpdateResponse handles update response.
func (mm *MainMenu) handleUpdateResponse(resp response.Update) {
	serial.Reset()
	if mm.mod == nil {
		mm.mod = module.New()
		mm.mod.Apply(resp.Module)
		return
	}
	mm.mod.Apply(resp.Module)
}
