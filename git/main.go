package main

import (
	"github.com/seetohjinwei/ccfyi/git/internal/pkg/commands"
	"github.com/seetohjinwei/ccfyi/git/internal/pkg/logging"
)

func main() {
	logging.Init()

	command := commands.Parse()
	command.Exec()
}
