package chatgpt

import (
	"fmt"
	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
	"strings"
	"time"
)

type ChatGPT struct {
	b       *core.Butlertron
	client  *Client
	enabled bool
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

func (c *ChatGPT) Execute(tc telebot.Context) error {
	c.enabled = !c.enabled
	if c.enabled {
		c.b.SetOnText(c.OnTextHandler, 1*time.Hour, false)
		return tc.Send("Entering chat mode.")
	} else {
		c.b.CancelOnText()
		return tc.Send("Exiting chat mode.")
	}
}

func (c *ChatGPT) OnTextHandler(tc telebot.Context) error {
	prompt := tc.Text()

	if strings.ToLower(prompt) == "reset" {
		c.client.ResetHistory()
		return tc.Send("Reset conversation history.")
	}

	if strings.ToLower(prompt) == "switch" {
		c.client.SwapPrompt()
		if c.client.systemPrompt == standardPropmt {
			return tc.Send("Swapped to the standard prompt.")
		} else {
			return tc.Send("Swapped to the system prompt.")
		}
	}

	msg, err := c.client.SendMessage(prompt)
	if err != nil {
		fmt.Printf("error sending message: %s", err.Error())
		return tc.Send("Error talking to chatGPT")
	}

	return tc.Send(msg)
}
