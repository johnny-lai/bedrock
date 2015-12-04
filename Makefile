COMMIT ?= $(shell git log --pretty=format:'%h' -n 1)

DOCKER ?= docker

IMAGE_NAME = johnnylai/bedrock-dev

build:
	go build cli/bedrock.go

fmt:
	go fmt

image:
	$(DOCKER) build -t $(IMAGE_NAME) -f docker/dev/Dockerfile .
	$(DOCKER) tag -f $(IMAGE_NAME) $(IMAGE_NAME):$(COMMIT)

deploy: image
	$(DOCKER) push $(IMAGE_NAME)
	$(DOCKER) push $(IMAGE_NAME):$(COMMIT)


.PHONY: image
