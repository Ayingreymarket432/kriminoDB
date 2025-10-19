package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

type DB interface {
	Set(key string, value []byte)
	Get(key string) ([]byte, bool)
	ListenAndAccept() error
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
	}
}

func (s *Store) ListenAndAccept() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.Host, s.Port))
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("TCP accept error: %s", err)
		}

		fmt.Printf("new incoming connection %+v\n", conn)

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
		fmt.Printf("Message incoming: %s\n", string(msg))
	}
}

func (s *Store) Set(key string, value []byte) {
}

func (s *Store) Get(key string) ([]byte, bool) {
	return nil, false
}

func main() {
	conf := Config{
		Host: "localhost",
		Port: "3000",
	}

	s := NewStore(conf)

	err := s.ListenAndAccept()
	if err != nil {
		log.Fatal(err)
	}
}
