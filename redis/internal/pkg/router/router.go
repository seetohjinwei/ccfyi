package router

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/handler"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

var EmptyBodyErr string = messages.NewError("request body cannot be empty").Serialise()
var BodyParsingErr string = messages.NewError("request body could not be parsed").Serialise()

type Router struct {
	srv *http.Server
}

func handlerFunc(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprint(w, BodyParsingErr)
		return
	} else if len(body) == 0 {
		fmt.Fprint(w, EmptyBodyErr)
		return
	}

	res := handler.Handle(string(body))
	fmt.Fprint(w, res)
}

func New(port string) *Router {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: http.HandlerFunc(handlerFunc),
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
