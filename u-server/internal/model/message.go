package model

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Response string `json:"response"`
	Text     string `json:"text"`
}

func (m *Message) AsJson() []byte {
	res, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("failed to marshall message: %v", err)
		return nil
	}
	return res
}

func PrepareMessage(response string, text string) Message {
	return Message{
		Response: response,
		Text:     text,
	}
}
