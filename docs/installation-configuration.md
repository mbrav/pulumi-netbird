---
title: NetBird Pulumi Provider Installation & Configuration
meta_desc: Information on how to install the NetBird Pulumi provider.
layout: installation
---

## Installation

The NetBird Pulumi provider is available for multiple Pulumi-supported languages. Use the appropriate package for your stack:

- **Python**: [`pulumi_netbird`](https://pypi.org/project/pulumi-netbird/)
- **Go**: [`github.com/mbrav/pulumi-netbird/sdk/go/netbird`](https://pkg.go.dev/github.com/mbrav/pulumi-netbird/sdk)
- **.NET**: [`Mbrav.PulumiNetbird`](https://www.nuget.org/packages/Mbrav.PulumiNetbird)

The provider plugin is distributed via GitHub Releases. Pulumi will install it automatically on first use, or you can install it explicitly:

```sh
pulumi plugin install resource netbird v0.3.1 \
  --server github://api.github.com/mbrav/pulumi-netbird
```

## Configuration

Configure the provider using the Pulumi CLI. The `token` value is sensitive and should always be set as a secret:

```sh
pulumi config set netbird:url https://api.netbird.io
pulumi config set --secret netbird:token <YOUR_API_TOKEN>
```

Alternatively, set the equivalent environment variables before running `pulumi up`:

```sh
export NETBIRD_URL=https://api.netbird.io
export NETBIRD_TOKEN=<YOUR_API_TOKEN>
```

### Configuration Reference

| Option  | Environment Variable | Required | Default | Description |
| ------- | -------------------- | -------- | ------- | ----------- |
| `url`   | `NETBIRD_URL`        | Yes      | `https://api.netbird.io` | URL of your NetBird management API |
| `token` | `NETBIRD_TOKEN`      | Yes      | —       | API token for authentication (mark as secret) |

## Supported Resources

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
| Setup key | `netbird:resource:SetupKey` |
| User | `netbird:resource:User` |
