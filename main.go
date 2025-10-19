package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type DB interface {
	Set(key string, value []byte)
	Get(key string) ([]byte, bool)
	Start() error
}

type Config struct {
	Host string
	Port string
}

type Store struct {
	Config

	mu   sync.RWMutex
	data map[string][]byte
}

func NewStore(conf Config) DB {
	return &Store{
		Config: conf,
		data:   make(map[string][]byte),
	}
}

func (s *Store) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.Host, s.Port))
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("TCP accept error: %s", err)
		}

		log.Printf("[CONNECTION] New client from %s", conn.RemoteAddr())

		go s.handleConn(conn)
	}
}

func (s *Store) handleConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			return
		}

		cmd, args, err := s.parseCommands(msg)
		if err != nil {
			fmt.Printf("TCP parse error: %s\n", err)
		}

		switch cmd {
		case "SET", "set":
			if len(args) < 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'set' command\r\n"))
			} else {
				s.Set(args[0], []byte(args[1]))
				conn.Write([]byte("+OK\r\n"))
			}
		case "GET", "get":
			val, found := s.Get(args[0])
			if found {
				conn.Write(fmt.Appendf(nil, "%s=%s\n", args[0], string(val)))
			}
		default:
			conn.Write(fmt.Appendf(nil, "-ERR unknown command!\n"))
		}

		log.Printf("[COMMAND] %s %v", cmd, args)
	}
}

func (s *Store) parseCommands(msg string) (string, []string, error) {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return "", nil, fmt.Errorf("empty command")
	}

	splittedStr := strings.Split(msg, " ")
	if len(splittedStr) < 2 {
		return "", nil, fmt.Errorf("invalid Command")
	}
	cmd := strings.ToUpper(splittedStr[0])
	args := splittedStr[1:]

	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}

	return cmd, args, nil
}

func (s *Store) Set(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	if !ok {
		return nil, false
	}
	return val, true
}

func main() {
	conf := Config{
		Host: "localhost",
		Port: "3000",
	}

	s := NewStore(conf)

	err := s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
