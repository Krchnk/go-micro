package notify

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type UserEvent struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func SendUserRegistered(broker, topic string, evt UserEvent) error {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Async:    false,
	})
	defer w.Close()
	value, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	msg := kafka.Message{Value: value}
	if err := w.WriteMessages(context.Background(), msg); err != nil {
		log.Printf("Kafka write error: %v", err)
		return err
	}
	return nil
}
