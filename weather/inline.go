package weather

import (
	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
)

var inlineHandlers = [][]core.InlineCommand{
	{
		core.InlineCommand{
			Name:        "foo",
			Description: "🔧 foo",
			Handler:     FooHandler,
		},
		core.InlineCommand{
			Name:        "bar",
			Description: "🔨 bar",
			Handler:     BarHandler,
		},
	},
	{
		core.InlineCommand{
			Name:        "baz",
			Description: "🪛 baz",
			Handler:     BazHandler,
		},
	},
}

// Inline Keyboard Handlers
func FooHandler(c telebot.Context) error {
	return c.Send("Fooey!")
}

func BarHandler(c telebot.Context) error {
	return c.Send("Babar!")
}

func BazHandler(c telebot.Context) error {
	return c.Send("Bazz Hands!")
}
