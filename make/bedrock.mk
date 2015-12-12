BEDROCK = $(SRCROOT)/tmp/dist/bedrock
BEDROCK_D = $(SRCROOT_D)/tmp/dist/bedrock
		
BEDROCK_DOCKER_IMAGE = johnnylai/bedrock-dev-golang:1.5
BEDROCK_DOCKER_OPTS = $(DOCKER_OPTS) \
			  -v $(BEDROCK_ROOT):$(BEDROCK_ROOT_D) \
			  -e "GO15VENDOREXPERIMENT=1"

RUN_BEDROCK_D = $(DOCKER) run --rm $(BEDROCK_DOCKER_OPTS) $(BEDROCK_DOCKER_IMAGE) $(BEDROCK_D) --fixtures $(FIXTURES_ROOT_D)

$(BEDROCK): $(BEDROCK_ROOT)/cli/bedrock.go $(GO_BASE_DEPENDENCIES)
	$(DOCKER) run --rm \
	          $(BEDROCK_DOCKER_OPTS) \
			  $(BEDROCK_DOCKER_IMAGE) \
			  go build -o $(BEDROCK_D) $(BEDROCK_ROOT_D)/cli/bedrock.go
