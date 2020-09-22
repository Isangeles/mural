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
	"bufio"
	"fmt"
	"net"

	flameres "github.com/isangeles/flame/data/res"

	"github.com/isangeles/fire/request"
	"github.com/isangeles/fire/response"

	"github.com/isangeles/mural/log"
)

// Struct for server connection.
type Server struct {
	conn       net.Conn
	authorized bool
	onResponse func(r response.Response)
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
	s.conn = conn
	go s.handleResponses()
	return s, nil
}

// Address returns server address.
func (s *Server) Address() string {
	return s.conn.RemoteAddr().String()
}

// Authorized checks if server connection is authorized.
func (s *Server) Authorized() bool {
	return s.authorized
}

// SetOnResponseFunc sets function triggered on each server response.
func (s *Server) SetOnResponseFunc(f func(r response.Response)) {
	s.onResponse = f
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
	_, err = s.conn.Write([]byte(t))
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
	_, err = s.conn.Write([]byte(t))
	if err != nil {
		return fmt.Errorf("Unable to write new char request: %v", err)
	}
	return nil
}

// handleResponses handles responses from the server connection and
// calls onResponse function for each response.
func (s *Server) handleResponses() {
	out := bufio.NewScanner(s.conn)
	for out.Scan() {
		if out.Err() != nil {
			log.Err.Printf("Server: %s: Unable to read from server: %v",
				s.Address(), out.Err())
			continue
		}
		resp, err := response.Unmarshal(out.Text())
		if err != nil {
			log.Err.Printf("Server: %v: Unable to unmarshal server resonse: %v",
				s.Address(), err)
			continue
		}
		s.authorized = !resp.Logon
		if s.onResponse != nil {
			go s.onResponse(resp)
		}
	}
}
