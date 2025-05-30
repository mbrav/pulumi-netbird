# Pulumi NetBird Native Provider

<p align="center">
    <a href="https://github.com/mbrav/pulumi-netbird" target="_blank" rel="noopener noreferrer">
        <img width="100" src="./assets/logo.webp" title="pulumi-netbird"">
    </a>
</p>

[NetBird](https://github.com/netbirdio/netbird) is a modern, WireGuard-based mesh VPN. This provider integrates NetBird into Pulumi for seamless infrastructure automation.

This repository contains the **Pulumi NetBird Provider**, a native Pulumi provider built in Go using the [`pulumi-go-provider`](https://github.com/pulumi/pulumi-go-provider) SDK. It enables you to manage **NetBird** resources—like networks, peers, groups, and access rules—declaratively using Pulumi's infrastructure-as-code framework.

## ✨ Features

- Manage NetBird resources using Pulumi in Go or YAML
- Built natively with Pulumi's Go SDK
- Includes example configurations for local testing

## 📦 Installing plugin

To manually install the Pulumi NetBird resource plugin replace the version number (`0.0.19`) with the desired release if needed. The plugin will be downloaded from the specified GitHub repository.

```bash
pulumi plugin install resource netbird 0.0.19 --server github://api.github.com/mbrav/pulumi-netbird
````

## 🧪 Build and Test

```bash
make help                 # View available build/test commands
````

## 🚀 Example Usage with Pulumi YAML

You can use this provider with **Pulumi YAML** to manage NetBird infrastructure declaratively.

### 1. Setup

Navigate to the YAML example directory:

```bash
cd examples/yaml
```

Initialize a new stack and configure your credentials:

```bash
pulumi stack init test
pulumi config set netbird:netbirdToken YOUR_TOKEN
pulumi config set netbird:netbirdUrl https://nb.domain:33073
```

### 2. Deploy

```bash
pulumi up
```

This deploys a sample NetBird environment with networks, groups, network resources, a router, and a policy.

### Example `Pulumi.yaml`

```yaml
name: provider-netbird
runtime: yaml
plugins:
  providers:
    - name: netbird
      path: ../../bin

# You can also define creds here
config:
  netbird:netbirdToken: token
  netbird:netbirdUrl: https://nb.domain:33073

outputs:
  networkManagement:
    value:
      name: ${net-management.name}
      id: ${net-management.Id}

resources:
  net-management:
    type: netbird:resource:Network
    properties:
      name: Management
      description: Network for Management

  net-r1:
    type: netbird:resource:Network
    properties:
      name: R1
      description: Network for Region 1

  group-devops:
    type: netbird:resource:Group
    properties:
      name: DevOps
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

  netres-r1-net-01:
    type: netbird:resource:NetworkResource
    properties:
      name: S1 Franfurt Net 01
      description: Network 01 in S1 Franfurt
      network_id: ${net-r1.id}
      address: 10.10.1.0/24
      enabled: true
      group_ids:
        - ${group-devops.id}

  netres-r1-net-02:
    type: netbird:resource:NetworkResource
    properties:
      name: S1 Franfurt Net 02
      description: Network 02 in S1 Franfurt
      network_id: ${net-r1.id}
      address: 10.10.2.0/24
      enabled: true
      group_ids:
        - ${group-devops.id}

  netres-r1-net-03:
    type: netbird:resource:NetworkResource
    properties:
      name: S1 Franfurt Net 03
      description: Network 03 in S1 Franfurt
      network_id: ${net-r1.id}
      address: 10.10.3.0/24
      enabled: true
      group_ids:
        - ${group-devops.id}

  router-r1:
    type: netbird:resource:NetworkRouter
    properties:
      network_id: ${net-r1.id}
      enabled: true
      masquerade: true
      metric: 10
      peer: ""
      peer_groups:
        - ${group-devops.id}

  test-import-peer:
    type: netbird:resource:Peer
    properties:
      inactivity_expiration_enabled: false
      login_expiration_enabled: false
      name: test-import-peer
      sshEnabled: true
    options:
      protect: true

  test-ssh-policy:
    type: netbird:resource:Policy
    properties:
      name: "Test SSH Policy"
      description: "Allow SSH access from admin group to servers"
      enabled: true
      posture_checks: []
      rules:
        - name: "SSH Access"
          description: "Allow inbound SSH from Admins to Servers"
          bidirectional: false
          action: accept
          enabled: true
          protocol: tcp
          ports:
            - "22"
          # sources:
          #   - "group-admins"
          # destinations:
          #   - "group-servers"
```

## 🦫 Example Usage with Pulumi Go

You can use this provider with **Pulumi Go** to manage NetBird infrastructure declaratively.

The SDK is accessible through the generated `github.com/mbrav/pulumi-netbird/sdk/go/netbird` module.

SDK versions are available to Go with tags that are prefixed with `sdk/vx.x.x` and can be listed with the following command:

```bash
go list -m -versions github.com/mbrav/pulumi-netbird/sdk
```

Output:

```text
github.com/mbrav/pulumi-netbird/sdk v0.0.11
```

### 1. Setup

Navigate to the Go example directory:

```bash
cd examples/go
```

Initialize a new stack and configure your credentials:

```bash
pulumi stack init test
pulumi config set netbird:netbirdToken YOUR_TOKEN
pulumi config set netbird:netbirdUrl https://nb.domain:33073
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
pip install sdk/python/bin/dist/pulumi_netbird-0.0.12a1747213794+dirty.tar.gz
```

Navigate to the Python example directory:

```bash
cd examples/python
```

Initialize a new stack and configure your credentials:

```bash
pulumi stack init test
pulumi config set netbird:netbirdToken YOUR_TOKEN
pulumi config set netbird:netbirdUrl https://nb.domain:33073
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
