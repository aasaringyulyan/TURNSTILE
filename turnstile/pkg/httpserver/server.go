package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout  = 15 * time.Second
	_defaultWriteTimeout = 15 * time.Second
	//_defaultAddr            = ":80"
	//_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server *http.Server
	notify chan error
	//shutdownTimeout time.Duration
}

func New(handler http.Handler, port string) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		Addr:         ":" + port,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
	}

	s := &Server{
		server: httpServer,
		notify: make(chan error, 1),
		//shutdownTimeout: _defaultShutdownTimeout,
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	//ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	//defer cancel()

	//return s.server.Shutdown(ctx)
	return s.server.Shutdown(context.Background())
}
