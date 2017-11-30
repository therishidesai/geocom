package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

const PORT = 5000

var (
	connections = make(map[string]net.Conn)
	nick        string
)

func StartServer() error {
	IP := fmt.Sprintf("127.0.0.1:%d", PORT)
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr).IP
	fmt.Printf("[*] Handling connection from %s\n", remoteAddr)
	go receive(conn)
}

func receive(conn net.Conn) {
	for {
		msg, err := readFromConn(conn)
		if err != nil {
			return
		}
		switch msg.Kind {
		case MESSAGE_CONNECT:
			fmt.Printf("[*] Initializing connection to %s\n", msg.Contents)
			connections[msg.Contents] = conn
			response := CreateMessage(MESSAGE_CONNECT, "", nick)
			response.Send(connections)
		case MESSAGE_PUBLIC:
			fmt.Printf("[*] %s said: %s\n", msg.Author, msg.Contents)
		default:
			fmt.Printf("[*] Bad message.")
		}
	}
}
func readFromConn(conn net.Conn) (*Message, error) {
	dec := json.NewDecoder(conn)
	msg := new(Message)
	if err := dec.Decode(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func createConnection(ip string) (net.Conn, error) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func ConnectToServer(ip string) error {
	fmt.Printf("[*] Connecting to %s\n", ip)
	conn, err := createConnection(ip)
	if err != nil {
		fmt.Println("[*] Failed to connect to server")
		return err
	}

	msg := CreateMessage(MESSAGE_CONNECT, "", nick)
	msg.SendTo(conn)

	msg, err = readFromConn(conn)
	if err != nil {
		return err
	}
	connections[msg.Author] = conn
	go receive(conn)
	return nil
}

func HandleInput() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Send message: ")
		text, _ := reader.ReadString('\n')
		msg := CreateMessage(MESSAGE_PUBLIC, text, nick)
		msg.Send(connections)
	}
}
