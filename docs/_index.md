---
title: netbird
meta_desc: Manage NetBird resources declaratively using Pulumi's infrastructure-as-code framework.
layout: overview
---

The NetBird Pulumi Provider enables you to manage [NetBird](https://netbird.io) resources declaratively using Pulumi infrastructure as code.
It supports all major Pulumi-supported languages and works with both NetBird Cloud (`https://api.netbird.io`) and self-hosted management servers.

The provider covers 16 resource types including groups, peers, policies, setup keys, DNS, networks, posture checks, users, and reverse proxy services, plus 6 **invoke functions** (data sources) for reading existing NetBird objects by name, email, or CIDR without managing them as resources. It also ships two **experimental components** (`NetworkBundle`, `DNSZoneBundle`) as a proof of concept — see the note below before using them. See the [Installation & Configuration](installation-configuration/) page for the full list and setup instructions.

## Example

{{< chooser language "go,python,typescript,csharp" >}}

{{% choosable language go %}}

```go
package main

import (
    netbird "github.com/mbrav/pulumi-netbird/sdk/go/netbird/resource"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        _, err := netbird.NewGroup(ctx, "group-devops", &netbird.GroupArgs{
            Name:  pulumi.String("DevOps"),
            Peers: pulumi.StringArray{},
        })
        return err
    })
}
```

{{% /choosable %}}

{{% choosable language python %}}

```python
import pulumi
from pulumi_netbird import resource

devops_group = resource.Group("group-devops",
    name="DevOps",
    peers=[])
```

{{% /choosable %}}

{{% choosable language typescript %}}

```typescript
import * as netbird from "@mbrav/pulumi-netbird";

const devOpsGroup = new netbird.resource.Group("group-devops", {
    name: "DevOps",
    peers: [],
});
```

{{% /choosable %}}

{{% choosable language csharp %}}

```csharp
using Pulumi;
using Mbrav.PulumiNetbird.Resource;

return await Deployment.RunAsync(() =>
{
    var devOpsGroup = new Group("group-devops", new GroupArgs
    {
        Name = "DevOps",
        Peers = new InputList<string>(),
    });
});
```

{{% /choosable %}}

{{< /chooser >}}

## Invoke Functions

Invoke functions query live NetBird state and return data without creating or modifying resources. They are the equivalent of Terraform data sources.

| Function | Looks up by |
| -------- | ----------- |
| `netbird:function:getPeers` | All peers, with an optional group ID filter |
| `netbird:function:lookupGroup` | Group name |
| `netbird:function:lookupPeer` | Peer name |
| `netbird:function:lookupRoute` | Network CIDR |
| `netbird:function:lookupSetupKey` | Setup key name |
| `netbird:function:lookupUser` | User email address |

{{< chooser language "go,python,typescript,csharp" >}}

{{% choosable language go %}}

```go
devopsGroup, err := netbird.LookupGroup(ctx, &netbird.LookupGroupArgs{
    Name: "DevOps",
}, nil)
if err != nil {
    return err
}
// devopsGroup.GroupId contains the resolved NetBird group ID
```

{{% /choosable %}}

{{% choosable language python %}}

```python
import pulumi_netbird as netbird

devops_group = netbird.lookup_group(name="DevOps")
# devops_group.group_id contains the resolved NetBird group ID
```

{{% /choosable %}}

{{% choosable language typescript %}}

```typescript
import * as netbird from "@mbrav/pulumi-netbird";

const devopsGroup = netbird.lookupGroup({ name: "DevOps" });
// devopsGroup.then(g => g.groupId) contains the resolved NetBird group ID
```

{{% /choosable %}}

{{% choosable language csharp %}}

```csharp
var devopsGroup = await Netbird.LookupGroup.InvokeAsync(new LookupGroupArgs { Name = "DevOps" });
// devopsGroup.GroupId contains the resolved NetBird group ID
```

{{% /choosable %}}

{{< /chooser >}}

## Experimental Components

> **Proof of concept only.** These components explore the `pulumi-go-provider` component API. Their interface may change without notice; do not use them in production.

Components bundle multiple NetBird resources into a single declaration, wiring resource IDs automatically.

| Component | Pulumi type | Child resources created |
| --------- | ----------- | ----------------------- |
| Network bundle | `netbird:component:NetworkBundle` | `Network` + `NetworkRouter` + N `NetworkResource` subnets |
| DNS zone bundle | `netbird:component:DNSZoneBundle` | `DNSZone` + N `DNSRecord`s |

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
```
