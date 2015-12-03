# These are local paths
SRCROOT ?= $(abspath .)
BUILD_ROOT ?= $(SRCROOT)
DOCKER_ROOT ?= $(SRCROOT)/docker
TEST_CONFIG_YML ?= $(SRCROOT)/config/test.yml

# The current version
MAJOR_VERSION ?= 0
MINOR_VERSION ?= 0
BUILD_NUMBER ?= 0
COMMIT ?= $(shell git log --pretty=format:'%h' -n 1)
VERSION = $(MAJOR_VERSION).$(MINOR_VERSION).$(BUILD_NUMBER)

# Application settings
APP_NAME ?= unset
APP_DOCKER_LABEL ?= $(APP_NAME)
APP_GO_LINKING ?= static
APP_GO_SOURCES ?= main.go
APP_GO_PACKAGES ?= $(APP_NAME) $(APP_NAME)/core/service
APP_DOCKER_PUSH ?= yes
APP_SECRETS_ROOT ?= $(HOME)/.secrets/$(APP_NAME)
APP_ITEST_ENV_ROOT ?= $(SRCROOT)/itest/env

# These are paths used in the docker image
SRCROOT_D = /go/src/$(APP_NAME)
BUILD_ROOT_D = $(SRCROOT_D)/tmp/dist
TEST_CONFIG_YML_D = $(SRCROOT_D)/config/production.yml
APP_SECRETS_ROOT_D = /etc/secrets

# Docker Labels
APP_DOCKER_LABEL_VERSION = $(APP_DOCKER_LABEL):$(MAJOR_VERSION).$(MINOR_VERSION)
APP_DOCKER_LABEL_COMMIT = $(APP_DOCKER_LABEL):$(COMMIT)

TESTDB_DOCKER_LABEL ?= $(APP_DOCKER_LABEL)-testdb
TESTDB_DOCKER_LABEL_VERSION = $(TESTDB_DOCKER_LABEL):$(MAJOR_VERSION).$(MINOR_VERSION)
TESTDB_DOCKER_LABEL_COMMIT = $(TESTDB_DOCKER_LABEL):$(COMMIT)

# Docker commands
DOCKER_DEVIMAGE ?= johnnylai/bedrock-dev:17cbf55
DOCKER_DEV_UID ?= $(shell which docker-machine &> /dev/null || id -u)
DOCKER_DEV_GID ?= $(shell which docker-machine &> /dev/null || id -g)
DOCKER_OPTS ?= -v $(SRCROOT):$(SRCROOT_D) \
               -v $(KUBERNETES_CONFIG):/home/dev/.kube/config \
               -v $(KUBERNETES_CONFIG):/root/.kube/config \
               -v $(APP_SECRETS_ROOT):$(APP_SECRETS_ROOT_D) \
               -w $(SRCROOT_D) \
							 -e BUILD_ROOT=$(BUILD_ROOT_D) \
               -e APP_SECRETS_ROOT=$(APP_SECRETS_ROOT_D) \
							 -e BUILD_NUMBER=$(BUILD_NUMBER) \
               -e DEV_UID=$(DOCKER_DEV_UID) \
               -e DEV_GID=$(DOCKER_DEV_GID)
ifneq ($(findstring gcr.io/,$(APP_DOCKER_LABEL)),)
	DOCKER_PUSH ?= gcloud docker push
else
	DOCKER_PUSH ?= docker push
endif

# Kubernetes config
KUBERNETES_CONFIG ?= $(BEDROCK_ROOT)/make/kubernetes.config.default

# Executables
GLIDE = $(GOPATH)/bin/glide
BEDROCK = $(BUILD_ROOT)/bedrock
CLUSTER_SH = $(BEDROCK_ROOT)/scripts/cluster.sh

# Directory of gin. Used to detect if `glide update` is needed
GIN_ROOT = $(SRCROOT)/vendor/github.com/gin-gonic/gin

# Basic dependencies to build go programs
GO_BASE_DEPENDENCIES = $(GLIDE) $(BUILD_ROOT) $(GIN_ROOT)
deps: $(GO_BASE_DEPENDENCIES)

$(GLIDE): $(SRCROOT)/glide.yaml
	go get github.com/Masterminds/glide

$(BEDROCK): $(BEDROCK_ROOT)/cli/bedrock.go $(GO_BASE_DEPENDENCIES)
	GO15VENDOREXPERIMENT=1 go build -o $(BEDROCK) $(BEDROCK_ROOT)/cli/bedrock.go

$(BUILD_ROOT):
	mkdir -p $(BUILD_ROOT)

$(APP_SECRETS_ROOT):
	mkdir -p $@

$(GIN_ROOT): $(SRCROOT)/glide.yaml
	$(GLIDE) update

.PHONY: deps