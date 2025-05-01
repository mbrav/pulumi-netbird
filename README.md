# Pulumi NetBird Native Provider

![netbird_logo](./assets/logo.webp)

## Project still in WIP

This repository contains the **Pulumi NetBird Provider**, a native Pulumi provider built in Go using the [`pulumi-go-provider`](https://github.com/pulumi/pulumi-go-provider) SDK. It enables you to manage **NetBird** resources—like networks, peers, groups, and access rules—declaratively using Pulumi's infrastructure-as-code framework.

NetBird is a modern, WireGuard-based mesh VPN. This provider integrates NetBird into Pulumi for seamless infrastructure automation.

## Features

* Manage NetBird resources using Pulumi in Go
* Built natively with Pulumi's Go SDK
* Includes example configurations for local testing

## Prerequisites

Ensure the following are installed and available in your `$PATH`:

* [Go 1.24+](https://go.dev/dl/)
* [`pulumictl`](https://github.com/pulumi/pulumictl#installation)
* [Pulumi CLI](https://www.pulumi.com/docs/get-started/install/)

## Build and Test

```bash
make build install        # Build and install the provider locally
make gen_examples         # Generate Go examples from YAML source
make up                   # Deploy the example Pulumi stack
make down                 # Destroy the example stack
```

## Creating Your Own Provider (Optional)

This repository was based on Pulumi’s native provider boilerplate. If you're building your own provider:

```bash
make prepare NAME=yourprovider ORG=yourorg REPOSITORY=github.com/yourorg/pulumi-yourprovider
```

This command will:

* Rename relevant files and folders
* Update Go module paths
* Replace placeholder names

## Example Usage

Navigate to the example directory:

```bash
cd examples/simple
pulumi stack init test
pulumi up
```

This deploys a sample NetBird configuration using the provider.

## Repository Structure

* `provider/` – Go implementation of the provider
* `sdk/go/netbird/` – Go SDK for the NetBird provider
* `examples/` – Example Pulumi projects using the provider
* `Makefile` – Task runner for build, install, and test operations

## References

* [Pulumi Go Provider Docs](https://github.com/pulumi/pulumi-go-provider)
* [NetBird Documentation](https://docs.netbird.io/)
* [Pulumi Command Provider (example implementation)](https://github.com/pulumi/pulumi-command)
