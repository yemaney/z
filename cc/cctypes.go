package cc

// CC represents a conventional commit
type CC struct {
	typ     string
	scope   string
	subject string
	body    string
	footer  string
	message string
}

// An array of strings corrosponding to a type used in a conventional commit.
//
// Each string is in the format of: `<number>. <type> <description>`. With the type being
// wrapped in cyan foreground color escape codes.
var types = [11]string{
	"0.  \033[36;1mbuild\033[0m:      Changes that affect the build system or external dependencies\n",
	"1.  \033[36;1mci\033[0m:         Changes to our CI configuration files and scripts\n",
	"2.  \033[36;1mdocs\033[0m:       Documentation only changes\n",
	"3.  \033[36;1mfeat\033[0m:       A new feature\n",
	"4.  \033[36;1mfix\033[0m:        A bug fix\n",
	"5.  \033[36;1mperf\033[0m:       A code change that improves performance\n",
	"6.  \033[36;1mrefactor\033[0m:   A code change that neither fixes a bug nor adds a feature\n",
	"7.  \033[36;1mstyle\033[0m:      Styling changes that don't affect the code performance or behavior\n",
	"8.  \033[36;1mtest\033[0m:       Adding missing tests or correcting existing tests\n",
	"9.  \033[36;1mrevert\033[0m:     If changes are reverted\n",
	"10. \033[36;1mchore\033[0m:      Ad-hoc task that doesn't match other types\n",
}

// CCTypeMap associates a number with each name of a conventional commit type.
// Meant to be used by the CLI to get a type name given a number input.
var CCTypeMap = map[string]string{
	"0":  "build",
	"1":  "ci",
	"2":  "docs",
	"3":  "feat",
	"4":  "fix",
	"5":  "perf",
	"6":  "refactor",
	"7":  "style",
	"8":  "test",
	"9":  "revert",
	"10": "chore",
}
