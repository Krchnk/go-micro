package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Server struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.mu.Unlock()
		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
			s.mu.Unlock()
		case message := <-s.broadcast:
			s.mu.Lock()
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256)}
	s.register <- client
	go client.writePump()
	client.readPump(s)
}

func (c *Client) readPump(s *Server) {
	defer func() {
		s.unregister <- c
		c.conn.Close()
	}()
	for {
		msgType, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		// Пример: обработка разных типов сообщений
		if msgType == websocket.TextMessage {
			log.Printf("Text message: %s", message)
			s.broadcast <- message
		} else if msgType == websocket.BinaryMessage {
			log.Printf("Binary message: %v bytes", len(message))
			s.broadcast <- message
		}
	}
}

func (c *Client) writePump() {
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
	c.conn.Close()
}
