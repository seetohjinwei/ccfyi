package commands

import (
	"flag"
	"os"

	"github.com/seetohjinwei/ccfyi/git/internal/pkg/operations/init_op"
)

type Init struct {
	Flags
	Args []string

	quiet bool
}

func (c *Init) parseFlags() {
	f := flag.NewFlagSet("git-init", flag.ExitOnError)

	f.BoolVar(&c.quiet, "q", false, "only print error and warning messages (shorthand)")
	f.BoolVar(&c.quiet, "quiet", false, "only print error and warning messages")

	f.Parse(c.Args)

	c.Args = f.Args()
}

func (c *Init) Exec() {
	c.parseFlags()

	directory := "."
	if len(c.Args) >= 1 {
		directory = c.Args[0]
	}

	options := init_op.Options{
		Quiet: c.quiet,
	}

	err := init_op.Perform(directory, options)
	if err != nil {
		os.Exit(1)
	}
}
