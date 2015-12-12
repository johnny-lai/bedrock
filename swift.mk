DOCKER_DEVIMAGE ?= johnnylai/bedrock-dev-swift:2.2
FIXTURES_ROOT_D = $(BEDROCK_ROOT_D)/fixtures/swift

APP_SWIFT_LINKING ?= dynamic
APP_SWIFT_SOURCES ?=

$(APP): $(APP_SWIFT_SOURCES)
	swiftc -o $@ $(APP_SWIFT_SOURCES)

clean: clean.swift

clean.swift:
	echo "clean swift"