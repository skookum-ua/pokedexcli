package main
import "testing"
func TestCleanInput(t *testing.T) {
	cases:= []struct {
		input string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Helloworld ",
			expected: []string{"helloworld"},
		},
		{
			input:    "Helloworld\n ",
			expected: []string{"helloworld"},
		},
		// add more cases here
	}

	for _, c := range cases {
	actual := cleanInput(c.input)
		if len(actual) == 0{
			t.Errorf("empty function output")
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("words dont metch")
			}
		}
}
}