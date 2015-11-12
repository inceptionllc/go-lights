package lights

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// WorkerFunc handles incoming string messages returning an error if the
// message could not be handled.
type WorkerFunc func(message string) error

// Worker takes care of the common nsq worker tasks that all message
// driven agents must carry out. The worker takes care of bootstrapping
// the system, and automatically configures nsq according to the current
// environment.
type Worker struct {
	agent string
}

// NewWorker creates a new worker ready for configuration. Call Start() on
// the worker to begin processing messages. Returns an error if there was a
// problem creating the worker. If an ID is provided the worker will use it,
// otherwise an ID will automatically be generated using lights.NewID().
func NewWorker(agent string) (*Worker, error) {
	w := &Worker{agent: agent}
	port := w.agentPort(agent)
	if len(port) == 0 {
		return nil, errors.New("Agent " + agent + " not supported")
	}
	return w, nil
}

// Start begins processing commands blocking the thread.
func (w *Worker) Start() error {
	host := ":" + w.agentPort(w.agent)
	log.Println("Listening for HTTP API", host)
	return http.ListenAndServe(host, nil)
}

// Consumer creates a new API command Consumer for the worker.
func (w *Worker) Consumer(handler WorkerFunc) error {
	http.HandleFunc("/command", func(resp http.ResponseWriter, r *http.Request) {
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

// Send transmits a message to an agent.
func (w *Worker) Send(agent, message string) {
	url := "http://127.0.0.1:" + w.agentPort(agent) + "/command"
	log.Println("->", agent, message)
	resp, err := http.Post(url, "text/plain", strings.NewReader(message))
	if err != nil {
		log.Println("Error sending message", agent, message, err)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body", err)
		} else {
			log.Println("  ", string(body))
		}
		resp.Body.Close()
	}
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

// agentPort looks up the correct port for an agent by name.
func (w *Worker) agentPort(agent string) string {
	switch agent {
	case "gateway":
		return "8001"
	case "controller":
		return "8002"
	case "gatekeeper":
		return "8003"
	case "scheduler":
		return "8004"
	case "updater":
		return "8005"
	default:
		// Default is the gateway agent
		return ""
	}
}
