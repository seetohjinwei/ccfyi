package commands

import (
	"flag"
	"fmt"
	"os"
)

type Help struct {
	Flags
}

func (c *Help) Exec() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}
