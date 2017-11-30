package main

import (
	"encoding/json"
	"net"
)

// readFromConn reads an incoming message from a connection
func readFromConn(conn net.Conn) (*Message, error) {
	dec := json.NewDecoder(conn)
	msg := new(Message)
	if err := dec.Decode(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// createConnection establishes a connection between the peer and ip
func createConnection(ip string) (net.Conn, error) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
