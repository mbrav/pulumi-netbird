package main

import (
	"github.com/mbrav/pulumi-netbird/sdk/go/netbird/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// cfg := config.New(ctx, "")
		// netbirdNetbirdToken := cfg.RequireObject("netbird:netbirdToken")
		// netbirdNetbirdUrl := cfg.RequireObject("netbird:netbirdUrl")
		groupDevops, err := resource.NewGroup(ctx, "group-devops", &resource.GroupArgs{
			Name:  pulumi.String("DevOps"),
			Peers: pulumi.StringArray{},
		})
		if err != nil {
			return err
		}

		groupDev, err := resource.NewGroup(ctx, "group-dev", &resource.GroupArgs{
			Name:  pulumi.String("Dev"),
			Peers: pulumi.StringArray{},
		})
		if err != nil {
			return err
		}

		groupBackoffice, err := resource.NewGroup(ctx, "group-backoffice", &resource.GroupArgs{
			Name:  pulumi.String("Backoffice"),
			Peers: pulumi.StringArray{},
		})
		if err != nil {
			return err
		}

		_, err = resource.NewGroup(ctx, "group-hr", &resource.GroupArgs{
			Name:  pulumi.String("HR"),
			Peers: pulumi.StringArray{},
		})
		if err != nil {
			return err
		}

		netR1, err := resource.NewNetwork(ctx, "net-r1", &resource.NetworkArgs{
			Name:        pulumi.String("R1"),
			Description: pulumi.String("Network for Region 1"),
		})
		if err != nil {
			return err
		}

		_, err = resource.NewNetworkResource(ctx, "netres-r1-net-01", &resource.NetworkResourceArgs{
			Name:        pulumi.String("Region 1 Net 01"),
			Description: pulumi.String("Network 01 in Region 1"),
			Network_id:  netR1.ID(),
			Address:     pulumi.String("10.10.1.0/24"),
			Enabled:     pulumi.Bool(true),
			Group_ids: pulumi.StringArray{
				groupDevops.ID(),
			},
		})
		if err != nil {
			return err
		}

		netresR1Net02, err := resource.NewNetworkResource(ctx, "netres-r1-net-02", &resource.NetworkResourceArgs{
			Name:        pulumi.String("Region 1 Net 02"),
			Description: pulumi.String("Network 02 in S1 Region 1"),
			Network_id:  netR1.ID(),
			Address:     pulumi.String("10.10.2.0/24"),
			Enabled:     pulumi.Bool(true),
			Group_ids: pulumi.StringArray{
				groupDevops.ID(),
			},
		})
		if err != nil {
			return err
		}

		_, err = resource.NewNetworkResource(ctx, "netres-r1-net-03", &resource.NetworkResourceArgs{
			Name:        pulumi.String("Region 1 Net 03"),
			Description: pulumi.String("Network 03 in Region 1"),
			Network_id:  netR1.ID(),
			Address:     pulumi.String("10.10.3.0/24"),
			Enabled:     pulumi.Bool(true),
			Group_ids: pulumi.StringArray{
				groupDevops.ID(),
			},
		})
		if err != nil {
			return err
		}

		_, err = resource.NewNetworkRouter(ctx, "router-r1", &resource.NetworkRouterArgs{
			Network_id: netR1.ID(),
			Enabled:    pulumi.Bool(true),
			Masquerade: pulumi.Bool(true),
			Metric:     pulumi.Int(10),
			Peer:       pulumi.String(""),
			Peer_groups: pulumi.StringArray{
				groupDevops.ID(),
			},
		})
		if err != nil {
			return err
		}

		_, err = resource.NewPolicy(ctx, "policy-ssh-grp-src-net-dest", &resource.PolicyArgs{
			Name:           pulumi.String("SSH Policy - Group to Subnet"),
			Description:    pulumi.String("Allow SSH (22/TCP) from DevOps and Dev groups to Region 1 Net 02"),
			Enabled:        pulumi.Bool(true),
			Posture_checks: pulumi.StringArray{},
			Rules: resource.PolicyRuleArgsArray{
				&resource.PolicyRuleArgsArgs{
					Name:          pulumi.String("SSH Access - Group → Subnet"),
					Description:   pulumi.String("Allow unidirectional SSH from DevOps & Dev groups to Net 02"),
					Bidirectional: pulumi.Bool(false),
					Action:        resource.RuleActionAccept,
					Enabled:       pulumi.Bool(true),
					Protocol:      resource.ProtocolTcp,
					Ports: pulumi.StringArray{
						pulumi.String("22"),
					},
					Sources: pulumi.StringArray{
						groupDevops.ID(),
						groupDev.ID(),
					},
					DestinationResource: &resource.ResourceArgs{
						Type: resource.TypeSubnet,
						Id:   netresR1Net02.ID(),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		_, err = resource.NewPolicy(ctx, "policy-ssh-grp-src-grp-dest", &resource.PolicyArgs{
			Name:           pulumi.String("SSH Policy - Group to Group"),
			Description:    pulumi.String("Allow SSH (22/TCP) from DevOps to Backoffice group resources"),
			Enabled:        pulumi.Bool(true),
			Posture_checks: pulumi.StringArray{},
			Rules: resource.PolicyRuleArgsArray{
				&resource.PolicyRuleArgsArgs{
					Name:          pulumi.String("SSH Access - Group → Group"),
					Description:   pulumi.String("SSH from DevOps group to Backoffice group"),
					Bidirectional: pulumi.Bool(false),
					Action:        resource.RuleActionAccept,
					Enabled:       pulumi.Bool(true),
					Protocol:      resource.ProtocolTcp,
					Ports: pulumi.StringArray{
						pulumi.String("22"),
					},
					Sources: pulumi.StringArray{
						groupDevops.ID(),
					},
					Destinations: pulumi.StringArray{
						groupBackoffice.ID(),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("networkR1", pulumi.StringMapMap{
			"value": pulumi.StringMap{
				"name": netR1.Name,
				"id":   netR1.ID(),
			},
		})

		return nil
	})
}
