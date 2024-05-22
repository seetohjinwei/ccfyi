package handler

import (
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

const LRangeCommand = "LRANGE"

func LRange(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{LRangeCommand}) {
		return "", false
	}

	if len(commands) != 4 {
		return invalidArgNum()
	}

	s := store.GetSingleton()
	key := commands[1]

	start, err := strconv.Atoi(commands[2])
	if err != nil {
		msg := "value is not an integer or out of range"
		log.Error().Any("start", commands[2]).Msg(msg)
		return messages.GetErrorString(msg), true
	}
	stop, err := strconv.Atoi(commands[3])
	if err != nil {
		msg := "value is not an integer or out of range"
		log.Error().Any("stop", commands[3]).Msg(msg)
		return messages.GetErrorString(msg), true
	}

	item, ok := s.Get(key)
	if !ok {
		return messages.NewArray([]messages.Message{}).Serialise(), true
	}

	ret, ok := item.LRange(start, stop)
	if !ok {
		return wrongTypeError(item)
	}

	return messages.NewArrayBulkString(ret).Serialise(), true
}
