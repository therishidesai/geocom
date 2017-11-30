package main

import (
	"encoding/json"
	"net"
)

const (
	MESSAGE_CONNECT    = "CONNECT"
	MESSAGE_DISCONNECT = "DISCONNECT"
	MESSAGE_PUBLIC     = "PUBLIC"
)

type Message struct {
	Kind     string
	Contents string
	Author   string
}

func CreateMessage(kind string, contents string, author string) *Message {
	return &Message{
		Kind:     kind,
		Contents: contents,
		Author:   author,
	}
}

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

func (this *Message) SendTo(conn net.Conn) error {
	enc := json.NewEncoder(conn)
	if err := enc.Encode(this); err != nil {
		return err
	}
	return nil
}
