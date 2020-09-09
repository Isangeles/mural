/*
 * server.go
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

package game

import (
	"fmt"
	"net"

	flameres "github.com/isangeles/flame/data/res"

	"github.com/isangeles/fire/request"
)

// Struct for server connection.
type Server struct {
	net.Conn
}

// NewServer creates new server struct with connection to
// server with specified host and port number.
func NewServer(host, port string) (*Server, error) {
	s := new(Server)
	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("Unable to dial server: %v", err)
	}
	s.Conn = conn
	return s, nil
}

// Login sends login request to the server.
func (s *Server) Login(id, pass string) error {
	login := request.Login{
		ID:   id,
		Pass: pass,
	}
	req := request.Request{Login: []request.Login{login}}
	t, err := request.Marshal(&req)
	if err != nil {
		return fmt.Errorf("Unable to marshal login request: %v",
			err)
	}
	t = fmt.Sprintf("%s\r\n", t)
	_, err = s.Write([]byte(t))
	if err != nil {
		return fmt.Errorf("Unable to write login request: %v", err)
	}
	return nil
}

// NewCharacter sends new character request to the server.
func (s *Server) NewCharacter(charData flameres.CharacterData) error {
	req := request.Request{NewChar: []flameres.CharacterData{charData}}
	t, err := request.Marshal(&req)
	if err != nil {
		return fmt.Errorf("Unable to marshal new char request: %v",
			err)
	}
	t = fmt.Sprintf("%s\r\n", t)
	_, err = s.Write([]byte(t))
	if err != nil {
		return fmt.Errorf("Unable to write new char request: %v", err)
	}
	return nil
}
	
