package internal

import (
	"fmt"
	"net"
	"runtime"
)

type TopicMessage []byte

type Topic struct {
	name        string
	connections []*Connections
	msgch       chan TopicMessage
}

type PublishService struct {
	tmap map[string]*Topic
}

func NewPublishService() *PublishService {
	return &PublishService{
		tmap: map[string]*Topic{},
	}
}

func (ps *PublishService) AddConnectionToTopic(topic string, addr net.Addr) {
	if ad, ok := addr.(*net.TCPAddr); ok {
		if _, k := ps.tmap[topic]; !k {
			t := &Topic{
				name:        topic,
				connections: []*Connections{},
				msgch:       make(chan TopicMessage, 1),
			}
			ps.tmap[topic] = t
			go ps.tmap[topic].ReceiveMsg()
		}
		ps.tmap[topic].connections = append(ps.tmap[topic].connections, &Connections{
			Addr: ad.IP.String(),
			Port: ":5959",
		})

	}
	fmt.Println(ps.tmap)

}

// Publish data to subscribed remote Addresses
// Data should be serialised for each topic

func (ps *PublishService) Publish(topic string, data []byte) error {

	if _, ok := ps.tmap[topic]; !ok {

		td := &Topic{
			name:        topic,
			connections: []*Connections{},
			msgch:       make(chan TopicMessage, 1),
		}
		ps.tmap[topic] = td
		go ps.tmap[topic].ReceiveMsg()
	}
	ps.tmap[topic].msgch <- TopicMessage(data)
	conn := ps.tmap[topic].connections
	for i := range conn {
		conn[i].Send(data)
	}
	return nil
}

func (t *Topic) ReceiveMsg() {
	for msg := range t.msgch {
		for i := range t.connections {
			t.connections[i].Send(msg)
			runtime.Gosched()
		}
	}
}
