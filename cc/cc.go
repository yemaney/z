// Package cc provides functionality for making conventional commits.
package cc

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// CLI defines the cli for this package.
type CLI struct {
	Out io.Writer
	In  *bufio.Scanner
	cc  *CC
	ce  CmdExecutor
}

// NewCLI creates a CLI for creating conventional commits
func NewCLI(out io.Writer, in io.Reader, ce CmdExecutor) *CLI {
	return &CLI{
		Out: out,
		In:  bufio.NewScanner(in),
		cc:  &CC{},
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
	typesString += "\nEnter a number between 0 and 7: "
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
			fmt.Fprint(c.Out, "Enter a valid number between 0 and 7: ")
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

	if c.cc.breaking {
		message += "!"
	}

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

// writeConfirmationPrompt writes a message to the user asking them to confirm if they
// want to make a commit with the message that was built.
func (c *CLI) writeConfirmationPrompt() {
	start := "\n\nPotential commit message:\n\n"
	end := "\n\nCommit these changes with the message [y/N]: "
	fmt.Fprint(c.Out, start+"\033[36;1m"+c.cc.message+"\033[0m"+end)
}

// makeCommit firsts prompts the user to confirm if they want to make a commit with the message.
// If the user responds with either a "y" or "yes" it will build the  CmdExecutor *exec.Cmd
// and run it to make a conventional commit with git
func (c *CLI) makeCommit() {
	c.writeConfirmationPrompt()

	input := strings.ToLower(c.readLine())

	if input == "y" || input == "yes" {
		cmd := c.ce.build(c.cc.message, c.cc.signed)
		cmd.Stdout = c.Out
		cmd.Run()
	}

}

// parseParams loops through all the parameters passed to the command
// and updates the state of the CC accordingly.
func (c *CLI) parseParams(args []string) {
	for i := 0; i < len(args); i++ {
		param := args[i]

		switch param {
		case "signed":
			c.cc.signed = true
		case "breaking":
			c.cc.breaking = true
		}
	}
}
