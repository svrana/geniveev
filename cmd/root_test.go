package cmd

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"

	"github.com/svrana/geniveev"
)

func ptr(s string) *string {
	return &s
}

func teardown() {
	AppFs = afero.NewOsFs()
}

func TestConstructFilename(t *testing.T) {
	tests := []struct {
		inputFilename geniveev.Filename
		inputValues   map[string]*string
		want          string
	}{
		{inputFilename: "foo.go", inputValues: map[string]*string{}, want: "foo.go"},
		{inputFilename: "protos/{{.service_name}}.proto", inputValues: map[string]*string{"service_name": ptr("user")}, want: "protos/user.proto"},
		{inputFilename: "services/v1/{{.service_name}}/{{.service_name}}.go", inputValues: map[string]*string{"service_name": ptr("user")}, want: "services/v1/user/user.go"},
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

func TestCreatePath(t *testing.T) {
	AppFs = afero.NewMemMapFs()
	defer teardown()

	directory := "/foo/bar"
	filename := filepath.Join(directory, "baz.go")
	// non-existent path should succeed
	if err := createPath(filename); err != nil {
		t.Fatalf("expected success got %s", err)
	}
	_, err := AppFs.Stat(directory)
	if err != nil {
		t.Fatalf("file does not exist")
	}
	if err = afero.WriteFile(AppFs, filename, []byte("file b"), 0644); err != nil {
		t.Fatalf("unexpected error setting up test: %s", err)
	}
	// now try again with an existing file and make sure it fails as we do not want
	// to overwrite any existing files.
	err = createPath(filename)
	if err == nil {
		t.Fatalf("unexpected success for existing file")
	}
}

func TestIntegration(t *testing.T) {
	// first read in example geniveev configuration file from disk
	absFilepath, err := filepath.Abs(filepath.Join("..", "example", cfgFile))
	if err != nil {
		t.Fatalf("failed to get absolute filename of example config file: %s", err)

	}
	b, err := afero.ReadFile(AppFs, absFilepath)
	if err != nil {
		t.Fatalf("failed to read config file: %s", err)
	}
	AppFs = afero.NewMemMapFs()
	defer teardown()
	if err = afero.WriteFile(AppFs, cfgFile, b, 0644); err != nil {
		t.Fatalf("failed to write test config to in memory filesystem for test setup: %s", err)
	}
	if err = Initialize(); err != nil {
		t.Fatalf("initialize returned an error: %s", err)
	}
	if v := config.Generator["service-stubs"]; v == nil {
		t.Fatalf("service-stubs not parsed")
	}

	config.Generator["service-stubs"].TemplateValues["service_name"] = ptr("user") // match the readme

	if err = start(config.Generator["service-stubs"]); err != nil {
		t.Fatalf("failed to generate code: %s", err)
	}
	userProto := "protos/user/v1/user/user.proto"
	_, err = AppFs.Stat(userProto)
	if err != nil {
		t.Fatalf("could not locate generated file")
	}
	userProtoMem, err := afero.ReadFile(AppFs, userProto)
	if err != nil {
		t.Fatalf("failed to read generated proto: %s", err)
	}
	expected := `syntax = "proto3";

package user.v1;

import "validate/validate.proto";

// note that Title is a geniveev builtin that we use to capitalize the Name of the service, i.e.,
// UserService
service UserService {
}
`
	if string(userProtoMem) != expected {
		t.Fatalf("expected\n%s\ngot\n%s\n", expected, userProtoMem)
	}

	implProto := "services/v1/user/user.go"
	_, err = AppFs.Stat(implProto)
	if err != nil {
		t.Fatalf("could not locate generated file")
	}
	userImplMem, err := afero.ReadFile(AppFs, implProto)
	if err != nil {
		t.Fatalf("failed to read generated user implementation: %s", err)
	}
	expected = `package userv1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"

	"github.com/bommie/b6/config"
	"github.com/bommie/b6/db"
)

type UserServer struct {
	db  *db.DB
	cfg *config.AuthConfig
}

var _ userv1connect.UserServiceClient = (*UserServer)(nil)

func NewServer(_ context.Context, db *db.DB, cfg *config.AuthConfig) *UserServer {
    return &UserServer{
		db:  db,
		cfg: cfg,
	}
}
`
	if string(userImplMem) != expected {
		t.Fatalf("expected \n%s\ngot \n%s\n", expected, userImplMem)
	}
}
