package internal

import (
	"net"
	"sync"
)

type Connections struct {
	conn *net.Conn
	name string
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
	conn, err := net.Dial(addr.Network(), addr.String())
	if err != nil {
		return err
	}

	s.connectionMap[addr.String()] = &Connections{
		conn: &conn,
		name: addr.String(),
	}

	return nil
}

func (c *Connections) Send(data []byte) error {
	_, e := (*c.conn).Write(data)
	return e
}
