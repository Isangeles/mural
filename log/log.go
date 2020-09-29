/*
 * log.go
 *
 * Copyright 2018-2020 Dariusz Sikora <dev@isangeles.pl>
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

// Package with loggers for flame engine log.
package log

import (
	"log"

	flamelog "github.com/isangeles/flame/log"
)

var (
	Inf *log.Logger = log.New(flamelog.InfWriter, "mural: ", 0)
	Err *log.Logger = log.New(flamelog.ErrWriter, "mural-error: ", 0)
	Dbg *log.Logger = log.New(flamelog.DbgWriter, "mural-debug: ", 0)
	Cli *log.Logger = log.New(flamelog.InfWriter, "mural-cli: ", 0)
)

// PrintStdOut enables/disables printing log messages
// on the standard out.
func PrintStdOut(print bool) {
	flamelog.PrintStdOut = print
}
