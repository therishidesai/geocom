package main

import (
	"fmt"
	"os"
	"time"
	
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
		//fmt.Printf("[*] Connecting to %s\n", ip)
		err := client.connectToServer(ip)
		if err != nil {
			//fmt.Printf("[*] Failed to connect to %s\n", ip)
			os.Exit(1)
		}
		//fmt.Printf("[*] Connected to %s\n", ip)
	}

	ui.input.OnSubmit(func(e *tui.Entry) {
		ui.history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", client.nick))),
			tui.NewLabel(e.Text()),
			tui.NewSpacer(),
		))
		client.handleInput(e.Text())
		ui.input.SetText("")
	})
	
	inputBox := tui.NewHBox(ui.input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(ui.history, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)
	
	root := tui.NewHBox(chat)

	view := tui.New(root)
	view.SetKeybinding("Esc", func() { view.Quit() })

	if err := view.Run(); err != nil {
		panic(err)
	}

}
