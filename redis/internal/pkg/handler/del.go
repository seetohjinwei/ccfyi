package handler

import (
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

const DelCommand = "DEL"

func Del(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{DelCommand}) {
		return "", false
	}

	if len(commands) < 2 {
		return invalidArgNum()
	}

	keys := commands[1:]

	s := store.GetSingleton()
	count := s.DeleteMany(keys)

	return messages.NewInteger(count).Serialise(), true
}
