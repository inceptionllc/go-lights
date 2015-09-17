package lights

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bitly/go-nsq"
)

// Worker takes care of the common nsq worker tasks that all message
// driven agents must carry out. The worker takes care of bootstrapping
// the system, and automatically configures nsq according to the current
// environment.
type Worker struct {
	ID   string // The ID for the hardware the worker is running on
	Mode string // The mode for the current worker (DEV | "")

	config    *nsq.Config     // The nsq configuration for the worker
	consumers []*nsq.Consumer // Consumers that have been created
	started   bool            // Flag true if the worker has been started
	stopped   []bool          // Flags for consumers that have been stopped
	stopOnce  sync.Once       // We only want to stop the consumers once!
	stopChan  chan bool       // Channel receives true when the worker should stop
}

// NewWorker creates a new worker ready for configuration. Call Start() on
// the worker to begin processing messages. Returns an error if there was a
// problem creating the worker.
func NewWorker(id string) (*Worker, error) {
	// Allow setting agent mode via environmental variable
	mode := os.Getenv("MODE")

	config := nsq.NewConfig()

	return &Worker{ID: id, Mode: mode, config: config}, nil
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
	var lookupd string
	if w.Mode == "DEV" {
		lookupd = "localhost:4161"
	} else {
		lookupd = "nsqlookupd.local:4161"
	}
	for _, consumer := range w.consumers {
		err := consumer.ConnectToNSQLookupd(lookupd)
		if err != nil {
			return err
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

// Consumer creates a new nsq Consumer for the worker. The name is used
// to subscribe to the topic name which is:
//
// worker.ID + '.' + name
//
// and the channel `ID`
func (w *Worker) Consumer(name string) (*nsq.Consumer, error) {
	if w.started {
		return nil, fmt.Errorf("Could not create consumer, worker is already started")
	}

	consumer, err := nsq.NewConsumer(w.ID+"."+name, w.ID, w.config)
	if err != nil {
		//log.Println("Error creating NSQ consumer", err)
		return nil, err
	}

	w.consumers = append(w.consumers, consumer)

	return consumer, nil
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
