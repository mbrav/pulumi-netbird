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

- Manage 28 NetBird resource types declaratively using Pulumi (Go, Python, YAML, TypeScript, C#)
- 10 read-only **invoke functions** (data sources) for referencing existing NetBird objects by name, email, CIDR, or country
- Built natively with Pulumi's Go SDK
- Works with NetBird Cloud (`https://api.netbird.io`) and self-hosted management servers

## 📦 Installing plugin

To install the Pulumi NetBird resource plugin, replace the version number with the desired release if needed. The plugin will be downloaded from the specified GitHub repository.

```bash
pulumi plugin install resource netbird 0.5.11 --server github://api.github.com/mbrav/pulumi-netbird
````

## 🧪 Build and Test

```bash
make help                 # View available build/test commands
````

## 🗂️ Examples

All runnable examples live in [`examples/`](./examples/README.md). The table below summarises what is available:

| Example | Runtime | Description |
| --------- | --------- | ------------- |
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
pulumi plugin install resource netbird 0.5.11 --server github://api.github.com/mbrav/pulumi-netbird
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
github.com/mbrav/pulumi-netbird/sdk v0.3.6 v0.3.7 v0.3.8 v0.4.1 v0.5.0 v0.5.1 v0.5.2 v0.5.3 v0.5.4 v0.5.5 v0.5.6 v0.5.7 v0.5.8 v0.5.9 v0.5.10 v0.5.11 # ... and so on
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
pip install sdk/python/bin/dist/pulumi_netbird-0.5.11.tar.gz
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

Every resource below is documented with its inputs, read-only outputs, and a YAML example. Click a type to expand it.

### Networks and routing

<details>
<summary><code>netbird:resource:Network</code></summary>

A NetBird network.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `name` | `string` | yes | The name of the NetBird network. |
| `description` | `string` | no | An optional description of the network. |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
  net-r1:
    type: netbird:resource:Network
    properties:
      name: R1
      description: Network for Region 1
```

</details>

<details>
<summary><code>netbird:resource:NetworkResource</code></summary>

A NetBird network resource, such as a CIDR range assigned to the network. Import ID format: &lt;networkID&gt;/&lt;resourceID&gt;.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `address` | `string` | yes | CIDR or IP address block assigned to the resource. |
| `enabled` | `boolean` | yes | Whether the resource is enabled. |
| `groupIDs` | `string[]` | yes | List of group IDs associated with this resource. |
| `name` | `string` | yes | Name of the network resource. |
| `networkID` | `string` | yes | ID of the network this resource belongs to. |
| `description` | `string` | no | Optional description of the resource. |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
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
```

</details>

<details>
<summary><code>netbird:resource:NetworkRouter</code></summary>

A NetBird network router resource. Import ID format: &lt;networkID&gt;/&lt;routerID&gt;.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enabled` | `boolean` | yes | Whether the router is enabled. |
| `masquerade` | `boolean` | yes | Whether masquerading is enabled. |
| `metric` | `integer` | yes | Routing metric value. |
| `networkID` | `string` | yes | ID of the network this router belongs to. |
| `peer` | `string` | no | Optional peer ID associated with this router. |
| `peerGroups` | `string[]` | no | Optional list of peer group IDs associated with this router. |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
  router-r1:
    type: netbird:resource:NetworkRouter
    properties:
      networkID: ${net-r1.id}
      enabled: true
      masquerade: true
      metric: 10
      peerGroups:
        - ${group-devops.id}
```

</details>

<details>
<summary><code>netbird:resource:Route</code></summary>

A NetBird route resource for directing traffic through exit nodes or routing peers.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `description` | `string` | yes | Route description. |
| `enabled` | `boolean` | yes | Whether the route is enabled. |
| `groups` | `string[]` | yes | Group IDs whose peers will use this route. |
| `keepRoute` | `boolean` | yes | Keep the route after a domain no longer resolves to the IP. |
| `masquerade` | `boolean` | yes | Whether masquerading is enabled for this route. |
| `metric` | `integer` | yes | Route metric; lower value means higher priority. |
| `networkId` | `string` | yes | Route network identifier, used to group HA routes. |
| `accessControlGroups` | `string[]` | no | Access control group IDs associated with this route. |
| `domains` | `string[]` | no | Domain list for dynamic resolution (conflicts with network). |
| `network` | `string` | no | Network CIDR range (conflicts with domains). |
| `peer` | `string` | no | Peer ID acting as the routing peer (conflicts with peerGroups). |
| `peerGroups` | `string[]` | no | Peer group IDs acting as routing peers (conflicts with peer). |
| `skipAutoApply` | `boolean` | no | Skip auto-application for exit-node (0.0.0.0/0) routes. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `networkType` | `string` | Network type (IPv4, IPv6, or domain) — computed by the API. |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
  route-r1-mgmt:
    type: netbird:resource:Route
    properties:
      networkId: route-r1-mgmt
      description: Management subnet route via Region 1
      enabled: true
      network: 192.168.10.0/24
      masquerade: true
      metric: 100
      keepRoute: true
      groups:
        - ${group-devops.id}
      peerGroups:
        - ${group-devops.id}
```

</details>

<details>
<summary><code>netbird:resource:Group</code></summary>

A NetBird group, which represents a collection of peers.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `name` | `string` | yes | The name of the NetBird group. |
| `peers` | `string[]` | no | An optional list of peer IDs to associate with this group. |
| `resources` | `Resource[]` | no | An optional list of resources to associate with this group. |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
  group-devops:
    type: netbird:resource:Group
    properties:
      name: DevOps
      peers: []
```

</details>

### Peers, users and access control

<details>
<summary><code>netbird:resource:Peer</code></summary>

A NetBird peer representing a connected device.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `name` | `string` | yes | The name of the peer. |
| `approvalRequired` | `boolean` | no | **Deprecated:** Cloud only, not maintained in this provider |
| `inactivityExpirationEnabled` | `boolean` | no | Whether Inactivity Expiration is enabled. |
| `loginExpirationEnabled` | `boolean` | no | Whether Login Expiration is enabled. |
| `sshEnabled` | `boolean` | no | Whether SSH is enabled. |

**Example**

```yaml
  # Peers cannot be created through the NetBird API - they register themselves
  # with a setup key. Import an existing peer, then manage its settings here.
  #   pulumi import netbird:resource:Peer peer-gw <peerId>
  peer-gw:
    type: netbird:resource:Peer
    properties:
      name: gateway-01
      sshEnabled: false
      loginExpirationEnabled: true
      inactivityExpirationEnabled: false
      approvalRequired: false
```

</details>

<details>
<summary><code>netbird:resource:IngressPeer</code></summary>

A NetBird ingress peer: an existing peer designated to receive forwarded ingress traffic.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enabled` | `boolean` | yes | Whether the ingress peer is enabled. |
| `fallback` | `boolean` | yes | Whether this ingress peer may be used as a fallback when no ingress peer exists in the forwarded peer's region. |
| `peerId` | `string` | yes | ID of the peer used as an ingress peer. Changing this forces a replacement. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `availablePorts` | `IngressAvailablePorts` | Forwarding ports remaining on the ingress peer. |
| `connected` | `boolean` | Whether the ingress peer is connected to the management server. |
| `ingressIp` | `string` | Ingress IP address where forwarded traffic arrives. |
| `region` | `string` | Region of the ingress peer. |

**Example**

```yaml
  ingress-gw:
    type: netbird:resource:IngressPeer
    properties:
      peerId: ${peer-gw.id}
      enabled: true
      fallback: false
```

</details>

<details>
<summary><code>netbird:resource:SetupKey</code></summary>

Manages a NetBird setup key.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `autoGroups` | `string[]` | yes | Group IDs to auto-assign to peers created with this key. |
| `expiresIn` | `integer` | yes | Time-to-live in seconds from creation; use 0 for no expiration if supported by the API. |
| `name` | `string` | yes | Setup key display name. |
| `type` | `SetupKeyType` | yes | Setup key type: 'one-off' (single use) or 'reusable'. |
| `usageLimit` | `integer` | yes | Maximum uses for reusable keys; 0 = unlimited. |
| `allowExtraDnsLabels` | `boolean` | no | Allow peers to add extra DNS labels beyond the base peer name. |
| `ephemeral` | `boolean` | no | Whether peers registered with this key are ephemeral (auto-expire). |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `expires` | `string` | |
| `key` | `string` | |
| `lastUsed` | `string` | |
| `revoked` | `boolean` | |
| `state` | `string` | |
| `updatedAt` | `string` | |
| `usedTimes` | `integer` | |
| `valid` | `boolean` | |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
  setup-key-devops:
    type: netbird:resource:SetupKey
    properties:
      name: DevOps Onboarding
      type: reusable
      expiresIn: 0
      usageLimit: 0
      ephemeral: false
      allowExtraDnsLabels: false
      autoGroups:
        - ${group-devops.id}
```

</details>

<details>
<summary><code>netbird:resource:User</code></summary>

A NetBird user that receives an invite and is optionally assigned groups and roles.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `role` | `string` | yes | NetBird account role (e.g., 'admin', 'user'). |
| `autoGroups` | `string[]` | no | List of group IDs to auto-assign this user’s peers to. |
| `blocked` | `boolean` | no | Indicates whether the user is blocked from accessing the system. Used only on update, not create. |
| `email` | `string` | no | Email address to send user invite to. |
| `isServiceUser` | `boolean` | no | Whether this user is a service identity. |
| `name` | `string` | no | Full name of the user. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `status` | `UserStatus` | Current account status of the user: 'active', 'blocked', or 'invited'. |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
  user-ci-bot:
    type: netbird:resource:User
    properties:
      role: admin
      isServiceUser: true
      name: ci-bot
      autoGroups:
        - ${group-devops.id}
```

</details>

<details>
<summary><code>netbird:resource:Token</code></summary>

A NetBird personal access token (PAT) for a user. The plaintext token is returned only once, on creation, and is exposed as a secret output. The token cannot be modified after creation; any input change forces a replacement.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `expiresIn` | `integer` | yes | Token lifetime in days. |
| `name` | `string` | yes | Display name of the token. |
| `userId` | `string` | yes | ID of the user the token is issued for. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `createdAt` | `string` | Timestamp the token was created. |
| `createdBy` | `string` | User ID of the principal that created the token. |
| `expirationDate` | `string` | Timestamp the token expires. |
| `lastUsed` | `string` | Timestamp the token was last used, if ever. |
| `token` | `string` | Plaintext token value. Only populated on creation; never returned by the API afterwards. **Secret.** |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
  token-ci-bot:
    type: netbird:resource:Token
    properties:
      userId: ${user-ci-bot.id}
      name: ci-bot-pat
      expiresIn: 90
```

</details>

<details>
<summary><code>netbird:resource:Policy</code></summary>

A NetBird policy defining rules for communication between peers.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enabled` | `boolean` | yes | Enabled Policy status |
| `name` | `string` | yes | Name Policy name identifier |
| `rules` | `PolicyRuleArgs[]` | yes | Rules Policy rule object for policy UI editor |
| `description` | `string` | no | Description Policy friendly description, optional |
| `postureChecks` | `string[]` | no | SourcePostureChecks Posture checks ID's applied to policy source groups, optional |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
  policy-ssh-grp-src-net-dest:
    type: netbird:resource:Policy
    properties:
      name: "SSH Policy - Group to Subnet"
      description: "Allow SSH (22/TCP) from DevOps and Dev groups to Region 1 Net 02"
      enabled: true
      postureChecks:
        - ${posture-devops.id}
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
```

</details>

<details>
<summary><code>netbird:resource:PostureCheck</code></summary>

A NetBird posture check used to validate peer properties before granting policy access.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `checks` | `PostureChecksConfig` | yes | List of checks to perform against peer properties. |
| `name` | `string` | yes | Posture check unique name identifier. |
| `description` | `string` | no | Posture check friendly description. |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
  posture-devops:
    type: netbird:resource:PostureCheck
    properties:
      name: DevOps Posture
      description: "Enforce client version, OS, geo, network range, and process checks"
      checks:
        nbVersionCheck:
          minVersion: "0.28.0"
        osVersionCheck:
          darwin:
            minVersion: "13.0"
          linux:
            minKernelVersion: "5.15"
          windows:
            minKernelVersion: "10.0"
        geoLocationCheck:
          action: allow
          locations:
            - countryCode: DE
            - countryCode: US
              cityName: New York
        peerNetworkRangeCheck:
          action: deny
          ranges:
            - 10.0.0.0/8
        processCheck:
          processes:
            - linuxPath: /usr/bin/netbird
              macPath: /usr/local/bin/netbird
              windowsPath: 'C:\Program Files\NetBird\netbird.exe'
```

</details>

### DNS

<details>
<summary><code>netbird:resource:DNS</code></summary>

A NetBird network.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `description` | `string` | yes | Description of the nameserver group |
| `domains` | `string[]` | yes | Domains Match domain list. It should be empty only if primary is true. |
| `enabled` | `boolean` | yes | Enabled Nameserver group status |
| `groups` | `string[]` | yes | Groups Distribution group IDs that defines group of peers that will use this nameserver group |
| `name` | `string` | yes | Name of nameserver group name |
| `nameservers` | `Nameserver[]` | yes | Nameservers Nameserver list |
| `primary` | `boolean` | yes | Primary Defines if a nameserver group is primary that resolves all domains. It should be true only if domains list is empty. |
| `searchDomainsEnabled` | `boolean` | yes | SearchDomainsEnabled Search domain status for match domains. It should be true only if domains list is not empty. |

**Example**

```yaml
  dns-internal:
    type: netbird:resource:DNS
    properties:
      name: Internal resolvers
      description: Corporate DNS servers for the DevOps group
      enabled: true
      primary: false
      searchDomainsEnabled: false
      domains:
        - corp.example.com
      groups:
        - ${group-devops.id}
      nameservers:
        - ip: 10.10.1.53
          type: udp
          port: 53
```

</details>

<details>
<summary><code>netbird:resource:DNSZone</code></summary>

A NetBird DNS zone.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `distributionGroups` | `string[]` | yes | Group IDs that define groups of peers that will resolve this zone. |
| `domain` | `string` | yes | Zone domain (FQDN). |
| `enableSearchDomain` | `boolean` | yes | Enable this zone as a search domain. |
| `enabled` | `boolean` | yes | Zone status. |
| `name` | `string` | yes | Zone name identifier. |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
  dns-zone-corp:
    type: netbird:resource:DNSZone
    properties:
      name: corp-internal
      domain: corp.example.com
      enabled: true
      enableSearchDomain: true
      distributionGroups:
        - ${group-devops.id}
        - ${group-dev.id}
```

</details>

<details>
<summary><code>netbird:resource:DNSRecord</code></summary>

A DNS record within a NetBird DNS zone.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `content` | `string` | yes | DNS record content (IP address for A/AAAA, domain for CNAME). |
| `name` | `string` | yes | FQDN for the DNS record. Must be a subdomain within or match the zone's domain. |
| `ttl` | `integer` | yes | Time to live in seconds. |
| `type` | `DNSRecordType` | yes | DNS record type. |
| `zoneID` | `string` | yes | ID of the DNS zone this record belongs to. |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
  dns-record-gw-a:
    type: netbird:resource:DNSRecord
    properties:
      zoneID: ${dns-zone-corp.id}
      name: gw.corp.example.com
      type: A
      content: 10.10.1.1
      ttl: 300
```

</details>

<details>
<summary><code>netbird:resource:DNSSettings</code></summary>

NetBird global DNS settings. This is a singleton resource — only one instance exists per account.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `disabledManagementGroups` | `string[]` | yes | Group IDs whose DNS management is disabled. |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
  dns-settings:
    type: netbird:resource:DNSSettings
    properties:
      disabledManagementGroups: []
```

</details>

### Reverse proxy

<details>
<summary><code>netbird:resource:ReverseProxyDomain</code></summary>

A NetBird reverse proxy custom domain.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `domain` | `string` | yes | Domain name for the reverse proxy. |
| `targetCluster` | `string` | yes | The proxy cluster this domain should be validated against. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `requireSubdomain` | `boolean` | Whether a subdomain label is required in front of this domain. |
| `supportsCustomPorts` | `boolean` | Whether the cluster supports binding arbitrary TCP/UDP ports. |
| `type` | `ReverseProxyDomainType` | Type of the reverse proxy domain (custom or free). |
| `validated` | `boolean` | Whether the domain has been validated. A `ReverseProxyService` can only use a validated custom domain. Since netbird v0.79.0 a pending registration expires 48 hours after it is created, and an expired registration is removed unless services are still attached to it. |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
  rp-domain-corp:
    type: netbird:resource:ReverseProxyDomain
    properties:
      domain: proxy.corp.example.com
      targetCluster: eu-central-1
```

</details>

<details>
<summary><code>netbird:resource:ReverseProxyService</code></summary>

A NetBird reverse proxy service.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `domain` | `string` | yes | Domain for the service. A custom domain must already be validated (see `ReverseProxyDomain.validated`) before a service may use it, or before a service may be moved onto it — since netbird v0.79.0 the API rejects both otherwise. |
| `enabled` | `boolean` | yes | Whether the service is enabled. |
| `name` | `string` | yes | Service name. |
| `targets` | `ReverseProxyTarget[]` | yes | List of target backends for this service. |
| `accessGroups` | `string[]` | no | NetBird group IDs whose peers may reach this private service over the tunnel. Required when private=true; ignored otherwise. |
| `accessRestrictions` | `ReverseProxyAccessRestrictions` | no | Connection-level access restrictions based on IP address or geography. Applies to both HTTP and L4 services. |
| `auth` | `ReverseProxyAuth` | no | Authentication configuration for the service (bearer/header/link/password/pin). Mutually exclusive with private=true. |
| `listenPort` | `integer` | no | Port the proxy listens on (L4/TLS only). Set to 0 for auto-assignment. |
| `mode` | `ReverseProxyServiceMode` | no | Service mode: "http" for L7 reverse proxy, "tcp"/"udp"/"tls" for L4 passthrough. |
| `passHostHeader` | `boolean` | no | When true, the original client Host header is passed through to the backend. |
| `private` | `boolean` | no | When true, the service is NetBird-only: peers authenticate via WireGuard tunnel identity and an ACL policy is auto-generated from accessGroups. Requires mode=http. Mutually exclusive with SSO/bearer auth. |
| `rewriteRedirects` | `boolean` | no | When true, Location headers in backend responses are rewritten to the public-facing domain. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `portAutoAssigned` | `boolean` | Whether the listen port was auto-assigned. |
| `proxyCluster` | `string` | The proxy cluster handling this service (derived from domain). |
| `status` | `ReverseProxyServiceStatus` | Current status of the service. |
| `terminated` | `boolean` | Whether the service has been terminated. Terminated services cannot be updated. |

**Example** (from [`examples/yaml/Pulumi.yaml`](examples/yaml/Pulumi.yaml))

```yaml
  rp-svc-api:
    type: netbird:resource:ReverseProxyService
    properties:
      name: api-service
      domain: ${rp-domain-corp.domain}
      enabled: true
      mode: http
      passHostHeader: true
      rewriteRedirects: false
      auth:
        bearerAuth:
          enabled: true
          distributionGroups:
            - ${group-devops.id}
      accessRestrictions:
        allowedCidrs:
          - 10.10.0.0/16
        crowdsecMode: observe
      targets:
        - enabled: true
          host: 10.10.1.10
          port: 8080
          protocol: http
          targetType: host
          targetId: ""
          path: /api
          options:
            pathRewrite: preserve
            requestTimeout: 30s
            skipTlsVerify: false
            customHeaders:
              X-Forwarded-Proto: https
```

</details>

### Agent Network (AI/LLM gateway)

<details>
<summary><code>netbird:resource:AgentNetworkSettings</code></summary>

Per-account NetBird Agent Network gateway settings. This is a singleton resource — only one instance exists per account. Creating it bootstraps the account's gateway endpoint from exactly one of proxyAddress (the server allocates a label beneath that cluster) or endpoint (the hostname is claimed verbatim as a dedicated endpoint). The assigned endpoint is immutable; changing either field replaces the resource, which releases the endpoint and allocates a new one.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enableLogCollection` | `boolean` | yes | Whether per-request access-log entries are collected for this account's agent-network traffic. |
| `enablePromptCollection` | `boolean` | yes | Master switch for request/response prompt capture. Capture runs only when this is on AND a policy guardrail also enables it. |
| `redactPii` | `boolean` | yes | Whether captured prompts have PII redacted. |
| `accessLogRetentionDays` | `integer` | no | Days to retain full access-log rows; older rows are swept. 0 or less means keep indefinitely. Defaults to 30 when omitted. |
| `endpoint` | `string` | no | Hostname to claim as the account's self-addressed (dedicated) endpoint, served only by a proxy declaring exactly that address. Mutually exclusive with proxyAddress; exactly one of the two is required. Immutable — changing it replaces the resource. |
| `proxyAddress` | `string` | no | Cluster address to allocate a labeled endpoint beneath: the server assigns a label and the endpoint becomes `<label>.<proxyAddress>`. Mutually exclusive with endpoint; exactly one of the two is required. Immutable — changing it replaces the resource. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `createdAt` | `string` | Timestamp when the settings row was created. Absent until bootstrapped. |
| `dedicated` | `boolean` | Whether the account's gateway is served by a proxy dedicated to it (endpoint equals proxyAddress). |
| `updatedAt` | `string` | Timestamp when the settings row was last updated. Absent until bootstrapped. |

**Example**

```yaml
  agentnet-settings:
    type: netbird:resource:AgentNetworkSettings
    properties:
      proxyAddress: gateway.netbird.io
      enableLogCollection: true
      enablePromptCollection: false
      redactPii: true
      accessLogRetentionDays: 30
```

</details>

<details>
<summary><code>netbird:resource:AgentNetworkProvider</code></summary>

A NetBird Agent Network provider: an upstream AI/LLM API or gateway (OpenAI, Anthropic, Bedrock, LiteLLM, a custom OpenAI-compatible endpoint, ...) that policies can route traffic to.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `apiKey` | `string` | yes | Upstream provider API key. Sealed at rest on the management server and never returned in responses. **Secret.** |
| `name` | `string` | yes | Display name shown in the dashboard. |
| `providerId` | `string` | yes | Catalog identifier for the upstream AI provider (e.g. openai_api, anthropic_api, azure_openai_api, bedrock_api, vertex_ai_api, mistral_api, custom). See getAgentNetworkCatalogProviders. Changing this forces a replacement. |
| `upstreamUrl` | `string` | yes | Full upstream URL (with scheme) that NetBird forwards traffic to. |
| `bootstrapCluster` | `string` | no | Deprecated and ignored. NetBird removed bootstrap_cluster from the provider API; bootstrap the account's gateway endpoint with netbird:resource:AgentNetworkSettings (proxyAddress or endpoint) instead. **Deprecated:** bootstrapCluster is ignored. Bootstrap the account gateway with AgentNetworkSettings (proxyAddress or endpoint). |
| `enabled` | `boolean` | no | Whether the provider is enabled. |
| `extraValues` | `map[string]string` | no | Operator-typed values for catalog-declared extra headers (see getAgentNetworkCatalogProviders). |
| `identityHeaderGroups` | `string` | no | Wire header name the proxy stamps with the caller's NetBird groups (comma-separated), when the catalog entry supports customizable identity headers. |
| `identityHeaderUserId` | `string` | no | Wire header name the proxy stamps with the caller's display identity, when the catalog entry supports customizable identity headers. |
| `metadataDisabled` | `boolean` | no | Disable identity metadata injection (caller's user + authorizing group) for this provider. |
| `models` | `AgentNetworkProviderModelPricing[]` | no | Models exposed through this endpoint, with operator per-1k price overrides. Empty means all catalog models are allowed at catalog prices. |
| `skipTlsVerification` | `boolean` | no | Skip upstream TLS certificate verification. For self-hosted / internal gateways with a private or self-signed certificate. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `createdAt` | `string` | Timestamp when the provider was created. |
| `updatedAt` | `string` | Timestamp when the provider was last updated. |

**Example**

```yaml
  agentnet-openai:
    type: netbird:resource:AgentNetworkProvider
    properties:
      name: OpenAI prod
      providerId: openai_api
      upstreamUrl: https://api.openai.com
      apiKey:
        fn::secret: ${openaiApiKey}
      enabled: true
      models:
        - id: gpt-4o
          inputPer1k: 0.005
          outputPer1k: 0.015
```

</details>

<details>
<summary><code>netbird:resource:AgentNetworkPolicy</code></summary>

A NetBird Agent Network policy: authorizes source groups to reach a set of Agent Network providers, subject to attached guardrails and token/budget limits.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `destinationProviderIds` | `string[]` | yes | AgentNetworkProvider IDs the source groups can reach. |
| `name` | `string` | yes | Display name for the policy. |
| `sourceGroups` | `string[]` | yes | NetBird group IDs whose members are allowed to call the destination providers. |
| `description` | `string` | no | Optional human-readable description. |
| `enabled` | `boolean` | no | Whether the policy is enabled. |
| `guardrailIds` | `string[]` | no | AgentNetworkGuardrail IDs to attach to this policy. |
| `limits` | `AgentNetworkLimits` | no | Token and budget caps attached directly to the policy. These compose with any guardrail-level checks. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `createdAt` | `string` | Timestamp when the policy was created. |
| `updatedAt` | `string` | Timestamp when the policy was last updated. |

**Example**

```yaml
  agentnet-policy-devops:
    type: netbird:resource:AgentNetworkPolicy
    properties:
      name: DevOps to OpenAI
      description: Lets the DevOps group reach the OpenAI provider.
      enabled: true
      sourceGroups:
        - ${group-devops.id}
      destinationProviderIds:
        - ${agentnet-openai.id}
      guardrailIds:
        - ${agentnet-guardrail.id}
```

</details>

<details>
<summary><code>netbird:resource:AgentNetworkGuardrail</code></summary>

A NetBird Agent Network guardrail: a reusable set of checks (model allowlist, prompt capture) attachable to one or more AgentNetworkPolicy resources.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `checks` | `AgentNetworkGuardrailChecks` | yes | Guardrail check parameters. Each entry has an enabled flag plus per-check configuration; disabled entries are inert. |
| `name` | `string` | yes | Display name for the guardrail. |
| `description` | `string` | no | Optional human-readable description. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `createdAt` | `string` | Timestamp when the guardrail was created. |
| `updatedAt` | `string` | Timestamp when the guardrail was last updated. |

**Example**

```yaml
  agentnet-guardrail:
    type: netbird:resource:AgentNetworkGuardrail
    properties:
      name: Approved models only
      checks:
        modelAllowlist:
          enabled: true
          models:
            - gpt-4o
            - claude-sonnet-4-5
        promptCapture:
          enabled: true
          redactPii: true
```

</details>

<details>
<summary><code>netbird:resource:AgentNetworkBudgetRule</code></summary>

A NetBird Agent Network account-level budget rule: a limit-only rule bound to groups and/or users that applies across all policies as a min-wins ceiling. Empty targets means it applies to every caller.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `limits` | `AgentNetworkLimits` | yes | Token and budget caps attached directly to the rule. |
| `name` | `string` | yes | Display name for the budget rule. |
| `enabled` | `boolean` | no | Whether the rule is enforced. |
| `targetGroups` | `string[]` | no | NetBird group IDs the rule binds. Empty plus empty targetUsers means account-wide. |
| `targetUsers` | `string[]` | no | NetBird user IDs the rule binds directly. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `createdAt` | `string` | Timestamp when the budget rule was created. |
| `updatedAt` | `string` | Timestamp when the budget rule was last updated. |

**Example**

```yaml
  agentnet-budget:
    type: netbird:resource:AgentNetworkBudgetRule
    properties:
      name: Account-wide ceiling
      enabled: true
      targetGroups:
        - ${group-devops.id}
      limits:
        tokenLimit:
          enabled: true
          groupCap: 5000000
          userCap: 500000
          windowSeconds: 86400
        budgetLimit:
          enabled: true
          groupCapUsd: 250
          userCapUsd: 25
          windowSeconds: 86400
```

</details>

### Identity providers and SCIM

<details>
<summary><code>netbird:resource:IdentityProvider</code></summary>

A NetBird identity provider (OIDC) configuration for self-hosted authentication.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `clientId` | `string` | yes | OAuth2 client ID. |
| `clientSecret` | `string` | yes | OAuth2 client secret. **Secret.** |
| `issuer` | `string` | yes | OIDC issuer URL. |
| `name` | `string` | yes | Human-readable name for the identity provider. |
| `type` | `IdentityProviderType` | yes | Type of identity provider. |

**Example**

```yaml
  idp-okta:
    type: netbird:resource:IdentityProvider
    properties:
      name: Okta
      type: okta
      issuer: https://example.okta.com/oauth2/default
      clientId: 0oa1example
      clientSecret:
        fn::secret: ${oktaClientSecret}
```

</details>

<details>
<summary><code>netbird:resource:GoogleIDP</code></summary>

A NetBird Google Workspace identity-provider sync integration.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `customerId` | `string` | yes | Customer ID from Google Workspace account settings. |
| `serviceAccountKey` | `string` | yes | Base64-encoded Google service account key. **Secret.** |
| `connectorId` | `string` | no | DEX connector ID for embedded IdP setups. |
| `enabled` | `boolean` | no | Whether the integration is enabled. |
| `groupPrefixes` | `string[]` | no | start_with patterns for groups to sync. |
| `syncInterval` | `integer` | no | Sync interval in seconds (minimum 300). |
| `userGroupPrefixes` | `string[]` | no | start_with patterns for groups whose users to sync. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `lastSyncedAt` | `string` | Timestamp of the last synchronization. |

**Example**

```yaml
  idp-google-sync:
    type: netbird:resource:GoogleIDP
    properties:
      customerId: C01example
      serviceAccountKey:
        fn::secret: ${googleServiceAccountKey}
      enabled: true
      syncInterval: 300
```

</details>

<details>
<summary><code>netbird:resource:AzureIDP</code></summary>

A NetBird Azure AD (Entra ID) identity-provider sync integration.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `clientId` | `string` | yes | Azure AD application (client) ID. |
| `clientSecret` | `string` | yes | Base64-encoded Azure AD client secret. **Secret.** |
| `host` | `AzureHost` | yes | Azure host domain for the Graph API. Changing this forces a replacement. |
| `tenantId` | `string` | yes | Azure AD tenant ID. |
| `connectorId` | `string` | no | DEX connector ID for embedded IdP setups. |
| `enabled` | `boolean` | no | Whether the integration is enabled. |
| `groupPrefixes` | `string[]` | no | start_with patterns for groups to sync. |
| `syncInterval` | `integer` | no | Sync interval in seconds (minimum 300). |
| `userGroupPrefixes` | `string[]` | no | start_with patterns for groups whose users to sync. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `lastSyncedAt` | `string` | Timestamp of the last synchronization. |

**Example**

```yaml
  idp-azure-sync:
    type: netbird:resource:AzureIDP
    properties:
      host: microsoft.com
      tenantId: 00000000-0000-0000-0000-000000000000
      clientId: 11111111-1111-1111-1111-111111111111
      clientSecret:
        fn::secret: ${azureClientSecret}
      enabled: true
```

</details>

<details>
<summary><code>netbird:resource:OktaScimIDP</code></summary>

A NetBird Okta SCIM identity-provider sync integration.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `connectionName` | `string` | yes | The Okta enterprise connection name on Auth0. Changing this forces a replacement. |
| `connectorId` | `string` | no | DEX connector ID for embedded IdP setups. |
| `enabled` | `boolean` | no | Whether the integration is enabled. |
| `groupPrefixes` | `string[]` | no | start_with patterns for groups to sync. |
| `userGroupPrefixes` | `string[]` | no | start_with patterns for groups whose users to sync. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `authToken` | `string` | SCIM API token. Returned in full only on creation; masked afterwards, so the created value is preserved. **Secret.** |
| `lastSyncedAt` | `string` | Timestamp of the last synchronization. |

**Example**

```yaml
  idp-okta-scim:
    type: netbird:resource:OktaScimIDP
    properties:
      connectionName: okta-scim
      enabled: true
```

</details>

<details>
<summary><code>netbird:resource:ScimIntegration</code></summary>

A NetBird generic SCIM identity-provider integration.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `prefix` | `string` | yes | The connection prefix used for the SCIM provider. |
| `provider` | `string` | yes | Name of the SCIM identity provider. Changing this forces a replacement. |
| `connectorId` | `string` | no | DEX connector ID for embedded IdP setups. |
| `enabled` | `boolean` | no | Whether the integration is enabled. |
| `groupPrefixes` | `string[]` | no | start_with patterns for groups to sync. |
| `userGroupPrefixes` | `string[]` | no | start_with patterns for groups whose users to sync. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `authToken` | `string` | SCIM API token. Returned in full only on creation; masked afterwards, so the created value is preserved. **Secret.** |
| `lastSyncedAt` | `string` | Timestamp of the last synchronization. |

**Example**

```yaml
  scim-okta:
    type: netbird:resource:ScimIntegration
    properties:
      provider: okta
      prefix: netbird
      enabled: true
```

</details>

## 🔍 Invoke Functions (Data Sources)

Invoke functions are **read-only** — they query live NetBird state and return data without managing any resources. Use them to reference existing objects by a human-readable key rather than a hardcoded ID.

<details>
<summary><code>netbird:function:getAgentNetworkCatalogProviders</code></summary>

List the catalog of upstream AI providers supported by NetBird Agent Network, with their default models and pricing. Useful for looking up a providerId and model ids for AgentNetworkProvider.

Takes no arguments.

**Returns**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `providers` | `AgentNetworkCatalogProviderEntry[]` | The list of catalog providers. |

</details>

<details>
<summary><code>netbird:function:getCountries</code></summary>

List all countries known to NetBird's geo-location database. Useful for populating PostureCheck geo-location rules.

Takes no arguments.

**Returns**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `countries` | `Country[]` | The list of countries. |

</details>

<details>
<summary><code>netbird:function:getCountryCities</code></summary>

List the cities within a country from NetBird's geo-location database. Useful for populating PostureCheck geo-location rules with specific cities.

**Arguments**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `countryCode` | `string` | yes | 2-letter ISO 3166-1 alpha-2 country code to list cities for. |

**Returns**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `cities` | `City[]` | The list of cities in the country. |

</details>

<details>
<summary><code>netbird:function:getPeers</code></summary>

List all NetBird peers, optionally filtered to those belonging to a specific group ID.

**Arguments**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `groupId` | `string` | no | Optional group ID to filter peers. When set, only peers that belong to this group are returned. |

**Returns**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `peers` | `PeerSummary[]` | The list of peers matching the filter criteria. |

</details>

<details>
<summary><code>netbird:function:getReverseProxyClusters</code></summary>

List all NetBird reverse proxy clusters. Clusters are auto-provisioned server-side (there is no create endpoint); use this to discover a cluster address to pin a ReverseProxyDomain against via its targetCluster input.

**Arguments**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `type` | `string` | no | Optional cluster type filter: 'account' (owned/operated by the account, BYOP) or 'shared' (operated by NetBird and shared across accounts). When set, only clusters of this type are returned. |

**Returns**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `clusters` | `ProxyClusterSummary[]` | The list of reverse proxy clusters matching the filter criteria. |

</details>

<details>
<summary><code>netbird:function:lookupGroup</code></summary>

Look up an existing NetBird group by name and return its ID, peer list, and resource list.

**Arguments**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `name` | `string` | yes | The name of the group to look up. |

**Returns**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `groupId` | `string` | The NetBird group ID. |
| `name` | `string` | The group name. |
| `peers` | `string[]` | IDs of peers belonging to the group. |
| `peersCount` | `integer` | Number of peers in the group. |
| `resources` | `ResourceRef[]` | Resources associated with the group. |
| `resourcesCount` | `integer` | Number of resources in the group. |

</details>

<details>
<summary><code>netbird:function:lookupPeer</code></summary>

Look up an existing NetBird peer by name and return its ID, IP address, and group memberships.

**Arguments**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `name` | `string` | yes | The name (hostname) of the peer to look up. |

**Returns**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `connected` | `boolean` | Whether the peer is currently connected to the management server. |
| `dnsLabel` | `string` | The DNS label used to form the peer's FQDN. |
| `groups` | `string[]` | IDs of groups the peer belongs to. |
| `hostname` | `string` | The OS hostname of the machine. |
| `ip` | `string` | The WireGuard IP address assigned to the peer. |
| `name` | `string` | The peer name. |
| `os` | `string` | Operating system string reported by the peer. |
| `peerId` | `string` | The NetBird peer ID. |

</details>

<details>
<summary><code>netbird:function:lookupRoute</code></summary>

Look up an existing NetBird route by network CIDR and return its ID, routing peers, and configuration.

**Arguments**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `network` | `string` | yes | The network CIDR (e.g. '10.0.0.0/8') to look up. Returns the first matching route. |

**Returns**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `description` | `string` | The route description. |
| `domains` | `string[]` | Domain list for dynamic resolution (if this is a domain route). |
| `enabled` | `boolean` | Whether the route is enabled. |
| `groups` | `string[]` | Group IDs that have access to this route. |
| `masquerade` | `boolean` | Whether the routing peer masquerades traffic to this prefix. |
| `metric` | `integer` | Route metric; lower number means higher priority. |
| `network` | `string` | The network CIDR range. |
| `peer` | `string` | The peer ID acting as the routing peer (mutually exclusive with peerGroups). |
| `peerGroups` | `string[]` | Group IDs whose peers act as routing peers (mutually exclusive with peer). |
| `routeId` | `string` | The NetBird route ID. |

</details>

<details>
<summary><code>netbird:function:lookupSetupKey</code></summary>

Look up an existing NetBird setup key by name and return its metadata (not the secret key value).

**Arguments**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `name` | `string` | yes | The name of the setup key to look up. |

**Returns**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `autoGroups` | `string[]` | Group IDs automatically assigned to peers that register with this key. |
| `ephemeral` | `boolean` | Whether peers registered with this key are ephemeral. |
| `expires` | `string` | Key expiration timestamp in RFC3339 format. |
| `lastUsed` | `string` | Timestamp of last key use in RFC3339 format. |
| `name` | `string` | The setup key name. |
| `revoked` | `boolean` | Whether the setup key has been revoked. |
| `setupKeyId` | `string` | The NetBird setup key ID. |
| `state` | `string` | The setup key state: 'valid', 'overused', 'expired', or 'revoked'. |
| `type` | `string` | The setup key type: 'one-off' or 'reusable'. |
| `usageLimit` | `integer` | Maximum number of times the key can be used (0 = unlimited for reusable). |

</details>

<details>
<summary><code>netbird:function:lookupUser</code></summary>

Look up an existing NetBird user by email and return their ID, role, and auto-group assignments.

**Arguments**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `email` | `string` | yes | The email address of the user to look up. |

**Returns**

| Name | Type | Description |
| ---- | ---- | ----------- |
| `autoGroups` | `string[]` | Group IDs automatically assigned to peers registered by this user. |
| `email` | `string` | The user's email address. |
| `isBlocked` | `boolean` | Whether the user is blocked from accessing the system. |
| `name` | `string` | The user's display name. |
| `role` | `string` | The user's role in the NetBird account (e.g. admin, user). |
| `userId` | `string` | The NetBird user ID. |

</details>

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

<details>
<summary><code>netbird:component:DNSZoneBundle</code></summary>

Experimental. Declares a DNSZone and one DNSRecord per records[] entry as a single unit. The created zoneID is wired into every record automatically.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `distributionGroups` | `string[]` | yes | Group IDs whose peers receive this DNS zone. |
| `domain` | `string` | yes | DNS domain (e.g. corp.example.com). |
| `enableSearchDomain` | `boolean` | yes | Whether to enable the zone as a search domain for peers. |
| `enabled` | `boolean` | yes | Whether the DNS zone is active. |
| `name` | `string` | yes | Logical name for the DNS zone resource. |
| `records` | `DNSRecordSpec[]` | yes | DNS records to create within the zone. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `recordIds` | `string[]` | IDs of the created DNSRecord resources, in declaration order. |
| `zoneId` | `string` | ID of the created DNSZone resource. |

**Example**

```yaml
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
```

</details>

<details>
<summary><code>netbird:component:NetworkBundle</code></summary>

Experimental. Declares a Network, a NetworkRouter, and one NetworkResource per subnets[] entry as a single unit. The created networkID is wired into the router and every subnet automatically.

**Inputs**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `name` | `string` | yes | Name of the overlay network. |
| `router` | `NetworkRouterSpec` | yes | Router configuration attached to the network. |
| `subnets` | `NetworkSubnetSpec[]` | yes | Subnet resources to attach to the network. |
| `description` | `string` | no | Optional description for the network. |

**Outputs** (read-only, in addition to the inputs above)

| Name | Type | Description |
| ---- | ---- | ----------- |
| `networkId` | `string` | ID of the created Network resource. |
| `routerId` | `string` | ID of the created NetworkRouter resource. |
| `subnetIds` | `string[]` | IDs of the created NetworkResource (subnet) resources, in declaration order. |

**Example**

```yaml
  r1:
    type: netbird:component:NetworkBundle
    properties:
      name: Region1
      description: Network for Region 1
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
```

</details>

## 🧩 Nested Types

Object and enum types referenced from the resource tables above.

<details>
<summary><strong>Object types</strong> (42)</summary>

<details>
<summary><code>netbird:component:DNSRecordSpec</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `content` | `string` | yes | Record value: an IP address for A/AAAA or a hostname for CNAME. |
| `name` | `string` | yes | Fully-qualified DNS record name (e.g. api.corp.example.com). |
| `ttl` | `integer` | yes | Time-to-live in seconds. |
| `type` | `string` | yes | Record type: A, AAAA, or CNAME. |

</details>

<details>
<summary><code>netbird:component:NetworkRouterSpec</code></summary>

Router configuration for a NetworkBundle. At least one of peer or peerGroups must be set - the underlying NetworkRouter rejects a router with neither.

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enabled` | `boolean` | yes | Whether the router is enabled. |
| `masquerade` | `boolean` | yes | Whether to masquerade traffic through the router. |
| `metric` | `integer` | yes | Route metric; lower values have higher priority. |
| `peer` | `string` | no | Specific peer to use as router. |
| `peerGroups` | `string[]` | no | Peer groups to use as router peers. |

</details>

<details>
<summary><code>netbird:component:NetworkSubnetSpec</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `address` | `string` | yes | CIDR block for the subnet (e.g. 10.10.1.0/24). |
| `enabled` | `boolean` | yes | Whether the subnet resource is enabled. |
| `groupIDs` | `string[]` | yes | Group IDs that have access to this subnet. |
| `name` | `string` | yes | Display name for the subnet resource. |
| `description` | `string` | no | Optional description for the subnet resource. |

</details>

<details>
<summary><code>netbird:function:AgentNetworkCatalogModel</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `contextWindow` | `integer` | yes | Maximum context window in tokens. |
| `id` | `string` | yes | Catalog model identifier as exposed by the upstream provider. |
| `inputPer1k` | `number` | yes | Default input token price per 1k tokens, in USD. |
| `label` | `string` | yes | Human-friendly model name. |
| `outputPer1k` | `number` | yes | Default output token price per 1k tokens, in USD. |
| `cacheCreationPer1k` | `number` | no | Anthropic-shape cache rate: default cost per 1k cache-creation tokens, in USD. |
| `cacheReadPer1k` | `number` | no | Anthropic-shape cache rate: default cost per 1k cache-read tokens, in USD. |
| `cachedInputPer1k` | `number` | no | OpenAI-shape cache rate: default cost per 1k cached prompt tokens, in USD. |

</details>

<details>
<summary><code>netbird:function:AgentNetworkCatalogProviderEntry</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `authHeaderTemplate` | `string` | yes | Template the proxy uses to inject the API key. |
| `brandColor` | `string` | yes | Hex brand color used to render the provider badge in the dashboard. |
| `defaultContentType` | `string` | yes | Default Content-Type for upstream requests. |
| `defaultHost` | `string` | yes | Default upstream host suggested when adding a provider of this type. |
| `description` | `string` | yes | Short description shown in the provider picker. |
| `id` | `string` | yes | Catalog provider identifier, used as AgentNetworkProvider.providerId. |
| `kind` | `string` | yes | Presentation grouping: "provider" (first-party vendor API), "gateway" (routing/aggregation layer), or "custom" (generic OpenAI-compatible endpoint). |
| `models` | `AgentNetworkCatalogModel[]` | yes | Catalog models available for this provider, with default pricing. |
| `name` | `string` | yes | Display name for the provider. |
| `pricingSurfaces` | `string[]` | no | Cost-meter pricing surfaces this provider's traffic is metered under ("openai", "anthropic", "bedrock"). |

</details>

<details>
<summary><code>netbird:function:City</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `cityName` | `string` | yes | Commonly used English name of the city. |
| `geonameId` | `integer` | yes | Integer ID of the record in the GeoNames database. |

</details>

<details>
<summary><code>netbird:function:Country</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `countryCode` | `string` | yes | 2-letter ISO 3166-1 alpha-2 country code. |
| `countryName` | `string` | yes | Commonly used English name of the country. |

</details>

<details>
<summary><code>netbird:function:PeerSummary</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `connected` | `boolean` | yes | Whether the peer is currently connected to the management server. |
| `dnsLabel` | `string` | yes | The DNS label used to form the peer's FQDN. |
| `groups` | `string[]` | yes | IDs of groups the peer belongs to. |
| `hostname` | `string` | yes | The OS hostname of the machine. |
| `id` | `string` | yes | The peer ID. |
| `ip` | `string` | yes | The WireGuard IP address assigned to the peer. |
| `name` | `string` | yes | The peer name. |

</details>

<details>
<summary><code>netbird:function:ProxyClusterSummary</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `address` | `string` | yes | Cluster address used for CNAME targets; the value to pass as a ReverseProxyDomain targetCluster. |
| `connectedProxies` | `integer` | yes | Number of proxy nodes currently connected. |
| `id` | `string` | yes | Unique identifier of the proxy cluster. |
| `online` | `boolean` | yes | Whether at least one proxy in the cluster has heartbeated within the active window. |
| `private` | `boolean` | yes | True when at least one connected proxy is embedded in a netbird client and serving over a WireGuard tunnel. |
| `requireSubdomain` | `boolean` | yes | Whether services on this cluster must include a subdomain label. |
| `supportsCrowdsec` | `boolean` | yes | Whether all active proxies in the cluster have CrowdSec configured. |
| `supportsCustomPorts` | `boolean` | yes | Whether the cluster supports binding arbitrary TCP/UDP ports. |
| `type` | `string` | yes | Source of the cluster: 'account' (BYOP) or 'shared' (operated by NetBird). |

</details>

<details>
<summary><code>netbird:function:ResourceRef</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `id` | `string` | yes | The unique identifier of the resource. |
| `type` | `string` | yes | The type of resource: 'domain', 'host', or 'subnet'. |

</details>

<details>
<summary><code>netbird:resource:AgentNetworkBudgetLimit</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enabled` | `boolean` | yes | Whether the budget limit is enforced. |
| `groupCapUsd` | `number` | yes | USD allowed per source group within the window (each group has its own bucket). 0 means uncapped. |
| `userCapUsd` | `number` | yes | USD allowed per individual user within the window. 0 means uncapped. |
| `windowSeconds` | `integer` | yes | Reset frequency in seconds. Minimum 60 when the limit is enabled. |

</details>

<details>
<summary><code>netbird:resource:AgentNetworkGuardrailChecks</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `modelAllowlist` | `AgentNetworkGuardrailModelAllowlist` | yes | Restricts requests to an explicit set of catalog model IDs. |
| `promptCapture` | `AgentNetworkGuardrailPromptCapture` | yes | Controls request/response prompt capture. |

</details>

<details>
<summary><code>netbird:resource:AgentNetworkGuardrailModelAllowlist</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enabled` | `boolean` | yes | Whether the model allowlist check is enforced. |
| `models` | `string[]` | yes | Allowed catalog model IDs. Requests for any other model are denied. |

</details>

<details>
<summary><code>netbird:resource:AgentNetworkGuardrailPromptCapture</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enabled` | `boolean` | yes | Whether prompt/response capture is enforced for requests passing through this guardrail. |
| `redactPii` | `boolean` | yes | Whether captured prompts have PII redacted. |

</details>

<details>
<summary><code>netbird:resource:AgentNetworkLimits</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `budgetLimit` | `AgentNetworkBudgetLimit` | yes | USD budget cap composed with any guardrail-level checks. |
| `tokenLimit` | `AgentNetworkTokenLimit` | yes | Token cap composed with any guardrail-level checks. |

</details>

<details>
<summary><code>netbird:resource:AgentNetworkProviderModelPricing</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `id` | `string` | yes | Catalog model identifier (e.g. "gpt-4o-mini"). |
| `inputPer1k` | `number` | yes | Cost per 1k input tokens, in USD. |
| `outputPer1k` | `number` | yes | Cost per 1k output tokens, in USD. |
| `cacheCreationPer1k` | `number` | no | Anthropic-shape cache rate: cost per 1k cache-creation tokens, in USD. |
| `cacheReadPer1k` | `number` | no | Anthropic-shape cache rate: cost per 1k cache-read tokens, in USD. |
| `cachedInputPer1k` | `number` | no | OpenAI-shape cache rate: cost per 1k cached prompt tokens, in USD. |

</details>

<details>
<summary><code>netbird:resource:AgentNetworkTokenLimit</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enabled` | `boolean` | yes | Whether the token limit is enforced. |
| `groupCap` | `integer` | yes | Tokens allowed per source group within the window (each group has its own bucket). 0 means uncapped. |
| `userCap` | `integer` | yes | Tokens allowed per individual user within the window. 0 means uncapped. |
| `windowSeconds` | `integer` | yes | Reset frequency in seconds. Minimum 60 when the limit is enabled. |

</details>

<details>
<summary><code>netbird:resource:IngressAvailablePorts</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `tcp` | `integer` | yes | Number of available TCP ports left on the ingress peer. |
| `udp` | `integer` | yes | Number of available UDP ports left on the ingress peer. |

</details>

<details>
<summary><code>netbird:resource:Nameserver</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `ip` | `string` | yes | IP of Nameserver |
| `port` | `integer` | yes | Port Nameserver Port |
| `type` | `NameserverNsType` | yes | NsType Nameserver Type |

</details>

<details>
<summary><code>netbird:resource:PolicyRuleArgs</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `action` | `RuleAction` | yes | Action Policy rule accept or drops packets |
| `bidirectional` | `boolean` | yes | Bidirectional Define if the rule is applicable in both directions, sources, and destinations. |
| `enabled` | `boolean` | yes | Enabled Policy rule status |
| `name` | `string` | yes | Name Policy rule name identifier |
| `protocol` | `Protocol` | yes | Protocol Policy rule type of the traffic |
| `authorizedGroups` | `map[string]string[]` | no | Map of user group IDs to a list of local users for network access authorization |
| `description` | `string` | no | Description Policy rule friendly description |
| `destinationResource` | `Resource` | no | DestinationResource for the rule |
| `destinations` | `string[]` | no | Destinations Policy rule destination group IDs |
| `id` | `string` | no | ID Policy rule. |
| `portRanges` | `RulePortRange[]` | no | PortRanges Policy rule affected ports ranges list |
| `ports` | `string[]` | no | Ports Policy rule affected ports |
| `sourceResource` | `Resource` | no | SourceResource for the rule |
| `sources` | `string[]` | no | Sources Policy rule source group IDs |

</details>

<details>
<summary><code>netbird:resource:PolicyRuleState</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `action` | `RuleAction` | yes | Action Policy rule accept or drops packets |
| `bidirectional` | `boolean` | yes | Bidirectional Define if the rule is applicable in both directions, sources, and destinations. |
| `enabled` | `boolean` | yes | Enabled Policy rule status |
| `name` | `string` | yes | Name Policy rule name identifier |
| `protocol` | `Protocol` | yes | Protocol Policy rule type of the traffic |
| `authorizedGroups` | `map[string]string[]` | no | AuthorizedGroups Map of user group IDs to a list of local users |
| `description` | `string` | no | Description Policy rule friendly description |
| `destinationResource` | `Resource` | no | DestinationResource for the rule |
| `destinations` | `RuleGroup[]` | no | Destinations Policy rule destination group IDs |
| `id` | `string` | no | ID Policy rule. |
| `portRanges` | `RulePortRange[]` | no | PortRanges Policy rule affected ports ranges list |
| `ports` | `string[]` | no | Ports Policy rule affected ports |
| `sourceResource` | `Resource` | no | SourceResource for the rule |
| `sources` | `RuleGroup[]` | no | Sources Policy rule source group IDs |

</details>

<details>
<summary><code>netbird:resource:PostureChecksConfig</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `geoLocationCheck` | `PostureGeoLocationCheck` | no | Posture check for geo location. |
| `nbVersionCheck` | `PostureMinVersionCheck` | no | Posture check for the minimum NetBird client version. |
| `osVersionCheck` | `PostureOSVersionCheck` | no | Posture check for the minimum operating system version. |
| `peerNetworkRangeCheck` | `PosturePeerNetworkRangeCheck` | no | Posture check based on peer local network addresses. |
| `processCheck` | `PostureProcessCheck` | no | Posture check for required binaries running on the peer. |

</details>

<details>
<summary><code>netbird:resource:PostureGeoLocationCheck</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `action` | `PostureGeoLocationAction` | yes | Action to take upon geo location match (allow or deny). |
| `locations` | `PostureLocation[]` | yes | List of geo locations to which the check applies. |

</details>

<details>
<summary><code>netbird:resource:PostureLocation</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `countryCode` | `string` | yes | 2-letter ISO 3166-1 alpha-2 country code. |
| `cityName` | `string` | no | Commonly used English name of the city. |

</details>

<details>
<summary><code>netbird:resource:PostureMinKernelVersionCheck</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `minKernelVersion` | `string` | yes | Minimum acceptable kernel version string. |

</details>

<details>
<summary><code>netbird:resource:PostureMinVersionCheck</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `minVersion` | `string` | yes | Minimum acceptable version string. |

</details>

<details>
<summary><code>netbird:resource:PostureOSVersionCheck</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `android` | `PostureMinVersionCheck` | no | Minimum version check for Android. |
| `darwin` | `PostureMinVersionCheck` | no | Minimum version check for macOS. |
| `ios` | `PostureMinVersionCheck` | no | Minimum version check for iOS. |
| `linux` | `PostureMinKernelVersionCheck` | no | Minimum kernel version check for Linux. |
| `windows` | `PostureMinKernelVersionCheck` | no | Minimum kernel version check for Windows. |

</details>

<details>
<summary><code>netbird:resource:PosturePeerNetworkRangeCheck</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `action` | `PosturePeerNetworkRangeAction` | yes | Action to take when the peer's network range matches (allow or deny). |
| `ranges` | `string[]` | yes | List of CIDR network ranges to match against. |

</details>

<details>
<summary><code>netbird:resource:PostureProcess</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `linuxPath` | `string` | no | Path to the process executable on Linux. |
| `macPath` | `string` | no | Path to the process executable on macOS. |
| `windowsPath` | `string` | no | Path to the process executable on Windows. |

</details>

<details>
<summary><code>netbird:resource:PostureProcessCheck</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `processes` | `PostureProcess[]` | yes | List of processes that must be running on the peer. |

</details>

<details>
<summary><code>netbird:resource:Resource</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `id` | `string` | yes | The unique identifier of the resource. |
| `type` | `Type` | yes | The type of resource: 'domain', 'host', or 'subnet'. |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyAccessRestrictions</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `allowedCidrs` | `string[]` | no | CIDR allowlist. If non-empty, only IPs matching these CIDRs are allowed. |
| `allowedCountries` | `string[]` | no | ISO 3166-1 alpha-2 country codes to allow. If non-empty, only these countries are permitted. |
| `blockedCidrs` | `string[]` | no | CIDR blocklist. Connections from these CIDRs are rejected. Evaluated after allowedCidrs. |
| `blockedCountries` | `string[]` | no | ISO 3166-1 alpha-2 country codes to block. |
| `crowdsecMode` | `ReverseProxyCrowdsecMode` | no | CrowdSec IP reputation mode. Only available when the proxy cluster supports CrowdSec. |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyAuth</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `bearerAuth` | `ReverseProxyBearerAuth` | no | Bearer (SSO) authentication configuration. |
| `headerAuths` | `ReverseProxyHeaderAuth[]` | no | Header-based authentication entries. |
| `linkAuth` | `ReverseProxyLinkAuth` | no | Link-based authentication configuration. |
| `passwordAuth` | `ReverseProxyPasswordAuth` | no | Password-based authentication configuration. |
| `pinAuth` | `ReverseProxyPINAuth` | no | PIN-based authentication configuration. |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyBearerAuth</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enabled` | `boolean` | yes | Whether bearer auth is enabled. |
| `distributionGroups` | `string[]` | no | List of group IDs that can use bearer auth. |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyHeaderAuth</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enabled` | `boolean` | yes | Whether this header auth entry is enabled. |
| `header` | `string` | yes | The header name to match. |
| `value` | `string` | yes | The header value to match. |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyLinkAuth</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enabled` | `boolean` | yes | Whether link auth is enabled. |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyPINAuth</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enabled` | `boolean` | yes | Whether PIN auth is enabled. |
| `pin` | `string` | yes | The PIN required to access the service. **Secret.** |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyPasswordAuth</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enabled` | `boolean` | yes | Whether password auth is enabled. |
| `password` | `string` | yes | The password required to access the service. **Secret.** |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyTarget</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `enabled` | `boolean` | yes | Whether this target is enabled. |
| `port` | `integer` | yes | Backend port for this target. |
| `protocol` | `ReverseProxyTargetProtocol` | yes | Protocol to use when connecting to the backend. |
| `targetId` | `string` | yes | Target ID (assigned by the server). |
| `targetType` | `ReverseProxyTargetType` | yes | Target type. |
| `host` | `string` | no | Backend IP or domain for this target. |
| `options` | `ReverseProxyTargetOptions` | no | Advanced per-target proxy options. |
| `path` | `string` | no | URL path prefix for this target (HTTP only). |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyTargetOptions</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `customHeaders` | `map[string]string` | no | Extra headers sent to the backend. Hop-by-hop and proxy-managed headers are rejected. |
| `directUpstream` | `boolean` | no | When true, the proxy dials this target via the host's network stack instead of through its embedded NetBird client. |
| `pathRewrite` | `ReverseProxyPathRewrite` | no | How the request path is rewritten before forwarding. Default strips the matched prefix; "preserve" keeps the full path. |
| `proxyProtocol` | `boolean` | no | Send PROXY Protocol v2 header to this backend (TCP/TLS only). |
| `requestTimeout` | `string` | no | Per-target response timeout as a Go duration string (e.g. "30s", "2m"). |
| `sessionIdleTimeout` | `string` | no | Idle timeout before a UDP session is reaped, as a Go duration string (e.g. "30s", "2m"). |
| `skipTlsVerify` | `boolean` | no | Skip TLS certificate verification for this backend. |

</details>

<details>
<summary><code>netbird:resource:RuleGroup</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `id` | `string` | yes | The unique identifier of the group. |
| `name` | `string` | yes | The name of the group. |

</details>

<details>
<summary><code>netbird:resource:RulePortRange</code></summary>

**Fields**

| Name | Type | Required | Description |
| ---- | ---- | -------- | ----------- |
| `end` | `integer` | yes | End of port range |
| `start` | `integer` | yes | Start of port range |

</details>

</details>

<details>
<summary><strong>Enum types</strong> (18)</summary>

<details>
<summary><code>netbird:resource:AzureHost</code></summary>

| Value | Description |
| ----- | ----------- |
| `microsoft.com` | Commercial Microsoft Graph host. |
| `microsoft.us` | US Government Microsoft Graph host. |

</details>

<details>
<summary><code>netbird:resource:DNSRecordType</code></summary>

| Value | Description |
| ----- | ----------- |
| `A` | IPv4 address record. |
| `AAAA` | IPv6 address record. |
| `CNAME` | Canonical name record. |

</details>

<details>
<summary><code>netbird:resource:IdentityProviderType</code></summary>

| Value | Description |
| ----- | ----------- |
| `adfs` | Microsoft AD FS. |
| `entra` | Microsoft Entra ID. |
| `google` | Google. |
| `microsoft` | Microsoft. |
| `oidc` | Generic OIDC provider. |
| `okta` | Okta. |
| `pocketid` | PocketID. |
| `zitadel` | Zitadel. |

</details>

<details>
<summary><code>netbird:resource:NameserverNsType</code></summary>

| Value | Description |
| ----- | ----------- |
| `udp` | UDP type |

</details>

<details>
<summary><code>netbird:resource:PostureGeoLocationAction</code></summary>

| Value | Description |
| ----- | ----------- |
| `allow` | Allow peers from the specified locations. |
| `deny` | Deny peers from the specified locations. |

</details>

<details>
<summary><code>netbird:resource:PosturePeerNetworkRangeAction</code></summary>

| Value | Description |
| ----- | ----------- |
| `allow` | Allow peers whose local network matches. |
| `deny` | Deny peers whose local network matches. |

</details>

<details>
<summary><code>netbird:resource:Protocol</code></summary>

| Value | Description |
| ----- | ----------- |
| `all` | All protocols |
| `icmp` | ICMP protocol |
| `tcp` | TCP protocol |
| `udp` | UDP protocol |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyCrowdsecMode</code></summary>

| Value | Description |
| ----- | ----------- |
| `enforce` | Block connections flagged by CrowdSec. |
| `observe` | Only log connections flagged by CrowdSec. |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyDomainType</code></summary>

| Value | Description |
| ----- | ----------- |
| `custom` | A custom domain managed by the user. |
| `free` | A free managed domain provided by NetBird. |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyPathRewrite</code></summary>

| Value | Description |
| ----- | ----------- |
| `preserve` | Keep the full original request path instead of stripping the matched prefix. |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyServiceMode</code></summary>

| Value | Description |
| ----- | ----------- |
| `http` | L7 HTTP reverse proxy mode. |
| `tcp` | L4 TCP passthrough mode. |
| `tls` | L4 TLS passthrough mode. |
| `udp` | L4 UDP passthrough mode. |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyServiceStatus</code></summary>

| Value | Description |
| ----- | ----------- |
| `active` | Service is provisioned and serving. |
| `certificate_failed` | TLS certificate issuance failed. |
| `certificate_pending` | TLS certificate issuance is in progress. |
| `error` | Service is in an error state. |
| `pending` | Service is being provisioned. |
| `tunnel_not_created` | Underlying tunnel has not been created yet. |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyTargetProtocol</code></summary>

| Value | Description |
| ----- | ----------- |
| `http` | HTTP protocol. |
| `https` | HTTPS protocol. |
| `tcp` | TCP protocol. |
| `udp` | UDP protocol. |

</details>

<details>
<summary><code>netbird:resource:ReverseProxyTargetType</code></summary>

| Value | Description |
| ----- | ----------- |
| `cluster` | Proxy-cluster target. |
| `domain` | Domain-based target. |
| `host` | Host-based target. |
| `peer` | Peer-based target. |
| `subnet` | Subnet-based target. |

</details>

<details>
<summary><code>netbird:resource:RuleAction</code></summary>

| Value | Description |
| ----- | ----------- |
| `accept` | Accept action |
| `drop` | Drop action |

</details>

<details>
<summary><code>netbird:resource:SetupKeyType</code></summary>

| Value | Description |
| ----- | ----------- |
| `reusable` | Reusable setup key that supports multiple peers. |
| `one-off` | One-off setup key that can be used only once. |

</details>

<details>
<summary><code>netbird:resource:Type</code></summary>

| Value | Description |
| ----- | ----------- |
| `domain` | A domain resource (e.g., example.com). |
| `host` | A host resource (e.g., peer or device). |
| `subnet` | A subnet resource (e.g., 192.168.0.0/24). |

</details>

<details>
<summary><code>netbird:resource:UserStatus</code></summary>

| Value | Description |
| ----- | ----------- |
| `active` | User has accepted the invite and is active. |
| `blocked` | User is blocked from accessing the system. |
| `invited` | User has been invited but has not yet joined. |

</details>

</details>

## 📁 Repository Structure

- `provider/` – Go implementation of the provider
- `sdk/go/netbird/` – Go SDK for the NetBird provider
- `examples/` – Example Pulumi projects using the provider

## 📚 References

- [Pulumi Go Provider Docs](https://github.com/pulumi/pulumi-go-provider)
- [NetBird Documentation](https://docs.netbird.io/)
- [Pulumi YAML Documentation](https://www.pulumi.com/docs/using-pulumi/yaml/)
