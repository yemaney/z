package cc

import (
	"bufio"
	"fmt"
	"io"
)

// Struct that defines the cli for this package.
type CLI struct {
	Out io.Writer
	In  *bufio.Scanner
	cc  CC
}

// NewCLI creates a CLI for creating conventional commits
func NewCLI(out io.Writer, in io.Reader) *CLI {
	return &CLI{
		Out: out,
		In:  bufio.NewScanner(in),
	}
}

// writes conventional commit type options
func (c *CLI) writeTypes() {

	for _, v := range types {
		fmt.Fprint(c.Out, v)
	}
}

func (c *CLI) readType() {
	input := c.readLine()
	c.cc.cctype = CCTypeMap[input]
}

// readLine reads a line from the CLI's input
func (c *CLI) readLine() string {
	c.In.Scan()
	return c.In.Text()
}
