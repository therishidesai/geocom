package main

import (
	"fmt"
	"time"

	"github.com/marcusolsson/tui-go"
)

type UI struct {
	history *tui.Box
	input   *tui.Entry
}

func CreateUI() *UI {
	history := tui.NewVBox()
	history.SetBorder(true)
	history.Append(tui.NewSpacer())

	history.Append(tui.NewHBox(
		tui.NewLabel(time.Now().Format("15:04")),
		tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", "root"))),
		tui.NewLabel("Welcome to geocom, press ESC to quit."),
		tui.NewSpacer(),
	))

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	
	return &UI{
		history: history,
		input: input} 
}

func (this *UI)  updateMessage(author string, text string) {
	this.history.Append(tui.NewHBox(
		tui.NewLabel(time.Now().Format("15:04")),
		tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", author))),
		tui.NewLabel(text),
		tui.NewSpacer(),
	))
}
