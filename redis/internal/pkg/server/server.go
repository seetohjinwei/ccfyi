package server

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/router"
)

// Server is a TCP server. To construct one, use `Server::New`.
// When a sigint is captured, if there are any ongoing events (e.g. connections), the server will wait for up to X seconds before forcefully shutting down; if there are no events, it will gracefully shutdown.
type Server struct {
	ctx      context.Context
	port     string
	wg       sync.WaitGroup
	stopOnce sync.Once
	r        *router.Router
	l        net.Listener
}

// New constructs a new Server with the specified port.
func New(port string) (*Server, error) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	s := &Server{
		ctx:      ctx,
		port:     ":" + port,
		wg:       sync.WaitGroup{},
		stopOnce: sync.Once{},
		r:        router.NewDefault(),
		l:        l,
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		cancelFunc()
		s.Stop()
	}()

	return s, nil
}

// Serve lets the server start accepting connections.
func (s *Server) Serve() error {
	if s == nil {
		return errors.New("tried to call *Server::Serve() on nil")
	}

	for {
		conn, err := s.l.Accept()
		if err != nil {
			select {
			case <-s.ctx.Done():
				return nil
			default:
				return err
			}
		}

		go s.handleConnection(conn)
	}
}

// Stops the server.
func (s *Server) Stop() {
	s.stopOnce.Do(func() {
		done := make(chan bool, 2)
		go func() {
			// TODO: increase timeout
			<-time.After(1 * time.Second)
			done <- false
		}()
		go func() {
			s.wg.Wait()
			// TODO: save store to disk?
			done <- true
		}()

		isGraceful := <-done
		if isGraceful {
			log.Info().Msg("server gracefully stopped")
		} else {
			log.Info().Msg("server abruptly stopped because of timeout")
		}

		s.l.Close()
	})
}

func (s *Server) handleConnection(conn net.Conn) {
	s.wg.Add(1)
	defer s.wg.Done()
	defer conn.Close()

	for {
		// connection loop

		buf := make([]byte, 0, 4096)
		for {
			// fetches a single command
			tmp := make([]byte, 2048)
			isEof := false

			n, err := conn.Read(tmp)
			if err != nil {
				if err == io.EOF {
					isEof = true
				} else {
					log.Err(err).Msg("reading from conn")
					return
				}
			}

			buf = append(buf, tmp[:n]...)
			if len(buf) == 0 && isEof {
				log.Debug().Msg("connection stopped because len(buf) == 0 && eof")
				return
			}
			req := string(buf)
			reply, ok := s.r.Handle(req)
			if ok || isEof {
				// return the reply
				reply := []byte(reply)
				log.Debug().Str("req", req).Bytes("reply", reply).Msg("raw")
				for len(reply) > 0 {
					n, err := conn.Write(reply)
					if err != nil {
						log.Err(err).Msg("writing to conn")
						return
					}
					reply = reply[n:]
				}

				break
			}
		}
	}
}
