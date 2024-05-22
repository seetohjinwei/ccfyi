package handler

import (
	"github.com/rs/zerolog/log"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/items"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

const RPushCommand = "RPUSH"

func RPush(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{RPushCommand}) {
		return "", false
	}

	if len(commands) < 3 {
		return invalidArgNum()
	}

	s := store.GetSingleton()
	key := commands[1]
	item, ok := s.Get(key)
	if !ok {
		item = items.NewList()
		err := s.Set(key, item)
		if err != nil {
			return messages.GetError(err), true
		}
	}

	log.Debug().Strs("values", commands[2:]).Msg("pushing into RPush")
	ret, ok := item.RPush(commands[2:])
	if !ok {
		return wrongTypeError(item)
	}

	return messages.NewInteger(ret).Serialise(), true
}
