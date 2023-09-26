package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

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
		s.process(msg[:n], conn)

	}
}

func (s *Server) process(msg []byte, conn net.Conn) error {
	var data struct {
		Type  string `json:"type"`
		Topic string `json:"topic"`
		Value string `json:"value"`
	}
	fmt.Println(string(msg))
	if err := json.Unmarshal(msg, &data); err != nil {

		return err
	}

	switch data.Type {
	case "publish":
		// How to add connection here

		s.ps.Publish(data.Topic, []byte(data.Value))
		conn.Write([]byte("ok"))
	case "subscribe":
		s.ps.AddConnectionToTopic(data.Topic, conn.RemoteAddr())
	}
	return nil
}
func main() {
	nt := internal.NewPublishService()
	server := NewServer("localhost:3000", nt)

	log.Fatal(server.Start())
}
