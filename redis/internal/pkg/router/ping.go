package router

import (
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

func ping(commands []string) (string, bool) {
	if len(commands) == 0 || commands[0] != "PING" {
		return "", false
	}

	if len(commands) == 1 {
		return messages.NewSimpleString("PONG").Serialise(), true
	}

	return messages.NewBulkString(commands[1]).Serialise(), true
}
