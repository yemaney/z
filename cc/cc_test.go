package cc

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
	"testing"
)

func TestTypes(t *testing.T) {
	t.Run("Types written correctly", func(t *testing.T) {
		buffer, cli, _ := mockCLI("")

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
			_, cli, _ := mockCLI(tC.in)

			cli.readType()

			if cli.cc.typ != tC.want {
				t.Errorf("got %q want %q", cli.cc.typ, tC.want)
			}
		}
	})

	t.Run("Type chosen incorrectly results in prompt", func(t *testing.T) {
		buffer, cli, _ := mockCLI("")

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

		buffer, cli, _ := mockCLI("dependency")
		cli.readScope()

		promptGot := buffer.String()
		promptWant := "Enter a scope: "
		if promptGot != promptWant {
			t.Errorf("got %q want %q", promptGot, promptWant)
		}

		got := cli.cc.scope
		want := "(dependency)"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("Empty scope not set", func(t *testing.T) {

		buffer, cli, _ := mockCLI("")
		cli.readScope()

		promptGot := buffer.String()
		promptWant := "Enter a scope: "
		if promptGot != promptWant {
			t.Errorf("got %q want %q", promptGot, promptWant)
		}

		got := cli.cc.scope
		want := ""

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestSubject(t *testing.T) {

	t.Run("Subject set correctly", func(t *testing.T) {

		subject := "add github action to run tests"
		buffer, cli, _ := mockCLI(subject)

		cli.readSubject()

		promptGot := buffer.String()
		promptWant := "Enter a subject: "
		if promptGot != promptWant {
			t.Errorf("got %q want %q", promptGot, promptWant)
		}

		got := cli.cc.subject
		want := subject

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("Empty subject not set", func(t *testing.T) {
		buffer, cli, _ := mockCLI("")

		cli.readSubject()

		promptGot := buffer.String()
		promptWant := "Enter a subject: Enter a subject: Enter a subject: "
		if promptGot != promptWant {
			t.Errorf("got %q want %q", promptGot, promptWant)
		}

		got := cli.cc.subject
		want := ""

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestBodyAndFooter(t *testing.T) {

	t.Run("Body and Footer Correctly Set", func(t *testing.T) {

		dummyBody := "body"
		dummyFooter := "footer"
		buffer, cli, _ := mockCLI(dummyBody, dummyFooter)

		cli.readBodyAndFooter()

		promptGot := buffer.String()
		promptWant := "Enter a body: Enter a footer: "
		if promptGot != promptWant {
			t.Errorf("got %q want %q", promptGot, promptWant)
		}

		gotBody := cli.cc.body
		wantBody := dummyBody

		if gotBody != wantBody {
			t.Errorf("got %q want %q", gotBody, wantBody)
		}

		gotFooter := cli.cc.footer
		wantFooter := dummyFooter

		if gotFooter != wantFooter {
			t.Errorf("got %q want %q", gotFooter, wantFooter)
		}
	})

}

func TestMessage(t *testing.T) {
	t.Run("Complete conventional commit message created", func(t *testing.T) {
		typ := "1"
		scope := "dummy scope"
		subject := "dummy subject"
		body := "dummy body"
		footer := "dummmy footer"
		_, cli, _ := mockCLI(typ, scope, subject, body, footer)

		cli.readType()
		cli.readScope()
		cli.readSubject()
		cli.readBodyAndFooter()
		cli.buildMessage()

		got := cli.cc.message
		want := CCTypeMap[typ] + "(" + scope + "): " + subject + "\n\n" + body + "\n\n" + footer

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("Empty body and footer ignored", func(t *testing.T) {
		typ := "1"
		scope := "dummy scope"
		subject := "dummy subject"
		body := ""
		footer := ""
		_, cli, _ := mockCLI(typ, scope, subject, body, footer)

		cli.readType()
		cli.readScope()
		cli.readSubject()
		cli.readBodyAndFooter()
		cli.buildMessage()

		got := cli.cc.message
		want := CCTypeMap[typ] + "(" + scope + "): " + subject

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestConfirmationPrompt(t *testing.T) {
	typ := "1"
	scope := "dummy scope"
	subject := "dummy subject"
	body := "dummy body"
	footer := "dummmy footer"
	buffer, cli, _ := mockCLI(typ, scope, subject, body, footer)

	cli.readType()
	cli.readScope()
	cli.readSubject()
	cli.readBodyAndFooter()
	cli.buildMessage()
	cli.writeConfirmationPrompt()

	// slice to ignore the other prompts
	got := buffer.String()[62:]
	wantMsg := CCTypeMap[typ] + "(" + scope + "): " + subject + "\n\n" + body + "\n\n" + footer
	want := "\n\nPotential commit message:\n\n" + "\033[36;1m" + wantMsg + "\033[0m" + "\n\nCommit these changes with the message [y/N]: "

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestCommandExecuted(t *testing.T) {
	typ := "1"
	scope := "dummy scope"
	subject := "dummy subject"
	body := "dummy body"
	footer := "dummmy footer"

	testCases := []struct {
		confirm string
		want    string
	}{
		{
			confirm: "y",
			want:    "execute",
		},
		{
			confirm: "YES",
			want:    "execute",
		}, {
			confirm: "No",
			want:    "",
		}, {
			confirm: "random",
			want:    "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.confirm + " confirmation", func(t *testing.T) {
			_, cli, ce := mockCLI(typ, scope, subject, body, footer, tC.confirm)

			cli.readType()
			cli.readScope()
			cli.readSubject()
			cli.readBodyAndFooter()
			cli.buildMessage()
			cli.makeCommit()

			got := ce.command
			want := tC.want

			if got != want {
				t.Errorf("got %q want %q", got, want)
			}
		})
	}

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

func mockCLI(messages ...string) (*bytes.Buffer, *CLI, *mockCommandExecutor) {
	buffer := bytes.Buffer{}
	in := userSends(messages...)
	ce := mockCommandExecutor{}
	cli := NewCLI(&buffer, in, &ce)
	return &buffer, cli, &ce
}

type mockCommandExecutor struct {
	command string
}

func (mce *mockCommandExecutor) build(message string) *exec.Cmd {
	mce.command = "execute"
	return exec.Command(mce.command)
}
