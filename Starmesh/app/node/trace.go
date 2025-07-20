package node

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// SaveTraceToFile appends trace to a local file for viewer
func SaveTraceToFile(m *Message) {
	filename := "trace_log.json"

	entry := map[string]interface{}{
		"timestamp":  time.Now().Format(time.RFC3339),
		"message_id": m.MessageID,
		"path":       m.Trace,
		"payload":    m.Payload,
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("[Trace] Failed to open trace file:", err)
		return
	}
	defer file.Close()

	data, _ := json.Marshal(entry)
	file.Write(append(data, '\n'))
}
