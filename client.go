package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
)

const PORT = 5000

type Client struct {
	connections map[string]net.Conn
	nick        string
	server      bool
}

func CreateClient(nick string, isServer bool) *Client {
	return &Client{
		connections: make(map[string]net.Conn),
		nick:        nick,
		server:      isServer,
	}
}

// StartServer initializes the client as a server
func (this *Client) StartServer() error {
	if this.server {
		IP := fmt.Sprintf(":%d", PORT)
		tcpAddr, err := net.ResolveTCPAddr("tcp", IP)
		if err != nil {
			return err
		}

		listener, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			return err
		}

		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go this.handleConnection(conn)
		}
	} else {
		return errors.New("Cannot start server as client.")
	}
}

// handleConnection handles an incoming connection
func (this *Client) handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr).IP
	fmt.Printf("[*] Handling connection from %s\n", remoteAddr)
	go this.receive(conn)
}

// receieve continuously receieves messages from a connection
func (this *Client) receive(conn net.Conn) {
	for {
		msg, err := readFromConn(conn)
		if err != nil {
			return
		}
		switch msg.Kind {
		case MESSAGE_CONNECT:
			fmt.Printf("[*] Initializing connection to %s\n", msg.Author)
			err := this.connectToPeer(msg, conn)
			if err != nil {
				fmt.Printf("[*] Failed to connect to %s\n", msg.Author)
			} else {
				fmt.Printf("[*] Established connection to %s\n", msg.Author)
			}
		case MESSAGE_PUBLIC:
			// Server receives request to send message to everyone
			fmt.Printf("[*] %s said: %s\n", msg.Author, msg.Contents)
			msg.Kind = MESSAGE_SHOW
			msg.Send(this.connections)
		case MESSAGE_SHOW:
			// Client receieves request from server to show a message
			fmt.Printf("[*] %s said: %s\n", msg.Author, msg.Contents)
		default:
			fmt.Printf("[*] Bad message.")
		}
	}
}

// ConnectToServer connects the server to the client
func (this *Client) connectToPeer(msg *Message, conn net.Conn) error {
	this.connections[msg.Author] = conn
	response := CreateMessage(MESSAGE_CONNECT, "", this.nick)
	response.SendTo(conn)
	return nil
}

// ConnectToServer connects the client to the server
func (this *Client) connectToServer(ip string) error {
	conn, err := createConnection(ip)
	if err != nil {
		return err
	}

	msg := CreateMessage(MESSAGE_CONNECT, "", this.nick)
	msg.SendTo(conn)

	msg, err = readFromConn(conn)
	if err != nil {
		return err
	}
	this.connections[msg.Author] = conn
	go this.receive(conn)
	return nil
}

// HandleInput reads user input, encodes it as a message, and sends it
func (this *Client) handleInput() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Send message: ")
		text, _ := reader.ReadString('\n')
		if this.server {
			// We are the server, so send the message to all peers
			fmt.Printf("[*] %s said: %s\n", this.nick, text)
			msg := CreateMessage(MESSAGE_SHOW, text, this.nick)
			msg.Send(this.connections)
		} else {
			// We are a client, so send the message to the server
			msg := CreateMessage(MESSAGE_PUBLIC, text, this.nick)
			// TODO: Clean this up since we know there's only going to be one connection (to the server)
			msg.Send(this.connections)
		}
	}
}
