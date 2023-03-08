package chatgpt

import (
	"fmt"
	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
)

func New(b *core.Butlertron) (*ChatGPT, error) {
	if b.Config.OpenAIToken == "" {
		return nil, fmt.Errorf("OpenAI token not set")
	}
	return &ChatGPT{b: b}, nil
}

type ChatGPT struct {
	b *core.Butlertron
}

func (c ChatGPT) Name() string {
	return "ChatGPT"
}

func (c ChatGPT) Description() string {
	return "🤖 ChatGPT"
}

func (c ChatGPT) Command() string {
	return "/chatgpt"
}

func (c ChatGPT) Execute(tc telebot.Context) error {
	return tc.Send("Not implemented yet")
}
