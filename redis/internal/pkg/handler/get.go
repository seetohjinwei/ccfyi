package handler

import (
	"github.com/rs/zerolog/log"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

const GetCommand = "GET"

func Get(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{GetCommand}) {
		return "", false
	}

	if len(commands) != 2 {
		return invalidArgNum()
	}

	s := store.GetSingleton()

	key := commands[1]

	item, ok := s.Get(key)
	if !ok {
		return messages.NewNullBulkString().Serialise(), true
	}

	val, ok := item.Do("get", nil)
	if !ok {
		log.Error().Msg("get does not exist!")
		return internalServerError()
	}

	return messages.NewBulkString(val).Serialise(), true
}
