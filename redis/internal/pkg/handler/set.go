package handler

import (
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

func Set(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{"SET"}) {
		return "", false
	}

	if len(commands) < 3 || len(commands) > 6 {
		return invalidArgNum()
	}

	s := store.GetSingleton()

	key := commands[1]
	value := commands[2]

	err := s.Set(key, store.NewString(value))
	if err != nil {
		return messages.GetError(err), true
	}

	return messages.NewSimpleString("OK").Serialise(), true
}
