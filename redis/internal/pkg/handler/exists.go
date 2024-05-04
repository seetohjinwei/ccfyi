package handler

import (
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

func exists(s *store.Store, cache map[string]bool, key string) bool {
	if has, ok := cache[key]; ok {
		return has
	}

	_, has := s.Get(key)
	cache[key] = has

	return has
}

const ExistsCommand = "EXISTS"

func Exists(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{ExistsCommand}) {
		return "", false
	}

	if len(commands) < 2 {
		return invalidArgNum()
	}

	cache := make(map[string]bool)
	s := store.GetSingleton()
	count := int64(0)
	for i := 1; i < len(commands); i++ {
		key := commands[i]
		if exists(s, cache, key) {
			count++
		}
	}

	return messages.NewInteger(count).Serialise(), true
}
