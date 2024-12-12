package lark

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/go-joe/joe"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
	"go.uber.org/zap"
)

func Adapter(appId, secret string) joe.Module {
	return joe.ModuleFunc(func(conf *joe.Config) error {
		adapter := &larkAdapter{
			joeConf: conf,
			appId:   appId,
			secret:  secret,
			logger:  conf.Logger("lark"),
			client:  lark.NewClient(appId, secret),
		}

		conf.SetAdapter(adapter)
		return nil
	})
}

type larkAdapter struct {
	joeConf *joe.Config

	appId, secret string
	logger        *zap.Logger
	client        *lark.Client
}

func (a larkAdapter) RegisterAt(brain *joe.Brain) {
	eh := dispatcher.NewEventDispatcher("", "")
	eh.OnP2MessageReceiveV1(func(ctx context.Context, data *larkim.P2MessageReceiveV1) error {
		larkcore.Prettify(data)
		ev := data.Event
		eventID := data.EventV2Base.Header.EventID
		rawMsg := ev.Message
		if rawMsg == nil {
			return nil
		}

		var m Message
		if err := json.Unmarshal([]byte(*ev.Message.Content), &m); err != nil {
			return err
		}

		m.ID = *rawMsg.MessageId
		a.logger.Debug("Received message",
			zap.String("event_id", eventID),
			zap.Stringp("message_id", rawMsg.MessageId),
		)
		brain.Emit(joe.ReceiveMessageEvent{
			ID:       m.ID,
			Text:     m.Text,
			AuthorID: *ev.Sender.SenderId.OpenId,
			Channel:  *rawMsg.ChatId,
			Data:     data.Body,
		})
		return nil
	})
	wsClient := larkws.NewClient(a.appId, a.secret, larkws.WithEventHandler(eh), larkws.WithLogLevel(larkcore.LogLevelDebug))
	go func() {
		if err := wsClient.Start(a.joeConf.Context); err != nil {
			a.joeConf.Logger("lark").Error("start failed", zap.Error(err))
			return
		}
	}()
}

func (a larkAdapter) Send(text, channel string) error {
	msgType := larkim.MsgTypeText
	content := larkim.NewMessageTextBuilder().Text(text).Build()
	if strings.HasPrefix(text, "{") && strings.HasSuffix(text, "}") {
		// support messageCard
		msgType = larkim.MsgTypeInteractive
		content = text
	}

	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(makeReceiveIdType(channel)).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(msgType).
			ReceiveId(channel).
			Content(content).
			Build(),
		).Build()
	resp, err := a.client.Im.Message.Create(context.Background(), req)
	if err != nil {
		return err
	} else if !resp.Success() {
		return resp.CodeError
	}

	return err
}

func makeReceiveIdType(channel string) string {
	if strings.HasPrefix(channel, "ou_") {
		return larkim.ReceiveIdTypeOpenId
	}

	return larkim.ReceiveIdTypeChatId
}

func (a larkAdapter) Close() error {
	return nil
}
