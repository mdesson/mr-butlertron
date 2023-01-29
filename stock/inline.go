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

func sendGraph(c telebot.Context, symbol string, start *time.Time, end *time.Time, interval datetime.Interval, dateFormatString string) error {
	graphBuff, err := getChart(symbol, start, end, interval, dateFormatString)
	if err != nil {
		fmt.Println(err)
		return c.Send("Something went wrong getting your stock")
	}

	// return results
	a := &telebot.Photo{File: telebot.FromReader(graphBuff)}

	return c.SendAlbum(telebot.Album{a})
}

func makeDeleteHandler(stock *Stock, symbol string, i int) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		for j := 0; j < len(stock.symbols); j++ {
			if stock.symbols[j] == symbol {
				stock.symbols = append(stock.symbols[:i], stock.symbols[i+1:]...)
			}
		}
		return c.Send("Deleted!")
	}
}

func makeAddHandler(stock *Stock) func(c telebot.Context) error {
	return func(c telebot.Context) error {
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
	}
}

func makeInlineHandler(stock *Stock, s string, i int) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		// keyboard and handlers setup
		deleteHandler := makeDeleteHandler(stock, s, i)

		// get yesterday as starting point for all historical data
		tz, err := time.LoadLocation("America/New_York")
		if err != nil {
			fmt.Println(err)
			return c.Send("Something went wrong getting your stock")
		}
		yesterday := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 59, 59, 59, 59, tz)
		oneMonth := yesterday.AddDate(0, -1, 0)
		threeMonths := yesterday.AddDate(0, -3, 0)
		sixMonths := yesterday.AddDate(0, -6, 0)
		oneYear := yesterday.AddDate(-1, 0, 0)

		// add inline keyboard for historical data and delete
		selector := stock.b.RegisterInlineKeyboard([][]core.InlineCommand{
			{
				core.InlineCommand{
					Name:        fmt.Sprintf("monthly-%s", s),
					Description: "📅 1mo",
					Handler: func(c telebot.Context) error {
						return sendGraph(c, s, &oneMonth, &yesterday, datetime.OneDay, "2006-01-02")
					},
				},
				core.InlineCommand{
					Name:        fmt.Sprintf("quarterly-%s", s),
					Description: "📅 3mo",
					Handler: func(c telebot.Context) error {
						return sendGraph(c, s, &threeMonths, &yesterday, datetime.FiveDay, "2006-01-02")
					},
				},
				core.InlineCommand{
					Name:        fmt.Sprintf("biannually-%s", s),
					Description: "📅 6mo",
					Handler: func(c telebot.Context) error {
						return sendGraph(c, s, &sixMonths, &yesterday, datetime.FiveDay, "2006-01-02")
					},
				},
				core.InlineCommand{
					Name:        fmt.Sprintf("annualy-%s", s),
					Description: "📅 1y",
					Handler: func(c telebot.Context) error {
						return sendGraph(c, s, &oneYear, &yesterday, datetime.OneMonth, "2006-01-02")
					},
				},
			},
			{
				core.InlineCommand{
					Name:        "delete",
					Description: "🗑️ Delete",
					Handler:     deleteHandler,
				},
			},
		})

		// get text to return to client
		details, err := getDetails(s)
		if err != nil {
			fmt.Println(err)
			return c.Send("Something went wrong getting your stock")
		}

		if err := c.Send(details, selector); err != nil {
			fmt.Println(err)
			c.Send("Something went wrong getting your stock")
		}

		return sendGraph(c, s, nil, nil, datetime.OneHour, "3:04pm")
	}
}

func createInlineHandlers(stock *Stock) ([][]core.InlineCommand, error) {
	// get descriptions in parallel
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	var descErr error

	descriptions := make(map[string]string)

	for _, symbol := range stock.symbols {
		wg.Add(1)
		go func(s string) {
			d, err := getDescription(s)
			if err != nil {
				descErr = err
			} else {
				mu.Lock()
				descriptions[s] = d
				mu.Unlock()
			}
			wg.Done()
		}(symbol)
	}
	wg.Wait()

	// each symbol must have both a quote and a corresponding ticker
	if descErr != nil {
		return nil, descErr
	}

	// add inline commands
	cmds := make([][]core.InlineCommand, 0)
	for i, symbol := range stock.symbols {
		if i%2 == 0 {
			cmds = append(cmds, []core.InlineCommand{})
		}

		cmds[len(cmds)-1] = append(cmds[len(cmds)-1], core.InlineCommand{
			Name:        fmt.Sprintf("stock-%d", i),
			Description: descriptions[symbol],
			Handler:     makeInlineHandler(stock, symbol, i),
		})
	}

	// The bottom row is always the Add button, full width
	addCommand := core.InlineCommand{
		Name:        "add",
		Description: "➕ Add",
		Handler:     makeAddHandler(stock),
	}
	cmds = append(cmds, []core.InlineCommand{addCommand})

	return cmds, nil
}
