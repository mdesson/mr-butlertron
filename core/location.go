package core

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

type Location struct {
	butlertron *Butlertron
}

func NewLocation(b *Butlertron) *Location {
	return &Location{b}
}

func (w Location) Name() string {
	return "Location"
}

func (w Location) Description() string {
	return "🗺️ location"
}

func (w Location) Command() string {
	return telebot.OnLocation
}

func (w *Location) Execute(c telebot.Context) error {
	loc := c.Message().Location
	w.butlertron.Location = loc
	msg := fmt.Sprintf("Updated your location to %v, %v.", loc.Lat, loc.Lng)
	return c.Send(msg)
}

func (l *Location) LocationMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		c.Set("location", l.butlertron.Location)
		return next(c)
	}
}
