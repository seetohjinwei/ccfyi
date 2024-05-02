package store

import (
	"github.com/rs/zerolog/log"
)

type String struct {
	str string
}

func NewString(str string) *String {
	return &String{str: str}
}

func (s *String) Do(command string, args []string) (string, bool) {
	switch command {
	case "get":
		if len(args) != 0 {
			log.Error().Strs("args", args).Int("args len", len(args)).Msg("store.string wrong args len")
		}
		return s.get(), true
	}

	return "", false
}

func (s *String) get() string {
	return s.str
}
