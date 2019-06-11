package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Server struct
// Inspired from: https://meshstudio.io/blog/2017-11-06-serving-html-with-golang/
type Server struct {
	srv *http.Server
	wg  sync.WaitGroup
}

// New creates new Server instance
func New(s *http.Server) *Server {

	return &Server{
		srv: s,
	}
}

// Start server
func (s *Server) Start() {
	// Setup Context
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Add to the WaitGroup for the listener goroutine
	s.wg.Add(1)

	// Start the listener
	go func() {
		fmt.Println("Starting server...")
		s.srv.ListenAndServe()
		s.wg.Done()
	}()
}

// Stop server
func (s *Server) Stop() {
	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		if err = s.srv.Close(); err != nil {
			fmt.Printf(fmt.Sprintf("Service stopping: Error %v\n", err))
			return
		}
	}

	s.wg.Wait()
	fmt.Println("Server has stopped.")
}
