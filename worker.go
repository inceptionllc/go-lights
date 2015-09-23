package lights

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/bitly/go-nsq"
)

// WorkerFunc handles incoming string messages returning an error if the
// message could not be handled.
type WorkerFunc func(message string) error

// Worker takes care of the common nsq worker tasks that all message
// driven agents must carry out. The worker takes care of bootstrapping
// the system, and automatically configures nsq according to the current
// environment.
type Worker struct {
	ID        *DeviceID       // The ID for the hardware the worker is running on
	config    *nsq.Config     // The nsq configuration for the worker
	consumers []*nsq.Consumer // Consumers that have been created
	started   bool            // Flag true if the worker has been started
	stopped   []bool          // Flags for consumers that have been stopped
	stopOnce  sync.Once       // We only want to stop the consumers once!
	stopChan  chan bool       // Channel receives true when the worker should stop
}

// NewWorker creates a new worker ready for configuration. Call Start() on
// the worker to begin processing messages. Returns an error if there was a
// problem creating the worker. If an ID is provided the worker will use it,
// otherwise an ID will automatically be generated using lights.NewID().
func NewWorker(id ...*DeviceID) (*Worker, error) {

	config := nsq.NewConfig()

	if len(id) > 0 {
		return &Worker{ID: id[0], config: config}, nil
	}
	// Automatically add an ID if none was provided.
	did, err := NewID()
	if err != nil {
		return nil, err
	}
	return &Worker{ID: did, config: config}, nil
}

// Started returns true if the Worker has already been started by calling Start().
func (w *Worker) Started() bool {
	return w.started
}

// Start begins processing the queue and stops only when the nsq channels
// signal they are stopping, or the program receives a SIGINT or SIGTERM.
func (w *Worker) Start() error {
	w.started = true

	// Connect consumers to the local nsqlookupd
	var lookupds []string
	lookupd := os.Getenv("INC_NSQLOOKUPD")
	if lookupd == "" {
		lookupds = []string{"nsqlookupd.local:4161"}
	} else {
		lookupds = strings.Split(lookupd, ",")
	}
	for _, consumer := range w.consumers {
		for _, lookup := range lookupds {
			err := consumer.ConnectToNSQLookupd(lookup)
			if err != nil {
				return err
			}
		}
	}

	// Monitor each channel for stopping
	for i, consumer := range w.consumers {
		go func(i int, consumer *nsq.Consumer) {
			select {
			case <-consumer.StopChan:
				// Stop everyone else - but only once
				w.stopOnce.Do(func() { w.stopConsumers(i) })
				w.stopped[i] = true
				log.Println("Stopped channel", i)
				if w.allStopped() {
					w.stopChan <- true
				}
			}
		}(i, consumer)
	}

	// Create a system signal channel to notify us when the system bumps us
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for the command channels to stop or for the system to receive a signal
	for {
		select {
		case <-w.stopChan:
			log.Println("Stopping worker")
			return nil
		case <-sigChan:
			w.stopOnce.Do(func() { w.stopConsumers(-1) })
		}
	}
}

// Consumer creates a new nsq Consumer for the worker. The topic is used
// to subscribe to the topic name which is:
//
// worker.ID + '.' + topic
//
// and the channel as provided.
func (w *Worker) Consumer(topic, channel string, handler WorkerFunc) error {
	if w.started {
		return fmt.Errorf("Could not create consumer, worker is already started")
	}

	consumer, err := nsq.NewConsumer(w.ID.ID+"."+topic, channel, w.config)
	if err != nil {
		//log.Println("Error creating NSQ consumer", err)
		return err
	}

	consumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		log.Println("Msg p", msg.Body)
		cmd := string(msg.Body)
		log.Println("Got program", cmd)
		return handler(cmd)
	}))

	w.consumers = append(w.consumers, consumer)

	return nil
}

// allStopped checks if all the channels are stopped.
func (w *Worker) allStopped() bool {
	for _, stopped := range w.stopped {
		if !stopped {
			return false
		}
	}
	return true
}

// stopConsumers stops all the consumers that have not already received a
// stop signal. If a channel corresponds to the index number provided, that
// channel is skipped.
func (w *Worker) stopConsumers(index int) {
	w.stopped = make([]bool, len(w.consumers))

	for i, ch := range w.consumers {
		if index == i {
			log.Println("Skipping Stop() for channel", index)
		} else {
			ch.Stop()
		}
	}
}
