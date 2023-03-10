package main

import (
	"github.com/mdesson/mr-butlertron/chatgpt"
	"github.com/mdesson/mr-butlertron/core"
	"github.com/mdesson/mr-butlertron/etymology"
	"github.com/mdesson/mr-butlertron/stock"
	"github.com/mdesson/mr-butlertron/weather"
	"gopkg.in/telebot.v3"

	"fmt"
	"log"
)

var (
	b *core.Butlertron
)

func init() {
	fmt.Println("Initializing Butlertron.")
	var err error
	commands := make([]core.Command, 0)
	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}

	//// init bot ////
	b, err = core.NewButlertron()
	if err != nil {
		log.Fatal(err)
	}

	//// Custom Commands ////
	// init weather command
	weatherCmd, err := weather.New(b)
	if err != nil {
		log.Fatal(err)
	}
	commands = append(commands, weatherCmd)

	// init etymology command
	etymologyCmd := etymology.New(b)
	commands = append(commands, etymologyCmd)

	// init stock command
	stockCmd := stock.New(b)
	commands = append(commands, stockCmd)

	// init chatgpt command
	chatgptCmd, err := chatgpt.New(b)
	if err != nil {
		log.Fatal(err)
	}
	commands = append(commands, chatgptCmd)

	bot := b.Bot
	setCommandsArgs := make([]telebot.Command, 0)

	//// register commands ////
	for _, command := range commands {
		// Add text handler
		bot.Handle(command.Command(), command.Execute)

		// Add as menu button
		btn := menu.Text(command.Name())
		menu.Reply(menu.Row(btn))

		// ensure command is set
		setCommandsArgs = append(setCommandsArgs, telebot.Command{Text: command.Command(), Description: command.Description()})
	}

	//// set commands ////
	if err := bot.SetCommands(setCommandsArgs); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Starting Butlertron.")
	b.Bot.Start()
}
