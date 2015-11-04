itest: itest.env itest.run

itest.run:
	TEST_HOST="http://$(SERVER):$(PORT)" go test $(APP_NAME)/itest

itest.env: itest.env.stop itest.env.start

itest.env.start:
	for n in $(APP_ITEST_ENV_ROOT)/*-secrets.yml $(APP_ITEST_ENV_ROOT)/*-controller.yml $(APP_ITEST_ENV_ROOT)/*-service.yml; do \
		cat $$n | kubectl create -f - ; \
	done
	-wait-for-pod.sh $(APP_NAME)

itest.env.stop:
	-kubectl delete all -lapp=$(APP_NAME)
	-kubectl delete secrets -lapp=$(APP_NAME)

distitest:
	docker run --rm --net=host \
	           -v $(SRCROOT):$(SRCROOT_D) \
 	           -w $(SRCROOT_D) \
	           -e DEV_UID=$(DOCKER_DEV_UID) \
	           -e DEV_GID=$(DOCKER_DEV_GID) \
	           $(DOCKER_DEVIMAGE) \
						 make itest

distitest.env:
	docker run --rm --net=host \
	           -v $(SRCROOT):$(SRCROOT_D) \
 	           -w $(SRCROOT_D) \
	           -e DEV_UID=$(DOCKER_DEV_UID) \
	           -e DEV_GID=$(DOCKER_DEV_GID) \
	           $(DOCKER_DEVIMAGE) \
						 make itest.env

distitest.run:
	docker run --rm --net=host \
	           -v $(SRCROOT):$(SRCROOT_D) \
 	           -w $(SRCROOT_D) \
	           -e DEV_UID=$(DOCKER_DEV_UID) \
	           -e DEV_GID=$(DOCKER_DEV_GID) \
	           $(DOCKER_DEVIMAGE) \
	           make itest.run
