package handler

import (
	"github.com/rs/zerolog/log"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/items"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

const IncrCommand = "INCR"

func Incr(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{IncrCommand}) {
		return "", false
	}

	if len(commands) != 2 {
		return invalidArgNum()
	}

	s := store.GetSingleton()
	key := commands[1]

	item, ok := s.Get(key)
	if !ok {
		item = items.NewString("0")
		err := s.Set(key, item)
		if err != nil {
			return messages.GetError(err), true
		}
	}

	ret, ok := item.Incr()
	if !ok {
		msg := "value is not an integer or out of range"
		log.Error().Any("item", item).Msg(msg)
		return messages.GetErrorString(msg), true
	}

	return messages.NewInteger(ret).Serialise(), true
}
