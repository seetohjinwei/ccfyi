package handler

import (
	"errors"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

type setFlag int

const (
	setNone setFlag = iota
	setNX
	setXX
)

func getDuration(commands []string) (time.Duration, []string, error) {
	// first value is the type (ignore), second value is the time
	if len(commands) < 2 {
		err := errors.New("not enough values")
		log.Error().Err(err).Msg("getDuration")
		return 0, commands, err
	}

	val, err := strconv.ParseInt(commands[1], 10, 64)
	if err != nil {
		return 0, commands, err
	}
	if val <= 0 {
		err := errors.New("val must be positive")
		log.Error().Err(err).Msg("getDuration")
		return 0, commands, err
	}

	return time.Duration(val), commands[1:], nil
}

type setArgs struct {
	flag      setFlag
	shouldGet bool
	expiry    time.Time
}

func parseSetArguments(commands []string) (setArgs, error) {
	args := setArgs{
		flag:      setNone,
		shouldGet: false,
		expiry:    time.Time{},
	}
	var err error

	var d time.Duration

	commands = commands[3:]

	// supports arguments being out of order
	for len(commands) > 0 {
		switch commands[0] {
		case "NX":
			args.flag = setNX
		case "XX":
			args.flag = setXX
		case "GET":
			args.shouldGet = true
		case "EX":
			d, commands, err = getDuration(commands)
			if err != nil {
				return args, err
			}
			args.expiry = time.Now().Add(d * time.Second)
		case "PX":
			d, commands, err = getDuration(commands)
			if err != nil {
				return args, err
			}
			args.expiry = time.Now().Add(d * time.Millisecond)
		case "EXAT":
			d, commands, err = getDuration(commands)
			if err != nil {
				return args, err
			}
			args.expiry = time.Unix(int64(d), 0).UTC()
		case "PXAT":
			d, commands, err = getDuration(commands)
			if err != nil {
				return args, err
			}
			args.expiry = time.UnixMilli(int64(d)).UTC()
		case "KEEPTTL":
			// no-op
		}

		commands = commands[1:]
	}

	return args, nil
}

func Set(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{"SET"}) {
		return "", false
	}

	if len(commands) < 3 || len(commands) > 7 {
		return invalidArgNum()
	}

	s := store.GetSingleton()

	key := commands[1]
	value := commands[2]

	err := s.Set(key, store.NewString(value))
	if err != nil {
		return messages.GetError(err), true
	}

	return messages.NewSimpleString("OK").Serialise(), true
}
