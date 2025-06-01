---
title: Netbird
meta_desc: Provides an overview of the Netbird Provider for Pulumi.
layout: overview
---

The NetBird Pulumi Provider enables you to manage [NetBird](https://netbird.io) resources using Pulumi infrastructure as code across your preferred programming language.
You must configure the provider with the proper credentials and endpoint to interact with your NetBird instance.

## Example

{{< chooser language "go,python,typescript,csharp" >}}

{{% choosable language go %}}

```go
package main

import (
    "github.com/mbrav/pulumi-netbird/sdk/go/netbird/resource"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        _, err := resource.NewGroup(ctx, "group-devops", &resource.GroupArgs{
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
import * as pulumi from "@pulumi/pulumi";
import * as netbird from "@mbrav/pulumi-netbird";

const devOpsGroup = new netbird.Group("group-devops", {
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
        Peers = {},
    });
});
```

{{% /choosable %}}

{{< /chooser >}}
