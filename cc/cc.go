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
	cli.readScope()
	cli.readSubject()
}

// CLI defines the cli for this package.
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

// writeTypesPrompt writes conventional commit type options
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

// readType will try to set a conventional commit type from the number input from user.
// Will retry after an invalid input for three times before exiting the program.
func (c *CLI) readType() {

	fails := 0
	var cctype string
	for {
		input := c.readLine()
		val, ok := CCTypeMap[input]

		if ok {
			cctype = val
			break
		} else {
			if fails > 1 {
				break
			}
			fails++
			fmt.Fprint(c.Out, "Enter a valid number between 0 and 10: ")
		}
	}

	c.cc.typ = cctype
}

// readScope takes input and sets a scope for the conventional commit
func (c *CLI) readScope() {
	fmt.Fprint(c.Out, "Enter a scope: ")

	input := c.readLine()
	if input != "" {
		input = "(" + input + ")"
	}

	c.cc.scope = input
}

// readSubject will try to set a conventional commit subject from the user input.
// Will retry after an invalid input for three times before exiting the program.
func (c *CLI) readSubject() {
	fmt.Fprint(c.Out, "Enter a subject: ")

	fails := 0
	var subject string
	for {
		input := c.readLine()

		if input != "" {
			subject = input
			break
		} else {
			if fails > 1 {
				break
			}
			fails++
			fmt.Fprint(c.Out, "Enter a subject: ")
		}
	}

	c.cc.subject = subject

}

// readLine reads a line from the CLI's input
func (c *CLI) readLine() string {
	c.In.Scan()
	return c.In.Text()
}
