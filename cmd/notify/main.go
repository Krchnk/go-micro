package main

import (
	"log"
	"os"

	"github.com/Krchnk/go-micro/internal/notify"
)

func main() {
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	topic := getEnv("KAFKA_TOPIC", "user-registered")
	group := getEnv("KAFKA_GROUP", "notify-service")
	log.Printf("Notify service listening on topic %s", topic)
	notify.StartConsumer(broker, topic, group, func(userID int64, email, name string) {
		log.Printf("Send notification to %s <%s> (id=%d)", name, email, userID)
		// Здесь может быть отправка email, push и т.д.
	})
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
