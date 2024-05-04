package handler

import (
	"errors"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/pkg/delay"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

type setFlag int

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
	NX        bool
	XX        bool
	shouldGet bool
	expiry    time.Time
}

func parseSetArguments(commands []string) (setArgs, error) {
	args := setArgs{
		NX:        false,
		XX:        false,
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
			args.NX = true
		case "XX":
			args.XX = true
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

const SetCommand = "SET"

func Set(commands []string) (string, bool) {
	if len(commands) == 0 || !commandsStartWith(commands, []string{SetCommand}) {
		return "", false
	}

	if len(commands) < 3 || len(commands) > 7 {
		return invalidArgNum()
	}

	s := store.GetSingleton()

	key := commands[1]
	value := commands[2]
	args, err := parseSetArguments(commands)
	if err != nil {
		return messages.GetError(err), true
	}

	var oldKey string
	var exists bool
	if args.NX || args.XX || args.shouldGet {
		// only get the key if required
		var item store.Item
		var ok bool
		item, exists = s.Get(key)
		if exists {
			oldKey, ok = item.Do("get", nil)
			if !ok {
				err := errors.New("item was not a string")
				log.Error().Err(err).Msg("getDuration")
				return messages.GetError(err), true
			}
		}
	}

	if (args.NX && exists) || (args.XX && !exists) {
		// key was NOT set
		return messages.NewNullBulkString().Serialise(), true
	}

	if args.expiry.IsZero() {
		err = s.Set(key, store.NewString(value))
	} else {
		err = s.SetWithDelay(key, store.NewString(value), delay.NewDelay(args.expiry))
	}
	if err != nil {
		return messages.GetError(err), true
	}

	if args.shouldGet {
		// key was set (with GET)
		if !exists {
			return messages.NewNullBulkString().Serialise(), true
		} else {
			return messages.NewBulkString(oldKey).Serialise(), true
		}
	}
	// key was set (without GET)
	return messages.NewSimpleString("OK").Serialise(), true
}
