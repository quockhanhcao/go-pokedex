package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "     hello     world   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU  ",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, testCase := range cases {
		actual := cleanInput(testCase.input)
		for i := range actual {
			word := actual[i]
			expectedWord := testCase.expected[i]
            if word != expectedWord {
                t.Errorf("Expected %s, but got %s", expectedWord, word)
                t.Fail()
            }
		}
	}
}
