COMMIT ?= $(shell git log --pretty=format:'%h' -n 1)

IMAGE_NAME = johnnylai/bedrock-dev

build:
	go build cli/bedrock.go

fmt:
	go fmt

image:
	docker build -t $(IMAGE_NAME) -f docker/dev/Dockerfile .
	docker tag -f $(IMAGE_NAME) $(IMAGE_NAME):$(COMMIT)

deploy: image
	docker push $(IMAGE_NAME)
	docker push $(IMAGE_NAME):$(COMMIT)


.PHONY: image
