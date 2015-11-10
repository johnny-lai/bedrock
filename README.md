# bedrock

Bedrock is meant to provide a base on which to build your Go microservice. It comes with a `Makefile` that
can build, test and benchmark your service. The rules are written so that the build machine only needs `docker`
to be installed.

## Installation

1. Add bedrock as a submodule
   
   ```
   $ git submodule add -f git@github.com:johnny-lai/bedrock vendor/github.com/johnny-lai/bedrock
   ```

2. Add a [`glide.yaml`](https://github.com/Masterminds/glide) to your project

   ```
   $ glide create
   ```
   
3. Add dependencies needed by bedrock into your `glide.yaml`.

   ```
   - package: gopkg.in/yaml.v2
   - package: github.com/codegangsta/cli
   - package: github.com/gin-gonic/gin
   ```

4. Include the `boot.mk` into your Makefile to get all the bedrock build rules

   ```
   BEDROCK_ROOT = $(realpath vendor/github.com/johnny-lai/bedrock)
   include $(BEDROCK_ROOT)/boot.mk 

   APP_NAME = your-app-name
   APP_DOCKER_LABEL = your-docker-label  # Used for generating docker container labels
   APP_DOCKER_PUSH = yes       # Set to no to avoid publishing your docker image. Default is yes
   APP_GO_PACKAGES = packages  # Set to all the go package names that make up your service
   APP_GO_SOURCES = file.go    # Set to all the go source files used to build your main service
                               # Defaults to $(APP_NAME).go
   ```
   
5. Create your docker images. Place your `Dockerfile`s in the `docker/dist` and `docker/testdb` folder.

6. Create your integration tests. Place your tests in the `itest` folder, and your Kubernetes pod and service definitions in `itest/env`.

## Integrating with Jenkins

To initialize and build your project on Jenkins, you should use:

```
git submodule init
git submodule update
make deploy
```

## Companion container

The Makefile depends on scripts and custom behavior provided in the `johnnylai/bedrock-dev` docker images in order to
function. To build those images, use:

```
$ make deploy
```

### Starting Kubernetes

The companion container includes scripts to make it easier to start your own kubernetes cluster locally using docker.

```
# Enter the container image
$ make devconsole
$ cd /go
$ make kubernetes-start
```

Kubernetes will start and will keep running even after you exist the container. It will listen on port `8080` of the host.

If you are on a Mac, because docker runs in a host VM, it will actually be listening on the host VM's 8080. If you want
to use `kubectl` in the host itself, then you will need to forward port `8080`, using something like:

```
$ ssh -i  ~/.docker/machine/machines/default/id_rsa docker@$(MACHINE_DEFAULT_IP) -L8080:localhost:8080
```

Alternatively, you can always enter the container image using `make devconsole` and then run your `kubectl` command there
instead.

### Debugging Go

The companion container also contains [delve](https://github.com/derekparker/delve). So you can debug your program using
something like:

```
$ make devconsole
$ dlv debug
```

## boot.mk

* `deploy`: Build rule for Jenkins
* `dist`: Builds all the docker images
* `distutest`: Runs the unit tests in docker
* `distitest`: Runs the integration tests in docker
* `distibench`: Runs the benchmark tests in docker
* `fmt`: Runs `go fmt` on your Go packages
* `devconsole`: Enters the container image. Useful for starting kubernetes or running delve.

## Basing your Go service on bedrock

There is a sample Go service based on bedrock called [go-service-basic](https://github.com/johnny-lai/go-service-basic).

### main

Your main program would generally be something sort like:

```
package main

import (
	"github.com/johnny-lai/bedrock"
	"go-service-basic/core/service"
	"os"
)

var version = "unset"

func main() {
	app := bedrock.NewApp(&service.TodoService{})
	app.Name = "go-service-basic"
	app.Version = version
	app.Run(os.Args)
}
```

The `main.version` will be filled in by the `book.mk` during build.

### service

The service itself would need to implement the following interface:

```
type AppServicer interface {
  Config() interface{}
  Migrate() error
  Build(r *gin.Engine) error
  Run(r *gin.Engine) error
}
```
