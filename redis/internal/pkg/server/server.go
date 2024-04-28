package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/router"
)

type Server struct {
	srv *http.Server
}

func New(port string) *Server {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router.NewDefault(),
	}

	return &Server{
		srv: srv,
	}
}

func (s *Server) Serve() error {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		err := s.srv.Shutdown(context.Background())
		if err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	err := s.srv.ListenAndServe()
	return err
}
