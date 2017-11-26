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

var bastion *gobastion.Bastion

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	res := struct {
		Message string `json:"message"`
	}{"world"}
	utils.Send(w, res)
}

func main() {
	bastion = gobastion.NewBastion("")
	bastion.APIRouter.Get("/hello", helloHandler)
	bastion.Serve(":8080")
}
```

## Configuration
Represents the configuration for bastion. Config are used to define how the application should run.

###YAML
```yaml
api:
  base_path: "/api/"

```
###JSON
```json
{
  "api": {
    "base_path": "/"
  }
}
```


### `basePath`
basePath value where the application is going to be mounted. Default `/`.

```json
"base_path": "/foo/test",
```

```
http://localhost/foo/test
```
