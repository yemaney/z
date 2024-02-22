package ssh

import (
	"os"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	compcmd "github.com/yemaney/z/compcmd"
)

// exported bonzai leaf that uses uses ssh.CLI to edit an ssh config
var Cmd = &Z.Cmd{
	Name:     `ssh`,
	Summary:  `edit your ssh config file`,
	Params:   []string{"add", "delete"},
	Comp:     compcmd.New(),
	Commands: []*Z.Cmd{help.Cmd, addCmd, delCmd},
}

var addCmd = &Z.Cmd{
	Name:    `add`,
	Summary: `add a section to your ssh config file`,
	Usage:   `-n example -h example.com -u root`,
	Params:  []string{"name", "host", "user", "identityFile", "port"},
	Comp:    compcmd.New(),
	Description: `
		The {{aka}} command provides the ability to update your ssh config file
		through the command line.

		Options:

		-n	 	:	name to identify this ssh section or hostname that should be used to establish the connection

		-host	:	hostname that should be used to establish the connection 

		-u 		:	username to be used for the connection

		-i 		:	private key that the client should use for authentication when connecting to the ssh server

		-p 		:	port that the remote SSH daemon is running on. only necessary if the remote SSH instance is not running on the default port 22
		`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {

		c := NewCLI(os.Stdout)

		err := c.createSection(args)
		if err != nil {
			os.Exit(1)
		}

		err = c.loadConfig()
		if err != nil {
			os.Exit(1)
		}

		c.parseConfig()
		c.updateSections()
		c.createConfig()

		err = c.backupAndSave()
		if err != nil {
			os.Exit(1)
		}
		return nil
	},
}

var delCmd = &Z.Cmd{
	Name:     `delete`,
	Summary:  `delete a section from your ssh config file`,
	Usage:    ` section-name-1 section-name-2 ...`,
	Comp:     compcmd.New(),
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {

		c := NewCLI(os.Stdout)

		err := c.loadConfig()
		if err != nil {
			os.Exit(1)
		}

		c.parseConfig()
		c.deleteSections(args)
		c.createConfig()

		err = c.backupAndSave()
		if err != nil {
			os.Exit(1)
		}
		return nil
	},
}
