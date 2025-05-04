package main

import (
	"github.com/mbrav/pulumi-netbird/sdk/go/netbird"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create the "Management" network
		_, err := netbird.NewNetwork(ctx, "net-management", &netbird.NetworkArgs{
			Name:        pulumi.String("Management"),
			Description: pulumi.String("Network for Management"),
		})
		if err != nil {
			return err
		}
		//
		// // Create a policy to allow SSH access to the Management network
		// _, err = netbird.NewPolicy(ctx, "allow-ssh-management", &netbird.PolicyArgs{
		// 	Name:        pulumi.String("Allow SSH to Management"),
		// 	Description: pulumi.String("Allow SSH (TCP port 22) from all sources to Management"),
		// 	Enabled:     pulumi.Bool(true),
		// 	Rules: []api.PolicyRuleUpdate{
		// 		{
		// 			Name:          pulumi.String("SSH Rule"),
		// 			Action:        pulumi.String("accept"),
		// 			Protocol:      pulumi.String("tcp"),
		// 			Ports:         pulumi.StringArray{pulumi.String("22")},
		// 			Bidirectional: pulumi.Bool(false),
		// 			Sources:       pulumi.StringArray{pulumi.String("all")}, // might be a group ID in real API
		// 			Destinations:  pulumi.StringArray{network.NbId},
		// 		},
		// 	},
		// 	SourcePostureChecks: pulumi.StringArray{}, // optional, leave empty
		// })
		// if err != nil {
		// 	return err
		// }

		return nil
	})
}
