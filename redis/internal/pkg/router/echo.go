package router

import (
  "github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

func echo(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{"ECHO"}) {
		return "", false
	}

	if len(commands) != 2 {
		return invalidArgNum()
	}

	return messages.NewBulkString(commands[1]).Serialise(), true
}
