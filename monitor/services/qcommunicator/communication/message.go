package communication

import (
	"encoding/json"
	"strings"
)

type Message struct {
	Type         string `json:"type"`
	UserID       string `json:"id"`
	Body         []byte `json:"body"`
	Name         string `json:"name"`
	Extension    string `json:"extension"`
	NewName      string `json:"new_name"`
	NewExtension string `json:"new_extension"`
}

type MessageType string

const (
	LOOKUP = MessageType("lookup")
	REMOVE = MessageType("remove")
	RENAME = MessageType("rename")
	CREATE = MessageType("create")
)

func NewMessage(t MessageType, userID, name string, body []byte) Message {
	var extension = ""
	if arr := strings.Split(name, "."); len(arr) > 1 {
		extension = arr[1]
		name = arr[0]
	}

	return Message{
		Type:      string(t),
		UserID:    userID,
		Name:      name,
		Extension: extension,
		Body:      body,
	}
}

func Marshal(m Message) ([]byte, error) {
	return json.Marshal(m)
}

func Unmarshal(body []byte, m *Message) error {
	return json.Unmarshal(body, m)
}
