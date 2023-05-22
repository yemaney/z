package cc

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestTypes(t *testing.T) {
	t.Run("Types written correctly", func(t *testing.T) {
		buffer := bytes.Buffer{}
		in := userSends("")
		cli := NewCLI(&buffer, in)

		cli.writeTypesPrompt()

		got := buffer.String()
		want := getTypesPrompt()

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("Types chosen correctly", func(t *testing.T) {
		testCases := []struct {
			in   string
			want string
		}{
			{
				in:   "0",
				want: "build",
			},
			{
				in:   "1",
				want: "ci",
			},
			{
				in:   "2",
				want: "docs",
			},
			{
				in:   "3",
				want: "feat",
			},
			{
				in:   "4",
				want: "fix",
			},
			{
				in:   "5",
				want: "perf",
			},
			{
				in:   "6",
				want: "refactor",
			},
			{
				in:   "7",
				want: "style",
			},
			{
				in:   "8",
				want: "test",
			},
			{
				in:   "9",
				want: "revert",
			},
			{
				in:   "10",
				want: "chore",
			},
		}
		for _, tC := range testCases {
			buffer := bytes.Buffer{}
			in := userSends(tC.in)
			cli := NewCLI(&buffer, in)

			cli.readType()

			if cli.cc.cctype != tC.want {
				t.Errorf("got %q want %q", cli.cc.cctype, tC.want)
			}
		}
	})

}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}
