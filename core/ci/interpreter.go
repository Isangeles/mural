/*
 * interpreter.go
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

// ci package provides GUI specific command line tools
// for Burn CI.
package ci

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/isangeles/burn"
	"github.com/isangeles/burn/ash"

	"github.com/isangeles/mural/log"
)

const (
	GUI_MAN = "guiman"
)

// On init.
func init() {
	burn.AddToolHandler(GUI_MAN, handleGUICommand)
}

// RunScriptsDir runs in background all Ash scripts
// in directory with specified path.
func RunScriptsDir(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("fail_to_read_dir:%v", err)
	}
	for _, finfo := range files {
		if !strings.HasSuffix(finfo.Name(), ash.SCRIPT_FILE_EXT) {
			continue
		}
		filepath := filepath.FromSlash(path + "/" + finfo.Name())
		err := RunScript(filepath)
		if err != nil {
			log.Err.Printf("ci:script:%s:%v", err)
		}
	}
	return nil
}

// RunScript runs in background Asg script
// from specified path.
func RunScript(path string, args ...string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("fail_to_open_file:%v", err)
	}
	text, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("fail_to_read_file:%v", err)
	}
	script, err := ash.NewScript(fmt.Sprintf("%s", text), args...)
	if err != nil {
		return fmt.Errorf("fail_to_create_ash_script:%v", err)
	}
	go runScript(script)
	return nil
}

// runScript executes specified script,
// in case of error sends err message to
// Mural log.
func runScript(s *ash.Script) {
	err := ash.Run(s)
	if err != nil {
		log.Err.Printf("ci:fail_to_run_script:%v", err)
		return
	}
}
