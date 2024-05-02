package main

import (
	"github.com/rs/zerolog/log"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/logging"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/server"
)

func init() {
	logging.Init()
}

func main() {
	// TODO:

	// test my Serialise and Deserialise more? (esp error paths!)
	// create some integration tests with the go redis client library
	// zerolog display the file + line

	// find some way to propagate errors? (not sure what i meant by this)

	// protocol description: https://redis.io/docs/latest/develop/reference/protocol-spec/#resp-protocol-description

	router, err := server.New("localhost:6379") // TODO: take port as flag
	if err != nil {
		log.Fatal().Err(err).Msg("server init")
	}
	err = router.Serve()
	if err != nil {
		log.Fatal().Err(err).Msg("server serve error")
	}
}
