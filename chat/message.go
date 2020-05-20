package chat

import "fmt"

// Message message format structure
type Message struct {
	UserName  string `json:"userName"` // Email
	Body      string `json:"body"`
	Timestamp string `json:"timestamp"`
}

func (m *Message) String() string {
	return fmt.Sprintf("%s at %s sent %s", m.UserName, m.Timestamp, m.Body)
}
