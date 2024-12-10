<h1 align="center">Joe Bot - Lark Adapter</h1>
<p align="center">Connecting joe with the Lark chat application. https://github.com/go-joe/joe</p>
<p align="center">
	<a href="https://github.com/saltbo/joe-lark-adapter/releases"><img src="https://img.shields.io/github/tag/go-joe/lark-adapter.svg?label=version&color=brightgreen"></a>
	<a href="https://github.com/saltbo/joe-lark-adapter/actions/workflows/test.yml"><img src="https://github.com/saltbo/joe-lark-adapter/actions/workflows/test.yml/badge.svg"></a>
	<a href="https://goreportcard.com/report/github.com/saltbo/joe-lark-adapter"><img src="https://goreportcard.com/badge/github.com/saltbo/joe-lark-adapter"></a>
	<a href="https://codecov.io/gh/saltbo/joe-lark-adapter"><img src="https://codecov.io/gh/saltbo/joe-lark-adapter/branch/master/graph/badge.svg"/></a>
	<a href="https://pkg.go.dev/github.com/saltbo/joe-lark-adapter"><img src="https://img.shields.io/badge/godoc-reference-blue.svg?color=blue"></a>
	<a href="https://github.com/saltbo/joe-lark-adapter/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-BSD--3--Clause-blue.svg"></a>
</p>

---

This repository contains a module for the [Joe Bot library][joe].

## Getting Started

This library is packaged as [Go module][go-modules]. You can get it via:

```
go get github.com/saltbo/joe-lark-adapter
```

### Example usage

In order to connect your bot to lark you can simply pass it as module when
creating a new bot:

```go
package main

import (
	"fmt"
	"log"

	"github.com/go-joe/joe"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"go.uber.org/zap"

	lark "github.com/saltbo/joe-lark-adapter"
)

func main() {
	bot := joe.New("example", joe.WithLogLevel(zap.DebugLevel),
		lark.Adapter("cli_a7d28f1f86b89013", "WTF1wrfQzyWk0kbsghTXzfgQx7rxuKqw"),
	)
	bot.Respond("ping", func(message joe.Message) error {
		card := larkcard.NewMessageCard()
		card.Config(larkcard.NewMessageCardConfig().WideScreenMode(true))
		card.Header(larkcard.NewMessageCardHeader().Title(larkcard.NewMessageCardPlainText().Content("pong")))
		card.Elements([]larkcard.MessageCardElement{larkcard.NewMessageCardMarkdown().Content("**Name**: abc")})
		content, err := card.String()
		if err != nil {
			return err
		}

		// reply a card message
		bot.Say(message.Channel, content)
		// reply a text message
		bot.Say(message.Channel, "pong")
		return nil
	})
	bot.Brain.RegisterHandler(func(ev joe.ReceiveMessageEvent) {
		fmt.Println(ev.Text, string(ev.Data.([]byte)))
	})
	if err := bot.Run(); err != nil {
		log.Fatalln(err)
		return
	}
}

```

## Built With

* [larksuite/oapi-sdk-go](github.com/larksuite/oapi-sdk-go/v3) - Lark API in Go
* [zap](https://github.com/uber-go/zap) - Blazing fast, structured, leveled logging in Go
* [testify](https://github.com/stretchr/testify) - A simple unit test library

## Contributing

The current implementation is rather minimal and there are many more features
that could be implemented on the lark adapter so you are highly encouraged to
contribute. If you want to hack on this repository, please read the short
[CONTRIBUTING.md](CONTRIBUTING.md) guide first.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available,
see the [tags on this repository][tags].

## Authors

- **Ambor** - *Initial work* - [saltbo](https://github.com/saltbo)

See also the list of [contributors][contributors] who participated in this project.

## License

This project is licensed under the BSD-3-Clause License - see the [LICENSE](LICENSE) file for details.

[go-modules]: https://github.com/golang/go/wiki/Modules
[tags]: https://github.com/saltbo/joe-lark-adapter/tags
[contributors]: https://github.com/saltbo/joe-lark-adapter/contributors