# Pulumi NetBird Native Provider

<p align="center">
    <a href="https://github.com/mbrav/pulumi-netbird" target="_blank" rel="noopener noreferrer">
        <img width="100" src="./assets/logo.webp" title="pulumi-netbird"">
    </a>
</p>

[![Go Report Card](https://goreportcard.com/badge/github.com/mbrav/pulumi-netbird)](https://goreportcard.com/report/github.com/mbrav/pulumi-netbird)

[NetBird](https://github.com/netbirdio/netbird) is a modern, WireGuard-based mesh VPN. This provider integrates NetBird into Pulumi for seamless infrastructure automation.

This repository contains the **Pulumi NetBird Provider**, a native Pulumi provider built in Go using the [`pulumi-go-provider`](https://github.com/pulumi/pulumi-go-provider) SDK. It enables you to manage **NetBird** resources—like networks, peers, groups, and access rules—declaratively using Pulumi's infrastructure-as-code framework.

## ✨ Features

- Manage 15 NetBird resource types declaratively using Pulumi (Go, Python, YAML, TypeScript, C#)
- Built natively with Pulumi's Go SDK
- Works with NetBird Cloud (`https://api.netbird.io`) and self-hosted management servers

## 📦 Installing plugin

To install the Pulumi NetBird resource plugin, replace the version number with the desired release if needed. The plugin will be downloaded from the specified GitHub repository.

```bash
pulumi plugin install resource netbird 0.3.1 --server github://api.github.com/mbrav/pulumi-netbird
````

## 🧪 Build and Test

```bash
make help                 # View available build/test commands
````

## 🚀 Example Usage with Pulumi YAML

You can use this provider with **Pulumi YAML** to manage NetBird infrastructure declaratively.

### 1. Setup

Install the plugin (required before first `pulumi up`):

```bash
pulumi plugin install resource netbird 0.3.1 --server github://api.github.com/mbrav/pulumi-netbird
```

> **Note:** `--server` is a CLI-only flag. Do **not** add it to `Pulumi.yaml` — the `plugins.providers` block only accepts `name`, `path` (local binary), and `version`. For GitHub-hosted plugins, the CLI install above is sufficient.

Navigate to the YAML example directory:

```bash
cd examples/yaml
```

Initialize a new stack. If you are using the **local file backend** (`pulumi login --local`), the organization is always the literal string `organization` — use a simple name or the `organization/<project>/<stack>` form:

```bash
# Pulumi Cloud
pulumi stack init myorg/myproject/dev

# Local backend
pulumi stack init dev
# or fully qualified:
pulumi stack init organization/myproject/dev
```

Configure your credentials. Always use `--secret` for the token so it is encrypted in the stack config file:

```bash
pulumi config set --secret netbird:token YOUR_TOKEN
pulumi config set netbird:url https://nb.domain:33073
```

### 2. Deploy

```bash
pulumi up
```

This deploys a sample NetBird environment with networks, groups, network resources, a router, and a policy.

### 3. Sync existing resources

If a resource already exists in NetBird and you add it to `Pulumi.yaml`, import it first to avoid creating a duplicate:

```bash
# Get the resource ID from the NetBird UI or API, then:
pulumi import netbird:resource:Group group-admin <GROUP_ID>
```

After import, `pulumi up` will converge any property differences between your YAML and the live state.

For ongoing drift detection:

```bash
pulumi refresh   # pull live NetBird state into Pulumi state (detects out-of-band changes)
pulumi preview   # show what pulumi up would change
pulumi up        # apply
```

### Example `Pulumi.yaml` (published plugin)

When using the released plugin installed via `pulumi plugin install`, no `plugins:` block is needed. Use typed `config` entries so the token is encrypted in the stack config file:

```yaml
name: netbird
description: NetBird infrastructure managed via Pulumi
runtime: yaml

config:
  netbird:token:
    type: string
    secret: true
  netbird:url:
    type: string
    default: https://api.netbird.io

resources:
  group-admin:
    type: netbird:resource:Group
    properties:
      name: Admin
      peers: []

outputs: {}
```

### Example `Pulumi.yaml` (local dev build)

When developing the provider locally, point `plugins.providers` at the compiled binary. The `path` field is the only supported alternative to CLI-installed plugins — `server` is not a valid key here:

```yaml
name: provider-netbird
runtime: yaml
plugins:
  providers:
    - name: netbird
      path: ../../bin   # path to locally compiled binary

config:
  netbird:token:
    type: string
    secret: true
  netbird:url:
    type: string
    default: https://api.netbird.io

outputs:
  networkR1:
    value:
      name: ${net-r1.name}
      id: ${net-r1.id}

resources:
  group-devops:
    type: netbird:resource:Group
    properties:
      name: DevOps
      peers: []

  group-dev:
    type: netbird:resource:Group
    properties:
      name: Dev
      peers: []

  group-backoffice:
    type: netbird:resource:Group
    properties:
      name: Backoffice
      peers: []

  group-hr:
    type: netbird:resource:Group
    properties:
      name: HR
      peers: []

  net-r1:
    type: netbird:resource:Network
    properties:
      name: R1
      description: Network for Region 1

  netres-r1-net-01:
    type: netbird:resource:NetworkResource
    properties:
      name: Region 1 Net 01
      description: Network 01 in Region 1
      networkID: ${net-r1.id}
      address: 10.10.1.0/24
      enabled: true
      groupIDs:
        - ${group-devops.id}

  netres-r1-net-02:
    type: netbird:resource:NetworkResource
    properties:
      name: Region 1 Net 02
      description: Network 02 in S1 Region 1
      networkID: ${net-r1.id}
      address: 10.10.2.0/24
      enabled: true
      groupIDs:
        - ${group-devops.id}

  netres-r1-net-03:
    type: netbird:resource:NetworkResource
    properties:
      name: Region 1 Net 03
      description: Network 03 in Region 1
      networkID: ${net-r1.id}
      address: 10.10.3.0/24
      enabled: true
      groupIDs:
        - ${group-devops.id}

  router-r1:
    type: netbird:resource:NetworkRouter
    properties:
      networkID: ${net-r1.id}
      enabled: true
      masquerade: true
      metric: 10
      peer: ""
      peerGroups:
        - ${group-devops.id}

  policy-ssh-grp-src-net-dest:
    type: netbird:resource:Policy
    properties:
      name: "SSH Policy - Group to Subnet"
      description: "Allow SSH (22/TCP) from DevOps and Dev groups to Region 1 Net 02"
      enabled: true
      postureChecks: []
      rules:
        - name: "SSH Access - Group → Subnet"
          description: "Allow unidirectional SSH from DevOps & Dev groups to Net 02"
          bidirectional: false
          action: accept
          enabled: true
          protocol: tcp
          ports:
            - "22"
          sources:
            - ${group-devops.id}
            - ${group-dev.id}
          destinationResource:
            type: subnet
            id: ${netres-r1-net-02.id}

  policy-ssh-grp-src-grp-dest:
    type: netbird:resource:Policy
    properties:
      name: "SSH Policy - Group to Group"
      description: "Allow SSH (22/TCP) from DevOps to Backoffice group resources"
      enabled: true
      postureChecks: []
      rules:
        - name: "SSH Access - Group → Group"
          description: "SSH from DevOps group to Backoffice group"
          bidirectional: false
          action: accept
          enabled: true
          protocol: tcp
          ports:
            - "22"
          sources:
            - ${group-devops.id}
          destinations:
            - ${group-backoffice.id}

```

## 🦫 Example Usage with Pulumi Go

You can use this provider with **Pulumi Go** to manage NetBird infrastructure declaratively.

The SDK is accessible through the generated `github.com/mbrav/pulumi-netbird/sdk/go/netbird` module.

SDK versions are available to Go with tags that are prefixed with `sdk/vx.x.x` and can be listed with the following command:

```bash
go list -m -versions github.com/mbrav/pulumi-netbird/sdk
```

Output:

```bash
github.com/mbrav/pulumi-netbird/sdk v0.3.0 v0.3.1 # and so on
```

### 1. Setup

Navigate to the Go example directory:

```bash
cd examples/go
```

Initialize a new stack and configure your credentials:

```bash
pulumi stack init test
pulumi config set --secret netbird:token YOUR_TOKEN
pulumi config set netbird:url https://nb.domain:33073
```

### 2. Deploy

```bash
pulumi up
```

## 🐍 Example Usage with Pulumi Python

You can use this provider with **Pulumi Python** to manage NetBird infrastructure declaratively.

### 1. Setup

First, you must generate the python SDK:

```bash
make provider
make sdk_python
```

Then install the wheel:

```bash
pip install sdk/python/bin/dist/pulumi_netbird-0.3.1.tar.gz
```

Navigate to the Python example directory:

```bash
cd examples/python
```

Initialize a new stack and configure your credentials:

```bash
pulumi stack init test
pulumi config set --secret netbird:token YOUR_TOKEN
pulumi config set netbird:url https://nb.domain:33073
```

### 2. Deploy

```bash
pulumi up
```

## 📁 Repository Structure

- `provider/` – Go implementation of the provider
- `sdk/go/netbird/` – Go SDK for the NetBird provider
- `examples/` – Example Pulumi projects using the provider

## 📚 References

- [Pulumi Go Provider Docs](https://github.com/pulumi/pulumi-go-provider)
- [NetBird Documentation](https://docs.netbird.io/)
- [Pulumi YAML Documentation](https://www.pulumi.com/docs/using-pulumi/yaml/)
