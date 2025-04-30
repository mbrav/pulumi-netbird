PROJECT_NAME      := Pulumi NetBird Resource Provider

PACK              := netbird
PACKDIR           := sdk
PROJECT           := github.com/mbrav/pulumi-netbird

GOPATH            := $(shell go env GOPATH)
WORKING_DIR       := $(shell pwd)
EXAMPLES_DIR      := $(WORKING_DIR)/examples/yaml
TESTPARALLELISM   := 4

PROVIDER          := pulumi-resource-$(PACK)
PROVIDER_BIN      := $(WORKING_DIR)/bin/$(PROVIDER)
VERSION           ?= $(shell pulumictl get version)
PROVIDER_PATH     := provider
VERSION_PATH      := ${PROVIDER_PATH}.Version

OS                := $(shell uname)
SHELL             := /bin/bash
GO_TEST           := go test -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM}

.PHONY: help prepare ensure build only_build lint install install_provider_symlink \
        provider provider_symlink provider_debug go_sdk test_provider test_all \
        install_go_sdk pulumi_init up down

help: ## Show help
	@echo "Available Makefile targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

prepare:: ## Prepare renamed project structure
	@if test -z "${NAME}"; then echo "NAME not set"; exit 1; fi
	@if test -z "${REPOSITORY}"; then echo "REPOSITORY not set"; exit 1; fi
	@if test -z "${ORG}"; then echo "ORG not set"; exit 1; fi
	@if test ! -d "provider/cmd/pulumi-resource-${PACK}"; then "Project already prepared"; exit 1; fi # SED_SKIP
	mv "provider/cmd/pulumi-resource-${PACK}" provider/cmd/pulumi-resource-${NAME} # SED_SKIP
	if [[ "${OS}" != "Darwin" ]]; then \
		find . \( -path './.git' -o -path './sdk' \) -prune -o -not -name 'go.sum' -type f -exec sed -i '/SED_SKIP/!s,github.com/pulumi/pulumi-[x]yz,${REPOSITORY},g' {} \; &> /dev/null; \
		find . \( -path './.git' -o -path './sdk' \) -prune -o -not -name 'go.sum' -type f -exec sed -i '/SED_SKIP/!s/[xX]yz/${NAME}/g' {} \; &> /dev/null; \
		find . \( -path './.git' -o -path './sdk' \) -prune -o -not -name 'go.sum' -type f -exec sed -i '/SED_SKIP/!s/[aA]bc/${ORG}/g' {} \; &> /dev/null; \
	fi
	# In MacOS the -i parameter needs an empty string to execute in place.
	if [[ "${OS}" == "Darwin" ]]; then \
		find . \( -path './.git' -o -path './sdk' \) -prune -o -not -name 'go.sum' -type f -exec sed -i '' '/SED_SKIP/!s,github.com/pulumi/pulumi-[x]yz,${REPOSITORY},g' {} \; &> /dev/null; \
		find . \( -path './.git' -o -path './sdk' \) -prune -o -not -name 'go.sum' -type f -exec sed -i '' '/SED_SKIP/!s/[xX]yz/${NAME}/g' {} \; &> /dev/null; \
		find . \( -path './.git' -o -path './sdk' \) -prune -o -not -name 'go.sum' -type f -exec sed -i '' '/SED_SKIP/!s/[aA]bc/${ORG}/g' {} \; &> /dev/null; \
	fi

ensure: ## Ensure Go modules are tidy
	cd provider && go mod tidy
	cd sdk && go mod tidy
	cd tests && go mod tidy

provider: $(PROVIDER_BIN) ## Build provider binary
$(PROVIDER_BIN): $(shell find provider -name "*.go")
	go build -o $(PROVIDER_BIN) -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER)

	@echo "Installing local plugin to Pulumi plugin cache..."
	@mkdir -p ~/.pulumi/plugins/${PACK}/resource/${VERSION}
	@cp $(PROVIDER_BIN) ~/.pulumi/plugins/${PACK}/resource/${VERSION}/pulumi-resource-${PACK}

provider_debug: ## Build provider with debug flags
	(cd provider && go build -o $(PROVIDER_BIN) -gcflags="all=-N -l" -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

provider_symlink: ## Symlink provider to Pulumi plugin cache
	@mkdir -p ~/.pulumi/plugins/${PACK}/resource/${VERSION}
	@ln -sf $(PROVIDER_BIN) ~/.pulumi/plugins/${PACK}/resource/${VERSION}/pulumi-resource-${PACK}

go_sdk: $(PROVIDER_BIN) ## Generate Go SDK from provider binary
	rm -rf sdk/go
	pulumi package gen-sdk $(PROVIDER_BIN) --language go

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
	cd ${EXAMPLES_DIR} && \
	pulumi cancel --stack ${PACK}-dev || true && \
	(pulumi stack init ${PACK}-dev || pulumi stack select ${PACK}-dev) && \
	pulumi up -y


down: ## Destroy stack
	$(call pulumi_login) \
	cd ${EXAMPLES_DIR} && \
	pulumi stack select ${PACK}-dev && \
	pulumi destroy -y && \
	pulumi stack rm ${PACK}-dev -y

lint: ## Run Go linters
	for DIR in "provider" "sdk" "tests" ; do \
		pushd $$DIR && golangci-lint run -c ../.golangci.yml --timeout 10m && popd ; \
	done

install: ## Install provider into $GOPATH/bin
	cp $(PROVIDER_BIN) $(GOPATH)/bin

build: provider go_sdk ## Build provider binary and SDK

only_build: build ## Alias for build used by CI
