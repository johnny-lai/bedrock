COMMIT ?= $(shell git log --pretty=format:'%h' -n 1)

IMAGE_NAME = johnnylai/bedrock-dev

image:
	docker build -t $(IMAGE_NAME) -f docker/dev/Dockerfile docker/dev
	docker tag -f $(IMAGE_NAME) $(IMAGE_NAME):$(COMMIT)

deploy: image
	docker push $(IMAGE_NAME)
	docker push $(IMAGE_NAME):$(COMMIT)


.PHONY: image
