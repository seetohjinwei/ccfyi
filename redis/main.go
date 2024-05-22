package main

import (
	"github.com/rs/zerolog/log"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/logging"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/server"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
)

func main() {
	logging.Init()

	// TODO:

	// zerolog display the file + line

	// protocol description: https://redis.io/docs/latest/develop/reference/protocol-spec/#resp-protocol-description

	s := store.GetSingleton()
	if err := s.LoadFromDisk(); err != nil {
		panic(err)
	}

	router, err := server.New("localhost:6379") // TODO: take port as flag
	if err != nil {
		log.Fatal().Err(err).Msg("server init")
	}
	err = router.Serve()
	if err != nil {
		log.Fatal().Err(err).Msg("server serve error")
	}
}
