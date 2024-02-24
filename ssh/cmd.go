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
	Comp:     compcmd.New(),
	Commands: []*Z.Cmd{help.Cmd, addCmd, delCmd, getCmd, patchCmd},
}

var addCmd = &Z.Cmd{
	Name:    `add`,
	Summary: `add a section to your ssh config file`,
	Usage:   `sectionName host host.com [field1 value1 field2 value2 ...]`,
	Params:  []string{"host", "user", "identityFile", "port"},
	Comp:    compcmd.New(),
	Description: `
		The {{aka}} command provides the ability to update your ssh config file
		through the command line.

		Options: 

		host [required]	:	hostname that should be used to establish the connection 

		user 			:	username to be used for the connection

		port 			:	port that the remote SSH daemon is running on. only necessary if the remote SSH instance is not running on the default port 22

		identityFile	:	private key that the client should use for authentication when connecting to the ssh server

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
	Usage:    `sectionName1 [sectionName2 ...]`,
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
var getCmd = &Z.Cmd{
	Name:     `get`,
	Summary:  `get sections from your ssh config file in YAML format`,
	Usage:    `sectionName1 [sectionName2 ...]`,
	Comp:     compcmd.New(),
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {

		c := NewCLI(os.Stdout)

		err := c.loadConfig()
		if err != nil {
			os.Exit(1)
		}

		c.parseConfig()
		s := c.getSections(args)
		c.printSections(s)

		return nil
	},
}

var patchCmd = &Z.Cmd{
	Name:    `patch`,
	Summary: `patch a section in your ssh config file`,
	Usage:   `sectionName field1 value1 [field2 value2 ...]`,
	Params:  []string{"host", "user", "identityFile", "port"},
	Comp:    compcmd.New(),
	Description: `
		The {{aka}} command provides the ability to patch a section in your ssh config file
		through the command line.

		Options: 

		host [required]	:	hostname that should be used to establish the connection 

		user 			:	username to be used for the connection

		port 			:	port that the remote SSH daemon is running on. only necessary if the remote SSH instance is not running on the default port 22

		identityFile	:	private key that the client should use for authentication when connecting to the ssh server

		`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {

		c := NewCLI(os.Stdout)

		err := c.loadConfig()
		if err != nil {
			os.Exit(1)
		}

		c.parseConfig()
		c.patchSection(args)
		c.createConfig()

		err = c.backupAndSave()
		if err != nil {
			os.Exit(1)
		}
		return nil
	},
}
