init: gen-itest gen-docker

export APP_NAME
export APP_DOCKER_LABEL

gen-app: $(BEDROCK)
	$(BEDROCK) --base $(BEDROCK_ROOT) generate app $(SRCROOT)

gen-itest: $(BEDROCK)
	$(BEDROCK) --base $(BEDROCK_ROOT) generate itest

gen-docker: $(BEDROCK)
	$(BEDROCK) --base $(BEDROCK_ROOT) generate docker

gen-config: $(BEDROCK)
	$(BEDROCK) --base $(BEDROCK_ROOT) generate config