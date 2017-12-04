package main

import (
	"fmt"
	"os"

	"github.com/marcusolsson/tui-go"
)

func main() {

	ui := CreateUI()
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
		client = CreateClient(nick, true, ui)
		go client.StartServer()
	} else {
		client = CreateClient(nick, false, ui)
		err := client.connectToServer(ip)
		if err != nil {
			os.Exit(1)
		}
	}

	ui.input.OnSubmit(func(e *tui.Entry) {
		if e.Text() == "" {
			return
		}
		ui.updateMessage(nick, e.Text())
		client.handleInput(e.Text())
		ui.input.SetText("")
	})

	if err := ui.view.Run(); err != nil {
		panic(err)
	}

}
