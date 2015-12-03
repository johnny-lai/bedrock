KUBE_SECRETS = $(wildcard $(APP_ITEST_ENV_ROOT)/*-secret.yml)
KUBE_CONTROLLERS = $(wildcard $(APP_ITEST_ENV_ROOT)/*-controller.yml)
KUBE_SERVICES = $(wildcard $(APP_ITEST_ENV_ROOT)/*-service.yml)

export APP_NAME
export VERSION
export COMMIT
export APP_DOCKER_LABEL
export APP_SECRETS_ROOT

export PO_APP_NAME = $(APP_NAME)
export PO_APP_IMAGE = $(APP_DOCKER_LABEL_COMMIT)
export SVC_APP_NAME = $(PO_APP_NAME)

export PO_DB_NAME = $(APP_NAME)-db
export PO_DB_IMAGE = $(TESTDB_DOCKER_LABEL_COMMIT)
export SVC_DB_NAME = $(PO_DB_NAME)

export SECRET_DB_NAME = $(APP_NAME)-db-secret

itest: itest.env itest.run

itest.run:
	. $(CLUSTER_SH) env $(SVC_APP_NAME) && \
	  go test $(APP_NAME)/itest

itest.env: itest.env.stop itest.env.start

itest.env.start: $(BEDROCK)
	for n in $(KUBE_SECRETS); do \
		$(BEDROCK) dump $$n | kubectl create -f - ; \
	done
	for n in $(KUBE_CONTROLLERS); do \
		$(BEDROCK) dump $$n | kubectl create -f - ; \
	done
	for n in $(KUBE_SERVICES); do \
		$(BEDROCK) dump $$n | kubectl create -f - ; \
	done
	-$(CLUSTER_SH) wait $(SVC_APP_NAME)

itest.env.stop:
	-kubectl delete all -lapp=$(APP_NAME)
	-kubectl delete secrets -lapp=$(APP_NAME)

distitest:
	docker run --rm --net=host \
	           $(DOCKER_OPTS) \
	           $(DOCKER_DEVIMAGE) \
						 make itest

distitest.env:
	docker run --rm --net=host \
	           $(DOCKER_OPTS) \
	           $(DOCKER_DEVIMAGE) \
						 make itest.env

distitest.run:
	docker run --rm --net=host \
	           $(DOCKER_OPTS) \
	           $(DOCKER_DEVIMAGE) \
	           make itest.run
