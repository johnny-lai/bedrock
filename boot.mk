default: build

BEDROCK_ROOT := $(dir $(lastword $(MAKEFILE_LIST)))
include $(BEDROCK_ROOT)/make/env.mk
include $(BEDROCK_ROOT)/make/build.mk
include $(BEDROCK_ROOT)/make/itest.mk
include $(BEDROCK_ROOT)/make/utest.mk
include $(BEDROCK_ROOT)/make/ibench.mk

migrate:
	./cmd/server/server --config config.yaml migratedb

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


.PHONY: build clean default deploy deps dist distbuild fmt migrate itest utest


