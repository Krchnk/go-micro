package main

import (
	"log"
	"net/http"

	"github.com/Krchnk/go-micro/internal/ws"
)

func main() {
	server := ws.NewServer()
	go server.Run()

	http.HandleFunc("/ws", server.HandleWS)
	log.Println("WebSocket server started on :8080/ws")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
