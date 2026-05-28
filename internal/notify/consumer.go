package notify

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(broker, topic, groupID string, notifyFunc func(userID int64, email, name string)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		GroupID:  groupID,
		Topic:    topic,
	})
	defer r.Close()
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			continue
		}
		// Пример: ожидаем JSON {"id":int64, "email":string, "name":string}
		var evt struct {
			ID    int64  `json:"id"`
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		if err := json.Unmarshal(m.Value, &evt); err != nil {
			log.Printf("Invalid event: %v", err)
			continue
		}
		notifyFunc(evt.ID, evt.Email, evt.Name)
	}
}
