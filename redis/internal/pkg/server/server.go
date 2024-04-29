package server

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/router"
)

type Server struct {
	ctx  context.Context
	port string
	wg   sync.WaitGroup
	r    *router.Router
	l    net.Listener
}

func New(port string) (*Server, error) {
	port = ":" + port
	l, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	s := &Server{
		ctx:  ctx,
		port: ":" + port,
		wg:   sync.WaitGroup{},
		r:    router.NewDefault(),
		l:    l,
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

func (s *Server) Stop() bool {
	done := make(chan bool, 2)
	go func() {
		// TODO: increase timeout
		<-time.After(1 * time.Second)
		done <- false
	}()
	go func() {
		s.wg.Wait()
		done <- true
	}()

	isGraceful := <-done
	if isGraceful {
		log.Printf("server gracefully stopped")
	} else {
		log.Printf("server abruptly stopped because of timeout")
	}

	s.l.Close()

	return isGraceful
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
					log.Printf("err while reading from conn")
					return
				}
			}

			buf = append(buf, tmp[:n]...)
			if len(buf) == 0 && isEof {
				log.Printf("connection stopped because len(buf) == 0 && eof")
				return
			}
			req := string(buf)
			reply, ok := s.r.Handle(req)
			if ok || isEof {
				// return the reply
				reply := []byte(reply)
				for len(reply) > 0 {
					n, err := conn.Write(reply)
					if err != nil {
						log.Printf("err while writing from conn")
						return
					}
					reply = reply[n:]
				}

				log.Printf("debug raw request: %q, raw reply: %q", req, reply)
				break
			}
		}
	}
}
