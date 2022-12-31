package core

import "gopkg.in/telebot.v3"

type Butlertron struct {
	Bot      *telebot.Bot
	Location *telebot.Location
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
