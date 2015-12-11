COMMIT ?= $(shell git log --pretty=format:'%h' -n 1)

DOCKER ?= docker

IMAGE_NAME = johnnylai/bedrock-dev

build:
	go build cli/bedrock.go

fmt:
	go fmt

image.dev.golang:
	$(DOCKER) build -t $(IMAGE_NAME)-golang -f docker/dev/golang.dockerfile .
	$(DOCKER) tag -f $(IMAGE_NAME)-golang $(IMAGE_NAME)-golang:1.5
	$(DOCKER) tag -f $(IMAGE_NAME)-golang $(IMAGE_NAME)-golang:$(COMMIT)

image.swift:
	$(DOCKER) build -t johnnylai/swift:2.2 -f docker/swift/Dockerfile .

image.dev.swift: image.swift
	$(DOCKER) build -t $(IMAGE_NAME)-swift -f docker/dev/swift.dockerfile .
	$(DOCKER) tag -f $(IMAGE_NAME)-swift $(IMAGE_NAME)-swift:2.2
	$(DOCKER) tag -f $(IMAGE_NAME)-swift $(IMAGE_NAME)-swift:$(COMMIT)

deploy: image.dev.golang image.dev.swift
	$(DOCKER) push $(IMAGE_NAME)-golang
	$(DOCKER) push $(IMAGE_NAME)-golang:1.5
	$(DOCKER) push $(IMAGE_NAME)-swift
	$(DOCKER) push $(IMAGE_NAME)-swift:2.2

.PHONY: image
