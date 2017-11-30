package main

import (
	"fmt"
	"os"
)

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

	if ip == "" {
		go StartServer()
	} else {
		err := ConnectToServer(ip)
		if err != nil {
			fmt.Printf("[*] Failed to connect to %s\n", ip)
			os.Exit(1)
		}
		fmt.Printf("[*] Connected to %s\n", ip)
	}

	HandleInput()
}
