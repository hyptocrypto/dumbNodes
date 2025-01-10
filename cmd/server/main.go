package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/hyptocrypto/dumbNodes/internal/config"
	"github.com/hyptocrypto/dumbNodes/internal/mproto"
	"github.com/hyptocrypto/dumbNodes/internal/types"
	"github.com/hyptocrypto/dumbNodes/internal/util"
)

func main() {
	s, err := NewServer()
	if err != nil {
		fmt.Printf("Error creating server %v", err)
		os.Exit(1)
	}
	s.Listen()
}

type Server struct {
	mu          sync.RWMutex
	sock        net.Listener
	connections types.ClientConnections
	closeChn    chan bool
}

// Create new server with default values
func NewServer() (*Server, error) {
	conf := config.NewConfig()
	conn, err := net.Listen(conf.Protocol, conf.FullHost())
	if err != nil {
		return nil, fmt.Errorf("could not create listener for server: %v", err)
	}
	return &Server{mu: sync.RWMutex{}, sock: conn, connections: make(types.ClientConnections), closeChn: make(chan bool)}, nil
}

func (s *Server) GetOrAddClientConnection(host string) (*types.ClientConn, error) {
	clientID := util.GenerateUUIDForClient(host)
	s.mu.RLock()
	cliCon, ok := s.connections[clientID]
	s.mu.RUnlock()
	if ok {
		return cliCon, nil
	}
	clientConn := types.ClientConn{
		ClientId: clientID,
		Key:      "This is a key",
	}
	s.mu.Lock()
	s.connections[clientID] = &clientConn
	s.mu.Unlock()
	fmt.Printf("Added new client connection: %v\n", clientID)
	return &clientConn, nil
}

// Close net connection and reset connections
func (s *Server) Listen() error {
	fmt.Printf("Listening on: %v\n", s.sock.Addr().String())
	for {
		select {
		case <-s.closeChn:
			fmt.Println("Server is shutting down...")
			s.close()
			return nil
		default:
			conn, err := s.sock.Accept()
			if os.IsTimeout(err) {
				fmt.Printf("Error accepting connection: %v\n", err)
				continue
			}
			fmt.Printf("Accepted new connection from: %v. ClientID: %v\n", conn.RemoteAddr().String(), s.praseClientFromConn(conn))

			go s.handleClientConnection(conn)

			// Avoid hot looping
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (s *Server) praseClientFromConn(c net.Conn) *types.ClientConn {
	id := util.GenerateUUIDForClient(c.RemoteAddr().String())
	if c, ok := s.connections[id]; ok {
		return c
	}
	clientConn := types.ClientConn{
		ClientId: id,
		Key:      "This is a key",
	}
	s.connections[id] = &clientConn
	return &clientConn
}

func (s *Server) handleClientConnection(conn net.Conn) error {
	defer conn.Close()
	for {
		req, err := mproto.DeserializetReqeustFromConn(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected.")
				return nil
			}
			fmt.Printf("Error reading request: %v\n", err)
			return err
		}

		fmt.Printf("Request parsed: %v\n", *req)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		resp, err := mproto.Serialize(types.Response{
			Headers: map[string]string{"hello": "world"},
			Data:    []byte("hello"),
			Source:  req.Destination,
		})
		if err != nil {
			return fmt.Errorf("failed to serialize response: %w", err)
		}

		_, err = conn.Write(resp)
		if err != nil {
			fmt.Printf("Error writing response: %v\n", err)
			return err
		}
	}
}

// Close tcp listener and client channels
func (s *Server) close() error {
	if err := s.sock.Close(); err != nil {
		return err
	}
	return nil
}
