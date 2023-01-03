package core

import (
	"gopkg.in/telebot.v3"
)

type Butlertron struct {
	Bot      *telebot.Bot
	Location *telebot.Location
}

func (b *Butlertron) RegisterInlineKeyboard(commands [][]InlineCommand) *telebot.ReplyMarkup {
	selector := &telebot.ReplyMarkup{}
	keyboardRows := make([]telebot.Row, 0)

	for _, row := range commands {
		keyboardRow := make([]telebot.Btn, 0)
		for _, command := range row {
			btn := selector.Data(command.Description, command.Name)
			keyboardRow = append(keyboardRow, btn)
			b.Bot.Handle(&btn, command.Handler)
		}
		keyboardRows = append(keyboardRows, keyboardRow)
	}

	selector.Inline(keyboardRows...)

	return selector
}

// Command is the interface that all bot commands must implement.
type Command interface {
	// Name returns the name of the command.
	Name() string
	// Description returns a brief description of the command.
	Description() string
	// Command is the string which will call the command
	Command() string
	// Execute is called when the command is triggered.
	Execute(c telebot.Context) error
}

type InlineCommand struct {
	Name        string
	Description string
	Selector    *telebot.ReplyMarkup
	Handler     telebot.HandlerFunc
}
