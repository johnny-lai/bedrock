init: gen-itest gen-docker

export APP_NAME
export APP_DOCKER_LABEL

gen-all: gen-app gen-itest gen-docker gen-config gen-api

gen-app: $(BEDROCK) $(SRCROOT)/glide.yaml
	$(RUN_BEDROCK_D) generate app $(SRCROOT_D)

gen-itest: $(BEDROCK)
	$(RUN_BEDROCK_D) generate itest

gen-docker: $(BEDROCK)
	$(RUN_BEDROCK_D) generate docker

gen-config: $(BEDROCK)
	$(RUN_BEDROCK_D) generate config

gen-api: $(BEDROCK)
	$(RUN_BEDROCK_D) generate api

gen-secret: $(BEDROCK)
	mkdir -p $(APP_SECRETS_ROOT)
	$(RUN_BEDROCK_D) generate secret $(APP_SECRETS_ROOT_D)

$(SRCROOT)/glide.yaml:
	cp $(BEDROCK_ROOT)/fixtures/glide.yaml $@
