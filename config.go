/*
 * config.go
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

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/isangeles/flame/core/data/text"
)

const (
	CONF_FILE_NAME string = ".mural"
)

var (
	fullscreen bool
)

// LoadConfig loads configuration file.
func LoadConfig() error {
	confValues, err := text.ReadConfigValue(CONF_FILE_NAME, "fullscreen")
	if err != nil {
		return err
	}

	fullscreen = confValues[0] == "true"
	return nil
}

// SaveConfig saves current configuration to file
func SaveConfig() error {
	f, err := os.Create(CONF_FILE_NAME)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(fmt.Sprintf("%s\n", "#Mural GUI config file.")) // default header
	w.WriteString(fmt.Sprintf("fullscreen:%v;\n", fullscreen))
	w.Flush()

	return nil
}
