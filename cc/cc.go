package cc

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// CCommit uses the CLI to to create a conventional commit
func CCommit() {
	cli := NewCLI(os.Stdout, os.Stdin)
	cli.writeTypesPrompt()
	cli.readType()
}

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
func (c *CLI) writeTypesPrompt() {
	prompt := getTypesPrompt()
	fmt.Fprint(c.Out, prompt)
}

// getTypesPrompt combines all conventional commit types options into one prompt
func getTypesPrompt() string {
	typesString := ""
	for _, v := range types {
		typesString += v
	}
	typesString += "\nEnter a number between 0 and 10: "
	return typesString
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
