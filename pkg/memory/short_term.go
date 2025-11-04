package memory

import (
	"time"
)

type Message struct {
	UserID         int64
	ChatID         int64
	FromCurrentBot bool
	Text           string
	Timestamp      time.Time
}

type ShortTermMemory struct {
	buffer map[int64][]Message
	limit  int
}

func NewShortTermMemory(limit int) *ShortTermMemory {
	return &ShortTermMemory{
		buffer: make(map[int64][]Message),
		limit:  limit,
	}
}

func (m *ShortTermMemory) Add(msg Message) {
	msgs := append(m.buffer[msg.ChatID], msg)
	if len(msgs) > m.limit {
		msgs = msgs[len(msgs)-m.limit:]
	}
	m.buffer[msg.ChatID] = msgs
}

func (m *ShortTermMemory) ChatHistory(chatID int64) []Message {
	// return copy of short memory state
	return append(m.buffer[chatID])
}
