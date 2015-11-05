GLIDE = $(GOPATH)/bin/glide

MAJOR_VERSION ?= 0
MINOR_VERSION ?= 0
BUILD_NUMBER ?= 0
COMMIT ?= $(shell git log --pretty=format:'%h' -n 1)
VERSION = $(MAJOR_VERSION).$(MINOR_VERSION).$(BUILD_NUMBER)

KUBERNETES_CONFIG ?= $(BEDROCK_ROOT)/make/kubernetes.config.default

DOCKER_DEVIMAGE ?= johnnylai/bedrock-dev
DOCKER_DEV_UID ?= $(shell which docker-machine &> /dev/null || id -u)
DOCKER_DEV_GID ?= $(shell which docker-machine &> /dev/null || id -g)
DOCKER_OPTS ?= -v $(SRCROOT):$(SRCROOT_D) \
               -v $(KUBERNETES_CONFIG):/home/dev/.kube/config \
               -v $(KUBERNETES_CONFIG):/root/.kube/config \
               -w $(SRCROOT_D) \
               -e DEV_UID=$(DOCKER_DEV_UID) \
               -e DEV_GID=$(DOCKER_DEV_GID)

APP_DOCKER_LABEL ?= unset
APP_GO_LINKING ?= dynamic
APP_GO_SOURCES ?= $(APP_NAME).go
APP_DOCKER_PUSH ?= yes
APP_ITEST_ENV_ROOT ?= $(SRCROOT)/itest/env



# These are local paths
SRCROOT ?= $(abspath .)

DOCKER_ROOT ?= $(SRCROOT)/docker
TEST_CONFIG_YML ?= $(SRCROOT)/config/test.yml
PRODUCT_PATH = tmp/dist/$(APP_NAME)

# These are paths used in the docker image
SRCROOT_D = /go/src/$(APP_NAME)
BUILD_ROOT_D = $(SRCROOT_D)/tmp/dist
TEST_CONFIG_YML_D = $(SRCROOT_D)/config/production.yml

#
SERVER ?= $(shell kubectl get svc $(APP_NAME) -o json | jq -r '.status.loadBalancer.ingress[0].ip')
PORT ?= $(shell kubectl get svc $(APP_NAME) -o json | jq '.spec.ports[0].targetPort')


deps: $(GLIDE) $(BUILD_ROOT)
	if [ ! -d vendor/github.com/gin-gonic/gin ]; then $(GLIDE) update; fi

$(GLIDE):
	go get github.com/Masterminds/glide

$(BUILD_ROOT):
	mkdir -p $@
