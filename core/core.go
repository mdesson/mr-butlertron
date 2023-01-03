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
			b.Bot.Handle(&btn, command.Handler)
			keyboardRow = append(keyboardRow, btn)
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
	// Name of the Inline Command
	Name string
	// Description of the Inline Command, this will be displayed on the button
	Description string
	// The Selector that stores this command
	Selector *telebot.ReplyMarkup
	// The Handler which is executed
	Handler telebot.HandlerFunc
	// The Parent Command. The inline command buttons will be displayed below this
	Parent Command
	// Conditional refers to if it should be loaded by default as a button. If true, it will not be added on init
	Conditional bool
}
