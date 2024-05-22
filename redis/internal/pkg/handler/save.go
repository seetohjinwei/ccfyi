package handler

import (
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

const SaveCommand = "SAVE"

func Save(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{SaveCommand}) {
		return "", false
	}

	if len(commands) != 1 {
		return invalidArgNum()
	}

	s := store.GetSingleton()
	err := s.SaveToDisk()
	if err != nil {
		return messages.GetError(err), true
	}

	return messages.NewSimpleString("OK").Serialise(), true
}
