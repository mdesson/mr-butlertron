package etymology

import (
	"fmt"
	"time"

	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
)

type Etymology struct {
	b *core.Butlertron
}

func (e Etymology) Name() string {
	return "Etymology"
}

func (e Etymology) Description() string {
	return "📖 Etymology"
}

func (e Etymology) Command() string {
	return "/etymology"
}

func (e *Etymology) Execute(c telebot.Context) error {
	e.b.SetOnText(onTextHandler, 1*time.Minute)
	return c.Send("What word are you looking for?")
}

func New(b *core.Butlertron) *Etymology {
	return &Etymology{b}
}

func onTextHandler(c telebot.Context) error {
	txt := fmt.Sprintf("You sent %s", c.Message().Text)
	return c.Send(txt)
}
