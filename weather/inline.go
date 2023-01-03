package weather

import (
	"fmt"
	"strings"
	"time"

	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
)

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
			Name:        "alert-weather",
			Description: "🚨 Weather Alert 🚨",
			Handler:     AlertHandler,
			Conditional: true,
		},
	},
}

// Inline Keyboard Handlers
func WeeklyHandler(c telebot.Context) error {
	data, err := getWeatherForInline(c)
	if err != nil {
		fmt.Println(err)
		return c.Send("sorry, something went wrong")
	}

	// format and return response
	msg := "The weekly forecast is:\n\n"
	for _, day := range data.Daily {
		weekday := time.Unix(int64(day.Dt), 0).Weekday().String()
		weather := day.Weather[0]
		emoji := conditionIDToEmoji(weather.ID)

		msg += fmt.Sprintf("%s *%s*: %s\n", emoji, weekday, weather.Description)
		msg += fmt.Sprintf("☀ %.1f°C (feels like %.1f°C)\n", day.Temp.Day, day.FeelsLike.Day)
		msg += fmt.Sprintf("🌙 %.1f°C (feels like %.1f°C)\n", day.Temp.Eve, day.FeelsLike.Night)
		msg += "\n"
	}
	return c.Send(msg)
}

func HourlyHandler(c telebot.Context) error {
	data, err := getWeatherForInline(c)
	if err != nil {
		fmt.Println(err)
		return c.Send("sorry, something went wrong")
	}

	// format and return response
	msg := "Your next eight hours are:\n\n"
	for _, hour := range data.Hourly[:8] {
		localTime := time.Unix(int64(hour.Dt), 0).Local().Format("3:04pm")
		weather := hour.Weather[0]
		emoji := conditionIDToEmoji(weather.ID)

		msg += fmt.Sprintf("*%s*\n", localTime)
		msg += fmt.Sprintf("%s %s\n", emoji, weather.Description)
		msg += fmt.Sprintf("🌡 %.1f°C (feels like %.1f°C)\n", hour.Temp, hour.FeelsLike)
		msg += fmt.Sprintf("💧 %.1fmm\n", hour.Pop)
		msg += fmt.Sprintf("🌬️ %.2f km/h\n", hour.WindSpeed)
		msg += "\n"
	}
	return c.Send(msg)
}

func AlertHandler(c telebot.Context) error {
	data, err := getWeatherForInline(c)
	if err != nil {
		fmt.Println(err)
		return c.Send("sorry, something went wrong")
	}

	msg := ""
	for _, alert := range data.Alerts {
		msg += "\n\n"
		msg += fmt.Sprintf("🚨 %s Alert 🚨\n%s\n", strings.Title(alert.Event), alert.Description)
	}

	return c.Send(msg)
}

func getWeatherForInline(c telebot.Context) (weatherData, error) {
	// Get variables from context
	loc, ok := c.Get("location").(*telebot.Location)
	if !ok {
		return weatherData{}, fmt.Errorf("Couldn't cast location middleware result to location type")
	}

	if loc == nil {
		return weatherData{}, fmt.Errorf("Please share your location to get the weather.")
	}

	token, ok := c.Get("token").(string)
	if !ok {
		return weatherData{}, fmt.Errorf("Couldn't cast token middleware result to string")
	}

	if token == "" {
		return weatherData{}, fmt.Errorf("weather token wasn't set")
	}

	// get latest weather data
	data, err := getWeather(loc, token)
	if err != nil {
		return weatherData{}, err
	}
	return data, nil
}
