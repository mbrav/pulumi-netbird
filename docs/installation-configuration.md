---
title: NetBird Pulumi Provider Installation & Configuration
meta_desc: Information on how to install the NetBird Pulumi provider.
layout: installation
---

## Installation

The NetBird Pulumi provider is available for multiple Pulumi-supported languages. Use the appropriate package for your stack:

- **JavaScript/TypeScript**: [`@mbrav/pulumi-netbird`](#) *(Link pending)*
- **Python**: [`pulumi_netbird`](https://pypi.org/project/pulumi-netbird/)
- **Go**: [`github.com/mbrav/pulumi-netbird/sdk/go/netbird`](https://pkg.go.dev/github.com/mbrav/pulumi-netbird/sdk)
- **.NET**: [`Mbrav.PulimiNetbird`](https://www.nuget.org/packages/Mbrav.PulimiNetbird)

## Setup

Before provisioning any resources, ensure that the provider is correctly configured to communicate with your NetBird instance. This requires setting up endpoint details and authentication credentials.

## Configuration

You can configure the provider using Pulumi CLI:

```sh
pulumi config set netbird:<option> [--secret]
```

### Available Configuration Options

| Option  | Environment Variable | Required | Description                            |
| ------- | -------------------- | -------- | -------------------------------------- |
| `url`   | `NETBIRD_URL`        | Yes      | The URL of your NetBird management service |
| `token` | `NETBIRD_TOKEN`      | Yes      | A valid API token for authentication       |
