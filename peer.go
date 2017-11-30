package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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
		if strings.HasPrefix(msg, "MSG:") {
			msg = msg[4:]
			fmt.Printf("[*] Received: %s\n", msg)
		} else {
			fmt.Printf("[*] Initializing connection to %s\n", msg)
			connections[msg] = conn
			conn.Write([]byte(fmt.Sprintln(nick)))
		}
	}
}
func readFromConn(conn net.Conn) (string, error) {
	return bufio.NewReader(conn).ReadString('\n')
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
	conn.Write([]byte(fmt.Sprintln(nick)))
	msg, err := readFromConn(conn)
	if err != nil {
		return err
	}
	connections[msg] = conn
	go receive(conn)
	return nil
}

func HandleInput() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Send message: ")
		text, _ := reader.ReadString('\n')
		for name, conn := range connections {
			_, err := conn.Write([]byte(fmt.Sprintf("MSG:%s", text)))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
