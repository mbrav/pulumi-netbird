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

- Manage 16 NetBird resource types declaratively using Pulumi (Go, Python, YAML, TypeScript, C#)
- 6 read-only **invoke functions** (data sources) for referencing existing NetBird objects by name, email, or CIDR
- Built natively with Pulumi's Go SDK
- Works with NetBird Cloud (`https://api.netbird.io`) and self-hosted management servers

## 📦 Installing plugin

To install the Pulumi NetBird resource plugin, replace the version number with the desired release if needed. The plugin will be downloaded from the specified GitHub repository.

```bash
pulumi plugin install resource netbird 0.5.1 --server github://api.github.com/mbrav/pulumi-netbird
````

## 🧪 Build and Test

```bash
make help                 # View available build/test commands
````

## 🗂️ Examples

All runnable examples live in [`examples/`](./examples/README.md). The table below summarises what is available:

| Example | Runtime | Description |
|---------|---------|-------------|
| [`yaml`](./examples/yaml/) | Pulumi YAML | All resources in a single `Pulumi.yaml` |
| [`yaml-yq`](./examples/yaml-yq/) | Pulumi YAML + `yq` | Resources split across `src/*.yaml`, assembled by `make build` |
| [`go`](./examples/go/) | Pulumi Go | Provider usage via the generated Go SDK |
| [`python`](./examples/python/) | Pulumi Python | Provider usage via the generated Python SDK |

See **[examples/README.md](./examples/README.md)** for setup instructions for each example.

## 🚀 Example Usage with Pulumi YAML

You can use this provider with **Pulumi YAML** to manage NetBird infrastructure declaratively.

### 1. Setup

Install the plugin (required before first `pulumi up`):

```bash
pulumi plugin install resource netbird 0.5.1 --server github://api.github.com/mbrav/pulumi-netbird
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

### 3. Import existing resources

If resources already exist in NetBird, import them before running `pulumi up` so Pulumi adopts the live resources instead of creating duplicates.

1. Define the resource in your Pulumi program with the same logical name you will use in the import command.
2. Find the resource ID in the NetBird UI, API, or exported state from the tool that currently manages it.
3. Run `pulumi import <type> <name> <id>`.
4. Run `pulumi preview` and adjust the program until there are no unintended changes.
5. Run `pulumi up` to persist the reconciled inputs.

For most resources, the import ID is the NetBird resource ID:

```bash
pulumi import netbird:resource:Group group-admin <GROUP_ID>
pulumi import netbird:resource:Network net-r1 <NETWORK_ID>
pulumi import netbird:resource:Policy policy-admin <POLICY_ID>
pulumi import netbird:resource:Peer peer-mp1 <PEER_ID>
```

`NetworkRouter` and `NetworkResource` belong to a NetBird network, so their import IDs must include both the parent network ID and the child resource ID:

```bash
pulumi import netbird:resource:NetworkRouter router-r1 <NETWORK_ID>/<ROUTER_ID>
pulumi import netbird:resource:NetworkResource netres-r1-net-01 <NETWORK_ID>/<RESOURCE_ID>
```

Peers must be imported. They cannot be created through the NetBird management API, so `pulumi up` for a new `Peer` resource will fail unless the peer already exists in state. A minimal YAML declaration for an imported peer can look like this:

```yaml
resources:
  peer-mp1:
    type: netbird:resource:Peer
    properties:
      name: mp1
    options:
      protect: true
```

For policies, keep the intended `rules` in your Pulumi program after import. The provider can reconstruct rule inputs during import, but declaring the rules explicitly keeps future previews understandable and makes drift intentional.

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
      # peers and resources are optional. Omit them to let membership be
      # managed externally; declare them to manage membership from here.
      peers: []
      resources: []

  group-dev:
    type: netbird:resource:Group
    properties:
      name: Dev

  group-backoffice:
    type: netbird:resource:Group
    properties:
      name: Backoffice

  group-hr:
    type: netbird:resource:Group
    properties:
      name: HR

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
github.com/mbrav/pulumi-netbird/sdk v0.3.6 v0.3.7 v0.3.8 v0.4.1 v0.5.0 v0.5.1 # ... and so on
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
pip install sdk/python/bin/dist/pulumi_netbird-0.5.1.tar.gz
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

## 📋 Supported Resources

| Resource | Pulumi type |
| -------- | ----------- |
| DNS nameserver group | `netbird:resource:DNS` |
| DNS record | `netbird:resource:DNSRecord` |
| DNS settings | `netbird:resource:DNSSettings` |
| DNS zone | `netbird:resource:DNSZone` |
| Group | `netbird:resource:Group` |
| Network | `netbird:resource:Network` |
| Network resource | `netbird:resource:NetworkResource` |
| Network router | `netbird:resource:NetworkRouter` |
| Peer | `netbird:resource:Peer` |
| Policy | `netbird:resource:Policy` |
| Posture check | `netbird:resource:PostureCheck` |
| Reverse proxy domain | `netbird:resource:ReverseProxyDomain` |
| Reverse proxy service | `netbird:resource:ReverseProxyService` |
| Route | `netbird:resource:Route` |
| Setup key | `netbird:resource:SetupKey` |
| User | `netbird:resource:User` |

## 🔍 Invoke Functions (Data Sources)

Invoke functions are **read-only** — they query live NetBird state and return data without managing any resources. Use them to reference existing objects by a human-readable key rather than a hardcoded ID.

| Function | Pulumi type | Looks up by | Key output fields |
| -------- | ----------- | ----------- | ----------------- |
| Get peers | `netbird:function:getPeers` | optional group ID filter | `peers[]` (id, name, ip, connected, groups) |
| Lookup group | `netbird:function:lookupGroup` | group name | `groupId`, `peers[]`, `resources[]` |
| Lookup peer | `netbird:function:lookupPeer` | peer name | `peerId`, `ip`, `dnsLabel`, `connected`, `groups[]` |
| Lookup route | `netbird:function:lookupRoute` | network CIDR | `routeId`, `peerGroups[]`, `groups[]` |
| Lookup setup key | `netbird:function:lookupSetupKey` | key name | `setupKeyId`, `state`, `expires` |
| Lookup user | `netbird:function:lookupUser` | email address | `userId`, `role`, `autoGroups[]` |

### Example: cross-referencing an existing group in YAML

```yaml
variables:
  devopsGroup:
    fn::invoke:
      function: netbird:function:lookupGroup
      arguments:
        name: DevOps

resources:
  setup-key-devops:
    type: netbird:resource:SetupKey
    properties:
      name: devops-onboarding
      type: reusable
      expiresIn: 86400
      autoGroups:
        - ${devopsGroup.groupId}
```

### Example: cross-referencing in Go

```go
devopsGroup, err := netbird.LookupGroup(ctx, &netbird.LookupGroupArgs{
    Name: "DevOps",
}, nil)
if err != nil {
    return err
}

_, err = netbird.NewSetupKey(ctx, "setup-key-devops", &netbird.SetupKeyArgs{
    Name:       pulumi.String("devops-onboarding"),
    Type:       pulumi.String("reusable"),
    ExpiresIn:  pulumi.Int(86400),
    AutoGroups: pulumi.StringArray{pulumi.String(devopsGroup.GroupId)},
})
```

## 🧪 Experimental Components

> **These components are a proof of concept.** They explore the `pulumi-go-provider` component API and demonstrate how multiple resources can be bundled into a single declaration. The interface may change without notice. Do not rely on them in production.

Components are higher-level abstractions that create several NetBird resources together and wire them automatically. They appear in the schema under the `netbird:component:*` token prefix.

| Component | Pulumi type | Creates |
| --------- | ----------- | ------- |
| Network bundle | `netbird:component:NetworkBundle` | `Network` + `NetworkRouter` + N `NetworkResource` subnets |
| DNS zone bundle | `netbird:component:DNSZoneBundle` | `DNSZone` + N `DNSRecord`s |

### Example: NetworkBundle in YAML

```yaml
resources:
  r1:
    type: netbird:component:NetworkBundle
    properties:
      name: Region1
      router:
        enabled: true
        masquerade: true
        metric: 10
        peerGroups:
          - ${group-devops.id}
      subnets:
        - name: Net01
          address: 10.10.1.0/24
          enabled: true
          groupIDs:
            - ${group-devops.id}
        - name: Net02
          address: 10.10.2.0/24
          enabled: true
          groupIDs:
            - ${group-devops.id}

outputs:
  networkId: ${r1.networkId}
  subnetIds: ${r1.subnetIds}
```

### Example: DNSZoneBundle in YAML

```yaml
resources:
  corp-zone:
    type: netbird:component:DNSZoneBundle
    properties:
      name: corp-internal
      domain: corp.example.com
      enabled: true
      enableSearchDomain: true
      distributionGroups:
        - ${group-devops.id}
      records:
        - name: gw.corp.example.com
          type: A
          content: 10.10.1.1
          ttl: 300
        - name: api.corp.example.com
          type: CNAME
          content: gw.corp.example.com
          ttl: 300

outputs:
  zoneId: ${corp-zone.zoneId}
```

## 📁 Repository Structure

- `provider/` – Go implementation of the provider
- `sdk/go/netbird/` – Go SDK for the NetBird provider
- `examples/` – Example Pulumi projects using the provider

## 📚 References

- [Pulumi Go Provider Docs](https://github.com/pulumi/pulumi-go-provider)
- [NetBird Documentation](https://docs.netbird.io/)
- [Pulumi YAML Documentation](https://www.pulumi.com/docs/using-pulumi/yaml/)
