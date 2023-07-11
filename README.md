# geniveev

An easy to use, language-agnostic code generation tool, powered by Go templates and an intuitive CLI.

## Why?

I enjoy working with protobufs. To create a new service, I need to generate a
number of files containing some boilerplate and I always need to go and look
at a previous service to remember where these files should be created, create those files
and copy/pasta a bit of boilerplate. This is made even worse if you use a dependency injection
tool that requires its own files, something I also sometimes do. This little tool was
written to automatically create those files with some simple boilerplate to speed up the
process of creating a new service.

# Configuration

To configure geniveev, you must construct a toml .geniveev configuration file in your
current working directory.

## geniveev configuration file

This is an example configuration file included in the example/directory.

The geniveev configuration file contains command definitions and any number of templates associated with the
with these definitions. The templates are associated with a filename that is itself a
template, allowing the filename and its contents to be specified per-run via the command
line.

In the following configuration, we have defined the `service-stubs` command which consists of two templates. When geniveev
is run, it generates a service-stubs command (and any others you add to the configuration
file) and adds the service_name as a required parameter for that command. Any template
value found the filename can be reused in the template itself.

Note that it is common to want to manipulate the filename in the file-content template itself, so some builtins are provided,
as shown below.

The following configuration file will generate two files for a protobuf setup. To run geniveev using this configuration file,
clone this repo and run `make && cd examples && ../build/gen service-stubs --service-name User`.

```toml
[service-stubs."protos/{{.service_name}}/v1/{{.service_name}}/{{.service_name}}.proto"]
    code = '''
syntax = "proto3";

package {{.service_name}}.v1;

import "validate/validate.proto";

// note that Title is a geniveev builtin that we use to capitalize the Name of the service, i.e.,
// UserService
service {{ .service_name | Title }}Service {
}
'''
[service-stubs.'services/v1/{{.service_name}}/{{.service_name}}.go']
    code = '''
package {{.service_name}}v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"

	"github.com/bommie/b6/config"
	"github.com/bommie/b6/db"
)

type {{.service_name | Title}}Server struct {
	db  *db.DB
	cfg *config.AuthConfig
}

var _ {{.service_name}}v1connect.{{.service_name | Title}}ServiceClient = (*{{.service_name | Title }}Server)(nil)

func NewServer(_ context.Context, db *db.DB, cfg *config.AuthConfig) *{{.service_name | Title }}Server {
    return &{{.service_name | Title}}Server{
		db:  db,
		cfg: cfg,
	}
}
'''
```

After running `../build/gen service-stubs --service-name User` run `git status` and verify
that the following files exist and look as you expect.

- protos/auth/v1/auth.proto was created
- services/v1/user/user.go was created

## Builtin template functions

Title
