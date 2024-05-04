package handler

import (
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

const EchoCommand = "ECHO"

func Echo(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{EchoCommand}) {
		return "", false
	}

	if len(commands) != 2 {
		return invalidArgNum()
	}

	return messages.NewBulkString(commands[1]).Serialise(), true
}
