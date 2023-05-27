package cc

import (
	"os"

	Z "github.com/rwxrob/bonzai/z"
)

// exported bonzai leaf that uses uses cc.CLI to to create a conventional commit
var Cmd = &Z.Cmd{
	Name:    `cc`,
	Summary: `git commit in the style of conventional commits`,
	Call: func(caller *Z.Cmd, none ...string) error {
		cli := NewCLI(os.Stdout, os.Stdin, &CCExecutor{})
		cli.checkSigned()
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
