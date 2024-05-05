package handler

import (
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

func del(s *store.Store, key string) bool {
	// TODO: call Delete
	_, has := s.Get(key)

	return has
}

const DelCommand = "DEL"

func Del(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{DelCommand}) {
		return "", false
	}

	if len(commands) < 2 {
		return invalidArgNum()
	}

	s := store.GetSingleton()
	count := int64(0)
	for i := 1; i < len(commands); i++ {
		key := commands[i]
		if del(s, key) {
			count++
		}
	}

	return messages.NewInteger(count).Serialise(), true
}
