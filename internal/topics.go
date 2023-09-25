package internal

import (
	"fmt"
	"net"
)

type Topic struct {
	name        string
	connections []*Connections
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
			t := &Topic{name: topic, connections: []*Connections{}}
			ps.tmap[topic] = t
		}
		ps.tmap[topic].connections = append(ps.tmap[topic].connections, &Connections{
			Addr: ad.IP.String(),
			Port: ":5959",
		})

	}
	fmt.Println(ps.tmap)

}

// Publish data to subscribed remote Addresses
func (ps *PublishService) Publish(topic string, data []byte) error {

	if _, ok := ps.tmap[topic]; !ok {

		td := &Topic{
			name:        topic,
			connections: []*Connections{},
		}
		ps.tmap[topic] = td
	}

	conn := ps.tmap[topic].connections
	for i := range conn {
		conn[i].Send(data)
	}
	return nil
}
