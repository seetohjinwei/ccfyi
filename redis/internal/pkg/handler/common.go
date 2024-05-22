package handler

import (
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/items"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

func invalidArgNum() (string, bool) {
	msg := "ERR wrong number of arguments for command"
	return messages.NewError(msg).Serialise(), true
}

func internalServerError() (string, bool) {
	msg := "ERR internal server error (check server logs)"
	return messages.NewError(msg).Serialise(), true
}

func wrongTypeError(item items.Item) (string, bool) {
	msg := "WRONGTYPE Operation against a key holding the wrong kind of value"
	log.Error().Any("item", item).Msg(msg)
	return messages.GetErrorString(msg), true
}

func commandsStartWith(commands []string, should []string) bool {
	if len(commands) < len(should) {
		return false
	}

	for i, s := range should {
		if !strings.EqualFold(s, commands[i]) {
			return false
		}
	}

	return true
}
