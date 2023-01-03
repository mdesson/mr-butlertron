package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/telebot.v3"
)

func getWeather(loc *telebot.Location, token string) (weatherData, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?units=metric&lat=%.2f&lon=%.2f&appid=%s", loc.Lat, loc.Lng, token)
	res, err := http.Get(url)
	if err != nil {
		return weatherData{}, err
	}

	if res.StatusCode >= 300 {
		msg := fmt.Sprintf("Weather request failed with code %d", res.StatusCode)
		return weatherData{}, errors.New(msg)
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
