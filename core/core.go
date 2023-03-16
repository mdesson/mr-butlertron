package core

import (
	"context"
	"sync"
	"time"

	"github.com/caarlos0/env"
	"gopkg.in/telebot.v3"
)

type Config struct {
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	WeatherToken     string `env:"WEATHER_TOKEN"`
	OpenAIToken      string `env:"OPENAI_TOKEN"`
}

type Butlertron struct {
	Bot            *telebot.Bot
	Config         Config
	Location       *telebot.Location
	onTextMetadata *onTextMetadata
}

func NewButlertron() (*Butlertron, error) {
	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, err
	}

	pref := telebot.Settings{
		Token:     c.TelegramBotToken,
		Poller:    &telebot.LongPoller{Timeout: 10 * time.Second},
		ParseMode: telebot.ModeMarkdown,
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	// init Mr. Butlertron
	b := &Butlertron{Bot: bot, Config: c}

	// init onText settings
	b.onTextMetadata = &onTextMetadata{mu: &sync.Mutex{}}

	//// Core Commands ////
	// init location command
	locationCmd := NewLocation(b)
	bot.Handle(locationCmd.Command(), locationCmd.Execute)
	bot.Use(locationCmd.LocationMiddleware)

	return b, nil
}

func (b *Butlertron) RegisterInlineKeyboard(commands [][]InlineCommand) *telebot.ReplyMarkup {
	selector := &telebot.ReplyMarkup{}
	keyboardRows := make([]telebot.Row, 0)

	for _, row := range commands {
		keyboardRow := make([]telebot.Btn, 0)
		for _, command := range row {
			btn := selector.Data(command.Description, command.Name)
			b.Bot.Handle(&btn, command.Handler)
			keyboardRow = append(keyboardRow, btn)
		}
		keyboardRows = append(keyboardRows, keyboardRow)
	}

	selector.Inline(keyboardRows...)

	return selector
}

// SetOnTextDefault sets the default responder if no OnText has been manually trigger
func (b *Butlertron) SetOnTextDefault(h telebot.HandlerFunc) {
	b.onTextMetadata.defaultHandler = h
}

// SetOnText will set the command that will run when text is next sent
// If the deadline is exceed or it has been cancelled, there will be no reply
func (b *Butlertron) SetOnText(h telebot.HandlerFunc, timeout time.Duration, cancelAfterHandling bool) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	b.onTextMetadata.ctx = &ctx
	b.onTextMetadata.cancel = &cancel

	b.Bot.Handle(telebot.OnText, func(c telebot.Context) error {
		b.onTextMetadata.mu.Lock()
		defer b.onTextMetadata.mu.Unlock()
		if ctx == nil || ctx.Err() != nil {
			return b.onTextMetadata.defaultHandler(c)
		}
		if !cancelAfterHandling {
			cancel()
		}
		return h(c)
	})
}

// CancelOnText cancels the current command assigned to OnText
// If there is no OnText handler, it has been cancelled, or it has timed out, it is a noop
func (b *Butlertron) CancelOnText() {
	ctx := b.onTextMetadata.ctx
	mu := b.onTextMetadata.mu
	cancel := *b.onTextMetadata.cancel
	if ctx != nil {
		mu.Lock()
		defer mu.Unlock()
		cancel()
	}
}

// Command is the interface that all bot commands must implement.
type Command interface { // Name returns the name of the command.
	Name() string
	// Description returns a brief description of the command.
	Description() string
	// Command is the string which will call the command
	Command() string
	// Execute is called when the command is triggered.
	Execute(c telebot.Context) error
}

type InlineCommand struct {
	// Name of the Inline Command
	Name string
	// Description of the Inline Command, this will be displayed on the button
	Description string
	// The Selector that stores this command
	Selector *telebot.ReplyMarkup
	// The Handler which is executed
	Handler telebot.HandlerFunc
	// The Parent Command. The inline command buttons will be displayed below this
	Parent Command
	// Conditional refers to if it should be loaded by default as a button. If true, it will not be added on init
	Conditional bool
}

type onTextMetadata struct {
	ctx            *context.Context
	cancel         *context.CancelFunc
	mu             *sync.Mutex
	defaultHandler telebot.HandlerFunc
}
