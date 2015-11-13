default: build

BEDROCK_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
include $(BEDROCK_ROOT)/make/env.mk
include $(BEDROCK_ROOT)/make/build.mk
include $(BEDROCK_ROOT)/make/itest.mk
include $(BEDROCK_ROOT)/make/utest.mk
include $(BEDROCK_ROOT)/make/ibench.mk
include $(BEDROCK_ROOT)/make/gen.mk

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
	           -e GO15VENDOREXPERIMENT=1 \
	           -it \
	           $(DOCKER_OPTS) \
	           $(DOCKER_DEVIMAGE)


.PHONY: build clean default deploy deps dist distbuild fmt migrate itest utest


