package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Server represents the chat server
type Server struct {
	mu         sync.RWMutex
	conns      map[*websocket.Conn]*connection
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

type connection struct {
	send   chan []byte
	ws     *websocket.Conn
	server *Server
}

// writer writes messages to the WebSocket connection
func (c *connection) writer() {
	defer func() {
		c.server.unregister <- c.ws
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The channel has been closed by the server
				return
			}

			err := c.ws.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("Error writing message:", err)
				return
			}
		}
	}
}

// reader reads messages from the WebSocket connection
func (c *connection) reader() {
	defer func() {
		c.server.unregister <- c.ws
		c.ws.Close()
	}()

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("Error reading message:", err)
			}
			break
		}

		c.server.broadcast <- message
	}
}

// run is the main loop for the chat server
func (s *Server) run() {
	for {
		select {
		case conn := <-s.register:
			c := &connection{
				send:   make(chan []byte, 256),
				ws:     conn,
				server: s,
			}

			s.mu.Lock()
			s.conns[conn] = c
			s.mu.Unlock()

			fmt.Println("New connection")
		case conn := <-s.unregister:
			if c, ok := s.conns[conn]; ok {
				close(c.send)
				s.mu.Lock()
				delete(s.conns, conn)
				s.mu.Unlock()

				fmt.Println("Closed connection")
			}
		case message := <-s.broadcast:
			s.mu.RLock()

			for conn, c := range s.conns {
				select {
				case c.send <- message:
				default:
					close(c.send)
					delete(s.conns, conn)
				}
			}

			s.mu.RUnlock()

			fmt.Println("Sent message")
		}
	}
}

// ServeHTTP handles incoming HTTP requests
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to websocket", http.StatusInternalServerError)
		return
	}

	// Register the new WebSocket connection and defer the unregister operation
	s.register <- ws
	defer func() { s.unregister <- ws }()

	c := s.conns[ws]
	go c.writer()
	c.reader()
}

// NewServer creates a new instance of Server
func NewServer() *Server {
	return &Server{
		conns:      make(map[*websocket.Conn]*connection),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	} 
} 
	
func main() {
  server := NewServer()
  go server.run()
// Serve the web page and WebSocket endpoint
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
})
    http.Handle("/ws", server)

// Start the server
    fmt.Println("Starting server...")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
	fmt.Println("Error starting server:", err)
}
}