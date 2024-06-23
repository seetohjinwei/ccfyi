package commands

import (
	"flag"
	"os"
)

type Command interface {
	Exec()
}

type Flags struct {
	Help bool
}

func parseFlags() Flags {
	ret := Flags{}

	flag.BoolVar(&ret.Help, "help", false, "prints usage info")

	flag.Parse()

	return ret
}

func Parse() (Command, error) {
	args := os.Args

	flags := parseFlags()

	if flags.Help || len(args) < 2 {
		return &Help{
			Flags: flags,
		}, nil
	}

	return &Init{
		Flags: flags,
	}, nil
}
