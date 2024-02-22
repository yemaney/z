package ssh

import (
	"bytes"
	"reflect"
	"testing"
)

func TestFlags(t *testing.T) {
	t.Run("Section created correctly", func(t *testing.T) {
		args := []string{"name", "example", "host", "example.com", "user", "root", "identityFile", "file.pem", "port", "22"}

		cli := NewCLI(&bytes.Buffer{})

		cli.createSection(args)
		want := &sshSection{
			host:         "example",
			hostName:     "example.com",
			user:         "root",
			identityFile: "file.pem",
			port:         22,
		}

		if !reflect.DeepEqual(cli.newSection, want) {
			t.Errorf("got %v want %v", cli.newSection, want)
		}
	})

	t.Run("Section missing flags raise error", func(t *testing.T) {
		args := []string{"executableName", "example.com", "user", "root", "identityFile", "file.pem"}

		buffer := bytes.Buffer{}

		cli := NewCLI(&buffer)
		got := cli.createSection(args)
		want := &SSHError{"Please provide the required flag: host"}

		if got == nil || got.Error() != want.Error() {
			t.Errorf("got %v want %v", got, want)
		}

		errorMessageGot := buffer.String()
		errorMessageWant := "Please provide the required flag: host\n"

		if errorMessageGot != errorMessageWant {
			t.Errorf("got %s want %s", errorMessageGot, errorMessageWant)
		}

	})
}

func TestParseConfig(t *testing.T) {
	fileContent := `# Read more about SSH config files: https://linux.die.net/man/5/ssh_config
Host test
	HostName test.com
	User testuser
	IdentityFile ~/Downloads/test.pem

Host build
	HostName build.com
	User builduser
	IdentityFile ~/Downloads/build.pem

Host sandbox
	HostName sandbox.com
	User sandboxuser

`
	sections := &[]sshSection{
		{host: "test", hostName: "test.com", user: "testuser", identityFile: "~/Downloads/test.pem"},
		{host: "build", hostName: "build.com", user: "builduser", identityFile: "~/Downloads/build.pem"},
		{host: "sandbox", hostName: "sandbox.com", user: "sandboxuser"},
	}

	t.Run("Parsing config works", func(t *testing.T) {
		cli := &CLI{
			sshConfig: &sshConfig{config: &fileContent},
		}

		cli.parseConfig()

		if !reflect.DeepEqual(cli.sections, sections) {
			t.Errorf("got %v, want %v", cli.sections, sections)
		}

	})

	t.Run("Creating config works", func(t *testing.T) {
		cli := &CLI{
			sshConfig: &sshConfig{config: &fileContent},
		}
		cli.parseConfig()
		cli.createConfig()

		want := fileContent
		got := *cli.config

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

	})
	t.Run("Adding section to config works", func(t *testing.T) {
		cli := &CLI{
			sshConfig: &sshConfig{config: &fileContent},
		}
		args := []string{"name", "example", "host", "example.com", "user", "root", "identityFile", "file.pem", "port", "22"}

		cli.createSection(args)
		cli.parseConfig()
		cli.updateSections()
		cli.createConfig()

		newConfigstr := `Host example
	HostName example.com
	User root
	IdentityFile file.pem
	Port  22

`
		want := fileContent + newConfigstr
		got := *cli.config

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		sectionWant := sshSection{host: "example", hostName: "example.com", user: "root", identityFile: "file.pem", port: 22}
		if !reflect.DeepEqual(cli.newSection, &sectionWant) {
			t.Errorf("got %v want %v", cli.newSection, sectionWant)
		}

		sectionsWant := append(*sections, sectionWant)
		if !reflect.DeepEqual(*cli.sections, sectionsWant) {
			t.Errorf("got %v want %v", cli.sections, sectionsWant)
		}
	})

	t.Run("Deleting section to config works", func(t *testing.T) {
		cli := &CLI{
			sshConfig: &sshConfig{config: &fileContent},
		}
		args := []string{"test", "build"}
		cli.parseConfig()
		cli.deleteSections(args)
		cli.createConfig()

		newConfigstr := `# Read more about SSH config files: https://linux.die.net/man/5/ssh_config
Host sandbox
	HostName sandbox.com
	User sandboxuser

`
		got := *cli.config

		if got != newConfigstr {
			t.Errorf("got %v, want %v", got, newConfigstr)
		}

		sectionsWant := []sshSection{{host: "sandbox", hostName: "sandbox.com", user: "sandboxuser"}}
		if !reflect.DeepEqual(*cli.sections, sectionsWant) {
			t.Errorf("got %v want %v", cli.sections, sectionsWant)
		}
	})
}
