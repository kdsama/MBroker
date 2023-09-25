package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	q          chan struct{}
}

func NewServer(addr string) *Server {
	return &Server{
		listenAddr: addr,
		q:          make(chan struct{}),
	}
}
func (s *Server) Start() error {
	ln, e := net.Listen("tcp", s.listenAddr)
	if e != nil {
		log.Fatal(e)
	}
	s.ln = ln
	go s.AcceptLoop()
	<-s.q
	return nil
}

func (s *Server) AcceptLoop() {

	defer s.ln.Close()

	for {
		c, e := s.ln.Accept()
		if e != nil {
			fmt.Println(e)
			continue
		}
		go s.ReadLoop(c)
	}
}

func (s *Server) ReadLoop(conn net.Conn) {
	defer conn.Close()
	msg := make([]byte, 1024)
	for {
		n, e := conn.Read(msg)
		if e != nil {
			continue
		}
		fmt.Println("data")
		fmt.Println(string(msg[:n]))

	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go fn()
	conn.Write([]byte("subscribe topic1"))
	time.Sleep(1 * time.Second)
	conn.Write([]byte("publish topic1 what the hell"))

	for {
		i := 0
		i++
		time.Sleep(1 * time.Second)
	}
}
func fn() {
	s := NewServer("localhost:5959")
	log.Fatal(s.Start())
}
