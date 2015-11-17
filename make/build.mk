BUILD_ROOT ?= $(SRCROOT)

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
	rm $(BEDROCK)
	rm -f $(BUILD_ROOT)/$(APP_NAME)
	rm -rf tmp

dist: image-dist image-testdb

distrun:
	docker run --rm \
	           --link $(APP_NAME)-testdb:$(APP_NAME)-db \
             -p 8080:8080 \
	           -v $(APP_SECRETS_ROOT):/etc/secrets \
	           $(APP_DOCKER_LABEL_COMMIT)

distrun.env:
	docker run --rm \
	           --link $(APP_NAME)-testdb:$(APP_NAME)-db \
	           -v $(APP_SECRETS_ROOT):/etc/secrets \
	           $(APP_DOCKER_LABEL_COMMIT) \
             env

distrun.sh:
	docker run --rm \
	           --link $(APP_NAME)-testdb:$(APP_NAME)-db \
	           -v $(APP_SECRETS_ROOT):/etc/secrets \
	           --entrypoint sh \
	           -it \
	           $(APP_DOCKER_LABEL_COMMIT)

distbuild: $(PRODUCT_PATH)

distpush: image-dist.push image-testdb.push

distpublish: image-dist.publish image-testdb.publish

deploy: image-testdb distutest image-dist distpush distitest

image-testdb:
	docker build -f $(DOCKER_ROOT)/testdb/Dockerfile -t $(TESTDB_DOCKER_LABEL_COMMIT) $(SRCROOT)
	docker tag -f $(TESTDB_DOCKER_LABEL_COMMIT) $(TESTDB_DOCKER_LABEL)
	docker tag -f $(TESTDB_DOCKER_LABEL_COMMIT) $(TESTDB_DOCKER_LABEL_VERSION)

image-testdb.push:
	if [ "$(APP_DOCKER_PUSH)" != "no" ]; then \
		$(DOCKER_PUSH) $(TESTDB_DOCKER_LABEL_COMMIT); \
	fi

image-testdb.publish:
	if [ "$(APP_DOCKER_PUSH)" != "no" ]; then \
		$(DOCKER_PUSH) $(TESTDB_DOCKER_LABEL); \
		$(DOCKER_PUSH) $(TESTDB_DOCKER_LABEL_VERSION); \
	fi

image-dist: distbuild
	docker build -f $(DOCKER_ROOT)/dist/Dockerfile -t $(APP_DOCKER_LABEL_COMMIT) $(SRCROOT)
	docker tag -f $(APP_DOCKER_LABEL_COMMIT) $(APP_DOCKER_LABEL)
	docker tag -f $(APP_DOCKER_LABEL_COMMIT) $(APP_DOCKER_LABEL_VERSION)

image-dist.push:
	if [ "$(APP_DOCKER_PUSH)" != "no" ]; then \
		$(DOCKER_PUSH) $(APP_DOCKER_LABEL_COMMIT); \
	fi

image-dist.publish:
	if [ "$(APP_DOCKER_PUSH)" != "no" ]; then \
		$(DOCKER_PUSH) $(APP_DOCKER_LABEL); \
		$(DOCKER_PUSH) $(APP_DOCKER_LABEL_VERSION); \
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


