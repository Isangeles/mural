/*
 * interpreter.go
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

// ci package provides GUI specific command line tools and
// connection to flame engine commands interpreter.
package ci

import (
	flameci "github.com/isangeles/flame/cmd/ci"
)

const (
	GUI_MAN = "guiman"
)

// Handles specified command,
// returns response code and message.
func HandleCommand(cmd flameci.Command) (int, string) {
	switch cmd.Tool() {
	case GUI_MAN:
		return handleGUICommand(cmd)
	case flameci.ENGINE_MAN:
		//return 2, "tool_unavalible:" + cmd.Tool()
		return flameci.HandleCommand(cmd)
	default:
		return flameci.HandleCommand(cmd)
	}
}

