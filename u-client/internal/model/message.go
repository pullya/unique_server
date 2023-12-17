package model

import (
	"encoding/json"
	"fmt"
)

const (
	MessageTypeRequest  = "Request"
	MessageTypeResponse = "Response"
)

type Message struct {
	MessageType   string `json:"messageType"`
	MessageString string `json:"messageString"`
}

func (m Message) AsJson() []byte {
	res, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("failed to marshall message: %v", err)
		return nil
	}
	return res
}

func PrepareMessage(mType string, mString string) Message {
	return Message{
		MessageType:   mType,
		MessageString: mString,
	}
}
