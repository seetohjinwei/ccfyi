package router

import (
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/handler"
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
		handler.Ping,
		handler.Echo,
		handler.Get,
		handler.Set,
	}

	return New(routes)
}

func (r *Router) Handle(request string) (string, bool) {
	command, err := messages.Deserialise(request)
	if err != nil {
		log.Debug().Err(err).Str("request", request).Msg("parsing request")
		return messages.GetError(err), false
	}

	commands, err := r.getCommands(command)
	if err != nil {
		log.Err(err).Str("request", request).Msg("getting commands from request")
		return messages.GetError(err), false
	}

	ret, ok := r.route(commands)
	if !ok {
		msg := "did not match any route"
		log.Error().Str("err", msg).Strs("commands", commands).Msg("getting commands from request")
		return messages.GetErrorString(msg), true
	}

	return ret, true
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
			log.Info().Strs("commands", commands).Str("resp", resp).Msg("matched route")
			return resp, true
		}
	}

	return "", false
}
