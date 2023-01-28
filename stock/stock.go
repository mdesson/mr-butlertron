package stock

import (
	"fmt"
	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
)

type Stock struct {
	b       *core.Butlertron
	symbols []string
}

func New(b *core.Butlertron) *Stock {
	return &Stock{b: b}
}

func (s Stock) Name() string {
	return "Stock"
}

func (s Stock) Description() string {
	return "📈 Stock"
}

func (s Stock) Command() string {
	return "/stock"
}

func (s *Stock) Execute(c telebot.Context) error {

	handlers, err := createInlineHandlers(s)
	if err != nil {
		fmt.Println(err)
		return c.Send("Something went wrong getting your socks")
	}

	msg := "Here's your watchlist:"
	if len(s.symbols) == 0 {
		msg = "Hit the add button to add a stock to your watchlist."
	}

	selector := s.b.RegisterInlineKeyboard(handlers)
	return c.Send(msg, selector)
}
