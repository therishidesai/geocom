package main

import (
	"encoding/json"
	"net"
)

const (
	MESSAGE_CONNECT    = "CONNECT"
	MESSAGE_DISCONNECT = "DISCONNECT"
	MESSAGE_PUBLIC     = "PUBLIC"
	MESSAGE_SHOW       = "SHOW"
)

// Message is a struct that represents a message
type Message struct {
	Kind     string
	Contents string
	Author   string
}

// CreateMessage creates a new message
func CreateMessage(kind string, contents string, author string) *Message {
	return &Message{
		Kind:     kind,
		Contents: contents,
		Author:   author,
	}
}

// Send sends the message to all the connections
func (this *Message) Send(peers map[string]net.Conn) error {
	for name, conn := range peers {
		if name != this.Author {
			if err := this.SendTo(conn); err != nil {
				return err
			}
		}
	}
	return nil
}

// SendTo sends the message to a specific connection
func (this *Message) SendTo(conn net.Conn) error {
	enc := json.NewEncoder(conn)
	if err := enc.Encode(this); err != nil {
		return err
	}
	return nil
}
