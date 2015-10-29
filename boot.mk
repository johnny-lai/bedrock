GLIDE = $(GOPATH)/bin/glide

MAJOR_VERSION ?= 0
MINOR_VERSION ?= 0
BUILD_NUMBER ?= 0
COMMIT ?= $(shell git log --pretty=format:'%h' -n 1)
VERSION = $(MAJOR_VERSION).$(MINOR_VERSION).$(BUILD_NUMBER)

DOCKER_DEVIMAGE ?= johnnylai/bedrock-dev:275ad78
DOCKER_DEV_UID ?= $(shell which docker-machine &> /dev/null || id -u)
DOCKER_DEV_GID ?= $(shell which docker-machine &> /dev/null || id -g)

APP_ITEST_ENV_ROOT ?= $(SRCROOT)/itest/env

# These are local paths
SRCROOT ?= $(realpath .)
BUILD_ROOT ?= $(SRCROOT)
DOCKER_ROOT ?= $(SRCROOT)/docker
TEST_CONFIG_YML ?= $(SRCROOT)/config/test.yml
PRODUCT_PATH = tmp/dist/$(APP_NAME)

# These are paths used in the docker image
SRCROOT_D = /go/src/$(APP_NAME)
BUILD_ROOT_D = $(SRCROOT_D)/tmp/dist
TEST_CONFIG_YML_D = $(SRCROOT_D)/config/production.yml

#
SERVER ?= $(shell kubectl get svc go-service-basic -o json | jq -r '.spec.clusterIP')
PORT ?= $(shell kubectl get svc go-service-basic -o json | jq '.spec.ports[0].targetPort')

# For itest
RUN_IN_DEV = docker run --rm --net=host -i $(DOCKER_DEVIMAGE)
KUBECTL = $(RUN_IN_DEV) kubectl

default: build

clean:
	go clean
	git clean -ffxd vendor
	rm -f $(BUILD_ROOT)/$(APP_NAME)
	rm -rf tmp

build: deps
	GO15VENDOREXPERIMENT=1 go build \
		-o $(BUILD_ROOT)/$(APP_NAME) \
		-ldflags "-X main.version=$(VERSION)-$(COMMIT)" \
		$(APP_NAME).go

deps: $(GLIDE) $(BUILD_ROOT)
	if [ ! -d vendor/github.com/gin-gonic/gin ]; then $(GLIDE) update; fi

migrate:
	./cmd/server/server --config config.yaml migratedb

utest: deps
	TEST_CONFIG_YML=$(TEST_CONFIG_YML) GO15VENDOREXPERIMENT=1 go test $(APP_GO_PACKAGES)

itest:
	TEST_HOST="http://$(SERVER):$(PORT)" go test $(APP_NAME)/itest

bench:
	TEST_HOST="http://$(SERVER):$(PORT)" go test -bench=. $(APP_NAME)/itest

itest.env.start:
	for n in $(APP_ITEST_ENV_ROOT)/*.yml; do \
		cat $$n | $(KUBECTL) create -f - ; \
	done
	$(RUN_IN_DEV) wait-for-pod.sh go-service-basic

itest.env.stop:
	-docker run --rm -i --net=host $(DOCKER_DEVIMAGE) kubectl delete all -lapp=$(APP_NAME)

fmt:
	GO15VENDOREXPERIMENT=1 go fmt $(APP_GO_PACKAGES)

devconsole:
	docker run --rm \
	           --net=host \
	           -v `which docker`:/bin/docker \
	           -v /var/run/docker.sock:/var/run/docker.sock \
	           -v /lib64/libdevmapper.so.1.02:/lib/libdevmapper.so.1.02 \
	           -v /lib64/libudev.so.0:/lib/libudev.so.0 \
	           -v $(SRCROOT):$(SRCROOT_D) \
	           -w $(SRCROOT_D) \
	           -e DEV_UID=$(DOCKER_DEV_UID) \
	           -e DEV_GID=$(DOCKER_DEV_GID) \
	           -e GO15VENDOREXPERIMENT=1 \
	           -it \
	           $(DOCKER_DEVIMAGE)

dist: image-dist image-testdb

distbuild: $(PRODUCT_PATH)

#-------------------------------------------------------------------------------
distitest: distitest.env distitest.run

distitest.env: itest.env.stop itest.env.start

distitest.run:
	docker run --rm --net=host \
	           -v $(SRCROOT):$(SRCROOT_D) \
 	           -w $(SRCROOT_D) \
	           -e DEV_UID=$(DOCKER_DEV_UID) \
	           -e DEV_GID=$(DOCKER_DEV_GID) \
	           $(DOCKER_DEVIMAGE) \
	           make itest

#-------------------------------------------------------------------------------
distibench: distibench.env distibench.run

distibench.env: itest.env.stop itest.env.start

distibench.run:
	docker run --rm --net=host \
	           -v $(SRCROOT):$(SRCROOT_D) \
 	           -w $(SRCROOT_D)/itest \
	           -e DEV_UID=$(DOCKER_DEV_UID) \
	           -e DEV_GID=$(DOCKER_DEV_GID) \
	           $(DOCKER_DEVIMAGE) \
	           make bench

#-------------------------------------------------------------------------------
distutest: distutest.env distutest.run

distutest.env:
	-docker rm -f $(APP_NAME)-testdb
	docker run -d --name $(APP_NAME)-testdb $(APP_DOCKER_LABEL)-testdb

distutest.run:
	docker run --rm \
	           --link $(APP_NAME)-testdb:$(APP_NAME)-db \
	           -v $(SRCROOT):$(SRCROOT_D) \
	           -w $(SRCROOT_D) \
	           -e DEV_UID=$(DOCKER_DEV_UID) \
	           -e DEV_GID=$(DOCKER_DEV_GID) \
	           -e DB_ENV_MYSQL_ROOT_PASSWORD=whatever \
	           -e TEST_CONFIG_YML=$(TEST_CONFIG_YML_D) \
	           $(DOCKER_DEVIMAGE) \
	           make utest

#-------------------------------------------------------------------------------
deploy: distutest dist distitest
	docker push $(APP_DOCKER_LABEL)

.PHONY: build clean default deploy deps dist distbuild fmt migrate itest utest

image-testdb:
	docker build -f $(DOCKER_ROOT)/testdb/Dockerfile -t $(APP_DOCKER_LABEL)-testdb .

image-dist: distbuild
	docker build -f $(DOCKER_ROOT)/dist/Dockerfile -t $(APP_DOCKER_LABEL) .

$(GLIDE):
	go get github.com/Masterminds/glide

$(BUILD_ROOT):
	mkdir -p $@

$(PRODUCT_PATH): $(wildcard *.go)
	docker run --rm \
	           -v $(SRCROOT):$(SRCROOT_D) \
	           -w $(SRCROOT_D) \
	           -e BUILD_ROOT=$(BUILD_ROOT_D) \
	           -e BUILD_NUMBER=$(BUILD_NUMBER) \
	           -e DEV_UID=$(DOCKER_DEV_UID) \
	           -e DEV_GID=$(DOCKER_DEV_GID) \
	           $(DOCKER_DEVIMAGE) \
	           make build
