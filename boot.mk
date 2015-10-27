GLIDE = $(GOPATH)/bin/glide

MAJOR_VERSION ?= 0
MINOR_VERSION ?= 0
BUILD_NUMBER ?= 0
COMMIT ?= $(shell git log --pretty=format:'%h' -n 1)
VERSION = $(MAJOR_VERSION).$(MINOR_VERSION).$(BUILD_NUMBER)

DOCKER_DEVIMAGE ?= johnnylai/golang-dev

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
	rm -f $(BUILD_ROOT)/$(APP_NAME)
	rm -rf tmp

build: deps
	GO15VENDOREXPERIMENT=1 go build \
		-o $(BUILD_ROOT)/$(APP_NAME) \
		-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT)" \
		$(APP_NAME).go

deps: $(GLIDE) $(BUILD_ROOT)
	if [ ! -d vendor ]; then $(GLIDE) update; fi

migrate:
	./cmd/server/server --config config.yaml migratedb

utest: deps
	TEST_CONFIG_YML=$(TEST_CONFIG_YML) GO15VENDOREXPERIMENT=1 go test $(APP_GO_PACKAGES)

itest:
	TEST_HOST="http://$(SERVER):$(PORT)" go test $(APP_NAME)/itest

bench:
	TEST_HOST="http://$(SERVER):$(PORT)" go test -bench=. $(APP_NAME)/itest

itestenv-restart: itestenv-stop itestenv-start
	
itestenv-start:
	for n in $(SRCROOT)/*.yml; do \
		cat $$n | $(KUBECTL) create -f - ; \
	done
	$(RUN_IN_DEV) wait-for-pod.sh go-service-basic

itestenv-stop:
	docker run --rm -i --net=host $(DOCKER_DEVIMAGE) kubectl delete all -lapp=$(APP_NAME)

fmt:
	GO15VENDOREXPERIMENT=1 go fmt $(APP_GO_PACKAGES)

devconsole:
	docker run --rm \
	           --net=host \
	           -v $(SRCROOT):$(SRCROOT_D) \
	           -w $(SRCROOT_D) \
	           -e GO15VENDOREXPERIMENT=1 \
	           -it \
	           $(DOCKER_DEVIMAGE)

dist: image-dist image-testdb

distbuild: clean build
	chown -R $(UID):$(GID) $(SRCROOT)

distitest:
	docker run --rm --net=host \
	           -v $(SRCROOT):$(SRCROOT_D) \
 	           -w $(SRCROOT_D) \
	           $(DOCKER_DEVIMAGE) \
	           make itest

distibench:
	docker run --rm --net=host \
	           -v $(SRCROOT):$(SRCROOT_D) \
 	           -w $(SRCROOT_D)/itest \
	           $(DOCKER_DEVIMAGE) \
	           make bench

distutest:
	-docker rm -f $(APP_NAME)-testdb
	docker run -d --name $(APP_NAME)-testdb $(APP_DOCKER_LABEL)-testdb
	docker run --rm \
	           --link $(APP_NAME)-testdb:$(APP_NAME)-db \
	           -v $(SRCROOT):$(SRCROOT_D) \
	           -w $(SRCROOT_D) \
	           -e DB_ENV_MYSQL_ROOT_PASSWORD=whatever \
	           -e TEST_CONFIG_YML=$(TEST_CONFIG_YML_D) \
	           $(DOCKER_DEVIMAGE) \
	           make utest

deploy: distutest dist itest-env-restart distitest
	docker push $(APP_DOCKER_LABEL)

.PHONY: build clean default deploy deps dist distbuild fmt migrate itest utest

image-testdb:
	docker build -f $(DOCKER_ROOT)/testdb/Dockerfile -t $(APP_DOCKER_LABEL)-testdb .

image-dist: $(PRODUCT_PATH)
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
	           -e UID=`id -u` \
	           -e GID=`id -g` \
	           $(DOCKER_DEVIMAGE) \
	           make distbuild
