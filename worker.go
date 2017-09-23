package lights

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// WorkerFunc handles incoming string messages returning an error if the
// message could not be handled.
type WorkerFunc func(message string) error

// Worker takes care of the common worker tasks that all message
// driven agents must carry out. The worker takes care of bootstrapping
// the system.
type Worker struct {
	agent  string
	queues map[string]map[string]chan<- (QMessage) // Message queues for each agent/route combination
}

// NewWorker creates a new worker ready for configuration. Call Start() on
// the worker to begin processing messages. Returns an error if there was a
// problem creating the worker. If an ID is provided the worker will use it,
// otherwise an ID will automatically be generated using lights.NewID().
func NewWorker(agent string) (*Worker, error) {
	w := &Worker{agent: agent, queues: make(map[string]map[string]chan<- (QMessage))}
	// TODO pull cached messages from disk
	port := AgentPort(agent)
	if len(port) == 0 {
		return nil, errors.New("Agent " + agent + " not supported")
	}
	return w, nil
}

// Start begins processing commands blocking the thread.
func (w *Worker) Start() error {
	// Start up sending go routine
	host := ":" + AgentPort(w.agent)
	log.Println("Listening for HTTP API", host)
	return http.ListenAndServe(host, nil)
}

// Consumer creates a new API /command Consumer for the worker.
func (w *Worker) Consumer(handler WorkerFunc) error {
	return w.Handler("/command", handler)
}

// Handler registers a new API route handler for the worker.
func (w *Worker) Handler(route string, handler WorkerFunc) error {
	http.HandleFunc(route, func(resp http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		log.Println("<-", string(body))
		if w.errorFree(err, resp) {
			err = handler(string(body))
			if w.errorFree(err, resp) {
				io.WriteString(resp, "OK")
			}
		}
	})
	return nil
}

// Transmit reliably sends a message to another agent using HTTP.
// Set consolidate to true if messages sent to the same agent and route should
// only transmit the last message when an agent is offline. If consolidate is
// false, messages queued for later delivery will all be delivered when the
// agent is reachable again.
func (w *Worker) Transmit(agent, route, message string, consolidate bool) {
	msg := QMessage{agent: agent, route: route, message: message}
	routes, ok := w.queues[msg.agent]
	if !ok {
		routes = map[string]chan<- (QMessage){}
		w.queues[msg.agent] = routes
	}
	queue, ok := routes[msg.route]
	if !ok {
		// Set up the QWorker
		_, queue = NewQWorker(msg)
		routes[msg.route] = queue
	}
	queue <- msg
}

// Send transmits a message to an agent.
func (w *Worker) Send(agent, message string) {
	w.Transmit(agent, "/command", message, false)
}

// Status transmits a status update to an agent.
func (w *Worker) Status(agent, message string) {
	w.Transmit(agent, "/status", message, true)
}

// errorFree will respond correctly to clients when an error occurs.
// Returns true if there was an error for easy handling.
func (w *Worker) errorFree(err error, resp http.ResponseWriter) bool {
	if err == nil {
		return true
	}
	resp.WriteHeader(http.StatusInternalServerError)
	io.WriteString(resp, err.Error())
	return false
}
