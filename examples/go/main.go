package main

import (
	"github.com/mbrav/pulumi-netbird/sdk/go/netbird"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		myRandomResource, err := netbird.NewRandom(ctx, "myRandomResource", &netbird.RandomArgs{
			Length: pulumi.Int(24),
		})
		if err != nil {
			return err
		}

		ctx.Export("output", pulumi.Map{
			"value": myRandomResource.Result,
		})

		return nil
	})
}
