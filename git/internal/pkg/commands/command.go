package commands

import (
	"flag"
	"fmt"
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

func Parse() Command {
	flags := parseFlags()
	args := flag.Args()

	if flags.Help || len(args) < 1 {
		return &Help{
			Flags: flags,
		}
	}

	commandName := args[0]

	switch commandName {
	case "init":
		return &Init{
			Flags: flags,
			Args:  args[1:],
		}
	}

	fmt.Fprintf(os.Stderr, "'%v' is not a git command. See `git --help`.\n", commandName)
	os.Exit(1)
	panic("unreachable because exit")
}
