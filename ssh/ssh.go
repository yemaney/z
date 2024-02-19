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

type SSHError struct {
	Message string
}

// Error returns the error message for the custom error type.
func (ce *SSHError) Error() string {
	return ce.Message
}

type sshSection struct {
	host         string
	hostName     string
	user         string
	identityFile string
	port         int
}

type CLI struct {
	Out io.Writer
	*sshConfig
}

func NewCLI(out io.Writer) *CLI {
	return &CLI{
		Out:       out,
		sshConfig: &sshConfig{},
	}
}

func (c *CLI) createSection(args []string) error {
	s := &sshSection{}

	for i := 0; i < len(args); i += 2 {
		// Access the current and next elements
		flag := args[i]
		value := ""

		if i+1 < len(args) {
			value = args[i+1]
		}

		switch flag {
		case "name":
			s.host = value
		case "host":
			s.hostName = value
			if s.host == "" {
				s.host = value
			}
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
		}
	}

	if s.hostName == "" {
		fmt.Fprintln(c.Out, "Please provide the required flag: host")
		return &SSHError{"Please provide the required flag: host"}
	}

	c.newSection = s
	return nil
}

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

type sshConfig struct {
	sections   *[]sshSection
	newSection *sshSection
	config     *string
}

func (s *sshConfig) updateSections() {
	new := append(*s.sections, *s.newSection)
	s.sections = &new
}

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
	s.sections = &sections
}

func (s *sshConfig) createConfig() {
	var builder strings.Builder

	// Add a comment at the beginning of the file
	builder.WriteString("# Read more about SSH config files: https://linux.die.net/man/5/ssh_config\n")

	// Iterate over sections and build the file content
	for _, s := range *s.sections {
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
