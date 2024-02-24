// Package ssh provides functionality for interacting with SSH configurations.
package ssh

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// SSHError is a custom error type for SSH-related errors.
type SSHError struct {
	Message string
}

// Error returns the error message for the custom error type.
func (ce *SSHError) Error() string {
	return ce.Message
}

// sshSection represents a section within an SSH configuration.
type sshSection struct {
	host         string
	hostName     string
	user         string
	identityFile string
	port         int
}

// toYAML converts an sshSection to a YAML formatted string.
func (s sshSection) toYAML() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s:\n", s.host))

	if s.hostName != "" {
		builder.WriteString(fmt.Sprintf("  hostName: %s\n", s.hostName))
	}

	if s.user != "" {
		builder.WriteString(fmt.Sprintf("  user: %s\n", s.user))
	}

	if s.identityFile != "" {
		builder.WriteString(fmt.Sprintf("  identityFile: %s\n", s.identityFile))
	}

	if s.port != 0 {
		builder.WriteString(fmt.Sprintf("  port: %v\n", s.port))
	}

	return builder.String()
}

// CLI represents a command-line interface for managing SSH configurations.
type CLI struct {
	Out io.Writer
	*sshConfig
}

// NewCLI creates a new CLI instance with the specified output writer.
func NewCLI(out io.Writer) *CLI {
	return &CLI{
		Out:       out,
		sshConfig: &sshConfig{},
	}
}

// createSection creates a new SSH section based on the provided command-line arguments.
func (c *CLI) createSection(args []string) error {
	s := &sshSection{}

	if len(args) < 3 {
		return &SSHError{"Not Enough Parameters."}
	}

	s.host = args[0]

	for i := 1; i < len(args); i += 2 {
		// Access the current and next elements
		param := args[i]
		value := ""

		if i+1 < len(args) {
			value = args[i+1]
		}

		switch param {
		case "host":
			s.hostName = value
		case "user":
			s.user = value
		case "identityFile":
			s.identityFile = value
		case "port":
			num, err := strconv.Atoi(value)
			if err != nil {
				fmt.Fprintf(c.Out, "Error with port %s for host %s. Reverting to default\n", value, s.host)
				s.port = 22
			} else {
				s.port = num
			}
		default:
			fmt.Fprintf(c.Out, "Unsupported Parameter: %s\n", param)
			return &SSHError{"Unsupported Parameter"}
		}
	}

	if s.hostName == "" {
		fmt.Fprintln(c.Out, "Please provide the required flag: host")
		return &SSHError{"Please provide the required flag: host"}
	}

	c.newSection = s
	return nil
}

// backupAndSave backs up the existing SSH config file and saves the new configuration.
func (c *CLI) backupAndSave() error {
	// Set the file paths
	currentUser, _ := user.Current()
	configFilePath := filepath.Join(currentUser.HomeDir, ".ssh", "config")
	backupFilePath := filepath.Join(currentUser.HomeDir, ".ssh", ".config.backup")

	// Read the existing content of the config file
	existingConfig, _ := os.ReadFile(configFilePath)

	// Backup the existing config file
	if err := os.WriteFile(backupFilePath, existingConfig, 0644); err != nil {
		fmt.Fprintf(c.Out, "Error writing backup: %s\n", err)
		return err
	}

	// Save the new config text to the config file
	if err := os.WriteFile(configFilePath, []byte(*c.config), 0644); err != nil {
		fmt.Fprintf(c.Out, "Error new config: %s\n", err)
		return err
	}

	return nil
}

// printSections prints the YAML representation of SSH sections.
func (c *CLI) printSections(sections []sshSection) {
	sc := ""

	for _, v := range sections {
		sc += v.toYAML()
	}

	fmt.Fprintln(c.Out, sc)
}

// sshConfig represents the overall SSH configuration and manages its sections.
type sshConfig struct {
	sections   []sshSection
	config     *string
	newSection *sshSection
}

// updateSections adds a new section to the SSH configuration.
func (s *sshConfig) updateSections() {
	s.sections = append(s.sections, *s.newSection)
}

// deleteSections removes specified sections from the SSH configuration.
func (s *sshConfig) deleteSections(args []string) {

	m := map[string]sshSection{}

	for _, sc := range s.sections {
		delete := false
		for _, v := range args {
			if sc.host == v {
				delete = true
				break
			}

		}
		if !delete {
			m[sc.host] = sc
		}
	}

	n := []sshSection{}
	for _, v := range m {
		n = append(n, v)
	}

	s.sections = n
}

// getSections retrieves specified sections from the SSH configuration.
func (s *sshConfig) getSections(args []string) []sshSection {
	n := []sshSection{}

	if len(args) == 0 {
		return n
	}
	if args[0] == "all" {
		return s.sections
	}

	m := map[string]sshSection{}

	for _, sc := range s.sections {
		found := false
		for _, v := range args {
			if sc.host == v {
				found = true
				break
			}

		}
		if found {
			m[sc.host] = sc
		}
	}

	for _, v := range m {
		n = append(n, v)
	}

	return n
}

// loadConfig reads the SSH configuration file and sets it as the sshConfig's config field value
func (s *sshConfig) loadConfig() error {

	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:\n", err)
		return err
	}

	// Construct the path to the ~/.ssh/config file
	sshConfigPath := filepath.Join(currentUser.HomeDir, ".ssh", "config")

	// Check if the file exists
	if _, err := os.Stat(sshConfigPath); os.IsNotExist(err) {
		// File does not exist
		fmt.Printf("The file %s does not exist.\n", sshConfigPath)
		fmt.Print("Do you want to create it? (y/n): ")

		// Read user input
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := strings.ToLower(scanner.Text())

		if input == "y" {
			// Create the file
			err := os.WriteFile(sshConfigPath, []byte{}, 0644)
			if err != nil {
				fmt.Println("Error creating ~/.ssh/config file:\n", err)
				return err
			}
			fmt.Printf("File %s created successfully.\n", sshConfigPath)
		} else {
			fmt.Println("File creation canceled.")
			os.Exit(1)
		}
	}

	// Read the contents of the file
	content, err := os.ReadFile(sshConfigPath)
	if err != nil {
		fmt.Println("Error reading ~/.ssh/config file:\n", err)
		return err
	}

	config := string(content)
	s.config = &config
	return nil
}

// parseConfig parses the content of the SSH configuration file into a []sshSection and sets it as the
// sshConfig's sections field value
func (s *sshConfig) parseConfig() {
	var sections []sshSection
	var currentSection sshSection

	scanner := bufio.NewScanner(strings.NewReader(*s.config))
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// Check for Host line to start a new section
		if strings.HasPrefix(line, "Host ") {
			// If there's a current section, add it to the slice
			if currentSection.host != "" {
				sections = append(sections, currentSection)
			}

			// Initialize a new section
			currentSection = sshSection{host: strings.TrimSpace(line[5:])}
		} else {
			// Parse other lines within a section
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				key := fields[0]
				value := strings.Join(fields[1:], " ")

				switch key {
				case "HostName":
					currentSection.hostName = value
				case "User":
					currentSection.user = value
				case "IdentityFile":
					currentSection.identityFile = value
				case "Port":
					num, err := strconv.Atoi(value)
					if err != nil {
						fmt.Printf("Error with port %s for host %s. Reverting to default\n", value, currentSection.host)
						currentSection.port = 22
					} else {
						currentSection.port = num
					}
				}

			}
		}
	}

	// Add the last section to the slice
	if currentSection.host != "" {
		sections = append(sections, currentSection)
	}
	s.sections = sections
}

// createConfig generates the content of the SSH configuration file from  the sshConfig's sections field
// and sets it as the config field value.
func (s *sshConfig) createConfig() {
	var builder strings.Builder

	// Add a comment at the beginning of the file
	builder.WriteString("# Read more about SSH config files: https://linux.die.net/man/5/ssh_config\n")

	// Iterate over sections and build the file content
	for _, s := range s.sections {
		builder.WriteString(fmt.Sprintf("Host %s\n", s.host))

		if s.hostName != "" {
			builder.WriteString(fmt.Sprintf("\tHostName %s\n", s.hostName))
		}

		if s.user != "" {
			builder.WriteString(fmt.Sprintf("\tUser %s\n", s.user))
		}

		if s.identityFile != "" {
			builder.WriteString(fmt.Sprintf("\tIdentityFile %s\n", s.identityFile))
		}

		if s.port != 0 {
			builder.WriteString(fmt.Sprintf("\tPort  %v\n", s.port))
		}

		builder.WriteString("\n")
	}

	config := builder.String()
	s.config = &config
}

// patchSection updates specified fields in an existing SSH section.
func (s *sshConfig) patchSection(args []string) error {

	if len(args) < 3 || len(args)%2 != 1 {
		fmt.Println("Usage: patchSection <host> <field1> <value1> [<field2> <value2> ...]")
		return &SSHError{}
	}

	host := args[0]
	index := -1
	var section sshSection

	// Find the section with the specified host
	for i, sc := range s.sections {
		if sc.host == host {
			index = i
			section = s.sections[i]
			break
		}
	}

	if index == -1 {
		fmt.Printf("Section %s does not exist.\n", host)
		return &SSHError{}
	}

	// Update the specified fields in the section
	for i := 1; i < len(args); i += 2 {
		field := args[i]
		value := args[i+1]

		switch field {
		case "hostName":
			section.hostName = value
		case "user":
			section.user = value
		case "identityFile":
			section.identityFile = value
		case "port":
			port, err := strconv.Atoi(value)
			if err != nil {
				fmt.Printf("Invalid port value for field %s.\n", field)
				return &SSHError{}
			}
			section.port = port
		default:
			fmt.Printf("Invalid field name: %s.\n", field)
			return &SSHError{}
		}
	}

	s.sections[index] = section

	return nil
}
