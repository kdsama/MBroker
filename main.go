package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/kdsama/mbroker/internal"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	q          chan struct{}
	ps         *internal.PublishService
}

func NewServer(addr string, ps *internal.PublishService) *Server {
	return &Server{
		listenAddr: addr,
		q:          make(chan struct{}),
		ps:         ps,
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
		e = s.process(string(msg[:n]))
		if e != nil {
			conn.Write([]byte("Fuck you\n"))
		}
		conn.Write([]byte("Thank you"))
	}
}

func (s *Server) process(msg string) error {
	msgFields := strings.Split(msg, "")
	if len(msgFields) <= 1 {

		return fmt.Errorf("What are you on about ")
	}
	switch msgFields[0] {
	case "add":
		// How to add connection here
		s.ps.AddTopic(msgFields[1], nil)
	}
	return nil
}
func main() {
	nt := internal.NewPublishService()
	server := NewServer("localhost:3000", nt)

	log.Fatal(server.Start())
}
