init: gen-itest gen-docker

export APP_NAME
export APP_DOCKER_LABEL

gen-all: gen-app gen-itest gen-docker gen-config gen-api

gen-app: $(BEDROCK) $(SRCROOT)/glide.yaml
	$(BEDROCK) --base $(BEDROCK_ROOT) generate app $(SRCROOT)

gen-itest: $(BEDROCK)
	$(BEDROCK) --base $(BEDROCK_ROOT) generate itest

gen-docker: $(BEDROCK)
	$(BEDROCK) --base $(BEDROCK_ROOT) generate docker

gen-config: $(BEDROCK)
	$(BEDROCK) --base $(BEDROCK_ROOT) generate config

gen-api: $(BEDROCK)
	$(BEDROCK) --base $(BEDROCK_ROOT) generate api

gen-secret: $(BEDROCK)
	$(BEDROCK) --base $(BEDROCK_ROOT) generate secret $(APP_SECRETS_ROOT)

$(SRCROOT)/glide.yaml:
	cp $(BEDROCK_ROOT)/fixtures/glide.yaml $@

distgen-secret:
	docker run -it --rm --net=host \
	           $(DOCKER_OPTS) \
	           $(DOCKER_DEVIMAGE) \
	           make gen-secret