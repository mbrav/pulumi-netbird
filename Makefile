PACK              := netbird
PACKDIR           := sdk
PROJECT           := github.com/mbrav/pulumi-netbird

GOPATH            := $(shell go env GOPATH)
WORKING_DIR       := $(shell pwd)
EXAMPLES_DIR      := $(WORKING_DIR)/examples
TESTPARALLELISM   := 4

PROVIDER          := pulumi-resource-$(PACK)
PROVIDER_BIN      := $(WORKING_DIR)/bin/$(PROVIDER)
VERSION           ?= $(shell pulumictl get version)
PROVIDER_PATH     := provider
VERSION_PATH      := ${PROVIDER_PATH}.Version
SCHEMA_PATH       := $(WORKING_DIR)/provider/cmd/pulumi-resource-${PACK}/schema.json

OS                := $(shell uname)
SHELL             := /bin/bash
GO_TEST           := go test -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM}

.PHONY: help ensure cov build only_build lint install install_provider_symlink \
        provider provider_symlink provider_debug go_sdk test_provider test_all \
        install_go_sdk pulumi_init up down

help: ## Show help
	@echo "Available Makefile targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

print: ## Print variable values
	@echo "PACK: ${PACK}"
	@echo "PACKDIR: ${PACKDIR}"
	@echo "PROJECT: ${PROJECT}"
	@echo "GOPATH: ${GOPATH}"
	@echo "WORKING_DIR: ${WORKING_DIR}"
	@echo "EXAMPLES_DIR: ${EXAMPLES_DIR}"
	@echo "PROVIDER_BIN: ${PROVIDER_BIN}"
	@echo "VERSION: ${VERSION}"
	@echo "OS: ${OS}"

ensure: ## Ensure Go modules are tidy
	cd provider && go mod tidy
	cd sdk && go mod tidy
	cd tests && go mod tidy
	cd examples/go && go mod tidy

ensure-update: ## Ensure Go modules are tidy and update
	cd provider && go get -u ./... && go mod tidy
	cd sdk && go get -u ./... && go mod tidy
	cd tests && go get -u ./... && go mod tidy
	cd examples/go && go get -u ./... && go mod tidy

cov: ## Generate coverage report
	rm ./coverage.out || true
	go test -coverprofile=coverage.out ./...
	go tool cover -func coverage.out

cov-html: cov ## View coverage HTML report
	go tool cover -html=coverage.out

schema: $(PROVIDER_BIN) ## Generate schema.json from provider binary
	@echo "Generating schema.json..."
	pulumi package get-schema $(PROVIDER_BIN) > $(SCHEMA_PATH)
	@echo "Wrote schema.json to $(SCHEMA_PATH)"

provider: $(PROVIDER_BIN) ## Build provider binary
$(PROVIDER_BIN): $(shell find provider -name "*.go")
	go build -o $(PROVIDER_BIN) -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER)
	@echo "Generating schema to $(SCHEMA_PATH)"
	pulumi package get-schema $(PROVIDER_BIN) > $(SCHEMA_PATH)
	@echo "Installing local plugin to Pulumi plugin cache..."
	@mkdir -p ~/.pulumi/plugins/${PACK}/resource/${VERSION}
	@cp $(PROVIDER_BIN) ~/.pulumi/plugins/${PACK}/resource/${VERSION}/pulumi-resource-${PACK}

provider_debug: ## Build provider with debug flags
	(cd provider && go build -o $(PROVIDER_BIN) -gcflags="all=-N -l" -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

provider_symlink: ## Symlink provider to Pulumi plugin cache
	@mkdir -p ~/.pulumi/plugins/${PACK}/resource/${VERSION}
	@ln -sf $(PROVIDER_BIN) ~/.pulumi/plugins/${PACK}/resource/${VERSION}/pulumi-resource-${PACK}

sdk_go: $(PROVIDER_BIN) ## Generate Go SDK from provider binary
	rm -rf sdk/go
	pulumi package gen-sdk $(PROVIDER_BIN) --language go

sdk_python: PYPI_VERSION := $(shell pulumictl get version --language python)
sdk_python: $(PROVIDER_BIN) ## Generate Python SDK from provider binary
	rm -rf sdk/python
	@echo "Generating SDK Pversion $(PYPI_VERSION)"
	pulumi package gen-sdk $(PROVIDER_BIN) --language python
	cp README.md sdk/python/

test_provider: ## Run provider tests
	cd tests && $(GO_TEST) ./...

test_all: test_provider ## Run all tests
	cd tests/sdk/go && $(GO_TEST) ./...


define pulumi_login
    export PULUMI_CONFIG_PASSPHRASE=test; \
    pulumi login --local;
endef

up: ## Deploy stack
	$(call pulumi_login) \
	cd ${EXAMPLES_DIR}/yaml && \
	pulumi cancel --stack ${PACK}-dev --yes >/dev/null 2>&1 || true && \
	(pulumi stack init ${PACK}-dev || pulumi stack select ${PACK}-dev) && \
	pulumi up --yes -d -v 3

refresh: ## Refresh the stack state from the actual resources
	$(call pulumi_login) \
	cd ${EXAMPLES_DIR}/yaml && \
	pulumi cancel --stack ${PACK}-dev --yes >/dev/null 2>&1 || true && \
	(pulumi stack init ${PACK}-dev || pulumi stack select ${PACK}-dev) && \
	pulumi refresh --yes -d -v 3

plan: ## Preview stack changes without applying
	$(call pulumi_login) \
	cd ${EXAMPLES_DIR}/yaml && \
	pulumi cancel --stack ${PACK}-dev --yes >/dev/null 2>&1 || true && \
	(pulumi stack init ${PACK}-dev || pulumi stack select ${PACK}-dev) && \
	pulumi preview -d -v 3

down: ## Destroy stack
	$(call pulumi_login)/yaml \
	cd ${EXAMPLES_DIR} && \
	pulumi stack select ${PACK}-dev && \
	pulumi destroy --yes && \
	pulumi stack rm ${PACK}-dev --yes

lint: ## Run Go linters
	GOFLAGS=-buildvcs=false golangci-lint run -c ./.golangci.yml

lint-fix: ## Run Go linters and fix
	GOFLAGS=-buildvcs=false golangci-lint run -c ./.golangci.yml --fix

build: provider sdk_go ## Build provider binary and SDK


cross_build: ## Build binaries for multiple OS/ARCH targets
	@mkdir -p dist
	@mkdir -p bin
	@rm -rf dist/* bin/*
	@for os in linux darwin; do \
		for arch in amd64 arm64; do \
			BIN_NAME=pulumi-resource-${PACK}-v${VERSION}-$$os-$$arch; \
			OUT_DIR=dist/$$BIN_NAME; \
			echo "Building $$BIN_NAME..."; \
			mkdir -p $$OUT_DIR; \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 \
			go build -o $$OUT_DIR/pulumi-resource-${PACK} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" \
				$(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER); \
			cp README.md $$OUT_DIR/; \
			cp LICENSE $$OUT_DIR/; \
			tar -czf dist/$$BIN_NAME.tar.gz -C $$OUT_DIR .; \
			rm -rf $$OUT_DIR; \
		done; \
	done
# echo "# Generated on $$(date -u +\"%Y-%m-%dT%H:%M:%SZ\")" > $$OUT_DIR/PulumiPlugin.yaml; \
# echo "name: ${PACK}" >> $$OUT_DIR/PulumiPlugin.yaml; \
# echo "kind: resource" >> $$OUT_DIR/PulumiPlugin.yaml; \
# echo "version: ${VERSION}" >> $$OUT_DIR/PulumiPlugin.yaml; \
# echo "runtime: go" >> $$OUT_DIR/PulumiPlugin.yaml; \

install: build ## Install provider into $GOPATH/bin
	cp $(PROVIDER_BIN) $(GOPATH)/bin
	pulumi plugin rm resource $$PACK -y || true
	pulumi plugin install resource $(PACK) $(VERSION) -f $(PROVIDER_BIN) || exit 1
	@echo "✅ Installed plugin $(PACK)@$(VERSION)"

only_build: build ## Alias for build used by CI
