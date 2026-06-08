# Examples

This directory contains example Pulumi projects that demonstrate different ways to use the NetBird provider.

## Overview

| Example | Runtime | Description |
|---------|---------|-------------|
| [`yaml`](./yaml/) | Pulumi YAML | All resources in a single `Pulumi.yaml` — good starting point |
| [`yaml-yq`](./yaml-yq/) | Pulumi YAML + `yq` | Resources split across multiple `src/*.yaml` files, assembled by `make build` |
| [`go`](./go/) | Pulumi Go | Provider usage via the generated Go SDK |
| [`python`](./python/) | Pulumi Python | Provider usage via the generated Python SDK |

---

## `yaml` — Single-file YAML

A single `Pulumi.yaml` containing all resources: groups, networks, network resources, a router, policies, DNS, routes, setup keys, and a reverse proxy service.

**Prerequisites:** `pulumi` CLI, provider plugin installed.

```bash
pulumi plugin install resource netbird 0.3.8 --server github://api.github.com/mbrav/pulumi-netbird

cd yaml
pulumi stack init dev
pulumi config set --secret netbird:token YOUR_TOKEN
pulumi config set netbird:url https://nb.example.com:33073
pulumi up
```

## `yaml-yq` — Multi-file YAML with `yq`

Resources are split into focused files under `src/` and deep-merged into a single `Pulumi.yaml` at build time using [`yq`](https://github.com/mikefarah/yq). The `Makefile` drives the workflow.

```
yaml-yq/
  Pulumi.base.yaml   # stack metadata and config schema
  Makefile           # build / preview / up / destroy targets
  src/
    peers.yaml       # Peer resources (imported only)
    groups.yaml      # Group resources
    policies.yaml    # Policy resources
    routes.yaml      # Route resources
```

**How the build works:**

```bash
yq eval-all '. as $item ireduce ({}; . * $item)' Pulumi.base.yaml src/*.yaml > Pulumi.yaml
```

`make preview` and `make up` both run `make build` first, so `Pulumi.yaml` is always up to date. **Never edit `Pulumi.yaml` directly** — it is overwritten on every build.

**Prerequisites:** `pulumi` CLI, `yq` v4+, provider plugin installed.

```bash
# Install yq (macOS)
brew install yq

# Install the provider plugin
cd yaml-yq
make setup

pulumi stack init dev
pulumi config set --secret netbird:token YOUR_TOKEN
pulumi config set netbird:url https://nb.example.com:33073

make preview   # build + dry-run
make up        # build + apply
```

**When to use this pattern:**

- Many resources that benefit from separation by type (peers, groups, policies, routes)
- Team environments where different people own different resource files
- Large stacks where a single `Pulumi.yaml` becomes unwieldy

## `go` — Go SDK

Uses the generated `github.com/mbrav/pulumi-netbird/sdk/go/netbird` module to manage NetBird resources from a Go Pulumi program.

**Prerequisites:** Go 1.21+, `pulumi` CLI.

```bash
cd go
pulumi stack init dev
pulumi config set --secret netbird:token YOUR_TOKEN
pulumi config set netbird:url https://nb.example.com:33073
pulumi up
```

To list available SDK versions:

```bash
go list -m -versions github.com/mbrav/pulumi-netbird/sdk
```

## `python` — Python SDK

Uses the generated `pulumi_netbird` Python package to manage NetBird resources from a Python Pulumi program.

**Prerequisites:** Python 3.9+, `pulumi` CLI, generated Python SDK wheel.

Build and install the SDK first (from the repo root):

```bash
make provider
make sdk_python
pip install sdk/python/bin/dist/pulumi_netbird-*.tar.gz
```

Then run the example:

```bash
cd python
pulumi stack init dev
pulumi config set --secret netbird:token YOUR_TOKEN
pulumi config set netbird:url https://nb.example.com:33073
pulumi up
```
