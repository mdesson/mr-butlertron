package dalle2

import (
	"fmt"
	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
	"time"
)

type DALLE2 struct {
	b      *core.Butlertron
	client *Client
}

func New(b *core.Butlertron) *DALLE2 {
	return &DALLE2{b: b, client: NewClient(b.Config.OpenAIToken)}
}

func (d DALLE2) Name() string {
	return "ChatGPT"
}

func (d DALLE2) Description() string {
	return "🖼️ DALL-E 2"
}

func (d DALLE2) Command() string {
	return "/dalle"
}

func (d *DALLE2) Execute(c telebot.Context) error {
	d.b.SetOnText(d.OnTextHandler, 5*time.Minute, true)
	return c.Send("Enter your prompt.")
}

func (d *DALLE2) OnTextHandler(c telebot.Context) error {
	prompt := c.Text()
	URLs, err := d.client.RequestImages(prompt)
	if err != nil {
		fmt.Printf("error sending message: %s", err.Error())
		return c.Send("Error getting images")
	}
	album := URLsToAlbum(URLs)

	selector := d.b.RegisterInlineKeyboard([][]core.InlineCommand{{{
		Name:        "retry-images",
		Description: "🔁",
		Handler:     RegenerateImages(prompt, d),
	},
	}})

	c.SendAlbum(album)
	return c.Send("Retry?", selector)
}

func RegenerateImages(prompt string, dalle *DALLE2) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		URLs, err := dalle.client.RequestImages(prompt)
		if err != nil {
			fmt.Printf("error sending message: %s", err.Error())
			return c.Send("Error getting images")
		}
		album := URLsToAlbum(URLs)

		selector := dalle.b.RegisterInlineKeyboard([][]core.InlineCommand{{{
			Name:        fmt.Sprintf("retry-images-%d", time.Now().Unix()),
			Description: "🔁",
			Handler:     RegenerateImages(prompt, dalle),
		},
		}})

		c.SendAlbum(album)
		return c.Send("Retry?", selector)
	}

}

func URLsToAlbum(URLs []string) telebot.Album {
	album := telebot.Album{}
	for _, URL := range URLs {
		album = append(album, &telebot.Photo{File: telebot.FromURL(URL)})
	}
	return album
}
