package cc

import (
	"fmt"
	"io"
)

// Struct that defines the cli for this package. 
type CLI struct {
	Out io.Writer 
}

// writes conventional commit type options
func (c *CLI) writeTypes() {

	for _, v := range types {
		fmt.Fprint(c.Out, v)
	}
}

// An array of strings corrosponding to a type used in a conventional commit.
// 
// Each string is in the format of: `<number>. <type> <description>`. With the type being
// wrapped in cyan foreground color escape codes.
var types = [11]string{
	"0.  \033[36;1mbuild\033[0m:      Changes that affect the build system or external dependencies\n",
	"1.  \033[36;1mchore\033[0m:      Ad-hoc task that doesn't match other types\n",
	"2.  \033[36;1mci\033[0m:         Changes to our CI configuration files and scripts\n",
	"3.  \033[36;1mdocs\033[0m:       Documentation only changes\n",
	"4.  \033[36;1mfeat\033[0m:       A new feature\n",
	"5.  \033[36;1mfix\033[0m:        A bug fix\n",
	"6.  \033[36;1mperf\033[0m:       A code change that improves performance\n",
	"7.  \033[36;1mrefactor\033[0m:   A code change that neither fixes a bug nor adds a feature\n",
	"8.  \033[36;1mrevert\033[0m:     If changes are reverted\n",
	"9.  \033[36;1mstyle\033[0m:      Styling changes that don't affect the code performance or behavior\n",
	"10. \033[36;1mtest\033[0m:       Adding missing tests or correcting existing tests\n",
}
