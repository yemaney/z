package cc

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// CCommit uses the CLI to to create a conventional commit
func CCommit() {
	cli := NewCLI(os.Stdout, os.Stdin, &CCExecutor{})
	cli.writeTypesPrompt()
	cli.readType()
	cli.readScope()
	cli.readSubject()
	cli.readBodyAndFooter()
	cli.buildMessage()
	cli.makeCommit()
}

// CLI defines the cli for this package.
type CLI struct {
	Out io.Writer
	In  *bufio.Scanner
	cc  CC
	ce  CmdExecutor
}

// NewCLI creates a CLI for creating conventional commits
func NewCLI(out io.Writer, in io.Reader, ce CmdExecutor) *CLI {
	return &CLI{
		Out: out,
		In:  bufio.NewScanner(in),
		ce:  ce,
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

// readBodyAndFooter will try to set a body and footer for the conventional commit.
func (c *CLI) readBodyAndFooter() {
	fmt.Fprint(c.Out, "Enter a body: ")
	inputBody := c.readLine()
	c.cc.body = inputBody

	fmt.Fprint(c.Out, "Enter a footer: ")
	inputFooter := c.readLine()
	c.cc.footer = inputFooter
}

// readLine reads a line from the CLI's input
func (c *CLI) readLine() string {
	c.In.Scan()
	return c.In.Text()
}

// buildMessage uses all the CC fields to create a conventional commit message
func (c *CLI) buildMessage() {
	message := ""
	message += c.cc.typ
	message += c.cc.scope
	message += ": "
	message += c.cc.subject

	if c.cc.body != "" {
		message += "\n\n"
		message += c.cc.body
	}

	if c.cc.footer != "" {
		message += "\n\n"
		message += c.cc.footer
	}

	c.cc.message = message
}

// makeCommit runs the CmdExecutor *exec.Cmd to make a conventional commit with git
func (c *CLI) makeCommit() {
	cmd := c.ce.build(c.cc.message)
	cmd.Stdout = c.Out
	cmd.Run()
}
