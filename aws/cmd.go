package aws

import (
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	compcmd "github.com/yemaney/z/compcmd"
)

var Cmd = &Z.Cmd{
	Name:     `aws`,
	Summary:  `Helpful AWS commands`,
	Commands: []*Z.Cmd{help.Cmd, listCmd, getCmd, startCmd, stopCmd},
	Comp:     compcmd.New(),
}

var listCmd = &Z.Cmd{
	Name:     `list`,
	Summary:  `display the details and status of ec2 instances`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		list()
		return nil
	},
}

var getCmd = &Z.Cmd{
	Name:     `get`,
	Summary:  `display details for an instance that can be used for SSH connections`,
	Usage:    `instanceName`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		get(args[0])
		return nil
	},
	MinArgs: 1,
}

var startCmd = &Z.Cmd{
	Name:     `start`,
	Summary:  `start an ec2 instance`,
	Usage:    `instanceName`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		start(args[0])
		return nil
	},
	MinArgs: 1,
}

var stopCmd = &Z.Cmd{
	Name:     `stop`,
	Summary:  `stop and ec2 instance`,
	Usage:    `instanceName`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		stop(args[0])
		return nil
	},
	MinArgs: 1,
}
