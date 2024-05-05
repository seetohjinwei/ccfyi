package handler

import (
	"github.com/rs/zerolog/log"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

const LLenCommand = "LLEN"

func LLen(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{LLenCommand}) {
		return "", false
	}

	if len(commands) != 2 {
		return invalidArgNum()
	}

	s := store.GetSingleton()
	key := commands[1]
	item, ok := s.Get(key)
	if !ok {
		return messages.NewInteger(0).Serialise(), true
	}

	ret, ok := item.LLen()
	if !ok {
		msg := "WRONGTYPE Operation against a key holding the wrong kind of value"
		log.Error().Any("item", item).Msg(msg)
		return messages.GetErrorString(msg), true
	}

	return messages.NewInteger(ret).Serialise(), true
}
