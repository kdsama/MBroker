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
		s.process(string(msg[:n]), conn)

	}
}

func (s *Server) process(msg string, conn net.Conn) error {
	msgFields := strings.Split(msg, " ")
	if len(msgFields) <= 1 {

		return fmt.Errorf("What are you on about ")
	}
	fmt.Println(msgFields[0], msg)
	switch msgFields[0] {
	case "publish":
		// How to add connection here
		fmt.Println("We are gonna publish from here these things ", msgFields[1], []byte(msg[len(msgFields[1])+len(msgFields[2])-1:]))
		s.ps.Publish(msgFields[1], []byte(msg[len(msgFields[1])+len(msgFields[2])-1:]))
		conn.Write([]byte("ok"))
	case "subscribe":
		s.ps.AddConnectionToTopic(msgFields[1], conn.RemoteAddr())
	}
	return nil
}
func main() {
	nt := internal.NewPublishService()
	server := NewServer("localhost:3000", nt)

	log.Fatal(server.Start())
}
