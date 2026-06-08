package main

import (
	"github.com/mbrav/pulumi-netbird/sdk/go/netbird/function"
	"github.com/mbrav/pulumi-netbird/sdk/go/netbird/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// cfg := config.New(ctx, "")
		// netbirdNetbirdToken := cfg.RequireObject("netbird:netbirdToken")
		// netbirdNetbirdUrl := cfg.RequireObject("netbird:netbirdUrl")

		// ── Groups ────────────────────────────────────────────────────────────
		// Peer groups used to scope policies, network resources, and DNS zones.

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

		// ── Networks ──────────────────────────────────────────────────────────
		// Overlay network for Region 1 that groups related subnets and routers.

		netR1, err := resource.NewNetwork(ctx, "net-r1", &resource.NetworkArgs{
			Name:        pulumi.String("R1"),
			Description: pulumi.String("Network for Region 1"),
		})
		if err != nil {
			return err
		}

		// ── Network Resources ─────────────────────────────────────────────────
		// Subnet resources attached to net-r1; accessible by the DevOps group.

		_, err = resource.NewNetworkResource(ctx, "netres-r1-net-01", &resource.NetworkResourceArgs{
			Name:        pulumi.String("Region 1 Net 01"),
			Description: pulumi.StringPtr("Network 01 in Region 1"),
			NetworkID:   netR1.ID(),
			Address:     pulumi.String("10.10.1.0/24"),
			Enabled:     pulumi.Bool(true),
			GroupIDs: pulumi.StringArray{
				groupDevops.ID(),
			},
		})
		if err != nil {
			return err
		}

		netresR1Net02, err := resource.NewNetworkResource(ctx, "netres-r1-net-02", &resource.NetworkResourceArgs{
			Name:        pulumi.String("Region 1 Net 02"),
			Description: pulumi.StringPtr("Network 02 in Region 1"),
			NetworkID:   netR1.ID(),
			Address:     pulumi.String("10.10.2.0/24"),
			Enabled:     pulumi.Bool(true),
			GroupIDs: pulumi.StringArray{
				groupDevops.ID(),
			},
		})
		if err != nil {
			return err
		}

		_, err = resource.NewNetworkResource(ctx, "netres-r1-net-03", &resource.NetworkResourceArgs{
			Name:        pulumi.String("Region 1 Net 03"),
			Description: pulumi.StringPtr("Network 03 in Region 1"),
			NetworkID:   netR1.ID(),
			Address:     pulumi.String("10.10.3.0/24"),
			Enabled:     pulumi.Bool(true),
			GroupIDs: pulumi.StringArray{
				groupDevops.ID(),
			},
		})
		if err != nil {
			return err
		}

		// ── Network Router ────────────────────────────────────────────────────
		// Masquerading router for net-r1; uses DevOps group as peer group.

		_, err = resource.NewNetworkRouter(ctx, "router-r1", &resource.NetworkRouterArgs{
			NetworkID:  netR1.ID(),
			Enabled:    pulumi.Bool(true),
			Masquerade: pulumi.Bool(true),
			Metric:     pulumi.Int(10),
			PeerGroups: pulumi.StringArray{
				groupDevops.ID(),
			},
		})
		if err != nil {
			return err
		}

		// ── Posture Checks ────────────────────────────────────────────────────
		// PostureCheck validates peer properties before granting policy access.
		// The check below requires a minimum NetBird client version, enforces a
		// minimum OS kernel version on Linux/Windows, restricts by geo location,
		// blocks peers on private RFC-1918 ranges, and requires a specific binary.

		postureCheck, err := resource.NewPostureCheck(ctx, "posture-devops", &resource.PostureCheckArgs{
			Name:        pulumi.String("DevOps Posture"),
			Description: pulumi.StringPtr("Enforce client version, OS, geo, network range, and process checks"),
			Checks: resource.PostureChecksConfigArgs{
				// Require minimum NetBird client version 0.28.0.
				NbVersionCheck: resource.PostureMinVersionCheckArgs{
					MinVersion: pulumi.String("0.28.0"),
				},
				// Require minimum OS versions per platform.
				OsVersionCheck: resource.PostureOSVersionCheckArgs{
					Darwin: resource.PostureMinVersionCheckArgs{
						MinVersion: pulumi.String("13.0"),
					},
					Linux: resource.PostureMinKernelVersionCheckArgs{
						MinKernelVersion: pulumi.String("5.15"),
					},
					Windows: resource.PostureMinKernelVersionCheckArgs{
						MinKernelVersion: pulumi.String("10.0"),
					},
				},
				// Allow peers only from Germany (DE) or United States (US).
				GeoLocationCheck: resource.PostureGeoLocationCheckArgs{
					Action: resource.PostureGeoLocationActionAllow,
					Locations: resource.PostureLocationArray{
						resource.PostureLocationArgs{
							CountryCode: pulumi.String("DE"),
						},
						resource.PostureLocationArgs{
							CountryCode: pulumi.String("US"),
							CityName:    pulumi.StringPtr("New York"),
						},
					},
				},
				// Deny peers whose local network is a private RFC-1918 /8 range.
				PeerNetworkRangeCheck: resource.PosturePeerNetworkRangeCheckArgs{
					Action: resource.PosturePeerNetworkRangeActionDeny,
					Ranges: pulumi.StringArray{
						pulumi.String("10.0.0.0/8"),
					},
				},
				// Require the NetBird agent binary to be present on each platform.
				ProcessCheck: resource.PostureProcessCheckArgs{
					Processes: resource.PostureProcessArray{
						resource.PostureProcessArgs{
							LinuxPath:   pulumi.StringPtr("/usr/bin/netbird"),
							MacPath:     pulumi.StringPtr("/usr/local/bin/netbird"),
							WindowsPath: pulumi.StringPtr("C:\\Program Files\\NetBird\\netbird.exe"),
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// ── Policies ──────────────────────────────────────────────────────────
		// Policy: DevOps/Dev → Region 1 Net 02 (subnet destination, with posture).

		_, err = resource.NewPolicy(ctx, "policy-ssh-grp-src-net-dest", &resource.PolicyArgs{
			Name:        pulumi.String("SSH Policy - Group to Subnet"),
			Description: pulumi.String("Allow SSH (22/TCP) from DevOps and Dev groups to Region 1 Net 02"),
			Enabled:     pulumi.Bool(true),
			PostureChecks: pulumi.StringArray{
				postureCheck.ID(),
			},
			Rules: resource.PolicyRuleArgsArray{
				&resource.PolicyRuleArgsArgs{
					Name:          pulumi.String("SSH Access - Group → Subnet"),
					Description:   pulumi.StringPtr("Allow unidirectional SSH from DevOps & Dev groups to Net 02"),
					Bidirectional: pulumi.Bool(false),
					Action:        resource.RuleActionAccept,
					Enabled:       pulumi.Bool(true),
					Protocol:      resource.ProtocolTCP,
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

		// Policy: DevOps → Backoffice (group-to-group destination).

		_, err = resource.NewPolicy(ctx, "policy-ssh-grp-src-grp-dest", &resource.PolicyArgs{
			Name:          pulumi.String("SSH Policy - Group to Group"),
			Description:   pulumi.String("Allow SSH (22/TCP) from DevOps to Backoffice group resources"),
			Enabled:       pulumi.Bool(true),
			PostureChecks: pulumi.StringArray{},
			Rules: resource.PolicyRuleArgsArray{
				&resource.PolicyRuleArgsArgs{
					Name:          pulumi.String("SSH Access - Group → Group"),
					Description:   pulumi.StringPtr("SSH from DevOps group to Backoffice group"),
					Bidirectional: pulumi.Bool(false),
					Action:        resource.RuleActionAccept,
					Enabled:       pulumi.Bool(true),
					Protocol:      resource.ProtocolTCP,
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

		// ── DNS Settings ──────────────────────────────────────────────────────
		// Singleton resource — only one exists per account.
		// Disables DNS management for the HR group so their peers use default DNS.

		_, err = resource.NewDNSSettings(ctx, "dns-settings", &resource.DNSSettingsArgs{
			DisabledManagementGroups: pulumi.StringArray{},
		})
		if err != nil {
			return err
		}

		// ── DNS Zones ─────────────────────────────────────────────────────────
		// Internal DNS zone for the corp.example.com domain distributed to DevOps.

		dnsZone, err := resource.NewDNSZone(ctx, "dns-zone-corp", &resource.DNSZoneArgs{
			Name:               pulumi.String("corp-internal"),
			Domain:             pulumi.String("corp.example.com"),
			Enabled:            pulumi.Bool(true),
			EnableSearchDomain: pulumi.Bool(true),
			DistributionGroups: pulumi.StringArray{
				groupDevops.ID(),
				groupDev.ID(),
			},
		})
		if err != nil {
			return err
		}

		// ── DNS Records ───────────────────────────────────────────────────────
		// A record pointing the gateway hostname to an internal IPv4 address.

		_, err = resource.NewDNSRecord(ctx, "dns-record-gw-a", &resource.DNSRecordArgs{
			ZoneID:  dnsZone.ID(),
			Name:    pulumi.String("gw.corp.example.com"),
			Type:    resource.DNSRecordTypeA,
			Content: pulumi.String("10.10.1.1"),
			Ttl:     pulumi.Int(300),
		})
		if err != nil {
			return err
		}

		// CNAME record aliasing api.corp.example.com to the gateway hostname.
		_, err = resource.NewDNSRecord(ctx, "dns-record-api-cname", &resource.DNSRecordArgs{
			ZoneID:  dnsZone.ID(),
			Name:    pulumi.String("api.corp.example.com"),
			Type:    resource.DNSRecordTypeCNAME,
			Content: pulumi.String("gw.corp.example.com"),
			Ttl:     pulumi.Int(300),
		})
		if err != nil {
			return err
		}

		// ── Setup Key ─────────────────────────────────────────────────────────
		// Reusable setup key for onboarding new peers into the DevOps group.
		// UsageLimit 0 = unlimited uses; ExpiresIn 0 = no expiry.

		_, err = resource.NewSetupKey(ctx, "setup-key-devops", &resource.SetupKeyArgs{
			Name:                pulumi.String("DevOps Onboarding"),
			Type:                resource.SetupKeyTypeReusable,
			ExpiresIn:           pulumi.Int(0),
			UsageLimit:          pulumi.Int(0),
			Ephemeral:           pulumi.BoolPtr(false),
			AllowExtraDnsLabels: pulumi.BoolPtr(false),
			AutoGroups: pulumi.StringArray{
				groupDevops.ID(),
			},
		})
		if err != nil {
			return err
		}

		// ── Route ─────────────────────────────────────────────────────────────
		// Network route advertising 192.168.10.0/24 through the DevOps peer group.
		// Masquerade hides the source IP behind the router's address.

		_, err = resource.NewRoute(ctx, "route-r1-mgmt", &resource.RouteArgs{
			NetworkId:   pulumi.String("route-r1-mgmt"),
			Description: pulumi.String("Management subnet route via Region 1"),
			Enabled:     pulumi.Bool(true),
			Network:     pulumi.StringPtr("192.168.10.0/24"),
			Masquerade:  pulumi.Bool(true),
			Metric:      pulumi.Int(100),
			KeepRoute:   pulumi.Bool(true),
			Groups: pulumi.StringArray{
				groupDevops.ID(),
			},
			PeerGroups: pulumi.StringArray{
				groupDevops.ID(),
			},
		})
		if err != nil {
			return err
		}

		// ── Service User ──────────────────────────────────────────────────────
		// Automation service user with admin role; placed in the DevOps group.

		_, err = resource.NewUser(ctx, "user-ci-bot", &resource.UserArgs{
			Role:          pulumi.String("admin"),
			IsServiceUser: pulumi.BoolPtr(true),
			Name:          pulumi.StringPtr("ci-bot"),
			AutoGroups: pulumi.StringArray{
				groupDevops.ID(),
			},
		})
		if err != nil {
			return err
		}

		// ── Reverse Proxy Domain ──────────────────────────────────────────────
		// Custom domain validated against a specific proxy cluster.
		// Domain changes trigger resource replacement (no Update endpoint).

		rpDomain, err := resource.NewReverseProxyDomain(ctx, "rp-domain-corp", &resource.ReverseProxyDomainArgs{
			Domain:        pulumi.String("proxy.corp.example.com"),
			TargetCluster: pulumi.String("eu-central-1"),
		})
		if err != nil {
			return err
		}

		// ── Reverse Proxy Service ─────────────────────────────────────────────
		// HTTP (L7) reverse proxy service routing traffic to an internal backend.
		// PassHostHeader preserves the original Host header at the backend.

		_, err = resource.NewReverseProxyService(ctx, "rp-svc-api", &resource.ReverseProxyServiceArgs{
			Name:             pulumi.String("api-service"),
			Domain:           rpDomain.Domain,
			Enabled:          pulumi.Bool(true),
			Mode:             resource.ReverseProxyServiceModeHttp.ToReverseProxyServiceModePtrOutput(),
			PassHostHeader:   pulumi.BoolPtr(true),
			RewriteRedirects: pulumi.BoolPtr(false),
			Targets: resource.ReverseProxyTargetArray{
				resource.ReverseProxyTargetArgs{
					Enabled:    pulumi.Bool(true),
					Host:       pulumi.StringPtr("10.10.1.10"),
					Port:       pulumi.Int(8080),
					Protocol:   resource.ReverseProxyTargetProtocolHttp,
					TargetType: resource.ReverseProxyTargetTypeHost,
					TargetId:   pulumi.String(""),
				},
			},
		})
		if err != nil {
			return err
		}

		// ── Invoke Functions (Data Sources) ──────────────────────────────────
		// Functions are read-only: they query live NetBird state without
		// managing any resources. Use them to reference objects that exist
		// outside this stack rather than hardcoding IDs.

		// Look up the built-in "All" group that NetBird creates automatically.
		// Every peer joins this group on registration; it cannot be created
		// via Pulumi, so a lookup function is the correct way to reference it.
		allGroup, err := function.LookupGroup(ctx, &function.LookupGroupArgs{
			Name: "All",
		}, nil)
		if err != nil {
			return err
		}

		// List all peers currently registered in the account.
		allPeers, err := function.GetPeers(ctx, &function.GetPeersArgs{}, nil)
		if err != nil {
			return err
		}

		// List only peers belonging to the DevOps group.
		// groupDevops.ID() returns a pulumi.IDOutput; use ApplyT to resolve
		// the string before passing it to a synchronous invoke function.
		devopsPeers, err := function.GetPeers(ctx, &function.GetPeersArgs{
			GroupId: pulumi.StringRef(allGroup.GroupId),
		}, nil)
		if err != nil {
			return err
		}

		// ── Outputs ───────────────────────────────────────────────────────────

		ctx.Export("networkR1", pulumi.StringMapMap{
			"value": pulumi.StringMap{
				"name": netR1.Name,
				"id":   netR1.ID(),
			},
		})

		ctx.Export("dnsZoneCorp", pulumi.StringMapMap{
			"value": pulumi.StringMap{
				"name":   dnsZone.Name,
				"domain": dnsZone.Domain,
				"id":     dnsZone.ID(),
			},
		})

		// Export invoke function results.
		ctx.Export("allGroupId", pulumi.String(allGroup.GroupId))
		ctx.Export("allGroupPeersCount", pulumi.Int(allGroup.PeersCount))
		ctx.Export("allPeerCount", pulumi.Int(len(allPeers.Peers)))
		ctx.Export("devopsPeerCount", pulumi.Int(len(devopsPeers.Peers)))

		return nil
	})
}
