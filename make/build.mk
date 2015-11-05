BUILD_ROOT ?= $(SRCROOT)

APP_DOCKER_LABEL_VERSION = $(APP_DOCKER_LABEL):$(MAJOR_VERSION).$(MINOR_VERSION)
APP_DOCKER_LABEL_COMMIT = $(APP_DOCKER_LABEL):$(COMMIT)

TESTDB_DOCKER_LABEL ?= $(APP_DOCKER_LABEL)-testdb
TESTDB_DOCKER_LABEL_VERSION = $(TESTDB_DOCKER_LABEL):$(MAJOR_VERSION).$(MINOR_VERSION)
TESTDB_DOCKER_LABEL_COMMIT = $(TESTDB_DOCKER_LABEL):$(COMMIT)

# GO flags
ifeq ($(APP_GO_LINKING), static)
	GO_ENV ?= GO15VENDOREXPERIMENT=1 CGO_ENABLED=0
	GO_CFLAGS ?= -a
else
	GO_ENV ?= GO15VENDOREXPERIMENT=1
	GO_CFLAGS ?=
endif

build: deps
	$(GO_ENV) go build $(GO_CFLAGS) \
		-o $(BUILD_ROOT)/$(APP_NAME) \
		-ldflags "-X main.version=$(VERSION)-$(COMMIT)" \
		$(APP_GO_SOURCES)

clean:
	go clean
	git clean -ffxd vendor
	rm -f $(BUILD_ROOT)/$(APP_NAME)
	rm -rf tmp

dist: image-dist image-testdb

distbuild: $(PRODUCT_PATH)

distpush: image-dist.push image-testdb.push

deploy: image-testdb distutest image-dist distitest distpush

image-testdb:
	docker build -f $(DOCKER_ROOT)/testdb/Dockerfile -t $(TESTDB_DOCKER_LABEL_COMMIT) $(SRCROOT)
	docker tag -f $(TESTDB_DOCKER_LABEL_COMMIT) $(TESTDB_DOCKER_LABEL)
	docker tag -f $(TESTDB_DOCKER_LABEL_COMMIT) $(TESTDB_DOCKER_LABEL_VERSION)

image-testdb.push:
	if [ "$(APP_DOCKER_PUSH)" == "yes" ]; then \
		$(DOCKER_PUSH) $(TESTDB_DOCKER_LABEL); \
		$(DOCKER_PUSH) $(TESTDB_DOCKER_LABEL_VERSION); \
		$(DOCKER_PUSH) $(TESTDB_DOCKER_LABEL_COMMIT); \
	fi

image-dist: distbuild
	docker build -f $(DOCKER_ROOT)/dist/Dockerfile -t $(APP_DOCKER_LABEL_COMMIT) $(SRCROOT)
	docker tag -f $(APP_DOCKER_LABEL_COMMIT) $(APP_DOCKER_LABEL)
	docker tag -f $(APP_DOCKER_LABEL_COMMIT) $(APP_DOCKER_LABEL_VERSION)

image-dist.push:
	if [ "$(APP_DOCKER_PUSH)" == "yes" ]; then \
		$(DOCKER_PUSH) $(APP_DOCKER_LABEL); \
		$(DOCKER_PUSH) $(APP_DOCKER_LABEL_VERSION); \
		$(DOCKER_PUSH) $(APP_DOCKER_LABEL_COMMIT); \
	fi

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


