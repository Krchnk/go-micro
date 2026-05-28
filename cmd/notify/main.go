package main

import (
	"log"

	"github.com/Krchnk/go-micro/internal/config"
	"github.com/Krchnk/go-micro/internal/notify"
)

func main() {
	cfg := config.Load()
	log.Printf("Notify service listening on topic %s", cfg.KafkaTopic)
	notify.StartConsumer(cfg.KafkaBroker, cfg.KafkaTopic, cfg.KafkaGroup, func(userID int64, email, name string) {
		log.Printf("Send notification to %s <%s> (id=%d)", name, email, userID)
		// Здесь может быть отправка email, push и т.д.
	})
}
