/*
 * scripts.go
 *
 * Copyright 2021 Dariusz Sikora <dev@isangeles.pl>
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

// data package contains functions for loading
// graphic and audio data.
package data

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/isangeles/burn/ash"

	"github.com/isangeles/mural/log"
)

var (
	ashScriptExt = ".ash"
)

// ScriptsDir returns all scripts from directory with
// specified path.
func ScriptsDir(path string) ([]*ash.Script, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir: %v", err)
	}
	scripts := make([]*ash.Script, 0)
	for _, info := range files {
		if !strings.HasSuffix(info.Name(), ashScriptExt) {
			continue
		}
		scriptPath := filepath.FromSlash(path + "/" + info.Name())
		s, err := Script(scriptPath)
		if err != nil {
			log.Err.Printf("data: %s: unable to retrieve script: %v",
				path, err)
			continue
		}
		scripts = append(scripts, s)
	}
	return scripts, nil
}

// Script parses file with specified path to
// Ash scirpt.
func Script(path string) (*ash.Script, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	text, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %v", err)
	}
	scriptName := filepath.Base(path)
	script, err := ash.NewScript(scriptName, fmt.Sprintf("%s", text))
	if err != nil {
		return nil, fmt.Errorf("unable to parse script text: %v", err)
	}
	return script, nil
}
