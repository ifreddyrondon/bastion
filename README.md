# gobastion

Defend your API from the sieges. Bastion offers an "augmented" Router instance.

It has the minimal necessary to create an API with default handlers and middleware that help you raise your API easy and fast.
Allows to have commons handlers and middleware between projects with the need for each one to do so.

## Router
Bastion use go-chi router to modularize the applications. Each instance of Bastion, will have the possibility
of mounting an api router, it will define the routes and middleware of the application with the app logic.

### Example

```go
package main

import (
	"net/http"

	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gobastion/utils"
)

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	utils.Send(w, res)
}

func main() {
	bastion := gobastion.New("")
	bastion.APIRouter.Get("/hello", helloHandler)
	bastion.Serve()
}
```

## Configuration
Represents the configuration for bastion. Config are used to define how the application should run.

###YAML
```yaml
api:
  base_path: "/"
server:
  address: ":8080"

```
###JSON
```json
{
  "api": {
    "base_path": "/"
  },
  "server": {
    "address": ":8080"
  }
}
```

### api
#### `api.base_path`
base path value where the application is going to be mounted. Default `/`.

```json
"base_path": "/foo/test",
```

```
http://localhost/foo/test
```

### `server`
#### `server.address`
Address is the host and port where the app is serve. Default `127.0.0.1:8080`.
When `server.address` is not provided it'll search the ADDR and PORT environment variables 
before set the default.

## Middlewares

Bastion comes equipped with a set of commons middlewares, providing a suite of standard
`net/http` middlewares.

Name | Description
---- | -----------
Logger | Logs the start and end of each request with the elapsed processing time
Recovery | Gracefully absorb panics and prints the stack trace
RequestID | Injects a request ID into the context of each request
