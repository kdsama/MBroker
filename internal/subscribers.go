package internal

import (
	"fmt"
	"net"
	"sync"
)

type Connections struct {
	Addr string
	Port string
}

type ConnectionService struct {
	conMux        *sync.Mutex
	connectionMap map[string]*Connections
}

func New() *ConnectionService {
	return &ConnectionService{
		conMux:        &sync.Mutex{},
		connectionMap: map[string]*Connections{},
	}
}

func (s *ConnectionService) Add(addr net.Addr) error {

	if ad, ok := addr.(*net.TCPAddr); ok {
		s.connectionMap[addr.String()] = &Connections{
			Addr: ad.IP.String(),
			Port: "5959",
		}
		return nil
	}

	return nil
}

func (c *Connections) Send(data []byte) error {
	fmt.Println("Are we coming here or not ???/")
	conn, e := net.Dial("tcp", c.Addr+c.Port)
	if e != nil {
		fmt.Println(e)
		return e
	}
	_, e = conn.Write(data)
	return e
}
