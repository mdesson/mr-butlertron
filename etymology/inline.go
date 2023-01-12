package etymology

import (
	"fmt"
	"github.com/mdesson/mr-butlertron/core"
	"gopkg.in/telebot.v3"
)

func createInlineHandlers(definitions map[string]string) [][]core.InlineCommand {
	cmds := make([][]core.InlineCommand, 0)
	counter := 0
	for word, definition := range definitions {
		if counter%3 == 0 {
			cmds = append(cmds, []core.InlineCommand{})
		}

		name := fmt.Sprintf("etym-%d", counter)
		// capture loop variable values instead of final reference
		handler := func(w, d string) func(c telebot.Context) error {
			return func(c telebot.Context) error {
				txt := fmt.Sprintf("*%s*\n\n%s", w, d)
				return c.Send(txt)
			}
		}(word, definition)

		cmds[len(cmds)-1] = append(cmds[len(cmds)-1], core.InlineCommand{
			Name:        name,
			Description: word,
			Handler:     handler,
		})

		counter += 1
	}
	return cmds
}
