/*
 * character_test.go
 *
 * Copyright 2023 Dariusz Sikora <ds@isangeles.dev>
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

package game

import (
	"testing"

	"github.com/isangeles/flame"
	"github.com/isangeles/flame/character"
	"github.com/isangeles/flame/data/res"
)

// TestCharSetPosition tests setting position for the game character.
func TestCharSetPosition(t *testing.T) {
	// Create game.
	mod := flame.NewModule(res.ModuleData{})
	game := New(mod)
	// Create character.
	char := character.New(res.CharacterData{ID: "char", Level: 1})
	gameChar := NewCharacter(char, game)
	// Test.
	gameChar.SetPosition(10, 10)
	x, y := gameChar.Position()
	if x != 10 || y != 10 {
		t.Errorf("Character position invalid: %fx%f != 10x10",
			x, y)
	}
	x, y = gameChar.DestPoint()
	if x != 10 || y != 10 {
		t.Errorf("Character destination point invalid: %fx%f != 10x10",
			x, y)
	}
}
