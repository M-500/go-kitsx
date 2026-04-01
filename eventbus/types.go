package eventbus

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Message 是统一的事件信封。
// 传输层只关心 Message，本地、channel、kafka 使用相同的数据结构。
type Message struct {
	ID        string            `json:"id"`
	Topic     string            `json:"topic"`
	Key       string            `json:"key,omitempty"`
	Source    string            `json:"source,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	Headers   map[string]string `json:"headers,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Payload   json.RawMessage   `json:"payload,omitempty"`
}

// NewMessage 创建原始消息。
func NewMessage(topic string, payload []byte) *Message {
	return &Message{
		ID:        uuid.NewString(),
		Topic:     topic,
		Timestamp: time.Now().UTC(),
		Payload:   append([]byte(nil), payload...),
		Headers:   make(map[string]string),
		Metadata:  make(map[string]string),
	}
}

// NewJSONMessage 创建 JSON 消息。
func NewJSONMessage(topic string, payload any) (*Message, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return NewMessage(topic, raw), nil
}

// Clone 复制消息，避免多订阅者共享可变底层数据。
func (m *Message) Clone() *Message {
	if m == nil {
		return nil
	}

	clone := &Message{
		ID:        m.ID,
		Topic:     m.Topic,
		Key:       m.Key,
		Source:    m.Source,
		Timestamp: m.Timestamp,
		Payload:   append([]byte(nil), m.Payload...),
	}
	if len(m.Headers) > 0 {
		clone.Headers = make(map[string]string, len(m.Headers))
		for k, v := range m.Headers {
			clone.Headers[k] = v
		}
	} else {
		clone.Headers = make(map[string]string)
	}
	if len(m.Metadata) > 0 {
		clone.Metadata = make(map[string]string, len(m.Metadata))
		for k, v := range m.Metadata {
			clone.Metadata[k] = v
		}
	} else {
		clone.Metadata = make(map[string]string)
	}
	return clone
}

// DecodeJSON 将消息载荷解码到目标对象。
func (m *Message) DecodeJSON(target any) error {
	return json.Unmarshal(m.Payload, target)
}
