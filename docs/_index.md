---
title: netbird
meta_desc: Manage NetBird resources declaratively using Pulumi's infrastructure-as-code framework.
layout: overview
---

The NetBird Pulumi Provider enables you to manage [NetBird](https://netbird.io) resources declaratively using Pulumi infrastructure as code.
It supports all major Pulumi-supported languages and works with both NetBird Cloud (`https://api.netbird.io`) and self-hosted management servers.

The provider covers 15 resource types including groups, peers, policies, setup keys, DNS, networks, posture checks, users, and reverse proxy services. See the [Installation & Configuration](installation-configuration/) page for the full list and setup instructions.

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
