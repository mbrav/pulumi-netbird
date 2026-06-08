package component

import (
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NetworkRouterSpec holds the router configuration within a NetworkBundle.
type NetworkRouterSpec struct {
	Enabled    bool      `pulumi:"enabled"`
	Masquerade bool      `pulumi:"masquerade"`
	Metric     int       `pulumi:"metric"`
	PeerGroups *[]string `pulumi:"peerGroups,optional"`
	Peer       *string   `pulumi:"peer,optional"`
}

// Annotate adds schema descriptions to NetworkRouterSpec fields.
func (r *NetworkRouterSpec) Annotate(a infer.Annotator) {
	a.Describe(&r.Enabled, "Whether the router is enabled.")
	a.Describe(&r.Masquerade, "Whether to masquerade traffic through the router.")
	a.Describe(&r.Metric, "Route metric; lower values have higher priority.")
	a.Describe(&r.PeerGroups, "Peer groups to use as router peers.")
	a.Describe(&r.Peer, "Specific peer to use as router.")
}

// NetworkSubnetSpec holds configuration for a single subnet resource in a NetworkBundle.
type NetworkSubnetSpec struct {
	Name        string   `pulumi:"name"`
	Address     string   `pulumi:"address"`
	Enabled     bool     `pulumi:"enabled"`
	GroupIDs    []string `pulumi:"groupIDs"`
	Description *string  `pulumi:"description,optional"`
}

// Annotate adds schema descriptions to NetworkSubnetSpec fields.
func (s *NetworkSubnetSpec) Annotate(a infer.Annotator) {
	a.Describe(&s.Name, "Display name for the subnet resource.")
	a.Describe(&s.Address, "CIDR block for the subnet (e.g. 10.10.1.0/24).")
	a.Describe(&s.Enabled, "Whether the subnet resource is enabled.")
	a.Describe(&s.GroupIDs, "Group IDs that have access to this subnet.")
	a.Describe(&s.Description, "Optional description for the subnet resource.")
}

// NetworkBundleArgs are the inputs for a NetworkBundle component.
type NetworkBundleArgs struct {
	Name        string              `pulumi:"name"`
	Description *string             `pulumi:"description,optional"`
	Router      NetworkRouterSpec   `pulumi:"router"`
	Subnets     []NetworkSubnetSpec `pulumi:"subnets"`
}

// Annotate adds schema descriptions to NetworkBundleArgs fields.
func (n *NetworkBundleArgs) Annotate(a infer.Annotator) {
	a.Describe(&n.Name, "Name of the overlay network.")
	a.Describe(&n.Description, "Optional description for the network.")
	a.Describe(&n.Router, "Router configuration attached to the network.")
	a.Describe(&n.Subnets, "Subnet resources to attach to the network.")
}

// NetworkBundleState holds the outputs of a NetworkBundle component.
type NetworkBundleState struct {
	pulumi.ResourceState

	NetworkID pulumi.StringOutput      `pulumi:"networkId"`
	RouterID  pulumi.StringOutput      `pulumi:"routerId"`
	SubnetIDs pulumi.StringArrayOutput `pulumi:"subnetIds"`
}

// Annotate adds schema descriptions to NetworkBundleState fields.
func (s *NetworkBundleState) Annotate(a infer.Annotator) {
	a.Describe(&s.NetworkID, "ID of the created Network resource.")
	a.Describe(&s.RouterID, "ID of the created NetworkRouter resource.")
	a.Describe(&s.SubnetIDs, "IDs of the created NetworkResource (subnet) resources, in declaration order.")
}

// NetworkBundle is the ComponentResource anchor for the NetworkBundle component.
type NetworkBundle struct{}

// Construct implements infer.ComponentResource and creates the child resources.
func (*NetworkBundle) Construct(
	ctx *pulumi.Context, name, typ string,
	args NetworkBundleArgs, opts pulumi.ResourceOption,
) (*NetworkBundleState, error) {
	return newNetworkBundle(ctx, name, typ, args, opts)
}

func newNetworkBundle(
	ctx *pulumi.Context,
	name, typ string,
	args NetworkBundleArgs,
	opts ...pulumi.ResourceOption,
) (*NetworkBundleState, error) {
	comp := &NetworkBundleState{} //nolint:exhaustruct

	err := ctx.RegisterComponentResource(typ, name, comp, opts...)
	if err != nil {
		return nil, fmt.Errorf("registering NetworkBundle component: %w", err)
	}

	networkInputs := pulumi.Map{
		"name": pulumi.String(args.Name),
	}
	if args.Description != nil {
		networkInputs["description"] = pulumi.String(*args.Description)
	}

	var net pulumi.CustomResourceState

	err = ctx.RegisterResource(tokenNetwork, name+"-network", networkInputs, &net, pulumi.Parent(comp))
	if err != nil {
		return nil, fmt.Errorf("creating Network: %w", err)
	}

	routerInputs := pulumi.Map{
		"networkID":  net.ID().ToStringOutput(),
		"enabled":    pulumi.Bool(args.Router.Enabled),
		"masquerade": pulumi.Bool(args.Router.Masquerade),
		"metric":     pulumi.Int(args.Router.Metric),
	}

	if args.Router.PeerGroups != nil {
		peerGroups := make(pulumi.StringArray, len(*args.Router.PeerGroups))
		for j, pg := range *args.Router.PeerGroups {
			peerGroups[j] = pulumi.String(pg)
		}

		routerInputs["peerGroups"] = peerGroups
	}

	if args.Router.Peer != nil {
		routerInputs["peer"] = pulumi.String(*args.Router.Peer)
	}

	var router pulumi.CustomResourceState

	err = ctx.RegisterResource(tokenNetworkRouter, name+"-router", routerInputs, &router, pulumi.Parent(comp))
	if err != nil {
		return nil, fmt.Errorf("creating NetworkRouter: %w", err)
	}

	subnetIDs := make(pulumi.StringArray, len(args.Subnets))

	for subnetIdx, subnet := range args.Subnets {
		groupIDs := make(pulumi.StringArray, len(subnet.GroupIDs))
		for j, gid := range subnet.GroupIDs {
			groupIDs[j] = pulumi.String(gid)
		}

		subnetInputs := pulumi.Map{
			"name":      pulumi.String(subnet.Name),
			"networkID": net.ID().ToStringOutput(),
			"address":   pulumi.String(subnet.Address),
			"enabled":   pulumi.Bool(subnet.Enabled),
			"groupIDs":  groupIDs,
		}

		if subnet.Description != nil {
			subnetInputs["description"] = pulumi.String(*subnet.Description)
		}

		var sub pulumi.CustomResourceState

		err = ctx.RegisterResource(
			tokenNetworkResource,
			name+"-subnet-"+subnet.Name,
			subnetInputs,
			&sub,
			pulumi.Parent(comp),
		)
		if err != nil {
			return nil, fmt.Errorf("creating NetworkResource %q: %w", subnet.Name, err)
		}

		subnetIDs[subnetIdx] = sub.ID().ToStringOutput()
	}

	comp.NetworkID = net.ID().ToStringOutput()
	comp.RouterID = router.ID().ToStringOutput()
	comp.SubnetIDs = subnetIDs.ToStringArrayOutput()

	return comp, nil
}
