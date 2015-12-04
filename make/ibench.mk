ibench: ibench.env ibench.run

ibench.run:
	TEST_HOST="http://$(CLUSTER_SERVER):$(CLUSTER_PORT)" go test -bench=. $(APP_NAME)/itest

ibench.env: itest.env

distibench: distibench.env distibench.run

distibench.env: distitest.env

distibench.run:
	$(DOCKER) run --rm --net=host \
	           -v $(SRCROOT):$(SRCROOT_D) \
 	           -w $(SRCROOT_D) \
	           -e DEV_UID=$(DOCKER_DEV_UID) \
	           -e DEV_GID=$(DOCKER_DEV_GID) \
	           $(DOCKER_DEVIMAGE) \
	           make ibench.run