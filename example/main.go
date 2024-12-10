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
