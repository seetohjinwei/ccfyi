package handler

import (
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

func Ping(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{"PING"}) {
		return "", false
	}

	if len(commands) == 1 {
		return messages.NewSimpleString("PONG").Serialise(), true
	}

	if len(commands) != 2 {
		return invalidArgNum()
	}

	return messages.NewBulkString(commands[1]).Serialise(), true
}
