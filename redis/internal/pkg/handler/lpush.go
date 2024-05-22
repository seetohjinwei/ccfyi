package handler

import (
	"github.com/rs/zerolog/log"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/items"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

const LPushCommand = "LPUSH"

func LPush(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{LPushCommand}) {
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

	log.Debug().Strs("values", commands[2:]).Msg("pushing into LPush")
	ret, ok := item.LPush(commands[2:])
	if !ok {
		return wrongTypeError(item)
	}

	return messages.NewInteger(ret).Serialise(), true
}
