.PHONY: build preview up destroy setup

PROVIDER_NAME    := netbird
PROVIDER_VERSION := 0.3.8
PROVIDER_SERVER  := github://api.github.com/mbrav/pulumi-netbird
PROVIDER_CACHE   := $(HOME)/.pulumi/plugins/resource-$(PROVIDER_NAME)-v$(PROVIDER_VERSION)

setup:
	pulumi plugin install resource $(PROVIDER_NAME) $(PROVIDER_VERSION) --server $(PROVIDER_SERVER)

# Deep-merge Pulumi.base.yaml with every src/*.yaml file using yq.
# The resulting Pulumi.yaml is what Pulumi actually reads — never edit it directly.
build:
	yq eval-all '. as $$item ireduce ({}; . * $$item)' Pulumi.base.yaml src/*.yaml > Pulumi.yaml

preview: build
	pulumi preview

up: build
	pulumi up

destroy: build
	pulumi destroy
