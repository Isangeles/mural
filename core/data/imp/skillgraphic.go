/*
 * skillgraphic.go
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

package imp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/isangeles/mural/core/data/parsexml"
	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/log"
)

var (
	SKILLS_GRAPHIC_FILE_EXT = ".graphic"
)

// ImportSkillsGraphics imports all skills graphics from
// base file with specified path.
func ImportSkillsGraphics(path string) ([]*res.SkillGraphicData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail_to_open_base_file:%v", err)
	}
	skills, err := parsexml.UnmarshalSkillsGraphicsBase(f)
	if err != nil {
		return nil, fmt.Errorf("fail_to_parse_xml:%v", err)
	}
	return skills, nil
}

// ImportSkillsGraphicsDir imports all files with skills graphics from
// directory with specified path.
func ImportSkillsGraphicsDir(path string) ([]*res.SkillGraphicData, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("fail_to_read_dir:%v", err)
	}
	skills := make([]*res.SkillGraphicData, 0)
	for _, finfo := range files {
		if !strings.HasSuffix(finfo.Name(), SKILLS_GRAPHIC_FILE_EXT) {
			continue
		}
		basePath := filepath.FromSlash(path + "/" + finfo.Name())
		impSkills, err := ImportSkillsGraphics(basePath)
		if err != nil {
			log.Err.Printf("data_skills_graphic_import:%s:fail_to_parse_file:%v",
				basePath, err)
			continue
		}
		skills = append(skills, impSkills...)
	}
	return skills, nil
}
