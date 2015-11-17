# bedrock

Bedrock is meant to provide a base on which to build your Go microservice. It comes with a `Makefile` that
can build, test and benchmark your service. The rules are written so that the build machine only needs `docker`
to be installed.

## Installation

1. Add bedrock as a submodule
   
   ```
   $ git submodule add -f git@github.com:johnny-lai/bedrock vendor/github.com/johnny-lai/bedrock
   ```

2. Include the `boot.mk` into your Makefile to get all the bedrock build rules

   This is the minimal sample:
   ```
   APP_NAME = your-app-name
   include vendor/github.com/johnny-lai/bedrock/boot.mk
   ```

   This is a sample with all the options:
   ```
   APP_NAME = your-app-name
   APP_DOCKER_LABEL = your-docker-label  # Used for generating docker container labels
   APP_DOCKER_PUSH = yes       # Set to no to avoid publishing your docker image. Default is yes
   APP_GO_PACKAGES = packages  # Set to all the go package names that make up your service
   APP_GO_SOURCES = file.go    # Set to all the go source files used to build your main service
                               # Defaults to main.go
   include vendor/github.com/johnny-lai/bedrock/boot.mk 
   ```
	 
3. Commit your changes now. If you don't have a commit number, you will get a lot of warning messages during build.
 
4. Generate your application. You can generate the portions piece-meal:

   ```
	 # Generate a sample app
   $ make gen-app

   # Generate sample configs for the app
	 $ make gen-itest

   # Generate docker images
	 $ make gen-docker

   # Generate an integration test environment
	 $ make gen-itest

	 # Generate a Swagger API doc
	 $ make gen-api
   ```
	 
	 or all at once:
	 
	 ```
	 # make gen-all
	 ```

5. Generate your secrets
   ```
	 $ make gen-secret
	 ```
	 The script will ask for secrets like the Airbrake and New Relic keys and
	 put this into files in the `$HOME/.secrets/$APP_NAME` folder. This folder
	 will be mounted into the docker images, and used to generate kubernetes
	 secrets. The code and images themselves will not have these secrets.
	 
6. The generated README.md should contain more information on how to build and
   test your app.

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
