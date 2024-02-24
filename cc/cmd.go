package cc

import (
	"os"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"github.com/yemaney/z/compcmd"
)

// exported bonzai leaf that uses uses cc.CLI to to create a conventional commit
var Cmd = &Z.Cmd{
	Name:     `cc`,
	Summary:  `git commit in the style of conventional commits`,
	Params:   []string{"signed", "breaking"},
	Usage:    `[signed] [breaking]`,
	Comp:     compcmd.New(),
	Commands: []*Z.Cmd{help.Cmd},
	Description: `
		The {{aka}} command provides the ability to make git commits in the style of conventional commits. (https://www.conventionalcommits.org)


		Options: 

		[OPTIONAL]

		signed		:	the commit will be signed <https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits>

		breaking	:	the commit message will start with: type! or type(scope)!

		`,
	Call: func(caller *Z.Cmd, args ...string) error {
		cli := NewCLI(os.Stdout, os.Stdin, &CCExecutor{})
		cli.parseParams(args)
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
