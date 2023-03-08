package chatgpt

import (
	"fmt"
	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
)

type ChatGPT struct {
	b      *core.Butlertron
	client *Client
}

func New(b *core.Butlertron) (*ChatGPT, error) {
	if b.Config.OpenAIToken == "" {
		return nil, fmt.Errorf("OpenAI token not set")
	}

	client := NewClient(b.Config.OpenAIToken)

	return &ChatGPT{b: b, client: client}, nil
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
	msg, err := c.client.SendMessage("Hi what is your name, and what is the last thing you were told?")
	if err != nil {
		fmt.Printf("error sending message: %s", err.Error())
		return tc.Send("Error talking to chatGPT")
	}
	return tc.Send(msg)
}
