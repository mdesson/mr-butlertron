package main

import (
	"log"
	"os"

	"github.com/mdesson/mr-butlertron/core"
	"github.com/mdesson/mr-butlertron/etymology"
	"github.com/mdesson/mr-butlertron/stock"
	"github.com/mdesson/mr-butlertron/weather"
	telebot "gopkg.in/telebot.v3"
)

var (
	b        *core.Butlertron
	commands []core.Command
	menu     *telebot.ReplyMarkup
)

func init() {
	// init bot
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	var err error
	b, err = core.NewButlertron(telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	// init menu
	menu = &telebot.ReplyMarkup{ResizeKeyboard: true}

	//// Custom Commands ////
	// init weather command
	weatherToken := os.Getenv("WEATHER_TOKEN")
	weatherCmd := weather.New(weatherToken, b)
	commands = append(commands, weatherCmd)

	// init etymology command
	etymologyCmd := etymology.New(b)
	commands = append(commands, etymologyCmd)

	// init stock command
	stockCmd := stock.New(b)
	commands = append(commands, stockCmd)

}

func main() {
	bot := b.Bot
	setCommandsArgs := make([]telebot.Command, 0)

	// register commands
	for _, command := range commands {
		// Add text handler
		bot.Handle(command.Command(), command.Execute)

		// Add as menu button
		btn := menu.Text(command.Name())
		menu.Reply(menu.Row(btn))

		// ensure command is set
		setCommandsArgs = append(setCommandsArgs, telebot.Command{Text: command.Command(), Description: command.Description()})
	}

	// set commands
	if err := bot.SetCommands(setCommandsArgs); err != nil {
		log.Fatal(err)
	}

	bot.Start()
}
