package chatgpt

import (
	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
)

func InlineHandlers(c *ChatGPT) [][]core.InlineCommand {
	return [][]core.InlineCommand{
		{
			{
				Name:        "swap-chatgpt-prompt",
				Description: "🔄 swap prompts",
				Handler:     SwapSystemPromptFunc(c),
			},
			{
				Name:        "reset-chatgpt-history",
				Description: "📔 reset history",
				Handler:     ResetChatFunc(c),
			},
		},
		{
			{
				Name:        "toggle-chatgpt",
				Description: "🔌 turn on/off",
				Handler:     ToggleFunc(c),
			},
		},
	}
}

func ToggleFunc(c *ChatGPT) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		c.enabled = !c.enabled
		if c.enabled {
			return ctx.Send("Enabled ChatGPT")
		} else {
			return ctx.Send("Disabled ChatGPT")
		}
	}
}

func SwapSystemPromptFunc(c *ChatGPT) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		c.client.SwapPrompt()
		if c.client.systemPrompt == standardPropmt {
			return ctx.Send("Swapped to standard prompt")
		} else {
			return ctx.Send("Swapped to DAN")
		}
	}
}

func ResetChatFunc(c *ChatGPT) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		c.client.ResetHistory()
		return ctx.Send("Reset chat history")
	}
}
