package weather

import (
	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
)

// TODO: (inline) Daily forecast
// TODO: (inline) Hourly forecast
// TODO: (inline) Minutely precipitation

var inlineHandlers = [][]core.InlineCommand{
	{
		core.InlineCommand{
			Name:        "weekly-weather",
			Description: "📆 weekly",
			Handler:     WeeklyHandler,
		},
		core.InlineCommand{
			Name:        "hourly-weather",
			Description: "⌚ next 48h",
			Handler:     HourlyHandler,
		},
	},
	{
		core.InlineCommand{
			Name:        "precipitation-weather",
			Description: "💧 rain next hour",
			Handler:     PrecipitationHandler,
		},
	},
}

// Inline Keyboard Handlers
func WeeklyHandler(c telebot.Context) error {
	return c.Send("Fooey!")
}

func HourlyHandler(c telebot.Context) error {
	return c.Send("Babar!")
}

func PrecipitationHandler(c telebot.Context) error {
	return c.Send("Bazz Hands!")
}
