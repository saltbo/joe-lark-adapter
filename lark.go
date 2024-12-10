package lark

type Message struct {
	ID string
	TextMessage
}

type TextMessage struct {
	Text string `json:"text"`
}
