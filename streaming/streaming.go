package streaming

import (
	"fmt"
	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
	"strings"
	"time"
)

type Streaming struct {
	b     *core.Butlertron
	token string
}

func New(b *core.Butlertron) *Streaming {
	return &Streaming{b: b, token: b.Config.StreamingToken}
}

func (s *Streaming) Name() string {
	return "Streaming"
}

func (s *Streaming) Description() string {
	return "🎬 Streaming"
}

func (s *Streaming) Command() string {
	return "/streaming"
}

func (s *Streaming) Execute(c telebot.Context) error {
	s.b.SetOnText(s.OnTextHandler, 2*time.Minute, true)
	return c.Send("What movie or TV show are you looking for?")
}

func (s *Streaming) OnTextHandler(c telebot.Context) error {
	title := c.Text()
	results, err := Search(title, s.token)
	if err != nil {
		fmt.Println(err)
		return c.Send("sorry, something went wrong searching streaming services")
	}

	inline := generateInlineHandlers(results)
	selector := s.b.RegisterInlineKeyboard(inline)

	return c.Send("Found these results:", selector)
}

func generateInlineHandlers(r []SearchResult) [][]core.InlineCommand {
	cmds := make([][]core.InlineCommand, 0)

	if len(r) > 15 {
		r = r[:15]
	}

	for i, res := range r {
		description := fmt.Sprintf("📺 %s (%d)", res.Title, res.FirstAirYear)
		if res.Type == "movie" {
			description = fmt.Sprintf("🎥 %s (%d)", res.Title, res.Year)
		}

		cmds = append(cmds, []core.InlineCommand{
			{
				Name:        fmt.Sprintf("streaming-%d", i),
				Description: description,
				Handler:     makeInlineHandler(res),
			},
		})
	}
	return cmds
}

func makeInlineHandler(r SearchResult) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if r.StreamingInfo.Ca == nil {
			msg := fmt.Sprintf("_%s_ is not available for streaming or rent.", r.Title)
			return c.Send(msg)
		}
		msg := fmt.Sprintf("Details for _%s_:\n", r.Title)

		for service, details := range r.StreamingInfo.Ca {
			serviceStr := fmt.Sprintf("\n*%s\n*", strings.Title(service))
			for _, detail := range details {
				m := ""
				if detail.Type == "subscription" {
					m = "✅ stream\n"
				} else if detail.Type == "addon" {
					m = fmt.Sprintf("✅ %s (%s)\n", detail.Type, detail.AddOn)
				} else {
					m = fmt.Sprintf("✅ %s ($%s)\n", detail.Type, detail.Price.Amount)
				}
				if !strings.Contains(serviceStr, m) {
					serviceStr += m
				}
			}
			msg += serviceStr
		}
		return c.Send(msg)
	}
}
