package main

import (
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	"github.com/yemaney/z/aws"
	"github.com/yemaney/z/cc"
	"github.com/yemaney/z/compcmd"
	"github.com/yemaney/z/ssh"
)

func main() {
	Cmd.Run()
}

var Cmd = &Z.Cmd{
	Name:    `z`,
	Summary: `yemane's bonzai command tree`,
	Version: `v1.4.0`,
	Source:  `git@github.com:yemaney/z.git`,
	Issues:  `github.com/yemaney/z/issues`,
	Comp:    compcmd.New(),
	Commands: []*Z.Cmd{
		help.Cmd, cc.Cmd, ssh.Cmd, aws.Cmd, conf.Cmd
	},

	Description: `
		Hi, I'm yemane and this {{cmd .Name }} is my Bonzai™ tree. I am
		slowly replacing all my shell scripts and other Go utilities with
		Bonzai branches that I graft into this {{cmd .Name}} command.
		`,
}
