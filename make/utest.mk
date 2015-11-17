utest: deps
	TEST_CONFIG_YML=$(TEST_CONFIG_YML) GO15VENDOREXPERIMENT=1 go test $(APP_GO_PACKAGES)

distutest: distutest.env distutest.run

distutest.env:
	-docker rm -f $(APP_NAME)-testdb
	docker run -d --name $(APP_NAME)-testdb $(APP_DOCKER_LABEL)-testdb

distutest.run:
	docker run --rm \
	           --link $(APP_NAME)-testdb:$(APP_NAME)-db \
	           -v $(SRCROOT):$(SRCROOT_D) \
	           -v $(APP_SECRETS_ROOT):/etc/secrets \
	           -w $(SRCROOT_D) \
	           -e DEV_UID=$(DOCKER_DEV_UID) \
	           -e DEV_GID=$(DOCKER_DEV_GID) \
	           -e TEST_CONFIG_YML=$(TEST_CONFIG_YML_D) \
	           $(DOCKER_DEVIMAGE) \
	           make utest