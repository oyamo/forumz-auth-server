package pkg

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"time"
)

type JSONSender struct {
	consumer string
}

func (s *JSONSender) Send(topic string, v any) error {
	m, err := json.Marshal(v)
	if err != nil {
		return err
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", s.consumer, topic, 0)
	if err != nil {
		return err
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return err
	}

	_, err = conn.WriteMessages(
		kafka.Message{Value: m})

	if err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		return err
	}

	return nil
}

func NewJSONSender(conf *Config) *JSONSender {
	return &JSONSender{conf.KafkaProducerServer}
}
