package repl

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input: " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input: "  GoLang  is  Awesome ",
			expected: []string{"golang", "is", "awesome"},
		},
		{
			input: "",
			expected: []string{},
		},
		{
			input: "   ",
			expected: []string{},
		},
	}
		//further test cases can be added here (check for empty strings as well!)

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("Incorrect length of returning slice, expected a length of %d, received %d", len(c.expected), len(actual))
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Expected the word %q in position %d, received %q, from the original %q", expectedWord, i + 1, word, c.input)
			}
		}
	}
}
