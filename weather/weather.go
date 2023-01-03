package weather

import (
	"fmt"
	"strings"

	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
)

type weatherData struct {
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Timezone       string  `json:"timezone"`
	TimezoneOffset int     `json:"timezone_offset"`
	Current        struct {
		Dt         int     `json:"dt"`
		Sunrise    int     `json:"sunrise"`
		Sunset     int     `json:"sunset"`
		Temp       float64 `json:"temp"`
		FeelsLike  float64 `json:"feels_like"`
		Pressure   int     `json:"pressure"`
		Humidity   int     `json:"humidity"`
		DewPoint   float64 `json:"dew_point"`
		Uvi        float64 `json:"uvi"`
		Clouds     int     `json:"clouds"`
		Visibility int     `json:"visibility"`
		WindSpeed  float64 `json:"wind_speed"`
		WindDeg    int     `json:"wind_deg"`
		Weather    []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"current"`
	Minutely []struct {
		Dt            int     `json:"dt"`
		Precipitation float64 `json:"precipitation"`
	} `json:"minutely"`
	Hourly []struct {
		Dt         int     `json:"dt"`
		Temp       float64 `json:"temp"`
		FeelsLike  float64 `json:"feels_like"`
		Pressure   int     `json:"pressure"`
		Humidity   int     `json:"humidity"`
		DewPoint   float64 `json:"dew_point"`
		Uvi        float64 `json:"uvi"`
		Clouds     int     `json:"clouds"`
		Visibility int     `json:"visibility"`
		WindSpeed  float64 `json:"wind_speed"`
		WindDeg    int     `json:"wind_deg"`
		WindGust   float64 `json:"wind_gust"`
		Weather    []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Pop  float64 `json:"pop"`
		Rain struct {
			OneH float64 `json:"1h"`
		} `json:"rain,omitempty"`
	} `json:"hourly"`
	Daily []struct {
		Dt        int     `json:"dt"`
		Sunrise   int     `json:"sunrise"`
		Sunset    int     `json:"sunset"`
		Moonrise  int     `json:"moonrise"`
		Moonset   int     `json:"moonset"`
		MoonPhase float64 `json:"moon_phase"`
		Temp      struct {
			Day   float64 `json:"day"`
			Min   float64 `json:"min"`
			Max   float64 `json:"max"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"temp"`
		FeelsLike struct {
			Day   float64 `json:"day"`
			Night float64 `json:"night"`
			Eve   float64 `json:"eve"`
			Morn  float64 `json:"morn"`
		} `json:"feels_like"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		DewPoint  float64 `json:"dew_point"`
		WindSpeed float64 `json:"wind_speed"`
		WindDeg   int     `json:"wind_deg"`
		WindGust  float64 `json:"wind_gust"`
		Weather   []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds int     `json:"clouds"`
		Pop    float64 `json:"pop"`
		Rain   float64 `json:"rain,omitempty"`
		Uvi    float64 `json:"uvi"`
	} `json:"daily"`
	Alerts []struct {
		SenderName  string   `json:"sender_name"`
		Event       string   `json:"event"`
		Start       int      `json:"start"`
		End         int      `json:"end"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	} `json:"alerts"`
}

// The Weather Command
type Weather struct {
	token         string
	butlertron    *core.Butlertron
	selector      *telebot.ReplyMarkup
	cachedWeather weatherData
}

func New(token string, butlertron *core.Butlertron) *Weather {
	w := &Weather{
		token:      token,
		butlertron: butlertron,
	}

	for i, row := range inlineHandlers {
		for j, handler := range row {
			h := handler.Handler
			inlineHandlers[i][j].Handler = func(ctx telebot.Context) error {
				ctx.Set("token", token)
				return h(ctx)
			}
		}
	}

	w.selector = w.butlertron.RegisterInlineKeyboard(inlineHandlers)

	return w
}

func (w Weather) Name() string {
	return "Weather"
}

func (w Weather) Description() string {
	return "⛅ Weather"
}

func (w Weather) Command() string {
	return "/weather"
}

func (w *Weather) Execute(c telebot.Context) error {
	// get and format the weather
	loc := w.butlertron.Location
	if loc == nil {
		return c.Send("Please share your location to get the weather.")
	}

	data, err := getWeather(loc, w.token)
	if err != nil {
		fmt.Printf("Error getting weather: %s", err.Error())
		return c.Send("Sorry, something went wrong fetching your weather for you!")
	}
	msg := currentConditionString(data)
	w.cachedWeather = data

	return c.Send(msg, w.selector)
}

func currentConditionString(d weatherData) string {
	currentWeather := d.Current.Weather[0]
	weatherEmoji := conditionIDToEmoji(currentWeather.ID)

	s := fmt.Sprintf("%s %s\n", weatherEmoji, currentWeather.Description)
	s += fmt.Sprintf("🌡️ feels like %.1f°C (actual %.1f°C)\n", d.Current.FeelsLike, d.Current.Temp)
	s += fmt.Sprintf("🌬️ %.2f km/h", d.Current.WindSpeed)

	if d.Alerts != nil {
		for _, alert := range d.Alerts {
			s += "\n\n"
			s += fmt.Sprintf("🚨 %s Alert 🚨\n%s\n", strings.Title(alert.Event), alert.Description)
		}
	}
	return s
}
