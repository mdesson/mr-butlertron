package main

import (
	"github.com/mdesson/mr-butlertron/chatgpt"
	"github.com/mdesson/mr-butlertron/core"
	"github.com/mdesson/mr-butlertron/dalle2"
	"github.com/mdesson/mr-butlertron/etymology"
	"github.com/mdesson/mr-butlertron/stock"
	"github.com/mdesson/mr-butlertron/streaming"
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

	// init streaming command
	streamingCmd := streaming.New(b)
	commands = append(commands, streamingCmd)

	// init stock command
	stockCmd := stock.New(b)
	commands = append(commands, stockCmd)

	// init DALLE-2 Command
	dalleCmd := dalle2.New(b)
	commands = append(commands, dalleCmd)

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

	//// set default onText handler ////
	b.SetOnTextDefault(chatgptCmd.OnTextHandler)
}

func main() {
	fmt.Println("Starting Butlertron.")
	b.Bot.Start()
}
