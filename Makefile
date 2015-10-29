IMAGE_NAME = johnnylai/bedrock-dev

image:
	docker build -t $(IMAGE_NAME) -f docker/dev/Dockerfile docker/dev

deploy: image
	docker push $(IMAGE_NAME)

.PHONY: image
