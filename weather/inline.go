package weather

import (
	"fmt"
	"time"

	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
)

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
			Description: "⌚ next 8h",
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
	// Get variables from context
	loc, ok := c.Get("location").(*telebot.Location)
	if !ok {
		fmt.Println("Couldn't cast location middleware result to location type")
		return c.Send("Sorry, something went wrong")
	}

	if loc == nil {
		return c.Send("Please share your location to get the weather.")
	}

	token, ok := c.Get("token").(string)
	if !ok {
		fmt.Println("Couldn't cast token middleware result to string")
		return c.Send("Sorry, something went wrong")
	}

	if token == "" {
		fmt.Println("weather token wasn't set")
		return c.Send("Sorry, something went wrong")
	}

	// get latest weather data
	data, err := getWeather(loc, token)
	if err != nil {
		fmt.Println(err)
		return c.Send("Sorry, something went wrong")
	}

	// format and return response
	msg := "The weekly forecast is:\n\n"
	for _, day := range data.Daily {
		weekday := time.Unix(int64(day.Dt), 0).Weekday().String()
		weather := day.Weather[0]
		emoji := conditionIDToEmoji(weather.ID)

		msg += fmt.Sprintf("%s %s: %s\n", emoji, weekday, weather.Description)
		msg += fmt.Sprintf("☀ %.1f°C (feels like %.1f°C)\n", day.Temp.Day, day.FeelsLike.Day)
		msg += fmt.Sprintf("🌙 %.1f°C (feels like %.1f°C)\n", day.Temp.Eve, day.FeelsLike.Night)
		msg += "\n"
	}
	return c.Send(msg)
}

func HourlyHandler(c telebot.Context) error {
	// Get variables from context
	loc, ok := c.Get("location").(*telebot.Location)
	if !ok {
		fmt.Println("Couldn't cast location middleware result to location type")
		return c.Send("Sorry, something went wrong")
	}

	if loc == nil {
		return c.Send("Please share your location to get the weather.")
	}

	token, ok := c.Get("token").(string)
	if !ok {
		fmt.Println("Couldn't cast token middleware result to string")
		return c.Send("Sorry, something went wrong")
	}

	if token == "" {
		fmt.Println("weather token wasn't set")
		return c.Send("Sorry, something went wrong")
	}

	// get latest weather data
	data, err := getWeather(loc, token)
	if err != nil {
		fmt.Println(err)
		return c.Send("Sorry, something went wrong")
	}

	// format and return response
	msg := "Your next eight hours are:\n\n"
	for _, hour := range data.Hourly[:8] {
		localTime := time.Unix(int64(hour.Dt), 0).Local().Format("3:04pm")
		weather := hour.Weather[0]
		emoji := conditionIDToEmoji(weather.ID)

		msg += localTime + "\n"
		msg += fmt.Sprintf("%s %s\n", emoji, weather.Description)
		msg += fmt.Sprintf("🌡 %.1f°C (feels like %.1f°C)\n", hour.Temp, hour.FeelsLike)
		msg += fmt.Sprintf("🌬️ %.2f km/h\n", hour.WindSpeed)
		msg += "\n"
	}
	return c.Send(msg)
}

func PrecipitationHandler(c telebot.Context) error {
	return c.Send("Bazz Hands!")
}
