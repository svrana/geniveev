# geniveev

An easy to use, language-agnostic code generation tool, powered by Go templates and an intuitive CLI.

## Why?

I enjoy working with protobufs. To create a new service, I need to generate a
number of files containing some boilerplate and I always need to go and look
at a previous service to remember where these files should be created, create those files
and copy/pasta a bit of boilerplate. This is made even worse if you use a dependency injection
tool that requires its own files, something I also sometimes do. This is my pain, but in any
project there is boilerplate and removing it can make programmers happy and I've found that to be
important.

# Configuration

To provide configuration details I will provide an example of a golang project, though geniveev is language agnostic.

## Example project setup

```
Project/
    .geniveev.toml      -- configuration file for this program
    protos/             -- all protobuf files here
      user/
        v1/
      auth/
        v1/

    gen/                -- all generated code from buf here
    services/           -- protobuf implementations here
        v1/
          auth/        -- implementation of AuthService
          user/        -- implementation of UserService
            .
            .
            .

```

## geniveev configuration file

Contents of .geniveev.json

(Note this config file is located in the example directory. If you wish to see geniveev do
her thing, run `make && cd examples && ../build/geniveev service-stubs --service-name User`

```toml
[service-stubs."protos/{{.service_name}}/v1/{{.service_name}}/{{.service_name}}.proto"]
    code = '''
syntax = "proto3";

package {{.service_name}}.v1;

import "validate/validate.proto";

// note that Title is a geniveev builtin that we use to capalize the Name of the service, i.e.,
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

If you run geniveev with the --help flag with the above .geniveev configuration
file in your current directory, you will see that a command called 'service-stubs'
is available to run. You just created that command by declaring it at the top level
of the toml configuration file. This name is arbitrary; name them as you see fit for your
project.

Following the top-level 'service-stubs' map is another another map, where the
key is a filename. The values between the brackets are also turned into
required CLI options. So, type `genvieev service-stubs --help`. You will see that
there is a required string option 'service_name'. Ok, let's provide one and see what
happens.

`geniveev service-stubs --service_name user`

Run git status and verify:

- protos/auth/v1/auth.proto was created
- services/v1/user/user.go was created

## Builtin template functions

Title
