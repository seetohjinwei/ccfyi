package commands

import "fmt"

type Init struct {
	Flags
	Args []string
}

func (c *Init) Exec() {
	fmt.Println(c.Args)
}
