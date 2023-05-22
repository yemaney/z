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

	t.Run("Type chosen incorrectly results in prompt", func(t *testing.T) {
		buffer := bytes.Buffer{}
		in := userSends("")
		cli := NewCLI(&buffer, in)

		cli.readType()

		got := buffer.String()
		want := typeErrorMsg()

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})
}

func TestScope(t *testing.T) {

	t.Run("Scope set correctly", func(t *testing.T) {
		buffer := bytes.Buffer{}
		in := userSends("0", "dependency")
		cli := NewCLI(&buffer, in)

		cli.readType()
		cli.readScope()

		got := cli.cc.scope
		want := "(dependency)"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("Empty scope not set", func(t *testing.T) {
		buffer := bytes.Buffer{}
		in := userSends("0", "")
		cli := NewCLI(&buffer, in)

		cli.readType()
		cli.readScope()

		got := cli.cc.scope
		want := ""

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}

func typeErrorMsg() string {
	msg := ""
	for i := 0; i < 2; i++ {
		msg += "Enter a valid number between 0 and 10: "
	}
	return msg
}
