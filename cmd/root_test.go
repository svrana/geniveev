package cmd

import (
	"testing"

	"github.com/svrana/geniveev"
)

func String(s string) *string {
	return &s
}

func TestConstructFilename(t *testing.T) {
	tests := []struct {
		inputFilename geniveev.Filename
		inputValues   map[string]*string
		want          string
	}{
		{inputFilename: "foo.go", inputValues: map[string]*string{}, want: "foo.go"},
		{inputFilename: "protos/{{.service_name}}.proto", inputValues: map[string]*string{"service_name": String("user")}, want: "protos/user.proto"},
		{inputFilename: "services/v1/{{.service_name}}/{{.service_name}}.go", inputValues: map[string]*string{"service_name": String("user")}, want: "services/v1/user/user.go"},
		{inputFilename: "{{Title \"protos\"}}", inputValues: map[string]*string{}, want: "Protos"},
	}

	for _, tc := range tests {
		filename, err := constructFilename(tc.inputFilename, tc.inputValues)
		if err != nil {
			t.Fatalf("expected success, got %s", err)
		}
		if filename != tc.want {
			t.Fatalf("expected %s, got %s", tc.want, filename)
		}
	}
}
