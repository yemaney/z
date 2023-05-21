package cc

import (
	"bytes"
	"testing"
)

func Test(t *testing.T) {
	t.Run("Types written correctly", func(t *testing.T) {
		buffer := bytes.Buffer{}
		cli := CLI{&buffer}

		cli.writeTypes()

		got := buffer.String()
		want := getTypes()

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

}

// combines all conventional commit type optoins into one string
func getTypes() string {
	typesString := ""
	for _, v := range types {
		typesString += v
	}

	return typesString
}
