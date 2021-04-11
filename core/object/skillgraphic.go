/*
 * skillgraphic.go
 *
 * Copyright 2019-2021 Dariusz Sikora <dev@isangeles.pl>
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

package object

import (
	"github.com/faiface/pixel"

	"github.com/faiface/beep"

	"github.com/isangeles/flame/skill"

	"github.com/isangeles/mural/core/data/res"
	"github.com/isangeles/mural/core/data/res/audio"
	"github.com/isangeles/mural/core/data/res/graphic"
)

// Graphical wrapper for skills.
type SkillGraphic struct {
	*skill.Skill
	icon                pixel.Picture
	activationAudio     *beep.Buffer
	iconName            string
	activationAudioName string
	activationAnim      AvatarAnimType
}

// NewSkillGraphic creates new graphical wrapper for specified skill.
func NewSkillGraphic(skill *skill.Skill, data *res.SkillGraphicData) *SkillGraphic {
	sg := new(SkillGraphic)
	sg.Skill = skill
	// Icon.
	sg.icon = graphic.Icons[data.Icon]
	if sg.icon != nil {
		sg.iconName = data.Icon
	}
	// Activation effect.
	sg.activationAudio = audio.Effects[data.ActivationAudio]
	if sg.activationAudio != nil {
		sg.activationAudioName = data.ActivationAudio
	}
	sg.activationAnim = AvatarAnimType(data.ActivationAnim)
	return sg
}

// Icon returns skill icon sprite.
func (sg *SkillGraphic) Icon() pixel.Picture {
	return sg.icon
}

// AudioEffect returns skill audio effect.
func (sg *SkillGraphic) ActivationAudio() *beep.Buffer {
	return sg.activationAudio
}

// ActivationAnim returns skill activation animation.
func (sg *SkillGraphic) ActivationAnim() AvatarAnimType {
	return sg.activationAnim
}
