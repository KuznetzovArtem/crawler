package http

import (
	"crawler/task"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net"
	"net/http"
)

// Server will perform operations over http.
type Server interface {
	// Open will setup a new listener on the given port.
	Open() error

	// Close will close the connection if it's open.
	Close()

	// Handler returns a http handler with all routes in place.
	Handler() http.Handler
}

// Server represents an HTTP server.
type server struct {
	listener    net.Listener
	taskService TaskSender
	encoder     *encoder
	addr        string
}

// Open will setup a tcp listener and serve the http requests.
func (s *server) Open() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	// Save listener so we can decide to close it later.
	s.listener = ln

	// Start HTTP server.
	server := http.Server{
		Handler: s.Handler(),
	}

	return server.Serve(s.listener)
}

// Close will close the socket if it's open.
func (s *server) Close() {
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
}

// route will setup all routes and return the http handler.
func (s *server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(
			middleware.Recoverer,
		)
		r.Route(
			task.SendCollectRequest,
			newTaskHandler(s.encoder, s.taskService).Routes,
		)
	})
	return r
}

// NewServer returns a new instance of Server.
func NewServer(
	addr string,
	taskService TaskSender,
) Server {
	return &server{
		addr:        addr,
		taskService: taskService,
	}
}
