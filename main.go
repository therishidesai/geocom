package main

import (
	"fmt"
	"os"
)

func main() {

	numArgs := len(os.Args)
	var ip string
	var nick string

	if numArgs == 2 {
		nick = os.Args[1]
	} else if numArgs == 3 {
		nick = os.Args[1]
		ip = fmt.Sprintf("%s:%d", os.Args[2], PORT)
	} else {
		fmt.Printf("Usage: %s [nick] [ip]\n", os.Args[0])
		os.Exit(1)
	}

	var client *Client
	if ip == "" {
		client = CreateClient(nick, true)
		go client.StartServer()
	} else {
		client = CreateClient(nick, false)
		fmt.Printf("[*] Connecting to %s\n", ip)
		err := client.connectToServer(ip)
		if err != nil {
			fmt.Printf("[*] Failed to connect to %s\n", ip)
			os.Exit(1)
		}
		fmt.Printf("[*] Connected to %s\n", ip)
	}
	client.handleInput()
}
