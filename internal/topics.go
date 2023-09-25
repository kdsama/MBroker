package internal

import "fmt"

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
func (ps *PublishService) AddTopic(topic string, connections []*Connections) {
	td := &Topic{
		name:        topic,
		connections: connections,
	}
	ps.tmap[topic] = td
}

func (ps *PublishService) AddConnectionToTopic(topic string, connection *Connections) {

	ps.tmap[topic].connections = append(ps.tmap[topic].connections, connection)
}

func (ps *PublishService) Publish(topic string, data []byte) error {

	if _, ok := ps.tmap[topic]; !ok {
		return fmt.Errorf("format, %d", 1)
	}
	conn := ps.tmap[topic].connections
	for i := range conn {
		conn[i].Send(data)
	}
	return nil
}
