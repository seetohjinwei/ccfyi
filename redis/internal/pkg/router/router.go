package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Router struct {
	srv *http.Server
}

func handler(w http.ResponseWriter, req *http.Request) {
	// TODO: take the req and pass to the actual handler
	// then write the response from the handler
	// the actual handler is some Handler(command string) (response string)
	fmt.Fprintf(w, "안녕\n")
}

func New(port string) *Router {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: http.HandlerFunc(handler),
	}

	return &Router{
		srv: srv,
	}
}

func (r *Router) Serve() error {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		err := r.srv.Shutdown(context.Background())
		if err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	err := r.srv.ListenAndServe()
	return err
}
