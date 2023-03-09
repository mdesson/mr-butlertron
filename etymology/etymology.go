package etymology

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
)

var (
	maxDefinitions                  = 6
	b              *core.Butlertron = nil
)

type Etymology struct{}

func New(butler *core.Butlertron) *Etymology {
	b = butler
	return &Etymology{}
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
	b.SetOnText(onTextHandler, 1*time.Minute, true)
	return c.Send("What word are you looking for?")
}

func onTextHandler(c telebot.Context) error {
	history, err := lookUpWord(c.Text())
	if err != nil {
		return err
	}

	handlers := createInlineHandlers(history)
	selector := b.RegisterInlineKeyboard(handlers)

	return c.Send("Here's the search results:", selector)
}

func lookUpWord(word string) (map[string]string, error) {
	// join word with + for http request
	word = strings.TrimSpace(word)
	word = strings.ReplaceAll(word, " ", "+")

	// make http request
	url := fmt.Sprintf("https://www.etymonline.com/search?q=%s", word)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:108.0) Gecko/20100101 Firefox/108.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-CA,en-US;q=0.7,en;q=0.3")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://www.etymonline.com/")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// read response and unzip it
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// parse html
	definitions, err := parseHTML(string(raw))
	if err != nil {
		return nil, err
	}

	return definitions, nil
}
