# mr-butlertron
An extensible do-anything telegram bot.

## Adding a New Command

Simply implement the `core.Command` interface and away you go! Note that `Description()` will appear in any menus as the text.

### Commands with Inline Keyboards

The recommended method would be to declare a series of `core.InlineCommand` variables as package variables (no need to export). Then, in your Command, create a `New` function which will call the butlertron's `RegisterInlineKeyboard` function.

Now, inside of the program's init function you can simply call `New` when you register it with the bot and this will all be done for you.

For example:

```go
// foo/foo.go (package foo)

type Foo struct {
	selector telebot.ReplyMarkup
}

// snip implementing the Command interface // 
func New(b *core.Butlertron) {
	f := Foo{}
	f.Selector =  b.RegisterInlineKeyboard(inlineHandlers)

	return &f
}

var inlineHandlers = [][]core.InlineCommand{
	{
		core.InlineCommand{
			Name:        "foo",
			Description: "🔧 foo",
			Handler:     FooHandler,
		},
		core.InlineCommand{
			Name:        "bar",
			Description: "🔨 bar",
			Handler:     BarHandler,
		},
	},
	{
		core.InlineCommand{
			Name:        "baz",
			Description: "🪛 baz",
			Handler:     BazHandler,
		},
	},
}

// Inline Keyboard Handlers
func FooHandler(c telebot.Context) error {
	return c.Send("Fooey!")
}

func BarHandler(c telebot.Context) error {
	return c.Send("Babar!")
}

func BazHandler(c telebot.Context) error {
	return c.Send("Bazz Hands!")
}
```

