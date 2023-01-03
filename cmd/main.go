package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/mdesson/mr-butlertron/core"
	"github.com/mdesson/mr-butlertron/weather"
	telebot "gopkg.in/telebot.v3"
)

var (
	b           *core.Butlertron
	commandsMux sync.RWMutex
	commands    []core.Command
	menu        *telebot.ReplyMarkup
)

func init() {
	// init bot
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	pref := telebot.Settings{
		Token:  telegramToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	// init Mr. Butlertron
	b = &core.Butlertron{Bot: bot}

	// init menu
	menu = &telebot.ReplyMarkup{ResizeKeyboard: true}

	//// Core Commands ////
	// init location command
	locationCmd := core.NewLocation(b)
	bot.Handle(locationCmd.Command(), locationCmd.Execute)
	bot.Use(locationCmd.LocationMiddleware)

	//// Custom Commands ////
	// init weather command
	weatherToken := os.Getenv("WEATHER_TOKEN")
	weatherCmd := weather.New(weatherToken, b)
	commands = append(commands, weatherCmd)
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
	bot.SetCommands(setCommandsArgs)

	bot.Start()
}
