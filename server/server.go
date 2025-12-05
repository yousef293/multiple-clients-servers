package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

func main() {
	server := NewServer()
	go server.RunBroadcaster()
	server.Listen()
}

type Server struct {
	clients   map[int]chan string
	messages  chan Message
	mutex     sync.Mutex
	nextID    int
}

type Message struct {
	senderID int
	content  string
}

func NewServer() *Server {
	return &Server{
		clients:  make(map[int]chan string),
		messages: make(chan Message, 10),
		nextID:   0,
	}
}

func (s *Server) Listen() {
	listener, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Listener error:", err)
	}
	defer listener.Close()

	fmt.Println("Server running on port 1234...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}

		s.mutex.Lock()
		clientID := s.nextID
		s.nextID++
		clientChannel := make(chan string, 10)
		s.clients[clientID] = clientChannel
		s.mutex.Unlock()

		fmt.Fprintf(conn, "YOUR_ID:%d\n", clientID)

		s.messages <- Message{
			senderID: clientID,
			content:  fmt.Sprintf("User %d joined", clientID),
		}

		log.Printf("Client %d connected from %s", clientID, conn.RemoteAddr())

		go s.HandleReader(conn, clientID)
		go s.HandleWriter(conn, clientChannel)
	}
}

func (s *Server) HandleReader(conn net.Conn, clientID int) {
	defer func() {
		s.mutex.Lock()
		if ch, exists := s.clients[clientID]; exists {
			delete(s.clients, clientID)
			close(ch)
		}
		s.mutex.Unlock()

		s.messages <- Message{
			senderID: clientID,
			content:  fmt.Sprintf("User %d left", clientID),
		}
		conn.Close()
		log.Printf("Client %d disconnected", clientID)
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := fmt.Sprintf("User %d: %s", clientID, scanner.Text())
		s.messages <- Message{
			senderID: clientID,
			content:  msg,
		}
	}
}

func (s *Server) HandleWriter(conn net.Conn, ch chan string) {
	for msg := range ch {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			return
		}
	}
}

func (s *Server) RunBroadcaster() {
	for msg := range s.messages {
		s.mutex.Lock()
		for id, ch := range s.clients {
			if id != msg.senderID {
				select {
				case ch <- msg.content:
				default:
				}
			}
		}
		s.mutex.Unlock()
	}
}
