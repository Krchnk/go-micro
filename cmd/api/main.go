package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Krchnk/go-micro/internal/config"
	"github.com/Krchnk/go-micro/internal/httpapi"
	"github.com/Krchnk/go-micro/internal/users"
)

func main() {
	cfg := config.Load()
	repo := users.NewInMemoryRepository()
	service := users.NewService(repo)
	handler := httpapi.NewHandler(service)

	server := &http.Server{
		Addr:              cfg.HTTPAddr(),
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("user service is running on %s", cfg.HTTPAddr())
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
