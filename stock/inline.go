package stock

import (
	"fmt"
	"github.com/mdesson/mr-butlertron/core"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/quote"
	"gopkg.in/telebot.v3"
	"strings"
	"sync"
	"time"
)

func createInlineHandlers(stock *Stock) ([][]core.InlineCommand, error) {
	// get descriptions in parallel
	wg := sync.WaitGroup{}
	var descErr error

	descriptions := make(chan string, len(stock.symbols))

	for _, symbol := range stock.symbols {
		wg.Add(1)
		go func(s string) {
			s, err := getDescription(s)
			if err != nil {
				descErr = err
				descriptions <- ""
			} else {
				descriptions <- s
			}
			wg.Done()
		}(symbol)
	}
	wg.Wait()
	close(descriptions)

	// each symbol must have both a quote and a corresponding ticker
	if descErr != nil {
		return nil, descErr
	}

	// add inline commands
	cmds := make([][]core.InlineCommand, 0)
	for i, symbol := range stock.symbols {
		d := <-descriptions

		if i%2 == 0 {
			cmds = append(cmds, []core.InlineCommand{})
		}

		name := fmt.Sprintf("stock-%d", i)
		handler := func(c telebot.Context) error {

			// keyboard and handlers setup
			deleteHandler := func(c telebot.Context) error {
				for j := 0; j < len(stock.symbols); j++ {
					if stock.symbols[j] == symbol {
						stock.symbols = append(stock.symbols[:j], stock.symbols[j+1:]...)
					}
				}
				return c.Send("Deleted!")
			}

			selector := stock.b.RegisterInlineKeyboard([][]core.InlineCommand{
				{
					core.InlineCommand{
						Name:        "delete",
						Description: "🗑️ Delete",
						Handler:     deleteHandler,
					},
				},
			})

			// get text to return to client
			details, err := getDetails(symbol)
			if err != nil {
				fmt.Println(err)
				return c.Send("Something went wrong getting your stock")
			}

			// get chart
			graphBuff, err := getChart(symbol, nil, nil, datetime.OneHour)
			if err != nil {
				fmt.Println(err)
				return c.Send("Something went wrong getting your stock")
			}

			// return results
			a := &telebot.Photo{File: telebot.FromReader(graphBuff)}

			if err := c.Send(details, selector); err != nil {
				fmt.Println(err)
				c.Send("Something went wrong getting your stock")
			}
			return c.SendAlbum(telebot.Album{a})
		}

		cmds[len(cmds)-1] = append(cmds[len(cmds)-1], core.InlineCommand{
			Name:        name,
			Description: d,
			Handler:     handler,
		})
	}

	// The bottom row is always the Add button, full width
	addCommand := core.InlineCommand{
		Name:        "add",
		Description: "➕ Add",
		Handler: func(c telebot.Context) error {
			stock.b.SetOnText(func(c telebot.Context) error {
				symbol := strings.ToUpper(c.Text())

				// library returns nil for non-existent stocks
				if q, err := quote.Get(symbol); err != nil || q == nil {
					return c.Send("Sorry, I couldn't find it.")
				}

				stock.symbols = append(stock.symbols, symbol)
				return c.Send("Added!")
			}, 1*time.Minute)
			return c.Send("What stock would you like to add?")
		},
	}
	cmds = append(cmds, []core.InlineCommand{addCommand})

	return cmds, nil
}