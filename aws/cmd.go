package aws

import (
	"fmt"
	"os"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	compcmd "github.com/yemaney/z/compcmd"
)

var Cmd = &Z.Cmd{
	Name:     `aws`,
	Summary:  `Helpful AWS commands`,
	Commands: []*Z.Cmd{help.Cmd, listCmd, getCmd, startCmd, stopCmd, createCmd, deleteCmd},
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

var createCmd = &Z.Cmd{
	Name:     `create`,
	Summary:  `create and ec2 instance`,
	Usage:    `name instanceName type t2.micro key sshkey securitygroup sg-0c00a01e168c82867`,
	Commands: []*Z.Cmd{help.Cmd, conf.Cmd},
	Params:   []string{"name", "key", "securitygroup", "type"},
	Comp:     compcmd.New(),
	Description: `
	{{cmd .Name}} is a command for creating an ubuntu based EC2 instance.

	name:	the name tag for the instance
	key:	the key pair name that will be used to make an SSH connection
	securitygroup:	the security group to associate with the EC2 instance
	type:	the instance type (ex t2.micro)


	The default values for the options can be configured using the **conf** command.
	The configuration is expected to be a YAML file with the following schema.

	---

	aws:

	{{ indent 4 "key: new_ssh" }}
	{{ indent 4 "type: t2.small" }}
	{{ indent 4 "securitygroup: sg-03718963ed05bc2f9" }}
	---

	`,
	Call: func(caller *Z.Cmd, args ...string) error {

		cArgs := CreateArgs{}

		// load defaults from config
		v, err := Z.Conf.Query(".aws.key")
		if err != nil {
			fmt.Println("Error loading conf")
		}
		if v != "null" {
			cArgs.key = v
		}
		v, _ = Z.Conf.Query(".aws.type")
		if v != "null" {
			cArgs.typ = v
		}
		v, _ = Z.Conf.Query(".aws.securitygroup")
		if v != "null" {
			cArgs.securitygroup = v
		}

		// overwrite defaults with any passed values
		for i := 0; i < len(args); i += 2 {
			if len(args) <= i+1 {
				break
			}
			if args[i] == "name" {
				cArgs.name = args[i+1]
			} else if args[i] == "key" {
				cArgs.key = args[i+1]
			} else if args[i] == "securitygroup" {
				cArgs.securitygroup = args[i+1]
			} else if args[i] == "type" {
				cArgs.typ = args[i+1]
			} else {
				fmt.Printf("Unsupported option %s\n", args[i])
				os.Exit(1)
			}

		}
		if cArgs.name == "" {
			fmt.Println("Require a name for the instance!")
			os.Exit(1)
		}
		ami := getLatestAmi()
		cArgs.ami = *ami.ImageId

		create(cArgs)

		return nil
	},
	MinArgs: 2,
}

var deleteCmd = &Z.Cmd{
	Name:     `delete`,
	Summary:  `delete an ec2 instance`,
	Usage:    `instanceName|instanceId`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		delete(args[0])
		return nil
	},
	MinArgs: 1,
}
