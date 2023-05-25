package cc

import (
	"os"

	Z "github.com/rwxrob/bonzai/z"
)

// exported leaf
//
// CCommit uses the CLI to to create a conventional commit
var Cmd = &Z.Cmd{
	Name: `cc`,
	Call: func(caller *Z.Cmd, none ...string) error {
		cli := NewCLI(os.Stdout, os.Stdin, &CCExecutor{})
		cli.writeTypesPrompt()
		cli.readType()
		cli.readScope()
		cli.readSubject()
		cli.readBodyAndFooter()
		cli.buildMessage()
		cli.makeCommit()
		return nil
	},
}
