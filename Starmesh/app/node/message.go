package node

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"starmesh/common/model"

	"github.com/google/uuid"
)

const (
	MaxRetry     = 3               // 최대 재시도 횟수
	RetryBackoff = 2 * time.Second // 재시도 간격
)

type Message struct {
	MessageID string       `json:"message_id"`
	Route     []model.Peer `json:"route"`
	Payload   string       `json:"payload"`
	HopIndex  int          `json:"hop_index"`
	Trace     []string     `json:"trace"` // 경유 노드 기록
}

// NewMessage creates a new message with a UUID
func NewMessage(route []model.Peer, payload string) *Message {
	return &Message{
		MessageID: uuid.New().String(),
		Route:     route,
		Payload:   payload,
		HopIndex:  0,
		Trace:     []string{},
	}
}

// SendToNextHop sends the message to the next node in the route
func (m *Message) SendToNextHop() error {
	if m.HopIndex >= len(m.Route) {
		fmt.Printf("[Relay] Reached end of route. Delivering payload (%s): %s\n", m.MessageID, m.Payload)
		SaveTraceToFile(m)
		return nil
	}

	next := m.Route[m.HopIndex]
	address := next.IP + ":" + next.Port

	var lastErr error
	for i := 0; i < MaxRetry; i++ {
		fmt.Printf("[Relay] [%s] Attempt %d: forwarding to %s (hop %d)\n", m.MessageID, i+1, address, m.HopIndex)

		conn, err := net.Dial("tcp", address)
		if err == nil {
			defer conn.Close()
			data, _ := json.Marshal(m)
			conn.Write(data)
			return nil
		}

		lastErr = err
		time.Sleep(RetryBackoff)
	}

	// 모든 재시도 실패
	fmt.Printf("[Relay] [%s] Failed to forward to %s after %d attempts: %v\n", m.MessageID, address, MaxRetry, lastErr)
	m.Trace = append(m.Trace, "FAILED:"+address)
	SaveTraceToFile(m)
	return lastErr
}
