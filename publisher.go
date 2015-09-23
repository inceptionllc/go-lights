package lights

import (
	"os"

	"github.com/bitly/go-nsq"
)

// Publisher creates a long lived queue massage publishing object.
type Publisher struct {
	ID       *DeviceID     // ID is the device ID for the publisher
	config   *nsq.Config   // The nsq configuration for the publisher
	producer *nsq.Producer // Underlying producer for this publisher
}

// NewPublisher creates a new queue publisher. If a device ID is
// provided, it will be used otherwise lights.NewID() will be used
// to automatically assign a device ID for the publisher.
func NewPublisher(id ...*DeviceID) (*Publisher, error) {

	config := nsq.NewConfig()

	var pub *Publisher

	if len(id) > 0 {
		pub = &Publisher{ID: id[0], config: config}
	} else {
		// Automatically add an ID if none was provided.
		did, err := NewID()
		if err != nil {
			return nil, err
		}
		pub = &Publisher{ID: did, config: config}
	}
	host := os.Getenv("INC_NSQD")
	if host == "" {
		host = "127.0.0.1:4150"
	}
	producer, err := nsq.NewProducer(host, config)
	if err != nil {
		return nil, err
	}
	pub.producer = producer
	return pub, nil
}

// Send posts a message to the message queue on the given topic. The
// topic name is automatically prepended with the device ID followed by a `.`
// e.g. `abc123.topic`.
func (p *Publisher) Send(topic, message string) {
	p.producer.PublishAsync(p.ID.ID+"."+topic, []byte(message), nil)
}
