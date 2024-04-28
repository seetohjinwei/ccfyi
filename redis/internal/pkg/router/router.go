package router

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

var EmptyBodyErr string = messages.GetErrorString("request body cannot be empty")
var BodyParsingErr string = messages.GetErrorString("request body could not be parsed")

type Route func(commands []string) (string, bool)

type Router struct {
	handlers []Route
}

func New(routes []Route) *Router {
	router := &Router{routes}
	return router
}

func NewDefault() *Router {
	routes := []Route{
		ping,
	}

	return New(routes)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO: Go router somehow returns 400 Bad Request before this is even called...
	// TODO: need to change to reading the raw TCP payload!!

	log.Printf("got request %v", req) // TODO: remove

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("err %q from parsing body: %q", BodyParsingErr, body)
		fmt.Fprint(w, BodyParsingErr)
		return
	} else if len(body) == 0 {
		log.Printf("err %q from parsing body: %q", EmptyBodyErr, body)
		fmt.Fprint(w, EmptyBodyErr)
		return
	}

	res := r.handle(string(body))
	fmt.Fprint(w, res)
}

func (r *Router) handle(request string) string {
	command, err := messages.Deserialise(request)
	if err != nil {
		log.Printf("err from parsing request %q: %v", request, err)
		return messages.GetError(err)
	}

	commands, err := r.getCommands(command)
	if err != nil {
		log.Printf("err from parsing request %q: %v", request, err)
		return messages.GetError(err)
	}

	ret, ok := r.route(commands)
	if !ok {
		msg := "did not match any route"
		log.Printf("err from matching command %q: %v", commands, msg)
		return messages.GetErrorString(msg)
	}

	return ret
}

func (r *Router) getCommands(request messages.Message) ([]string, error) {
	array, ok := request.(*messages.Array)
	if !ok {
		return nil, errors.New("request must be an array")
	}

	commands, err := array.GetCommands()
	if err != nil {
		return nil, err
	}

	return commands, nil
}

func (r *Router) AddRoute(route Route) {
	r.handlers = append(r.handlers, route)
}

func (r *Router) route(commands []string) (string, bool) {
	for _, route := range r.handlers {
		resp, ok := route(commands)
		if ok {
			log.Printf("received %v, response %q", commands, resp)
			return resp, true
		}
	}

	return "", false
}
