package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"unicode"

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
		Uvi        int     `json:"uvi"`
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
}

// The Weather Command
type Weather struct {
	token      string
	butlertron *core.Butlertron
}

func New(token string, butlertron *core.Butlertron) *Weather {
	return &Weather{
		token:      token,
		butlertron: butlertron,
	}
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

func (w Weather) RequiresLocation() bool {
	return true
}

func (w *Weather) Execute(c telebot.Context) error {
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
	return c.Send(msg)
}

// Helpers
func getWeather(loc *telebot.Location, token string) (weatherData, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?units=metric&lat=%.2f&lon=%.2f&exclude=hourly,daily&appid=%s", loc.Lat, loc.Lng, token)
	res, err := http.Get(url)
	if err != nil {
		return weatherData{}, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return weatherData{}, err
	}
	d := weatherData{}
	if err := json.Unmarshal(body, &d); err != nil {
		return d, err
	}
	return d, nil
}

func currentConditionString(d weatherData) string {
	currentWeather := d.Current.Weather[0]
	weatherEmoji := conditionIDToEmoji(currentWeather.ID)
	conditions := string(byte(unicode.ToUpper(rune(currentWeather.Description[0])))) + currentWeather.Description[1:]

	s := fmt.Sprintf("%s %s, feels like %.1f°C (actual %.1f°C).", weatherEmoji, conditions, d.Current.FeelsLike, d.Current.Temp)
	return s
}

func conditionIDToEmoji(id int) string {
	if id < 300 {
		return "⛈️"
	} else if id < 600 {
		return "🌧"
	} else if id < 700 {
		return "🌨️"
	} else if id < 800 {
		return "🌁"
	} else if id == 800 {
		return "☀"
	} else if id < 900 {
		return "🌥️"
	}

	fmt.Printf("Passed invalid weather condition id %d\n", id)
	return "?"
}
