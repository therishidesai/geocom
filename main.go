package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const PORT = 5000

var connections = make(map[string]net.Conn)
var nick string

func main() {

	numArgs := len(os.Args)
	var ip string

	if numArgs == 2 {
		nick = os.Args[1]
	} else if numArgs == 3 {
		nick = os.Args[1]
		ip = fmt.Sprintf("%s:%d", os.Args[2], PORT)
	} else {
		fmt.Printf("Usage: %s [nick] [ip]\n", os.Args[0])
		os.Exit(1)
	}

	fmt.Printf("Nick: %s\n", nick)
	if ip == "" {
		go StartServer()
	} else {
		err := connectToServer(ip)
		if err != nil {
			fmt.Printf("[*] Failed to connect to %s\n", ip)
			os.Exit(1)
		}
		fmt.Printf("[*] Connected to %s\n", ip)
	}

	HandleInput()
}

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

func connectToServer(ip string) error {
	fmt.Printf("[*] Connecting to %s\n", ip)
	conn, err := createConnection(ip)
	if err != nil {
		fmt.Println("[*] Failed to connect to server")
		return err
	}
	conn.Write([]byte(fmt.Sprintln(nick)))
	msg, err := readFromConn(conn)
	if err != nil {
		fmt.Println("[*] Could not read nick")
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
			fmt.Printf("[*] Writing to %s", name)
			_, err := conn.Write([]byte(fmt.Sprintf("MSG:%s", text)))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
