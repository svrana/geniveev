package builtins

import "testing"

func TestSnakeToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},
		{input: "foo", expected: "Foo"},
		{input: "user_service", expected: "UserService"},
	}

	for _, test := range tests {
		got := SnakeToCamelCase(test.input)
		if got != test.expected {
			t.Fatalf("got `%s`, expected: `%s`", got, test.expected)
		}
	}
}
